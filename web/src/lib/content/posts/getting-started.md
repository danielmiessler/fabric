---
title: Getting Started with SvelteKit
date: 2024-11-01
---

# Getting Started with SvelteKit

SvelteKit is a framework for building web applications of all sizes, with a beautiful development experience and flexible filesystem-based routing.

## Why SvelteKit?

- Zero-config setup
- Filesystem-based routing
- Server-side rendering
- Hot module replacement

```bash
npx sv create my-app
cd my-app
npm install
```

**Install SkeletonUI**

```bash
npm i -D @skeletonlabs/skeleton@next @skeletonlabs/skeleton-svelte@next
```

**Configure Tailwind CSS**

```tailwind.config
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