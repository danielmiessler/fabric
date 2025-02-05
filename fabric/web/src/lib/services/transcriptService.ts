export interface TranscriptResponse {
  transcript: string;
  title: string;
}

function decodeHtmlEntities(text: string): string {
  const textarea = document.createElement('textarea');
  textarea.innerHTML = text;
  return textarea.value;
}

export async function getTranscript(url: string): Promise<TranscriptResponse> {
  try {
    console.log('\n=== Getting Transcript ===');
    console.log('1. URL:', url);

    const response = await fetch('/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ url })
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    if (data.error) {
      throw new Error(data.error);
    }

    // Decode HTML entities in transcript
    data.transcript = decodeHtmlEntities(data.transcript);

    console.log('2. Got transcript, length:', data.transcript.length);
    console.log('3. First 100 chars:', data.transcript.substring(0, 100));
    return data;
  } catch (error) {
    console.error('Transcript fetch error:', error);
    throw error instanceof Error ? error : new Error('Failed to fetch transcript');
  }
}
