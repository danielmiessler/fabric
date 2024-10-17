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

[Updates](#updates) •
[What and Why](#whatandwhy) •
[Philosophy](#philosophy) •
[Installation](#Installation) •
[Usage](#Usage) •
[Examples](#examples) •
[Just Use the Patterns](#just-use-the-patterns) •
[Custom Patterns](#custom-patterns) •
[Helper Apps](#helper-apps) •
[Meta](#meta)

</div>

## Navigation

- [Updates](#updates)
- [What and Why](#what-and-why)
- [Philosophy](#philosophy)
  - [Breaking problems into components](#breaking-problems-into-components)
  - [Too many prompts](#too-many-prompts)
  - [The Fabric approach to prompting](#our-approach-to-prompting)
- [Installation](#Installation)
  - [Migration](#Migration)
  - [Upgrading](#Upgrading)
- [Usage](#Usage)
- [Examples](#examples)
  - [Just use the Patterns](#just-use-the-patterns)
- [Custom Patterns](#custom-patterns)
- [Helper Apps](#helper-apps)
- [Meta](#meta)
  - [Primary contributors](#primary-contributors)

<br />

## Updates

> [!NOTE] 
September 15, 2024 — Lots of new stuff!
> * Fabric now supports calling the new `o1-preview` model using the `-r` switch (which stands for raw. Normal queries won't work with `o1-preview` because they disabled System access and don't allow us to set `Temperature`.
> * We have early support for Raycast! Under the `/patterns` directory there's a `raycast` directory with scripts that can be called from Raycast. If you add a scripts directory within Raycast and point it to your `~/.config/fabric/patterns/raycast` directory, you'll then be able to 1) invoke Raycast, type the name of the script, and then 2) paste in the content to be passed, and the results will return in Raycast. There's currently only one script in there but I am (Daniel) adding more.
> * **Go Migration: The following command line options were changed during the migration to Go:**
> * You now need to use the -c option instead of -C to copy the result to the clipboard.
> * You now need to use the -s option instead of -S to stream results in realtime.
> * The following command line options have been removed `--agents` (-a), `--gui`, `--clearsession`, `--remoteOllamaServer`, and `--sessionlog`
> * You can now use (-S) to configure an Ollama server.
> * **We're working on a GUI rewrite in Go as well**

## Intro videos

Keep in mind that many of these were recorded when Fabric was Python-based, so remember to use the current [install instructions](#Installation) below.

* [Network Chuck](https://www.youtube.com/watch?v=UbDyjIIGaxQ)
* [David Bombal](https://www.youtube.com/watch?v=vF-MQmVxnCs)
* [My Own Intro to the Tool](https://www.youtube.com/watch?v=wPEyyigh10g)
* [More Fabric YouTube Videos](https://www.youtube.com/results?search_query=fabric+ai)

## What and why

Since the start of 2023 and GenAI we've seen a massive number of AI applications for accomplishing tasks. It's powerful, but _it's not easy to integrate this functionality into our lives._

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
- Creating social media posts from any content input
- And a million more…

## Installation

To install Fabric, you can use the latest release binaries or install it from the source.

### Get Latest Release Binaries

```bash
# Windows:
curl -L https://github.com/danielmiessler/fabric/releases/latest/download/fabric-windows-amd64.exe > fabric.exe && fabric.exe --version

# MacOS (arm64):
curl -L https://github.com/danielmiessler/fabric/releases/latest/download/fabric-darwin-arm64 > fabric && chmod +x fabric && ./fabric --version

# MacOS (amd64): 
curl -L https://github.com/danielmiessler/fabric/releases/latest/download/fabric-darwin-amd64 > fabric && chmod +x fabric && ./fabric --version

# Linux (amd64): 
curl -L https://github.com/danielmiessler/fabric/releases/latest/download/fabric-linux-amd64 > fabric && chmod +x fabric && ./fabric --version

# Linux (arm64):
curl -L https://github.com/danielmiessler/fabric/releases/latest/download/fabric-linux-arm64 > fabric && chmod +x fabric && ./fabric --version
```

### From Source

To install Fabric, [make sure Go is installed](https://go.dev/doc/install), and then run the following command.

```bash
# Install Fabric directly from the repo
go install github.com/danielmiessler/fabric@latest
```

### Environment Variables

You may need to set some environment variables in your `~/.bashrc` on linux or `~/.zshrc` file on mac to be able to run the `fabric` command. Here is an example of what you can add:

For Intel based macs or linux
```bash
# Golang environment variables
export GOROOT=/usr/local/go
export GOPATH=$HOME/go

# Update PATH to include GOPATH and GOROOT binaries
export PATH=$GOPATH/bin:$GOROOT/bin:$HOME/.local/bin:$PATH
```

for Apple Silicon based macs
```bash
# Golang environment variables
export GOROOT=$(brew --prefix go)/libexec
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$HOME/.local/bin:$PATH
```

### Setup
Now run the following command
```bash
# Run the setup to set up your directories and keys
fabric --setup
```
If everything works you are good to go.


### Migration

If you have the Legacy (Python) version installed and want to migrate to the Go version, here's how you do it. It's basically two steps: 1) uninstall the Python version, and 2) install the Go version.

```bash
# Uninstall Legacy Fabric
pipx uninstall fabric

# Clear any old Fabric aliases
(check your .bashrc, .zshrc, etc.)
# Install the Go version
go install github.com/danielmiessler/fabric@latest
# Run setup for the new version. Important because things have changed
fabric --setup
```

Then [set your environmental variables](#environmental-variables) as shown above.

### Upgrading

The great thing about Go is that it's super easy to upgrade. Just run the same command you used to install it in the first place and you'll always get the latest version.
```bash
go install -ldflags "-X main.version=$(git describe --tags --always)" github.com/danielmiessler/fabric@latest
```

## Usage
Once you have it all set up, here's how to use it.

```bash
fabric -h
```

```bash

Usage:
  fabric [OPTIONS]

Application Options:
  -p, --pattern=             Choose a pattern from the available patterns
  -v, --variable=            Values for pattern variables, e.g. -v=#role:expert -v=#points:30"
  -C, --context=             Choose a context from the available contexts
      --session=             Choose a session from the available sessions
  -S, --setup                Run setup for all reconfigurable parts of fabric
  -t, --temperature=         Set temperature (default: 0.7)
  -T, --topp=                Set top P (default: 0.9)
  -s, --stream               Stream
  -P, --presencepenalty=     Set presence penalty (default: 0.0)
  -r, --raw                  Use the defaults of the model without sending chat options (like temperature etc.) and use the user role instead of the system role for patterns.
  -F, --frequencypenalty=    Set frequency penalty (default: 0.0)
  -l, --listpatterns         List all patterns
  -L, --listmodels           List all available models
  -x, --listcontexts         List all contexts
  -X, --listsessions         List all sessions
  -U, --updatepatterns       Update patterns
  -c, --copy                 Copy to clipboard
  -m, --model=               Choose model
  -o, --output=              Output to file
      --output-session       Output the entire session (also a temporary one) to the output file
  -n, --latest=              Number of latest patterns to list (default: 0)
  -d, --changeDefaultModel   Change default model
  -y, --youtube=             YouTube video "URL" to grab transcript, comments from it and send to chat
      --transcript           Grab transcript from YouTube video and send to chat (it used per default).
      --comments             Grab comments from YouTube video and send to chat
  -g, --language=            Specify the Language Code for the chat, e.g. -g=en -g=zh
  -u, --scrape_url=          Scrape website URL to markdown using Jina AI
  -q, --scrape_question=     Search question using Jina AI
  -e, --seed=                Seed to be used for LMM generation
  -w, --wipecontext=         Wipe context
  -W, --wipesession=         Wipe session
      --printcontext=        Print context
      --printsession=        Print session
      --readability          Convert HTML input into a clean, readable view
      --dry-run              Show what would be sent to the model without actually sending it
      --version              Print current version

Help Options:
  -h, --help                 Show this help message

```

## Our approach to prompting

Fabric _Patterns_ are different than most prompts you'll see.

- **First, we use `Markdown` to help ensure maximum readability and editability**. This not only helps the creator make a good one, but also anyone who wants to deeply understand what it does. _Importantly, this also includes the AI you're sending it to!_

Here's an example of a Fabric Pattern.

```bash
https://github.com/danielmiessler/fabric/blob/main/patterns/extract_wisdom/system.md
```

<img width="1461" alt="pattern-example" src="https://github.com/danielmiessler/fabric/assets/50654/b910c551-9263-405f-9735-71ca69bbab6d">

- **Next, we are extremely clear in our instructions**, and we use the Markdown structure to emphasize what we want the AI to do, and in what order.

- **And finally, we tend to use the System section of the prompt almost exclusively**. In over a year of being heads-down with this stuff, we've just seen more efficacy from doing that. If that changes, or we're shown data that says otherwise, we will adjust.

## Examples

Now let's look at some things you can do with Fabric.

1. Run the `summarize` Pattern based on input from `stdin`. In this case, the body of an article.

```bash
pbpaste | fabric --pattern summarize
```

2. Run the `analyze_claims` Pattern with the `--stream` option to get immediate and streaming results.

```bash
pbpaste | fabric --stream --pattern analyze_claims
```

3. Run the `extract_wisdom` Pattern with the `--stream` option to get immediate and streaming results from any Youtube video (much like in the original introduction video).

```bash
fabric -y "https://youtube.com/watch?v=uXs-zPc63kM" --stream --pattern extract_wisdom
```

4. Create patterns- you must create a .md file with the pattern and save it to ~/.config/fabric/patterns/[yourpatternname].

## Just use the Patterns

<img width="1173" alt="fabric-patterns-screenshot" src="https://github.com/danielmiessler/fabric/assets/50654/9186a044-652b-4673-89f7-71cf066f32d8">

<br />
<br />

If you're not looking to do anything fancy, and you just want a lot of great prompts, you can navigate to the [`/patterns`](https://github.com/danielmiessler/fabric/tree/main/patterns) directory and start exploring!

We hope that if you used nothing else from Fabric, the Patterns by themselves will make the project useful.

You can use any of the Patterns you see there in any AI application that you have, whether that's ChatGPT or some other app or website. Our plan and prediction is that people will soon be sharing many more than those we've published, and they will be way better than ours.

The wisdom of crowds for the win.

## Custom Patterns

You may want to use Fabric to create your own custom Patterns—but not share them with others. No problem!

Just make a directory in `~/.config/custompatterns/` (or wherever) and put your `.md` files in there. 

When you're ready to use them, copy them into:

```
~/.config/fabric/patterns/
```

You can then use them like any other Patterns, but they won't be public unless you explicitly submit them as Pull Requests to the Fabric project. So don't worry—they're private to you.


This feature works with all openai and ollama models but does NOT work with claude. You can specify your model with the -m flag

## Helper Apps

Fabric also makes use of some core helper apps (tools) to make it easier to integrate with your various workflows. Here are some examples:

### `to_pdf`

`to_pdf` is a helper command that converts LaTeX files to PDF format. You can use it like this:

```bash
to_pdf input.tex
```

This will create a PDF file from the input LaTeX file in the same directory.

You can also use it with stdin which works perfectly with the `write_latex` pattern:

```bash
echo "ai security primer" | fabric --pattern write_latex | to_pdf
```

This will create a PDF file named `output.pdf` in the current directory.

### `to_pdf` Installation

To install `to_pdf`, install it the same way as you install Fabric, just with a different repo name.

```bash
go install github.com/danielmiessler/fabric/to_pdf@latest
```

Make sure you have a LaTeX distribution (like TeX Live or MiKTeX) installed on your system, as `to_pdf` requires `pdflatex` to be available in your system's PATH.

## Meta

> [!NOTE]
> Special thanks to the following people for their inspiration and contributions!

- _Jonathan Dunn_ for being the absolute MVP dev on the project, including spearheading the new Go version, as well as the GUI! All this while also being a full-time medical doctor!
- _Caleb Sima_ for pushing me over the edge of whether to make this a public project or not.
- _Eugen Eisler_ and _Frederick Ros_ for their invaluable contributions to the Go version
- _Joel Parish_ for super useful input on the project's Github directory structure..
- _Joseph Thacker_ for the idea of a `-c` context flag that adds pre-created context in the `./config/fabric/` directory to all Pattern queries.
- _Jason Haddix_ for the idea of a stitch (chained Pattern) to filter content using a local model before sending on to a cloud model, i.e., cleaning customer data using `llama2` before sending on to `gpt-4` for analysis.
- _Andre Guerra_ for assisting with numerous components to make things simpler and more maintainable.

### Primary contributors

<a href="https://github.com/danielmiessler"><img src="https://avatars.githubusercontent.com/u/50654?v=4" title="Daniel Miessler" width="50" height="50"></a>
<a href="https://github.com/xssdoctor"><img src="https://avatars.githubusercontent.com/u/9218431?v=4" title="Jonathan Dunn" width="50" height="50"></a>
<a href="https://github.com/sbehrens"><img src="https://avatars.githubusercontent.com/u/688589?v=4" title="Scott Behrens" width="50" height="50"></a>
<a href="https://github.com/agu3rra"><img src="https://avatars.githubusercontent.com/u/10410523?v=4" title="Andre Guerra" width="50" height="50"></a>

`fabric` was created by <a href="https://danielmiessler.com/subscribe" target="_blank">Daniel Miessler</a> in January of 2024.
<br /><br />
<a href="https://twitter.com/intent/user?screen_name=danielmiessler">![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/danielmiessler)</a>
