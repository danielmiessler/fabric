# Language Support Implementation Plan

## Overview
We will implement language support that allows users to switch between languages using qualifiers (--fr, --en) in chat input or through a global settings dropdown. The implementation will be simple but effective, working at multiple layers of the application.

## Implementation Steps

### 1. Create Language Store
Create `fabric/web/src/lib/store/language-store.ts`:
```typescript
import { writable } from 'svelte/store';
export const languageStore = writable<string>('en'); // Default to English
```

### 2. Modify ChatInput.svelte
1. Import language store:
```typescript
import { languageStore } from '$lib/store/language-store';
```

2. Add language qualifier detection in handleInput:
```typescript
function handleInput(event: Event) {
  const target = event.target as HTMLTextAreaElement;
  userInput = target.value;
  
  // Check for language qualifiers
  if (userInput.includes('--fr')) {
    languageStore.set('fr');
    userInput = userInput.replace(/--fr\s*/, '');
  } else if (userInput.includes('--en')) {
    languageStore.set('en');
    userInput = userInput.replace(/--en\s*/, '');
  }
  
  isYouTubeURL = detectYouTubeURL(userInput);
}
```

3. Reset language after message send in handleSubmit:
```typescript
try {
  await sendMessage(trimmedInput);
  languageStore.set('en'); // Reset to default
} catch (error) {
  console.error('Chat submission error:', error);
  // ... error handling
}
```

### 3. Modify ChatService.ts
1. Import language store:
```typescript
import { languageStore } from '$lib/store/language-store';
```

2. Update createChatPrompt to include language instruction:
```typescript
private createChatPrompt(userInput: string, systemPromptText?: string): ChatPrompt {
  const config = get(modelConfig);
  const language = get(languageStore);
  const languageInstruction = language !== 'en' 
    ? `. Please use the language '${language}' for the output.` 
    : '';

  return {
    userInput: userInput + languageInstruction,
    systemPrompt: systemPromptText ?? get(systemPrompt),
    model: config.model,
    patternName: get(selectedPatternName)
  };
}
```

### 4. Create Language Selector Component
Create `fabric/web/src/lib/components/settings/LanguageSelector.svelte`:
```svelte
<script lang="ts">
  import { Label } from "$lib/components/ui/label";
  import { Select } from "$lib/components/ui/select";
  import { languageStore } from '$lib/store/language-store';

  let selectedLanguage = $languageStore;
  $: languageStore.set(selectedLanguage);
</script>

<div class="flex flex-col gap-2">
  <Label>Language</Label>
  <Select bind:value={selectedLanguage}>
    <option value="en">English</option>
    <option value="fr">French</option>
  </Select>
</div>
```

### 5. Integration
1. Add LanguageSelector to ModelConfig.svelte:
```svelte
<script>
  import LanguageSelector from '../settings/LanguageSelector.svelte';
</script>

<!-- Add to existing settings UI -->
<LanguageSelector />
```

## Testing Plan

1. Test Language Qualifiers:
   - Send message with --fr qualifier
   - Verify message is sent without qualifier
   - Verify response is in French
   - Verify language resets to English after

2. Test Global Settings:
   - Change language in dropdown
   - Send message without qualifier
   - Verify response is in selected language
   - Verify language persists until changed

3. Test Error Cases:
   - Invalid qualifiers
   - Unknown languages
   - Error handling during message send

## Future Enhancements

1. Add more languages:
   - Spanish (--es)
   - German (--de)
   - Italian (--it)

2. Improve language persistence:
   - Option to keep language setting
   - Save preference in local storage

3. Add language detection:
   - Auto-detect input language
   - Suggest matching response language

## Notes

- Keep implementation simple and maintainable
- Focus on user experience
- Ensure clear error handling
- Document all changes
- Consider accessibility in UI components