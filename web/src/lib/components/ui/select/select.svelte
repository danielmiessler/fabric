<script lang="ts">
	import { cn } from "$lib/types/utils.ts";
	import { ChevronDown } from "lucide-svelte";
	
	export let value: any = undefined;
	export let disabled = false;
	let className = undefined;
	export { className as class };
</script>

<div class="relative">
	<select
		{disabled}
		bind:value
		class={cn(
			"flex h-9 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm shadow-lg ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50 appearance-none",
			className
		)}
		{...$$restProps}
	>
		<slot />
	</select>
	<ChevronDown class="absolute right-3 top-2.5 h-4 w-4 opacity-50 pointer-events-none" />
</div>