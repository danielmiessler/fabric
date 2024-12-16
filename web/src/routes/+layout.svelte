<script>
  import '../app.postcss';
  import { AppShell } from '@skeletonlabs/skeleton';
  import ToastContainer from '$lib/components/ui/toast/ToastContainer.svelte';
  import Footer from './Footer.svelte';
  import Header from './Header.svelte';
  import { initializeStores } from '@skeletonlabs/skeleton';
  import { page } from '$app/stores';
  import { fly } from 'svelte/transition';
  import { getToastStore } from '@skeletonlabs/skeleton';
  import { onMount } from 'svelte';

	// Initialize stores
	initializeStores();

  const toastStore = getToastStore();

	onMount(() => {
		toastStore.trigger({
      message: "👋 Welcome to the site! Tell people about yourself and what you do.",
			background: 'variant-filled-primary',
			timeout: 3333,
			hoverable: true
		});
	});
</script>

<ToastContainer />

{#key $page.url.pathname}
<AppShell class="relative">
  <div class="fixed inset-0 bg-gradient-to-br from-primary-500/20 via-tertiary-500/20 to-secondary-500/20 -z-10"></div>
  <svelte:fragment slot="header">
    <Header />
    
    <div class="h-2 py-4">
      </svelte:fragment>
        <div 
        in:fly={{ duration: 500, delay: 100, y: 100 }}
        >
      <main class="mx-auto p-4">
        <slot />
      </main>
    </div>
  
  <svelte:fragment slot="footer">
    <Footer />
  </svelte:fragment>
</AppShell>
{/key}

<style>
	main {
		padding: 2rem;
		box-sizing: border-box;
		overflow-y: auto;
	}
</style>
