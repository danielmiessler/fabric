import * as React from 'react';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { FabricText } from './FabricText';
import { FabricYoutube } from './FabricYoutube';
import { fetchFabricQuery, defaultFabricQueryProps } from './fetchFabricQuery';
import type { FabricQueryProps } from './fetchFabricQuery';

const MODES = [
  { key: 'text', title: 'Document / Query Input', desc: "", component: FabricText },
  { key: 'youtube', title: 'Youtube Transcript', desc: "", component: FabricYoutube }
]

export function ModeSelectTabs() {
  const [payload, setState] = React.useState<FabricQueryProps>(defaultFabricQueryProps)

  const runFabricQuery = async () => {
    const response = await fetchFabricQuery(payload)
    console.log({ payload, response })
  }

  console.log({ payload })
  return (
    <Tabs defaultValue={MODES[0].key}>
      <TabsList className="grid w-full grid-cols-2">
        {MODES.map(({ key, title }) => (
          <TabsTrigger value={key} key={key} className="capitalize">{title}</TabsTrigger>
        ))}
      </TabsList>
      {MODES.map(({ component, key, title, desc }) => (
        <TabsContent value={key} key={`tabcontent-${key}`}>
          <Card>
            <CardHeader>
              <CardTitle>{title}</CardTitle>
              <CardDescription>
                {desc}
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-2">
              {React.createElement(component, { onUpdate: setState })}
            </CardContent>
            <CardFooter>
              <Button onClick={runFabricQuery}>Run Fabric</Button>
            </CardFooter>
          </Card>
        </TabsContent>
      ))}
    </Tabs>
  )
}
