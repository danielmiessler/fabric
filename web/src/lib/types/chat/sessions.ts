import { createStorageAPI } from './base';
import type { Session } from '$lib/types/interfaces/session-interface';

export const sessionAPI = createStorageAPI<Session>('sessions');
