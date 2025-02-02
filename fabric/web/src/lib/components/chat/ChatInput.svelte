<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Textarea } from "$lib/components/ui/textarea";
  import { sendMessage, messageStore } from '$lib/store/chat-store';
  import { systemPrompt, selectedPatternName } from '$lib/store/pattern-store';
  import { getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';
  import { FileButton } from '@skeletonlabs/skeleton';
  import { Paperclip, Send, FileCheck } from 'lucide-svelte';
  import { onMount } from 'svelte';

  let userInput = "";
  let isYouTubeURL = false;
  const toastStore = getToastStore();

  let files: FileList | undefined = undefined;

  function detectYouTubeURL(input: string): boolean {
    const youtubePattern = /(?:https?:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)/i;
    return youtubePattern.test(input);
  }
  let uploadedFiles: string[] = [];
  let fileContents: string[] = [];
  let isProcessingFiles = false;

  function handleInput(event: Event) {
    const target = event.target as HTMLTextAreaElement;
    userInput = target.value;
    isYouTubeURL = detectYouTubeURL(userInput);
  }

  async function handleFileUpload(e: Event) {
    if (!files || files.length === 0) return;

    if (uploadedFiles.length >= 5 || (uploadedFiles.length + files.length) > 5) {
      toastStore.trigger({
        message: 'Maximum 5 files allowed',
        background: 'variant-filled-error'
      });
      return;
    }

    isProcessingFiles = true;
    try {
      for (let i = 0; i < files.length && uploadedFiles.length < 5; i++) {
        const file = files[i];
        const content = await readFileContent(file);
        fileContents.push(content);
        uploadedFiles = [...uploadedFiles, file.name];
      }
    } catch (error) {
      toastStore.trigger({
        message: 'Error processing files: ' + (error as Error).message,
        background: 'variant-filled-error'
      });
    } finally {
      isProcessingFiles = false;
    }
  }

  function readFileContent(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => resolve(e.target?.result as string);
      reader.onerror = (e) => reject(new Error('Failed to read file'));
      reader.readAsText(file);
    });
  }

  async function handleSubmit() {
    if (!userInput.trim()) return;

    try {
      let finalContent = "";
      if (fileContents.length > 0) {
        finalContent += '\n\nFile Contents:\n' + fileContents.map((content, index) => 
          `[${uploadedFiles[index]}]:\n${content}`
        ).join('\n\n');
      }

      const trimmedInput = userInput.trim() + '\n' + (finalContent || '');
      let messageHistory = JSON.stringify($messageStore);

      userInput = "";
      uploadedFiles = [];
      fileContents = [];

      // Send just the user input, let the ChatService handle the pattern
      await sendMessage(trimmedInput);
    } catch (error) {
      console.error('Chat submission error:', error);
      toastStore.trigger({
        message: 'Failed to send message. Please try again.',
        background: 'variant-filled-error'
      });
    }
  }

  // Handle keyboard shortcuts
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      handleSubmit();
    }
  }

  onMount(() => {
    console.log('ChatInput mounted, current system prompt:', $systemPrompt);
  });
</script>

<div class="flex flex-col gap-2 h-full">
  <div class="flex-1 relative shadow-lg">
    <Textarea 
      bind:value={$systemPrompt}
      on:input={(e) => $systemPrompt || ''}
      placeholder="Enter system instructions..."
      class="w-full h-full resize-none bg-primary-800/30 rounded-lg border-none"
    />
  </div>

  <div class="flex-1 relative shadow-lg">
    <Textarea
      bind:value={userInput}
      on:input={handleInput}
      on:keydown={handleKeydown}
      placeholder="Enter your message... (YouTube URLs will be automatically processed)"
      class="w-full h-full resize-none bg-primary-800/30 rounded-lg border-none"
    />
    {#if isYouTubeURL}
      <div class="absolute top-2 right-2">
        <span class="text-sm text-blue-400 bg-blue-900/30 px-2 py-1 rounded">
          YouTube URL detected - transcript will be processed
        </span>
      </div>
    {/if}
    <div class="absolute bottom-5 right-2 gap-2 flex justify-end end-7">
      <FileButton
        name="file-upload"
        button="btn variant-default"
        bind:files
        on:change={handleFileUpload}
        disabled={isProcessingFiles || uploadedFiles.length >= 5}
      >
        {#if uploadedFiles.length > 0}
          <FileCheck class="w-4 h-4" />
        {:else}
          <Paperclip class="w-4 h-4" />
        {/if}
      </FileButton>
      {#if uploadedFiles.length > 0}
        <span class="text-sm text-gray-500 space-x-2">
          {uploadedFiles.length} file{uploadedFiles.length > 1 ? 's' : ''} attached
        </span>
      {/if}
      <br>
      <Button
        type="button"
        variant="default"
        name="send"
        on:click={handleSubmit}
        disabled={isProcessingFiles || !userInput.trim()}
      >
        <Send class="w-4 h-4" />
      </Button>
    </div>
  </div>
</div>

<style>
.flex-col {
  min-height: 0;
}
</style>
