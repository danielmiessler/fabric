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
  const updatePatternsButton = document.getElementById("updatePatternsButton");
  const copyButton = document.createElement("button");

  window.electronAPI.on("patterns-ready", () => {
    console.log("Patterns are ready. Refreshing the pattern list.");
    loadPatterns();
  });
  window.electronAPI.on("request-api-key", () => {
    // Show the API key input section or modal to the user
    configSection.classList.remove("hidden"); // Assuming 'configSection' is your API key input area
  });
  copyButton.textContent = "Copy";
  copyButton.id = "copyButton";
  document.addEventListener("click", function (e) {
    if (e.target && e.target.id === "copyButton") {
      // Your copy to clipboard function
      copyToClipboard();
    }
  });
  window.electronAPI.on("no-api-key", () => {
    alert("API key is missing. Please enter your OpenAI API key.");
  });

  window.electronAPI.on("patterns-updated", () => {
    alert("Patterns updated. Refreshing the pattern list.");
    loadPatterns();
  });

  function htmlToPlainText(html) {
    // Create a temporary div element to hold the HTML
    var tempDiv = document.createElement("div");
    tempDiv.innerHTML = html;

    // Replace <br> tags with newline characters
    tempDiv.querySelectorAll("br").forEach((br) => br.replaceWith("\n"));

    // Replace block elements like <p> and <div> with newline characters
    tempDiv.querySelectorAll("p, div").forEach((block) => {
      block.prepend("\n"); // Add a newline before the block element's content
      block.replaceWith(...block.childNodes); // Replace the block element with its own contents
    });

    // Return the text content, trimming leading and trailing newlines
    return tempDiv.textContent.trim();
  }

  async function submitQuery(userInputValue) {
    userInput.value = ""; // Clear the input after submitting
    systemCommand = await window.electronAPI.invoke(
      "get-pattern-content",
      patternSelector.value
    );
    responseContainer.innerHTML = ""; // Clear previous responses
    if (responseContainer.classList.contains("hidden")) {
      console.log("contains hidden");
      responseContainer.classList.remove("hidden");
      responseContainer.appendChild(copyButton);
    }
    window.electronAPI.send(
      "start-query-openai",
      systemCommand,
      userInputValue
    );
  }

  function copyToClipboard() {
    const containerClone = responseContainer.cloneNode(true);
    // Remove the copy button from the clone
    const copyButtonClone = containerClone.querySelector("#copyButton");
    if (copyButtonClone) {
      copyButtonClone.parentNode.removeChild(copyButtonClone);
    }

    // Convert HTML to plain text, preserving newlines
    const plainText = htmlToPlainText(containerClone.innerHTML);

    // Use a temporary textarea for copying
    const textArea = document.createElement("textarea");
    textArea.style.position = "absolute";
    textArea.style.left = "-9999px";
    textArea.setAttribute("aria-hidden", "true");
    textArea.value = plainText;
    document.body.appendChild(textArea);
    textArea.select();

    try {
      document.execCommand("copy");
      console.log("Text successfully copied to clipboard");
    } catch (err) {
      console.error("Failed to copy text: ", err);
    }

    document.body.removeChild(textArea);
  }
  async function loadPatterns() {
    try {
      const patterns = await window.electronAPI.invoke("get-patterns");
      patternSelector.innerHTML = ""; // Clear existing options first
      patterns.forEach((pattern) => {
        const option = document.createElement("option");
        option.value = pattern;
        option.textContent = pattern;
        patternSelector.appendChild(option);
      });
    } catch (error) {
      console.error("Failed to load patterns:", error);
    }
  }

  function fallbackCopyTextToClipboard(text) {
    const textArea = document.createElement("textarea");
    textArea.value = text;
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    try {
      const successful = document.execCommand("copy");
      const msg = successful ? "successful" : "unsuccessful";
      console.log("Fallback: Copying text command was " + msg);
    } catch (err) {
      console.error("Fallback: Oops, unable to copy", err);
    }

    document.body.removeChild(textArea);
  }

  updatePatternsButton.addEventListener("click", () => {
    window.electronAPI.send("update-patterns");
  });

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

  // drag and drop
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
