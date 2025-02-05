import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { YoutubeTranscript } from 'youtube-transcript';

export const POST: RequestHandler = async ({ request }) => {
  try {
    const body = await request.json();
    console.log('\n=== Request Analysis ===');
    console.log('1. Raw request body:', JSON.stringify(body, null, 2));

    // Handle YouTube URL request
    if (body.url) {
      console.log('2. Processing YouTube URL:', body.url);
      
      // Extract video ID
      const match = body.url.match(/(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})/);
      const videoId = match ? match[1] : null;
      
      if (!videoId) {
        return json({ error: 'Invalid YouTube URL' }, { status: 400 });
      }

      console.log('3. Video ID:', videoId);
      const transcriptItems = await YoutubeTranscript.fetchTranscript(videoId);
      const transcript = transcriptItems
        .map(item => item.text)
        .join(' ');

      console.log('4. Transcript length:', transcript.length);
      return json({
        transcript,
        title: videoId
      });
    }

    // Handle pattern execution request
    console.log('\n=== Pattern Request ===');
    console.log('1. Request to fabric backend:', JSON.stringify(body, null, 2));

    // Log important fields
    console.log('2. Key fields:');
    console.log('- Pattern name:', body.prompts?.[0]?.patternName);
    console.log('- User input length:', body.prompts?.[0]?.userInput?.length);
    console.log('- System prompt length:', body.prompts?.[0]?.systemPrompt?.length);
    console.log('- Top level output path:', body.output?.path);
    console.log('- Prompt level output path:', body.prompts?.[0]?.output?.path);

    const fabricResponse = await fetch('http://localhost:8080/api/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body)
    });

    if (!fabricResponse.ok) {
      console.error('3. Fabric API error:', {
        status: fabricResponse.status,
        statusText: fabricResponse.statusText
      });
      throw new Error(`Fabric API error: ${fabricResponse.statusText}`);
    }

    console.log('3. Fabric response status:', fabricResponse.status);

    const stream = fabricResponse.body;
    if (!stream) {
      throw new Error('No response from fabric backend');
    }

    // Create a TransformStream to inspect the data without modifying it
    const transformStream = new TransformStream({
      transform(chunk, controller) {
        const text = new TextDecoder().decode(chunk);
        if (text.startsWith('data: ')) {
          try {
            const data = JSON.parse(text.slice(6));
            console.log('Stream chunk format:', {
              type: data.type,
              format: data.format,
              contentLength: data.content?.length
            });
          } catch (e) {
            console.log('Failed to parse stream chunk:', text);
          }
        }
        controller.enqueue(chunk);
      }
    });

    // Pipe through the transform stream
    const transformedStream = stream.pipeThrough(transformStream);

    // Return the transformed stream
    const response = new Response(transformedStream, {
      headers: {
        'Content-Type': 'text/event-stream',
        'Cache-Control': 'no-cache',
        'Connection': 'keep-alive'
      }
    });

    return response;

  } catch (error) {
    console.error('\n=== Error ===');
    console.error('Type:', error?.constructor?.name);
    console.error('Message:', error instanceof Error ? error.message : String(error));
    console.error('Stack:', error instanceof Error ? error.stack : 'No stack trace');
    return json(
      { error: error instanceof Error ? error.message : 'Failed to process request' },
      { status: 500 }
    );
  }
};
