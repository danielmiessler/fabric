import type { ModelConfig } from '$lib/types/interfaces/model-interface';
import { api } from './base';

export const configApi = {
    async get(): Promise<ModelConfig> {
      const response = await api.fetch<ModelConfig>('/config');
      if (response.error) throw new Error(response.error);
      return response.data || {
        model: [],
        temperature: 0.7,
        top_p: 0.9,
        frequency: 1,
        maxLength: 2000
      };
    },
  
    /* async update(config: Record<string, string>) {
      const response = await api.fetch('config/update', {
        method: 'POST',
        body: JSON.stringify(config),
      });
      if (response.error) throw new Error(response.error);
      return response;
    } */
};
