<script lang="ts">
  import ChatInput from "./ChatInput.svelte";
  import ChatMessages from "./ChatMessages.svelte";
  import ModelConfig from "./ModelConfig.svelte";
  import Models from "./Models.svelte";
  import Patterns from "./Patterns.svelte";
</script>

<div class="flex gap-4 p-2 w-full">
  <aside class="w-1/5">
    <div class="flex flex-col gap-2">
      <Patterns />
      <Models />
      <ModelConfig />
    </div>
  </aside>

  <div class="w-1/2">
    <ChatInput />
  </div>

  <div class="w-1/2">
    <ChatMessages />
  </div>
</div>

<style>
    .loading-message {
        animation: flash 1.5s ease-in-out infinite;
    }

    @keyframes flash {
        0% { opacity: 1; }
        50% { opacity: 0.5; }
        100% { opacity: 1; }
    }
</style>
