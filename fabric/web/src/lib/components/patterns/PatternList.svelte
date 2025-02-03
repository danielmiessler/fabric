<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';
  import type { Pattern } from '$lib/types';
  import { favorites } from '$lib/store/favorites-store';
  import { patternAPI } from '$lib/store/pattern-store';

  const dispatch = createEventDispatcher<{
    close: void;
    select: string;
  }>();

  let patterns: Pattern[] = [];
  let patternsContainer: HTMLDivElement;
  let sortBy: 'alphabetical' | 'favorites' = 'alphabetical';

  $: sortedPatterns = sortBy === 'alphabetical'
    ? [...patterns].sort((a, b) => a.patternName.localeCompare(b.patternName))
    : [
        ...patterns.filter(p => $favorites.includes(p.patternName)).sort((a, b) => a.patternName.localeCompare(b.patternName)),
        ...patterns.filter(p => !$favorites.includes(p.patternName)).sort((a, b) => a.patternName.localeCompare(b.patternName))
      ];

  onMount(async () => {
    try {
      const response = await fetch('/data/pattern_descriptions.json');
      const data = await response.json();
      patterns = data.patterns;
    } catch (error) {
      console.error('Error loading patterns:', error);
    }
  });

  function toggleFavorite(patternName: string) {
    favorites.toggleFavorite(patternName);
  }
</script>

<div class="bg-primary-800/70 rounded-lg flex flex-col h-[85vh] w-[600px] shadow-lg">
  <div class="flex flex-col border-b border-primary-700/30">
    <div class="flex justify-between items-center p-4">
      <b class="text-lg text-muted-foreground font-bold">Pattern Descriptions</b>
      <button
        on:click={() => dispatch('close')}
        class="text-muted-foreground hover:text-primary-300 transition-colors"
      >
        ✕
      </button>
    </div>
    
    <div class="px-4 pb-4 flex gap-4">
      <label class="flex items-center gap-2 text-sm text-muted-foreground">
        <input
          type="radio"
          bind:group={sortBy}
          value="alphabetical"
          class="radio"
        >
        Alphabetical
      </label>
      <label class="flex items-center gap-2 text-sm text-muted-foreground">
        <input
          type="radio"
          bind:group={sortBy}
          value="favorites"
          class="radio"
        >
        Favorites First
      </label>
    </div>
  </div>

  <div
    class="patterns-container p-4 flex-1 overflow-y-auto"
    bind:this={patternsContainer}
  >
    <div class="patterns-list space-y-2">
      {#each sortedPatterns as pattern}
        <div class="pattern-item bg-primary/10 rounded-lg p-3">
          <div class="flex justify-between items-start gap-4 mb-2">
            <h3
              class="text-xl font-bold text-primary-300 hover:text-primary-100 cursor-pointer transition-colors"
              on:click={() => {
                patternAPI.selectPattern(pattern.patternName);
                dispatch('select', pattern.patternName);
                dispatch('close');
              }}
            >
              {pattern.patternName}
            </h3>
            <button
              class="text-muted-foreground hover:text-primary-300 transition-colors"
              on:click={() => toggleFavorite(pattern.patternName)}
            >
              {#if $favorites.includes(pattern.patternName)}
                ★
              {:else}
                ☆
              {/if}
            </button>
          </div>
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