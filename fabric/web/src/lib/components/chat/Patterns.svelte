<script lang="ts">
  import { onMount } from 'svelte';
  import { Select } from "$lib/components/ui/select";
  import { patterns, patternAPI, systemPrompt, selectedPatternName } from "$lib/store/pattern-store";
  import { get } from 'svelte/store';

  let selectedPreset = "";

  // Only update pattern selection when selectedPreset changes from user selection
  $: if (selectedPreset) {
    console.log('Pattern selected from dropdown:', selectedPreset);
    try {
      patternAPI.selectPattern(selectedPreset);
      // Verify the selection
      const currentSystemPrompt = get(systemPrompt);
      const currentPattern = get(selectedPatternName);
      console.log('After dropdown selection - Pattern:', currentPattern);
      console.log('After dropdown selection - System Prompt length:', currentSystemPrompt?.length);
      
      if (!currentPattern || !currentSystemPrompt) {
        console.error('Pattern selection verification failed:');
        console.error('- Selected Pattern:', currentPattern);
        console.error('- System Prompt:', currentSystemPrompt);
      }
    } catch (error) {
      console.error('Error in pattern selection:', error);
    }
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
