export interface ChatRequest {
    prompts: {
        userInput: string;
        systemPrompt: string;
        model: string;
        vendor: string;
        // contextName: string;
        patternName: string;
        // sessionName: string;
    }[];
    temperature: number;
    top_p: number;
    frequency_penalty: number;
    presence_penalty: number;
}
  
export interface Message {
    role: 'system' | 'user' | 'assistant';
    content: string;
}

export interface ChatState {
    messages: Message[];
    isStreaming: boolean;
}
  
  export interface StreamResponse {
    type: 'content' | 'error' | 'complete';
    format: 'markdown' | 'mermaid' | 'plain';
    content: string;
}