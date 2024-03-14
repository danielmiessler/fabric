# Fabric Helpers

These are helper tools to work with Fabric. Examples include things like getting transcripts from media files, getting metadata about media, etc.

## yt (YouTube)

`yt` is a command that uses the YouTube API to pull transcripts, get video duration, and other functions. It's primary function is to get a transcript from a video that can then be stitched (piped) into other Fabric Patterns.

```bash
usage: yt [-h] [--duration] [--transcript] [url]

vm (video meta) extracts metadata about a video, such as the transcript and the video's duration. By Daniel Miessler.

positional arguments:
  url           YouTube video URL

options:
  -h, --help    show this help message and exit
  --duration    Output only the duration
  --transcript  Output only the transcript
```

## ts (Audio transcriptions)

'ts' is a command that uses the OpenApi Whisper API to transcribe audio files. Due to the context window, this tool uses pydub to split the files into 10 minute segments. for more information on pydub, please refer https://github.com/jiaaro/pydub

### installation

```bash

mac:
brew install ffmpeg

linux:
apt install ffmpeg

windows:
download instructions https://www.ffmpeg.org/download.html
```

````bash
ts -h
usage: ts [-h] audio_file

Transcribe an audio file.

positional arguments:
  audio_file  The path to the audio file to be transcribed.

options:
  -h, --help  show this help message and exit

## save

`save` is a "tee-like" utility to pipeline saving of content, while keeping the output stream intact. Can optionally generate "frontmatter" for PKM utilities like Obsidian via the
"FABRIC_FRONTMATTER" environment variable



If you'd like to default variables, set them in `~/.config/fabric/.env`. `FABRIC_OUTPUT_PATH` needs to be set so `save` where to write. `FABRIC_FRONTMATTER_TAGS` is optional, but useful for tracking how tags have entered your PKM, if that's important to you.

### usage
```bash
usage: save [-h] [-t, TAG] [-n] [-s] [stub]

save: a "tee-like" utility to pipeline saving of content, while keeping the output stream intact. Can optionally generate "frontmatter" for PKM utilities like Obsidian via the
"FABRIC_FRONTMATTER" environment variable

positional arguments:
  stub                stub to describe your content. Use quotes if you have spaces. Resulting format is YYYY-MM-DD-stub.md by default

options:
  -h, --help          show this help message and exit
  -t, TAG, --tag TAG  add an additional frontmatter tag. Use this argument multiple timesfor multiple tags
  -n, --nofabric      don't use the fabric tags, only use tags from --tag
  -s, --silent        don't use STDOUT for output, only save to the file
````

### example

```bash
echo test | save --tag extra-tag stub-for-name
test

$ cat ~/obsidian/Fabric/2024-03-02-stub-for-name.md
---
generation_date: 2024-03-02 10:43
tags: fabric-extraction stub-for-name extra-tag
---
test
```
