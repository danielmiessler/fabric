<script lang="ts">
  import { onMount } from 'svelte';
  import { Select } from "$lib/components/ui/select";
  import { Label } from "$lib/components/ui/label";
  import { Slider } from "$lib/components/ui/slider";
  import { modelConfig, availableModels, loadAvailableModels } from "$lib/store/model-config";
  import Transcripts from "./Transcripts.svelte";
  import { patterns } from "$lib/types/chat/patterns";
  import { patternAPI } from "$lib/types/chat/patterns";

	let selectedPreset = "";
	
	$: if (selectedPreset) {
		console.log('Pattern selected:', selectedPreset);
		patternAPI.selectPattern(selectedPreset);
	}

    onMount(async () => {
      await loadAvailableModels();
      await patternAPI.loadPatterns();
    });
</script>

<div class="space-y-2 max-w-full">
    <div class="flex flex-col gap-2">
        <div class="min-w-0">
  	    	<Select 
  	    		bind:value={selectedPreset}
  	    	> 
  	    	<option value="">Load a pattern...</option>
  	    		{#each $patterns as pattern}
  	    			<option value={pattern.Name}>{pattern.Description}</option>
  	    		{/each}
  	    	</Select>
        </div>
        <div class="min-w-0">
            <Select 
                bind:value={$modelConfig.model}
            >
            <option value="">Default Model</option>
                {#each $availableModels as model (model.name)}
                  <option value={model.name}>{model.vendor} - {model.name}</option>
                {/each}
            </Select>
        </div>
    </div>
  
    <div class="space-y-1">
        <Label class="p-1 font-bold">Temperature ({$modelConfig.temperature.toFixed(1)})</Label>
        <Slider
            bind:value={$modelConfig.temperature}
            min={0}
            max={2}
            step={0.1}
        />
    </div>

    <div class="space-y-1">
        <Label class="p-1 font-bold">Maximum Length ({$modelConfig.maxLength})</Label>
        <Slider
            bind:value={$modelConfig.maxLength}
            min={1}
            max={4000}
            step={1}
        />
    </div>

    <div class="space-y-1">
        <Label class="p-1 font-bold">Top P ({$modelConfig.top_p.toFixed(2)})</Label>
        <Slider
            bind:value={$modelConfig.top_p}
            min={0}
            max={1}
            step={0.01}
        />
    </div>

    <div class="space-y-1">
        <Label class="p-1 font-bold">Frequency Penalty ({$modelConfig.frequency.toFixed(2)})</Label>
        <Slider
            bind:value={$modelConfig.frequency}
            min={0}
            max={1}
            step={0.01}
        />
    </div>

    <div class="space-y-1">
        <Transcripts />
    </div>
</div>
