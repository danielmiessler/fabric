// import type { ModelConfig } from '$lib/types/model-types';
import type { StorageEntity } from '$lib/types/interfaces/storage-interface';

interface APIErrorResponse {
  error: string;
}

interface APIResponse<T> {
  data?: T;
  error?: string;
}

// Define and export the base api object
export const api = {
  async fetch<T>(endpoint: string, options: RequestInit = {}): Promise<APIResponse<T>> {
    try {
      const response = await fetch(`/api${endpoint}`, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*',
          ...options.headers,
        },
      });

      if (!response.ok) {
        const errorData = await response.json() as APIErrorResponse;
        return { error: errorData.error || response.statusText };
      }

      const data = await response.json();
      return { data };
    } catch (error) {
      return { error: error instanceof Error ? error.message : 'Unknown error occurred' };
    }
  },

  get: <T>(fetch: typeof window.fetch, endpoint: string) => api.fetch<T>(endpoint),
  
  post: <T>(fetch: typeof window.fetch, endpoint: string, data: unknown) => 
    api.fetch<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    }),
    
  put: <T>(fetch: typeof window.fetch, endpoint: string, data?: unknown) =>
    api.fetch<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    }),
    
  delete: <T>(fetch: typeof window.fetch, endpoint: string) =>
    api.fetch<T>(endpoint, {
      method: 'DELETE',
    }),

  stream: async function*(fetch: typeof window.fetch, endpoint: string, data: unknown): AsyncGenerator<string> {
    const response = await fetch(`/api${endpoint}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const reader = response.body?.getReader();
    if (!reader) throw new Error('Response body is null');

    const decoder = new TextDecoder();
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      yield decoder.decode(value);
    }
  }
};

export function createStorageAPI<T extends StorageEntity>(entityType: string) {
  return {
    // Get a specific entity by name
    async get(name: string): Promise<T> {
      const response = await api.fetch<T>(`/api${entityType}/${name}`);
      if (response.error) throw new Error(response.error);
      return response.data as T;
    },

    // Get all entity names
    async getNames(): Promise<string[]> {
      const response = await api.fetch<string[]>(`/api${entityType}/names`);
      if (response.error) throw new Error(response.error);
      return response.data || [];
    },

    // Delete an entity
    async delete(name: string): Promise<void> {
      const response = await api.fetch(`/api${entityType}/${name}`, {
        method: 'DELETE',
      });
      if (response.error) throw new Error(response.error);
    },

    // Check if an entity exists
    async exists(name: string): Promise<boolean> {
      const response = await api.fetch<boolean>(`/api${entityType}/exists/${name}`);
      if (response.error) throw new Error(response.error);
      return response.data || false;
    },

    // Rename an entity
    async rename(oldName: string, newName: string): Promise<void> {
      const response = await api.fetch(`/api${entityType}/rename/${oldName}/${newName}`, {
        method: 'PUT',
      });
      if (response.error) throw new Error(response.error);
    },

    // Save an entity
    async save(name: string, content: string | object): Promise<void> {
      const body = typeof content === 'string' ? content : JSON.stringify(content);
      const response = await api.fetch(`/api${entityType}/${name}`, {
        method: 'POST',
        body,
      });
      if (response.error) throw new Error(response.error);
    },
  };
}
