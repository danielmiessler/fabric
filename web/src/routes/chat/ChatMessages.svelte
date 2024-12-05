<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { RotateCcw, Trash2, Save, Copy } from 'lucide-svelte';
    import { chatState, clearMessages, revertLastMessage, currentSession } from '$lib/store/chat';
    import { afterUpdate } from 'svelte';
    import { marked } from 'marked';
    import { getToastStore } from '@skeletonlabs/skeleton';
    import { Toast } from '@skeletonlabs/skeleton';

    let sessionName: string | null = null;
    let messagesContainer: HTMLDivElement;

    currentSession.subscribe(value => {
        sessionName = value;
    });

    afterUpdate(() => {
        if (messagesContainer) {
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }
    });

    function saveChat() {
        const chatData = JSON.stringify($chatState.messages, null, 2);
        const blob = new Blob([chatData], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'chat-history.json';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }

    marked.setOptions({
        gfm: true,
        breaks: true,
    });

    function renderMarkdown(content: string, isAssistant: boolean) {
        if (!isAssistant) return content;
        try {
            return marked.parse(content);
        } catch (error) {
            console.error('Error rendering markdown:', error);
            return content;
        }
      }

    const toastStore = getToastStore();

    async function copyToClipboard() {
      try {
        await navigator.clipboard.writeText($chatState.messages.map(m => m.content).join('\n'));
        toastStore.trigger({
            message: 'Chat copied to clipboard!',
            background: 'variant-filled-success'
        });
      } catch (err) {
        toastStore.trigger({
            message: 'Failed to copy transcript',
            background: 'variant-filled-error'
        });
      }
    }
</script>

<div class="chat-messages-wrapper flex flex-col h-full">
    <div class="flex justify-between items-center mb-4 flex-none">
        <span class="text-sm font-medium">Chat History</span>
        <div class="flex gap-2">
            <Button class="variant-glass-tertiary" variant="outline" size="icon" aria-label="Revert Last Message" on:click={revertLastMessage}>
                <RotateCcw class="h-4 w-4" />
            </Button>
            <Button class="variant-glass-tertiary" variant="outline" size="icon" aria-label="Clear Chat" on:click={clearMessages}>
                <Trash2 class="h-4 w-4" />
            </Button>
            <Button class="variant-glass-tertiary" variant="outline" size="icon" aria-label="Copy Chat" on:click={copyToClipboard}>
                <Copy class="h-4 w-4" />
            </Button>
            <Toast position="b" />
            <Button class="variant-glass-tertiary" variant="outline" size="icon" aria-label="Save Chat" on:click={saveChat}>
                <Save class="h-4 w-4" />
            </Button>
        </div>
    </div>

    <div class="messages-container" bind:this={messagesContainer}>
        <div class="messages-content">
            {#each $chatState.messages as message}
                <div class="message-item {message.role === 'assistant' ? 'pl-4' : 'font-medium'} transition-all">
                    <span class="text-xs tertiary uppercase">{message.role}:</span>
                    {#if message.role === 'assistant'}
                        {@html renderMarkdown(message.content, true)}
                    {:else}
                        <div class="whitespace-pre-wrap text-sm">
                            {message.content}
                        </div>
                    {/if}
                </div>
            {/each}
            {#if $chatState.isStreaming}
                <div class="pl-4 text-tertiary-700 animate-pulse">▌</div>
            {/if}
        </div>
    </div>
</div>

<style>
    .chat-messages-wrapper {
        display: flex;
        flex-direction: column;
        min-height: 0;
    }

    .messages-container {
        position: relative;
        overflow-y: auto;
        scrollbar-width: thin;
        scrollbar-color: var(--color-primary-500) transparent;
    }

    .messages-content {
        padding-bottom: 1rem;
        padding-right: 0.5rem;
    }

    .message-item {
        margin-bottom: 0.5rem;
    }

    .messages-container::-webkit-scrollbar {
        width: 2px;
    }

    .messages-container::-webkit-scrollbar-track {
        background: transparent;
    }

    .messages-container::-webkit-scrollbar-thumb {
        background-color: var(--color-primary-500);
        border-radius: 2px;
    }

    /* Markdown content styles */
    :global(.message-item.pl-4) {
        font-size: 0.875rem;
    }

    :global(.message-item pre) {
        background-color: rgb(50, 50, 50);
        padding: 1rem;
        border-radius: 0.5rem;
        margin: 0.5rem 0;
        overflow-x: auto;
    }

    :global(.message-item code) {
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
        font-size: 0.875rem;
        padding: 0.2rem 0.4rem;
        border-radius: 0.25rem;
        background-color: rgb(50, 50, 60);
    }

    :global(.message-item h1) {
        margin: 0.5rem 0;
        font: bold 1.5rem/1.5 system-ui, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
    }

    :global(.message-item h2) {
        margin: 0.5rem 0;
        font: bold 1.25rem/1.5 system-ui, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
    }

    :global(.message-item h3) {
        margin: 0.5rem 0;
        font: bold 1rem/1.5 system-ui, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
    }

    :global(.message-item h4) {
        margin: 0.5rem 0;
        font: bold 0.875rem/1.5 system-ui, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
    }

    :global(.message-item h5) {
        margin: 0.5rem 0;
        font: bold 0.75rem/1.5 system-ui, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
    }

    :global(.message-item h6) {
        margin: 0.5rem 0;
        font: bold 0.625rem/1.5 system-ui, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
    }
    
    :global(.message-item p) {
        margin: 0.5rem 0;
    }

    :global(.message-item ul, .message-item ol) {
        margin: 0.5rem 0;
        padding-left: 1.5rem;
    }

    :global(.message-item li) {
        margin: 0.25rem 0;
    }

    :global(.message-item a) {
        color: rgb(var(--color-primary-600));
        text-decoration: underline;
    }

    :global(.message-item blockquote) {
        border-left: 4px solid rgb(var(--color-secondary-200));
        margin: 0.5rem 0;
        padding-left: 1rem;
        font-style: italic;
    }
</style>
