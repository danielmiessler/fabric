import * as React from 'react'
import { Textarea } from '@/components/ui/textarea'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { FabricPatternSelect } from '@/components/FabricPatternSelect'
import type { FabricQueryProps } from './fetchFabricQuery'

type Props = { onUpdate: (p: FabricQueryProps) => void }

export function FabricText({ onUpdate }: Props) {
  const [data, setData] = React.useState({
    query: '',
    apiurl: 'api/query',
    pattern: 'ai',
  })
  const update = (change: Partial<FabricQueryProps>) => {
    const changes = { ...data, ...change }
    setData(changes)
    onUpdate(changes)
  }

  return (
    <>
      <div className="space-y-1">
        <Label htmlFor="textinput">Document/Query</Label>
        <Textarea
          placeholder="Input your query here"
          id="textinput"
          onChangeCapture={({ currentTarget }) => update({ query: currentTarget.value })}
        />
      </div>
      <div className="space-y-1">
        <FabricPatternSelect onChange={(v) => update({ pattern: v })} />
      </div>
    </>
  )
}
