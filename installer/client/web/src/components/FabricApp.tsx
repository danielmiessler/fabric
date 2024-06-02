import * as React from 'react';
import { ModeSelectTabs } from "./fabricModes/ModeSelect";
import type { ExecuteOutput } from '@/lib/execute';
import { FabricOutput } from './FabricOutput';


export const FabricApp = () => {
  const [output, setOutput] = React.useState(null as ExecuteOutput | null)
  const onResult = async (response: ExecuteOutput) => {
    setOutput(response)
  }
  return (
    <>
      <ModeSelectTabs onResult={onResult} />
      {output && <FabricOutput output={output} />}
    </>
  )
}
