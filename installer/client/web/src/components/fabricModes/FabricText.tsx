import * as React from 'react';
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { FabricPatternSelect } from "@/components/FabricPatternSelect";
import type { FabricQueryProps } from './fetchFabricQuery';

export function FabricText({ onUpdate }: { onUpdate: (p: FabricQueryProps) => void }) {
  const [data, setData] = React.useState({ query: "", apiurl: 'api/query', pattern: 'ai' })
  const update = (change: Partial<FabricQueryProps>) => {
    setData({ ...data, ...change })
    onUpdate(data)
  }

  return (
    <>
      <div className="space-y-1">
        <Label htmlFor="textinput">Document/Query</Label>
        <Textarea placeholder="Input your query here" id="textinput"
          onChange={({ target }) => update({ query: target.value })} />
      </div>
      <div className="space-y-1">
        <FabricPatternSelect onChange={(v) => update({ pattern: v })} />
      </div>
    </>
  )
}