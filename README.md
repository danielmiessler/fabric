<div align="center">

<img src="./images/fabric-logo-gif.gif" alt="fabriclogo" width="400" height="400"/>

# `fabric`

![Static Badge](https://img.shields.io/badge/mission-human_flourishing_via_AI_augmentation-purple)
<br />
![GitHub top language](https://img.shields.io/github/languages/top/danielmiessler/fabric)
![GitHub last commit](https://img.shields.io/github/last-commit/danielmiessler/fabric)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

<p class="align center">
<h4><code>fabric</code> is an open-source framework for augmenting humans using AI.</h4>
</p>

[What and Why](#whatandwhy) •
[Philosophy](#philosophy) •
[Quickstart](#quickstart) •
[Structure](#structure) •
[Examples](#examples) •
[Meta](#meta)

</div>

> [!NOTE]  
> We're still working on both code and documentation. Please be patient with us as we roll things out in the next week.

## What and why

Since the start of 2023 and GenAI we've seen a massive number of AI applications for accomplishing tasks. It's powerful, but **it's not easy to integrate this functionality into our lives.**

<div align="center">
<h4>In other words, AI doesn't have a capabilities problem—it has an <em>integration</em> problem.</h4>
</div>

Fabric was created to address this by enabling everyone to granularly apply AI to everyday challenges.

## Philosophy

> AI isn't a thing; it's a _magnifier_ of a thing. And that thing is **human creativity**.

We believe the purpose of technology is to help humans flourish, so when we talk about AI we start with the **human** problems we want to solve.

### Breaking problems into components

Our approach is to break problems into individual pieces (see below) and then apply AI to them one at a time. See below for some examples.

<img width="2078" alt="augmented_challenges" src="https://github.com/danielmiessler/fabric/assets/50654/31997394-85a9-40c2-879b-b347e4701f06">

### Too many prompts

Prompts are good for this, but the biggest challenge I faced in 2023——which still exists today—is **the sheer number of AI prompts out there**. We all have prompts that are useful, but it's hard to discover new ones, know if they are good or not, _and manage different versions of the ones we like_.

One of <code>fabric</code>'s primary features is helping people collect and integrate prompts, which we call _Patterns_, into various parts of their lives.

Fabric has Patterns for all sorts of life and work activities, including:

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

### Our approach to prompting

Fabric _Patterns_ are different than most prompts you'll see. 

- **First, we use `Markdown` to help ensure maximum readability and editability**. This not only helps the creator make a good one, but also anyone who wants to deeply understand what it does. _Importantly, this also includes the AI you're sending it to!_

Here's an example of a Fabric Pattern.

```bash
https://github.com/danielmiessler/fabric/blob/main/patterns/extract_wisdom/system.md
```

<img width="1461" alt="pattern-example" src="https://github.com/danielmiessler/fabric/assets/50654/b910c551-9263-405f-9735-71ca69bbab6d">

- **Next, we are extremely clear in our instructions**, and we use the Markdown structure to emphasize what we want the AI to do, and in what order.

- **And finally, we tend to use the System section of the prompt almost exclusively**. In over a year of being heads-down with this stuff, we've just seen more efficacy from doing that. If that changes, or we're shown data that says otherwise, we will adjust.

## Quickstart

There are three (3) main ways to get started with Fabric.

<img width="1173" alt="fabric-patterns-screenshot" src="https://github.com/danielmiessler/fabric/assets/50654/9186a044-652b-4673-89f7-71cf066f32d8">

### 1. Just use the Patterns

If you're not looking to do anything fancy, and you just want a lot of great prompts, you can navigate to the [`/patterns`](https://github.com/danielmiessler/fabric/tree/main/patterns) directory and start exploring!

We hope that if you used nothing else from Fabric, the Patterns by themselves will make the project useful.

You can use any of the Patterns you see there in any AI application that you have, whether that's ChatGPT or some other app or website. Our plan and prediction is that people will soon be sharing many more than those we've published, and they will be way better than ours. 

The wisdom of crowds for the win.

### 2. Create your own Fabric Mill (Server)

<img width="2070" alt="fabric_mill_architecture" src="https://github.com/danielmiessler/fabric/assets/50654/ec3bd9b5-d285-483d-9003-7a8e6d842584">

<br />

But we go beyond just providing Patterns. We provide code for you to build your very own Fabric server and personal AI infrastructure! 

To get started, head over to the [`/server/`](https://github.com/danielmiessler/fabric/tree/main/server) directory and set up your own Fabric Mill with your own Patterns running! You can then use the [`/client/standalone_client_examples`](https://github.com/danielmiessler/fabric/tree/main/client/standalone_client_examples) to connect to it.

### 3. The standalone client

We're almost done with a `fabric` client that will let you do all sorts of cool stuff, including:

- Calling patterns without connecting to a Fabric server (direct to OpenAI).
- Streaming mode to get instant and dynamic results.
- Other cool stuff…

We expect this client to be ready very within a week or so, and we'll update the Quickstart as soon as it is.

## Structure

Fabric is themed off of, well… _fabric_—as in…woven materials. So, think blankets, quilts, patterns, etc. Here's the concept and structure:

### Components

The Fabric ecosystem has three primary components, all named within this textile theme.

- The **Mill** is the (optional) server that makes **Patterns** available.
- **Patterns** are the actual granular AI use cases (prompts).
- **Stitches** are chained together _Patterns_ that create advanced functionality (see below).
- **Looms** are the client-side apps that call a specific **Pattern** hosted by a **Mill**.

### CLI-native

One of the coolest parts of the project is that it's **command-line native**!

Each Pattern you see in the `/patterns` directory can be used in any AI application you use, but you can also set up your own server using the `/server` code and then call APIs directly!

Once you're set up, you can do things like:

```bash
# Take any idea from `stdin` and send it to the `/write_essay` API!
cat "An idea that coding is like speaking with rules." | write_essay
```

### Directly calling Patterns

One key feature of `fabric` and its Markdown-based format is the ability to ** directly reference** (and edit) individual [patterns](https://github.com/danielmiessler/fabric/tree/main#naming) directly—on their own—without surrounding code.

As an example, here's how to call _the direct location_ of the `extract_wisdom` pattern.

```bash
https://github.com/danielmiessler/fabric/blob/main/patterns/extract_wisdom/system.md
```

This means you can cleanly, and directly reference any pattern for use in a web-based AI app, your own code, or wherever!

Even better, you can also have your [Mill](https://github.com/danielmiessler/fabric/tree/main#naming) functionality directly call **system** and **user** prompts from `fabric`, meaning you can have your personal AI ecosystem automatically kept up to date with the latest version of your favorite [Patterns](https://github.com/danielmiessler/fabric/tree/main#naming).

Here's what that looks like in code:

```bash
https://github.com/danielmiessler/fabric/blob/main/server/fabric_api_server.py
```

```python
# /extwis
@app.route("/extwis", methods=["POST"])
@auth_required  # Require authentication
def extwis():
    data = request.get_json()

    # Warn if there's no input
    if "input" not in data:
        return jsonify({"error": "Missing input parameter"}), 400

    # Get data from client
    input_data = data["input"]

    # Set the system and user URLs
    system_url = "https://raw.githubusercontent.com/danielmiessler/fabric/main/patterns/extract_wisdom/system.md"
    user_url = "https://raw.githubusercontent.com/danielmiessler/fabric/main/patterns/extract_wisdom/user.md"

    # Fetch the prompt content
    system_content = fetch_content_from_url(system_url)
    user_file_content = fetch_content_from_url(user_url)

    # Build the API call
    system_message = {"role": "system", "content": system_content}
    user_message = {"role": "user", "content": user_file_content + "\n" + input_data}
    messages = [system_message, user_message]
    try:
        response = openai.chat.completions.create(
            model="gpt-4-1106-preview",
            messages=messages,
            temperature=0.0,
            top_p=1,
            frequency_penalty=0.1,
            presence_penalty=0.1,
        )
        assistant_message = response.choices[0].message.content
        return jsonify({"response": assistant_message})
    except Exception as e:
        return jsonify({"error": str(e)}), 500
```

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

## Meta

`fabric` was created by <a href="https://danielmiessler.com/" target="_blank">Daniel Miessler</a> in January of 2024.

> [!NOTE]  
> Special thanks to the following people for their inspiration and contributions!

- **Caleb Sima** for pushing me over the edge of whether to make this a public project or not.
- **Joel Parish** for super useful input on the project's Github directory structure.
- **Jonathan Dunn** for spectacular work on the soon-to-be-released universal client.

### Contributors

## Contributors

<!---
npx contributor-faces --exclude "*bot*" --limit 70 --repo "https://github.com/amplication/amplication"

change the height and width for each of the contributors from 80 to 50.
--->

[//]: contributor-faces
<a href="https://github.com/sbehrens"><img src="https://avatars.githubusercontent.com/u/43705455?v=4" title="yuval-hazaz" width="50" height="50"></a>
<a href="https://github.com/xssdoctor"><img src="https://avatars.githubusercontent.com/u/12671072?v=4" title="iddan" width="50" height="50"></a>

[//]: contributor-faces

<a href="https://twitter.com/intent/user?screen_name=danielmiessler">![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/danielmiessler)</a>
