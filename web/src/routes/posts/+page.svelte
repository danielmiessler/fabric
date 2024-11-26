<script lang="ts">
	import { formatDistance } from 'date-fns';
	import type { PageData } from './$types';
	import { Paginator } from '@skeletonlabs/skeleton'
	// import Spinner from '$lib/components/ui/spinner/spinner.svelte';

	export let data: PageData;

	$: posts = data.posts;
	let visible: boolean = true;
	let message: string = "No posts found";
</script>

<div class="container py-12">
	<h1 class="mb-4 text-3xl font-bold">Blog Posts</h1>
		<p class="text-sm mb-8 font-small">This blog is maintained in an Obsidian Vault</p>
	{#if posts.length === 0}
	
	{#if !visible}
	<aside class="alert variant-ghost">
		<div>(icon)</div>
		<slot:fragment href="./+error.svelte" />
		<div class="alert-actions">(buttons)</div>
	</aside>
	{/if}
		
	{:else}
	<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
		{#each posts as post}
		<article class="card card-hover group relative rounded-lg border p-6 hover:bg-muted/50">
			<a href="/posts/{post.slug}" class="absolute inset-0">
				<span class="sr-only">View {post.meta.title}</span>
			</a>
			<div class="flex flex-col justify-between space-y-4">
				<div class="space-y-2">
					<h2 class="text-xl font-semibold tracking-tight">{post.meta.title}</h2>
					<p class="text-muted-foreground">{post.meta.description}</p>
				</div>
				<div class="flex items-center space-x-4 text-sm text-muted-foreground">
					<time datetime={post.meta.date}>
						{formatDistance(new Date(post.meta.date), new Date(), { addSuffix: true })}
					</time>
					{#if post.meta.tags.length > 0}
					<span class="text-xs">•</span>
					<div class="flex flex-wrap gap-2">
						{#each post.meta.tags as tag}
						<a
						href="/tags/{tag}"
						class="inline-flex items-center rounded-md border px-2 py-0.5 text-xs font-semibold transition-colors hover:bg-secondary"
						>
						{tag}
					</a>
					{/each}
				</div>
				{/if}
			</div>
		</div>
		
	</article>
	{/each}
<!-- 	<Paginator records={posts} limit={6} buttonClass="btn" /> -->
</div>
	{/if}
</div>