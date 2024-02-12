const { OpenAI } = require("openai");
require("dotenv").config({
  path: require("os").homedir() + "/.config/fabric/.env",
});

let openaiClient = null;

// Function to initialize and get the OpenAI client
function getOpenAIClient() {
  if (!process.env.OPENAI_API_KEY) {
    throw new Error(
      "The OPENAI_API_KEY environment variable is missing or empty."
    );
  }
  return new OpenAI({ apiKey: process.env.OPENAI_API_KEY });
}

async function queryOpenAI(system, user, callback) {
  const openai = getOpenAIClient(); // Ensure the client is initialized here
  const messages = [
    { role: "system", content: system },
    { role: "user", content: user },
  ];
  try {
    const stream = await openai.chat.completions.create({
      model: "gpt-4-1106-preview", // Adjust the model as necessary.
      messages: messages,
      temperature: 0.0,
      top_p: 1,
      frequency_penalty: 0.1,
      presence_penalty: 0.1,
      stream: true,
    });

    for await (const chunk of stream) {
      const message = chunk.choices[0]?.delta?.content || "";
      callback(message); // Process each chunk of data
    }
  } catch (error) {
    console.error("Error querying OpenAI:", error);
    callback("Error querying OpenAI. Please try again.");
  }
}

module.exports = { queryOpenAI };
