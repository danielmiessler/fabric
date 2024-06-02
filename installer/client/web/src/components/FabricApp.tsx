import * as React from 'react';
import { ModeSelectTabs } from "./fabricModes/ModeSelect";
import type { ExecuteOutput } from '@/lib/execute';
import type { FabricQueryProps } from './fabricModes/fetchFabricQuery';
import { getMemory, saveMemory } from '@/lib/localStorage';
import { nanoid } from 'nanoid';
import { FabricMemory } from './FabricMemory';
import { Separator } from './ui/separator';
import { Button } from './ui/button';


export const FabricApp = () => {
  const [output, setOutput] = React.useState(null as ExecuteOutput | null)
  const [latest, setLatest] = React.useState(nanoid() as string)
  const memory = React.useMemo(() => {
    return getMemory()
  }, [latest])
  const onResult = async (query: FabricQueryProps, response: ExecuteOutput) => {
    setLatest(nanoid())
    memory.push({
      id: latest,
      date: new Date(),
      title: query.query || query.youtubeUrl,
      pattern: query.pattern,
      command: response.command,
      output: response.ok ? response.data : response.error
    })
    saveMemory(memory)
    setOutput(response)
  }
  return (
    <>
      <ModeSelectTabs onResult={onResult} />
      <Separator className="my-4" />
      <FabricMemory memory={memory} />
      <Separator className="my-4" />
      <ClearFabricMemory />
    </>
  )
}

const ClearFabricMemory = () => {
  const clear = () => {
    saveMemory([])
  }
  return <Button onClick={clear} className="block">Clear all Memory</Button>
}