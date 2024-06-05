# Getting Started

## Table of Contents
- [Using Pipe and Helper Apps](#using-pipe-and-helper-apps)
  - [Pipes](#pipes)
  - [Helper Apps](#helper-apps)
  - [pbpaste](#pbpaste)
- [Using Fabric with GUI](#using-fabric-with-gui)
- [API Keys](#api-keys)
  - [YouTube API](#youtube-api)
- [Running AI Models with Ollama](#running-ai-models-with-ollama)
  - [Ollama Local](#ollama-local)
  - [Ollama Server](#ollama-server)




## Using Pipe and Helper Apps

### Pipes
Pipes are used to pass the output of one command as the input to another command. This is useful when you want to chain multiple commands together. For example, you can use the `|` pipe to pass the output of the `ls` command to the `grep` command to search for a specific file. 

```bash
ls | grep "file.txt"
```

Fabric supports pipes, so you can use them to chain multiple commands together. For example, you can use the `|` pipe to pass the output of the `yt` command to the `fabric summarize` command to summarize the transcript of a YouTube video.

```bash
yt --transcrip https://www.youtube.com/watch?v=VIDEOID | fabric summarize
```

### Helper Apps
Helper apps are small programs that help you perform specific tasks. Fabric comes with a set of helper apps that you can use to enhance your workflow. For example, you can use the `yt --transcript` helper app to pull the transcript of a YouTube video.

```bash
yt --transcript https://www.youtube.com/watch?v=VIDEOID
```

or you can get the transcript of a audio file using the `ts` helper app.

```bash
ts audio.mp3
```

For more information on the helper apps available in Fabric, check the [Helper Apps](helper-apps.md) page.

### pbpaste
`pbpaste` is a command-line utility that allows you to access the contents of the clipboard on macOS. `pbpaste` prints the contents of the clipboard to the standard output combined with the `|` pipe, you can use `pbpaste` to pass the contents of the clipboard as input to another command.

Daniel Miessler uses `p` as an alias for `pbpaste`.


`pbpaste` is not available on Linux, but you can install it using the `xsel` package. Network Chuck has a great video on how to set up `pbpaste` on Linux. You can watch it on [YouTube](https://www.youtube.com/watch?v=aMzdeZ8vGXQ).


&nbsp;




## Using Fabric with GUI
For all the users who prefer a GUI over the command line, Fabric provides a web interface that you can use to interact with the tool. To use the web interface, follow the steps below:

```bash	
fabric --gui
```

&nbsp;



## API Keys
To use Fabric, you need an AI Model to talk to. If you don't have the nessessary Hardware to [run a model locally](#running-ai-models-with-ollama), you can use [OpenAI's API](https://openai.com/api/), [Google's API](https://ai.google.dev/), or any other API that Fabric supports.

### YouTube API
The YouTube API is used for the `yt` command in Fabric. This command allows you to pull transcripts and to pull user comments.

#### Getting a YouTube API Key
To use the YouTube API, you need to get an API key from Google. You can get one by following this tutorial: [YouTube Data API Overview](https://developers.google.com/youtube/v3/getting-started).

&nbsp;



## Running AI Models with [Ollama](https://ollama.com/)
Ollama allows you to run AI models locally or on a server and integrate them with Fabric. Follow the instructions below to get started.

### Ollama Local
To run AI models locally and use them with Fabric, follow these steps:

1. **Install Ollama**: Download and install Ollama from [Ollama's website](https://ollama.com/).
2. **Download a Model**: Download a model of your choice from Ollama's website.
3. **Run Fabric with the Model**: Execute the following command in your terminal, replacing `NAME_OF_MODEL` with the name of your chosen model.

```bash
fabric --model NAME_OF_MODEL

# Example
fabric --model phi3:latest
```   

**You may need to go through additional steps to set up Ollama and configure it to run on your GPU.**


### Ollama Server
1. **Install Ollama**: Download and install Ollama from [Ollama's website](https://ollama.com/).
2. **Download a Model**: Choose and download the AI model you want to use.
3. **Run Fabric with the Model on a Server**: Execute the following command in your terminal, replacing `NAME_OF_MODEL` with the name of your chosen model and `SERVER_IP` with the IP address of your server.

```bash
fabric --model NAME_OF_MODEL --server SERVER_IP

# Example
fabric --model phi3:latest --server 192.168.1.10
```

**You may need to go through additional steps to set up Ollama and configure it to run on your GPU.**