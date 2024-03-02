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
