<div align="center">

<img src="https://beehiiv-images-production.s3.amazonaws.com/uploads/asset/file/2012aa7c-a939-4262-9647-7ab614e02601/extwis-logo-miessler.png?t=1704502975" alt="extwislogo" width="400" height="400"/>

# `/extractwisdom`

`extract-wisdom` is a [Fabric](https://github.com/danielmiessler/fabric) pattern that _extracts wisdom_ from any text.

<br />

[Description](#description) •
[Installation](#installation) •
[Options](#options) •
[Example](#example) •
[Meta](#meta)

</div>

(NOTE: This is sample readme.md content below. Do not do anything with this yet.)

## Description

`extract-wisdom` (`extwis`) is a [Fabric](https://github.com/danielmiessler/fabric) pattern that _extracts wisdom_ from any text. For example:

- Podcast transcripts
- Academic papers
- Essays
- Blog posts
- Anything you can get into text!

## Installation

You can install `extractwisdom` by executing this command.

```sh
curl -sS https://raw.githubusercontent.com/danielmiessler/fabric/main/install.sh | bash
```

## Configuration

### Flags

When calling `extractwisdom`, the following flags are available:

- `--cmd`
  - Changes the prefix of the `z` and `zi` commands.
  - `--cmd j` would change the commands to (`j`, `ji`).
  - `--cmd cd` would replace the `cd` command (doesn't work on Nushell / POSIX shells).
- `--hook <HOOK>`
  - Changes how often zoxide increments a directory's score:
    | Hook | Description |
    | -------- | --------------------------------- |
    | `none` | Never |
    | `prompt` | At every shell prompt |
    | `pwd` | Whenever the directory is changed |
- `--no-cmd`
  - Prevents zoxide from defining the `z` and `zi` commands.
  - These functions will still be available in your shell as `__zoxide_z` and
    `__zoxide_zi`, should you choose to redefine them.

## Options

## Example

```sh
z foo              # cd into highest ranked directory matching foo
z foo bar          # cd into highest ranked directory matching foo and bar
z foo /            # cd into a subdirectory starting with foo

z ~/foo            # z also works like a regular cd command
z foo/             # cd into relative path
z ..               # cd one level up
z -                # cd into previous directory

zi foo             # cd with interactive selection (using fzf)

z foo<SPACE><TAB>  # show interactive completions (zoxide v0.8.0+, bash 4.4+/fish/zsh only)
```

## Meta

- Author: Daniel Miessler
- Published: January 5, 2024

Read more about the matching algorithm [here][algorithm-matching].
