<script lang="ts">
	import { cn } from "$lib/types/utils";
	let className = undefined;
	export { className as class };
</script>

<label
	class={cn(
		"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
		className
	)}
	{...$$restProps}
>
	<slot />
</label>