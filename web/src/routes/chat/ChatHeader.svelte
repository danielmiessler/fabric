<script lang="ts">
	import { patterns } from "$lib/types/chat/patterns";
	import { patternAPI } from "$lib/types/chat/patterns";
	import { Select } from "$lib/components/ui/select";
	import { onMount } from "svelte";

	let selectedPreset = "";
	
	$: if (selectedPreset) {
		console.log('Pattern selected:', selectedPreset);
		patternAPI.selectPattern(selectedPreset);
	}

	onMount(async () => {
		await patternAPI.loadPatterns();
	});
</script>

<header class="border-b">
	<div class="container mx-auto p-3">
		<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
			<h1 class="text-xl font-semibold">Chat</h1>
			<div class="flex flex-wrap items-center gap-2 w-full sm:w-auto">
				<Select 
					class="w-full font-bold sm:w-[200px]"
					bind:value={selectedPreset}
				> 
				<option value="">Load a pattern...</option>
					{#each $patterns as pattern}
						<option value={pattern.Name}>{pattern.Description}</option>
					{/each}
				</Select>
			</div>
		</div>
	</div>
</header>
