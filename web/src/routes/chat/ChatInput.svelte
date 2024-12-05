<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Textarea } from "$lib/components/ui/textarea";
    // import { Label } from "$lib/components/ui/label";
    import { FileButton } from '@skeletonlabs/skeleton';
    import { Paperclip, Send } from 'lucide-svelte';
    import { sendMessage } from '$lib/store/chat';
    // import { currentSession, setSession } from '$lib/store/chat';
    import { systemPrompt } from '$lib/types/chat/patterns';
    import { onMount } from 'svelte';

    let userInput = "";
    let files: FileList;

    async function handleSubmit() {
        if (!userInput.trim()) return;

        try {
            await sendMessage($systemPrompt.trim() + userInput.trim());
        } catch (error) {
            console.error('Chat submission error:', error);
        }
    }

    /* function handleSetSession(name: string | null) {
        setSession(name);
    } */

    onMount(() => {
        console.log('ChatInput mounted, current system prompt:', $systemPrompt);
    });
</script>

<div class="flex flex-col gap-2">
    <div class="flex-none">
        <Textarea 

            value={$systemPrompt}
            on:input={(e) => $systemPrompt}
            placeholder="Enter system instructions..."
            class="min-h-[330px] resize-none bg-background"
        />
    </div>

    <div class="flex-1 py-2 relative">
        <Textarea
            bind:value={userInput}
            on:input={(e) => userInput}
            placeholder="Enter your message..."
            class="min-h-[350px] resize-none bg-background"
        />
        <div class="absolute bottom-5 right-2 gap-2 flex justify-end end-7">
        <!-- TODO: Session Management. Move this to a new component, possibly ChatMessages -->
            <!-- <div class="flex gap-2">
                <button type="button" class="btn btn-sm variant-soft-tertiary"
                    on:click={() => handleSetSession(null)}
                >
                    Clear Session
                </button>
                {#if !$currentSession}
                    <button type="button" class="btn btn-sm variant-glass-tertiary"
                        on:click={() => handleSetSession('new-session-' + Date.now())}
                    >
                        New Session
                    </button>
                {/if}
            </div> -->
            <FileButton
                name="file-upload"
                button="btn btn-sm variant-soft-surface"
                bind:files={files}
                on:change={(e) => {
                    // Workin on the file selection
                    // Check out `https://www.skeleton.dev/components/file-buttons` for more info
                }}
            >
                <Paperclip class="w-4" />
            </FileButton>
            <Button type="button" name="submit" class="btn btn-sm variant-filled-secondary" on:click={handleSubmit}>
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