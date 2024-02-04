# The `fabric` client

This is the primary `fabric` client, which has multiple modes of operation.

## Client modes

You can use the client in three different modes:

1. **Local Only:** You can use the client without a server, and it will use patterns it's downloaded from this repository, or ones that you specify.
2. **Local Server:** You can run your own version of a Fabric Mill locally (on a private IP), which you can then connect to and use.
3. **Remote Server:** You can specify a remote server that your client commands will then be calling.

## Client features

1. Standalone Mode: Run without needing a server.
2. Clipboard Integration: Copy responses to the clipboard.
3. File Output: Save responses to files for later reference.
4. Pattern Module: Utilize specific patterns for different types of analysis.
5. Server Mode: Operate the tool in server mode to control your own patterns and let your other apps access it.

## Installation

1. If you have this repository downloaded, you already have the client.
   `git clone git@github.com:danielmiessler/fabric.git`
2. Navigate to the client's directory:
   `cd client`
3. Install poetry (if you don't have it already)
   `pip3 install poetry`
4. Install the required packages:
   `poetry install`
5. Activate the virtual environment:
   `poetry shell`
6. Copy to path:
   `echo export PATH=$PATH:$(pwd) >> ~/.bashrc` # or .zshrc
7. Copy your OpenAI API key to the `.env` file in your `nvim ~/.config/fabric/` directory (or create that file and put it in)
   `OPENAI_API_KEY=[Your_API_Key]`

## Usage

To use `fabric`, call it with your desired options:

fabric [options]
Options include:

--pattern, -p: Select the module for analysis.
--stream, -s: Stream output to another application.
--output, -o: Save the response to a file.
--copy, -c: Copy the response to the clipboard.

Example:

```bash
# Pasting in an article about LLMs
pbpaste | fabric --pattern extract_wisdom --output wisdom.txt | fabric --pattern summarize --stream
```

```markdown
ONE SENTENCE SUMMARY:

- The content covered the basics of LLMs and how they are used in everyday practice.

MAIN POINTS:

1. LLMs are large language models, and typically use the transformer architecture.
2. LLMs used to be used for story generation, but they're now used for many AI applications.
3. They are vulnerable to hallucination if not configured correctly, so be careful.

TAKEAWAYS:

1. It's possible to use LLMs for multiple AI use cases.
2. It's important to validate that the results you're receiving are correct.
3. The field of AI is moving faster than ever as a result of GenAI breakthroughs.
```

## Contributing

We welcome contributions to Fabric, including improvements and feature additions to this client.

## Credits

The `fabric` client was created by Jonathan Dunn and Daniel Meissler.
