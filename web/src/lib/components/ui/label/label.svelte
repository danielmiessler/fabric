<script lang="ts">
  import { cn } from "$lib/utils/utils";
  let className = undefined;
  export { className as class };
</script>

<label
	class={cn(
		"p-1 text-sm font-bold leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
		className
	)}
	{...$$restProps}
>
	<slot />
</label>
