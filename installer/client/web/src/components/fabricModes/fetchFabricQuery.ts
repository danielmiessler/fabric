export type FabricTextData = { query: string, apiurl: string, pattern: string }
export type FabricYoutubeData = { youtubeUrl: string, apiurl: string, pattern: string }
export type FabricQueryProps = FabricTextData | FabricYoutubeData

export const defaultFabricQueryProps = { query: 'Describe in 2 sentences or less the meaning of life.', apiurl: 'api/query', pattern: 'ai' }

export const fetchFabricQuery = async (payload: FabricQueryProps): Promise<string[]> => {
  console.log({ fetchFabricQuery: payload })
  const response = await fetch(payload.apiurl, { method: "POST", body: JSON.stringify(payload) });
  const body = await response.json();
  console.log({ api_response: true, body })
  if (body.ok) {
    return body.data
  }
  return []
}
