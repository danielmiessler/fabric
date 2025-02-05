import { writable, get } from 'svelte/store';
import { featureFlags } from '../config/features';

export interface ObsidianSettings {
  saveToObsidian: boolean;
  noteName: string;
}

const defaultSettings: ObsidianSettings = {
  saveToObsidian: false,
  noteName: ''
};

// Initialize store with default settings
export const obsidianSettings = writable<ObsidianSettings>(defaultSettings);

// Update settings only if Obsidian integration is enabled
export function updateObsidianSettings(settings: Partial<ObsidianSettings>) {
  const enabled = get(featureFlags).enableObsidianIntegration;
  console.log('Updating Obsidian settings:', settings, 'Integration enabled:', enabled);
  
  if (!enabled) {
    console.log('Obsidian integration disabled, not updating settings');
    return;
  }
  
  obsidianSettings.update(current => {
    const updated = {
      ...current,
      ...settings
    };
    console.log('Updated Obsidian settings:', updated);
    return updated;
  });
}

// Reset settings to default
export function resetObsidianSettings() {
  const enabled = get(featureFlags).enableObsidianIntegration;
  if (!enabled) return;
  
  obsidianSettings.set(defaultSettings);
}

// Helper to get file path
export function getObsidianFilePath(noteName: string): string | undefined {
  const enabled = get(featureFlags).enableObsidianIntegration;
  if (!enabled || !noteName) return undefined;

  return `/Users/jmdb/Documents/GitHub/FABRIC2/fabric/myfiles/Fabric_obsidian/${
    new Date().toISOString().split('T')[0]
  }-${noteName.trim()}.md`;
}