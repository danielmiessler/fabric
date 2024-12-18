import { writable, get } from 'svelte/store';
import type { ChatRequest, StreamResponse, ChatState, Message } from '$lib/types/interfaces/chat-interface';
import { chatApi } from '$lib/types/chat/chat';
import { modelConfig } from './model-config';
import { systemPrompt } from '$lib/types/chat/patterns';

export const currentSession = writable<string | null>(null);
export const chatState = writable<ChatState>({
    messages: [],
    isStreaming: false
});

export const setSession = (sessionName: string | null) => {
    currentSession.set(sessionName);
    if (!sessionName) {
        clearMessages();
    }
};

export const clearMessages = () => {
    chatState.update(state => ({ ...state, messages: [] }));
};

export const revertLastMessage = () => {
    chatState.update(state => ({
        ...state,
        messages: state.messages.slice(0, -1)
    }));
};

export async function sendMessage(userInput: string, systemPromptText?: string) {
    // Guard against streaming state
    const currentState = get(chatState);
    if (currentState.isStreaming) {
        console.log('Message submission blocked - already streaming');
        return;
    }

    // Update chat state
    chatState.update((state) => ({
        ...state,
        messages: [...state.messages, { role: 'user', content: userInput }],
        isStreaming: true
    }));
    
    try {
        const config = get(modelConfig);
        const sessionName = get(currentSession);
        
        const request: ChatRequest = {
            prompts: [{
                userInput: userInput,
                systemPrompt: systemPromptText || get(systemPrompt),
                model: Array.isArray(config.model) ? config.model.join(',') : config.model,
                vendor: '',
                patternName: '',
            }],
            temperature: config.temperature,
            top_p: config.top_p,
            frequency_penalty: 0,
            presence_penalty: 0
        };

        const stream = await chatApi.streamChat(request);
        const reader = stream.getReader();

        let assistantMessage: Message = {
            role: 'assistant',
            content: ''
        };

        let isCancelled = false;

        while (!isCancelled) {
            const { done, value } = await reader.read();
            if (done) break;

            // Check if we're still streaming before processing
            const currentState = get(chatState);
            if (!currentState.isStreaming) {
                isCancelled = true;
                break;
            }

            const response = value as StreamResponse;
            switch (response.type) {
                case 'content':
                    assistantMessage.content += response.content += `\n`;
                    chatState.update(state => ({
                        ...state,
                        messages: [
                            ...state.messages.slice(0, -1),
                            {...assistantMessage} 
                        ]
                    }));
                    break;
                case 'error':
                    throw new Error(response.content);
                case 'complete':
                    break;
            }
        }

        if (isCancelled) {
            throw new Error('Stream cancelled');
        }

    } catch (error) {
        console.error('Chat error:', error);
        // Only add error message if still streaming
        const currentState = get(chatState);
        if (currentState.isStreaming) {
            chatState.update(state => ({
                ...state,
                messages: [...state.messages, {
                    role: 'assistant',
                    content: `Error: ${error instanceof Error ? error.message : 'Unknown error occurred'}`
                }]
            }));
        }
    } finally {
        chatState.update(state => ({
            ...state,
            isStreaming: false
        }));
    }
}

export type { StreamResponse };