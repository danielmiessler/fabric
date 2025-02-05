# Obsidian Integration in Browser

## Overview
The Obsidian integration allows saving pattern output directly to Obsidian markdown files from the browser interface. This is implemented across several components:

## Components

### 1. Obsidian Store (src/lib/store/obsidian-store.ts)
```typescript
export interface ObsidianSettings {
  saveToObsidian: boolean;
  noteName: string;
}

const defaultSettings = {
  saveToObsidian: false,
  noteName: ''
};

export const obsidianSettings = writable<ObsidianSettings>(defaultSettings);
```
Manages Obsidian-related settings and state.

### 2. Model Configuration UI (src/lib/components/chat/ModelConfig.svelte)
```typescript
{#if showObsidian}
  <div class="mt-4 space-y-4 border-t pt-4">
    <Label class="font-bold">Obsidian Settings</Label>
    
    <div class="flex items-center space-x-2">
      <Checkbox
        bind:checked={$obsidianSettings.saveToObsidian}
        id="save-to-obsidian"
      />
      <Label for="save-to-obsidian">Save to Obsidian</Label>
    </div>

    {#if $obsidianSettings.saveToObsidian}
      <div class="space-y-2">
        <Label for="note-name">Note Name</Label>
        <Input
          id="note-name"
          bind:value={$obsidianSettings.noteName}
          placeholder="Enter note name..."
        />
      </div>
    {/if}
  </div>
{/if}
```
Provides UI controls for enabling Obsidian saving and setting note name.

### 3. Chat Input Handler (src/lib/components/chat/ChatInput.svelte)
```typescript
async function saveToObsidian(content: string) {
  try {
    const response = await fetch('/obsidian', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        pattern: $selectedPatternName,
        noteName: $obsidianSettings.noteName,
        content
      })
    });

    const responseData = await response.json();
    
    if (!response.ok) {
      throw new Error(responseData.error || 'Failed to save to Obsidian');
    }

    // Show success message with file name
    toastStore.trigger({
      message: responseData.message || `Saved to Obsidian: ${responseData.fileName}`,
      background: 'variant-filled-success'
    });
  } catch (error) {
    console.error('Failed to save to Obsidian:', error);
    throw error;
  }
}
```
Handles saving pattern output to Obsidian when processing YouTube URLs.

### 4. Obsidian API Endpoint (src/routes/obsidian/+server.ts)
```typescript
export const POST: RequestHandler = async ({ request }) => {
  let tempFile: string | undefined;

  try {
    const body = await request.json() as ObsidianRequest;
    
    // Format content with markdown code blocks
    const formattedContent = `\`\`\`markdown\n${body.content}\n\`\`\``;
    const escapedFormattedContent = escapeShellArg(formattedContent);

    // Generate file name and path
    const fileName = `${new Date().toISOString().split('T')[0]}-${body.noteName}.md`;
    const filePath = `${obsidianDir}/${fileName}`;

    // Write content to file
    await execAsync(`echo ${escapedFormattedContent} > "${filePath}"`);

    return json({
      success: true,
      fileName,
      message: `Successfully saved to ${fileName}`
    });
  } catch (error) {
    return json(
      { error: error instanceof Error ? error.message : 'Failed to process request' },
      { status: 500 }
    );
  }
}
```
Handles file saving on the server side.

## Flow

1. User Interface:
   - Enable Obsidian saving via checkbox
   - Enter note name in input field
   - Both controls in ModelConfig.svelte

2. YouTube Processing:
   - User pastes YouTube URL
   - System gets transcript
   - Processes with selected pattern
   - Shows output in browser

3. File Saving:
   - If Obsidian saving enabled:
     - Sends pattern output to /obsidian endpoint
     - Backend saves file with markdown formatting
     - Frontend shows success/error message

4. File Structure:
   - Files saved to: fabric/myfiles/Fabric_obsidian/
   - File naming: YYYY-MM-DD-{note-name}.md
   - Content wrapped in markdown code blocks

## Implementation Details

1. File Safety:
   - Shell argument escaping for safety
   - Proper file path handling
   - Error handling at all levels

2. User Feedback:
   - Success messages with file names
   - Error messages with details
   - Toast notifications for visibility

3. Code Organization:
   - Store for settings management
   - UI components for user interaction
   - API endpoint for file handling
   - Clear separation of concerns

## Usage

1. Enable Obsidian saving in Model Configuration panel
2. Enter desired note name
3. Process YouTube URL with pattern
4. File is automatically saved with pattern output
5. Success message shows saved file name

The integration provides a seamless way to save pattern output directly to Obsidian markdown files while maintaining the browser-based workflow.
