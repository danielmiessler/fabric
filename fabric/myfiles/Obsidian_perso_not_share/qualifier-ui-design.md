# Simplified Qualifier UI Design

## Overview
Add a qualifier selector under the model dropdown that allows users to add qualifiers, showing active ones as tags. Language qualifier can persist, while others reset with pattern execution.

## UI Components

### Language Selector (Persistent)
```svelte
<!-- In ModelConfig.svelte -->
<div class="language-select">
  <Label>Language</Label>
  <Select
    value={selectedLanguage}
    onChange={handleLanguageChange}
    options={[
      { value: '', label: 'Default' },
      { value: 'en', label: 'English' },
      { value: 'fr', label: 'French' },
      { value: 'zh', label: 'Chinese' }
      // Add more languages
    ]}
  />
</div>
```

### Pattern-Specific Qualifiers
```svelte
<!-- Under pattern selection -->
<div class="qualifier-section">
  <!-- Active Qualifiers -->
  <div class="active-qualifiers">
    {#each activeQualifiers as qualifier}
      <QualifierTag 
        type={qualifier.type}
        value={qualifier.value}
        onRemove={() => removeQualifier(qualifier.id)}
      />
    {/each}
  </div>

  <!-- Add Qualifier Dropdown -->
  <div class="add-qualifier">
    <Select
      placeholder="Add qualifier..."
      options={[
        { value: 'u', label: 'Scrape URL (-u)', hasValue: true },
        { value: 'a', label: 'Attachment (-a)', hasValue: true },
        { value: 'q', label: 'Search Question (-q)', hasValue: true },
        { value: 'readability', label: 'Readability', hasValue: false },
        { value: 'c', label: 'Copy to Clipboard (-c)', hasValue: false },
        { value: 'v', label: 'Pattern Variable (-v)', hasValue: true },
        { value: 'r', label: 'Raw Mode (-r)', hasValue: false }
      ]}
      bind:value={selectedQualifier}
    />

    <!-- Value Input (shows only when qualifier needs a value) -->
    {#if selectedQualifier?.hasValue}
      <div class="qualifier-value">
        {#if selectedQualifier.value === 'v'}
          <div class="variable-pair">
            <Input placeholder="key" bind:value={varKey} />
            <Input placeholder="value" bind:value={varValue} />
          </div>
        {:else}
          <Input 
            placeholder={getPlaceholder(selectedQualifier.value)}
            bind:value={qualifierValue}
          />
        {/if}
      </div>
    {/if}

    <!-- Action Buttons -->
    <div class="qualifier-actions">
      <Button 
        disabled={!isValidQualifier()} 
        on:click={addQualifier}
      >
        Add
      </Button>
    </div>
  </div>
</div>
```

### Qualifier Tag Component
```svelte
<div class="qualifier-tag" class:persistent={type === 'g'}>
  <span class="type">{getDisplayType(type)}</span>
  {#if value}
    <span class="value">{getDisplayValue(type, value)}</span>
  {/if}
  <button class="remove" on:click={onRemove}>×</button>
</div>

<style>
  .qualifier-tag {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.25rem 0.5rem;
    background: var(--accent-2);
    border-radius: 4px;
    font-size: 0.875rem;
  }

  .persistent {
    background: var(--accent-3);
    border: 1px solid var(--accent-4);
  }
</style>
```

## Store Implementation

```typescript
// Split into two stores: language and pattern-specific qualifiers
const languageStore = writable<string>('');

interface PatternQualifier {
  id: string;
  type: string;
  value?: string;
  key?: string;  // For pattern variables
}

function createQualifierStore() {
  const { subscribe, update, set } = writable<PatternQualifier[]>([]);
  
  return {
    subscribe,
    add: (type: string, value?: string, key?: string) => update(quals => [
      ...quals,
      { 
        id: crypto.randomUUID(),
        type,
        value,
        key
      }
    ]),
    remove: (id: string) => update(quals => 
      quals.filter(q => q.id !== id)
    ),
    reset: () => set([]), // Called when pattern executes
    
    getCommandString: () => {
      const quals = get(qualifierStore);
      const lang = get(languageStore);
      
      const commands = [];
      
      // Add language if set
      if (lang) {
        commands.push(`-g=${lang}`);
      }
      
      // Add pattern-specific qualifiers
      commands.push(...quals.map(q => {
        switch (q.type) {
          case 'u':
            return `-u=${q.value}`;
          case 'a':
            return `-a=${q.value}`;
          case 'q':
            return `-q=${q.value}`;
          case 'readability':
            return '--readability';
          case 'c':
            return '-c';
          case 'v':
            return `-v=${q.key}:${q.value}`;
          case 'r':
            return '-r';
          default:
            return '';
        }
      }));
      
      return commands.filter(Boolean).join(' ');
    }
  };
}
```

## Integration

```typescript
// In ChatInput.svelte
import { qualifierStore } from '$lib/store/qualifier-store';
import { languageStore } from '$lib/store/language-store';

async function handleSubmit() {
  const qualifiers = qualifierStore.getCommandString();
  const message = qualifiers 
    ? `${qualifiers} ${userInput.trim()}`
    : userInput.trim();
    
  await sendMessage(message);
  qualifierStore.reset(); // Reset pattern qualifiers only
}

// Reset qualifiers when pattern changes
$: selectedPatternName, qualifierStore.reset();
```

## Benefits

1. Clear separation between persistent language selection and pattern-specific qualifiers
2. Simple tag-based interface for active qualifiers
3. Easy to add/remove qualifiers
4. Automatic reset of pattern-specific qualifiers
5. Language selection persists across pattern executions
6. Clean, focused UI that matches existing components