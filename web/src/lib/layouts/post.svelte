<script lang="ts">
	import { formatDistance } from 'date-fns';
	import TagList from '$components/ui/tag-list/TagList.svelte';

	/** @type {string} */
	export let title;
	/** @type {string} */
	export let date;
	/** @type {string} */
	export let description;
	/** @type {string} */
	export let tags = [];
	/** @type {string}*/
	export let updated;
	/** @type {string}**/
	export let reference;
</script>

<article class="prose prose-slate mx-auto max-w-3xl dark:prose-invert py-12">
	<header class="mb-8 not-prose">
		{#if title}
		<h1 class="mb-2 text-4xl font-bold">{title}</h1>
		{/if}
		{#if description}
			<p class="mb-4 text-lg text-muted-foreground">{description}</p>
		{/if}
		{#if date}
		<div class="flex items-center space-x-4 text-sm text-muted-foreground">
			<time datetime={date}>{formatDistance(new Date(date), new Date(), { addSuffix: true })}</time>
			
			{#if tags?.length}
			<span class="text-xs">•</span>
			<TagList {tags} className="flex-1" />
			{/if}
			{#if updated}
			<span class="text-xs">•</span>
			<time datetime={updated}>Updated {formatDistance(new Date(updated), new Date(), { addSuffix: true })}</time>
			{/if}
			{#if reference}
			<span class="text-xs">•</span>
			<a href={reference}>Reference</a>
			{/if}
		</div>
		{/if}
	</header>

	<div class="mt-8">
		<slot />
	</div>
</article>