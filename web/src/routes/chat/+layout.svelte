<script>
	import { initializeStores } from "@skeletonlabs/skeleton";
	initializeStores();
</script>

<div class="chat-layout">
	<slot />
</div>

<style>
	.chat-layout {
		display: flex;
		flex-direction: column;
		/* height: 100vh;
		max-height: 100vh;
		overflow: hidden; */
	}
</style>
