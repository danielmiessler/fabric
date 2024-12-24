// import { api } from '$lib/types/base';
import type { ChatRequest, StreamResponse } from '$lib/types/interfaces/chat-interface';

// Create a chat API client
export const chatApi = {
  async streamChat(request: ChatRequest): Promise<ReadableStream<StreamResponse>> {
    const response = await fetch('/api/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const reader = response.body?.getReader();
    if (!reader) throw new Error('Response body is null');

    let buffer = '';
    
    return new ReadableStream({
      async start(controller) {
        try {
          while (true) {
            const { done, value } = await reader.read();
            if (done) break;
          
            buffer += new TextDecoder().decode(value);
            const messages = buffer
              .split('\n\n')
              .filter(msg => msg.startsWith('data: '));
          
            // Process complete messages
            if (messages.length > 1) {
              // Keep the last (potentially incomplete) chunk in the buffer
              buffer = messages.pop() || '';
            
              for (const msg of messages) {
                controller.enqueue(JSON.parse(msg.slice(6)) as StreamResponse);
              }
            }
          }
        } catch (error) {
          controller.error(error);
        } finally {
          // Process any remaining complete messages in the buffer
          if (buffer.startsWith('data: ')) {
            controller.enqueue(JSON.parse(buffer.slice(6)) as StreamResponse);
          }
          controller.close();
          reader.releaseLock();
        }
      },
    });
  },
};