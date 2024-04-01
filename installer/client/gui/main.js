const { app, BrowserWindow, ipcMain, dialog } = require("electron");
const fs = require("fs");
const path = require("path");
const os = require("os");
const OpenAI = require("openai");
const Ollama = require("ollama");
const Anthropic = require("@anthropic-ai/sdk");

let fetch, allModels;

import("node-fetch").then((module) => {
  fetch = module.default;
});
const unzipper = require("unzipper");

let win;
let openai;
let ollama;

function ensureFabricFoldersExist() {
  return new Promise(async (resolve, reject) => {
    const fabricPath = path.join(os.homedir(), ".config", "fabric");
    const patternsPath = path.join(fabricPath, "patterns");

    try {
      if (!fs.existsSync(fabricPath)) {
        fs.mkdirSync(fabricPath, { recursive: true });
      }

      if (!fs.existsSync(patternsPath)) {
        fs.mkdirSync(patternsPath, { recursive: true });
        await downloadAndUpdatePatterns(patternsPath);
      }
      resolve(); // Resolve the promise once everything is set up
    } catch (error) {
      console.error("Error ensuring fabric folders exist:", error);
      reject(error); // Reject the promise if an error occurs
    }
  });
}

async function downloadAndUpdatePatterns() {
  try {
    // Download the zip file
    const response = await axios({
      method: "get",
      url: "https://github.com/danielmiessler/fabric/archive/refs/heads/main.zip",
      responseType: "arraybuffer",
    });

    const zipPath = path.join(os.tmpdir(), "fabric.zip");
    fs.writeFileSync(zipPath, response.data);
    console.log("Zip file written to:", zipPath);

    // Prepare for extraction
    const tempExtractPath = path.join(os.tmpdir(), "fabric_extracted");
    await fsExtra.emptyDir(tempExtractPath);

    // Extract the zip file
    await fs
      .createReadStream(zipPath)
      .pipe(unzipper.Extract({ path: tempExtractPath }))
      .promise();
    console.log("Extraction complete");

    const extractedPatternsPath = path.join(
      tempExtractPath,
      "fabric-main",
      "patterns"
    );

    // Compare and move folders
    const existingPatternsPath = path.join(
      os.homedir(),
      ".config",
      "fabric",
      "patterns"
    );
    if (fs.existsSync(existingPatternsPath)) {
      const existingFolders = await fsExtra.readdir(existingPatternsPath);
      for (const folder of existingFolders) {
        if (!fs.existsSync(path.join(extractedPatternsPath, folder))) {
          await fsExtra.move(
            path.join(existingPatternsPath, folder),
            path.join(extractedPatternsPath, folder)
          );
          console.log(
            `Moved missing folder ${folder} to the extracted patterns directory.`
          );
        }
      }
    }

    // Overwrite the existing patterns directory with the updated extracted directory
    await fsExtra.copy(extractedPatternsPath, existingPatternsPath, {
      overwrite: true,
    });
    console.log("Patterns successfully updated");

    // Inform the renderer process that the patterns have been updated
    // win.webContents.send("patterns-updated");
  } catch (error) {
    console.error("Error downloading or updating patterns:", error);
  }
}
function getPatternFolders() {
  const patternsPath = path.join(os.homedir(), ".config", "fabric", "patterns");
  return fs
    .readdirSync(patternsPath, { withFileTypes: true })
    .filter((dirent) => dirent.isDirectory())
    .map((dirent) => dirent.name);
}

function checkApiKeyExists() {
  const configPath = path.join(os.homedir(), ".config", "fabric", ".env");
  return fs.existsSync(configPath);
}

function loadApiKeys() {
  const configPath = path.join(os.homedir(), ".config", "fabric", ".env");
  let keys = { openAIKey: null, claudeKey: null };

  if (fs.existsSync(configPath)) {
    const envContents = fs.readFileSync(configPath, { encoding: "utf8" });
    const openAIMatch = envContents.match(/^OPENAI_API_KEY=(.*)$/m);
    const claudeMatch = envContents.match(/^CLAUDE_API_KEY=(.*)$/m);

    if (openAIMatch && openAIMatch[1]) {
      keys.openAIKey = openAIMatch[1];
      openai = new OpenAI({ apiKey: keys.openAIKey });
    }
    if (claudeMatch && claudeMatch[1]) {
      keys.claudeKey = claudeMatch[1];
      claude = new Anthropic({ apiKey: keys.claudeKey });
    }
  }
  return keys;
}

