<script>
    import Content from './README.md';
    import { onMount } from 'svelte';
    
    let toc = [];
    
    onMount(() => {
        // Get all headings from the content
        const article = document.querySelector('article');
        if (article) {
            const headings = article.querySelectorAll('h1, h2, h3, h4, h5, h6');
            toc = Array.from(headings).map(heading => ({
                id: heading.id,
                text: heading.textContent,
                level: parseInt(heading.tagName.charAt(1))
            }));
        }
    });

    function scrollToSection(id) {
        const element = document.getElementById(id);
        if (element) {
            element.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
    }
</script>

{#if Content}
    <div class="items-center mx-auto grid-cols-[80%_20%] grid gap-8 max-w-7xl relative">
        <article class="prose max-w-3xl flex-1">
            <div class="space-y-4">
                <svelte:component this={Content} />
            </div>
        </article>
        
        <nav class="hidden lg:block w-64 fixed top-24 right-[max(0px,calc(50%-45rem))] max-h-[calc(80vh-5rem)] overflow-y-auto">
            <div class="p-4 bg-card text-card-foreground">
                <h4 class="font-semibold mb-4">On this page</h4>
                <ul class="space-y-2">
                    {#each toc as item}
                        <li style="margin-left: {(item.level - 1) * 1}rem">
                            <a
                                href="#{item.id}"
                                class="text-xs hover:text-primary transition-colors"
                                on:click|preventDefault={() => scrollToSection(item.id)}
                            >
                                {item.text}
                            </a>
                        </li>
                    {/each}
                </ul>
            </div>
        </nav>
    </div>
{:else}
<div class="container py-12">
    <h1 class="mb-8 text-3xl font-bold">Sorry</h1>
    <div class="flex min-h-[400px] items-center justify-center text-center">
        <p class="text-lg font-medium">Nothing found</p>
        <p class="mt-2 text-sm text-muted-foreground">Check back later for new content.</p>
    </div>
</div>

{/if}

<!-- <style>
    :global(h1, h2, h3, h4, h5, h6) {
        scroll-margin-top: 5rem;
    }
</style> -->