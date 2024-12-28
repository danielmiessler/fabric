<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Textarea } from "$lib/components/ui/textarea";
    import { sendMessage, messageStore } from '$lib/store/chat';
    import { systemPrompt } from '$lib/store/pattern-store';
    import { getToastStore } from '@skeletonlabs/skeleton';
    import { FileButton } from '@skeletonlabs/skeleton';
    import { Paperclip, Send } from 'lucide-svelte';
    import { onMount } from 'svelte';

    let userInput = "";
    let files: FileList;
    const toastStore = getToastStore();

    async function handleSubmit() {
        if (!userInput.trim()) return;

        try {
            const trimmedInput = userInput.trim();
            const trimmedSystemPrompt = $systemPrompt.trim();
            
            // Clear input before sending to improve perceived performance
            userInput = "";
            
            await sendMessage(trimmedSystemPrompt + '\n' + trimmedInput);
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

<div class="h-full">
  <div class="flex flex-col gap-2 h-screen">
    <div class="flex-1 rounded-lg border-current">
      <Textarea 
        bind:value={$systemPrompt}
        on:input={(e) => $systemPrompt || ''}
        placeholder="Enter system instructions..."
        class="w-full h-full resize-none bg-primary-800/30 rounded-lg border-none"
      />
    </div>

    <div class="flex-1 py-2 relative shadow-lg">
      <Textarea
        bind:value={userInput}
        on:input={(e) => userInput}
        on:keydown={handleKeydown}
        placeholder="Enter your message..."
        class="w-full h-full resize-none bg-primary-800/30 rounded-lg border-none"
      />
      <div class="absolute bottom-5 right-2 gap-2 flex justify-end end-7">

        <FileButton
          name="file-upload"
          button="btn btn-sm variant-soft-surface"
          bind:files={files}
          on:change={(e) => {
            // Workin on the file selection
            // Check out `https://www.skeleton.dev/components/file-buttons` for more info
            // Check 24-12-08 for half-baked implementation
          }}
        >
          <Paperclip class="w-4" />
        </FileButton>
        <Button type="button" name="submit" variant="secondary" on:click={handleSubmit}>
          <Send class="w-4 h-4" />
        </Button>
      </div>
    </div>
  </div>
</div>

<style>
    .flex-col {
        min-height: 0;
    }
</style>
