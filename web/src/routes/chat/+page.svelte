<script lang="ts">
    import ChatHeader from './ChatHeader.svelte';
	import ChatInput from "./ChatInput.svelte";
	import ChatMessages from "./ChatMessages.svelte";
	import ModelConfig from "./ModelConfig.svelte";

	//import { Button } from "$lib/components/ui/button";
	//import { setSystemPrompt } from "$lib/store/chat";
	// import { Select } from "$lib/components/ui/select";
	// import { Save, Code, Share2, MoreHorizontal } from 'lucide-svelte';
	// import { onMount } from "svelte";

</script>

<ChatHeader />

<div class="flex-1 container mx-auto p-4">
	<div class="grid grid-cols-1 lg:grid-cols-[minmax(400px,_600px),minmax(400px,_600px),300px] gap-4">
		<div class="space-y-4 order-2 lg:order-1 h-full overflow-y flex-1"> 
			<ChatInput />
		</div>
		<div class="border container max-h-[680px] rounded-lg bg-muted/50 p-4 order-1 lg:order-2 h-full">
			<ChatMessages />
		</div>
		<div class="space-y-2 order-3">
			<ModelConfig />
		</div>
	</div>
</div>

<!-- NOTE : This is for the mobile view 
 <div class="border container max-h-[660px] rounded-lg bg-muted/50 p-4 order-1 lg:order-2 h-full overflow-auto overflow-y-hidden overflow-x-hidden">
	<ChatMessages />
</div> -->
