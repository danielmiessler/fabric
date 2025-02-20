<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Textarea } from "$lib/components/ui/textarea";
  import { sendMessage, messageStore } from '$lib/store/chat-store';
  import { systemPrompt, selectedPatternName } from '$lib/store/pattern-store';
  import { getToastStore } from '@skeletonlabs/skeleton';
  import { FileButton } from '@skeletonlabs/skeleton';
  import { Paperclip, Send, FileCheck } from 'lucide-svelte';
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';
  import { getTranscript } from '$lib/services/transcriptService';
  import { ChatService } from '$lib/services/ChatService';
  // import { obsidianSettings } from '$lib/store/obsidian-store';
  import { languageStore } from '$lib/store/language-store';
  import { obsidianSettings, updateObsidianSettings } from '$lib/store/obsidian-store';

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
      // Add this after successful save
      updateObsidianSettings({ 
      saveToObsidian: false,  // Reset the save flag
      noteName: ''           // Clear the note name
      });
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
          // Add processing message first
          messageStore.update(messages => [...messages, {
              role: 'system',
              content: 'Processing YouTube video...',
              format: 'loading'
          }]);
        
          // Get transcript but don't display it
          const { transcript } = await getTranscript(input);
        
          // Log system prompt BEFORE createChatRequest
          console.log('System prompt BEFORE createChatRequest in YouTube flow:', $systemPrompt);

          // Log system prompt BEFORE streamChat
          console.log(`System prompt BEFORE streamChat in YouTube flow: ${$systemPrompt}`);

          const stream = await chatService.streamChat(transcript, $systemPrompt);
          await chatService.processStream(
              stream,
              (content, response) => {
                  messageStore.update(messages => {
                      const newMessages = [...messages];
                      // Replace the processing message with actual content
                      const lastMessage = newMessages[newMessages.length - 1];
                      if (lastMessage?.format === 'loading') {
                          newMessages.pop();
                      }
                      newMessages.push({
                          role: 'assistant',
                          content,
                          format: response?.format
                      });
                      return newMessages;
                  });
              },
              (error) => {
                  messageStore.update(messages => 
                      messages.filter(m => m.format !== 'loading')
                  );
                  throw error;
              }
          );

          // Handle Obsidian saving if needed
          if ($obsidianSettings.saveToObsidian) {
              let lastContent = '';
              messageStore.subscribe(messages => {
                  const lastMessage = messages[messages.length - 1];
                  if (lastMessage?.role === 'assistant') {
                      lastContent = lastMessage.content;
                  }
              })();
              if (lastContent) await saveToObsidian(lastContent);
          }

          userInput = "";
          uploadedFiles = [];
          fileContents = [];
      } catch (error) {
          console.error('Error processing YouTube URL:', error);
          messageStore.update(messages => 
              messages.filter(m => m.format !== 'loading')
          );
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

<div class="h-full flex flex-col p-2">
  <div class="relative flex-1 min-h-0 bg-primary-800/30 rounded-lg">
    <Textarea
      bind:value={userInput}
      on:input={handleInput}
      on:keydown={handleKeydown}
      placeholder="Enter your message (YouTube URLs will be automatically processed)..."
      class="w-full h-full resize-none bg-transparent border-none text-sm focus:ring-0 transition-colors p-3 pb-[48px]"
    />
    <div class="absolute bottom-3 right-3 flex items-center gap-2">
      <div class="flex items-center gap-2">
        {#if uploadedFiles.length > 0}
          <span class="text-xs text-white/70">
            {uploadedFiles.length} file{uploadedFiles.length > 1 ? 's' : ''} attached
          </span>
        {/if}
        <FileButton
          name="file-upload"
          button="btn-icon variant-ghost"
          bind:files
          on:change={handleFileUpload}
          disabled={isProcessingFiles || uploadedFiles.length >= 5}
          class="h-10 w-10 bg-primary-800/30 hover:bg-primary-800/50 rounded-full transition-colors"
        >
          {#if uploadedFiles.length > 0}
            <FileCheck class="w-5 h-5" />
          {:else}
            <Paperclip class="w-5 h-5" />
          {/if}
        </FileButton>
        <Button
          type="button"
          variant="ghost"
          size="icon"
          name="send"
          on:click={handleSubmit}
          disabled={isProcessingFiles || !userInput.trim()}
          class="h-10 w-10 bg-primary-800/30 hover:bg-primary-800/50 rounded-full transition-colors disabled:opacity-30"
        >
          <Send class="w-5 h-5" />
        </Button>
      </div>
    </div>
  </div>
</div>

<style>
  :global(textarea) {
    scrollbar-width: thin;
    scrollbar-color: rgba(255, 255, 255, 0.2) transparent;
  }

  :global(textarea::-webkit-scrollbar) {
    width: 6px;
  }

  :global(textarea::-webkit-scrollbar-track) {
    background: transparent;
  }

  :global(textarea::-webkit-scrollbar-thumb) {
    background-color: rgba(255, 255, 255, 0.2);
    border-radius: 3px;
  }

  :global(textarea::-webkit-scrollbar-thumb:hover) {
    background-color: rgba(255, 255, 255, 0.3);
  }

  :global(textarea::selection) {
    background-color: rgba(255, 255, 255, 0.1);
  }
</style>
