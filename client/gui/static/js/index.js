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

  // Submit button click handler
  submitButton.addEventListener("click", async () => {
    responseContainer.innerHTML = ""; // Clear previous responses
    responseContainer.classList.remove("hidden");
    const selectedPattern = patternSelector.value;
    const systemCommand = await window.electronAPI.invoke(
      "get-pattern-content",
      selectedPattern
    );
    const userInputValue = userInput.value;
    window.electronAPI.send(
      "start-query-openai",
      systemCommand,
      userInputValue
    );
  });

  // Theme changer click handler
  themeChanger.addEventListener("click", function (e) {
    e.preventDefault();
    console.log("Theme changer clicked");
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
});
