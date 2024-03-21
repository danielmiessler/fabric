import requests
import os
from openai import OpenAI, APIConnectionError
import asyncio
import pyperclip
import sys
import platform
from dotenv import load_dotenv
import zipfile
import tempfile
import re
import shutil

current_directory = os.path.dirname(os.path.realpath(__file__))
config_directory = os.path.expanduser("~/.config/fabric")
env_file = os.path.join(config_directory, ".env")


class Standalone:
    def __init__(self, args, pattern="", env_file="~/.config/fabric/.env"):
        """        Initialize the class with the provided arguments and environment file.

        Args:
            args: The arguments for initialization.
            pattern: The pattern to be used (default is an empty string).
            env_file: The path to the environment file (default is "~/.config/fabric/.env").

        Returns:
            None

        Raises:
            KeyError: If the "OPENAI_API_KEY" is not found in the environment variables.
            FileNotFoundError: If no API key is found in the environment variables.
        """

        # Expand the tilde to the full path
        if args is None:
            args = type('Args', (), {})()
        env_file = os.path.expanduser(env_file)
        load_dotenv(env_file)
        assert 'OPENAI_API_KEY' in os.environ, "Error: OPENAI_API_KEY not found in environment variables. Please run fabric --setup and add a key."
        api_key = os.environ['OPENAI_API_KEY']
        base_url = os.environ.get(
            'OPENAI_BASE_URL', 'https://api.openai.com/v1/')
        self.client = OpenAI(api_key=api_key, base_url=base_url)
        self.local = False
        self.config_pattern_directory = config_directory
        self.pattern = pattern
        self.args = args
        self.model = getattr(args, 'model', None)
        if not self.model:
            self.model = os.environ.get('DEFAULT_MODEL', None)
            if not self.model:
                self.model = 'gpt-4-turbo-preview'
        self.claude = False
        sorted_gpt_models, ollamaList, claudeList = self.fetch_available_models()
        self.local = self.model in ollamaList
        self.claude = self.model in claudeList

    async def localChat(self, messages, host=''):
        from ollama import AsyncClient
        response = None
        if host:
            response = await AsyncClient(host=host).chat(model=self.model, messages=messages, host=host)
        else:
            response = await AsyncClient().chat(model=self.model, messages=messages)
        print(response['message']['content'])
        copy = self.args.copy
        if copy:
            pyperclip.copy(response['message']['content'])

    async def localStream(self, messages, host=''):
        from ollama import AsyncClient
        if host:
            async for part in await AsyncClient(host=host).chat(model=self.model, messages=messages, stream=True, host=host):
                print(part['message']['content'], end='', flush=True)
        else:
            async for part in await AsyncClient().chat(model=self.model, messages=messages, stream=True):
                print(part['message']['content'], end='', flush=True)

    async def claudeStream(self, system, user):
        from anthropic import AsyncAnthropic
        self.claudeApiKey = os.environ["CLAUDE_API_KEY"]
        Streamingclient = AsyncAnthropic(api_key=self.claudeApiKey)
        async with Streamingclient.messages.stream(
            max_tokens=4096,
            system=system,
            messages=[user],
            model=self.model, temperature=0.0, top_p=1.0
        ) as stream:
            async for text in stream.text_stream:
                print(text, end="", flush=True)
            print()

        message = await stream.get_final_message()

    async def claudeChat(self, system, user, copy=False):
        from anthropic import Anthropic
        self.claudeApiKey = os.environ["CLAUDE_API_KEY"]
        client = Anthropic(api_key=self.claudeApiKey)
        message = client.messages.create(
            max_tokens=4096,
            system=system,
            messages=[user],
            model=self.model,
            temperature=0.0, top_p=1.0
        )
        print(message.content[0].text)
        copy = self.args.copy
        if copy:
            pyperclip.copy(message.content[0].text)

    def streamMessage(self, input_data: str, context="", host=''):
        """        Stream a message and handle exceptions.

        Args:
            input_data (str): The input data for the message.

        Returns:
            None: If the pattern is not found.

        Raises:
            FileNotFoundError: If the pattern file is not found.
        """

        wisdomFilePath = os.path.join(
            config_directory, f"patterns/{self.pattern}/system.md"
        )
        user_message = {"role": "user", "content": f"{input_data}"}
        wisdom_File = os.path.join(current_directory, wisdomFilePath)
        system = ""
        buffer = ""
        if self.pattern:
            try:
                with open(wisdom_File, "r") as f:
                    if context:
                        system = context + '\n\n' + f.read()
                    else:
                        system = f.read()
                    system_message = {"role": "system", "content": system}
                messages = [system_message, user_message]
            except FileNotFoundError:
                print("pattern not found")
                return
        else:
            if context:
                messages = [
                    {"role": "system", "content": context}, user_message]
            else:
                messages = [user_message]
        try:
            if self.local:
                if host:
                    asyncio.run(self.localStream(messages, host=host))
                else:
                    asyncio.run(self.localStream(messages))
            elif self.claude:
                from anthropic import AsyncAnthropic
                asyncio.run(self.claudeStream(system, user_message))
            else:
                stream = self.client.chat.completions.create(
                    model=self.model,
                    messages=messages,
                    temperature=0.0,
                    top_p=1,
                    frequency_penalty=0.1,
                    presence_penalty=0.1,
                    stream=True,
                )
                for chunk in stream:
                    if chunk.choices[0].delta.content is not None:
                        char = chunk.choices[0].delta.content
                        buffer += char
                        if char not in ["\n", " "]:
                            print(char, end="")
                        elif char == " ":
                            print(" ", end="")  # Explicitly handle spaces
                        elif char == "\n":
                            print()  # Handle newlines
                    sys.stdout.flush()
        except Exception as e:
            if "All connection attempts failed" in str(e):
                print(
                    "Error: cannot connect to llama2. If you have not already, please visit https://ollama.com for installation instructions")
            if "CLAUDE_API_KEY" in str(e):
                print(
                    "Error: CLAUDE_API_KEY not found in environment variables. Please run --setup and add the key")
            if "overloaded_error" in str(e):
                print(
                    "Error: Fabric is working fine, but claude is overloaded. Please try again later.")
            else:
                print(f"Error: {e}")
                print(e)
        if self.args.copy:
            pyperclip.copy(buffer)
        if self.args.output:
            with open(self.args.output, "w") as f:
                f.write(buffer)

    def sendMessage(self, input_data: str, context="", host=''):
        """        Send a message using the input data and generate a response.

        Args:
            input_data (str): The input data to be sent as a message.

        Returns:
            None

        Raises:
            FileNotFoundError: If the specified pattern file is not found.
        """

        wisdomFilePath = os.path.join(
            config_directory, f"patterns/{self.pattern}/system.md"
        )
        user_message = {"role": "user", "content": f"{input_data}"}
        wisdom_File = os.path.join(current_directory, wisdomFilePath)
        system = ""
        if self.pattern:
            try:
                with open(wisdom_File, "r") as f:
                    if context:
                        system = context + '\n\n' + f.read()
                    else:
                        system = f.read()
                    system_message = {"role": "system", "content": system}
                messages = [system_message, user_message]
            except FileNotFoundError:
                print("pattern not found")
                return
        else:
            if context:
                messages = [
                    {'role': 'system', 'content': context}, user_message]
            else:
                messages = [user_message]
        try:
            if self.local:
                if host:
                    asyncio.run(self.localChat(messages, host=host))
                else:
                    asyncio.run(self.localChat(messages))
            elif self.claude:
                asyncio.run(self.claudeChat(system, user_message))
            else:
                response = self.client.chat.completions.create(
                    model=self.model,
                    messages=messages,
                    temperature=0.0,
                    top_p=1,
                    frequency_penalty=0.1,
                    presence_penalty=0.1,
                )
                print(response.choices[0].message.content)
                if self.args.copy:
                    pyperclip.copy(response.choices[0].message.content)
                if self.args.output:
                    with open(self.args.output, "w") as f:
                        f.write(response.choices[0].message.content)
        except Exception as e:
            if "All connection attempts failed" in str(e):
                print(
                    "Error: cannot connect to llama2. If you have not already, please visit https://ollama.com for installation instructions")
            if "CLAUDE_API_KEY" in str(e):
                print(
                    "Error: CLAUDE_API_KEY not found in environment variables. Please run --setup and add the key")
            if "overloaded_error" in str(e):
                print(
                    "Error: Fabric is working fine, but claude is overloaded. Please try again later.")
            if "Attempted to call a sync iterator on an async stream" in str(e):
                print("Error: There is a problem connecting fabric with your local ollama installation. Please visit https://ollama.com for installation instructions. It is possible that you have chosen the wrong model. Please run fabric --listmodels to see the available models and choose the right one with fabric --model <model> or fabric --changeDefaultModel. If this does not work. Restart your computer (always a good idea) and try again. If you are still having problems, please visit https://ollama.com for installation instructions.")
            else:
                print(f"Error: {e}")
                print(e)

    def fetch_available_models(self):
        gptlist = []
        fullOllamaList = []
        claudeList = ['claude-3-opus-20240229',
                      'claude-3-sonnet-20240229',
                      'claude-3-haiku-20240307',
                      'claude-2.1']
        try:
            models = [model.id.strip()
                      for model in self.client.models.list().data]
        except APIConnectionError as e:
            if getattr(e.__cause__, 'args', [''])[0] == "Illegal header value b'Bearer '":
                print("Error: Cannot connect to the OpenAI API Server because the API key is not set. Please run fabric --setup and add a key.")

            else:
                print(
                    f"Error: {e.message} trying to access {e.request.url}: {getattr(e.__cause__, 'args', [''])}")
            sys.exit()
        except Exception as e:
            print(f"Error: {getattr(e.__context__, 'args', [''])[0]}")
            sys.exit()
        if "/" in models[0] or "\\" in models[0]:
            # lmstudio returns full paths to models. Iterate and truncate everything before and including the last slash
            gptlist = [item[item.rfind(
                "/") + 1:] if "/" in item else item[item.rfind("\\") + 1:] for item in models]
        else:
            # Keep items that start with "gpt"
            gptlist = [item.strip()
                       for item in models if item.startswith("gpt")]
        gptlist.sort()
        import ollama
        try:
            default_modelollamaList = ollama.list()['models']
            for model in default_modelollamaList:
                fullOllamaList.append(model['name'])
        except:
            fullOllamaList = []
        return gptlist, fullOllamaList, claudeList

    def get_cli_input(self):
        """ aided by ChatGPT; uses platform library
        accepts either piped input or console input
        from either Windows or Linux

        Args:
            none
        Returns:
            string from either user or pipe
        """
        system = platform.system()
        if system == 'Windows':
            if not sys.stdin.isatty():  # Check if input is being piped
                return sys.stdin.read().strip()  # Read piped input
            else:
                # Prompt user for input from console
                return input("Enter Question: ")
        else:
            return sys.stdin.read()


