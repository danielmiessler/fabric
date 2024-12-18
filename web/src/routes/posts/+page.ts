import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
// import type { PostMetadata } from '$lib/types';

export const load: PageLoad = async () => {
	try {
		const postFiles = import.meta.glob('/src/lib/content/posts/*.{md,svx}', { eager: true });
		
		if (Object.keys(postFiles).length === 0) {
			return {
				posts: []
			};
		}

		const posts = Object.entries(postFiles).map(([path, post]: [string, any]) => {
			const slug = path.split('/').pop()?.replace(/\.(md|svx)$/, '');
			return {
				slug,
				meta: {
					title: post.metadata?.title || 'Untitled',
					date: post.metadata?.date || new Date().toISOString(),
					created: post.metadata?.created || new Date().toISOString(),
					description: post.metadata?.description || '',
					tags: post.metadata?.tags || [],
					aliases: post.metadata?.aliases || [],
					lead: post.metadata?.lead || '',
					updated: post.metadata?.updated || new Date().toISOString(),
					author: post.metadata?.author || 'John Connor',
				}
			};
		});

		posts.sort((a, b) => new Date(b.meta.date).getTime() - new Date(a.meta.date).getTime());

		return {
			posts
		};
	} catch (e) {
		console.error('Failed to load posts:', e);
		throw error(500, 'Failed to load posts');
	}
};