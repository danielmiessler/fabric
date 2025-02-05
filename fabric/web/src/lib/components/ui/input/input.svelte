<script lang="ts">
  export let value: string = '';
  export let placeholder: string | undefined = undefined;
  export let id: string | undefined = undefined;
  export let disabled: boolean = false;
  export let required: boolean = false;
</script>

<input
  type="text"
  {id}
  bind:value
  {placeholder}
  {disabled}
  {required}
  class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm disabled:cursor-not-allowed disabled:opacity-50"
/>
