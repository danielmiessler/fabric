import { api } from './base';
import type { Context } from '$lib/types/interfaces/context-interface';

export const contextAPI = {
  async getAvailable(): Promise<Context[]> {
    const response = await api.fetch<Context[]>('/contexts/names');
    return response.data || [];
  }
}

// TODO: add context element somewhere in the UI