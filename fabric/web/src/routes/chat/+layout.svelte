<script>
  import { disableScrollHandling } from "$app/navigation";
  import { languageStore } from '$lib/store/language-store';
  import { browser } from '$app/environment';
  import LanguageDisplay from '$lib/components/LanguageDisplay.svelte';
  import { onMount } from "svelte";

  onMount(() => {
    disableScrollHandling();
  });

  // Reactive statement for lang attribute
  $: if (browser) {
    document.documentElement.lang = $languageStore;
  }
</script>

<div id="page" class="page-wrapper">
  <div class="viewport-container flex h-[calc(100vh-8rem)]">
    <LanguageDisplay />
    <slot />
  </div>
</div>

<style>
  /* Container that enforces viewport bounds */
  .viewport-container {
    width: 100vw;  /* Full viewport width */
    overflow: hidden; /* Prevent scrolling */
    position: fixed; /* Fix position to viewport */
    left: 0;
  }

  /* Ensure the wrapper doesn't introduce scrolling */
  .page-wrapper :global(#page) {
    display: block;
    flex: none;
    overflow: hidden;
  }

  :global(#page-content) {
    flex: none;
    overflow: hidden;
  }

  /* Ensure any nested content doesn't cause scrolling */
  :global(.viewport-container *) {
    overflow: hidden;
  }
</style>
