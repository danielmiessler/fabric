# Fabric Helpers

These are helper tools to work with Fabric. Examples include things like getting transcripts from media files, getting metadata about media, etc.

## yt (YouTube)

`yt` is a command that uses the YouTube API to pull transcripts, get video duration, and other functions. It's primary function is to get a transcript from a video that can then be stitched (piped) into other Fabric Patterns.

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

```bash
ts -h
usage: ts [-h] audio_file

Transcribe an audio file.

positional arguments:
  audio_file  The path to the audio file to be transcribed.

options:
  -h, --help  show this help message and exit
```
