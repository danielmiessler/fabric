document.addEventListener("DOMContentLoaded", function () {
  const socket = io.connect("/");
  const patternSelector = document.getElementById("patternSelector");
  const userInput = document.getElementById("userInput");
  const submitButton = document.getElementById("submit");
  const responseContainer = document.getElementById("responseContainer");
  const themeChanger = document.getElementById("themeChanger");

  fetch("/patterns/get")
    .then((res) => res.json())
    .then((patterns) => {
      for (const name in patterns) {
        let option = document.createElement("option");
        option.value = name;
        option.textContent = name;
        patternSelector.appendChild(option);
      }
    });

  submitButton.addEventListener("click", () => {
    responseContainer.classList.remove("hidden");
    const input = userInput.value;
    const pattern = patternSelector.value;
    socket.emit("fabric", { module: pattern, input_data: input });
    userInput.value = ""; // Clear input field after submission
  });
  themeChanger.addEventListener("click", function (e) {
    e.preventDefault();
    document.body.classList.toggle("light-theme");
    if (themeChanger.innerText == "Dark") {
      themeChanger.innerText = "Light";
    } else {
      themeChanger.innerText = "Dark";
    }
  });

  socket.on("message", function (data) {
    responseContainer.innerHTML += data.replace(/\n/g, "<br>");
  });

  socket.on("error", function (error) {
    responseContainer.innerHTML += `<div class="error">Error: ${error}</div>`;
  });
});
