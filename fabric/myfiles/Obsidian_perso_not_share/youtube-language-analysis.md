# YouTube vs Text Pattern Processing

## Text Processing (Working)
```typescript
// User enters: "Here's a news article --fr"
handleInput() {
  languageStore.set('fr')  // Set language
}

sendMessage(articleText) {
  // Goes through ChatService.createChatPrompt
  // Gets language from store
  // Adds language instruction
  // Pattern executes with language context
}
```

## YouTube Processing (Issue)
```typescript
// User enters: "https://youtube.com/... --fr"
handleInput() {
  languageStore.set('fr')  // Set language
}

handleSubmit() {
  // First message
  sendMessage("Processing...")  // Goes through proper flow

  // Get transcript
  const transcript = await getTranscript()

  // Direct to chat service - BYPASSES language setup
  chatService.streamChat(transcript)  // Missing language context
}
```

## Key Difference
For YouTube, we're bypassing the normal sendMessage flow that would add language instructions. Instead, we're sending the transcript directly to chatService.streamChat.

## Solution Direction
We need to ensure the transcript goes through the same flow as regular text:
1. Get transcript
2. Send through sendMessage to get proper language setup
3. Let pattern execute with language context

This way both text and transcripts will get proper language instructions before pattern execution.