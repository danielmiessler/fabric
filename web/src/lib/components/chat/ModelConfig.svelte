<script lang="ts">
  import { Label } from "$lib/components/ui/label";
  import { Slider } from "$lib/components/ui/slider";
  import { modelConfig } from "$lib/store/model-store";
  import Transcripts from "./Transcripts.svelte";
  import NoteDrawer from '$lib/components/ui/noteDrawer/NoteDrawer.svelte';
  import { getDrawerStore } from '@skeletonlabs/skeleton';
  import { Button } from '$lib/components/ui/button';
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

<div class="p-2">
  <div class="space-y-1">
    <Label>Maximum Length ({$modelConfig.maxLength})</Label>
    <Slider
      bind:value={$modelConfig.maxLength}
      min={1}
      max={4000}
      step={1}
    />
  </div>

  <div class="space-y-1">
    <Label>Temperature ({$modelConfig.temperature.toFixed(1)})</Label>
    <Slider
      bind:value={$modelConfig.temperature}
      min={0}
      max={2}
      step={0.1}
    />
  </div>

  <div class="space-y-1">
    <Label>Top P ({$modelConfig.top_p.toFixed(2)})</Label>
    <Slider
      bind:value={$modelConfig.top_p}
      min={0}
      max={1}
      step={0.01}
    />
  </div>

  <div class="space-y-1">
    <Label>Frequency Penalty ({$modelConfig.frequency.toFixed(2)})</Label>
    <Slider
      bind:value={$modelConfig.frequency}
      min={0}
      max={1}
      step={0.01}
    />
  </div>

  <div class="space-y-1">
    <Label>Presence Penalty ({$modelConfig.presence.toFixed(2)})</Label>
    <Slider
      bind:value={$modelConfig.presence}
      min={0}
      max={1}
      step={0.01}
    />
  </div>
  <br>
  <div class="space-y-1">
    <Transcripts />
  </div>

  <div class="flex flex-col gap-2">
    {#if isVisible}
      <div class="flex text-inherit justify-start mt-2">
        <Button
          variant="primary"
          class="btn border variant-filled-primary text-align-center"
          on:click={openDrawer}
        >Open Drawer
        </Button>
      </div>
      <NoteDrawer />
    {/if}
  </div>
</div>
