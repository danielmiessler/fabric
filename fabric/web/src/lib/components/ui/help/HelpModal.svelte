<script lang="ts">
  import { X } from 'lucide-svelte';
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher<{
    close: void;
  }>();

  let modalContent: HTMLDivElement;

  function handleMouseLeave() {
    dispatch('close');
  }
</script>

<div 
  class="bg-primary-800 rounded-lg flex flex-col h-[85vh] w-[600px] shadow-lg"
  on:mouseleave={handleMouseLeave}
>
  <!-- Header -->
  <div class="flex justify-between items-center p-4 border-b border-primary-700/30">
    <h2 class="text-lg font-bold text-muted-foreground">Help</h2>
    <button
      class="text-muted-foreground hover:text-primary-300 transition-colors"
      on:click={() => dispatch('close')}
    >
      ✕
    </button>
  </div>

  <!-- Content -->
  <div 
    class="p-6 flex-1 overflow-y-auto"
    bind:this={modalContent}
  >
    <div class="space-y-4">
      <section>
        <h3 class="text-base font-bold text-primary-300 mb-2">Getting Started</h3>
        <p class="text-sm text-muted-foreground">
          Click on "Pattern Description" at the top to select a pattern or select a pattern directly from the dropdown menu.
        </p>
      </section>

      <section>
        <h3 class="text-base font-bold text-primary-300 mb-2">YouTube & Other website URL</h3>
        <p class="text-sm text-muted-foreground">
          Paste a YouTube URL or an article URL in the message box and the link will be processed automatically. Youtube transcripts will be fetch and website html will be converted to markdown for LLM processing. Make sure to setup your Jina API key.
        </p>
      </section>

      <section>
        <h3 class="text-base font-bold text-primary-300 mb-2">Language Support</h3>
        <p class="text-sm text-muted-foreground">
          Select your preferred language from the dropdown menu. The AI will respond in the selected language.
        </p>
      </section>

      <section>
        <h3 class="text-base font-bold text-primary-300 mb-2">Model Configuration</h3>
        <p class="text-sm text-muted-foreground">
          Adjust model parameters like temperature and max length to control the AI's output. Higher temperature means more creative but potentially less focused responses.
        </p>
      </section>

      <section>
        <h3 class="text-base font-bold text-primary-300 mb-2">I only want a transcript or quick summary</h3>
        <p class="text-sm text-muted-foreground">
          Then don't select a pattern and simply ask in the message box "Give me a transcript" and paste the URL. it's that simple.
        </p>
      </section>

      <section>
        <h3 class="text-base font-bold text-primary-300 mb-2">File Support</h3>
        <p class="text-sm text-muted-foreground">
          You can attach files to your messages using the paperclip icon. The AI will analyze the content of text files.
        </p>
      </section>

      <section>
        <h3 class="text-base font-bold text-primary-300 mb-2">Taking notes and saving to Obsidian</h3>
        <p class="text-sm text-muted-foreground">
          Make sure to setup your files path for those options before using.
        </p>
      </section>
    </div>
  </div>
</div>
