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
    aliases: module.metadata?.aliases || [],
    date: module.metadata?.date || new Date().toISOString().split('T')[0],
    description: module.metadata?.description || '',
    tags: module.metadata?.tags || [],
    updated: module.metadata?.updated || new Date().toISOString(),
    author: module.metadata?.author || '',
    lead: module.metadata?.lead || '',
    reference: module.metadata?.reference || '',
    content: module.default
  };
});

export async function getPost(slug: string) {
  const post = posts.find(p => p.slug === slug);
  if (!post) return null;
  
  if (typeof post.content === 'string') {
    const compiled = await compile(post.content);
    post.content = compiled?.code || post.content;
  }
  
  return post;
}
