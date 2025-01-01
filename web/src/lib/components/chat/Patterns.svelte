<script lang="ts">
  import { onMount } from 'svelte';
  import { Select } from "$lib/components/ui/select";
  import { patterns, patternAPI } from "$lib/store/pattern-store";

  let selectedPreset = "";

  $: if (selectedPreset) {
    console.log('Pattern selected:', selectedPreset);
    patternAPI.selectPattern(selectedPreset);
  }

    onMount(async () => {
      await patternAPI.loadPatterns();
    });
</script>

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
