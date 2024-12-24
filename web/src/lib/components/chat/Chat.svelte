<script lang="ts">
	import ChatInput from "./ChatInput.svelte";
	import ChatMessages from "./ChatMessages.svelte";
	import ModelConfig from "./ModelConfig.svelte";
	import Models from "./Models.svelte";
	import Patterns from "./Patterns.svelte";
  import NoteDrawer from '$lib/components/ui/noteDrawer/NoteDrawer.svelte';
  import { getDrawerStore } from '@skeletonlabs/skeleton';

  import { page } from '$app/stores';
  import { beforeNavigate } from '$app/navigation';

  const drawerStore = getDrawerStore();
  function openDrawer() {
    drawerStore.open({});
  }

  beforeNavigate(() => {
    drawerStore.close();
  });

  $: isVisible = $page.url.pathname.startsWith('/chat');
</script>

<div class="flex-1 mx-auto p-4 min-h-screen">
  <div class="grid grid-cols-1 auto-fit lg:grid-cols-[250px_minmax(250px,_1.5fr)_minmax(250px,_1.5fr)] gap-4 h-[calc(100vh-2rem)]">
    <div class="flex flex-col space-y-1 order-3 lg:order-1">
      <div class="space-y-2 max-w-full">
        <div class="flex flex-col gap-2">
          <Patterns />
          <Models />
          <ModelConfig />
        </div>

        <div class="flex flex-col gap-2">

          {#if isVisible}
            <div class="flex justify-start mt-2">
              <button type="button"
                class="btn btn-sm border variant-filled-primary"
                on:click={openDrawer}
              >Open Drawer
              </button>
            </div>
            <NoteDrawer />
          {/if}
        </div>
      </div>

      <!-- <button class="primary" on:click={openDrawer}>Open Drawer</button> --> 
    </div>
    <div class="flex flex-col space-y-4 order-2 lg:order-2"> 
      <ChatInput />
    </div>
    <div class="flex flex-col border rounded-lg bg-muted/50 p-4 order-1 lg:order-3 max-h-[695px]">
      <ChatMessages />
    </div>
  </div>
</div>
