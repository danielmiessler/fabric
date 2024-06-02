import * as React from 'react';
import { ModeSelectTabs } from "./fabricModes/ModeSelect";
import type { ExecuteOutput } from '@/lib/execute';
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer"
import { Button } from './ui/button';

type Props = { output: ExecuteOutput }

export const FabricOutput = ({ output }: Props) => {
  return (
    <Drawer>
      <DrawerTrigger asChild>
        <Button variant="outline">Open Response</Button>
      </DrawerTrigger>
      <DrawerContent>
        <div className="mx-auto w-full max-w-sm">
          <DrawerHeader>
            <DrawerTitle>Fabric Response</DrawerTitle>
            <DrawerDescription><code>{output.command}</code></DrawerDescription>
          </DrawerHeader>
          <div className="p-4 pb-0">
            <div className="flex items-center justify-center space-x-2 w-1/2 md:w-full">
              {output.ok
                ? (<Success data={output.data} />)
                : (<Error error={output.error} />)
              }
            </div>
          </div>
          <DrawerFooter>
            <DrawerClose asChild>
              <Button variant="outline">Close</Button>
            </DrawerClose>
          </DrawerFooter>
        </div>
      </DrawerContent>
    </Drawer>
  )
}

const Error = ({ error }: { error: string[] }) => {
  return (
    <pre className="text-xl max-w-prose bg-red-950 text-white">
      {error.join("\n")}
    </pre>
  )
}
const Success = ({ data }: { data: string[] }) => {
  return (
    <pre className="text-xl max-w-prose bg-green-950 text-white">
      {data.join("\n")}
    </pre>
  )
}