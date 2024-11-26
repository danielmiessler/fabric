import { writable } from 'svelte/store';

function createThemeStore() {
  const { subscribe, set, update } = writable<'light' | 'dark'>('dark');

  return {
    subscribe,
    toggleTheme: () => update(theme => {
      const newTheme = theme === 'light' ? 'dark' : 'light';
      if (typeof document !== 'undefined') {
        document.documentElement.classList.toggle('dark', newTheme === 'dark');
      }
      return newTheme;
    }),
    setTheme: (theme: 'light' | 'dark') => {
      set(theme);
      if (typeof document !== 'undefined') {
        document.documentElement.classList.toggle('dark', theme === 'dark');
      }
    }
  };
}

export const theme = createThemeStore();
export const toggleTheme = theme.toggleTheme;