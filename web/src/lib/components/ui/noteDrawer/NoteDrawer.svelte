<script lang="ts">
  import { Drawer, getDrawerStore, getToastStore } from '@skeletonlabs/skeleton';
  import type { DrawerSettings, DrawerStore } from '@skeletonlabs/skeleton';
  import { onMount } from 'svelte';
  import { noteStore } from '$lib/store/noteStore';
  import { afterNavigate, beforeNavigate } from '$app/navigation'; 
  
  const drawerStore = getDrawerStore();
  const toastStore = getToastStore();
  
  let textareaEl: HTMLTextAreaElement;
  let saving = false;

  let content = '';
  
  // Auto-resize textarea
  function adjustTextareaHeight() {
    if (textareaEl) {
      textareaEl.style.height = 'auto';
      textareaEl.style.height = textareaEl.scrollHeight + 'px';
    }
  }
  
  async function saveContent() {
    if (!$noteStore.content.trim()) {
      toastStore.trigger({
        message: 'Cannot save empty note',
        background: 'variant-filled-warning'
      });
      return;
    }

    try {
      saving = true;
      await noteStore.save();
      
      toastStore.trigger({
        message: `Note saved successfully!`,
        background: 'variant-filled-success'
      });
    } catch (error) {
      console.error('Failed to save note:', error);
      toastStore.trigger({
        message: error instanceof Error ? error.message : 'Failed to save notes',
        background: 'variant-filled-error'
      });
    } finally {
      saving = false;
    }
  }
  
  // Prompt user if trying to close with unsaved changes
  $: if ($drawerStore.open === false && $noteStore.isDirty) {
    if (confirm('You have unsaved changes. Are you sure you want to close?')) {
      noteStore.reset();
    } else {
      drawerStore.open({});
    }
  }
  
  // Load saved content when drawer opens
  $: if ($drawerStore.open) {
    const savedContent = localStorage.getItem('savedText');
    if (savedContent) {
      noteStore.updateContent(savedContent);
      noteStore.save();
    }
  }
  
  // Keyboard shortcuts
  function handleKeydown(event: KeyboardEvent) {
    if ((event.ctrlKey || event.metaKey) && event.key === 's') {
      event.preventDefault();
      saveContent();
    }
  }
  
  onMount(() => {
    adjustTextareaHeight();
  });
</script>

<Drawer width="w-[40%]" padding="p-4">
  {#if $drawerStore.open}
    <div class="space-y-4">
      <header class="flex justify-between items-center">
        <h2 class="m-2 p-1 h2">Notes</h2>
        <p class="p-2 opacity-70">Notes are saved to <code>`src/lib/content/inbox`</code></p>
        <p class="p-2 opacity-70">Ctrl + S to save</p>
        {#if $noteStore.lastSaved}
          <span class="text-sm opacity-70">
            Last saved: {$noteStore.lastSaved.toLocaleTimeString()}
          </span>
        {/if}
      </header>
      <div class="m-4">
      <textarea
        bind:this={textareaEl}
        bind:value={$noteStore.content}
        on:input={adjustTextareaHeight}
        on:keydown={handleKeydown}
        class="w-full min-h-screen p-2 rounded-container-token text-primary-800 resize-none"
        placeholder="Enter your text here..."
      />
      </div>
      <footer class="flex justify-between pb-4 items-center">
        <span class="text-sm opacity-70">
          {#if $noteStore.isDirty}
            Unsaved changes
          {/if}
        </span>
        <button
          class="btn px-4 variant-filled-primary"
          on:click={saveContent}
        >
          {#if saving}
            Saving...
          {:else}
            Save
          {/if}
        </button>
      </footer>
    </div>
  {/if}
</Drawer>
