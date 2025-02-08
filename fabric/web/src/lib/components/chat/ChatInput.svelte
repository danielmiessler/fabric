<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Textarea } from "$lib/components/ui/textarea";
  import { sendMessage, messageStore } from '$lib/store/chat-store';
  import { systemPrompt, selectedPatternName } from '$lib/store/pattern-store';
  import { getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';
  import { FileButton } from '@skeletonlabs/skeleton';
  import { Paperclip, Send, FileCheck } from 'lucide-svelte';
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';
  import { getTranscript } from '$lib/services/transcriptService';
  import { ChatService } from '$lib/services/ChatService';
  import type { StreamResponse } from '$lib/interfaces/chat-interface';
  import { obsidianSettings } from '$lib/store/obsidian-store';
  import { languageStore } from '$lib/store/language-store';

  const chatService = new ChatService();
  let userInput = "";
  let isYouTubeURL = false;
  const toastStore = getToastStore();
  let files: FileList | undefined = undefined;
  let uploadedFiles: string[] = [];
  let fileContents: string[] = [];
  let isProcessingFiles = false;

  function detectYouTubeURL(input: string): boolean {
    const youtubePattern = /(?:https?:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)/i;
    const isYoutube = youtubePattern.test(input);
    if (isYoutube) {
      console.log('YouTube URL detected:', input);
      console.log('Current system prompt:', $systemPrompt?.length);
      console.log('Selected pattern:', $selectedPatternName);
    }
    return isYoutube;
  }

  function handleInput(event: Event) {
    console.log('\n=== Handle Input ===');
    const target = event.target as HTMLTextAreaElement;
    userInput = target.value;
    
    const currentLanguage = get(languageStore);
    
    const languageQualifiers = {
      '--en': 'en',
      '--fr': 'fr',
      '--es': 'es',
      '--de': 'de',
      '--zh': 'zh',
      '--ja': 'ja'
    };

    let detectedLang = '';
    for (const [qualifier, lang] of Object.entries(languageQualifiers)) {
      if (userInput.includes(qualifier)) {
        detectedLang = lang;
        languageStore.set(lang);
        userInput = userInput.replace(new RegExp(`${qualifier}\\s*`), '');
        break;
      }
    }

    console.log('2. Language state:', {
      previousLanguage: currentLanguage,
      currentLanguage: get(languageStore),
      detectedOverride: detectedLang,
      inputAfterLangRemoval: userInput
    });

    isYouTubeURL = detectYouTubeURL(userInput);
    console.log('3. URL detection:', {
      isYouTube: isYouTubeURL,
      pattern: $selectedPatternName,
      systemPromptLength: $systemPrompt?.length
    });
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

  async function saveToObsidian(content: string) {
    if (!$obsidianSettings.saveToObsidian) {
      console.log('Obsidian saving is disabled');
      return;
    }
    
    if (!$obsidianSettings.noteName) {
      toastStore.trigger({
        message: 'Please enter a note name in Obsidian settings',
        background: 'variant-filled-error'
      });
      return;
    }

    if (!$selectedPatternName) {
      toastStore.trigger({
        message: 'No pattern selected',
        background: 'variant-filled-error'
      });
      return;
    }

    if (!content) {
      toastStore.trigger({
        message: 'No content to save',
        background: 'variant-filled-error'
      });
      return;
    }

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

      toastStore.trigger({
        message: responseData.message || `Saved to Obsidian: ${responseData.fileName}`,
        background: 'variant-filled-success'
      });
    } catch (error) {
      console.error('Failed to save to Obsidian:', error);
      toastStore.trigger({
        message: error instanceof Error ? error.message : 'Failed to save to Obsidian',
        background: 'variant-filled-error'
      });
    }
  }

  async function processYouTubeURL(input: string) {
    console.log('\n=== YouTube Flow Start ===');
    const originalLanguage = get(languageStore);
    
    try {
        // Get transcript first
        const { transcript } = await getTranscript(input);
        
        // Process with current language and pattern
        await sendMessage(transcript, $systemPrompt);
        
        // Get last message for Obsidian
        let lastContent = '';
        messageStore.subscribe(messages => {
            const lastMessage = messages[messages.length - 1];
            if (lastMessage?.role === 'assistant') {
                lastContent = lastMessage.content;
            }
        })();

        if ($obsidianSettings.saveToObsidian && lastContent) {
            await saveToObsidian(lastContent);
        }

        userInput = "";
        uploadedFiles = [];
        fileContents = [];
    } catch (error) {
        console.error('Error processing YouTube URL:', error);
        messageStore.update(messages => messages.slice(0, -1));
        throw error;
    }
  }

  async function handleSubmit() {
    if (!userInput.trim()) return;

    try {
        console.log('\n=== Submit Handler Start ===');
        
        if (isYouTubeURL) {
            console.log('2a. Starting YouTube flow');
            await processYouTubeURL(userInput);
            return;
        }
        
        const finalContent = fileContents.length > 0 
            ? userInput + '\n\nFile Contents:\n' + fileContents.join('\n\n')
            : userInput;
            
        await sendMessage(finalContent);
        
        userInput = "";
        uploadedFiles = [];
        fileContents = [];
    } catch (error) {
        console.error('Chat submission error:', error);
    }
  }

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
    <div class="messages-container">
      <Textarea
        bind:value={$systemPrompt}
        readonly={true}
        placeholder="Enter system instructions..."
        class="w-full h-[300px] bg-primary-800/30 rounded-lg border-none whitespace-pre-wrap overflow-y-auto"
      />
    </div>
  </div>

  <div class="flex-1 relative shadow-lg">
    <div class="text-xs text-gray-400 mb-1">Enter your message (YouTube URLs will be automatically processed)</div>
    <Textarea
      bind:value={userInput}
      on:input={handleInput}
      on:keydown={handleKeydown}
      placeholder="Enter your message..."
      class="w-full h-full resize-none bg-primary-800/30 rounded-lg border-none"
    />
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

  .pattern-textarea::selection {
    background-color: rgba(155, 155, 155, 0.3);
  }

  :global(textarea) {
    scrollbar-width: thin;
    -ms-overflow-style: thin;
  }

  :global(textarea::-webkit-scrollbar) {
    width: 8px;
  }

  :global(textarea::-webkit-scrollbar-track) {
    background: transparent;
  }

  :global(textarea::-webkit-scrollbar-thumb) {
    background-color: rgba(155, 155, 155, 0.5);
    border-radius: 4px;
  }

  :global(textarea::-webkit-scrollbar-thumb:hover) {
    background-color: rgba(155, 155, 155, 0.7);
  }
</style>
