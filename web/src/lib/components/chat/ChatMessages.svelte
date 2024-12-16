<script lang="ts">
    import { chatState, errorStore, streamingStore } from '$lib/store/chat';
    import { afterUpdate } from 'svelte';
    import { marked } from 'marked';
    import SessionManager from './SessionManager.svelte';
    import { fade, slide } from 'svelte/transition';

    let messagesContainer: HTMLDivElement;

    afterUpdate(() => {
        if (messagesContainer) {
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }
    });

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
</script>

<div class="chat-messages-wrapper flex flex-col h-full">
    <div class="flex justify-between items-center mb-4 flex-none">
        <span class="text-sm font-medium">Chat History</span>
        <SessionManager />
    </div>

    {#if $errorStore}
        <div class="error-message" transition:slide>
            <div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-4" role="alert">
                <p>{$errorStore}</p>
            </div>
        </div>
    {/if}

    <div class="messages-container" bind:this={messagesContainer}>
        <div class="messages-content">
            {#each $chatState.messages as message}
                <div 
                    class="message-item {message.role === 'assistant' ? 'pl-4' : 'font-medium'}"
                    transition:fade
                >
                    <div class="message-header flex items-center gap-2 mb-1">
                        <span class="text-xs border rounded-lg p-1 variant-glass-secondary font-bold uppercase">{message.role}</span>
                        {#if message.role === 'assistant' && $streamingStore}
                            <span class="loading-indicator">
                                <span class="dot">.</span>
                                <span class="dot">.</span>
                                <span class="dot">.</span>
                            </span>
                        {/if}
                    </div>

                    {#if message.role === 'assistant'}
                        <div class="prose prose-sm text-inherit max-w-none">
                            {@html renderMarkdown(message.content, true)}
                        </div>
                    {:else}
                        <div class="whitespace-pre-wrap text-sm">
                            {message.content}
                        </div>
                    {/if}
                </div>
            {/each}
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
        flex: 1;
        overflow-y: auto;
        padding-right: 1rem;
    }

    .messages-content {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
    }

    .message-item {
        position: relative;
    }

    .loading-indicator {
        display: inline-flex;
        gap: 2px;
    }

    .dot {
        animation: blink 1.4s infinite;
        opacity: 0;
    }

    .dot:nth-child(2) {
        animation-delay: 0.2s;
    }

    .dot:nth-child(3) {
        animation-delay: 0.4s;
    }

    @keyframes blink {
        0%, 100% { opacity: 0; }
        50% { opacity: 1; }
    }

    :global(.prose pre) {
        background-color: rgb(40, 44, 52);
        color: rgb(171, 178, 191);
        padding: 1rem;
        border-radius: 0.375rem;
        margin: 1rem 0;
    }

    :global(.prose code) {
        color: rgb(171, 178, 191);
        background-color: rgba(40, 44, 52, 0.1);
        padding: 0.2em 0.4em;
        border-radius: 0.25rem;
    }
</style>
