import { compile } from 'mdsvex';

export interface Post {
  slug: string;
  title: string;
  date: string;
  content?: any;
}

const modules = import.meta.glob('../content/posts/*.md' + '../../routes/**/*.md', { eager: true });

export const posts: Post[] = Object.entries(modules).map(([path, module]: [string, any]) => {
  const slug = path.split('/').pop()?.replace('.md', '') || '';
  return {
    slug,
    title: module.metadata?.title || slug,
    date: module.metadata?.date || new Date().toISOString().split('T')[0],
    content: module.default
  };
});

export async function getPost(slug: string) {
  return posts.find(p => p.slug === slug) || null;
}