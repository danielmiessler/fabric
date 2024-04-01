document.addEventListener("DOMContentLoaded", async function () {
  const patternSelector = document.getElementById("patternSelector");
  const modelSelector = document.getElementById("modelSelector");
  const userInput = document.getElementById("userInput");
  const submitButton = document.getElementById("submit");
  const responseContainer = document.getElementById("responseContainer");
  const themeChanger = document.getElementById("themeChanger");
  const configButton = document.getElementById("configButton");
  const configSection = document.getElementById("configSection");
  const saveApiKeyButton = document.getElementById("saveApiKey");
  const openaiApiKeyInput = document.getElementById("apiKeyInput");
  const claudeApiKeyInput = document.getElementById("claudeApiKeyInput");
  const updatePatternsButton = document.getElementById("updatePatternsButton");
  const updatePatternButton = document.getElementById("createPattern");
  const patternCreator = document.getElementById("patternCreator");
  const submitPatternButton = document.getElementById("submitPattern");
  const fineTuningButton = document.getElementById("fineTuningButton");
  const fineTuningSection = document.getElementById("fineTuningSection");
  const temperatureSlider = document.getElementById("temperatureSlider");
  const temperatureValue = document.getElementById("temperatureValue");
  const topPSlider = document.getElementById("topPSlider");
  const topPValue = document.getElementById("topPValue");
  const frequencyPenaltySlider = document.getElementById(
    "frequencyPenaltySlider"
  );
  const frequencyPenaltyValue = document.getElementById(
    "frequencyPenaltyValue"
  );
  const presencePenaltySlider = document.getElementById(
    "presencePenaltySlider"
  );
  const presencePenaltyValue = document.getElementById("presencePenaltyValue");
  const myForm = document.getElementById("my-form");
  const copyButton = document.createElement("button");

  window.electronAPI.on("patterns-ready", () => {
    console.log("Patterns are ready. Refreshing the pattern list.");
    loadPatterns();
  });
  window.electronAPI.on("request-api-key", () => {
    configSection.classList.remove("hidden");
  });
  copyButton.textContent = "Copy";
  copyButton.id = "copyButton";
  document.addEventListener("click", function (e) {
    if (e.target && e.target.id === "copyButton") {
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
    var tempDiv = document.createElement("div");
    tempDiv.innerHTML = html;

    tempDiv.querySelectorAll("br").forEach((br) => br.replaceWith("\n"));

    tempDiv.querySelectorAll("p, div").forEach((block) => {
      block.prepend("\n");
      block.replaceWith(...block.childNodes);
    });

    return tempDiv.textContent.trim();
  }

  async function submitQuery(userInputValue) {
    const temperature = parseFloat(temperatureSlider.value);
    const topP = parseFloat(topPSlider.value);
    const frequencyPenalty = parseFloat(frequencyPenaltySlider.value);
    const presencePenalty = parseFloat(presencePenaltySlider.value);
    userInput.value = ""; // Clear the input after submitting
    const systemCommand = await window.electronAPI.invoke(
      "get-pattern-content",
      patternSelector.value
    );
    const selectedModel = modelSelector.value;
    responseContainer.innerHTML = ""; // Clear previous responses
    if (responseContainer.classList.contains("hidden")) {
      responseContainer.classList.remove("hidden");
      responseContainer.appendChild(copyButton);
    }
    window.electronAPI.send(
      "start-query",
      systemCommand,
      userInputValue,
      selectedModel,
      temperature,
      topP,
      frequencyPenalty,
      presencePenalty
    );
  }

  async function submitPattern(patternName, patternText) {
    try {
      const response = await window.electronAPI.invoke(
        "create-pattern",
        patternName,
        patternText
      );
      if (response.status === "success") {
        console.log(response.message);
        // Show success message
        const patternCreatedMessage = document.getElementById(
          "patternCreatedMessage"
        );
        patternCreatedMessage.classList.remove("hidden");
        setTimeout(() => {
          patternCreatedMessage.classList.add("hidden");
        }, 3000); // Hide the message after 3 seconds

        // Update pattern list
        loadPatterns();
      } else {
        console.error(response.message);
        // Handle failure (e.g., showing an error message to the user)
      }
    } catch (error) {
      console.error("IPC error:", error);
    }
  }

  function copyToClipboard() {
    const containerClone = responseContainer.cloneNode(true);
    const copyButtonClone = containerClone.querySelector("#copyButton");
    if (copyButtonClone) {
      copyButtonClone.parentNode.removeChild(copyButtonClone);
    }

    const plainText = htmlToPlainText(containerClone.innerHTML);

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

  async function loadModels() {
    try {
      const models = await window.electronAPI.invoke("get-models");
      modelSelector.innerHTML = ""; // Clear existing options first
      models.gptModels.forEach((model) => {
        const option = document.createElement("option");
        option.value = model.id;
        option.textContent = model.id;
        modelSelector.appendChild(option);
      });
      models.claudeModels.forEach((model) => {
        const option = document.createElement("option");
        option.value = model;
        option.textContent = model;
        modelSelector.appendChild(option);
      });
      models.ollamaModels.forEach((model) => {
        const option = document.createElement("option");
        option.value = model;
        option.textContent = model;
        modelSelector.appendChild(option);
      });
    } catch (error) {
      console.error("Failed to load models:", error);
      alert(
        "Failed to load models. Please check the console for more details."
      );
    }
  }

  // Load patterns and models on startup
  loadPatterns();
  loadModels();

  // Listen for model responses
  window.electronAPI.on("model-response", (message) => {
    const formattedMessage = message.replace(/\n/g, "<br>");
    responseContainer.innerHTML += formattedMessage; // Append new data as it arrives
  });

  window.electronAPI.on("model-response-end", (message) => {
    // Handle the end of the model response if needed
  });

  window.electronAPI.on("model-response-error", (message) => {
    alert(message);
  });

  window.electronAPI.on("file-response", (message) => {
    if (message.startsWith("Error")) {
      alert(message);
      return;
    }
    submitQuery(message);
  });

  window.electronAPI.on("api-keys-saved", async () => {
    try {
      await loadModels();
      alert("API Keys saved successfully.");
      configSection.classList.add("hidden");
      openaiApiKeyInput.value = "";
      claudeApiKeyInput.value = "";
    } catch (error) {
      console.error("Failed to reload models:", error);
      alert("Failed to reload models.");
    }
  });
  updatePatternsButton.addEventListener("click", async () => {
    window.electronAPI.send("update-patterns");
  });

  // Submit button click handler
  submitButton.addEventListener("click", async () => {
    const userInputValue = userInput.value;
    submitQuery(userInputValue);
  });

  fineTuningButton.addEventListener("click", function (e) {
    e.preventDefault();
    fineTuningSection.classList.toggle("hidden");
  });

  temperatureSlider.addEventListener("input", function () {
    temperatureValue.textContent = this.value;
  });

  topPSlider.addEventListener("input", function () {
    topPValue.textContent = this.value;
  });

  frequencyPenaltySlider.addEventListener("input", function () {
    frequencyPenaltyValue.textContent = this.value;
  });

  presencePenaltySlider.addEventListener("input", function () {
    presencePenaltyValue.textContent = this.value;
  });

  submitPatternButton.addEventListener("click", async () => {
    const patternName = document.getElementById("patternName").value;
    const patternText = document.getElementById("patternBody").value;
    document.getElementById("patternName").value = "";
    document.getElementById("patternBody").value = "";
    submitPattern(patternName, patternText);
  });

  // Theme changer click handler
  themeChanger.addEventListener("click", function (e) {
    e.preventDefault();
    document.body.classList.toggle("light-theme");
    themeChanger.innerText =
      themeChanger.innerText === "Dark" ? "Light" : "Dark";
  });

  updatePatternButton.addEventListener("click", function (e) {
    e.preventDefault();
    patternCreator.classList.toggle("hidden");
    myForm.classList.toggle("hidden");

    // window.electronAPI.send("create-pattern");
  });

  // Config button click handler - toggles the config section visibility
  configButton.addEventListener("click", function (e) {
    e.preventDefault();
    configSection.classList.toggle("hidden");
  });

  // Save API Key button click handler
  saveApiKeyButton.addEventListener("click", () => {
    const openAIKey = openaiApiKeyInput.value;
    const claudeKey = claudeApiKeyInput.value;
    window.electronAPI
      .invoke("save-api-keys", { openAIKey, claudeKey })
      .catch((err) => {
        console.error("Error saving API keys:", err);
        alert("Failed to save API Keys.");
      });
  });

  // Handler for pattern selection change
  patternSelector.addEventListener("change", async () => {
    const selectedPattern = patternSelector.value;
    const systemCommand = await window.electronAPI.invoke(
      "get-pattern-content",
      selectedPattern
    );
    // Use systemCommand as part of the input for querying the model
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
