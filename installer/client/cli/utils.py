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
import subprocess
import shutil
from youtube_transcript_api import YouTubeTranscriptApi

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
        self.client = None
        load_dotenv(env_file)
        if "OPENAI_API_KEY" in os.environ:
            api_key = os.environ['OPENAI_API_KEY']
            self.client = OpenAI(api_key=api_key)
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
        sorted_gpt_models, ollamaList, claudeList, googleList = self.fetch_available_models()
        self.sorted_gpt_models = sorted_gpt_models
        self.ollamaList = ollamaList
        self.claudeList = claudeList
        self.googleList = googleList
        self.local = self.model in ollamaList
        self.claude = self.model in claudeList
        self.google = self.model in googleList

    async def localChat(self, messages, host=''):
        from ollama import AsyncClient
        response = None
        if host:
            response = await AsyncClient(host=host).chat(model=self.model, messages=messages)
        else:
            response = await AsyncClient().chat(model=self.model, messages=messages)
        print(response['message']['content'])
        copy = self.args.copy
        if copy:
            pyperclip.copy(response['message']['content'])
        if self.args.output:
            with open(self.args.output, "w") as f:
                f.write(response['message']['content'])

    async def localStream(self, messages, host=''):
        from ollama import AsyncClient
        buffer = ""
        if host:
            async for part in await AsyncClient(host=host).chat(model=self.model, messages=messages, stream=True):
                buffer += part['message']['content']
                print(part['message']['content'], end='', flush=True)
        else:
            async for part in await AsyncClient().chat(model=self.model, messages=messages, stream=True):
                buffer += part['message']['content']
                print(part['message']['content'], end='', flush=True)
        if self.args.output:
            with open(self.args.output, "w") as f:
                f.write(buffer)
        if self.args.copy:
            pyperclip.copy(buffer)

    async def claudeStream(self, system, user):
        from anthropic import AsyncAnthropic
        self.claudeApiKey = os.environ["CLAUDE_API_KEY"]
        Streamingclient = AsyncAnthropic(api_key=self.claudeApiKey)
        buffer = ""
        async with Streamingclient.messages.stream(
            max_tokens=4096,
            system=system,
            messages=[user],
            model=self.model, temperature=self.args.temp, top_p=self.args.top_p
        ) as stream:
            async for text in stream.text_stream:
                buffer += text
                print(text, end="", flush=True)
            print()
        if self.args.copy:
            pyperclip.copy(buffer)
        if self.args.output:
            with open(self.args.output, "w") as f:
                f.write(buffer)
        if self.args.session:
            from .helper import Session
            session = Session()
            session.save_to_session(
                system, user, buffer, self.args.session)
        message = await stream.get_final_message()

    async def claudeChat(self, system, user, copy=False):
        from anthropic import Anthropic
        self.claudeApiKey = os.environ["CLAUDE_API_KEY"]
        client = Anthropic(api_key=self.claudeApiKey)
        message = None
        message = client.messages.create(
            max_tokens=4096,
            system=system,
            messages=[user],
            model=self.model,
            temperature=self.args.temp, top_p=self.args.top_p
        )
        print(message.content[0].text)
        copy = self.args.copy
        if copy:
            pyperclip.copy(message.content[0].text)
        if self.args.output:
            with open(self.args.output, "w") as f:
                f.write(message.content[0].text)
        if self.args.session:
            from .helper import Session
            session = Session()
            session.save_to_session(
                system, user, message.content[0].text, self.args.session)

    async def googleChat(self, system, user, copy=False):
        import google.generativeai as genai
        self.googleApiKey = os.environ["GOOGLE_API_KEY"]
        genai.configure(api_key=self.googleApiKey)
        model = genai.GenerativeModel(
            model_name=self.model, system_instruction=system)
        response = model.generate_content(user)
        print(response.text)
        if copy:
            pyperclip.copy(response.text)
        if self.args.output:
            with open(self.args.output, "w") as f:
                f.write(response.text)
        if self.args.session:
            from .helper import Session
            session = Session()
            session.save_to_session(
                system, user, response.text, self.args.session)

    async def googleStream(self, system, user, copy=False):
        import google.generativeai as genai
        buffer = ""
        self.googleApiKey = os.environ["GOOGLE_API_KEY"]
        genai.configure(api_key=self.googleApiKey)
        model = genai.GenerativeModel(
            model_name=self.model, system_instruction=system)
        response = model.generate_content(user, stream=True)
        for chunk in response:
            buffer += chunk.text
            print(chunk.text)
        if copy:
            pyperclip.copy(buffer)
        if self.args.output:
            with open(self.args.output, "w") as f:
                f.write(buffer)
        if self.args.session:
            from .helper import Session
            session = Session()
            session.save_to_session(
                system, user, buffer, self.args.session)

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
        session_message = ""
        user = ""
        if self.args.session:
            from .helper import Session
            session = Session()
            session_message = session.read_from_session(
                self.args.session)
        if session_message:
            user = session_message + '\n' + input_data
        else:
            user = input_data
        user_message = {"role": "user", "content": f"{input_data}"}
        wisdom_File = wisdomFilePath
        buffer = ""
        system = ""
        if self.pattern:
            try:
                with open(wisdom_File, "r") as f:
                    if context:
                        system = context + '\n\n' + f.read()
                        if session_message:
                            system = session_message + '\n' + system
                    else:
                        system = f.read()
                        if session_message:
                            system = session_message + '\n' + system
                    system_message = {"role": "system", "content": system}
                messages = [system_message, user_message]
            except FileNotFoundError:
                print("pattern not found")
                return
        else:
            if session_message:
                user_message['content'] = session_message + \
                    '\n' + user_message['content']
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
            elif self.google:
                if system == "":
                    system = " "
                asyncio.run(self.googleStream(system, user_message['content']))
            else:
                stream = self.client.chat.completions.create(
                    model=self.model,
                    messages=messages,
                    temperature=self.args.temp,
                    top_p=self.args.top_p,
                    frequency_penalty=self.args.frequency_penalty,
                    presence_penalty=self.args.presence_penalty,
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
        if self.args.session:
            from .helper import Session
            session = Session()
            session.save_to_session(
                system, user, buffer, self.args.session)

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
        user = input_data
        user_message = {"role": "user", "content": f"{input_data}"}
        wisdom_File = os.path.join(current_directory, wisdomFilePath)
        system = ""
        session_message = ""
        if self.args.session:
            from .helper import Session
            session = Session()
            session_message = session.read_from_session(
                self.args.session)
        if self.pattern:
            try:
                with open(wisdom_File, "r") as f:
                    if context:
                        if session_message:
                            system = session_message + '\n' + context + '\n\n' + f.read()
                        else:
                            system = context + '\n\n' + f.read()
                    else:
                        if session_message:
                            system = session_message + '\n' + f.read()
                        else:
                            system = f.read()
                    system_message = {"role": "system", "content": system}
                messages = [system_message, user_message]
            except FileNotFoundError:
                print("pattern not found")
                return
        else:
            if session_message:
                user_message['content'] = session_message + \
                    '\n' + user_message['content']
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
            elif self.google:
                if system == "":
                    system = " "
                asyncio.run(self.googleChat(system, user_message['content']))
            else:
                response = self.client.chat.completions.create(
                    model=self.model,
                    messages=messages,
                    temperature=self.args.temp,
                    top_p=self.args.top_p,
                    frequency_penalty=self.args.frequency_penalty,
                    presence_penalty=self.args.presence_penalty,
                )
                print(response.choices[0].message.content)
                if self.args.copy:
                    pyperclip.copy(response.choices[0].message.content)
                if self.args.output:
                    with open(self.args.output, "w") as f:
                        f.write(response.choices[0].message.content)
                if self.args.session:
                    from .helper import Session
                    session = Session()
                    session.save_to_session(
                        system, user, response.choices[0], self.args.session)
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
        googleList = []
        if "CLAUDE_API_KEY" in os.environ:
            claudeList = ['claude-3-opus-20240229', 'claude-3-sonnet-20240229',
                          'claude-3-haiku-20240307', 'claude-2.1']
        else:
            claudeList = []

        try:
            if self.client:
                models = [model.id.strip()
                          for model in self.client.models.list().data]
                if "/" in models[0] or "\\" in models[0]:
                    gptlist = [item[item.rfind(
                        "/") + 1:] if "/" in item else item[item.rfind("\\") + 1:] for item in models]
                else:
                    gptlist = [item.strip()
                               for item in models if item.startswith("gpt")]
                gptlist.sort()
        except APIConnectionError as e:
            pass
        except Exception as e:
            print(f"Error: {getattr(e.__context__, 'args', [''])[0]}")
            sys.exit()

        import ollama
        try:
            remoteOllamaServer = getattr(self.args, 'remoteOllamaServer', None)
            if remoteOllamaServer:
                client = ollama.Client(host=self.args.remoteOllamaServer)
                default_modelollamaList = client.list()['models']
            else:
                default_modelollamaList = ollama.list()['models']
            for model in default_modelollamaList:
                fullOllamaList.append(model['name'])
        except:
            fullOllamaList = []
        try:
            import google.generativeai as genai
            genai.configure(api_key=os.environ["GOOGLE_API_KEY"])
            for m in genai.list_models():
                if 'generateContent' in m.supported_generation_methods:
                    googleList.append(m.name)
        except:
            googleList = []

        return gptlist, fullOllamaList, claudeList, googleList

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

    def agents(self, userInput):
        from praisonai import PraisonAI
        model = self.model
        os.environ["OPENAI_MODEL_NAME"] = model
        if model in self.sorted_gpt_models:
            os.environ["OPENAI_API_BASE"] = "https://api.openai.com/v1/"
        elif model in self.ollamaList:
            os.environ["OPENAI_API_BASE"] = "http://localhost:11434/v1"
            os.environ["OPENAI_API_KEY"] = "NA"

        elif model in self.claudeList:
            print("Claude is not supported in this mode")
            sys.exit()
        print("Starting PraisonAI...")
        praison_ai = PraisonAI(auto=userInput, framework="autogen")
        praison_ai.main()


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
        self.googleList = []
        self.claudeList = ['claude-3-opus-20240229']
        load_dotenv(self.env_file)
        try:
            openaiapikey = os.environ["OPENAI_API_KEY"]
            self.openaiapi_key = openaiapikey
        except:
            pass

    def __ensure_env_file_created(self):
        """        Ensure that the environment file is created.

        Returns:
            None

        Raises:
            OSError: If the environment file cannot be created.
        """
        print("Creating empty environment file...")
        if not os.path.exists(self.env_file):
            with open(self.env_file, "w") as f:
                f.write("#No API key set\n")
        print("Environment file created.")

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

    def google_key(self, google_key):
        """        Set the Google API key in the environment file.

        Args:
            google_key (str): The API key to be set.

        Returns:
            None

        Raises:
            OSError: If the environment file does not exist or cannot be accessed.
        """
        google_key = google_key.strip()
        if os.path.exists(self.env_file) and google_key:
            with open(self.env_file, "r") as f:
                lines = f.readlines()
            with open(self.env_file, "w") as f:
                for line in lines:
                    if "GOOGLE_API_KEY" not in line:
                        f.write(line)
                f.write(f"GOOGLE_API_KEY={google_key}\n")
        elif google_key:
            with open(self.env_file, "w") as f:
                f.write(f"GOOGLE_API_KEY={google_key}\n")

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
        gpt, ollama, claude, google = standalone.fetch_available_models()
        allmodels = gpt + ollama + claude + google
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
        print("Please enter your Google API key. If you do not have one, or if you have already entered it, press enter.\n")
        googlekey = input()
        self.google_key(googlekey)
        print("Please enter your YouTube API key. If you do not have one, or if you have already entered it, press enter.\n")
        youtubekey = input()
        self.youtube_key(youtubekey)
        self.patterns()
        self.update_shconfigs()
        self.__ensure_env_file_created()


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


def run_electron_app():
    # Step 1: Set CWD to the directory of the script
    os.chdir(os.path.dirname(os.path.realpath(__file__)))

    # Step 2: Check for the './installer/client/gui' directory
    target_dir = '../gui'
    if not os.path.exists(target_dir):
        print(f"""The directory {
              target_dir} does not exist. Please check the path and try again.""")
        return

    # Step 3: Check for NPM installation
    try:
        subprocess.run(['npm', '--version'], check=True,
                       stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    except subprocess.CalledProcessError:
        print("NPM is not installed. Please install NPM and try again.")
        return

    # If this point is reached, NPM is installed.
    # Step 4: Change directory to the Electron app's directory
    os.chdir(target_dir)

    # Step 5: Run 'npm install' and 'npm start'
    try:
        print("Running 'npm install'... This might take a few minutes.")
        subprocess.run(['npm', 'install'], check=True)
        print(
            "'npm install' completed successfully. Starting the Electron app with 'npm start'...")
        subprocess.run(['npm', 'start'], check=True)
    except subprocess.CalledProcessError as e:
        print(f"An error occurred while executing NPM commands: {e}")
