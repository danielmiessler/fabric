import * as React from 'react'
import { Textarea } from '@/components/ui/textarea'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { FabricPatternSelect } from '@/components/FabricPatternSelect'
import type { FabricQueryProps } from './fetchFabricQuery'

type Props = { onUpdate: (p: FabricQueryProps) => void }
export function FabricYoutube({ onUpdate }: Props) {
  const [data, setData] = React.useState({
    youtubeUrl: '',
    apiurl: 'api/youtube',
    pattern: 'extract_wisdom',
  })
  const update = (change: Partial<FabricQueryProps>) => {
    const changes = { ...data, ...change }
    setData(changes)
    onUpdate(changes)
  }

  return (
    <>
      <div className="space-y-1">
        <Label htmlFor="ytinput">Youtube video</Label>
        <Input
          title="Youtube URL"
          id="ytinput"
          onChangeCapture={({ currentTarget }) => update({ youtubeUrl: currentTarget.value })}
        />
      </div>
      <div className="space-y-1">
        <FabricPatternSelect onChange={(v) => update({ ...data, pattern: v })} />
      </div>
    </>
  )
}
