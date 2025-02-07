# Language Configuration Fix Documentation

## Current State

### Language Implementation
- Language support is implemented through qualifiers (--fr, --en, etc.)
- Language state is managed in `languageStore`
- Language instruction is appended to user input: `. Please use the language '${language}' for the output.`
- Language resets to English after each message send

### Two Processing Flows
1. Regular Text Flow:
```
User Input -> Language Qualifier Detection -> Add Language Instruction -> Pattern Processing -> Output
```

2. YouTube Flow:
```
User Input -> Get Transcript -> Pattern Processing -> Stream Output in Chunks -> Try to Add Language to Each Chunk
```

## Observed Problem

### Issue Description
- Language functionality works consistently for regular text content
- Language functionality is inconsistent for YouTube content
- Works once but fails on subsequent attempts, even after browser refresh

### Root Cause Analysis
1. Language Reset:
   - Language store automatically resets to English after each message
   - This reset happens before streaming is complete for YouTube content

2. Streaming Chunks:
   - YouTube content is processed and streamed in chunks
   - Current approach tries to add language instruction to output chunks
   - Language state may be lost between chunks
   - Chunks are already generated in English before language instruction is attempted

3. Key Difference:
   - Regular text: Language instruction added to input before processing
   - YouTube: Attempting to add language instruction after processing

## Proposed Fix

### Solution
Add language instruction to transcript before pattern processing:
```typescript
// In processYouTubeURL():
const { transcript } = await getTranscript(input);
const currentLanguage = get(languageStore);
    
// Add language instruction to transcript before processing
const languageInstruction = currentLanguage !== 'en'
  ? `. Please use the language '${currentLanguage}' for the output.`
  : '';
const transcriptWithLanguage = transcript + languageInstruction;
    
const stream = await chatService.streamChat(transcriptWithLanguage, $systemPrompt);
```

This minimal change ensures:
- Language instruction is added before processing, like regular text flow
- No changes to core language functionality
- Maintains existing behavior for other content types

### Implementation Plan
1. Document current state and solution (this document)
2. Add language instruction to transcript before processing
3. Test YouTube content with different languages

### Expected Outcome
- YouTube content processes language consistently like other content
- No changes to existing language functionality
- Simple, focused fix for the specific issue

## Status
- [x] Documentation Created
- [x] Transcript Processing Modified
- [ ] Testing Completed

## Updates
2/7/2025 3:40 PM:
- Added detailed logging in ChatService to verify request structure
- Fixed message handling to properly show pattern output
- Keep "Processing..." message while showing pattern output
- Added logging to verify language instruction is included

## Status
- [x] Documentation Created
- [x] Transcript Processing Modified
- [ ] Testing Completed

## Testing Instructions
1. Use a YouTube URL with language qualifier (e.g. --fr)
2. Check browser console to verify:
   - Language is set correctly
   - Request includes language instruction
   - Pattern name is included
3. In browser UI verify:
   - "Processing..." message appears
   - Pattern output appears (not transcript)
   - Output is in correct language
4. Test multiple videos to ensure consistent behavior
