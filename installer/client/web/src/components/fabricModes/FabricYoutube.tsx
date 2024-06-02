import * as React from 'react';
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { FabricPatternSelect } from "@/components/FabricPatternSelect";
import type { FabricQueryProps } from './fetchFabricQuery';

export function FabricYoutube({ onUpdate }: { onUpdate: (p: FabricQueryProps) => void }) {
  const [data, setData] = React.useState({ youtubeUrl: "", apiurl: 'api/youtube', pattern: 'extract_wisdom' })
  const update = (change: Partial<FabricQueryProps>) => {
    setData({ ...data, ...change })
    onUpdate(data)
  }

  return (
    <>
      <div className="space-y-1">
        <Label htmlFor="ytinput">Youtube video</Label>
        <Input title="Youtube URL" id="ytinput"
          onChangeCapture={({ currentTarget }) => update({ youtubeUrl: currentTarget.value })} />

      </div>
      <div className="space-y-1">
        <FabricPatternSelect onChange={(v) => update({ ...data, pattern: v })} />
      </div>
    </>
  )
}