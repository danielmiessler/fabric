<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';
  import type { Pattern } from '$lib/types';

  const dispatch = createEventDispatcher<{
    close: void;
  }>();

  let patterns: Pattern[] = [];
  let patternsContainer: HTMLDivElement;

  onMount(async () => {
    try {
      const response = await fetch('/data/pattern_descriptions.json');
      const data = await response.json();
      patterns = data.patterns;
    } catch (error) {
      console.error('Error loading patterns:', error);
    }
  });
</script>

<div class="bg-primary-800/70 rounded-lg flex flex-col h-[85vh] w-[600px] shadow-lg">
  <div class="flex justify-between items-center p-4 border-b border-primary-700/30">
    <b class="text-lg text-muted-foreground font-bold">Pattern Descriptions</b>
    <button
      on:click={() => dispatch('close')}
      class="text-muted-foreground hover:text-primary-300 transition-colors"
    >
      ✕
    </button>
  </div>

  <div
    class="patterns-container p-4 flex-1 overflow-y-auto"
    bind:this={patternsContainer}
  >
    <div class="patterns-list space-y-2">
      {#each patterns as pattern}
        <div class="pattern-item bg-primary/10 rounded-lg p-3">
          <h3 class="text-xl font-bold mb-2 text-primary-300">{pattern.patternName}</h3>
          <p class="text-sm text-muted-foreground break-words leading-relaxed">{pattern.description}</p>
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

  .patterns-list {
    display: flex;
    flex-direction: column;
    width: 100%;
    max-width: 560px;
    margin: 0 auto;
  }

  .pattern-item {
    display: flex;
    flex-direction: column;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .pattern-item:last-child {
    border-bottom: none;
  }
</style>