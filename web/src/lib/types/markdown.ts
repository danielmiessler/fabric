
// Declare the interface first
export interface Frontmatter {
  title: string;
  description: string;
  date: string;
  tags: string[];
  updated?: string;
  author?: string;
  layout?: string;
  aliases?: string;
  images?: string[];
}

export interface MdsvexCompileData {
  fm: Frontmatter;
  [key: string]: unknown;
}

// Then declare the module for .md files
//declare module '*.md' {
//  import type { SvelteComponent } from 'svelte';
//  import type { PostMetadata } from '$lib/interfaces/post-interface';
//  export const metadata: PostMetadata;
//  export const frontmatter: Frontmatter;
//  const component: SvelteComponent;
//  export default class MarkdownComponent extends SvelteComponent {}
//   export default component;
//}
