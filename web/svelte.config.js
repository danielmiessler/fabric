import adapter from '@sveltejs/adapter-auto';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import { mdsvex } from 'mdsvex';
import rehypeSlug from 'rehype-slug';
import rehypeAutolinkHeadings from 'rehype-autolink-headings';
import rehypeExternalLinks from 'rehype-external-links';
import rehypeUnwrapImages from 'rehype-unwrap-images';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

/** @type {import('mdsvex').MdsvexOptions} */
const mdsvexOptions = {
  extensions: ['.md', '.svx'],
  layout: {
    _: join(__dirname, './src/lib/layouts/post.svelte')
  },
  highlight: {
    theme: {
      dark: 'github-dark',
      light: 'github-light'
    }
  },
rehypePlugins: [
	rehypeSlug,
	rehypeUnwrapImages,
    [rehypeAutolinkHeadings, {
      behavior: 'wrap'
    }],
    [rehypeExternalLinks, {
      target: '_blank',
      rel: ['noopener', 'noreferrer']
    }]
  ]
};

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
	  adapter: adapter(),
	  alias: {
		$components: join(__dirname, 'src/lib/components'),
		$lib: join(__dirname, 'src/lib'),
		$styles: join(__dirname, 'src/styles'),
    
		$utils: join(__dirname, 'src/lib/utils')
	  }
	},
	extensions: ['.svelte', '.md', '.svx'],
	preprocess: [
	  vitePreprocess(),
	  mdsvex(mdsvexOptions)
	]
  };
  
  export default config;