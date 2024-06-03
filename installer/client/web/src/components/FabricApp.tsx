import * as React from 'react'
import { nanoid } from 'nanoid'
import type { ExecuteOutput } from '../lib/execute'
import { getMemory, saveMemory } from '../lib/localStorage'
import type { FabricQueryProps } from './fabricModes/fetchFabricQuery'
import { ModeSelectTabs } from './fabricModes/ModeSelect'
import { FabricMemory } from './FabricMemory'
import { Separator } from './ui/separator'
import { Button } from './ui/button'
import { FabricModelSelect } from './FabricModelSelect'

export const FabricApp = () => {
  const [loading, setLoading] = React.useState(true)

  return (
    <>
      <FabricAppActive onLoaded={() => setLoading(false)} />
    </>
  )
}

type AppProps = { onLoaded: () => void }
export const FabricAppActive = ({ onLoaded }: AppProps) => {
  const [output, setOutput] = React.useState(null as ExecuteOutput | null)
  const [latest, setLatest] = React.useState(nanoid() as string)
  const memory = React.useMemo(() => {
    const mem = getMemory()
    // onLoaded()
    return mem
  }, [latest])

  const onResult = async (query: FabricQueryProps, response: ExecuteOutput) => {
    setLatest(nanoid())
    memory.push({
      id: latest,
      date: new Date(),
      title: query.query || query.youtubeUrl,
      pattern: query.pattern,
      command: response.command,
      output: response.ok ? response.data : response.error,
    })
    saveMemory(memory)
    setOutput(response)
  }
  return (
    <>
      <FabricModelSelect />
      <Separator className="my-4" />
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
  return (
    <Button onClick={clear} className="block">
      Clear all Memory
    </Button>
  )
}
