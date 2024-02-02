<div align="center">

<img src="./images/fabric-logo-gif.gif" alt="fabriclogo" width="400" height="400"/>

# `fabric`

<h4><code>fabric</code> is an open-source framework for augmenting humans using AI.</h4>

[What and Why](#whatandwhy) •
[Quickstart](#quickstart) •
[Usage](#usage) •
[Examples](#examples) •
[Structure](#structure) •
[Naming](#naming) •
[Meta](#meta)

</div>

<br />

## What and Why

Since the start of 2023 and GenAI we've seen a massive number of AI applications for accomplishing tasks. It's powerful, but **it's not easy to integrate this functionality into our lives.**

_In other words, AI doesn't have a capabilities problem—it has an **integration** problem._

Fabric was created to address that problem by allowing everyone to leverage AI throughout our life and work.

### Too many prompts

The biggest challenge I faced in 2023——and still have today—is **the sheer number of AI prompts out there**. We all have prompts that are useful, but it's hard to manage them, discover new ones, _and manage the different versions of the ones we like_.

One of <code>fabric</code>'s main features is helping people collect and integrate modular AI functionality (in this case: prompts), which we call _Patterns_, into various parts of their lives.

Fabric has patterns (prompts) for all sorts of life and work activities, including:

- Extracting the most interesting parts of YouTube videos and podcasts
- Writing an essay in your own voice with just an idea as an input
- Summarizing opaque academic papers
- Creating perfectly matched AI art prompts for a piece of writing
- Rating the quality of content to see if you want to read/watch the whole thing
- Getting summaries of long, boring content
- Explaining code to you
- Turning bad documentation into usable documentation
- Create social media posts from any content input
- And a million more…

## Quickstart

There are three main ways to get started with Fabric.

<img width="1173" alt="fabric-patterns-screenshot" src="https://github.com/danielmiessler/fabric/assets/50654/9186a044-652b-4673-89f7-71cf066f32d8">

### 1. Just use the patterns (prompts)

If you're not looking to do anything fancy, and you just want a lot of great prompts, you can navigate to the [`/patterns`](https://github.com/danielmiessler/fabric/tree/main/patterns) directory and start exploring! 

You can use any of those in any AI application that you have!

### 2. Create your own Fabric Mill (server)

If you want your very own Fabric server, head over to the [`/server/`](https://github.com/danielmiessler/fabric/tree/main/server) directory and set up your own Fabric Mill with your own patterns running! You can then use the [`/client/standalone_client_examples`](https://github.com/danielmiessler/fabric/tree/main/client/standalone_client_examples) to connect to it.

### 3. The standalone client 

We're almost done with a universal client that will let you do all sorts of cool stuff, including:

- Calling patterns without connecting to a Fabric server (direct to OpenAI).
- Streaming mode to get instant and dynamic results.
- Other cool stuff…

We expect this client to be ready very within a day or two, and we'll update the Quickstart as soon as it is.

## Usage

`fabric`'s main function is to make **Patterns** available to everyone in an open ecosystem, i.e., to allow people to share and fork prompts in a transparent, scalable, and dependable way.

But it also includes two other components that make it possible for AI enthusiasts and developers to _build your own Personal AI Ecosystem_.

_In other words you can have your own server, with your own copy of `fabric`, running your own specific combination of **Patterns** for your specific use cases._

### Components

Here are the three `fabric` ecosystem pieces, and how they work together.

- The **Mill** is the (optional) server that makes **Patterns** available.
- **Patterns** are the actual AI use cases.
- **Looms** are the modular, client-side apps that call a specific **Pattern** hosted by a **Mill**.

One key feature of `fabric` and its Markdown-based format is the ability to ** directly reference** (and edit) individual [patterns](https://github.com/danielmiessler/fabric/tree/main#naming) directly—on their own—without surrounding code.

As an example, here's how to call _the direct location_ of the **system** prompt for the `extract_wisdom` pattern.

```
https://github.com/danielmiessler/fabric/blob/main/patterns/extract_wisdom/system.md
```

This means you can cleanly, and directly reference any pattern for use in a web-based AI app, your own code, or wherever!

Even better, you can also have your [Mill](https://github.com/danielmiessler/fabric/tree/main#naming) functionality directly call **system** and **user** prompts from `fabric`, meaning you can have your personal AI ecosystem automatically kept up to date with the latest version of your favorite [Patterns](https://github.com/danielmiessler/fabric/tree/main#naming).

## Examples

Here's an abridged output example from the <a href="https://github.com/danielmiessler/fabric/blob/main/patterns/extract_wisdom/system.md">`extract_wisdom`</a> pattern (limited to only 10 items per section).

```bash
# Paste in the transcript of a YouTube video of Riva Tez on David Perrel's podcast
pbpaste | extract_wisdom
```

```markdown
## SUMMARY:

The content features a conversation between two individuals discussing various topics, including the decline of Western culture, the importance of beauty and subtlety in life, the impact of technology and AI, the resonance of Rilke's poetry, the value of deep reading and revisiting texts, the captivating nature of Ayn Rand's writing, the role of philosophy in understanding the world, and the influence of drugs on society. They also touch upon creativity, attention spans, and the importance of introspection.

## IDEAS:

1. Western culture is perceived to be declining due to a loss of values and an embrace of mediocrity.
2. Mass media and technology have contributed to shorter attention spans and a need for constant stimulation.
3. Rilke's poetry resonates due to its focus on beauty and ecstasy in everyday objects.
4. Subtlety is often overlooked in modern society due to sensory overload.
5. The role of technology in shaping music and performance art is significant.
6. Reading habits have shifted from deep, repetitive reading to consuming large quantities of new material.
7. Revisiting influential books as one ages can lead to new insights based on accumulated wisdom and experiences.
8. Fiction can vividly illustrate philosophical concepts through characters and narratives.
9. Many influential thinkers have backgrounds in philosophy, highlighting its importance in shaping reasoning skills.
10. Philosophy is seen as a bridge between theology and science, asking questions that both fields seek to answer.

## QUOTES:

1. "You can't necessarily think yourself into the answers. You have to create space for the answers to come to you."
2. "The West is dying and we are killing her."
3. "The American Dream has been replaced by mass packaged mediocrity porn, encouraging us to revel like happy pigs in our own meekness."
4. "There's just not that many people who have the courage to reach beyond consensus and go explore new ideas."
5. "I'll start watching Netflix when I've read the whole of human history."
6. "Rilke saw beauty in everything... He sees it's in one little thing, a representation of all things that are beautiful."
7. "Vanilla is a very subtle flavor... it speaks to sort of the sensory overload of the modern age."
8. "When you memorize chapters [of the Bible], it takes a few months, but you really understand how things are structured."
9. "As you get older, if there's books that moved you when you were younger, it's worth going back and rereading them."
10. "She [Ayn Rand] took complicated philosophy and embodied it in a way that anybody could resonate with."

## HABITS:

1. Avoiding mainstream media consumption for deeper engagement with historical texts and personal research.
2. Regularly revisiting influential books from youth to gain new insights with age.
3. Engaging in deep reading practices rather than skimming or speed-reading material.
4. Memorizing entire chapters or passages from significant texts for better understanding.
5. Disengaging from social media and fast-paced news cycles for more focused thought processes.
6. Walking long distances as a form of meditation and reflection.
7. Creating space for thoughts to solidify through introspection and stillness.
8. Embracing emotions such as grief or anger fully rather than suppressing them.
9. Seeking out varied experiences across different careers and lifestyles.
10. Prioritizing curiosity-driven research without specific goals or constraints.

## FACTS:

1. The West is perceived as declining due to cultural shifts away from traditional values.
2. Attention spans have shortened due to technological advancements and media consumption habits.
3. Rilke's poetry emphasizes finding beauty in everyday objects through detailed observation.
4. Modern society often overlooks subtlety due to sensory overload from various stimuli.
5. Reading habits have evolved from deep engagement with texts to consuming large quantities quickly.
6. Revisiting influential books can lead to new insights based on accumulated life experiences.
7. Fiction can effectively illustrate philosophical concepts through character development and narrative arcs.
8. Philosophy plays a significant role in shaping reasoning skills and understanding complex ideas.
9. Creativity may be stifled by cultural nihilism and protectionist attitudes within society.
10. Short-term thinking undermines efforts to create lasting works of beauty or significance.

## REFERENCES:

1. Rainer Maria Rilke's poetry
2. Netflix
3. Underworld concert
4. Katy Perry's theatrical performances
5. Taylor Swift's performances
6. Bible study
7. Atlas Shrugged by Ayn Rand
8. Robert Pirsig's writings
9. Bertrand Russell's definition of philosophy
10. Nietzsche's walks
```

## Structure

<img width="2070" alt="fabric_mill_architecture" src="https://github.com/danielmiessler/fabric/assets/50654/ec3bd9b5-d285-483d-9003-7a8e6d842584">

There are multiple ways to use `fabric`.

1. You can just use the `/patterns` in other AI applications/websites.
2. You can build your own server and host these patterns there (plus your own) using the Mill code in `/infrastructure/server`.
3. You can use the `fabric` client to switch between these (coming soon)!

If you are self-hosting your own Mill, the image above shows you what's going on. Basically, you are sending your input to your Fabric Mill, and your Fabric Mill then sends the input and the pattern on to OpenAI. Local model options also being added soon.

## CLI-native

One of the coolest parts of the project is that it's **command-line native**!

Each pattern (prompt) you see in the `/patterns` directory can be used in any AI application you use, but you can also set up your own server using the `/server` code and then call APIs directly!

Once you're set up, you can do things like:
```bash
# Take any idea from `stdin` and send it to the `/write_essay` API!
cat "An idea that coding is like speaking with rules." | write_essay
```

## Naming

Fabric is themed off of, well… _fabric_—as in…woven materials. So, think blankets, quilts, patterns, etc. Here's the concept and structure:

- The project itself is called **Fabric**, and it's the parent concept.
- Individual AI modules (think prompts) are called **Patterns**.
- Chaining together _Patterns_ to create advanced functionality is called a **Stitch**.
- The optional server-side functionality of `fabric` is called the **Mill**.
- The optional client-side scripts within `fabric` are called **Looms**.


## More Documentation

> [!IMPORTANT]\
> We are pushing hard to add lots more functionality and documentation. Please be patient and let us know what you'd like to see in Issues. Thank you!

## Meta

`fabric` was created by <a href="https://danielmiessler.com/" target="_blank">Daniel Miessler</a> in January of 2024.

Special thanks to the following people for inspiration and contributions.

- **Caleb Sima** for pushing me over the edge of whether to make this a public project or not.
- **Joel Parish** for super useful input on the project's Github directory structure.
- **Jonathan Dunn** for spectacular work on the soon-to-be-released universal client.
