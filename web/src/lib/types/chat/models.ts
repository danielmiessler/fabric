import { api } from './base';
import type { VendorModel, ModelsResponse } from '$lib/types/interfaces/model-interface';

export const modelsApi = {
  async getAvailable(): Promise<VendorModel[]> {
    const response = await api.fetch<ModelsResponse>('/models/names');
    console.log("Client raw API response:", response)

    if (response.error) {
      console.error("Client couldn't fetch models:", response.error);
      throw new Error(response.error);
    }
    
    if (!response.data) {
      console.error('No data received from models API');
      return [];
    }
    
    const vendorsData = response.data.vendors || {};
    const result: VendorModel[] = [];
    
    for (const [vendor, models] of Object.entries(vendorsData)) {
      for (const model of models) {
        result.push({
          name: model,
          vendor: vendor
        });
      }
    }
    
    console.log('Available models:', result);
    return result;
  },

/*   async getNames(): Promise<string[]> {
    const response = await api.fetch<ModelsResponse>('/models/names');
    if (response.error) throw new Error(response.error);
    return response.data?.models || [];
  } */
};