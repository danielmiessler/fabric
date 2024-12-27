<script lang="ts">
	import ChatInput from "./ChatInput.svelte";
	import ChatMessages from "./ChatMessages.svelte";
	import ModelConfig from "./ModelConfig.svelte";
	import Models from "./Models.svelte";
	import Patterns from "./Patterns.svelte";
  //import NoteDrawer from '$lib/components/ui/noteDrawer/NoteDrawer.svelte';
  //import { getDrawerStore } from '@skeletonlabs/skeleton';
  //import { Button } from '$lib/components/ui/button';
  //import { page } from '$app/stores';
  //import { beforeNavigate } from '$app/navigation';
  //
  //const drawerStore = getDrawerStore();
  //function openDrawer() {
  //  drawerStore.open({});
  //}
  //
  //beforeNavigate(() => {
  //  drawerStore.close();
  //});
  //
  //$: isVisible = $page.url.pathname.startsWith('/chat');
</script>

<div content="width=device-width, height=device-height, initial-scale=1.0">
<div class="h-screen overflow-auto">
  <div class="h-full flex flex-col">
    <div class="flex-1 overflow-hidden mx-auto p-2">
      <!-- <div class="flex flex-col columns-3 m-2 p-1 h-screen"> --> 
      <div class="h-full flex gap-2">
        <aside class="w-1/5 overflow-y-auto">
          <!-- <div class="space-y-2"> --> 
          <!-- <div class="flex flex-col gap-2"> -->
          <Patterns />
          <Models />
          <ModelConfig />
          <!-- </div> -->

          <!-- </div> --> 
        </aside>
        <div class="w-1/2">
          <!-- <div class="flex flex-col space-y-4 order-2 lg:order-2 h-screen"> --> 
          <ChatInput />
        </div>
        <br>
        <div class="w-1/2">
          <!-- <div class="flex flex-col rounded-lg order-1 lg:order-3 h-screen w-full"> --> 
          <ChatMessages />
        </div>
      </div>
      <!-- </div> --> 
    </div>
  </div>
</div>
  </div>
