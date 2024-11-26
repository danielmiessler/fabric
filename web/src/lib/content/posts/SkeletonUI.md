---
title: SkeletonUI
tags: 
- svelte
- styling
- skeletonui
date: 2023-01-17
---
SkeletonUI is a comprehensive UI toolkit that integrates seamlessly with SvelteKit and Tailwind CSS, enabling developers to build adaptive and accessible web interfaces efficiently. 

SkeletonUI offers a comprehensive suite of components to enhance your Svelte applications. Below is a categorized list of these components, presented in Svelte syntax:

```svelte
<!-- Layout Components -->
<AppShell />
<AppBar />
<Sidebar />
<Footer />

<!-- Navigation Components -->
<NavMenu />
<Breadcrumbs />
<Tabs />
<Pagination />

<!-- Form Components -->
<Button />
<Input />
<Select />
<Textarea />
<Checkbox />
<Radio />
<Switch />
<Slider />
<FileUpload />

<!-- Data Display Components -->
<Card />
<Avatar />
<Badge />
<Chip />
<Divider />
<Table />
<List />
<Accordion />
<ProgressBar />
<Rating />
<Tag />

<!-- Feedback Components -->
<Alert />
<Modal />
<Toast />
<Popover />
<Tooltip />

<!-- Media Components -->
<Image />
<Video />
<Icon />

<!-- Utility Components -->
<Spinner />
<SkeletonLoader />
<Placeholder />
```

For detailed information on each component, including their properties and usage examples, please refer to the official SkeletonUI documentation.  
___
Below is an expanded cheat sheet to assist you in leveraging SkeletonUI within your SvelteKit projects.

**1\. Installation**

To set up SkeletonUI in a new SvelteKit project, follow these steps:

- **Create a new SvelteKit project**:

```bash
npx sv create my-skeleton-app
cd my-skeleton-app
npm install
```
- **Install SkeletonUI packages**:

```bash
npm install -D @skeletonlabs/skeleton@next @skeletonlabs/skeleton-svelte@next
```
- **Configure Tailwind CSS**:

In your `tailwind.config.js` file, add the following:

```javascript
import { skeleton, contentPath } from '@skeletonlabs/skeleton/plugin';
import * as themes from '@skeletonlabs/skeleton/themes';

export default {
  content: [
    './src/**/*.{html,js,svelte,ts}',
    contentPath(import.meta.url, 'svelte')
  ],
  theme: {
    extend: {},
  },
  plugins: [
    skeleton({
      themes: [themes.cerberus, themes.rose]
    })
  ]
};
```
- **Set the active theme**:

In your `src/app.html`, set the `data-theme` attribute on the `<body>` tag:

```html
<body data-theme="cerberus">
  <!-- Your content -->
</body>
```

For detailed installation instructions, refer to the official SkeletonUI documentation.

**2\. Components**

SkeletonUI offers a variety of pre-designed components to accelerate your development process. Here's how to use some of them:

- **Button**:

```svelte
<script>
  import { Button } from '@skeletonlabs/skeleton-svelte';
</script>

<Button on:click={handleClick}>Click Me</Button>
```
- **Card**:

```svelte
<script>
  import { Card } from '@skeletonlabs/skeleton-svelte';
</script>

<Card>
  <h2>Card Title</h2>
  <p>Card content goes here.</p>
</Card>
```
- **Form Input**:

```svelte
<script>
  import { Input } from '@skeletonlabs/skeleton-svelte';
  let inputValue = '';
</script>

<Input bind:value={inputValue} placeholder="Enter text" />
```

For a comprehensive list of components and their usage, consult the SkeletonUI components documentation.

**3\. Theming**

SkeletonUI's theming system allows for extensive customization:

- **Applying a Theme**:

Set the desired theme in your `tailwind.config.js` and `app.html` as shown in the installation steps above.
- **Switching Themes Dynamically**:

To enable dynamic theme switching, you can modify the `data-theme` attribute programmatically:

```svelte
<script>
  function switchTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
  }
</script>

<button on:click={() => switchTheme('rose')}>Switch to Rose Theme</button>
```
- **Customizing Themes**:

You can create custom themes by defining your own color palettes and settings in the `tailwind.config.js` file.

For more information on theming, refer to the SkeletonUI theming guide.

**4\. Utilities**

SkeletonUI provides several utility functions and actions to enhance your SvelteKit application:

- **Table of Contents**:

Automatically generate a table of contents based on page headings:

```svelte
<script>
  import { TableOfContents, tocCrawler } from '@skeletonlabs/skeleton-svelte';
</script>

<div use:tocCrawler>
  <TableOfContents />
  <!-- Your content with headings -->
</div>
```
- **Transitions and Animations**:

Utilize built-in transitions for smooth animations:

```svelte
<script>
  import { fade } from '@skeletonlabs/skeleton-svelte';
  let visible = true;
</script>

{#if visible}
  <div transition:fade>
    Fading content
  </div>
{/if}
```

For a full list of utilities and their usage, explore the SkeletonUI utilities documentation.

This cheat sheet provides a foundational overview to help you start integrating SkeletonUI into your SvelteKit projects. For more detailed information and advanced features, please refer to the official SkeletonUI documentation.

https://www.skeleton.dev/docs/introduction