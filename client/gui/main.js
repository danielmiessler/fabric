const { app, BrowserWindow, ipcMain, dialog } = require("electron");
const pdfParse = require("pdf-parse");
const mammoth = require("mammoth");
const fs = require("fs");
const path = require("path");
const os = require("os");
const { queryOpenAI } = require("./chatgpt.js");
const axios = require("axios");
const fsExtra = require("fs-extra");

let fetch;
import("node-fetch").then((module) => {
  fetch = module.default;
});
const unzipper = require("unzipper");

let win;

function promptUserForApiKey() {
  // Create a new window to prompt the user for the API key
  const promptWindow = new BrowserWindow({
    // Window configuration for the prompt
    width: 500,
    height: 200,
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false, // Consider security implications
    },
  });

  // Handle the API key submission from the prompt window
  ipcMain.on("submit-api-key", (event, apiKey) => {
    if (apiKey) {
      saveApiKey(apiKey);
      promptWindow.close();
      createWindow(); // Proceed to create the main window
    } else {
      // Handle invalid input or user cancellation
      promptWindow.close();
    }
  });
}

function loadApiKey() {
  const configPath = path.join(os.homedir(), ".config", "fabric", ".env");
  if (fs.existsSync(configPath)) {
    const envContents = fs.readFileSync(configPath, { encoding: "utf8" });
    const matches = envContents.match(/^OPENAI_API_KEY=(.*)$/m);
    if (matches && matches[1]) {
      return matches[1];
    }
  }
  return null;
}

function saveApiKey(apiKey) {
  const configPath = path.join(os.homedir(), ".config", "fabric");
  const envFilePath = path.join(configPath, ".env");

  if (!fs.existsSync(configPath)) {
    fs.mkdirSync(configPath, { recursive: true });
  }

  fs.writeFileSync(envFilePath, `OPENAI_API_KEY=${apiKey}`);
  process.env.OPENAI_API_KEY = apiKey; // Set for current session
}

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

async function downloadAndUpdatePatterns(patternsPath) {
  try {
    const response = await axios({
      method: "get",
      url: "https://github.com/danielmiessler/fabric/archive/refs/heads/main.zip",
      responseType: "arraybuffer",
    });

    const zipPath = path.join(os.tmpdir(), "fabric.zip");
    fs.writeFileSync(zipPath, response.data);
    console.log("Zip file written to:", zipPath);

    const tempExtractPath = path.join(os.tmpdir(), "fabric_extracted");
    fsExtra.emptyDirSync(tempExtractPath);

    await fsExtra.remove(patternsPath); // Delete the existing patterns directory

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

    await fsExtra.copy(extractedPatternsPath, patternsPath);
    console.log("Patterns successfully updated");

    // Inform the renderer process that the patterns have been updated
    win.webContents.send("patterns-updated");
  } catch (error) {
    console.error("Error downloading or updating patterns:", error);
  }
}

function checkApiKeyExists() {
  const configPath = path.join(os.homedir(), ".config", "fabric", ".env");
  return fs.existsSync(configPath);
}

function getPatternFolders() {
  const patternsPath = path.join(os.homedir(), ".config", "fabric", "patterns");
  return fs
    .readdirSync(patternsPath, { withFileTypes: true })
    .filter((dirent) => dirent.isDirectory())
    .map((dirent) => dirent.name);
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
ipcMain.on("process-complex-file", (event, filePath) => {
  const extension = path.extname(filePath).toLowerCase();
  let fileProcessPromise;

  if (extension === ".pdf") {
    const dataBuffer = fs.readFileSync(filePath);
    fileProcessPromise = pdfParse(dataBuffer).then((data) => data.text);
  } else if (extension === ".docx") {
    fileProcessPromise = mammoth
      .extractRawText({ path: filePath })
      .then((result) => result.value)
      .catch((err) => {
        console.error("Error processing DOCX file:", err);
        throw new Error("Error processing DOCX file.");
      });
  } else {
    event.reply("file-response", "Error: Unsupported file type");
    return;
  }

  fileProcessPromise
    .then((extractedText) => {
      // Sending the extracted text back to the frontend.
      event.reply("file-response", extractedText);
    })
    .catch((error) => {
      // Handling any errors during file processing and sending them back to the frontend.
      event.reply("file-response", `Error processing file: ${error.message}`);
    });
});

ipcMain.on("start-query-openai", async (event, system, user) => {
  if (system == null || user == null) {
    console.error("Received null for system or user message");
    event.reply("openai-response", "Error: System or user message is null.");
    return;
  }
  try {
    await queryOpenAI(system, user, (message) => {
      event.reply("openai-response", message);
    });
  } catch (error) {
    console.error("Error querying OpenAI:", error);
    event.reply("no-api-key", "Error querying OpenAI.");
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

ipcMain.handle("save-api-key", async (event, apiKey) => {
  try {
    const configPath = path.join(os.homedir(), ".config", "fabric");
    if (!fs.existsSync(configPath)) {
      fs.mkdirSync(configPath, { recursive: true });
    }

    const envFilePath = path.join(configPath, ".env");
    fs.writeFileSync(envFilePath, `OPENAI_API_KEY=${apiKey}`);
    process.env.OPENAI_API_KEY = apiKey;

    return "API Key saved successfully.";
  } catch (error) {
    console.error("Error saving API key:", error);
    throw new Error("Failed to save API Key.");
  }
});

app.whenReady().then(async () => {
  try {
    const apiKey = loadApiKey();
    if (!apiKey) {
      promptUserForApiKey();
    } else {
      process.env.OPENAI_API_KEY = apiKey;
      createWindow();
    }
    await ensureFabricFoldersExist(); // Ensure fabric folders exist
    createWindow(); // Create the application window

    // After window creation, check if the API key exists
    if (!checkApiKeyExists()) {
      console.log("API key is missing. Prompting user to input API key.");
      // Optionally, directly invoke a function here to show a prompt in the renderer process
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
