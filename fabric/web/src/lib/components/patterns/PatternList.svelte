<script lang="ts">
  import { onMount } from 'svelte';
  import type { Pattern } from '$lib/types';

  let patterns: Pattern[] = [];
  let patternsContainer: HTMLDivElement;

  onMount(async () => {
    try {
      const response = await fetch('/myfiles/pattern_descriptions.json');
      const data = await response.json();
      patterns = data.patterns;
    } catch (error) {
      console.error('Error loading patterns:', error);
    }
  });
</script>

<div class="bg-primary-800/30 rounded-lg flex flex-col h-full shadow-lg">
  <div class="flex justify-between items-center mb-1 mt-1 flex-none pl-4">
    <b class="text-sm text-muted-foreground font-bold">Pattern List</b>
  </div>

  <div 
    class="patterns-container p-4 flex-1 overflow-y-auto max-h-dvh"
    bind:this={patternsContainer}
  >
    <div class="patterns-grid">
      {#each patterns as pattern}
        <div class="pattern-item bg-primary/5 rounded-lg p-3">
          <h3 class="text-sm font-semibold mb-1 text-primary-300">{pattern.patternName}</h3>
          <p class="text-sm text-muted-foreground">{pattern.description}</p>
        </div>
      {/each}
    </div>
  </div>
</div>

<style>
  .patterns-container {
    flex: 1;
    overflow-y: auto;
    scrollbar-width: thin;
    -ms-overflow-style: thin;
  }

  .patterns-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
  }

  .pattern-item {
    display: flex;
    flex-direction: column;
  }

  @media (max-width: 768px) {
    .patterns-grid {
      grid-template-columns: 1fr;
    }
  }
</style>