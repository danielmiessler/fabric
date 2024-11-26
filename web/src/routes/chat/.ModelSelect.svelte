<!-- Moved to ModelConfig.svelte -->
 
<!-- <script lang="ts">
    import { onMount } from 'svelte';
    import { availableModels, modelConfig, loadAvailableModels } from '$lib/store/model-config';
    
    onMount(async () => {
        await loadAvailableModels();
    });
</script>

<select 
    class="flex h-9 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50 appearance-none"
    bind:value={$modelConfig.model}
>
    <option value="">Select a model...</option>
    {#each $availableModels as model}
        <option value={model.name}>
            {model.name} ({model.vendor})
        </option>
    {/each}
</select>  -->