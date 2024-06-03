'use client'
export const prerender = false

import * as React from 'react'
import { Check, ChevronsUpDown } from 'lucide-react'
import { useQuery, useQueryClient, QueryClient, QueryClientProvider } from '@tanstack/react-query'

import { Button } from '@/components/ui/button'
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Badge } from './ui/badge'
import { cn } from '@/lib/utils'

const queryClient = new QueryClient()

const fetchModelList = async (): Promise<string[]> => {
  const response = await fetch('/api/model_list', { method: 'GET' })
  const body = await response.json()
  if (body.ok) {
    return body.data.filter((item: string) => item !== '' && !item.endsWith(':'))
  }
  return []
}

const setDefaultModel = async (model: string): Promise<string[]> => {
  const response = await fetch('/api/set_model', { method: 'POST', body: JSON.stringify({ model }) })
  const body = await response.json()
  console.log({ defaultModel: body })
  if (body.ok) {
    return body.data
  }
  return []
}

type SelectEvents = { onChange: (v: string) => void }

export function FabricModelSelect() {
  return (
    <QueryClientProvider client={queryClient}>
      <ModelSelectCombo />
    </QueryClientProvider>
  )
}

export function ModelSelectCombo() {
  const [open, setOpen] = React.useState(false)
  const [response, setResponse] = React.useState([] as string[])
  const [value, setValue] = React.useState('')
  const queryClient = useQueryClient()

  const update = async (currentValue: string) => {
    setValue(currentValue === value ? '' : currentValue)
    setOpen(false)
    const resp = await setDefaultModel(currentValue)
    setResponse(resp)
  }
  const { isPending, isError, data, error } = useQuery({
    queryKey: ['patterns'],
    queryFn: fetchModelList,
  })

  if (isPending) {
    return <Badge>Loading...</Badge>
  }

  if (isError) {
    return <Badge>Error: {error.message}</Badge>
  }

  return (
    <>
      <Popover open={open} onOpenChange={setOpen} modal={true}>
        <PopoverTrigger asChild>
          <Button variant="outline" role="combobox" aria-expanded={open} className="w-[500px] justify-between">
            {value || 'Select model...'}
            <ChevronsUpDown className="w-4 h-4 ml-2 opacity-50 shrink-0" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-[500px] p-0">
          <Command>
            <CommandInput placeholder="Search models..." />
            <CommandEmpty>No pattern found.</CommandEmpty>
            <CommandList>
              <CommandGroup>
                {data.map((item) => (
                  <CommandItem key={item} value={item} onSelect={update}>
                    <Check className={cn('mr-2 h-4 w-4', value === item ? 'opacity-100' : 'opacity-0')} />
                    {item}
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
      {response.length ? <Badge>{response}</Badge> : null}
    </>
  )
}
