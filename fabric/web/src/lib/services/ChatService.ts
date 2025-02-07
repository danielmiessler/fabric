import type {
  ChatRequest,
  StreamResponse,
  ChatError as IChatError,
  ChatPrompt
} from '$lib/interfaces/chat-interface';
import { get } from 'svelte/store';
import { modelConfig } from '$lib/store/model-store';
import { systemPrompt, selectedPatternName } from '$lib/store/pattern-store';
import { chatConfig } from '$lib/store/chat-config';
import { messageStore } from '$lib/store/chat-store';
import { languageStore } from '$lib/store/language-store';

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

  private cleanPatternOutput(content: string): string {
    // Remove the initial "# OUTPUT" if present
    content = content.replace(/^# OUTPUT\s*\n/, '');
    
    // Remove markdown code block delimiters
    content = content.replace(/^```markdown\s*\n/, '');
    content = content.replace(/\n```\s*$/, '');

    // Clean up section headings (remove # but keep the text)
    content = content.replace(/^#\s+([A-Z]+):/gm, '$1:');
    content = content.replace(/^#\s+([A-Z]+)\s*$/gm, '$1');

    // Remove extra newlines at start and end
    content = content.trim();

    // Ensure consistent newlines between sections
    content = content.replace(/\n{3,}/g, '\n\n');

    return content;
  }

  private createMessageStream(reader: ReadableStreamDefaultReader<Uint8Array>): ReadableStream<StreamResponse> {
    let buffer = '';
    const cleanPatternOutput = this.cleanPatternOutput.bind(this);

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
                try {
                  const response = JSON.parse(msg.slice(6)) as StreamResponse;
                  console.log('Parsed stream response:', response);
                  
                  // Clean pattern output if it's a pattern response
                  if (get(selectedPatternName)) {
                    response.content = cleanPatternOutput(response.content);
                    response.format = 'plain';
                  }
                  
                  controller.enqueue(response);
                } catch (parseError) {
                  console.error('Error parsing stream message:', parseError);
                  console.log('Problematic message:', msg);
                }
              }
            }
          }

          if (buffer.startsWith('data: ')) {
            try {
              const response = JSON.parse(buffer.slice(6)) as StreamResponse;
              console.log('Parsed final stream response:', response);
              
              // Clean pattern output if it's a pattern response
              if (get(selectedPatternName)) {
                response.content = cleanPatternOutput(response.content);
                response.format = 'plain';
              }
              
              controller.enqueue(response);
            } catch (parseError) {
              console.error('Error parsing final stream message:', parseError);
              console.log('Problematic final message:', buffer);
            }
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
    const language = get(languageStore);
    const languageInstruction = language !== 'en'
      ? `. Please use the language '${language}' for the output.`
      : '';

    return {
      userInput: userInput + languageInstruction,
      systemPrompt: systemPromptText ?? get(systemPrompt),
      model: config.model,
      patternName: get(selectedPatternName)
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
    onContent: (content: string, response?: StreamResponse) => void,
    onError: (error: Error) => void
  ): Promise<void> {
    const reader = stream.getReader();
    let accumulatedContent = '';
    let lastResponse: StreamResponse | undefined;

    try {
      while (true) {
        const { done, value } = await reader.read();

        if (done) {
          // Send the final message with the last known format
          if (lastResponse) {
            onContent(accumulatedContent, lastResponse);
          }
          break;
        }

        if (value.type === 'error') {
          throw new ChatError(value.content, 'STREAM_CONTENT_ERROR');
        }

        accumulatedContent += value.content;
        lastResponse = value;
        console.log('Processing stream chunk:', { 
          type: value.type, 
          format: value.format, 
          contentLength: value.content.length 
        });
        onContent(accumulatedContent, value);
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
