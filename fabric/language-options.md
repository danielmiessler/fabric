# Language Support Implementation

## Overview
The language support allows switching between languages using qualifiers (--fr, --en) in the chat input. The implementation is simple and effective, working at multiple layers of the application.

## Components

### 1. Language Store (language-store.ts)
```typescript
// Manages language state
export const languageStore = writable<string>('');
```

### 2. Chat Input (ChatInput.svelte)
- Detects language qualifiers in user input
- Updates language store
- Strips qualifier from message
```typescript
// Language qualifier handling
if (qualifier === 'fr') {
  languageStore.set('fr');
  userInput = userInput.replace(/--fr\s*/, '');
} else if (qualifier === 'en') {
  languageStore.set('en');
  userInput = userInput.replace(/--en\s*/, '');
}
```

### 3. Chat Service (ChatService.ts)
- Adds language instruction to prompts
- Defaults to English if no language specified
```typescript
const language = get(languageStore) || 'en';
const languageInstruction = language !== 'en' 
  ? `. Please use the language '${language}' for the output.` 
  : '';
const fullInput = userInput + languageInstruction;
```

## How It Works

1. User Input:
   - User types message with language qualifier (e.g., "--fr Hello")
   - ChatInput detects qualifier and updates language store
   - Qualifier is stripped from message

2. Request Processing:
   - ChatService gets language from store
   - Adds language instruction to prompt
   - Sends to backend

3. Response:
   - AI responds in requested language
   - Response is displayed without modification

## Usage Examples

1. English (Default):
```
User: What is the weather?
AI: The weather information...
```

2. French:
```
User: --fr What is the weather?
AI: Voici les informations météo...
```

## Implementation Notes

1. Simple Design:
   - No complex language detection
   - No translation layer
   - Direct instruction to AI

2. Stateful:
   - Language persists until changed
   - Resets to English on page refresh

3. Extensible:
   - Easy to add new languages
   - Just add new qualifiers and store values

## Best Practices

1. Always reset language after message:
```typescript
// Reset stores
languageStore.set('en');
```

2. Default to English:
```typescript
const language = get(languageStore) || 'en';
```

3. Clear language instruction:
```typescript
const languageInstruction = language !== 'en' 
  ? `. Please use the language '${language}' for the output.` 
  : '';