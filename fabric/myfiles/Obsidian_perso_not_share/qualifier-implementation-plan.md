# Plan for Implementing Qualifier Support in Web Interface

## Current Issue
The web interface currently treats qualifiers (like -g=fr) as part of the message text instead of processing them as command flags.

## Implementation Plan

### 1. Add Qualifier Interface
```typescript
// Add to chat-interface.ts
export interface ChatQualifiers {
  language?: string;      // -g flag
  temperature?: number;   // -t flag
  topP?: number;         // -T flag
  presencePenalty?: number;  // -P flag
  frequencyPenalty?: number; // -F flag
  raw?: boolean;         // -r flag
  model?: string;        // -m flag
}

// Update ChatRequest to include qualifiers
export interface ChatRequest {
  prompts: ChatPrompt[];
  messages: Message[];
  qualifiers?: ChatQualifiers;
  // ... existing fields
}
```

### 2. Add Qualifier Parser
```typescript
// New file: src/lib/utils/qualifier-parser.ts
export function parseQualifiers(input: string): {
  text: string;
  qualifiers: ChatQualifiers;
} {
  const qualifierRegex = /-([a-zA-Z]+)(?:=([^\s]+))?/g;
  const qualifiers: ChatQualifiers = {};
  
  // Remove qualifiers and get clean text
  const text = input.replace(qualifierRegex, (match, flag, value) => {
    switch (flag) {
      case 'g':
        qualifiers.language = value;
        break;
      case 't':
        qualifiers.temperature = parseFloat(value);
        break;
      // ... handle other qualifiers
    }
    return '';
  }).trim();

  return { text, qualifiers };
}
```

### 3. Modify ChatInput Component
```typescript
// In ChatInput.svelte
import { parseQualifiers } from '$lib/utils/qualifier-parser';

async function handleSubmit() {
  if (!userInput.trim()) return;

  try {
    const { text, qualifiers } = parseQualifiers(userInput);
    
    // Add qualifiers to request
    const request = await chatService.createChatRequest(text);
    request.qualifiers = qualifiers;
    
    // Send to backend
    const stream = await chatService.streamChat(request);
    // ... rest of the code
  } catch (error) {
    // ... error handling
  }
}
```

### 4. Update ChatService
```typescript
// In ChatService.ts
public async createChatRequest(userInput: string): Promise<ChatRequest> {
  const config = get(chatConfig);
  return {
    prompts: [{
      userInput,
      systemPrompt: get(systemPrompt),
      model: config.model,
      patternName: get(selectedPatternName)
    }],
    messages: get(messageStore),
    ...config
  };
}
```

### 5. Update Backend API
The backend API at localhost:8080/api/chat already supports these qualifiers, so we just need to ensure they're properly included in the request body.

## Testing Plan

1. Test Basic Qualifier Parsing:
```typescript
const input = "hello -g=fr -t=0.8";
const { text, qualifiers } = parseQualifiers(input);
// Should return:
// text: "hello"
// qualifiers: { language: "fr", temperature: 0.8 }
```

2. Test Multiple Qualifiers:
```typescript
const input = "-g=fr -t=0.8 -P=0.2 hello world";
// Should parse all qualifiers and extract clean text
```

3. Test Invalid Qualifiers:
```typescript
const input = "hello -g=fr -invalid=123";
// Should ignore invalid qualifiers
```

## Next Steps

1. Implement the qualifier parser
2. Update the chat interfaces
3. Modify ChatInput to use the parser
4. Add error handling for invalid qualifiers
5. Add validation for qualifier values
6. Update documentation

## Future Enhancements

1. Add UI controls for common qualifiers
2. Add autocomplete for qualifier flags
3. Add validation feedback for invalid qualifiers
4. Add help text showing available qualifiers