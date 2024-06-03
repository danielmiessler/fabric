import * as React from 'react'
import Markdown from 'marked-react'
import type { ExecuteOutput } from '../lib/execute'
import { Button } from './ui/button'
import type { MemoryEntry } from '../lib/localStorage'
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion'
import { nanoid } from 'nanoid'

type Props = { memory: MemoryEntry[] }

export const FabricMemory = ({ memory }: Props) => {
  if (memory.length === 0) {
    return <code className="my-5">NO HISTORY FOUND</code>
  }
  return (
    <div>
      <Accordion type="single" collapsible className="w-full">
        <MemoryItems memory={memory} />
      </Accordion>
    </div>
  )
}

const truncate = (str: string, maxlength: number) => {
  return str.length > maxlength ? str.slice(0, maxlength - 1) + '…' : str
}

const title = (entry: MemoryEntry): string => {
  return truncate(entry.title || 'unknown', 50)
}

const MemoryItems = ({ memory }: Props) => {
  return memory.reverse().map((entry) => (
    <AccordionItem value={entry.id} key={entry.id}>
      <AccordionTrigger>
        {title(entry)} - {entry.pattern}
      </AccordionTrigger>
      <AccordionContent>
        <code className="block bg-indigo-950">{entry.command}</code>
        <code className="block opacity-50">Ran on {entry.date.toString()}</code>
        <div className="my-5 ">
          {entry.output && <Markdown>{entry.output.join('\n')}</Markdown>}
        </div>
      </AccordionContent>
    </AccordionItem>
  ))
}