class Update:
    def __init__(self):
        """Initialize the object with default values."""
        self.repo_zip_url = "https://github.com/danielmiessler/fabric/archive/refs/heads/main.zip"
        self.config_directory = os.path.expanduser("~/.config/fabric")
        self.pattern_directory = os.path.join(
            self.config_directory, "patterns")
        os.makedirs(self.pattern_directory, exist_ok=True)
        print("Updating patterns...")
        self.update_patterns()  # Start the update process immediately

    def update_patterns(self):
        """Update the patterns by downloading the zip from GitHub and extracting it."""
        with tempfile.TemporaryDirectory() as temp_dir:
            zip_path = os.path.join(temp_dir, "repo.zip")
            self.download_zip(self.repo_zip_url, zip_path)
            extracted_folder_path = self.extract_zip(zip_path, temp_dir)
            # The patterns folder will be inside "fabric-main" after extraction
            patterns_source_path = os.path.join(
                extracted_folder_path, "fabric-main", "patterns")
            if os.path.exists(patterns_source_path):
                # If the patterns directory already exists, remove it before copying over the new one
                if os.path.exists(self.pattern_directory):
                    old_pattern_contents = os.listdir(self.pattern_directory)
                    new_pattern_contents = os.listdir(patterns_source_path)
                    custom_patterns = []
                    for pattern in old_pattern_contents:
                        if pattern not in new_pattern_contents:
                            custom_patterns.append(pattern)
                    if custom_patterns:
                        for pattern in custom_patterns:
                            custom_path = os.path.join(
                                self.pattern_directory, pattern)
                            shutil.move(custom_path, patterns_source_path)
                    shutil.rmtree(self.pattern_directory)
                shutil.copytree(patterns_source_path, self.pattern_directory)
                print("Patterns updated successfully.")
            else:
                print("Patterns folder not found in the downloaded zip.")

    def download_zip(self, url, save_path):
        """Download the zip file from the specified URL."""
        response = requests.get(url)
        response.raise_for_status()  # Check if the download was successful
        with open(save_path, 'wb') as f:
            f.write(response.content)
        print("Downloaded zip file successfully.")

    def extract_zip(self, zip_path, extract_to):
        """Extract the zip file to the specified directory."""
        with zipfile.ZipFile(zip_path, 'r') as zip_ref:
            zip_ref.extractall(extract_to)
        print("Extracted zip file successfully.")
        return extract_to  # Return the path to the extracted contents


