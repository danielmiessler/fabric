<!-- WIP to be included in ChatInput or ChatHeader -->

<!--
<script lang="ts">
  import { onMount } from 'svelte';
  import { sessionAPI } from '$lib/types/chat/sessions';
  import { currentSession } from '$lib/store/chat';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  
  let sessions: string[] = [];
  let newSessionName = '';
  
  onMount(async () => {
    try {
      sessions = await sessionAPI.getNames();
    } catch (error) {
      console.error('Failed to load sessions:', error);
    }
  });
  
  async function createSession() {
    if (!newSessionName.trim()) return;
    
    try {
      await sessionAPI.save(newSessionName, JSON.stringify({
        name: newSessionName,
        messages: []
      }));
      sessions = await sessionAPI.getNames();
      currentSession.set(newSessionName);
      newSessionName = '';
    } catch (error) {
      console.error('Failed to create session:', error);
    }
  }
  
  async function deleteSession(name: string) {
    try {
      await sessionAPI.delete(name);
      sessions = await sessionAPI.getNames();
      if ($currentSession === name) {
        currentSession.set(null);
      }
    } catch (error) {
      console.error('Failed to delete session:', error);
    }
  }
</script>

<div class="space-y-4">
  <div class="flex gap-2">
    <Input bind:value={newSessionName} placeholder="New session name" />
    <Button on:click={createSession}>Create</Button>
  </div>
  
  <div class="space-y-2">
    {#each sessions as session}
      <div class="flex items-center justify-between p-2 border rounded">
        <button 
          class="flex-1 text-left"
          on:click={() => currentSession.set(session)}
        >
          {session}
        </button>
        <Button 
          variant="destructive" 
          size="sm"
          on:click={() => deleteSession(session)}
        >
          Delete
        </Button>
      </div>
    {/each}
  </div>
</div> -->