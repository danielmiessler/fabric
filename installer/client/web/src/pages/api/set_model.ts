export const prerender = false

import { execute } from '@/lib/execute'
import type { APIContext } from 'astro'

export async function POST(context: APIContext) {
  const body = await context.request.json()
  const output = await execute(`fabric --changeDefaultModel ${body.model}`)
  return new Response(JSON.stringify(output))
}
