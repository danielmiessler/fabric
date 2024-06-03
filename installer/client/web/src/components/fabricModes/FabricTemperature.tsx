'use client'
export const prerender = false

import * as React from 'react'
import { Check, ChevronsUpDown } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Badge } from '../ui/badge'
import { Slider } from '@/components/ui/slider'
import { Label } from '../ui/label'
import { cn } from '@/lib/utils'
import type { FabricQueryProps } from './fetchFabricQuery'

type SelectEvents = { onUpdate: (v: Partial<FabricQueryProps>) => void }

export function FabricTemperature({ onUpdate }: SelectEvents) {
  const [value, setValue] = React.useState(1)
  const update = (currentValue: number[]) => {
    setValue(currentValue[0])
    onUpdate({ temp: currentValue[0] })
  }

  return (
    <>
      <Label htmlFor="tempslider">Temperature</Label>
      <Slider
        id="tempslider"
        className="w-1/2 mx-8"
        defaultValue={[value]}
        onValueChange={update}
        max={1}
        min={0}
        step={0.1}
      />
    </>
  )
}
