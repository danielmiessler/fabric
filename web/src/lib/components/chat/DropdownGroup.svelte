<script lang="ts">
  import Patterns from "./Patterns.svelte";
  import Models from "./Models.svelte";
  import { Select } from "$lib/components/ui/select";
  import { languageStore } from '$lib/store/language-store';

  const languages = [
    { code: '', name: 'Default Language' },
    { code: 'en', name: 'English' },
    { code: 'fr', name: 'French' },
    { code: 'es', name: 'Spanish' },
    { code: 'de', name: 'German' },
    { code: 'zh', name: 'Chinese' },
    { code: 'ja', name: 'Japanese' }
  ];
</script>

<div class="flex flex-col gap-3">
  <div class="w-[50%]">
    <Patterns />
  </div>
  <div class="w-[50%]">
    <Models />
  </div>
  <div class="w-[50%]">
    <Select 
      bind:value={$languageStore}
      class="bg-primary-800/30 border-none hover:bg-primary-800/40 transition-colors"
    >
      {#each languages as lang}
        <option value={lang.code}>{lang.name}</option>
      {/each}
    </Select>
  </div>
</div>
