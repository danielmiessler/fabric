<script lang="ts">
	import { formatDistance } from 'date-fns';
	import type { PageData } from './$types';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';

	export let data: PageData;

	$: ({ content: Content, meta } = data);
	$: formattedDate = formatDistance(new Date(meta.date), new Date(), { addSuffix: true });
</script>

<article class="container max-w-3xl py-6 lg:py-2">
	{#await Content}
		<div class="flex min-h-[400px] items-center justify-center">
			<div class="flex items-center gap-2">
				<Spinner class="h-6 w-6" />
				<span class="text-sm text-muted-foreground">Loading post...</span>
			</div>
		</div>
	{:then Content}
		<!-- <div class="space-y-4">
			<h1 class="inline-block text-4xl font-bold lg:text-5xl">{meta.title}</h1>
			<div class="flex items-center space-x-4 text-sm text-muted-foreground">
				<time datetime={meta.date}>{formattedDate}</time>
				<span class="text-xs">•</span>
				<div class="flex items-center space-x-2">
					{#each meta.tags as tag}
						<a
							href={`/tags/${tag}`}
							class="inline-flex items-center rounded-md border px-2 py-0.5 text-xs font-semibold transition-colors hover:bg-secondary"
						>
							{tag}
						</a>
					{/each}
				</div>
			</div>
		</div> 
		<hr class="my-8" /> -->
		<div class="prose prose-slate dark:prose-invert">
			<svelte:component this={Content} />
		</div>
	{:catch error}
		<div class="flex min-h-[400px] flex-col items-center justify-center text-center">
			<p class="text-lg font-medium">Failed to load post</p>
			<p class="mt-2 text-sm text-muted-foreground">{error.message}</p>
			<a
				href="/posts"
				class="mt-4 inline-flex items-center justify-center rounded-md bg-primary px-8 py-2 text-sm font-medium text-primary-foreground ring-offset-background transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
			>
				Back to Posts
			</a>
		</div>
	{/await}
</article>