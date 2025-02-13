import { writable, derived, get } from 'svelte/store';
import type { ChatState, Message, StreamResponse } from '$lib/interfaces/chat-interface';
import { ChatService, ChatError } from '$lib/services/ChatService';
import { languageStore } from '$lib/store/language-store';
import { selectedPatternName } from '$lib/store/pattern-store';

// Initialize chat service
const chatService = new ChatService();

// Local storage key for persisting messages
const MESSAGES_STORAGE_KEY = 'chat_messages';

// Load initial messages from local storage
const initialMessages = typeof localStorage !== 'undefined' 
  ? JSON.parse(localStorage.getItem(MESSAGES_STORAGE_KEY) || '[]') 
  : [];

// Separate stores for different concerns
export const messageStore = writable<Message[]>(initialMessages);
export const streamingStore = writable<boolean>(false);
export const errorStore = writable<string | null>(null);
export const currentSession = writable<string | null>(null);

// Subscribe to messageStore changes to persist messages
if (typeof localStorage !== 'undefined') {
  messageStore.subscribe($messages => {
    localStorage.setItem(MESSAGES_STORAGE_KEY, JSON.stringify($messages));
  });
}

// Derived store for chat state
export const chatState = derived(
  [messageStore, streamingStore],
  ([$messages, $streaming]) => ({
    messages: $messages,
    isStreaming: $streaming
  })
);

// Error handling utility
function handleError(error: Error | string) {
  const errorMessage = error instanceof ChatError 
    ? `${error.code}: ${error.message}`
    : error instanceof Error 
      ? error.message 
      : error;

  errorStore.set(errorMessage);
  streamingStore.set(false);
  return errorMessage;
}

export const setSession = (sessionName: string | null) => {
  currentSession.set(sessionName);
  if (!sessionName) {
    clearMessages();
  }
};

export const clearMessages = () => {
  messageStore.set([]);
  errorStore.set(null);
  if (typeof localStorage !== 'undefined') {
    localStorage.removeItem(MESSAGES_STORAGE_KEY);
  }
};

export const revertLastMessage = () => {
  messageStore.update(messages => messages.slice(0, -1));
};

export async function sendMessage(content: string, systemPromptText?: string, isSystem: boolean = false) {
  try {
    console.log('\n=== Message Processing Start ===');
    console.log('1. Initial state:', {
      content: content.substring(0, 100) + '...',
      isSystem,
      hasSystemPrompt: !!systemPromptText,
      currentLanguage: get(languageStore),
      messageCount: get(messageStore).length,
      pattern: get(selectedPatternName)
    });

    const $streaming = get(streamingStore);
    if ($streaming) {
      throw new ChatError('Message submission blocked - already streaming', 'STREAMING_BLOCKED');
    }

    streamingStore.set(true);
    errorStore.set(null);

    // Add message
    messageStore.update(messages => [...messages, { 
      role: isSystem ? 'system' : 'user', 
      content 
    }]);

    console.log('2. Message added to store:', {
      role: isSystem ? 'system' : 'user',
      contentLength: content.length,
      totalMessages: get(messageStore).length,
      language: get(languageStore)
    });

    if (!isSystem) {
      console.log('3. Preparing chat stream:', {
        language: get(languageStore),
        pattern: get(selectedPatternName),
        hasSystemPrompt: !!systemPromptText,
        systemPromptLength: systemPromptText?.length
      });

      const stream = await chatService.streamChat(content, systemPromptText);
      console.log('4. Stream created, beginning processing');

      await chatService.processStream(
        stream,
        (content: string, response?: StreamResponse) => {
          console.log('5. Processing stream chunk:', {
            contentLength: content.length,
            format: response?.format,
            type: response?.type,
            language: get(languageStore)
          });

          messageStore.update(messages => {
            const newMessages = [...messages];
            const lastMessage = newMessages[newMessages.length - 1];

            if (lastMessage?.role === 'assistant') {
              lastMessage.content = content;
              // Always preserve format from response
              lastMessage.format = response?.format || lastMessage.format;
              console.log('6a. Updated existing message:', {
                role: 'assistant',
                contentLength: content.length,
                format: lastMessage.format
              });
            } else {
              // Ensure new messages have format from response
              newMessages.push({
                role: 'assistant',
                content,
                format: response?.format || 'markdown'  // Default to markdown for pattern responses
              });
              console.log('6b. Added new message:', {
                role: 'assistant',
                contentLength: content.length,
                format: response?.format || 'markdown'
              });
            }

            return newMessages;
          });
        },
        (error) => {
          handleError(error);
        }
      );
    }

    streamingStore.set(false);
  } catch (error) {
    if (error instanceof Error) {
      handleError(error);
    } else {
      handleError(String(error));
    }
    throw error;
  }
}

// Re-export types for convenience
export type { ChatState, Message };
