<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';
  import { get } from 'svelte/store';
  import type { Pattern } from '$lib/types';
  import { favorites } from '$lib/store/favorites-store';
  import { patternAPI, systemPrompt, selectedPatternName } from '$lib/store/pattern-store';
  import { Input } from "$lib/components/ui/input";

  const dispatch = createEventDispatcher<{
    close: void;
    select: string;
  }>();

  let patterns: Pattern[] = [];
  let patternsContainer: HTMLDivElement;
  let sortBy: 'alphabetical' | 'favorites' = 'alphabetical';
  let searchText = ""; // For pattern filtering

  // First filter patterns by search text
  $: filteredPatterns = patterns.filter(p =>
    p.patternName.toLowerCase().includes(searchText.toLowerCase())
  );

  // Then sort the filtered patterns
  $: sortedPatterns = sortBy === 'alphabetical'
    ? [...filteredPatterns].sort((a, b) => a.patternName.localeCompare(b.patternName))
    : [
        ...filteredPatterns.filter(p => $favorites.includes(p.patternName)).sort((a, b) => a.patternName.localeCompare(b.patternName)),
        ...filteredPatterns.filter(p => !$favorites.includes(p.patternName)).sort((a, b) => a.patternName.localeCompare(b.patternName))
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

<div class="bg-primary-800 rounded-lg flex flex-col h-[85vh] w-[600px] shadow-lg">
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
    
    <div class="px-4 pb-4 flex items-center justify-between">
      <div class="flex gap-4">
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
      <div class="w-64 mr-4">
        <Input
          bind:value={searchText}
          placeholder="Search patterns..."
          class_="text-emerald-700"
        />
      </div>
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
              role="button"
              tabindex="0"
              on:click={async () => {
                try {
                  console.log('Selecting pattern:', pattern.patternName);
                  // Update pattern selection
                  patternAPI.selectPattern(pattern.patternName);
                  // Verify the selection
                  const currentSystemPrompt = get(systemPrompt);
                  const currentPattern = get(selectedPatternName);
                  console.log('After selection - Pattern:', currentPattern);
                  console.log('After selection - System Prompt length:', currentSystemPrompt?.length);
                  
                  // Only close if selection was successful
                  if (currentPattern === pattern.patternName && currentSystemPrompt) {
                    searchText = ""; // Reset search before closing
                    dispatch('select', pattern.patternName);
                    dispatch('close');
                  } else {
                    console.error('Pattern selection failed - Pattern:', currentPattern);
                    console.error('Pattern selection failed - System Prompt:', currentSystemPrompt);
                  }
                } catch (error) {
                  console.error('Error selecting pattern:', error);
                }
              }}
              on:keydown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.preventDefault();
                  e.currentTarget.click();
                }
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