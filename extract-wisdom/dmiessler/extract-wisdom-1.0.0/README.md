<div align="center">

<img src="https://beehiiv-images-production.s3.amazonaws.com/uploads/asset/file/2012aa7c-a939-4262-9647-7ab614e02601/extwis-logo-miessler.png?t=1704502975" alt="extwislogo" width="400" height="400"/>

# `/extractwisdom`

`extract-wisdom` (`extwis`) is a [Fabric](https://github.com/danielmiessler/fabric) pattern that _extracts wisdom_ from any text.

[![crates.io][crates.io-badge]][crates.io]
[![Downloads][downloads-badge]][releases]
[![Built with Nix][builtwithnix-badge]][builtwithnix]

<br />

[Description](#description) •
[Installation](#installation) •
[Options](#options) •
[Example](#example) •
[Meta](#meta)

</div>

## Description

`extract-wisdom` (`extwis`) is a [Fabric](https://github.com/danielmiessler/fabric) pattern that _extracts wisdom_ from any text. For example:

- Podcast transcripts
- Academic papers
- Essays
- Blog posts
- Anything you can get into text!

## Installation

## Options

## Meta

- Author: Daniel Miessler
- Published: January 5, 2024

## Example

![Tutorial][tutorial]

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

Read more about the matching algorithm [here][algorithm-matching].

## Installation

zoxide can be installed in 4 easy steps:

1. **Install binary**

   zoxide runs on most major platforms. If your platform isn't listed below,
   please [open an issue][issues].

   <details>
   <summary>Linux</summary>

   > The recommended way to install zoxide is via the install script:
   >
   > ```sh
   > curl -sS https://raw.githubusercontent.com/ajeetdsouza/zoxide/main/install.sh | bash
   > ```
   >
   > Or, you can use a package manager:
   >
   > | Distribution        | Repository              | Instructions                                                                                          |
   > | ------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------- |
   > | **_Any_**           | **[crates.io]**         | `cargo install zoxide --locked`                                                                       |
   > | _Any_               | [asdf]                  | `asdf plugin add zoxide https://github.com/nyrst/asdf-zoxide.git` <br /> `asdf install zoxide latest` |
   > | _Any_               | [conda-forge]           | `conda install -c conda-forge zoxide`                                                                 |
   > | _Any_               | [Linuxbrew]             | `brew install zoxide`                                                                                 |
   > | _Any_               | [nixpkgs]               | `nix-env -iA nixpkgs.zoxide`                                                                          |
   > | Alpine Linux 3.13+  | [Alpine Linux Packages] | `apk add zoxide`                                                                                      |
   > | Arch Linux          | [Arch Linux Extra]      | `pacman -S zoxide`                                                                                    |
   > | CentOS 7+           | [Copr]                  | `dnf copr enable atim/zoxide` <br /> `dnf install zoxide`                                             |
   > | Debian 11+[^1]      | [Debian Packages]       | `apt install zoxide`                                                                                  |
   > | Devuan 4.0+[^1]     | [Devuan Packages]       | `apt install zoxide`                                                                                  |
   > | Fedora 32+          | [Fedora Packages]       | `dnf install zoxide`                                                                                  |
   > | Gentoo              | [GURU Overlay]          | `eselect repository enable guru` <br /> `emerge --sync guru` <br /> `emerge app-shells/zoxide`        |
   > | Manjaro             |                         | `pacman -S zoxide`                                                                                    |
   > | openSUSE Tumbleweed | [openSUSE Factory]      | `zypper install zoxide`                                                                               |
   > | Parrot OS[^1]       |                         | `apt install zoxide`                                                                                  |
   > | Raspbian 11+[^1]    | [Raspbian Packages]     | `apt install zoxide`                                                                                  |
   > | Rhino Linux         | [Pacstall Packages]     | `pacstall -I zoxide-deb`                                                                              |
   > | Slackware 15.0+     | [SlackBuilds]           | [Instructions][slackbuilds-howto]                                                                     |
   > | Ubuntu 21.04+[^1]   | [Ubuntu Packages]       | `apt install zoxide`                                                                                  |
   > | Void Linux          | [Void Linux Packages]   | `xbps-install -S zoxide`                                                                              |

   </details>

   <details>
   <summary>macOS</summary>

   > To install zoxide, use a package manager:
   >
   > | Repository      | Instructions                                                                                          |
   > | --------------- | ----------------------------------------------------------------------------------------------------- |
   > | **[crates.io]** | `cargo install zoxide --locked`                                                                       |
   > | **[Homebrew]**  | `brew install zoxide`                                                                                 |
   > | [asdf]          | `asdf plugin add zoxide https://github.com/nyrst/asdf-zoxide.git` <br /> `asdf install zoxide latest` |
   > | [conda-forge]   | `conda install -c conda-forge zoxide`                                                                 |
   > | [MacPorts]      | `port install zoxide`                                                                                 |
   > | [nixpkgs]       | `nix-env -iA nixpkgs.zoxide`                                                                          |
   >
   > Or, run this command in your terminal:
   >
   > ```sh
   > curl -sS https://raw.githubusercontent.com/ajeetdsouza/zoxide/main/install.sh | bash
   > ```

   </details>

   <details>
   <summary>Windows</summary>

   > The recommended way to install zoxide is via `winget`:
   >
   > ```sh
   > winget install ajeetdsouza.zoxide
   > ```
   >
   > Or, you can use an alternative package manager:
   >
   > | Repository      | Instructions                          |
   > | --------------- | ------------------------------------- |
   > | **[crates.io]** | `cargo install zoxide --locked`       |
   > | [Chocolatey]    | `choco install zoxide`                |
   > | [conda-forge]   | `conda install -c conda-forge zoxide` |
   > | [Scoop]         | `scoop install zoxide`                |
   >
   > If you're using Cygwin, Git Bash, or MSYS2, use the install script instead:
   >
   > ```sh
   > curl -sS https://raw.githubusercontent.com/ajeetdsouza/zoxide/main/install.sh | bash
   > ```

   </details>

   <details>
   <summary>BSD</summary>

   > To install zoxide, use a package manager:
   >
   > | Distribution  | Repository      | Instructions                    |
   > | ------------- | --------------- | ------------------------------- |
   > | **_Any_**     | **[crates.io]** | `cargo install zoxide --locked` |
   > | DragonFly BSD | [DPorts]        | `pkg install zoxide`            |
   > | FreeBSD       | [FreshPorts]    | `pkg install zoxide`            |
   > | NetBSD        | [pkgsrc]        | `pkgin install zoxide`          |

   </details>

   <details>
   <summary>Android</summary>

   > To install zoxide, use a package manager:
   >
   > | Repository | Instructions         |
   > | ---------- | -------------------- |
   > | [Termux]   | `pkg install zoxide` |

   </details>

## Configuration

### Flags

When calling `zoxide init`, the following flags are available:

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

### Environment variables

Environment variables[^2] can be used for configuration. They must be set before
`zoxide init` is called.

- `_ZO_DATA_DIR`
  - Specifies the directory in which the database is stored.
  - The default value varies across OSes:
    | OS | Path | Example |
    | ----------- | ---------------------------------------- | ------------------------------------------ |
    | Linux / BSD | `$XDG_DATA_HOME` or `$HOME/.local/share` | `/home/alice/.local/share` |
    | macOS | `$HOME/Library/Application Support` | `/Users/Alice/Library/Application Support` |
    | Windows | `%LOCALAPPDATA%` | `C:\Users\Alice\AppData\Local` |
- `_ZO_ECHO`
  - When set to 1, `z` will print the matched directory before navigating to
    it.
- `_ZO_EXCLUDE_DIRS`
  - Excludes the specified directories from the database.
  - This is provided as a list of [globs][glob], separated by OS-specific
    characters:
    | OS | Separator | Example |
    | ------------------- | --------- | ----------------------- |
    | Linux / macOS / BSD | `:` | `$HOME:$HOME/private/*` |
    | Windows | `;` | `$HOME;$HOME/private/*` |
  - By default, this is set to `"$HOME"`.
- `_ZO_FZF_OPTS`
  - Custom options to pass to [fzf] during interactive selection. See
    [`man fzf`][fzf-man] for the list of options.
- `_ZO_MAXAGE`
  - Configures the [aging algorithm][algorithm-aging], which limits the maximum
    number of entries in the database.
  - By default, this is set to 10000.
- `_ZO_RESOLVE_SYMLINKS`
  - When set to 1, `z` will resolve symlinks before adding directories to the
    database.
