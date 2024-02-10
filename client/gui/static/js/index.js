document.addEventListener("DOMContentLoaded", async function () {
  const patternSelector = document.getElementById("patternSelector");
  const userInput = document.getElementById("userInput");
  const submitButton = document.getElementById("submit");
  const responseContainer = document.getElementById("responseContainer");
  const themeChanger = document.getElementById("themeChanger");
  const configButton = document.getElementById("configButton");
  const configSection = document.getElementById("configSection");
  const saveApiKeyButton = document.getElementById("saveApiKey");
  const apiKeyInput = document.getElementById("apiKeyInput");
  const originalPlaceholder = userInput.placeholder;
  const copyButton = document.getElementById("copyButton");

  async function submitQuery(userInputValue) {
    userInput.value = ""; // Clear the input after submitting
    systemCommand = await window.electronAPI.invoke(
      "get-pattern-content",
      patternSelector.value
    );
    responseContainer.innerHTML = ""; // Clear previous responses
    responseContainer.classList.remove("hidden");
    window.electronAPI.send(
      "start-query-openai",
      systemCommand,
      userInputValue
    );
  }

  function fallbackCopyTextToClipboard(text) {
    const textArea = document.createElement("textarea");
    textArea.value = text;

    // Avoid scrolling to bottom
    textArea.style.top = "0";
    textArea.style.left = "0";
    textArea.style.position = "fixed";

    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    function copyToClipboard() {
      try {
        if (responseContainer.textContent) {
          text = responseContainer.textContent;
        }
        if (navigator.clipboard) {
          navigator.clipboard
            .writeText(text)
            .then(function () {
              console.log("Text successfully copied to clipboard");
            })
            .catch(function (err) {
              console.error("Error in copying text: ", err);
              // Optionally, use fallback method here
            });
        } else {
          fallbackCopyTextToClipboard(text);
        }
      } catch (err) {
        console.error("Error in copying text: ", err);
      }
    }

    try {
      const successful = document.execCommand("copy");
      const msg = successful ? "successful" : "unsuccessful";
      console.log("Fallback: Copying text command was " + msg);
    } catch (err) {
      console.error("Fallback: Oops, unable to copy", err);
    }

    document.body.removeChild(textArea);
  }

  // Load patterns on startup
  try {
    const patterns = await window.electronAPI.invoke("get-patterns");
    patterns.forEach((pattern) => {
      const option = document.createElement("option");
      option.value = pattern;
      option.textContent = pattern;
      patternSelector.appendChild(option);
    });
  } catch (error) {
    console.error("Failed to load patterns:", error);
  }

  // Listen for OpenAI responses
  window.electronAPI.on("openai-response", (message) => {
    const formattedMessage = message.replace(/\n/g, "<br>");
    responseContainer.innerHTML += formattedMessage; // Append new data as it arrives
  });

  window.electronAPI.on("file-response", (message) => {
    if (message.startsWith("Error")) {
      alert(message);
      return;
    }
    submitQuery(message);
  });

  // Submit button click handler
  submitButton.addEventListener("click", async () => {
    const userInputValue = userInput.value;
    submitQuery(userInputValue);
  });

  // Theme changer click handler
  themeChanger.addEventListener("click", function (e) {
    e.preventDefault();
    document.body.classList.toggle("light-theme");
    themeChanger.innerText =
      themeChanger.innerText === "Dark" ? "Light" : "Dark";
  });

  // Config button click handler - toggles the config section visibility
  configButton.addEventListener("click", function (e) {
    e.preventDefault();
    configSection.classList.toggle("hidden");
  });

  // Save API Key button click handler
  saveApiKeyButton.addEventListener("click", () => {
    const apiKey = apiKeyInput.value;
    window.electronAPI
      .invoke("save-api-key", apiKey)
      .then(() => {
        alert("API Key saved successfully.");
        // Optionally hide the config section and clear the input after saving
        configSection.classList.add("hidden");
        apiKeyInput.value = "";
      })
      .catch((err) => {
        console.error("Error saving API key:", err);
        alert("Failed to save API Key.");
      });
  });

  // Handler for pattern selection change
  patternSelector.addEventListener("change", async () => {
    const selectedPattern = patternSelector.value;
    const systemCommand = await window.electronAPI.invoke(
      "get-pattern-content",
      selectedPattern
    );
    // Use systemCommand as part of the input for querying OpenAI
  });

  userInput.addEventListener("dragover", (event) => {
    event.stopPropagation();
    event.preventDefault();
    // Add some visual feedback
    userInput.classList.add("drag-over");
    userInput.placeholder = "Drop file here";
  });

  userInput.addEventListener("dragleave", (event) => {
    event.stopPropagation();
    event.preventDefault();
    // Remove visual feedback
    userInput.classList.remove("drag-over");
    userInput.placeholder = originalPlaceholder;
  });

  userInput.addEventListener("drop", (event) => {
    event.stopPropagation();
    event.preventDefault();
    const file = event.dataTransfer.files[0];
    userInput.classList.remove("drag-over");
    userInput.placeholder = originalPlaceholder;
    processFile(file);
  });

  function processFile(file) {
    const fileType = file.type;
    const reader = new FileReader();
    let content = "";

    reader.onload = (event) => {
      content = event.target.result;
      userInput.value = content;
      submitQuery(content);
    };

    if (fileType === "text/plain" || fileType === "image/svg+xml") {
      reader.readAsText(file);
    } else if (
      fileType === "application/pdf" ||
      fileType.match(/wordprocessingml/)
    ) {
      // For PDF and DOCX, we need to handle them in the main process due to complexity
      window.electronAPI.send("process-complex-file", file.path);
    } else {
      console.error("Unsupported file type");
    }
  }
});
