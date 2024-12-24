---
title: Getting Started with SvelteKit
description: How to get started with SvelteKit
aliases: SvelteKit for Beginners
date: '2024-11-01'
updated:
tags: 
  - getting-started
  - sveltekit
---

SvelteKit is a framework for building web applications of all sizes, with a beautiful development experience and flexible filesystem-based routing.

## Why SvelteKit?

- Zero-config setup
- Filesystem-based routing
- Server-side rendering
- Hot module replacement

```shell
npx sv create my-app
cd my-app
npm install
```

**Install SkeletonUI**

```
npm i -D @skeletonlabs/skeleton@next @skeletonlabs/skeleton-svelte@next
```

**Configure Tailwind CSS**

``` ts
import type { Config } from 'tailwindcss';

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
            // NOTE: each theme included will be added to your CSS bundle
            themes: [ themes.cerberus, themes.rose ]
        })
    ]
} satisfies Config
```

**Start the dev server**

```bash
npm run dev
```

Read more at https://svelte.dev, https://next.skeleton.dev/docs/get-started/installation/sveltekit, and https://www.skeleton.dev/docs/introduction
