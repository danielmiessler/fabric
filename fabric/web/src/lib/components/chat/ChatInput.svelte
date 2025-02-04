<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Textarea } from "$lib/components/ui/textarea";
  import { sendMessage, messageStore } from '$lib/store/chat-store';
  import { systemPrompt, selectedPatternName } from '$lib/store/pattern-store';
  import { getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';
  import { FileButton } from '@skeletonlabs/skeleton';
  import { Paperclip, Send, FileCheck } from 'lucide-svelte';
  import { onMount } from 'svelte';
  import { getTranscript } from '$lib/services/transcriptService';
  import { ChatService } from '$lib/services/ChatService';
  import type { StreamResponse } from '$lib/interfaces/chat-interface';

  const chatService = new ChatService();

  let userInput = "";
  let isYouTubeURL = false;
  const toastStore = getToastStore();

  let files: FileList | undefined = undefined;

  function detectYouTubeURL(input: string): boolean {
    const youtubePattern = /(?:https?:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)/i;
    const isYoutube = youtubePattern.test(input);
    if (isYoutube) {
      console.log('YouTube URL detected:', input);
      // Log current pattern selection state
      console.log('Current system prompt:', $systemPrompt?.length);
      console.log('Selected pattern:', $selectedPatternName);
    }
    return isYoutube;
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

      if (isYouTubeURL) {
        console.log('Processing YouTube URL in handleSubmit');
        console.log('Current pattern state:');
        console.log('- Selected Pattern:', $selectedPatternName);
        console.log('- System Prompt length:', $systemPrompt?.length);
        console.log('- Message content:', trimmedInput);
        
        try {
          // Show processing message
          await sendMessage("Processing YouTube transcript...", $systemPrompt, true);
          
          // Get transcript but don't display it
          const { transcript } = await getTranscript(trimmedInput);
          console.log('Got transcript, length:', transcript.length);
          
          // Get stream from chat service
          const stream = await chatService.streamChat(transcript, $systemPrompt);
          
          // Process stream directly
          await chatService.processStream(
            stream,
            (content: string) => {
              messageStore.update(messages => {
                const newMessages = [...messages];
                const lastMessage = newMessages[newMessages.length - 1];
                
                if (lastMessage?.role === 'assistant') {
                  lastMessage.content = content;
                } else {
                  newMessages.push({
                    role: 'assistant',
                    content
                  });
                }
                
                return newMessages;
              });
            },
            (error: Error) => {
              console.error('Stream processing error:', error);
              toastStore.trigger({
                message: 'Error processing transcript',
                background: 'variant-filled-error'
              });
            }
          );
          
          userInput = "";
          uploadedFiles = [];
          fileContents = [];
        } catch (error) {
          console.error('Error processing YouTube URL:', error);
          toastStore.trigger({
            message: 'Failed to process YouTube video. Please try again.',
            background: 'variant-filled-error'
          });
        }
      } else {
        userInput = "";
        uploadedFiles = [];
        fileContents = [];
        
        // Send regular message
        await sendMessage(trimmedInput);
      }
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
    <div class="messages-container p-4">
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
