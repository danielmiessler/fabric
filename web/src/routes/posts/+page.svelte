<script lang="ts">
  //import Search from './Search.svelte';
  import type { PageData } from './$types';
	import Card from '$lib/components/ui/cards/card.svelte';
	import { Youtube } from 'svelte-youtube-lite';
	import PostCard from '$lib/components/posts/PostCard.svelte';
	import { InputChip } from '@skeletonlabs/skeleton';
  import Connections from '$lib/components/ui/connections/Connections.svelte';

  let searchQuery = '';
	let selectedTags: string[] = [];
	let allTags: string[] = [];

	export let data: PageData;
	$: posts = data.posts || [];
	
	// Extract all unique tags from posts
	$: {
		const tagSet = new Set<string>();
		posts?.forEach(post => {
			post.metadata?.tags?.forEach(tag => tagSet.add(tag));
		});
		allTags = Array.from(tagSet);
	}

	// Filter posts based on selected tags
	$: filteredPosts = posts?.filter(post => {
		if (selectedTags.length === 0) return true;
		return selectedTags.every(tag => 
			post.metadata?.tags?.some(postTag => postTag.toLowerCase() === tag.toLowerCase())
		);
	}) || [];

	// Filter posts based on search query
	$: searchResults = filteredPosts.filter(post => {
		if (!searchQuery) return true;
		const query = searchQuery.toLowerCase();
		return (
			post.metadata?.title?.toLowerCase().includes(query) ||
			post.metadata?.description?.toLowerCase().includes(query) ||
			post.metadata?.tags?.some(tag => tag.toLowerCase().includes(query))
		);
	});

	function validateTag(value: string): boolean {
		return allTags.some(tag => tag.toLowerCase() === value.toLowerCase());
	}
</script>

<!-- <Search /> -->

<div class="absolute inset-0 -z-10 overflow-hidden h-96">
  <Connections  particleCount={100} particleSize={3} particleSpeed={0.1} connectionDistance={100}/>
</div>

<div class="py-12">
  <h1 class="mb-4 text-3xl font-bold">Blog Posts</h1>
  <p class="text-sm mb-4 font-small">This blog is maintained in an Obsidian Vault</p>

  <div class="mx-auto ml-auto grid grid-cols-1 md:grid-cols-2 gap-4 justify-end">
    <div class="container mx-auto justify-left">
      <img src="https://img.shields.io/github/languages/top/danielmiessler/fabric" alt="Github top language">
      <img src="https://img.shields.io/github/last-commit/danielmiessler/fabric" alt="GitHub last commit">
      <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
      <br>
      <hr class="!border-t-4" />
      <br>
      <h4 class="h4"><b>Leverage Proven Patterns</b></h4>
      <br>
      <Youtube id="UbDyjIIGaxQ" title="Network Chuck Explains fabric" />
      <p>Post your favorite videos.</p>
      <br>

    </div>
    <div class="container mx-auto ml-auto grid grid-cols-1 md:grid-cols-2 gap-4 m-4">

      <div>
        <h4 class=""><b>Find your interests, build your knowledge</b></h4>
      </div>
      <div class="m-auto md:col-start-1">
        <p>Share it with people. Your experience is valuable. Write often and gain a better understanding of yourself.</p>
        <br>
        <p>Use the patterns to help you create posts and templates for future posts. AI can be a powerful tool. There is no right or wrong way to use it.</p>
        <div class="container mx-auto justify-right">
          <button type="button"><a href="/posts/tutorial" class="btn btn-primary">Get Started</a></button>
        </div>
      </div>
      <div class="md:col-start-2">
        <br>
        <Card
          header="Let Your Voice Be Heard"
          imageUrl="/brain.png"
          imageAlt="Blog post header image"
          title="Blogging, Podcasting, Videos, and More."
          content="What will you create?"
          authorName="Your Name Here"
          authorAvatarUrl=""
          link="/"
        />
      </div>
    </div>
    <div class="container mx-auto ml-auto grid grid-cols-1 md:grid-cols-2 gap-4 mt-8">
      <Card
        header="Curate Your Content"
        imageUrl="/electric.png"
        imageAlt="Blog post header image"
        title="Enter a new title here"
        content="What will you share"
        authorName="Your Name Here"
        authorAvatarUrl=""
        link="/"
      />
      <div class="container mx-auto justify-right">
        <blockquote class="blockquote">There are countless use cases for AI. What will you use if for?</blockquote>
      </div>
    </div>
    <div class="container mx-auto ml-auto grid grid-cols-1 md:grid-cols-2 gap-4 justify-end max-h-36 mt-8 pb-8">
      <div class="container mx-auto md:col-start-2 justify-left">
        <hr class="!border-t-4" />
        <br>
        <h4 class="h4">Showcase your interests. Tell people what you've been working on. Create your community.</h4>
      </div>
      <div class="md:col-start-1">
        <!-- This card should be replaced with explainer graphic or text -->
         <Card
header="Explore the Possibilities"
imageUrl="/obsidian-logo.png"
imageAlt="Blog post header image"
title="Enter a new title here"
content="What will you share?"
authorName="Your Name Here"
authorAvatarUrl=""
link="/"
/> 
      </div>
    </div>

  </div>
  <div class="container mx-auto p-12 m-24 justify-right">
    <blockquote class="blockquote">Turn this into a graphic that spans the page</blockquote>
  </div>

  <div class="rounded-tl-container-token m-auto grid grid-cols-1 gap-4 mt-8">
    <div class="mx-auto">something here</div>
    <!-- <Card
header="Curate Your Content"
imageUrl="/electric.png"
imageAlt="Blog post header image"
title="Enter a new title here"
content="What will you share"
authorName="Your Name Here"
authorAvatarUrl=""
link="/"
/> -->
    <!-- <div class="container mx-auto justify-right">
<blockquote class="blockquote">There are countless use cases for AI. What will you use if for?</blockquote>
</div> -->
  </div>
  <div class="container mx-auto justify-center grid mt-8">
    <div class="container mx-auto justify-center">
      <hr class="!border-t-4" />
      <br>
      <h4 class="h4">Showcase your interests. Tell people what you've been working on. Create your community.</h4>
    </div>

    <!--    <Card
header="Explore the Possibilities"
imageUrl=""
imageAlt="Blog post header image"
title="Enter a new title here"
content="What will you share?"
authorName="Your Name Here"
authorAvatarUrl=""
link="/"
/> -->
  </div>
</div>
<div class="container py-12">
  <div class="my-4">
    <InputChip
      name="tags"
      placeholder="Filter by tags..."
      validation={validateTag}
      bind:value={selectedTags}
    />
  </div>
  <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
    {#each searchResults as post}
      <PostCard {post} /> <!-- TODO: Add images to post metadata --> 
    {/each}
  </div>
</div>
