# Rate AI Result

This is an example of a Fabric Stitch, which is a chained Fabric command that pipes Fabric results into each other to achieve a result. So it's multiple Patterns…*stitched* together.

## Problem

The problem we're trying to solve with this Stitch is not being able to tell how smart given AI models are. I want to be able to rate their output vs. the output from a different model with the same instructions.

## Solution

What `rate_ai_result` does is run a result using AI 1, and then rate it with AI 2.

## Functionality

`rate_ai_result` accomplishes that like so:

1. Get the input that will be operated on by an AI.
2. Get the instruction/pattern/prompt that will be used by the AI.
3. Get the result of the instructions running against the AI.
4. Combine all three of those together as the input to another Fabric call.
4. Send that combined input to the most advanced model you have available to assess the quality of the AI result.

```
(echo "beginning of content input" ; f -u https://danielmiessler.com/p/framing-is-everything ; echo "end ofcontent input"; echo "beginning of AI instructions (prompt)"; cat ~/.config/fabric/patterns/extract_insights/system.md; echo "end of AI instructions (prompt)" ; echo "beginning of AI output" ; f -u https://danielmiessler.com/p/framing-is-everything | f -p extract_insights -m gpt-3.5-turbo ; echo "end of AI output. Now you should have all three." ) | f -rp rate_ai_result -m o1-preview-2024-09-12
```
In this case we're taking:

* A blog post as the input
* Getting the content of the extract_insights pattern
* Capturing the output of extract_insights on the blog post using `gpt-3.5-turbo`
* Sending all of that to `o1-preview` using the `rate_ai_result` prompt

NOTE: `rate_ai_result` is both a Pattern name and the name of this Stitch.

## Output 

The `rate_ai_result` Pattern is designed to judge the output of another AI on a human sophistication scale that roughly maps to educational and world-state achievement, with the assumption that higher stages require higher cognitive ability as well. These are:

- Superhuman
- Best humans in the world
- Ph.D
- Masters
- Bachelors
- High School
- Partially Educated
- Uneducated

## How to run it

To run it, just execute the code in the `rate_ai_result` file in this repository. And adjust the components as desired to change the input, the AI you're testing, and the AI you're using to judge.

### Blog Post

Here's a full blog post describing in even more detail.

[Using the Smartest AI to Rate Other AI](https://danielmiessler.com/p/using-the-smartest-ai-to-rate-other-ai)

#### Credit

Created by Daniel Miessler on November 7th, 2024.
