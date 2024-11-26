<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Textarea } from "$lib/components/ui/textarea";
    import { Label } from "$lib/components/ui/label";
    import { currentSession, setSession, sendMessage } from '$lib/store/chat';
    import { systemPrompt } from '$lib/types/chat/patterns';
    import { onMount } from 'svelte';

    let userInput = "";

    async function handleSubmit() {
        if (!userInput.trim()) return;

        try {
            await sendMessage($systemPrompt.trim() + userInput.trim());
        } catch (error) {
            console.error('Chat submission error:', error);
        }
    }

    function handleSetSession(name: string | null) {
        setSession(name);
    }

    onMount(() => {
        console.log('ChatInput mounted, current system prompt:', $systemPrompt);
    });
</script>

<div class="flex flex-col gap-2">
    <div class="flex gap-2">
        <Button 
            variant="outline" 
            size="sm"
            on:click={() => handleSetSession(null)}
        >
            Clear Session
        </Button>
        {#if !$currentSession}
            <Button 
                variant="outline" 
                size="sm"
                on:click={() => handleSetSession('new-session-' + Date.now())}
            >
                New Session
            </Button>
        {/if}
    </div>

    <div class="flex flex-col h-full">
        <div class="space-y-2 flex-none">
                <Label class="p-1 font-bold">System Prompt</Label>
                <Textarea
                    value={$systemPrompt}
                    on:input={(e) => systemPrompt}
                    placeholder="Enter system instructions..."
                    class="min-h-[200px] resize-none bg-background"
                />
        </div>

        <div class="space-y-2 flex-1 py-1">
            <Label class="p-1 font-bold">User Input</Label>
            <Textarea
                bind:value={userInput}
                placeholder="Enter your message..."
                class="h-[calc(80vh-24rem)] resize-none bg-background"
            />
        </div>
        <Button on:click={handleSubmit} variant="outline" size="sm" class="mt-2">Submit</Button>
    </div>
</div>

<style>
    .flex-col {
        min-height: 0;
    }
</style>