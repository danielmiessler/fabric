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

class LanguageValidator {
    constructor(private targetLanguage: string) {}

    enforceLanguage(content: string): string {
        if (this.targetLanguage === 'en') return content;
        return `[Language: ${this.targetLanguage}]\n${content}`;
    }
}

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
    private validator: LanguageValidator;

    constructor() {
        this.validator = new LanguageValidator(get(languageStore));
    }

    private async fetchStream(request: ChatRequest): Promise<ReadableStream<StreamResponse>> {
        try {
            console.log('\n=== ChatService Request Start ===');
            console.log('1. Request details:', {
                language: get(languageStore),
                pattern: get(selectedPatternName),
                promptCount: request.prompts?.length,
                messageCount: request.messages?.length
            });

            const response = await fetch('/api/chat', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(request),
            });

            if (!response.ok) {
                throw new ChatError(`HTTP error! status: ${response.status}`, 'HTTP_ERROR', { status: response.status });
            }

            const reader = response.body?.getReader();
            if (!reader) {
                throw new ChatError('Response body is null', 'NULL_RESPONSE');
            }

            return this.createMessageStream(reader);
        } catch (error) {
            if (error instanceof ChatError) throw error;
            throw new ChatError('Failed to fetch chat stream', 'FETCH_ERROR', error);
        }
    }

    private cleanPatternOutput(content: string): string {
        content = content.replace(/^# OUTPUT\s*\n/, '');
        content = content.replace(/^\s*\n/, '');
        content = content.replace(/\n\s*$/, '');
        content = content.replace(/^#\s+([A-Z]+):/gm, '$1:');
        content = content.replace(/^#\s+([A-Z]+)\s*$/gm, '$1');
        content = content.trim();
        content = content.replace(/\n{3,}/g, '\n\n');
        return content;
    }

    private createMessageStream(reader: ReadableStreamDefaultReader<Uint8Array>): ReadableStream<StreamResponse> {
        let buffer = '';
        const cleanPatternOutput = this.cleanPatternOutput.bind(this);
        const language = get(languageStore);
        const validator = new LanguageValidator(language);

        const processResponse = (response: StreamResponse) => {
            const pattern = get(selectedPatternName);
            if (pattern) {
                response.content = cleanPatternOutput(response.content);
                response.format = 'markdown';  // Set format for pattern responses
            }
            if (response.type === 'content') {
                response.content = validator.enforceLanguage(response.content);
            }
            return response;
        };

        return new ReadableStream({
            async start(controller) {
                try {
                    while (true) {
                        const { done, value } = await reader.read();
                        if (done) break;

                        buffer += new TextDecoder().decode(value);
                        const messages = buffer.split('\n\n').filter(msg => msg.startsWith('data: '));

                        if (messages.length > 1) {
                            buffer = messages.pop() || '';
                            for (const msg of messages) {
                                try {
                                    let response = JSON.parse(msg.slice(6)) as StreamResponse;
                                    response = processResponse(response);
                                    controller.enqueue(response);
                                } catch (parseError) {
                                    console.error('Error parsing stream message:', parseError);
                                }
                            }
                        }
                    }

                    if (buffer.startsWith('data: ')) {
                        try {
                            let response = JSON.parse(buffer.slice(6)) as StreamResponse;
                            response = processResponse(response);
                            controller.enqueue(response);
                        } catch (parseError) {
                            console.error('Error parsing final message:', parseError);
                        }
                    }
                } catch (error) {
                    controller.error(new ChatError('Stream processing error', 'STREAM_ERROR', error));
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
            ? `You MUST respond in ${language} language. All output must be in ${language}. `
            : '';
        
        const finalSystemPrompt = languageInstruction + (systemPromptText ?? get(systemPrompt));
        
        const finalUserInput = language !== 'en'
            ? `${userInput}\n\nIMPORTANT: Respond in ${language} language only.`
            : userInput;

        return {
            userInput: finalUserInput,
            systemPrompt: finalSystemPrompt,
            model: config.model,
            patternName: get(selectedPatternName)
        };
    }

    public async createChatRequest(userInput: string, systemPromptText?: string): Promise<ChatRequest> {
        const prompt = this.createChatPrompt(userInput, systemPromptText);
        const config = get(chatConfig);
        
        return {
            prompts: [prompt],
            messages: [],
            ...config
        };
    }

    public async streamPattern(userInput: string, systemPromptText?: string): Promise<ReadableStream<StreamResponse>> {
        const request = await this.createChatRequest(userInput, systemPromptText);
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
                    if (lastResponse) onContent(accumulatedContent, lastResponse);
                    break;
                }
                if (value.type === 'error') {
                    throw new ChatError(value.content, 'STREAM_CONTENT_ERROR');
                }
                accumulatedContent += value.content;
                lastResponse = value;
                onContent(accumulatedContent, value);
            }
        } catch (error) {
            onError(error instanceof ChatError ? error : new ChatError('Stream processing error', 'STREAM_ERROR', error));
        } finally {
            reader.releaseLock();
        }
    }
}
