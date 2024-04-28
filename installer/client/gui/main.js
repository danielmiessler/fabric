const { app, BrowserWindow, ipcMain, dialog } = require("electron");
const fs = require("fs").promises;
const fsp = require("fs");
const path = require("path");
const os = require("os");
const OpenAI = require("openai");
const Ollama = require("ollama");
const Anthropic = require("@anthropic-ai/sdk");
const axios = require("axios");
const fsExtra = require("fs-extra");
const fsConstants = require("fs").constants;

let fetch, allModels;

import("node-fetch").then((module) => {
  fetch = module.default;
});
const unzipper = require("unzipper");

let win;
let openai;
let ollama = new Ollama.Ollama();

async function ensureFabricFoldersExist() {
  const fabricPath = path.join(os.homedir(), ".config", "fabric");
  const patternsPath = path.join(fabricPath, "patterns");

  try {
    await fs
      .access(fabricPath, fsConstants.F_OK)
      .catch(() => fs.mkdir(fabricPath, { recursive: true }));
    await fs
      .access(patternsPath, fsConstants.F_OK)
      .catch(() => fs.mkdir(patternsPath, { recursive: true }));
    // Optionally download and update patterns after ensuring the directories exist
  } catch (error) {
    console.error("Error ensuring fabric folders exist:", error);
    throw error; // Make sure to re-throw the error to handle it further up the call stack if necessary
  }
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
    await fs.writeFile(zipPath, response.data);
    console.log("Zip file written to:", zipPath);

    // Prepare for extraction
    const tempExtractPath = path.join(os.tmpdir(), "fabric_extracted");
    await fsExtra.emptyDir(tempExtractPath);

    // Extract the zip file
    await fsp
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
    if (fsp.existsSync(existingPatternsPath)) {
      const existingFolders = await fsExtra.readdir(existingPatternsPath);
      for (const folder of existingFolders) {
        if (!fsp.existsSync(path.join(extractedPatternsPath, folder))) {
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
  return new Promise((resolve, reject) => {
    fs.readdir(patternsPath, { withFileTypes: true }, (error, dirents) => {
      if (error) {
        console.error("Failed to read pattern folders:", error);
        reject(error);
      } else {
        const folders = dirents
          .filter((dirent) => dirent.isDirectory())
          .map((dirent) => dirent.name);
        resolve(folders);
      }
    });
  });
}

async function checkApiKeyExists() {
  const configPath = path.join(os.homedir(), ".config", "fabric", ".env");
  try {
    await fs.access(configPath, fsConstants.F_OK);
    return true; // The file exists
  } catch (e) {
    return false; // The file does not exist
  }
}

async function loadApiKeys() {
  const configPath = path.join(os.homedir(), ".config", "fabric", ".env");
  let keys = { openAIKey: null, claudeKey: null };

  try {
    const envContents = await fs.readFile(configPath, { encoding: "utf8" });
    const openAIMatch = envContents.match(/^OPENAI_API_KEY=(.*)$/m);
    const claudeMatch = envContents.match(/^CLAUDE_API_KEY=(.*)$/m);

    if (openAIMatch && openAIMatch[1]) {
      keys.openAIKey = openAIMatch[1];
    }
    if (claudeMatch && claudeMatch[1]) {
      keys.claudeKey = claudeMatch[1];
      claude = new Anthropic({ apiKey: keys.claudeKey });
    }
  } catch (error) {
    console.error("Could not load API keys:", error);
  }
  return keys;
}

async function saveApiKeys(openAIKey, claudeKey) {
  const configPath = path.join(os.homedir(), ".config", "fabric");
  const envFilePath = path.join(configPath, ".env");

  try {
    await fs.access(configPath);
  } catch {
    await fs.mkdir(configPath, { recursive: true });
  }

  let envContent = "";

  // Read the existing .env file if it exists
  try {
    envContent = await fs.readFile(envFilePath, "utf8");
  } catch (err) {
    if (err.code !== "ENOENT") {
      throw err;
    }
    // If the file doesn't exist, create an empty .env file
    await fs.writeFile(envFilePath, "");
  }

  // Update the specific API key
  if (openAIKey) {
    envContent = updateOrAddKey(envContent, "OPENAI_API_KEY", openAIKey);
    process.env.OPENAI_API_KEY = openAIKey; // Set for current session
    openai = new OpenAI({ apiKey: openAIKey });
  }
  if (claudeKey) {
    envContent = updateOrAddKey(envContent, "CLAUDE_API_KEY", claudeKey);
    process.env.CLAUDE_API_KEY = claudeKey; // Set for current session
    claude = new Anthropic({ apiKey: claudeKey });
  }

  await fs.writeFile(envFilePath, envContent.trim());
  await loadApiKeys();
  win.webContents.send("api-keys-saved");
}

function updateOrAddKey(envContent, keyName, keyValue) {
  const keyPattern = new RegExp(`^${keyName}=.*$`, "m");
  if (keyPattern.test(envContent)) {
    // Update the existing key
    envContent = envContent.replace(keyPattern, `${keyName}=${keyValue}`);
  } else {
    // Add the new key
    envContent += `\n${keyName}=${keyValue}`;
  }
  return envContent;
}

async function getOllamaModels() {
  try {
    ollama = new Ollama.Ollama();
    const _models = await ollama.list();
    return _models.models.map((x) => x.name);
  } catch (error) {
    if (error.cause && error.cause.code === "ECONNREFUSED") {
      console.error(
        "Failed to connect to Ollama. Make sure Ollama is running and accessible."
      );
      return []; // Return an empty array instead of throwing an error
    } else {
      console.error("Error fetching models from Ollama:", error);
      throw error; // Re-throw the error for other types of errors
    }
  }
}

async function getModels() {
  allModels = {
    gptModels: [],
    claudeModels: [],
    ollamaModels: [],
  };

  let keys = await loadApiKeys();

  if (keys.claudeKey) {
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
    try {
      const response = await openai.models.list();
      allModels.gptModels = response.data;
    } catch (error) {
      console.error("Error fetching models from OpenAI:", error);
    }
  }

  // Check if ollama exists and has a list method
  if (
    typeof ollama !== "undefined" &&
    ollama.list &&
    typeof ollama.list === "function"
  ) {
    try {
      allModels.ollamaModels = await getOllamaModels();
    } catch (error) {
      console.error("Error fetching models from Ollama:", error);
    }
  } else {
    console.log("Ollama is not available or does not support listing models.");
  }

  return allModels;
}

async function getPatternContent(patternName) {
  const patternPath = path.join(
    os.homedir(),
    ".config",
    "fabric",
    "patterns",
    patternName,
    "system.md"
  );
  try {
    const content = await fs.readFile(patternPath, "utf8");
    return content;
  } catch (error) {
    console.error("Error reading pattern file:", error);
    return "";
  }
}

async function ollamaMessage(
  system,
  user,
  model,
  temperature,
  topP,
  frequencyPenalty,
  presencePenalty,
  event
) {
  ollama = new Ollama.Ollama();
  const userMessage = {
    role: "user",
    content: user,
  };
  const systemMessage = { role: "system", content: system };
  const response = await ollama.chat({
    model: model,
    messages: [systemMessage, userMessage],
    temperature: temperature,
    top_p: topP,
    frequency_penalty: frequencyPenalty,
    presence_penalty: presencePenalty,
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

async function openaiMessage(
  system,
  user,
  model,
  temperature,
  topP,
  frequencyPenalty,
  presencePenalty,
  event
) {
  const userMessage = { role: "user", content: user };
  const systemMessage = { role: "system", content: system };
  const stream = await openai.chat.completions.create(
    {
      model: model,
      messages: [systemMessage, userMessage],
      temperature: temperature,
      top_p: topP,
      frequency_penalty: frequencyPenalty,
      presence_penalty: presencePenalty,
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

async function claudeMessage(system, user, model, temperature, topP, event) {
  if (!claude) {
    event.reply(
      "model-response-error",
      "Claude API key is missing or invalid."
    );
    return;
  }

  const userMessage = { role: "user", content: user };
  const systemMessage = system;
  const response = await claude.messages.create({
    model: model,
    system: systemMessage,
    max_tokens: 4096,
    messages: [userMessage],
    stream: true,
    temperature: temperature,
    top_p: topP,
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

async function createPatternFolder(patternName, patternBody) {
  try {
    const patternsPath = path.join(
      os.homedir(),
      ".config",
      "fabric",
      "patterns"
    );
    const patternFolderPath = path.join(patternsPath, patternName);

    // Create the pattern folder using the promise-based API
    await fs.mkdir(patternFolderPath, { recursive: true });

    // Create the system.md file inside the pattern folder
    const filePath = path.join(patternFolderPath, "system.md");
    await fs.writeFile(filePath, patternBody);

    console.log(
      `Pattern folder '${patternName}' created successfully with system.md inside.`
    );
    return `Pattern folder '${patternName}' created successfully with system.md inside.`;
  } catch (err) {
    console.error(`Failed to create the pattern folder: ${err.message}`);
    throw err; // Ensure the error is thrown so it can be caught by the caller
  }
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

ipcMain.on(
  "start-query",
  async (
    event,
    system,
    user,
    model,
    temperature,
    topP,
    frequencyPenalty,
    presencePenalty
  ) => {
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
        await claudeMessage(system, user, model, temperature, topP, event);
      } else if (_gptModels.includes(model)) {
        await openaiMessage(
          system,
          user,
          model,
          temperature,
          topP,
          frequencyPenalty,
          presencePenalty,
          event
        );
      } else if (allModels.ollamaModels.includes(model)) {
        await ollamaMessage(
          system,
          user,
          model,
          temperature,
          topP,
          frequencyPenalty,
          presencePenalty,
          event
        );
      } else {
        event.reply("model-response-error", "Unsupported model: " + model);
      }
    } catch (error) {
      console.error("Error querying model:", error);
      event.reply("model-response-error", "Error querying model.");
    }
  }
);

ipcMain.handle("create-pattern", async (event, patternName, patternContent) => {
  try {
    const result = await createPatternFolder(patternName, patternContent);
    return { status: "success", message: result }; // Use a response object for more detailed responses
  } catch (error) {
    console.error("Error creating pattern:", error);
    return { status: "error", message: error.message }; // Return an error object
  }
});

// Example of using ipcMain.handle for asynchronous operations
ipcMain.handle("get-patterns", async (event) => {
  try {
    const patterns = await getPatternFolders();
    return patterns;
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
    const content = await getPatternContent(patternName);
    return content;
  } catch (error) {
    console.error("Failed to get pattern content:", error);
    return "";
  }
});

ipcMain.handle("save-api-keys", async (event, { openAIKey, claudeKey }) => {
  try {
    await saveApiKeys(openAIKey, claudeKey);
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
    const keys = await loadApiKeys();
    await ensureFabricFoldersExist(); // Ensure fabric folders exist
    await getModels(); // Fetch models after loading API keys
    createWindow(); // Keep this line
  } catch (error) {
    await ensureFabricFoldersExist(); // Ensure fabric folders exist
    createWindow(); // Keep this line
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
