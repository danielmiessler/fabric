// For the Youtube API
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getTranscript } from '$lib/services/transcriptService';

export const POST: RequestHandler = async ({ request }) => {
    try {
        const body = await request.json();
        console.log('Received request body:', body);
        
        const { url } = body;
        if (!url) {
            return json({ error: 'URL is required' }, { status: 400 });
        }

        console.log('Fetching transcript for URL:', url);
        const transcriptData = await getTranscript(url);
        
        console.log('Successfully fetched transcript, preparing response');
        const response = json(transcriptData);
        
        // Log the actual response being sent
        const responseText = JSON.stringify(transcriptData);
        console.log('Sending response (first 200 chars):', responseText.slice(0, 200) + '...');
        
        return response;
    } catch (error) {
        console.error('Server error:', error);
        return json(
            { error: error instanceof Error ? error.message : 'Failed to fetch transcript' },
            { status: 500 }
        );
    }
};