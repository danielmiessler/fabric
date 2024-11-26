<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { RotateCcw, Trash2, Save } from 'lucide-svelte';
    import { chatState, clearMessages, revertLastMessage, currentSession } from '$lib/store/chat';
    import { afterUpdate } from 'svelte';
    
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
</script>

<div class="chat-messages-wrapper flex flex-col h-full">
    <div class="flex justify-between items-center mb-4 flex-none">
        <span class="text-sm font-medium">Chat History</span>
        <div class="flex gap-2">
            <Button variant="outline" size="icon" on:click={revertLastMessage}>
                <RotateCcw class="h-4 w-4" />
            </Button>
            <Button variant="outline" size="icon" on:click={clearMessages}>
                <Trash2 class="h-4 w-4" />
            </Button>
            <Button variant="outline" size="icon" on:click={saveChat}>
                <Save class="h-4 w-4" />
            </Button>
        </div>
    </div>

    <div class="messages-container flex-1" bind:this={messagesContainer}>
        <div class="messages-content space-y-4">
            {#each $chatState.messages as message}
                <div class="message-item whitespace-pre-wrap text-sm {message.role === 'assistant' ? 'pl-4' : 'font-medium'} transition-all">
                    <span class="text-xs tertiary uppercase">{message.role}:</span>
                    {message.content}
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
</style>
