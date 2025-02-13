import type { StorageEntity } from './storage-interface';

// Interface matching the JSON structure from pattern_descriptions.json
export interface PatternDescription {
  patternName: string;
  description: string;
}

// Interface for storage compatibility - must use uppercase for StorageEntity
export interface Pattern extends StorageEntity {
  Name: string;        // maps to patternName from JSON
  Description: string; // maps to description from JSON
  Pattern: string;     // pattern content from API
}
