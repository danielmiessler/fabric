import type { 
  ChatRequest, 
  StreamResponse, 
  ChatError as IChatError,
  ChatPrompt
} from '$lib/interfaces/chat-interface';
import { get } from 'svelte/store';
import { modelConfig } from '$lib/store/model-store';
import { systemPrompt } from '$lib/store/pattern-store';
import { chatConfig } from '$lib/store/chat-config';
import { messageStore } from '$lib/store/chat-store'; // Import messageStore

export class ChatError extends Error implements IChatError {
  constructor(
    message: string,
    public readonly code: string = 'CHAT_ERROR',
    public readonly details?: unknown
  ) {
    super(message);
    this.name = 'ChatError';
  }
}

export class ChatService {
  private async fetchStream(request: ChatRequest): Promise<ReadableStream<StreamResponse>> {
    try {
      const response = await fetch('/api/chat', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        throw new ChatError(
          `HTTP error! status: ${response.status}`,
          'HTTP_ERROR',
          { status: response.status }
        );
      }

      const reader = response.body?.getReader();
      if (!reader) {
        throw new ChatError('Response body is null', 'NULL_RESPONSE');
      }

      return this.createMessageStream(reader);
    } catch (error) {
      if (error instanceof ChatError) {
        throw error;
      }
      throw new ChatError(
        'Failed to fetch chat stream',
        'FETCH_ERROR',
        error
      );
    }
  }

  private createMessageStream(reader: ReadableStreamDefaultReader<Uint8Array>): ReadableStream<StreamResponse> {
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

            if (messages.length > 1) {
              buffer = messages.pop() || '';

              for (const msg of messages) {
                controller.enqueue(JSON.parse(msg.slice(6)) as StreamResponse);
              }
            }
          }

          if (buffer.startsWith('data: ')) {
            controller.enqueue(JSON.parse(buffer.slice(6)) as StreamResponse);
          }
        } catch (error) {
          controller.error(new ChatError(
            'Error processing stream',
            'STREAM_PROCESSING_ERROR',
            error
          ));
        } finally {
          reader.releaseLock();
          controller.close();
        }
      },

      cancel() {
        reader.cancel();
      }
    });
  }

  private createChatPrompt(userInput: string, systemPromptText?: string): ChatPrompt {
    const config = get(modelConfig);
    return {
      userInput,
      systemPrompt: systemPromptText ?? get(systemPrompt),
      model: config.model,
      patternName: ''
    };
  }

  public async createChatRequest(userInput: string, systemPromptText?: string): Promise<ChatRequest> {
    const prompt = this.createChatPrompt(userInput, systemPromptText);
    const config = get(chatConfig);
    const messages = get(messageStore);

    return {
      prompts: [prompt],
      messages: messages,  
      ...config
    };
  }

  public async streamChat(userInput: string, systemPromptText?: string): Promise<ReadableStream<StreamResponse>> {
    const request = await this.createChatRequest(userInput, systemPromptText);
    return this.fetchStream(request);
  }

  public async processStream(
    stream: ReadableStream<StreamResponse>,
    onContent: (content: string) => void,
    onError: (error: Error) => void
  ): Promise<void> {
    const reader = stream.getReader();
    let accumulatedContent = '';

    try {
      while (true) {
        const { done, value } = await reader.read();

        if (done) break;

        if (value.type === 'error') {
          throw new ChatError(value.content, 'STREAM_CONTENT_ERROR');
        }

        accumulatedContent += value.content;
        onContent(accumulatedContent);
      }
    } catch (error) {
      if (error instanceof ChatError) {
        onError(error);
      } else {
        onError(new ChatError(
          'Error processing stream content',
          'STREAM_PROCESSING_ERROR',
          error
        ));
      }
    } finally {
      reader.releaseLock();
    }
  }
}
