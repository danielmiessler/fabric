
import { writable } from 'svelte/store';

interface NoteState {
  content: string;
  lastSaved: Date | null;
  isDirty: boolean;
}

function createNoteStore() {
  const { subscribe, set, update } = writable<NoteState>({
    content: '',
    lastSaved: null,
    isDirty: false
  });

  return {
    subscribe,
    updateContent: (content: string) => update(state => ({
      ...state,
      content,
      isDirty: true
    })),
    save: () => update(state => ({
      ...state,
      lastSaved: new Date(),
      isDirty: false
    })),
    reset: () => set({
      content: '',
      lastSaved: null,
      isDirty: false
    })
  };
}

export const noteStore = createNoteStore();
