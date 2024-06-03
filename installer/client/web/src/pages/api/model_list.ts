export const prerender = false

import { execute } from '@/lib/execute'
import type { APIContext } from 'astro'

export async function GET(context: APIContext) {
  const patternlist = await execute('fabric --listmodels')
  return new Response(JSON.stringify(patternlist))
}
