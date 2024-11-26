<script lang="ts">
  import { onMount } from 'svelte';
  import { Select } from "$lib/components/ui/select";
  import { Label } from "$lib/components/ui/label";
  import { Slider } from "$lib/components/ui/slider";
  import { modelConfig, availableModels, loadAvailableModels } from "$lib/store/model-config";
  import Transcripts from "./Transcripts.svelte";

  onMount(async () => {
    await loadAvailableModels();
  });

  //Debugging
  //$: console.log('Current available models:', $availableModels);
  //$: console.log('Current model config:', $modelConfig);
</script>

<div class="space-y-4">
  <div class="space-y-2">
    <Label>Model</Label>
    <Select bind:value={$modelConfig.model}>
      <option value="">Default Model</option>
      {#each $availableModels as model (model.name)}
        <option value={model.name}>{model.vendor} - {model.name}</option>
      {/each}
    </Select>
  </div>

  <div class="space-y-2">
    <Label>Temperature ({$modelConfig.temperature.toFixed(1)})</Label>
    <Slider
      bind:value={$modelConfig.temperature}
      min={0}
      max={2}
      step={0.1}
    />
  </div>

  <div class="space-y-2">
    <Label>Maximum Length ({$modelConfig.maxLength})</Label>
    <Slider
      bind:value={$modelConfig.maxLength}
      min={1}
      max={4000}
      step={1}
    />
  </div>

  <div class="space-y-2">
    <Label>Top P ({$modelConfig.top_p.toFixed(2)})</Label>
    <Slider
      bind:value={$modelConfig.top_p}
      min={0}
      max={1}
      step={0.01}
    />
  </div>
  <div class="space-y-2">
    <Transcripts />
  </div>
</div>
