<script lang="ts">
  import Terminal from './Terminal.svelte';
  import Fabric from './Fabric.svelte';
</script>

<div class="">
  <Terminal />
  <div class="absolute inset-0 -z-10 overflow-hidden h-96">
		<Fabric  />
	</div>
</div>


