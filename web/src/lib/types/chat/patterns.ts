import { createStorageAPI } from './base';
import type { Pattern } from '$lib/types/interfaces/pattern-interface';
import { get, writable } from 'svelte/store';

export const patterns = writable<Pattern[]>([]);
export const systemPrompt = writable<string>('');

export const setSystemPrompt = (prompt: string) => {
    console.log('Setting system prompt:', prompt);
    systemPrompt.set(prompt);
    console.log('Current system prompt:', get(systemPrompt));
};

export const patternAPI = {
    ...createStorageAPI<Pattern>('patterns'),

    async loadPatterns() {
        try {
            const response = await fetch(`/api/patterns/names`);
            const data = await response.json();
            console.log("Load Patterns:", data);
            
            // Create an array of promises to fetch all pattern contents
            const patternsPromises = data.map(async (pattern: string) => {
                try {
                    const patternResponse = await fetch(`/api/patterns/${pattern}`);
                    const patternData = await patternResponse.json();
                    return {
                        Name: pattern,
                        Description: pattern.charAt(0).toUpperCase() + pattern.slice(1),
                        Pattern: patternData.Pattern
                    };
                } catch (error) {
                    console.error(`Failed to load pattern ${pattern}:`, error);
                    return {
                        Name: pattern,
                        Description: pattern.charAt(0).toUpperCase() + pattern.slice(1),
                        Pattern: ""
                    };
                }
            });
            
            // Wait for all pattern contents to be fetched
            const loadedPatterns = await Promise.all(patternsPromises);
            console.log("Patterns with content:", loadedPatterns);
            patterns.set(loadedPatterns);
            return loadedPatterns;
        } catch (error) {
            console.error('Failed to load patterns:', error);
            patterns.set([]);
            return [];
        }
    },

    selectPattern(patternName: string) {
        const allPatterns = get(patterns);
        console.log('Selecting pattern:', patternName);
        const selectedPattern = allPatterns.find(p => p.Name === patternName);
        if (selectedPattern) {
            console.log('Found pattern content:', selectedPattern.Pattern);
            setSystemPrompt(selectedPattern.Pattern.trim());
        } else {
            console.log('No pattern found for name:', patternName);
            setSystemPrompt('');
        }
        console.log('System prompt store value after setting:', get(systemPrompt));
    }
};