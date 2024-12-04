<script>
	import { initializeStores } from "@skeletonlabs/skeleton";
	initializeStores();
</script>

<div class="chat-layout">
	<slot />
</div>

<style>
	.chat-layout {
		flex-direction: column;
	}
</style>