function saveApiKeys(openAIKey, claudeKey) {
  const configPath = path.join(os.homedir(), ".config", "fabric");
  const envFilePath = path.join(configPath, ".env");

  if (!fs.existsSync(configPath)) {
    fs.mkdirSync(configPath, { recursive: true });
  }

  let envContent = "";
  if (openAIKey) {
    envContent += `OPENAI_API_KEY=${openAIKey}\n`;
    process.env.OPENAI_API_KEY = openAIKey; // Set for current session
    openai = new OpenAI({ apiKey: openAIKey });
  }
  if (claudeKey) {
    envContent += `CLAUDE_API_KEY=${claudeKey}\n`;
    process.env.CLAUDE_API_KEY = claudeKey; // Set for current session
    claude = new Anthropic({ apiKey: claudeKey });
  }

  fs.writeFileSync(envFilePath, envContent.trim());
}

async function getOllamaModels() {
  ollama = new Ollama.Ollama();
  const _models = await ollama.list();
  return _models.models.map((x) => x.name);
}

async function getModels() {
  ollama = new Ollama.Ollama();
  allModels = {
    gptModels: [],
    claudeModels: [],
    ollamaModels: [],
  };

  let keys = loadApiKeys(); // Assuming loadApiKeys() is updated to return both keys

  if (keys.claudeKey) {
    // Assuming claudeModels do not require an asynchronous call to be fetched
    claudeModels = [
      "claude-3-opus-20240229",
      "claude-3-sonnet-20240229",
      "claude-3-haiku-20240307",
      "claude-2.1",
    ];
    allModels.claudeModels = claudeModels;
  }

  if (keys.openAIKey) {
    openai = new OpenAI({ apiKey: keys.openAIKey });
    // Wrap asynchronous call with a Promise to handle it in parallel
    gptModelsPromise = openai.models.list();
  }

  // Check if ollama exists and has a list method
  if (
    typeof ollama !== "undefined" &&
    ollama.list &&
    typeof ollama.list === "function"
  ) {
    // Assuming ollama.list() returns a Promise
    ollamaModelsPromise = getOllamaModels();
  } else {
    console.log("ollama is not available or does not support listing models.");
  }

  // Wait for all asynchronous operations to complete
  try {
    const results = await Promise.all(
      [gptModelsPromise, ollamaModelsPromise].filter(Boolean)
    ); // Filter out undefined promises
    allModels.gptModels = results[0]?.data || []; // Assuming the first promise is always GPT models if it exists
    allModels.ollamaModels = results[1] || []; // Assuming the second promise is always Ollama models if it exists
  } catch (error) {
    console.error("Error fetching models from OpenAI or Ollama:", error);
  }

  return allModels; // Return the aggregated results
}

function getPatternContent(patternName) {
  const patternPath = path.join(
    os.homedir(),
    ".config",
    "fabric",
    "patterns",
    patternName,
    "system.md"
  );
  try {
    return fs.readFileSync(patternPath, "utf8");
  } catch (error) {
    console.error("Error reading pattern file:", error);
    return "";
  }
}

async function ollamaMessage(system, user, model, event) {
  ollama = new Ollama.Ollama();
  const userMessage = {
    role: "user",
    content: user,
  };
  const systemMessage = { role: "system", content: system };
  const response = await ollama.chat({
    model: model,
    messages: [systemMessage, userMessage],
    stream: true,
  });
  let responseMessage = "";
  for await (const chunk of response) {
    const content = chunk.message.content;
    if (content) {
      responseMessage += content;
      event.reply("model-response", content);
    }
    event.reply("model-response-end", responseMessage);
  }
}

