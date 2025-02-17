<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';
  import { get } from 'svelte/store';
  import type { Pattern } from '$lib/interfaces/pattern-interface';
  import { favorites } from '$lib/store/favorites-store';
  import { patterns, patternAPI, systemPrompt, selectedPatternName } from '$lib/store/pattern-store';
  import { Input } from "$lib/components/ui/input";
  import TagFilterPanel from './TagFilterPanel.svelte';
  
  let tagFilterRef: TagFilterPanel;


  const dispatch = createEventDispatcher<{
    close: void;
    select: string;
    tagsChanged: string[];  // Add this line
}>();


let patternsContainer: HTMLDivElement;
let sortBy: 'alphabetical' | 'favorites' = 'alphabetical';
let searchText = ""; // For pattern filtering
let selectedTags: string[] = [];

// First filter patterns by both text and tags
// First filter patterns by both text and tags
$: filteredPatterns = $patterns
    .filter((p: Pattern) => 
        p.Name.toLowerCase().includes(searchText.toLowerCase())
    )
    .filter((p: Pattern) => 
        selectedTags.length === 0 || 
        (p.tags && selectedTags.every(tag => p.tags.includes(tag)))
    );

// Then sort the filtered patterns
$: sortedPatterns = sortBy === 'alphabetical'
    ? [...filteredPatterns].sort((a: Pattern, b: Pattern) => a.Name.localeCompare(b.Name))
    : [
        ...filteredPatterns.filter((p: Pattern) => $favorites.includes(p.Name)).sort((a: Pattern, b: Pattern) => a.Name.localeCompare(b.Name)),
        ...filteredPatterns.filter((p: Pattern) => !$favorites.includes(p.Name)).sort((a: Pattern, b: Pattern) => a.Name.localeCompare(b.Name))
    ];


function handleTagFilter(event: CustomEvent<string[]>) {
    selectedTags = event.detail;
}


  onMount(async () => {
    try {
      await patternAPI.loadPatterns();
    } catch (error) {
      console.error('Error loading patterns:', error);
    }
  });

  function toggleFavorite(name: string) {
    favorites.toggleFavorite(name);
  }

  
 



</script>

<div class="bg-primary-800 rounded-lg flex flex-col h-[85vh] w-[600px] shadow-lg relative">

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
                class="text-emerald-900"
            />
        </div>
    </div>

    <!-- New tag display section -->
    <div class="px-4 pb-2">
      <div class="text-sm text-white/70 bg-primary-700/30 rounded-md p-2 flex justify-between items-center">
          <div>Tags: {selectedTags.length ? selectedTags.join(', ') : 'none'}</div>
          <button 
              class="px-2 py-1 text-xs text-white/70 bg-primary-600/30 rounded hover:bg-primary-600/50 transition-colors"
              on:click={() => {
                  selectedTags = [];
                  dispatch('tagsChanged', selectedTags);
              }}
          >
              reset
          </button>
      </div>
  </div>
  
  
</div>


  <TagFilterPanel 
  patterns={$patterns} 
  on:tagsChanged={handleTagFilter}
  bind:this={tagFilterRef}
/>

  <div
    class="patterns-container p-4 flex-1 overflow-y-auto"
    bind:this={patternsContainer}
  >
    <div class="patterns-list space-y-2">
      {#each sortedPatterns as pattern}
      <div class="pattern-item bg-primary/10 rounded-lg p-3">
        <div class="flex justify-between items-start gap-4 mb-2">
            <button
                class="text-xl font-bold text-primary-300 hover:text-primary-100 cursor-pointer transition-colors text-left w-full"
                on:click={() => {
                    console.log('Selecting pattern:', pattern.Name);
                    patternAPI.selectPattern(pattern.Name);
                    searchText = ""; 
                    tagFilterRef.reset();
                    dispatch('select', pattern.Name);
                    dispatch('close');
                }}
                on:keydown={(e) => {
                    if (e.key === 'Enter' || e.key === ' ') {
                        e.preventDefault();
                        e.currentTarget.click();
                    }
                }}
            >
                {pattern.Name}
            </button>
            <button
                class="text-muted-foreground hover:text-primary-300 transition-colors"
                on:click={() => toggleFavorite(pattern.Name)}
            >
                {#if $favorites.includes(pattern.Name)}
                    ★
                {:else}
                    ☆
                {/if}
            </button>
        </div>
        <p class="text-sm text-muted-foreground break-words leading-relaxed">{pattern.Description}</p>
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
