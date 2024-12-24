<script>
  import '../app.postcss';
  import { AppShell, Toast } from '@skeletonlabs/skeleton';
  import Footer from './Footer.svelte';
  import Header from './Header.svelte';
  import { initializeStores } from '@skeletonlabs/skeleton';
  import { getToastStore } from '@skeletonlabs/skeleton';
  import { onMount } from 'svelte';

	// Initialize stores
	initializeStores();

  const toastStore = getToastStore();

	onMount(() => {
		toastStore.trigger({
      message: "👋 Welcome to the site! I'm still working on it...  ",
			background: 'variant-filled-primary',
			timeout: 15000,
			hoverable: true
		});
	});
</script>

<Toast position="t" />

<AppShell class="relative">
  <div class="fixed inset-0 bg-gradient-to-br from-primary-500/20 via-tertiary-500/20 to-secondary-500/20 -z-10"></div>
  <svelte:fragment slot="header">
    <Header />
    <div class="bg-gradient-to-b variant-gradient-primary-tertiary opacity-20 h-2 py-4">
    </div>
  </svelte:fragment>
  <main class="mx-auto p-4">
    <slot />
  </main>
  <svelte:fragment slot="footer">
    <Footer />
  </svelte:fragment>
</AppShell>

<style>
	main {
    /*height: calc( 100vh - 2rem ); /* Adjust based on header/footer height */
		padding: 2rem;
		box-sizing: border-box;
		overflow-y: auto;
	}
</style>
