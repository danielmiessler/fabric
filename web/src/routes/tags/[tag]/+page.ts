import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	const postFiles = import.meta.glob('/src/lib/content/posts/*.{md,svx}', { eager: true });
	
	const posts = Object.entries(postFiles).map(([path, post]: [string, any]) => {
		const slug = path.split('/').pop()?.replace(/\.(md|svx)$/, '');
		return {
			slug,
			meta: {
				title: post.metadata.title,
				date: post.metadata.date,
				description: post.metadata.description,
				tags: post.metadata.tags || []
			}
		};
	});

	const tagPosts = posts.filter((post) => post.meta.tags.includes(params.tag));

	if (tagPosts.length === 0) {
		throw error(404, `Tag "${params.tag}" not found`);
	}

	return {
		tag: params.tag,
		posts: tagPosts
	};
};