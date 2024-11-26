<script lang="ts">
	import { page } from '$app/stores';
	import { Sun, Moon, Menu, X, Github } from 'lucide-svelte';
	import { Avatar } from '@skeletonlabs/skeleton';
	import { fade } from 'svelte/transition';
	import { theme, toggleTheme } from '$lib/store/theme';
  	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import Fabric from './Fabric.svelte'

	let isMenuOpen = false;

	function goToGithub() {
		window.open('https://github.com/danielmiessler/fabric', '_blank');
	}

	function toggleMenu() {
		isMenuOpen = !isMenuOpen;
	}

	$: currentPath = $page.url.pathname;
	$: isDarkMode = $theme === 'dark';

	const navItems = [
		{ href: '/', label: 'Home' },
		{ href: '/posts', label: 'Posts' },
		{ href: '/tags', label: 'Tags' },
		{ href: '/chat', label: 'Chat' },
		//{ href: '/obsidian', label: 'Obsidian' },
		{ href: '/contact', label: 'Contact' },
		{ href: '/about', label: 'About' },
	];

	onMount(() => {
		const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
		theme.setTheme(prefersDark ? 'dark' : 'light');
	});
</script>

<header class="fixed top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
	<div class="container flex h-16 items-center justify-between px-4">
		<div class="flex items-center gap-4">
			<Avatar 
				src="/src/lib/images/fabric-logo.png" 
				width="w-10" 
				rounded="rounded-full" 
				class="border-2 border-primary/20"
			/>
			<a href="/" class="flex items-center">
				<span class="text-lg font-semibold">Fabric</span>
			</a>
		</div>

		<!-- Desktop Navigation -->
		<nav class="hidden flex-1 px-8 md:flex">
			<ul class="flex items-center space-x-8">
				{#each navItems as { href, label }}
					<li>
						<a
							{href}
							class="text-sm font-medium transition-colors hover:text-primary {currentPath === href ? 'text-primary' : 'text-foreground/60'}"
						>
							{label}
						</a>
					</li>
				{/each}
			</ul>
		</nav>
		
		<div class="flex items-center gap-2">
			<button
				on:click={goToGithub}
				class="inline-flex h-9 w-9 items-center justify-center rounded-md border bg-background text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground"
				aria-label="GitHub"
			>
				<Github class="h-4 w-4" />
				<span class="sr-only">GitHub</span>
			</button>

			<button
				on:click={toggleTheme}
				class="inline-flex h-9 w-9 items-center justify-center rounded-md border bg-background text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground"
				aria-label="Toggle theme"
			>
				{#if isDarkMode}
					<Sun class="h-4 w-4" />
				{:else}
					<Moon class="h-4 w-4" />
				{/if}
				<span class="sr-only">Toggle theme</span>
			</button>

			<!-- Mobile Menu Button -->
			<button
				class="inline-flex h-9 w-9 items-center justify-center rounded-md border bg-background text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground md:hidden"
				on:click={toggleMenu}
				aria-expanded={isMenuOpen}
				aria-label="Toggle menu"
			>
				{#if isMenuOpen}
					<X class="h-4 w-4" />
				{:else}
					<Menu class="h-4 w-4" />
				{/if}
			</button>
		</div>
	</div>

	<!-- Mobile Navigation -->
	{#if isMenuOpen}
		<div class="container md:hidden" transition:fade={{ duration: 200 }}>
			<nav class="flex flex-col space-y-4 p-4">
				{#each navItems as { href, label }}
					<a
						{href}
						class="text-base font-medium transition-colors hover:text-primary {currentPath === href ? 'text-primary' : 'text-foreground/60'}"
						on:click={() => (isMenuOpen = false)}
					>
						{label}
					</a>
				{/each}
			</nav>
		</div>
	{/if}
</header>