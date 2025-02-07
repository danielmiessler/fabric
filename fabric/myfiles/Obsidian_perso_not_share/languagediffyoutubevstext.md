# language procesing comparison youtuve ve article

YOUTUBE LAST TEXT - NO LANGUAGE - TRASNCRIPT STILL SHOWING IN BROWSER BEFORE PATTERN RESULT ARE DISPLAYED

BACKEND LOG
2025/02/07 17:09:22 Received chat request with 1 prompts
2025/02/07 17:09:22 Processing prompt 1: Model= Pattern=extract_wisdom Context=
[GIN] 2025/02/07 - 17:09:45 | 200 |    22.987131s |             ::1 | POST     "/chat"


BROWSER LOG

NPUT:

Patterns.svelte:25 After dropdown selection - Pattern: extract_wisdom
Patterns.svelte:26 After dropdown selection - System Prompt length: 3175
ChatInput.svelte:41 
=== Handle Input ===
ChatInput.svelte:44 1. Raw input: https://www.youtube.com/watch?v=LMyfq8ZxZc8
ChatInput.svelte:66 2. Language detection: {detectedLang: '', currentLanguage: 'en', inputAfterLangRemoval: 'https://www.youtube.com/watch?v=LMyfq8ZxZc8'}
ChatInput.svelte:29 YouTube URL detected: https://www.youtube.com/watch?v=LMyfq8ZxZc8
ChatInput.svelte:30 Current system prompt: 3175
ChatInput.svelte:31 Selected pattern: extract_wisdom
ChatInput.svelte:73 3. URL detection: {isYouTube: true, pattern: 'extract_wisdom', systemPromptLength: 3175}
ChatInput.svelte:290 
=== Submit Handler Start ===
ChatInput.svelte:299 1. Prepared input: {isYouTube: true, language: 'fr', pattern: 'extract_wisdom', inputLength: 44}
ChatInput.svelte:307 2a. Starting YouTube flow
ChatInput.svelte:200 
=== YouTube Flow Start ===
ChatInput.svelte:203 1. Initial state: {input: 'https://www.youtube.com/watch?v=LMyfq8ZxZc8\n...', language: 'fr', pattern: 'extract_wisdom', systemPromptLength: 3175}
ChatInput.svelte:219 2. Requesting transcript with language: fr
transcriptService.ts:17 
=== YouTube Transcript Service Start ===
transcriptService.ts:18 1. Request details: {url: 'https://www.youtube.com/watch?v=LMyfq8ZxZc8\n', endpoint: '/chat', method: 'POST', isYouTubeURL: true, currentLanguage: 'fr'}
transcriptService.ts:37 2. Server response: {status: 200, ok: true, type: 'basic', language: 'fr'}
transcriptService.ts:57 3. Processed transcript: {status: 200, transcriptLength: 8671, firstChars: 'if you own a vehicle long enough sooner or later it&#39;s going to experience transmission problems ', hasError: false, videoId: 'LMyfq8ZxZc8', …}
ChatInput.svelte:221 3. Got transcript: {length: 8671, originalLanguage: 'fr', currentLanguage: 'fr', firstChars: 'if you own a vehicle long enough sooner or later i'}
ChatInput.svelte:235 4b. Preparing to process transcript: {language: 'fr', pattern: 'extract_wisdom', transcriptLength: 8671, hasLanguageInstruction: true}
ChatInput.svelte:243 5. Sending to message processing...
chat-store.ts:74 
=== Message Processing Start ===
chat-store.ts:75 1. Initial state: {content: 'if you own a vehicle long enough sooner or later i…9;s going to experience transmission problems ...', isSystem: false, hasSystemPrompt: true, currentLanguage: 'fr', messageCount: 1, …}
chat-store.ts:98 2. Message added to store: {role: 'user', contentLength: 8671, totalMessages: 2, language: 'fr'}
chat-store.ts:106 3. Preparing chat stream: {language: 'fr', pattern: 'extract_wisdom', hasSystemPrompt: true, systemPromptLength: 3175}
ChatService.ts:240 
=== Creating Chat Request ===
ChatService.ts:241 1. Input: {userInput: 'if you own a vehicle long enough sooner or later i…e take care and I look forward to seeing you then', isPattern: false, language: 'fr'}
ChatService.ts:202 
=== Creating Chat Prompt ===
ChatService.ts:203 1. Current state: {language: 'fr', hasLanguageInstruction: true, instruction: ". Please use the language 'fr' for the output.", pattern: 'extract_wisdom', inputLength: 8671, …}
ChatService.ts:228 2. Created prompt: {finalInput: 'if you own a vehicle long enough sooner or later i…9;s going to experience transmission problems ...', hasLanguageInInput: true, languagePosition: 1926, pattern: 'extract_wisdom', language: 'fr'}
ChatService.ts:259 2. Final request: {promptCount: 1, messageCount: 2, firstPromptInput: "if you own a vehicle long enough sooner or later i…hen. Please use the language 'fr' for the output.", hasLanguageInPrompt: true}
ChatService.ts:28 
=== ChatService Request Start ===
ChatService.ts:29 1. Request details: {language: 'fr', pattern: 'extract_wisdom', promptCount: 1, messageCount: 2}
ChatService.ts:36 2. First prompt: {pattern: 'extract_wisdom', inputLength: 8717, hasLanguageInInput: true, systemPromptLength: 3175}
ChatService.ts:43 3. Full request: {
  "prompts": [
    {
      "userInput": "if you own a vehicle long enough sooner or later it&#39;s going to experience transmission problems and when it does should you try a product such as Lucas transmission fix now this product does not claim to fix a broken part but it does make quite a few claims regarding improving the performance of a worn transmission the test subject today is going to be a 2000 Honda Accord with over 220,000 miles on it and it has some pretty significant issues so let&#39;s get the testing underway and see if this product can help this is a 2000 Honda Accord that has a 3 liter VTEC and a four-speed automatic transmission the owner of the vehicle has changed a transmission fluid and it did not solve the problem also this vehicle was taken to a transmission shop in the transmission shop wanted to do a complete rebuild they indicated that the torque converter had gone bad and the only way to repair the transmission was just to do a complete rebuild regarding 2p0 seven-four-zero code I&#39;m told that the TCC solenoid was replaced by the owner I&#39;m not sure if the codes are reset when the solenoid was replaced so the fact that the check engine light was on and isn&#39;t on now doesn&#39;t mean much the owner of the vehicle does not want to spend $1,500 to have the transmission rebuilt so we&#39;re gonna try this lucas transmission fix i&#39;m going to go ahead and see what the transmission fluid looks like so the transmission fluid does have a pink color to it with that said we really can&#39;t trust the color of the transmission fluid to give us much information just because it has been changed so let&#39;s take the vehicle on a quick road test and we&#39;ll see what it&#39;s doing and then come back and see if the Lucas can help okay we&#39;re in first second not going into third there it goes then in the 4th 1st and 2nd is fine 2nd to 3rd is the problem 3rd to 4th is fine [Music] okay did not want a shift from second to third [Music] I&#39;m not going to be adding any transmission fluid to the Honda I&#39;m only using this AC Delco dexron six for testing purposes I&#39;ll be comparing it against that Lukas transmission fix let&#39;s go ahead and send off the Lucas transmission fix and the ATF to an oil testing lab to see how they compare we&#39;ll take a look at the lab test results later in the video before we add Lucas to the transmission to the Honda let&#39;s do some testing on it beginning with the evaporative loss test I&#39;ll first measure out 200 grams of each product into the oil containers and we&#39;ll heat the oil to around 300 degrees Fahrenheit for two hours I&#39;ll rotate the oil containers every 10 minutes just in case there are hot spots on the griddle I&#39;ll be monitoring the temperature just to make sure that the oil containers are nearly the same temperature during the entire test a transmission fluids ability to withstand heat has a huge impact on the performance and longevity of the transmission the nowak volatility test is one method of testing automatic transmission fluid well I don&#39;t have the test equipment to conduct a no act volatility test this process will give us some great information regarding the transmission fluid and the Lucas transmission fix it&#39;s been run at two hours I&#39;m gonna go ahead remove both of these containers from the griddle and we come back we&#39;ll see a much evaporative loss occurred with these product Lucas start off weighing 394 0.62 grams and now weighs 394 that&#39;s only a loss of 0.62 grams the ATF&#39;s start off at 430 grams it now weighs 400 28.7 that&#39;s a loss of 1.3 grams so regarding evaporative loss Lucas did a very good job Lucas has a very high viscosity so it&#39;ll be interesting to see how it compares to automatic transmission fluid when it&#39;s very cold I&#39;ll place the new and the cooked products in a freezer that set to 20 below zero Fahrenheit and we&#39;ll come back to this later in the video up next we&#39;ll compare the film strengths of Lucas and automatic transmission fluid I&#39;ll measure out 40 milliliters of the cooked products that we&#39;ll be using during the test the test will last right at ten minutes once finished testing the automatic transmission fluid I&#39;ll clean the tester and then we&#39;ll test Lucas when it comes to transmission line Jeff D having the right blend of detergents dispersants and anti wear additives can make a huge difference while lubricity tester doesn&#39;t simulate transmission operating conditions perfectly it&#39;ll provide us with a lot of great information [Music] [Music] [Music] [Music] automatic transmission fluid on the left and the Lucas is on the right there is a huge difference between the two with Lucas doing far better than just the automatic transmission fluid new automatic transmission fluid is in lane one cooked ATF lane two cooked Lucas lane three new Lucas Lane for an automatic transmission fluid is out of the gate quickly an
transcriptService.ts:26 Fetch finished loading: POST "http://localhost:5173/chat".
window.fetch @ fetcher.js?v=09b7188a:66
getTranscript @ transcriptService.ts:26
processYouTubeURL @ ChatInput.svelte:220
handleSubmit @ ChatInput.svelte:308
(anonymous) @ chunk-ICCIXEWA.js?v=09b7188a:1249
bubble @ chunk-ICCIXEWA.js?v=09b7188a:1249
click_handler @ button.svelte:7
ChatService.ts:53 4. Response received: {status: 200, ok: true, type: 'basic'}
ChatService.ts:72 5. Creating message stream
ChatService.ts:114 
=== Stream Processing Start ===
chat-store.ts:114 4. Stream created, beginning processing
ChatService.ts:133 Processing chunk: {contentLength: 13628, format: 'markdown', type: 'content', hasPattern: true}
ChatService.ts:306 Processing stream chunk: {type: 'content', format: 'markdown', contentLength: 13592}
chat-store.ts:119 5. Processing stream chunk: {contentLength: 13592, format: 'markdown', type: 'content', language: 'fr'}
chat-store.ts:146 6b. Added new message: {role: 'assistant', contentLength: 13592, format: 'markdown'}
ChatService.ts:118 Stream complete
ChatService.ts:158 Processing final chunk: {contentLength: 0, format: 'plain', type: 'complete', hasPattern: true}
ChatService.ts:306 Processing stream chunk: {type: 'complete', format: 'markdown', contentLength: 0}
chat-store.ts:119 5. Processing stream chunk: {contentLength: 13592, format: 'markdown', type: 'complete', language: 'fr'}
chat-store.ts:135 6a. Updated existing message: {role: 'assistant', contentLength: 13592, format: 'markdown'}
chat-store.ts:119 5. Processing stream chunk: {contentLength: 13592, format: 'markdown', type: 'complete', language: 'fr'}
chat-store.ts:135 6a. Updated existing message: {role: 'assistant', contentLength: 13592, format: 'markdown'}
ChatInput.svelte:247 6. Post-processing language check: {originalLanguage: 'fr', finalLanguage: 'fr', languageMatches: true}
ChatService.ts:45 Fetch finished loading: POST "http://localhost:5173/api/chat".
window.fetch @ fetcher.js?v=09b7188a:66
fetchStream @ ChatService.ts:45
streamChat @ ChatService.ts:276
await in streamChat
sendMessage @ chat-store.ts:113
processYouTubeURL @ ChatInput.svelte:244
await in processYouTubeURL
handleSubmit @ ChatInput.svelte:308
(anonymous) @ chunk-ICCIXEWA.js?v=09b7188a:1249
bubble @ chunk-ICCIXEWA.js?v=09b7188a:1249
click_handler @ button.svelte:7



NEWS ARTICLE LOG - PROPER LANGUAGE DETECTION AND PATTERN PROCESSING

INPUT:

Patterns.svelte:25 After dropdown selection - Pattern: extract_wisdom
Patterns.svelte:26 After dropdown selection - System Prompt length: 3175
ChatInput.svelte:41 
=== Handle Input ===
ChatInput.svelte:44 1. Raw input: https://www.washingtonpost.com/technology/2025/02/07/apple-encryption-backdoor-uk/
ChatInput.svelte:66 2. Language detection: {detectedLang: '', currentLanguage: 'fr', inputAfterLangRemoval: 'https://www.washingtonpost.com/technology/2025/02/07/apple-encryption-backdoor-uk/'}
ChatInput.svelte:73 3. URL detection: {isYouTube: false, pattern: 'extract_wisdom', systemPromptLength: 3175}
ChatInput.svelte:273 
=== Submit Handler Start ===
ChatInput.svelte:282 1. Prepared input: {isYouTube: false, language: 'fr', pattern: 'extract_wisdom', inputLength: 83}
ChatInput.svelte:293 2b. Starting regular text flow
chat-store.ts:74 
=== Message Processing Start ===
chat-store.ts:75 1. Initial state: {content: 'https://www.washingtonpost.com/technology/2025/02/07/apple-encryption-backdoor-uk/\n...', isSystem: false, hasSystemPrompt: false, currentLanguage: 'fr', messageCount: 0, …}
chat-store.ts:98 2. Message added to store: {role: 'user', contentLength: 83, totalMessages: 1, language: 'fr'}
chat-store.ts:106 3. Preparing chat stream: {language: 'fr', pattern: 'extract_wisdom', hasSystemPrompt: false, systemPromptLength: undefined}
ChatService.ts:227 
=== Creating Chat Request ===
ChatService.ts:228 1. Input: {userInput: 'https://www.washingtonpost.com/technology/2025/02/07/apple-encryption-backdoor-uk/\n', isPattern: false, language: 'fr'}
ChatService.ts:202 
=== Creating Chat Prompt ===
ChatService.ts:203 1. Current state: {language: 'fr', hasLanguageInstruction: true, instruction: ". Please use the language 'fr' for the output.", pattern: 'extract_wisdom'}
ChatService.ts:217 2. Created prompt: {finalInput: "https://www.washingtonpost.com/technology/2025/02/…k/\n. Please use the language 'fr' for the output.", hasLanguageInInput: true, pattern: 'extract_wisdom'}
ChatService.ts:246 2. Final request: {promptCount: 1, messageCount: 1, firstPromptInput: "https://www.washingtonpost.com/technology/2025/02/…k/\n. Please use the language 'fr' for the output.", hasLanguageInPrompt: true}
ChatService.ts:28 
=== ChatService Request Start ===
ChatService.ts:29 1. Request details: {language: 'fr', pattern: 'extract_wisdom', promptCount: 1, messageCount: 1}
ChatService.ts:36 2. First prompt: {pattern: 'extract_wisdom', inputLength: 129, hasLanguageInInput: true, systemPromptLength: 3175}
ChatService.ts:43 3. Full request: {
  "prompts": [
    {
      "userInput": "https://www.washingtonpost.com/technology/2025/02/07/apple-encryption-backdoor-uk/\n. Please use the language 'fr' for the output.",
      "systemPrompt": "# IDENTITY and PURPOSE\n\nYou extract surprising, insightful, and interesting information from text content. You are interested in insights related to the purpose and meaning of life, human flourishing, the role of technology in the future of humanity, artificial intelligence and its affect on humans, memes, learning, reading, books, continuous improvement, and similar topics.\n\nTake a step back and think step-by-step about how to achieve the best possible results by following the steps below.\n\n# STEPS\n\n- Extract a summary of the content in 25 words, including who is presenting and the content being discussed into a section called SUMMARY.\n\n- Extract 20 to 50 of the most surprising, insightful, and/or interesting ideas from the input in a section called IDEAS:. If there are less than 50 then collect all of them. Make sure you extract at least 20.\n\n- Extract 10 to 20 of the best insights from the input and from a combination of the raw input and the IDEAS above into a section called INSIGHTS. These INSIGHTS should be fewer, more refined, more insightful, and more abstracted versions of the best ideas in the content. \n\n- Extract 15 to 30 of the most surprising, insightful, and/or interesting quotes from the input into a section called QUOTES:. Use the exact quote text from the input.\n\n- Extract 15 to 30 of the most practical and useful personal habits of the speakers, or mentioned by the speakers, in the content into a section called HABITS. Examples include but aren't limited to: sleep schedule, reading habits, things they always do, things they always avoid, productivity tips, diet, exercise, etc.\n\n- Extract 15 to 30 of the most surprising, insightful, and/or interesting valid facts about the greater world that were mentioned in the content into a section called FACTS:.\n\n- Extract all mentions of writing, art, tools, projects and other sources of inspiration mentioned by the speakers into a section called REFERENCES. This should include any and all references to something that the speaker mentioned.\n\n- Extract the most potent takeaway and recommendation into a section called ONE-SENTENCE TAKEAWAY. This should be a 15-word sentence that captures the most important essence of the content.\n\n- Extract the 15 to 30 of the most surprising, insightful, and/or interesting recommendations that can be collected from the content into a section called RECOMMENDATIONS.\n\n# OUTPUT INSTRUCTIONS\n\n- Only output Markdown.\n\n- Write the IDEAS bullets as exactly 16 words.\n\n- Write the RECOMMENDATIONS bullets as exactly 16 words.\n\n- Write the HABITS bullets as exactly 16 words.\n\n- Write the FACTS bullets as exactly 16 words.\n\n- Write the INSIGHTS bullets as exactly 16 words.\n\n- Extract at least 25 IDEAS from the content.\n\n- Extract at least 10 INSIGHTS from the content.\n\n- Extract at least 20 items for the other output sections.\n\n- Do not give warnings or notes; only output the requested sections.\n\n- You use bulleted lists for output, not numbered lists.\n\n- Do not repeat ideas, quotes, facts, or resources.\n\n- Do not start items with the same opening words.\n\n- Ensure you follow ALL these instructions when creating your output.\n\n# INPUT\n\nINPUT:\n",
      "model": "",
      "patternName": "extract_wisdom"
    }
  ],
  "messages": [
    {
      "role": "user",
      "content": "https://www.washingtonpost.com/technology/2025/02/07/apple-encryption-backdoor-uk/\n"
    }
  ],
  "temperature": 0.7,
  "top_p": 1,
  "frequency_penalty": 0,
  "presence_penalty": 0
}
ChatService.ts:53 4. Response received: {status: 200, ok: true, type: 'basic'}
ChatService.ts:72 5. Creating message stream
ChatService.ts:114 
=== Stream Processing Start ===
chat-store.ts:114 4. Stream created, beginning processing
ChatService.ts:133 Processing chunk: {contentLength: 10830, format: 'markdown', type: 'content', hasPattern: true}
ChatService.ts:293 Processing stream chunk: {type: 'content', format: 'markdown', contentLength: 10829}
chat-store.ts:119 5. Processing stream chunk: {contentLength: 10829, format: 'markdown', type: 'content', language: 'fr'}
chat-store.ts:146 6b. Added new message: {role: 'assistant', contentLength: 10829, format: 'markdown'}
ChatService.ts:118 Stream complete
ChatService.ts:158 Processing final chunk: {contentLength: 0, format: 'plain', type: 'complete', hasPattern: true}
ChatService.ts:293 Processing stream chunk: {type: 'complete', format: 'markdown', contentLength: 0}
chat-store.ts:119 5. Processing stream chunk: {contentLength: 10829, format: 'markdown', type: 'complete', language: 'fr'}
chat-store.ts:135 6a. Updated existing message: {role: 'assistant', contentLength: 10829, format: 'markdown'}
chat-store.ts:119 5. Processing stream chunk: {contentLength: 10829, format: 'markdown', type: 'complete', language: 'fr'}
chat-store.ts:135 6a. Updated existing message: {role: 'assistant', contentLength: 10829, format: 'markdown'}
ChatService.ts:45 Fetch finished loading: POST "http://localhost:5173/api/chat".
window.fetch @ fetcher.js?v=09b7188a:66
fetchStream @ ChatService.ts:45
streamChat @ ChatService.ts:263
await in streamChat
sendMessage @ chat-store.ts:113
handleSubmit @ ChatInput.svelte:294
(anonymous) @ chunk-ICCIXEWA.js?v=09b7188a:1249
bubble @ chunk-ICCIXEWA.js?v=09b7188a:1249
click_handler @ button.svelte:7



YOUTUBE PROCESSING - NO LANGUAGE PICKED UP AND BPOTH TRANSCRIPTS AND PATTERN OUTPUT SHOW IN BROWSER

=== Handle Input ===
ChatInput.svelte:44 1. Raw input: https://www.youtube.com/watch?v=LMyfq8ZxZc8
ChatInput.svelte:66 2. Language detection: {detectedLang: '', currentLanguage: 'en', inputAfterLangRemoval: 'https://www.youtube.com/watch?v=LMyfq8ZxZc8'}
ChatInput.svelte:29 YouTube URL detected: https://www.youtube.com/watch?v=LMyfq8ZxZc8
ChatInput.svelte:30 Current system prompt: 3175
ChatInput.svelte:31 Selected pattern: extract_wisdom
ChatInput.svelte:73 3. URL detection: {isYouTube: true, pattern: 'extract_wisdom', systemPromptLength: 3175}
ChatInput.svelte:273 
=== Submit Handler Start ===
ChatInput.svelte:282 1. Prepared input: {isYouTube: true, language: 'en', pattern: 'extract_wisdom', inputLength: 44}
ChatInput.svelte:290 2a. Starting YouTube flow
ChatInput.svelte:200 
=== YouTube Flow Start ===
ChatInput.svelte:201 1. Initial state: {input: 'https://www.youtube.com/watch?v=LMyfq8ZxZc8\n...', language: 'en', pattern: 'extract_wisdom', systemPromptLength: 3175}
ChatInput.svelte:217 2. Requesting transcript...
transcriptService.ts:14 
=== YouTube Transcript Service Start ===
transcriptService.ts:15 1. Request details: {url: 'https://www.youtube.com/watch?v=LMyfq8ZxZc8\n', endpoint: '/chat', method: 'POST', isYouTubeURL: true}
transcriptService.ts:30 2. Server response: {status: 200, ok: true, type: 'basic'}
transcriptService.ts:49 3. Processed transcript: {status: 200, transcriptLength: 8671, firstChars: 'if you own a vehicle long enough sooner or later it&#39;s going to experience transmission problems ', hasError: false, videoId: 'LMyfq8ZxZc8'}
ChatInput.svelte:219 3. Got transcript: {length: 8671, language: 'en', firstChars: 'if you own a vehicle long enough sooner or later i'}
ChatInput.svelte:226 4. Preparing to process transcript: {language: 'en', pattern: 'extract_wisdom', transcriptLength: 8671}
ChatInput.svelte:233 5. Sending to message processing...
chat-store.ts:74 
=== Message Processing Start ===
chat-store.ts:75 1. Initial state: {content: 'if you own a vehicle long enough sooner or later i…9;s going to experience transmission problems ...', isSystem: false, hasSystemPrompt: true, currentLanguage: 'en', messageCount: 1, …}
chat-store.ts:98 2. Message added to store: {role: 'user', contentLength: 8671, totalMessages: 2, language: 'en'}
chat-store.ts:106 3. Preparing chat stream: {language: 'en', pattern: 'extract_wisdom', hasSystemPrompt: true, systemPromptLength: 3175}
ChatService.ts:227 
=== Creating Chat Request ===
ChatService.ts:228 1. Input: {userInput: 'if you own a vehicle long enough sooner or later i…e take care and I look forward to seeing you then', isPattern: false, language: 'en'}
ChatService.ts:202 
=== Creating Chat Prompt ===
ChatService.ts:203 1. Current state: {language: 'en', hasLanguageInstruction: false, instruction: '', pattern: 'extract_wisdom'}
transcriptService.ts:22 Fetch finished loading: POST "http://localhost:5173/chat".
window.fetch @ fetcher.js?v=09b7188a:66
getTranscript @ transcriptService.ts:22
processYouTubeURL @ ChatInput.svelte:218
handleSubmit @ ChatInput.svelte:291
(anonymous) @ chunk-ICCIXEWA.js?v=09b7188a:1249
bubble @ chunk-ICCIXEWA.js?v=09b7188a:1249
click_handler @ button.svelte:7
ChatService.ts:217 2. Created prompt: {finalInput: 'if you own a vehicle long enough sooner or later i…e take care and I look forward to seeing you then', hasLanguageInInput: true, pattern: 'extract_wisdom'}
ChatService.ts:246 2. Final request: {promptCount: 1, messageCount: 2, firstPromptInput: 'if you own a vehicle long enough sooner or later i…e take care and I look forward to seeing you then', hasLanguageInPrompt: true}
ChatService.ts:28 
=== ChatService Request Start ===
ChatService.ts:29 1. Request details: {language: 'en', pattern: 'extract_wisdom', promptCount: 1, messageCount: 2}
ChatService.ts:36 2. First prompt: {pattern: 'extract_wisdom', inputLength: 8671, hasLanguageInInput: true, systemPromptLength: 3175}
ChatService.ts:43 3. Full request: {
  "prompts": [
    {
      "userInput": "if you own a vehicle long enough sooner or later it&#39;s going to experience transmission problems and when it does should you try a product such as Lucas transmission fix now this product does not claim to fix a broken part but it does make quite a few claims regarding improving the performance of a worn transmission the test subject today is going to be a 2000 Honda Accord with over 220,000 miles on it and it has some pretty significant issues so let&#39;s get the testing underway and see if this product can help this is a 2000 Honda Accord that has a 3 liter VTEC and a four-speed automatic transmission the owner of the vehicle has changed a transmission fluid and it did not solve the problem also this vehicle was taken to a transmission shop in the transmission shop wanted to do a complete rebuild they indicated that the torque converter had gone bad and the only way to repair the transmission was just to do a complete rebuild regarding 2p0 seven-four-zero code I&#39;m told that the TCC solenoid was replaced by the owner I&#39;m not sure if the codes are reset when the solenoid was replaced so the fact that the check engine light was on and isn&#39;t on now doesn&#39;t mean much the owner of the vehicle does not want to spend $1,500 to have the transmission rebuilt so we&#39;re gonna try this lucas transmission fix i&#39;m going to go ahead and see what the transmission fluid looks like so the transmission fluid does have a pink color to it with that said we really can&#39;t trust the color of the transmission fluid to give us much information just because it has been changed so let&#39;s take the vehicle on a quick road test and we&#39;ll see what it&#39;s doing and then come back and see if the Lucas can help okay we&#39;re in first second not going into third there it goes then in the 4th 1st and 2nd is fine 2nd to 3rd is the problem 3rd to 4th is fine [Music] okay did not want a shift from second to third [Music] I&#39;m not going to be adding any transmission fluid to the Honda I&#39;m only using this AC Delco dexron six for testing purposes I&#39;ll be comparing it against that Lukas transmission fix let&#39;s go ahead and send off the Lucas transmission fix and the ATF to an oil testing lab to see how they compare we&#39;ll take a look at the lab test results later in the video before we add Lucas to the transmission to the Honda let&#39;s do some testing on it beginning with the evaporative loss test I&#39;ll first measure out 200 grams of each product into the oil containers and we&#39;ll heat the oil to around 300 degrees Fahrenheit for two hours I&#39;ll rotate the oil containers every 10 minutes just in case there are hot spots on the griddle I&#39;ll be monitoring the temperature just to make sure that the oil containers are nearly the same temperature during the entire test a transmission fluids ability to withstand heat has a huge impact on the performance and longevity of the transmission the nowak volatility test is one method of testing automatic transmission fluid well I don&#39;t have the test equipment to conduct a no act volatility test this process will give us some great information regarding the transmission fluid and the Lucas transmission fix it&#39;s been run at two hours I&#39;m gonna go ahead remove both of these containers from the griddle and we come back we&#39;ll see a much evaporative loss occurred with these product Lucas start off weighing 394 0.62 grams and now weighs 394 that&#39;s only a loss of 0.62 grams the ATF&#39;s start off at 430 grams it now weighs 400 28.7 that&#39;s a loss of 1.3 grams so regarding evaporative loss Lucas did a very good job Lucas has a very high viscosity so it&#39;ll be interesting to see how it compares to automatic transmission fluid when it&#39;s very cold I&#39;ll place the new and the cooked products in a freezer that set to 20 below zero Fahrenheit and we&#39;ll come back to this later in the video up next we&#39;ll compare the film strengths of Lucas and automatic transmission fluid I&#39;ll measure out 40 milliliters of the cooked products that we&#39;ll be using during the test the test will last right at ten minutes once finished testing the automatic transmission fluid I&#39;ll clean the tester and then we&#39;ll test Lucas when it comes to transmission line Jeff D having the right blend of detergents dispersants and anti wear additives can make a huge difference while lubricity tester doesn&#39;t simulate transmission operating conditions perfectly it&#39;ll provide us with a lot of great information [Music] [Music] [Music] [Music] automatic transmission fluid on the left and the Lucas is on the right there is a huge difference between the two with Lucas doing far better than just the automatic transmission fluid new automatic transmission fluid is in lane one cooked ATF lane two cooked Lucas lane three new Lucas Lane for an automatic transmission fluid is out of the gate quickly an
ChatService.ts:53 4. Response received: {status: 200, ok: true, type: 'basic'}
ChatService.ts:72 5. Creating message stream
ChatService.ts:114 
=== Stream Processing Start ===
chat-store.ts:114 4. Stream created, beginning processing
ChatService.ts:133 Processing chunk: {contentLength: 12737, format: 'markdown', type: 'content', hasPattern: true}
ChatService.ts:293 Processing stream chunk: {type: 'content', format: 'markdown', contentLength: 12702}
chat-store.ts:119 5. Processing stream chunk: {contentLength: 12702, format: 'markdown', type: 'content', language: 'en'}
chat-store.ts:146 6b. Added new message: {role: 'assistant', contentLength: 12702, format: 'markdown'}
ChatService.ts:118 Stream complete
ChatService.ts:158 Processing final chunk: {contentLength: 0, format: 'plain', type: 'complete', hasPattern: true}
ChatService.ts:45 Fetch finished loading: POST "http://localhost:5173/api/chat".
window.fetch @ fetcher.js?v=09b7188a:66
fetchStream @ ChatService.ts:45
streamChat @ ChatService.ts:263
await in streamChat
sendMessage @ chat-store.ts:113
processYouTubeURL @ ChatInput.svelte:234
await in processYouTubeURL
handleSubmit @ ChatInput.svelte:291
(anonymous) @ chunk-ICCIXEWA.js?v=09b7188a:1249
bubble @ chunk-ICCIXEWA.js?v=09b7188a:1249
click_handler @ button.svelte:7
ChatService.ts:293 Processing stream chunk: {type: 'complete', format: 'markdown', contentLength: 0}
chat-store.ts:119 5. Processing stream chunk: {contentLength: 12702, format: 'markdown', type: 'complete', language: 'en'}
chat-store.ts:135 6a. Updated existing message: {role: 'assistant', contentLength: 12702, format: 'markdown'}
chat-store.ts:119 5. Processing stream chunk: {contentLength: 12702, format: 'markdown', type: 'complete', language: 'en'}
chat-store.ts:135 6a. Updated existing message: {role: 'assistant', contentLength: 12702, format: 'markdown'}