async function openaiMessage(system, user, model, event) {
  const userMessage = { role: "user", content: user };
  const systemMessage = { role: "system", content: system };
  const stream = await openai.chat.completions.create(
    {
      model: model,
      messages: [systemMessage, userMessage],
      stream: true,
    },
    { responseType: "stream" }
  );

  let responseMessage = "";

  for await (const chunk of stream) {
    const content = chunk.choices[0].delta.content;
    if (content) {
      responseMessage += content;
      event.reply("model-response", content);
    }
  }

  event.reply("model-response-end", responseMessage);
}

async function claudeMessage(system, user, model, event) {
  const userMessage = { role: "user", content: user };
  const systemMessage = system;
  const response = await claude.messages.create({
    model: model,
    system: systemMessage,
    max_tokens: 4096,
    messages: [userMessage],
    stream: true,
    temperature: 0.0,
    top_p: 1.0,
  });
  let responseMessage = "";
  for await (const chunk of response) {
    if (chunk.delta && chunk.delta.text) {
      responseMessage += chunk.delta.text;
      event.reply("model-response", chunk.delta.text);
    }
  }
  event.reply("model-response-end", responseMessage);
}

function createWindow() {
  win = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      contextIsolation: true,
      nodeIntegration: false,
      preload: path.join(__dirname, "preload.js"),
    },
  });

  win.loadFile("index.html");

  win.on("closed", () => {
    win = null;
  });
}

ipcMain.on("start-query", async (event, system, user, model) => {
  if (system == null || user == null || model == null) {
    console.error("Received null for system, user message, or model");
    event.reply(
      "model-response-error",
      "Error: System, user message, or model is null."
    );
    return;
  }

  try {
    const _gptModels = allModels.gptModels.map((model) => model.id);
    if (allModels.claudeModels.includes(model)) {
      await claudeMessage(system, user, model, event);
    } else if (_gptModels.includes(model)) {
      await openaiMessage(system, user, model, event);
    } else if (allModels.ollamaModels.includes(model)) {
      await ollamaMessage(system, user, model, event);
    } else {
      event.reply("model-response-error", "Unsupported model: " + model);
    }
  } catch (error) {
    console.error("Error querying model:", error);
    event.reply("model-response-error", "Error querying model.");
  }
});

// Example of using ipcMain.handle for asynchronous operations
ipcMain.handle("get-patterns", async (event) => {
  try {
    return getPatternFolders();
  } catch (error) {
    console.error("Failed to get patterns:", error);
    return [];
  }
});

ipcMain.on("update-patterns", () => {
  const patternsPath = path.join(os.homedir(), ".config", "fabric", "patterns");
  downloadAndUpdatePatterns(patternsPath);
});

ipcMain.handle("get-pattern-content", async (event, patternName) => {
  try {
    return getPatternContent(patternName);
  } catch (error) {
    console.error("Failed to get pattern content:", error);
    return "";
  }
});

ipcMain.handle("save-api-keys", async (event, { openAIKey, claudeKey }) => {
  try {
    saveApiKeys(openAIKey, claudeKey);
    return "API Keys saved successfully.";
  } catch (error) {
    console.error("Error saving API keys:", error);
    throw new Error("Failed to save API Keys.");
  }
});

ipcMain.handle("get-models", async (event) => {
  try {
    const models = await getModels();
    return models;
  } catch (error) {
    console.error("Failed to get models:", error);
    return { gptModels: [], claudeModels: [], ollamaModels: [] };
  }
});

app.whenReady().then(async () => {
  try {
    const keys = loadApiKeys();
    if (!keys.openAIKey && !keys.claudeKey) {
      promptUserForApiKey();
    } else {
      createWindow();
    }
    await ensureFabricFoldersExist(); // Ensure fabric folders exist
    createWindow(); // Create the application window

    // After window creation, check if the API key exists
    if (!checkApiKeyExists()) {
      console.log("API key is missing. Prompting user to input API key.");
      win.webContents.send("request-api-key");
    }
  } catch (error) {
    console.error("Failed to initialize fabric folders:", error);
    // Handle initialization failure (e.g., close the app or show an error message)
  }
});

app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
  }
});

app.on("activate", () => {
  if (win === null) {
    createWindow();
  }
});
