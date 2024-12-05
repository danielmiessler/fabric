<script lang="ts">
  // import Terminal from './Terminal.svelte';
  import Fabric from './Fabric.svelte';
  import Card from '$lib/components/ui/cards/card.svelte';
  import { Youtube } from 'svelte-youtube-lite';

  import { slide } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';

  let augemented = false;
  let cards = false;
  // let showTerminal = false;

  setTimeout(() => {
    augemented = true;
    cards = true;
    // showTerminal = true;
  }, 1000);

</script>

<div class="absolute inset-0 -z-10 overflow-hidden h-96">
  <Fabric  particleCount={100} particleSize={3} particleSpeed={0.1} connectionDistance={100}/>
</div>

<div class="flex flex-col justify-between items-center">
  <h1 class="h1 text-8xl font-bold font-sans mt-8">
    <span class="bg-gradient-to-br to-blue-500 from-cyan-300 bg-clip-text text-transparent box-decoration-clone">fabric</span>
  </h1>
  {#if augemented}
  <div class="py-2 mb-6" transition:slide|local="{{delay: 250, duration: 3000, easing: quintOut }}">
  <h2 class="h2 text-2xl text-center font-extrabold bg-gradient-to-br to-blue-500 from-cyan-300 bg-clip-text text-transparent pb-2">Human Flourishing via AI Augmentation</h2>
    <div class="text-2xl">
      <p class="mt-2 font-bold">Fill in the blanks... </p>
      <i class="mt-2 font-medium justify-end">I believe one of the biggest problems in my world is <code class="code variant-filled-secondary">___________,</code> which is why I'm building/creating/doing <code class="code variant-filled-secondary">___________.</code>
  </div>
</div>
  {/if}
</div>

<br>

{#if cards}
<div transition:slide|local="{{delay: 2500, duration: 4000, easing: quintOut, axis: 'x'}}">
  <div class="container mx-auto ml-auto grid grid-cols-1 md:grid-cols-2 gap-4 justify-end">
    <div class="container mx-auto justify-left">
      <img src="https://img.shields.io/github/languages/top/danielmiessler/fabric" alt="Github top language">
      <img src="https://img.shields.io/github/last-commit/danielmiessler/fabric" alt="GitHub last commit">
      <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
      <br>
      <hr class="!border-t-4" />
      <br>
      <h4 class="h4"><b>Leverage Proven Patterns</b> - Use the extensive library of pre-built patterns, crafted to extract maximum value from AI interactions. From summarization to analysis, code review to creative writing, the patterns help you achieve consistent, high-quality results.<br>
      </h4>
      <br>
      <Youtube id="UbDyjIIGaxQ" title="Network Chuck Explains fabric" />
      <br>

    </div>

    <Card
        header="Let Your Voice Be Heard"
        imageUrl="https://media.beehiiv.com/cdn-cgi/image/fit=scale-down,format=auto,onerror=redirect,quality=80/uploads/asset/file/b40c4a73-3efa-401d-a172-aa545b65f088/image.png?t=1703300882"
        imageAlt="Blog post header image"
        title="99% of Earth's population holds the self-limiting belief that only certain people have worthy thoughts and opinions"
        content="Two ideas that belong in the past: 1) Smoking improves your health, and 2) Thinking and sharing is only for special people"
        authorName="Daniel Miessler"
        authorAvatarUrl="src/lib/images/dan.png"
        link="https://danielmiessler.com/p/blogging-podcasting-gatekeeping-concepts"
    />
  </div>

  <div class="container mx-auto ml-auto grid grid-cols-1 md:grid-cols-2 gap-4 mt-8">
    <Card
        header="Curate Your Content"
        imageUrl="https://media.beehiiv.com/cdn-cgi/image/fit=scale-down,format=auto,onerror=redirect,quality=80/uploads/asset/file/697388b0-0eeb-4854-8fb2-83651227ecc0/Dynamic-Content-Summaries-Miessler-2024.png?t=1720992541"
        imageAlt="Blog post header image"
        title="Dynamic Content Generation (DCG)"
        content="Dynamic Content Generation will change how all media is created and consumed"
        authorName="Daniel Miessler"
        authorAvatarUrl="src/lib/images/dan.png"
        link="https://danielmiessler.com/p/dynamic-content-summaries"
    />
    <div class="container mx-auto justify-right">
      <blockquote class="blockquote">There are countless use cases for AI. What will you use if for?</blockquote>
    </div>
  </div>
  <div class="container mx-auto ml-auto grid grid-cols-1 md:grid-cols-2 gap-4 justify-end mt-8">
    <div class="container mx-auto justify-left">
      <hr class="!border-t-4" />
      <br>
      <h4 class="h4">Showcase your interests. Tell people what you've been working on. Create your community.</h4>
    </div>

    <Card
        header="Explore the Possibilities"
        imageUrl="https://public-files.gumroad.com/a1xibdqveljasnxffudsqdp7bqpo"
        imageAlt="Blog post header image"
        title="Red Blue Purple AI - December 2024"
        content="A new course from Jason Haddix!"
        authorName="Jason Haddix"
        authorAvatarUrl="src/lib/images/haddix.jpeg"
        link="https://arcanuminfosec.gumroad.com/l/ygmlpe"
    />
  </div>
</div>
{/if}


