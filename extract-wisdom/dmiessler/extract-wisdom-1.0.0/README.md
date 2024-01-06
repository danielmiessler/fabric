<div align="center">

<img src="https://beehiiv-images-production.s3.amazonaws.com/uploads/asset/file/2012aa7c-a939-4262-9647-7ab614e02601/extwis-logo-miessler.png?t=1704502975" alt="extwislogo" width="400" height="400"/>

# `/extractwisdom`

`extract-wisdom` is a [Fabric](https://github.com/danielmiessler/fabric) pattern that _extracts wisdom_ from any text.

<br />

[Description](#description) •
[Functionality](#functionality) •
[Installation](#installation) •
[Examples](#examples) •
[Meta](#meta)

</div>

<br />

> ⚠️ These are not the correct instructions yet, so hold off until this note is removed.

## Description

**`extractwisdom` addresses the problem of **too much content** and too little time.**

Not only that, but it's also too easy to forget the stuff do read, watch, or listen to.

The tool _extracts wisdom_ from any content that can be translated into text, for example:

- Podcast transcripts
- Academic papers
- Essays
- Blog posts
- Really, anything you can get into text!

## Functionality

When you use `extractwisdom`, it pulls the following content from the input.

- `IDEAS`
  - Extracts the best ideas from the content, i.e., what you might have taken notes on if you were doing so manually.
- `QUOTES`
  - Some of the best quotes from the content.
- `REFERENCES`
  - External writing, art, and other content referenced positively during the content that might be worth following up on.
- `HABITS`
  - Habits of the speakers that could be worth replicating.
- `RECOMMENDATIONS`
  - A list of things that the content recommends Habits of the speakers.

## Use cases

The output can help you in multiple ways, including:

1. **Decision Support:**<br />
   Allows you to quickly learn what's in a particular piece of content so you can decide if you want to consume the full source material or not.
2. **Note Taking:**<br />
   Can also be used as your note taker. It's designed to replicate the type of capture that you might have done if you took notes manually.

## Examples

Here are some examples of how to use `extractwisdom`.

## Installation

You can install `extractwisdom` by executing this command.

```sh
curl -sS https://raw.githubusercontent.com/danielmiessler/fabric/main/install.sh | bash
```

```sh
z foo              # cd into highest ranked directory matching foo
```

```sh
z foo bar          # cd into highest ranked directory matching foo and bar
```

## Meta

- Author: Daniel Miessler
- Published: January 5, 2024
