import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Initialize from localStorage if available
const storedLanguage = browser ? localStorage.getItem('selectedLanguage') || '' : '';

// Manages language state
const languageStore = writable<string>(storedLanguage);

// Subscribe to changes and update localStorage
if (browser) {
    languageStore.subscribe(value => {
        localStorage.setItem('selectedLanguage', value);
    });
}

export { languageStore };