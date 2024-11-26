import { writable, derived } from 'svelte/store';
import { modelsApi } from '$lib/types/chat/models';
import { configApi } from '$lib/types/chat/config';
import type { VendorModel } from '$lib/types/interfaces/model-interface';
import type { ModelConfig } from '$lib/types/interfaces/model-interface';

export const modelConfig = writable<ModelConfig>({
  model: [],
  temperature: 0.7,
  maxLength: 2000,
  top_p: 0.9,
  frequency: 1
});

export const availableModels = writable<VendorModel[]>([]);

// Initialize available models
export async function loadAvailableModels() {
  try {
    const models = await modelsApi.getAvailable();
    console.log('Load models:', models);
    const uniqueModels = [...new Map(models.map(model => [model.name, model])).values()];
    availableModels.set(uniqueModels);
  } catch (error) {
    console.error('Client failed to load available models:', error);
    availableModels.set([]);
  }
}

// Initialize config
export async function initializeConfig() {
  try {
    const config = await configApi.get();
    const safeConfig: ModelConfig = {
      ...config,
      model: Array.isArray(config.model) ? config.model : 
             typeof config.model === 'string' ? (config.model as string).split(',') : []
    };
    modelConfig.set(safeConfig);
  } catch (error) {
    console.error('Failed to load config:', error);
  }
}

/* modelConfig.subscribe(async (config) => {
  try {
    const configRecord: Record<string, string> = {
      model: config.model.toString(),
      temperature: config.temperature.toString(),
      maxLength: config.maxLength.toString(),
      top_p: config.top_p.toString(),
      frequency: config.frequency.toString()
    };
    // await configApi.update(configRecord);
  } catch (error) {
    console.error('Failed to update config:', error);
  }
}); */