class Alias:
    def __init__(self):
        self.config_files = []
        self.home_directory = os.path.expanduser("~")
        patternsFolder = os.path.join(
            self.home_directory, ".config/fabric/patterns")
        self.patterns = os.listdir(patternsFolder)

    def execute(self):
        with open(os.path.join(self.home_directory, ".config/fabric/fabric-bootstrap.inc"), "w") as w:
            for pattern in self.patterns:
                w.write(f"alias {pattern}='fabric --pattern {pattern}'\n")


class Setup:
    def __init__(self):
        """        Initialize the object.

        Raises:
            OSError: If there is an error in creating the pattern directory.
        """

        self.config_directory = os.path.expanduser("~/.config/fabric")
        self.pattern_directory = os.path.join(
            self.config_directory, "patterns")
        os.makedirs(self.pattern_directory, exist_ok=True)
        self.shconfigs = []
        home = os.path.expanduser("~")
        if os.path.exists(os.path.join(home, ".bashrc")):
            self.shconfigs.append(os.path.join(home, ".bashrc"))
        if os.path.exists(os.path.join(home, ".bash_profile")):
            self.shconfigs.append(os.path.join(home, ".bash_profile"))
        if os.path.exists(os.path.join(home, ".zshrc")):
            self.shconfigs.append(os.path.join(home, ".zshrc"))
        self.env_file = os.path.join(self.config_directory, ".env")
        self.gptlist = []
        self.fullOllamaList = []
        self.claudeList = ['claude-3-opus-20240229']
        load_dotenv(self.env_file)
        try:
            openaiapikey = os.environ["OPENAI_API_KEY"]
            self.openaiapi_key = openaiapikey
        except:
            pass

    def update_shconfigs(self):
        bootstrap_file = os.path.join(
            self.config_directory, "fabric-bootstrap.inc")
        sourceLine = f'if [ -f "{bootstrap_file}" ]; then . "{bootstrap_file}"; fi'
        for config in self.shconfigs:
            lines = None
            with open(config, 'r') as f:
                lines = f.readlines()
            with open(config, 'w') as f:
                for line in lines:
                    if sourceLine not in line:
                        f.write(line)
                f.write(sourceLine)

    def api_key(self, api_key):
        """        Set the OpenAI API key in the environment file.

        Args:
            api_key (str): The API key to be set.

        Returns:
            None

        Raises:
            OSError: If the environment file does not exist or cannot be accessed.
        """
        api_key = api_key.strip()
        if not os.path.exists(self.env_file) and api_key:
            with open(self.env_file, "w") as f:
                f.write(f"OPENAI_API_KEY={api_key}\n")
            print(f"OpenAI API key set to {api_key}")
        elif api_key:
            # erase the line OPENAI_API_KEY=key and write the new key
            with open(self.env_file, "r") as f:
                lines = f.readlines()
            with open(self.env_file, "w") as f:
                for line in lines:
                    if "OPENAI_API_KEY" not in line:
                        f.write(line)
                f.write(f"OPENAI_API_KEY={api_key}\n")

    def claude_key(self, claude_key):
        """        Set the Claude API key in the environment file.

        Args:
            claude_key (str): The API key to be set.

        Returns:
            None

        Raises:
            OSError: If the environment file does not exist or cannot be accessed.
        """
        claude_key = claude_key.strip()
        if os.path.exists(self.env_file) and claude_key:
            with open(self.env_file, "r") as f:
                lines = f.readlines()
            with open(self.env_file, "w") as f:
                for line in lines:
                    if "CLAUDE_API_KEY" not in line:
                        f.write(line)
                f.write(f"CLAUDE_API_KEY={claude_key}\n")
        elif claude_key:
            with open(self.env_file, "w") as f:
                f.write(f"CLAUDE_API_KEY={claude_key}\n")

    def youtube_key(self, youtube_key):
        """        Set the YouTube API key in the environment file.

        Args:
            youtube_key (str): The API key to be set.

        Returns:
            None

        Raises:
            OSError: If the environment file does not exist or cannot be accessed.
        """
        youtube_key = youtube_key.strip()
        if os.path.exists(self.env_file) and youtube_key:
            with open(self.env_file, "r") as f:
                lines = f.readlines()
            with open(self.env_file, "w") as f:
                for line in lines:
                    if "YOUTUBE_API_KEY" not in line:
                        f.write(line)
                f.write(f"YOUTUBE_API_KEY={youtube_key}\n")
        elif youtube_key:
            with open(self.env_file, "w") as f:
                f.write(f"YOUTUBE_API_KEY={youtube_key}\n")

    def default_model(self, model):
        """Set the default model in the environment file.

        Args:
            model (str): The model to be set.
        """
        model = model.strip()
        env = os.path.expanduser("~/.config/fabric/.env")
        standalone = Standalone(args=[], pattern="")
        gpt, ollama, claude = standalone.fetch_available_models()
        allmodels = gpt + ollama + claude
        if model not in allmodels:
            print(
                f"Error: {model} is not a valid model. Please run fabric --listmodels to see the available models.")
            sys.exit()

        # Only proceed if the model is not empty
        if model:
            if os.path.exists(env):
                # Initialize a flag to track the presence of DEFAULT_MODEL
                there = False
                with open(env, "r") as f:
                    lines = f.readlines()

                # Open the file again to write the changes
                with open(env, "w") as f:
                    for line in lines:
                        # Check each line to see if it contains DEFAULT_MODEL
                        if "DEFAULT_MODEL=" in line:
                            # Update the flag and the line with the new model
                            there = True
                            f.write(f'DEFAULT_MODEL={model}\n')
                        else:
                            # If the line does not contain DEFAULT_MODEL, write it unchanged
                            f.write(line)

                    # If DEFAULT_MODEL was not found in the file, add it
                    if not there:
                        f.write(f'DEFAULT_MODEL={model}\n')

                print(
                    f"Default model changed to {model}. Please restart your terminal to use it.")
            else:
                print("No shell configuration file found.")

    def patterns(self):
        """        Method to update patterns and exit the system.

        Returns:
            None
        """

        Update()

    def run(self):
        """        Execute the Fabric program.

        This method prompts the user for their OpenAI API key, sets the API key in the Fabric object, and then calls the patterns method.

        Returns:
            None
        """

        print("Welcome to Fabric. Let's get started.")
        apikey = input(
            "Please enter your OpenAI API key. If you do not have one or if you have already entered it, press enter.\n")
        self.api_key(apikey)
        print("Please enter your claude API key. If you do not have one, or if you have already entered it, press enter.\n")
        claudekey = input()
        self.claude_key(claudekey)
        print("Please enter your YouTube API key. If you do not have one, or if you have already entered it, press enter.\n")
        youtubekey = input()
        self.youtube_key(youtubekey)
        self.patterns()
        self.update_shconfigs()


