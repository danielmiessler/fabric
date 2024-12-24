import { YoutubeTranscript } from 'youtube-transcript';

export interface TranscriptResponse {
    transcript: string;
    title: string;
}

export async function getTranscript(url: string): Promise<TranscriptResponse> {
    try {
        const videoId = extractVideoId(url);
        if (!videoId) {
            throw new Error('Invalid YouTube URL');
        }

        const transcriptItems = await YoutubeTranscript.fetchTranscript(videoId);
        const transcript = transcriptItems
            .map(item => item.text)
            .join(' ');

        const transcriptTitle = transcriptItems
            .map(item => item.text)
            .join('');

        // TODO: Add title fetching 
        return {
            transcript,
            title: videoId // Just returning the video ID as title
        };
    } catch (error) {
        console.error('Transcript fetch error:', error);
        throw new Error('Failed to fetch transcript');
    }
}

function extractVideoId(url: string): string | null {
    const match = url.match(/(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})/);
    return match ? match[1] : null;
}