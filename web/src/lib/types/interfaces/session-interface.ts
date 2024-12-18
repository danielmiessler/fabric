import type { Message } from "$lib/types/interfaces/chat-interface";

export interface Session {
    name: string;
    content: Message[];
}
