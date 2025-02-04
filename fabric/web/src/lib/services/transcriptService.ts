export interface TranscriptResponse {
  transcript: string;
  title: string;
}

export async function getTranscript(url: string): Promise<TranscriptResponse> {
  try {
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

    return data;
  } catch (error) {
    console.error('Transcript fetch error:', error);
    throw error instanceof Error ? error : new Error('Failed to fetch transcript');
  }
}
