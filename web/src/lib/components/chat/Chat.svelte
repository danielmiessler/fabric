<script lang="ts">
  import ChatInput from "./ChatInput.svelte";
  import ChatMessages from "./ChatMessages.svelte";
  import ModelConfig from "./ModelConfig.svelte";
  import DropdownGroup from "./DropdownGroup.svelte";
  import NoteDrawer from "$lib/components/ui/noteDrawer/NoteDrawer.svelte";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import Tooltip from "$lib/components/ui/tooltip/Tooltip.svelte";
  import { Textarea } from "$lib/components/ui/textarea";
  import { obsidianSettings } from "$lib/store/obsidian-store";
  import { featureFlags } from "$lib/config/features";
  import { getDrawerStore } from '@skeletonlabs/skeleton';
  import { systemPrompt, selectedPatternName } from "$lib/store/pattern-store";
  import { onMount } from "svelte";

  const drawerStore = getDrawerStore();
  function openDrawer() {
    drawerStore.open({});
  }

  // Column width state (percentage values)
  let leftColumnWidth = 50;
  let rightColumnWidth = 50;
  let isDragging = false;
  
  // Handle resize functionality
  function startResize(e: MouseEvent | KeyboardEvent) {
    isDragging = true;
    e.preventDefault();
    
    // Add event listeners for drag and release
    window.addEventListener('mousemove', handleResize);
    window.addEventListener('mouseup', stopResize);
  }
  
  // Handle keyboard events for accessibility
  function handleKeyDown(e: KeyboardEvent) {
    // Only respond to Enter or Space key
    if (e.key === 'Enter' || e.key === ' ') {
      startResize(e);
    }
  }
  
  function handleResize(e: MouseEvent) {
    if (!isDragging) return;
    
    // Get container dimensions
    const container = document.querySelector('.chat-container');
    if (!container) return;
    
    const containerRect = container.getBoundingClientRect();
    const containerWidth = containerRect.width;
    
    // Calculate percentage based on mouse position
    const percentage = ((e.clientX - containerRect.left) / containerWidth) * 100;
    
    // Apply constraints (left: 25-50%, right: 50-75%)
    leftColumnWidth = Math.min(Math.max(percentage, 25), 50);
    rightColumnWidth = 100 - leftColumnWidth;
  }
  
  function stopResize() {
    isDragging = false;
    window.removeEventListener('mousemove', handleResize);
    window.removeEventListener('mouseup', stopResize);
  }

  // Clean up event listeners when component is destroyed
  onMount(() => {
    return () => {
      window.removeEventListener('mousemove', handleResize);
      window.removeEventListener('mouseup', stopResize);
    };
  });

  $: showObsidian = $featureFlags.enableObsidianIntegration;
</script>

<div class="chat-container flex gap-0 p-2 w-full h-screen">
  <!-- Left Column -->
  <aside class="flex flex-col gap-2 pr-2" style="width: {leftColumnWidth}%">
    <!-- Dropdowns Group -->
    <div class="bg-background/5 p-2 rounded-lg">
      <div class="rounded-lg bg-background/10">
        <DropdownGroup />
      </div>
    </div>

    <!-- Model Config -->
    <div class="bg-background/5 p-2 rounded-lg">
      <ModelConfig />
    </div>

    <!-- Message Input -->
    <div class="h-[200px] bg-background/5 rounded-lg overflow-hidden">
      <ChatInput />
    </div>

    <!-- System Instructions -->
    <div class="flex-1 min-h-0 bg-background/5 p-2 rounded-lg">
      <div class="h-full flex flex-col">
        <Textarea
          bind:value={$systemPrompt}
          readonly={true}
          placeholder="System instructions will appear here when you select a pattern..."
          class="w-full flex-1 bg-primary-800/30 rounded-lg border-none whitespace-pre-wrap overflow-y-auto resize-none text-sm scrollbar-thin scrollbar-thumb-white/10 scrollbar-track-transparent hover:scrollbar-thumb-white/20"
        />
      </div>
    </div>
  </aside>

  <!-- Resize Handle -->
  <button 
    class="resize-handle" 
    on:mousedown={startResize}
    on:keydown={handleKeyDown}
    type="button"
    aria-label="Resize chat panels"
  ></button>

  <!-- Right Column -->
  <div class="flex flex-col gap-2" style="width: {rightColumnWidth}%">
    <!-- Header with Obsidian Settings -->
    <div class="flex items-center justify-between px-2 py-1">
      <div class="flex items-center gap-2">
        {#if showObsidian}
          <div class="flex items-center gap-2">
            <div class="flex items-center gap-1">
              <Checkbox
                bind:checked={$obsidianSettings.saveToObsidian}
                id="save-to-obsidian"
                class="h-3 w-3"
              />
              <Label for="save-to-obsidian" class="text-xs text-white/70">Save to Obsidian</Label>
            </div>
            {#if $obsidianSettings.saveToObsidian}
              <Input
                id="note-name"
                bind:value={$obsidianSettings.noteName}
                placeholder="Note name..."
                class="h-6 text-xs w-48 bg-white/5 border-none focus:ring-1 ring-white/20"
              />
            {/if}
          </div>
        {/if}
      </div>
      <Button variant="ghost" size="sm" class="h-6 px-2 text-xs opacity-70 hover:opacity-100" on:click={openDrawer}>
        <Tooltip text="Take Notes" position="left">
          <span>Take Notes</span>
        </Tooltip>
      </Button>
    </div>

    <!-- Chat Area -->
    <div class="flex-1 flex flex-col min-h-0">
      <!-- Chat History -->
      <div class="flex-1 min-h-0 bg-background/5 rounded-lg overflow-y-scroll scrollbar-thin scrollbar-thumb-white/10 scrollbar-track-transparent hover:scrollbar-thumb-white/20">
        <ChatMessages />
        <div class="h-32"></div> <!-- Spacer div to ensure scrolling works properly -->
      </div>
    </div>
  </div>
</div>

<NoteDrawer />

<style>
  .resize-handle {
    width: 6px;
    margin: 0 -3px;
    height: 100%;
    cursor: col-resize;
    position: relative;
    z-index: 10;
    transition: background-color 0.2s;
  }

  .resize-handle::after {
    content: "";
    position: absolute;
    top: 0;
    left: 50%;
    transform: translateX(-50%);
    height: 100%;
    width: 2px;
    background-color: rgba(255, 255, 255, 0.1);
    transition: background-color 0.2s, width 0.2s;
  }

  .resize-handle:hover::after,
  .resize-handle:focus::after {
    background-color: rgba(255, 255, 255, 0.3);
    width: 4px;
  }

  .resize-handle:focus {
    outline: none;
  }

  .resize-handle:focus-visible::after {
    background-color: rgba(255, 255, 255, 0.5);
    width: 4px;
  }

  @keyframes flash {
    0% { opacity: 1; }
    50% { opacity: 0.5; }
    100% { opacity: 1; }
  }
</style>
