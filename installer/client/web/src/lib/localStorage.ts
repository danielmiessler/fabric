export type MemoryEntry = {
  date: Date,
  id: string,
  pattern: string,
  command: string,
  output: string[] | null,
  title: string,
}

const STORAGE_SLOT = 'fabric_memory'

export const getMemory = (): MemoryEntry[] => {
  const mem = window.localStorage.getItem(STORAGE_SLOT) || "[]"
  return JSON.parse(mem)
}

export const saveMemory = (mem: MemoryEntry[]): void => {
  const str = JSON.stringify(mem)
  window.localStorage.setItem(STORAGE_SLOT, str)
} 