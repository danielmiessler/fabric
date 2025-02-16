# Language Configuration Fix for YouTube vs Text Processing

## Working Configuration

### ChatService.ts
- Simple language handling that works:
```typescript
// Add language instruction at end of input
userInput: userInput + languageInstruction
```

### ChatInput.svelte
- Store and restore language state:
```typescript
// Store original language
const originalLanguage = get(languageStore);

// Restore before processing
languageStore.set(originalLanguage);

// Pass isYouTubeTranscript flag
await sendMessage(transcript, $systemPrompt, false, true);
```

### chat/+server.ts
- Pass language through transcript service:
```typescript
const response = {
  transcript,
  title: videoId,
  language: body.language
};
```

## Issue to Fix: Transcript Display

The transcript is showing in browser before pattern output because it's being added to the message store. The fix needs to:

1. Keep the "Processing YouTube transcript..." message
2. Skip adding the raw transcript to message store
3. Only show the pattern output when it comes back

### Location of Fix
In chat-store.ts:
```typescript
// Add isYouTubeTranscript parameter to skip adding transcript
export async function sendMessage(content: string, systemPromptText?: string, isSystem: boolean = false, isYouTubeTranscript: boolean = false)

// Skip only the transcript, not system messages
if (isYouTubeTranscript && !isSystem) {
  // Skip only transcript
  return;
}
messageStore.update(messages => [...messages, { 
  role: isSystem ? 'system' : 'user', 
  content 
}]);
```

## Testing
1. Regular text with --fr:
   - Should show message in chat
   - Should process in French

2. YouTube with --fr:
   - Should show "Processing..." message (isSystem=true)
   - Should NOT show transcript (isYouTubeTranscript=true)
   - Should show pattern output in French (language preserved)
