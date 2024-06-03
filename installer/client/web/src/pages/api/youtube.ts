export const prerender = false

import { execute } from '@/lib/execute'
import type { APIContext } from 'astro'

export async function POST(context: APIContext) {
  const body = await context.request.json()
  const response = await execute(`yt --transcript ${body.youtubeUrl} | fabric -p ${body.pattern} --temp ${body.temp}`)
  return new Response(JSON.stringify(response))
}
