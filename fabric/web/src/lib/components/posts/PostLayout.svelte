<script>
  export let aliases = undefined;
  export let date = '';
  export let tags = [];
  export let title = '';
  export let description = '';
  export let author = '';
  export let updated = null;
</script>

<article class="prose prose-slate dark:prose-invert max-w-5xl flex-1">
  {#if aliases}
    <h1 class="inline-block text-4xl font-bold inherit-colors lg:text-5xl">{aliases}</h1>
  {/if}
  <slot />
</article>

<style lang="postcss">
  :global(h1) {
    font-size: 1.875rem;
    line-height: 2.25rem;
    font-family: var(--theme-font-family-heading);
  }

  @media (min-width: 768px) {
    :global(h1) {
      font-size: 3rem;
      line-height: 1;
    }
  }

  :global(h2) {
    font-size: 1.5rem;
    line-height: 2rem;
    font-family: var(--theme-font-family-heading);
  }

  @media (min-width: 768px) {
    :global(h2) {
      font-size: 2.25rem;
      line-height: 2.5rem;
    }
  }

  :global(p) {
    @layer p;
  }

  :global(ul) {
    @layer ul;
  }

  :global(li) {
    @layer li;
  }

  :global(a) {
    @layer a;
  }

  :global(img) {
    @layer img;
  }

  :global(blockquote) {
    @layer blockquote;
  }

  :global(code) {
    @layer code;
  }
</style>
