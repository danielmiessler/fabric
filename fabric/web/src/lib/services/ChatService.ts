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
      console.log('\n=== ChatService Request Start ===');
      console.log('1. Request details:', {
        language: get(languageStore),
        pattern: get(selectedPatternName),
        promptCount: request.prompts?.length,
        messageCount: request.messages?.length
      });

      console.log('2. First prompt:', {
        pattern: request.prompts?.[0]?.patternName,
        inputLength: request.prompts?.[0]?.userInput?.length,
        hasLanguageInInput: request.prompts?.[0]?.userInput?.includes(get(languageStore)),
        systemPromptLength: request.prompts?.[0]?.systemPrompt?.length
      });

      console.log('3. Full request:', JSON.stringify(request, null, 2));

      const response = await fetch('/api/chat', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      });

      console.log('4. Response received:', {
        status: response.status,
        ok: response.ok,
        type: response.type
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

      console.log('5. Creating message stream');
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
          console.log('\n=== Stream Processing Start ===');
          while (true) {
            const { done, value } = await reader.read();
            if (done) {
              console.log('Stream complete');
              break;
            }

            buffer += new TextDecoder().decode(value);
            const messages = buffer
              .split('\n\n')
              .filter(msg => msg.startsWith('data: '));

            if (messages.length > 1) {
              buffer = messages.pop() || '';

              for (const msg of messages) {
                try {
                  const response = JSON.parse(msg.slice(6)) as StreamResponse;
                  console.log('Processing chunk:', {
                    contentLength: response.content?.length,
                    format: response.format,
                    type: response.type,
                    hasPattern: !!get(selectedPatternName)
                  });
                  
                  // Clean pattern output if it's a pattern response and ensure markdown format
                  if (get(selectedPatternName)) {
                    response.content = cleanPatternOutput(response.content);
                    response.format = 'markdown';
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
              console.log('Processing final chunk:', {
                contentLength: response.content?.length,
                format: response.format,
                type: response.type,
                hasPattern: !!get(selectedPatternName)
              });
              
              // Clean pattern output if it's a pattern response and ensure markdown format
              if (get(selectedPatternName)) {
                response.content = cleanPatternOutput(response.content);
                response.format = 'markdown';
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

    console.log('\n=== Creating Chat Prompt ===');
    console.log('1. Current state:', {
      language,
      hasLanguageInstruction: language !== 'en',
      instruction: languageInstruction,
      pattern: get(selectedPatternName)
    });

    const prompt = {
      userInput: userInput + languageInstruction,
      systemPrompt: systemPromptText ?? get(systemPrompt),
      model: config.model,
      patternName: get(selectedPatternName)
    };

    console.log('2. Created prompt:', {
      finalInput: prompt.userInput.substring(0, 100) + '...',
      hasLanguageInInput: prompt.userInput.includes(language),
      pattern: prompt.patternName,
      language
    });

    return prompt;
  }

  public async createChatRequest(userInput: string, systemPromptText?: string, isPattern: boolean = false): Promise<ChatRequest> {
    console.log('\n=== Creating Chat Request ===');
    console.log('1. Input:', {
      userInput,
      isPattern,
      language: get(languageStore)
    });

    const prompt = this.createChatPrompt(userInput, systemPromptText);
    const config = get(chatConfig);
    
    // For pattern processing, don't include message history to ensure clean context
    const messages = isPattern ? [] : get(messageStore);

    const request = {
      prompts: [prompt],
      messages: messages,
      ...config
    };

    console.log('2. Final request:', {
      promptCount: request.prompts.length,
      messageCount: request.messages.length,
      firstPromptInput: request.prompts[0].userInput,
      hasLanguageInPrompt: request.prompts[0].userInput.includes(get(languageStore))
    });

    return request;
  }

  public async streamPattern(userInput: string, systemPromptText?: string): Promise<ReadableStream<StreamResponse>> {
    const request = await this.createChatRequest(userInput, systemPromptText, true);
    return this.fetchStream(request);
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
