const { app, BrowserWindow, ipcMain, dialog } = require("electron");
const pdfParse = require("pdf-parse");
const mammoth = require("mammoth");
const fs = require("fs");
const path = require("path");
const os = require("os");
const { queryOpenAI } = require("./chatgpt.js");

let win;

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

ipcMain.on("start-query-openai", (event, system, user) => {
  if (system == null || user == null) {
    console.error("Received null for system or user message");
    event.reply("openai-response", "Error: System or user message is null.");
    return;
  }

  queryOpenAI(system, user, (message) => {
    event.reply("openai-response", message);
  });
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

function checkApiKeyExists() {
  const configPath = path.join(os.homedir(), ".config", "fabric", ".env");
  return fs.existsSync(configPath);
}

function getPatternFolders() {
  const patternsPath = path.join(__dirname, "patterns");
  return fs
    .readdirSync(patternsPath, { withFileTypes: true })
    .filter((dirent) => dirent.isDirectory())
    .map((dirent) => dirent.name);
}

function getPatternContent(patternName) {
  const patternPath = path.join(
    __dirname,
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

app.whenReady().then(async () => {
  createWindow();

  // Show dialog if API key does not exist
  if (!checkApiKeyExists()) {
    // Note: Electron does not have a built-in showInputBox method.
    // You would need to implement a custom dialog or use a web-based input for this.
    console.log("API key is missing. Implement dialog to collect API key.");
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