class Transcribe:
    def youtube(video_id):
        """ 
        This method gets the transciption
        of a YouTube video designated with the video_id

        Input:
            the video id specifying a YouTube video
            an example url for a video: https://www.youtube.com/watch?v=vF-MQmVxnCs&t=306s
            the video id is vF-MQmVxnCs&t=306s

        Output:
            a transcript for the video

        Raises:
            an exception and prints error


        """
        try:
            transcript_list = YouTubeTranscriptApi.get_transcript(video_id)
            transcript = ""
            for segment in transcript_list:
                transcript += segment['text'] + " "
            return transcript.strip()
        except Exception as e:
            print("Error:", e)
            return None


class AgentSetup:
    def apiKeys(self):
        """Method to set the API keys in the environment file.

        Returns:
            None
        """

        print("Welcome to Fabric. Let's get started.")
        browserless = input("Please enter your Browserless API key\n").strip()
        serper = input("Please enter your Serper API key\n").strip()

        # Entries to be added
        browserless_entry = f"BROWSERLESS_API_KEY={browserless}"
        serper_entry = f"SERPER_API_KEY={serper}"

        # Check and write to the file
        with open(env_file, "r+") as f:
            content = f.read()

            # Determine if the file ends with a newline
            if content.endswith('\n'):
                # If it ends with a newline, we directly write the new entries
                f.write(f"{browserless_entry}\n{serper_entry}\n")
            else:
                # If it does not end with a newline, add one before the new entries
                f.write(f"\n{browserless_entry}\n{serper_entry}\n")
