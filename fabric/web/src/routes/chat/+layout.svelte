<script>
  import { disableScrollHandling } from "$app/navigation";
  import { onMount } from "svelte";

  onMount(() => {
    disableScrollHandling();
  });

  function updateLanguage(lang) {
    document.documentElement.lang = lang;
  }
</script>

<div id="page" class="page-wrapper">
  <div class="viewport-container flex h-[calc(100vh-8rem)]">
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
 






<div role="navigation" aria-label="Language selection">
  <label for="language-select">Language</label>
  <select id="language-select" name="language" aria-label="Select language">
    <option value="en">English</option>
    <!-- other language options -->
  </select>
</div>
