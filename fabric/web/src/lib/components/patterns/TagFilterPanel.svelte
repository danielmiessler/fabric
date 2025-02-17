<script lang="ts">
    import type { Pattern } from '$lib/interfaces/pattern-interface';
    import { createEventDispatcher } from 'svelte';
    
    const dispatch = createEventDispatcher<{
        tagsChanged: string[];
    }>();

    export let patterns: Pattern[];
    let selectedTags: string[] = [];
    let isExpanded = false;

    // Add console log to see what tags we're getting
    $: console.log('Available tags:', Array.from(new Set(patterns.flatMap(p => p.tags))));

     // Add these debug logs
     $: console.log('Patterns received:', patterns);
    $: console.log('Tags extracted:', patterns.map(p => p.tags));
    $: console.log('Panel expanded:', isExpanded);

    function toggleTag(tag: string) {
        selectedTags = selectedTags.includes(tag)
            ? selectedTags.filter(t => t !== tag)
            : [...selectedTags, tag];
        dispatch('tagsChanged', selectedTags);
    }

    function togglePanel() {
        isExpanded = !isExpanded;
    }

    export function reset() {
        selectedTags = [];
        isExpanded = false;
        dispatch('tagsChanged', selectedTags);
    }
</script>

<div class="tag-panel {isExpanded ? 'expanded' : ''}" style="z-index: 50">
    <div class="panel-header">
        <button class="close-btn" on:click={togglePanel}>
            {isExpanded ? 'Close Filter Tags ◀' : 'Open Filter Tags ▶'}
        </button>
        
    </div>
    
    <div class="panel-content">
        <div class="reset-container">
            <button 
                class="reset-btn"
                on:click={() => {
                    selectedTags = [];
                    dispatch('tagsChanged', selectedTags);
                }}
            >
                Reset All Tags
            </button>
        </div>
        {#each Array.from(new Set(patterns.flatMap(p => p.tags))).sort() as tag}
            <button 
                class="tag-brick {selectedTags.includes(tag) ? 'selected' : ''}"
                on:click={() => toggleTag(tag)}
            >
                {tag}
            </button>
        {/each}
    </div>
</div>
<style>
   .tag-panel {
    position: fixed;  /* Change to fixed positioning */
    left: calc(50% + 300px); /* Position starts after modal's right edge */
    top: 50%;
    transform: translateY(-50%);
    width: 300px;
    transition: left 0.3s ease;
}

.tag-panel.expanded {
    left: calc(50% + 360px); /* Final position just to the right of modal */
}


    .panel-content {
    display: none;
    padding: 12px;
    flex-wrap: wrap;
    gap: 6px;
    max-height: 80vh;
    overflow-y: auto;
    grid-template-columns: repeat(3, 1fr);
}

.tag-brick {
    padding: 4px 8px;
    font-size: 0.8rem;
    border-radius: 12px;
    background: rgba(255,255,255,0.1);
    cursor: pointer;
    white-space: nowrap;
    text-overflow: ellipsis;
    overflow: hidden;
}


.reset-container {
    width: 100%;
    padding-bottom: 8px;
    margin-bottom: 8px;
    border-bottom: 1px solid rgba(255,255,255,0.1);
}

.reset-btn {
    width: 100%;
    padding: 6px;
    font-size: 0.8rem;
    color: var(--primary-300);
    background: rgba(255,255,255,0.05);
    border-radius: 4px;
    transition: all 0.2s;
}

.reset-btn:hover {
    background: rgba(255,255,255,0.1);
}



    .expanded .panel-content {
        display: flex;
    }


   /*  .toggle-btn {
    position: absolute;
    left: -30px;
    top: 50%;
    transform: translateY(-50%);
    padding: 8px;
    background: var(--primary-800);
    border-radius: 4px 0 0 4px;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 0.9rem;
    box-shadow: -2px 0 5px rgba(0,0,0,0.2); 
} */


.panel-header {
    padding: 8px;
    border-bottom: 1px solid rgba(255,255,255,0.1);
}

.close-btn {
    width: auto;
    padding: 6px;
    position: absolute;
    font-size: 0.8rem;
    color: var(--primary-300);
    background: rgba(255,255,255,0.05);
    border-radius: 4px;
    transition: all 0.2s;
    text-align: left;
}

/* Position for 'Open Filter Tags' */
.tag-panel:not(.expanded) .close-btn {
    top: -290px;  /* Moves up to search bar level */
    margin-left: 10px;
}

/* Position for 'Close Filter Tags' */
.expanded .close-btn {
    position: relative;
    top: 0;
    margin-left: -50px;
}




.close-btn:hover {
    background: rgba(255,255,255,0.1);
}

   

    .tag-brick.selected {
        background: var(--primary-300);
    }
</style>

