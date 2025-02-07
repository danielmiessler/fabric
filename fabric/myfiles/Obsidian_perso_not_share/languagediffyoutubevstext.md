# language procesing comparison youtuve ve article

YOUTUBE PROCESSING LOG

INPUT:

Patterns.svelte:25 After dropdown selection - Pattern: extract_wisdom
Patterns.svelte:26 After dropdown selection - System Prompt length: 3175
ChatInput.svelte:29 YouTube URL detected: https://www.youtube.com/watch?v=ANk07xDe2Lo
ChatInput.svelte:30 Current system prompt: 3175
ChatInput.svelte:31 Selected pattern: extract_wisdom
ChatInput.svelte:282 Processing YouTube URL in handleSubmit
ChatInput.svelte:283 Current state:
ChatInput.svelte:284 - Selected Pattern: extract_wisdom
ChatInput.svelte:285 - System Prompt length: 3175
ChatInput.svelte:286 - Message content: https://www.youtube.com/watch?v=ANk07xDe2Lo

ChatInput.svelte:287 - Obsidian settings: {saveToObsidian: false, noteName: ''}
transcriptService.ts:14 
=== Getting Transcript ===
transcriptService.ts:15 1. URL: https://www.youtube.com/watch?v=ANk07xDe2Lo

transcriptService.ts:17 Fetch finished loading: POST "http://localhost:5173/chat".
window.fetch @ fetcher.js?v=09b7188a:66
getTranscript @ transcriptService.ts:17
processYouTubeURL @ ChatInput.svelte:194
handleSubmit @ ChatInput.svelte:289
(anonymous) @ chunk-ICCIXEWA.js?v=09b7188a:1249
bubble @ chunk-ICCIXEWA.js?v=09b7188a:1249
click_handler @ button.svelte:7
transcriptService.ts:38 2. Got transcript, length: 23946
transcriptService.ts:39 3. First 100 chars: Okay, so you got the perfect idea for an app. It could be a mobile app, a desktop app, but you&#39;r
ChatInput.svelte:195 Got transcript, length: 23946
ChatInput.svelte:199 Processing transcript with language: fr
ChatService.ts:105 Parsed stream response: {type: 'content', format: 'markdown', content: '```markdown\n# SUMMARY\n\nDave reviews Krisspy, an AI…project needs, optimizing cost-effectiveness.\n```'}content: "SUMMARY\nDave reviews Krisspy, an AI-powered toolformat: "markdown"type: "content"[[Prototype]]: Objectconstructor: ƒ Object()hasOwnProperty: ƒ hasOwnProperty()isPrototypeOf: ƒ isPrototypeOf()propertyIsEnumerable: ƒ propertyIsEnumerable()toLocaleString: ƒ toLocaleString()toString: ƒ toString()valueOf: ƒ valueOf()__defineGetter__: ƒ __defineGetter__()__defineSetter__: ƒ __defineSetter__()__lookupGetter__: ƒ __lookupGetter__()__lookupSetter__: ƒ __lookupSetter__()__proto__: (...)get __proto__: ƒ __proto__()set __proto__: ƒ __proto__()
ChatService.ts:28 Fetch finished loading: POST "http://localhost:5173/api/chat".
window.fetch @ fetcher.js?v=09b7188a:66
fetchStream @ ChatService.ts:28
streamChat @ ChatService.ts:186
await in streamChat
processYouTubeURL @ ChatInput.svelte:204
await in processYouTubeURL
handleSubmit @ ChatInput.svelte:289
(anonymous) @ chunk-ICCIXEWA.js?v=09b7188a:1249
bubble @ chunk-ICCIXEWA.js?v=09b7188a:1249
click_handler @ button.svelte:7
ChatService.ts:216 Processing stream chunk: {type: 'content', format: 'markdown', contentLength: 13275}contentLength: 13275format: "markdown"type: "content"[[Prototype]]: Object
ChatService.ts:125 Parsed final stream response: {type: 'complete', format: 'plain', content: ''}
ChatService.ts:216 Processing stream chunk: {type: 'complete', format: 'markdown', contentLength: 0}
ChatInput.svelte:225 Setting message format: markdown
ChatInput.svelte:225 Setting message format: markdown


NEWSPAPER URL PROCESSING LOG

INPUT:

Patterns.svelte:25 After dropdown selection - Pattern: extract_wisdom
Patterns.svelte:26 After dropdown selection - System Prompt length: 3175
ChatService.ts:105 Parsed stream response: {type: 'content', format: 'markdown', content: "```markdown\n# SUMMARY\n\nL'article du Washington Pos… la liberté d'expression et de la vie privée.\n```"}content: "SUMMARY\ the article transcript was here but not copied in this md file...."format: "markdown"type: "content"[[Prototype]]: Objectconstructor: ƒ Object()hasOwnProperty: ƒ hasOwnProperty()isPrototypeOf: ƒ isPrototypeOf()propertyIsEnumerable: ƒ propertyIsEnumerable()toLocaleString: ƒ toLocaleString()toString: ƒ toString()valueOf: ƒ valueOf()__defineGetter__: ƒ __defineGetter__()__defineSetter__: ƒ __defineSetter__()__lookupGetter__: ƒ __lookupGetter__()__lookupSetter__: ƒ __lookupSetter__()__proto__: (...)get __proto__: ƒ __proto__()set __proto__: ƒ __proto__()
ChatService.ts:28 Fetch finished loading: POST "http://localhost:5173/api/chat".
window.fetch @ fetcher.js?v=09b7188a:66
fetchStream @ ChatService.ts:28
streamChat @ ChatService.ts:186
await in streamChat
sendMessage @ chat-store.ts:87
handleSubmit @ ChatInput.svelte:292
(anonymous) @ chunk-ICCIXEWA.js?v=09b7188a:1249
bubble @ chunk-ICCIXEWA.js?v=09b7188a:1249
click_handler @ button.svelte:7
ChatService.ts:216 Processing stream chunk: {type: 'content', format: 'markdown', contentLength: 8968}contentLength: 8968format: "markdown"type: "content"[[Prototype]]: Objectconstructor: ƒ Object()hasOwnProperty: ƒ hasOwnProperty()isPrototypeOf: ƒ isPrototypeOf()propertyIsEnumerable: ƒ propertyIsEnumerable()toLocaleString: ƒ toLocaleString()toString: ƒ toString()valueOf: ƒ valueOf()__defineGetter__: ƒ __defineGetter__()__defineSetter__: ƒ __defineSetter__()__lookupGetter__: ƒ __lookupGetter__()__lookupSetter__: ƒ __lookupSetter__()__proto__: (...)get __proto__: ƒ __proto__()set __proto__: ƒ __proto__()
ChatService.ts:125 Parsed final stream response: {type: 'complete', format: 'plain', content: ''}content: ""format: "markdown"type: "complete"[[Prototype]]: Object
ChatService.ts:216 Processing stream chunk: {type: 'complete', format: 'markdown', contentLength: 0}contentLength: 0format: "markdown"type: "complete"[[Prototype]]: Object