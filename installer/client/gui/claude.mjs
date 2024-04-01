import Claude from "claude-ai";

export function MakeClaude(apiKey) {
  return new Claude({ sessionKey: apiKey });
}
