import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	try {
		const post = await import(`$lib/content/posts/${params.slug}.md`);
		
		return {
			content: post.default,
			meta: {
				title: post.metadata.title,
				date: post.metadata.date,
				description: post.metadata.description,
				tags: post.metadata.tags || [],
				updated: post.metadata.updated || new Date().toISOString(),
				author: post.metadata.author || '',
				lead: post.metadata.lead || '',
				reference: post.metadata.reference || '',
			}
		};
	} catch (e) {
		console.error(`Failed to load post ${params.slug}:`, e);
		throw error(404, `Could not find post ${params.slug}`);
	}
};