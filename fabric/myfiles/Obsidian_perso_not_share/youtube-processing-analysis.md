# YouTube Processing Flow Analysis

## Current Flow vs sendMessage Flow

### Current Direct Flow
```typescript
// Get transcript
const transcript = await getTranscript(url);

// Send directly to chat service
const stream = await chatService.streamChat(transcript);

// Process stream manually
await chatService.processStream(stream, ...);
```
**Issue:** Bypasses language setup in ChatService

### sendMessage Flow
```typescript
// In chat-store.ts
async function sendMessage(content, systemPrompt) {
  // 1. Add user message
  messageStore.update(...);

  // 2. Get stream WITH language instruction
  const stream = await chatService.streamChat(content, systemPrompt);

  // 3. Process stream and preserve format
  await chatService.processStream(stream, ...);
}
```
**Benefits:** 
- Handles streaming state
- Gets language instruction from ChatService
- Preserves message format
- Updates message store correctly

## Solution Path

We can use sendMessage for the transcript because it:
1. Already handles streaming properly
2. Gets language instruction from ChatService
3. Preserves format in message updates
4. Manages message store correctly

### Implementation Steps
1. Keep initial processing message:
```typescript
await sendMessage("Processing...", systemPrompt, true);
```

2. Get transcript as before:
```typescript
const { transcript } = await getTranscript(url);
```

3. Use sendMessage instead of direct chatService:
```typescript
// This will:
// - Get language instruction from ChatService
// - Handle streaming
// - Preserve format
// - Update messages
await sendMessage(transcript, systemPrompt);
```

4. Keep Obsidian save:
```typescript
if ($obsidianSettings.saveToObsidian) {
  await saveToObsidian(lastContent);
}
```

This approach ensures the transcript gets proper language handling while preserving all existing functionality:
- Streaming updates
- Format preservation
- Message store management
- Obsidian integration

The key insight is that sendMessage already handles everything we need, including getting language instructions from ChatService. We just need to use it for the transcript instead of going directly to chatService.streamChat.