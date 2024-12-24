---
title: mdsvex - svelte in markdown
description: mdsvex is a markdown preprocessor for svelte
date: 2024-12-18
tags: ["svelte", "markdown", "mdsvex"]
---


"They said it was free..."

URL Source: https://mdsvex.pngwn.io/docs

mdsvex is a markdown preprocessor for [Svelte](https://svelte.dev/) components. Basically [MDX](https://mdxjs.com/) for Svelte.

This preprocessor allows you to use Svelte components in your markdown, or markdown in your Svelte components.

mdsvex supports all Svelte syntax and _almost_ all markdown syntax. See [limitations](https://mdsvex.pngwn.io/docs/#limitations) for more information.

You can do this:

```svelte
<script>
        import { Chart } from "../components/Chart.svelte";
</script>

# Here’s a chart

The chart is rendered inside our MDsveX document.

<Chart />
```

It uses [unified](https://unifiedjs.com/), [remark](https://github.com/remarkjs) and [rehype](https://github.com/rehypejs/rehype) and you can use any [remark plugins](https://github.com/remarkjs/remark/blob/main/doc/plugins.md#list-of-plugins) or [rehype plugins](https://github.com/rehypejs/rehype/blob/main/doc/plugins.md#list-of-plugins) to enhance your experience.

[Try it](https://mdsvex.pngwn.io/playground)

Install it
----------

Install it as a dev-dependency.

With `npm`:

```bash
npm i --save-dev mdsvex
```

With `yarn`:

```bash
yarn add --dev mdsvex
```

Use it
------

There are two named exports from `mdsvex` that can be used to transform mdsvex documents, `mdsvex` and `compile`. `mdsvex` is a Svelte preprocessor and is the preferred way to use this library. The `compile` function is useful when you wish to compile mdsvex documents to Svelte components directly, without hooking into the Svelte compiler.

### `mdsvex`

The `mdsvex` preprocessor function is a named import from the `mdsvex` module. Add it as a preprocessor to your rollup or webpack config, and tell the Svelte plugin or loader to also handle `.svx` files.

With rollup and `rollup-plugin-svelte`:

```ts
import { mdsvex } from "mdsvex";

export default {
        ...boring_config_stuff,
        plugins: [
                svelte({
                        // these are the defaults. If you want to add more extensions, see https://mdsvex.pngwn.io/docs#extensions
                        extensions: [".svelte", ".svx"],
                        preprocess: mdsvex()
                })
        ]
};
```

With webpack and `svelte-loader`:

```ts
const { mdsvex } = require('mdsvex')

// add ".svx" to the extensions array
const extensions = ['.mjs', '.js', '.json', '.svelte', '.html', '.svx'];

module.exports = {
        ...boring_config_stuff,
        resolve: { alias, extensions, mainFields },
        module: {
                rules: [
                        {
                                // tell svelte-loader to handle svx files as well
                                test: /.(svelte|html|svx)$/,
                                use: {
                                        loader: 'svelte-loader',
                                        options: {
                                                ...svelte_options,
                                                preprocess: mdsvex()
                                        }
                                }
                        }
                ]
        }
};
```

If you want to use mdsvex without a bundler because you are your own person, then you can use `svelte.preprocess` directly:

```ts
const svelte = require('svelte/compiler');
const { mdsvex } = require('mdsvex');

// This will give you a valid svelte component
const preprocessed = await svelte.preprocess(
        source,
        mdsvex(mdsvex_opts)
);

// Now you can compile it if you wish
const compiled = svelte.compile(
        preprocessed,
        compiler_options
);
```

> If you don’t like the `.svx` file extension, fear not, it is easily customised.

### `compile`

This option performs a very similar task to the preprocessor but it can be used directly, without needing to hook into the Svelte compiler, either directly or via a bundler. The compile option will transform valid mdsvex code into valid svelte code, but it will perform no further actions such as resolving imports.

It supports all of the same options as the preprocessor although the function signature is slightly different. The first argument should be the mdsvex source code you wish to compile, the second argument is an object of options.

```svelte
import { compile } from 'mdsvex';

const transformed_code = await compile(`
<script>
  import Chart from './Chart.svelte';
</script>

# Hello friends

<Chart />
`,
        mdsvexOptions
);
```

In addition to the standard mdsvex options, the options object can also take an optional `filename` property which will be passed to mdsvex. There is no significant advantage to doing this but this provided filename may be used for error reporting in the future. The extension you give to this filename must match one of the extensions provided in the options (defaults to `['.svx']`).

Options
-------

The preprocessor function accepts an object of options, that allow you to customise your experience. The options are global to all parsed files.

```ts
interface MdsvexOptions {
        extensions: string[];
        smartypants: boolean | smartypantsOptions;
        layout: string | { [name: string]: string };
        remarkPlugins: Array<plugin> | Array<[plugin, plugin_options]>;
        rehypePlugins: Array<plugin> | Array<[plugin, plugin_options]>;
        highlight: { highlighter: Function, alias: { [alias]: lang } };
        frontmatter: { parse: Function; marker: string };
}
```

### `extensions`

```ts
extensions: string[] = [".svx"];
```

The `extensions` option allows you to set custom file extensions for files written in mdsvex; the default value is `['.svx']`. Whatever value you choose here must be passed to the `extensions` field of `rollup-plugin-svelte` or `svelte-loader`. If you do not change the default, you must still pass the extension name to the plugin or loader config.

```ts
export default {
        ...config,
        plugins: [
                svelte({
                        extensions: [".svelte", ".custom"],
                        preprocess: mdsvex({
                                extensions: [".custom"]
                        })
                })
        ]
};
```

To import markdown files as components, add `.md` to both the Svelte compiler and `mdsvex` extensions:

```js
// svelte.config.js
import { mdsvex } from 'mdsvex'

export default {
  extensions: ['.svelte', '.svx', '.md'],
  preprocess: mdsvex({ extensions: ['.svx', '.md'] }),
}
```

If you use TypeScript, you should also declare an ambient module:

```ts
declare module '*.md' {
        import type { SvelteComponent } from 'svelte'

        export default class Comp extends SvelteComponent{}

        export const metadata: Record<string, unknown>
}
```

Then you can do:

```svelte
<script>
  import Readme from '../readme.md'
</script>

<Readme />
```


### `smartypants`

```ts
smartypants: boolean | {
        quotes: boolean = true;
        ellipses: boolean = true;
        backticks: boolean | 'all' = true;
        dashes: boolean | 'oldschool' | 'inverted' = true;
} = true;
```

The `smartypants` option transforms ASCII punctuation into fancy typographic punctuation HTML entities.

It turns stuff like:

```
"They said it was free..."
```

into:

> “They said it was free…”

Notice the beautiful punctuation. It does other nice things.

`smartypants` can be either a `boolean` (pass `false` to disable it) or an options object (defaults to `true`). The possible options are as follows.

```
quotes: boolean = true;
```

Converts straight double and single quotes to smart double or single quotes.

*   `"words"` **becomes**: “words”
*   `'words'` **becomes** ‘words’

```
ellipses: boolean = true;
```

Converts triple-dot characters (with or without spaces) into a single Unicode ellipsis character.

*   `words...` **becomes** words…

```
backticks: boolean | 'all' = true;
```

When `true`, converts double back-ticks into an opening double quote, and double straight single quotes into a closing double quote.

*   ` ``words''` **becomes** “words”

When `'all'` it also converts single back-ticks into a single opening quote, and a single straight quote into a closing single, smart quote.

Note: Quotes can not be `true` when backticks is `'all'`;

```
dashes: boolean | 'oldschool' | 'inverted' = true;
```

When `true`, converts two dashes into an em-dash character.

*   `--` **becomes** —

When `'oldschool'`, converts two dashes into an en-dash, and three dashes into an em-dash.

*   `--` **becomes** –
*   `---` **becomes** —

When `'inverted'`, converts two dashes into an em-dash, and three dashes into an en-dash.

*   `--` **becomes** —
*   `---` **becomes** –

### `layout`

```
layout: string | Array<string | RegExp, string>;
```

The `layout` option allows you to provide a custom layout component that will wrap your mdsvex file like so:

```
<Layout>
 <MdsvexDocument />
<Layout>
```

> Layout components receive all frontmatter values as props, which should provide a great deal of flexibility when designing your layouts.

You can provide a `string`, which should be the path to your layout component. An absolute path is preferred but mdsvex tries to resolve relative paths based upon the current working directory.

```
import { join } from "path";

const path_to_layout = join(__dirname, "./src/Layout.svelte");

mdsvex({
        layout: path_to_layout
});
```

In some cases you may want different layouts for different types of document, to address this you may pass an object of named layouts instead. Each key should be a name for your layout, the value should be a path as described above. A fallback layout, or default, can be passed using `_` (underscore) as a key name.

```
mdsvex({
        layout: {
                blog: "./path/to/blog/layout.svelte",
                article: "./path/to/article/layout.svelte",
                _: "./path/to/fallback/layout.svelte"
        }
});
```

```
remarkPlugins: Array<plugin> | Array<[plugin, plugin_options]>;
rehypePlugins: Array<plugin> | Array<[plugin, plugin_options]>;
```

mdsvex has a simple pipeline. Your source file is first parsed into a Markdown AST (MDAST), this is where remark plugins would run. Then it is converted into an HTML AST (HAST), this is where rehype plugins would be run. After this it is converted (stringified) into a valid Svelte component ready to be compiled.

[remark](https://github.com/remarkjs) and [rehype](https://github.com/rehypejs/rehype) have a vibrant plugin ecosystem and mdsvex allows you to pass any [remark plugins](https://github.com/remarkjs/remark/blob/main/doc/plugins.md#list-of-plugins) or [rehype plugins](https://github.com/rehypejs/rehype/blob/main/doc/plugins.md#list-of-plugins) as options, which will run on the remark and rehype ASTs at the correct point in the pipeline.

These options take an array. If you do not wish to pass any options to a plugin then you can simply pass an array of plugins like so:

```
import containers from "remark-containers";
import github from "remark-github";

mdsvex({
        remarkPlugins: [containers, github]
});
```

If you _do_ wish to pass options to your plugins then those array items should be an array of `[plugin, options]`, like so:

```
import containers from "remark-containers";
import github from "remark-github";

mdsvex({
        remarkPlugins: [
                [containers, container_opts],
                [github, github_opts]
        ]
});
```

You can mix and match as needed, only providing an array when options are needed:

```
import containers from "remark-containers";
import github from "remark-github";

mdsvex({
        remarkPlugins: [
                [containers, container_opts],
                github,
                another_plugin,
                [yet_another_plugin, more_options]
        ]
});
```

While these examples use `remarkPlugins`, the `rehypePlugins` option works in exactly the same way. You are free to use one or both of these options as you wish.

Remark plugins work on the Markdown AST (MDAST) produced by remark, rehype plugins work on the HTML AST (HAST) produced by rehype and it is possible to write your own custom plugins if the existing ones do not satisfy your needs!

### `highlight`

```
highlight: {
        highlighter: (code: string, lang: string) => string | Promise<string>
        alias: { [lang : string]: string }
};
```

Without any configuration, mdsvex will automatically highlight the syntax of over 100 languages using [PrismJS](https://prismjs.com/), you simply need to add the language name to the fenced code block and import the CSS file for a Prism theme of your choosing. See [here for available options](https://github.com/PrismJS/prism-themes). Languages are loaded on-demand and cached for later use, this feature does not unnecessarily load all languages for highlighting purposes.

Custom aliases for language names can be defined via the `alias` property of the highlight option. This property takes an object of key-value pairs: the key should be the alias you wish to define, the value should be the language you wish to assign it to.

```
mdsvex({
        highlight: {
                alias: { yavascript: "javascript" }
        }
})
```

If you wish to handle syntax-highlighting yourself, you can provide a custom highlight function via the `highlighter` property. The function will receive two arguments, the `code` to be highlighted and the `lang` defined in the fenced code-block, both are strings. You can use this information to highlight as you wish. The function should return a string of highlighted code.

You can disable syntax highlighting by passing a function that does nothing:

```
function highlighter(code, lang) {
        return `<pre><code>${code}</code></pre>`;
}

mdsvex({
        highlight: {
                highlighter
        }
})
```

### `frontmatter`

```
frontmatter: { parse: Function, marker: string };
```

By default mdsvex supports yaml frontmatter, this is defined by enclosing the YAML in three hyphens (`---`). If you want to use a custom language or marker for frontmatter then you can use the `frontmatter` option.

`frontmatter` should be an object that can contain a `marker` and a `parse` property.

```
marker: string = '-';
```

The marker option defines the fence for your frontmatter. This defaults to `-` which corresponds to the standard triple-hyphen syntax (`---`) that you would normally use to define frontmatter. You can pass in a custom string to change this behaviour:

```
mdsvex({
        frontmatter: {
                marker: "+"
        }
});
```

Now you can use `+++` to mark frontmatter. Setting _only_ the marker will keep the default frontmatter parser which only supports YAML.

```
parse: (frontmatter, message) => Object | undefined
```

The `parse` property accepts a function which allows you to provide a custom parser for frontmatter. This is useful if you want to use a different language in your frontmatter.

The parse function gets the raw frontmatter as the first argument and a `messages` array as the second.

If parsing is successful, the function should return the parsed frontmatter (as an object of key-value pairs), if there is a problem the function should return `undefined` or `false` . Any parsing errors or warnings should be pushed into the `messages` array which will be printed to the console when mdsvex has finished parsing. If you would prefer to throw an error, you are free to do so but it will interrupt the parsing process.

In the following example, we will modify the frontmatter handling so we can write our frontmatter in TOML with a triple-`+` fence.

```
mdsvex({
        marker: "+",
        parse(frontmatter, messages) {
                try {
                        return toml.parse(frontmatter);
                } catch (e) {
                        messages.push(
                                "Parsing error on line " +
                                        e.line +
                                        ", column " +
                                        e.column +
                                        ": " +
                                        e.message
                        );
                }
        }
});
```

Now we will be able to write TOML frontmatter:

```
+++
title = "TOML Example"

[owner]
name = "some name"
dob = 1879-05-27T07:32:00-08:00
+++
```

Layouts
-------

Layouts are one of the more powerful features available in mdsvex and allow for a great deal of flexibility. At their simplest a layout is just a component that wraps an mdsvex document. Providing a string as the layout option will enable this behaviour:

```
mdsvex({
        layout: "./path/to/layout.svelte"
});
```

Layouts receive all values defined in frontmatter as props:

```
<Layout {...props} >
  <!-- mdsvex content here -->
</Layout>
```

You can then use these values in your layout however you wish, a typical use might be to define some fancy formatting for headings, authors, and dates. Although you could do all kinds of wonderful things. You just need to make sure you provide a default `slot` so the mdsvex content can be passed into your layout and rendered.

```
<script>
  export let title;
  export let author;
  export let date;
</script>

<h1>{ title }</h1>
<p class="date">on: { date }</p>
<p class="date">by: { author }</p>
<slot>
  <!-- the mdsvex content will be slotted in here -->
</slot>
```

### Named Layouts

In some cases you may want different layouts for different types of document. To address this you can pass an object of named layouts instead. Each key should be a name for your layout, the value should be the path to that layout file. A fallback layout, or default, can be passed using `_` (underscore) as a key name.

```
mdsvex({
        layout: {
                blog: "./path/to/blog/layout.svelte",
                article: "./path/to/article/layout.svelte",
                _: "./path/to/fallback/layout.svelte"
        }
});
```

If you pass an object of named layouts, you can decide which layout to use on a file-by-file basis by declaring it in the frontmatter. For example, if you wanted to force a document to be wrapped with the `blog` layout you would do the following:

```
---
layout: blog
---
```

If you are using named layouts and do not have a layout field in the frontmatter then mdsvex will try to pick the correct one based on the folder a file is stored in. Take the following folder structure:

```
.
├── blog
│   └── my-blog-post.svx
└── article
    └── my-article.svx
```

If there is a layout named `blog` and `article` then documents in the `blog` folder will use the `blog` layout, articles in the `articles` folder will use the `article` layout. mdsvex will try to check both singular and pluralised names, as you may have named a folder `events` but the matching layout could be named `event`, however, having the same folder and layout name will make this process more reliable. The current working directory is removed from the path when checking for matches but nested folders can still cause problems if there are conflicts. Shallow folder structures and unique folder and layout names will prevent these kinds of collisions.

If there is no matching layout then the fallback layout (`_`) will be applied, if there is no fallback then no layout will be applied.

### disabling layouts

If you are using layouts but wish to disable them for a specific component, then you can set the `layout` field to `false` to prevent the application of a layout.

```
---
layout: false
---
```

### Custom Components

Layouts also allow you to provide custom components to any mdsvex file they are applied to. Custom components replace the elements that markdown would normally generate.

```
# Title

Some text

- a
- short
- list
```

Would normally compile to:

```
<h1>Title</h1>
<p>Some text</p>
<ul>
  <li>a</li>
  <li>short</li>
  <li>list</li>
</ul>
```

Custom components allow you to replace these elements with components. You can define components by exporting named exports from the `context="module"` script of your Layout file:

```
<script context="module">
  import { h1, p, li } from './components.js';
  export { h1, p, li };
</script>
```

The named exports must be named after the actual element you want to replace (`p`, `blockquote`, etc.), the value must be the component you wish to replace them with. This makes certain named exports ‘protected’ API, make sure you don’t use html names as export names for other values. Named exports whose names do not correspond to an HTML element will be ignored, so feel free to continue using them for other purposes as well. As these are named exports it is possible for the bundler to treeshake unused custom components, even if they are exported.

The above custom components would generate:

```
<script>
  import * as Components from './Layout.svelte';
</script>

<Components.h1>Title</Components.h1>
<Components.p>Some text</Components.p>
<ul>
  <Components.li>a</Components.li>
  <Components.li>short</Components.li>
  <Components.li>list</Components.li>
</ul>
```

Notice that the `ul` is left intact: elements are replaced _after_ the markdown is parsed to HTML. This allows greater flexibility, for example, when using custom components to customise lists, tables or other markdown that compiles to a combination of different HTML elements.

You may also receive attributes of the normal HTML component. For example, to render a custom `<img>` tag you could do:

```
<script>
  export let src;
</script>

<img src={src} />
```

Frontmatter
-----------

YAML frontmatter is a common convention in blog posts and mdsvex supports it out of the box. If you want to use a custom language or marker for frontmatter than you can use the [`frontmatter`](https://mdsvex.pngwn.io/docs#frontmatter) option to modify the default behaviour.

Mdsvex integrates well with frontmatter providing additional flexibility when authoring documents.

All variables defined in frontmatter are available directly in the component, exactly as you wrote them:

```
---
title: My lovely article
author: Dr. Fabuloso the Fabulous
---

# {title} by {author}

Some amazing content.
```

Additionally, all of these variables are exported as a single object named `metadata` from the `context="module"` script, so they can easily be imported in javascript:

```
<script context="module">
  export let metadata = {
    title: "My lovely article",
    author: "Dr. Fabuloso the Fabulous"
  };
</script>
```

Due to how `context="module"` scripts work, this metadata can be imported like this:

```
import { metadata } from "./some-mdsvex-file.svx";
```

Frontmatter also interacts with layouts, you can find more details in the [Layout section](https://mdsvex.pngwn.io/docs#layouts).

Integrations
------------

### With shiki

You can use shiki for highlighting rather than prism by leveraging the `highlighter` option:

```ts
import { mdsvex, escapeSvelte } from 'mdsvex';
import { createHighlighter } from 'shiki';

const theme = 'github-dark';
const highlighter = await createHighlighter({
        themes: [theme],
        langs: ['javascript', 'typescript']
});

/** @type {import('mdsvex').MdsvexOptions} */
const mdsvexOptions = {
        highlight: {
                highlighter: async (code, lang = 'text') => {
                        const html = escapeSvelte(highlighter.codeToHtml(code, { lang, theme }));
                        return `{@html \`${html}\` }`;
                }
        },
}
```

Limitations
-----------

### Indentation

In markdown you can begin a code block by indenting 4 spaces. This doesn’t work in mdsvex as indentation is common with XML-based languages. Indenting 4 spaces will do nothing.

In general you have a lot more flexibility when it comes to indenting code in mdsvex than you do in markdown because of the above change, however, you need to be very careful when indenting fenced code blocks. By which I mean, don’t do it.

The following code block will break in a way that is both very bad and quite unexpected:

```js
                ```js
                                        console.log('Hello, World!')
                ```
```

The solution is to not do this. When working with fenced code blocks, do not indent them. This isn’t an issue that can really be worked around, even if the parser did make assumptions about what you meant. Because code blocks are designed to respect whitespace, any fix would simply result in a different but equally frustrating failure. Don’t indent code blocks.