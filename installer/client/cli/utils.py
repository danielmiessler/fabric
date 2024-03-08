import requests
import os
from openai import OpenAI
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
        env_file = os.path.expanduser(env_file)
        load_dotenv(env_file)
        try:
            apikey = os.environ["OPENAI_API_KEY"]
            self.client = OpenAI()
            self.client.api_key = apikey
        except FileNotFoundError:
            print("No API key found. Use the --apikey option to set the key")
            sys.exit()
        self.local = False
        self.config_pattern_directory = config_directory
        self.pattern = pattern
        self.args = args
        self.model = args.model
        self.claude = False
        sorted_gpt_models, ollamaList, claudeList = self.fetch_available_models()
        self.local = self.model.strip() in ollamaList
        self.claude = self.model.strip() in claudeList

    async def localChat(self, messages):
        from ollama import AsyncClient
        response = await AsyncClient().chat(model=self.model, messages=messages)
        print(response['message']['content'])

    async def localStream(self, messages):
        from ollama import AsyncClient
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

    async def claudeChat(self, system, user):
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

    def streamMessage(self, input_data: str, context=""):
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

    def sendMessage(self, input_data: str, context=""):
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
        if self.args.copy:
            pyperclip.copy(response.choices[0].message.content)
        if self.args.output:
            with open(self.args.output, "w") as f:
                f.write(response.choices[0].message.content)

    def fetch_available_models(self):
        gptlist = []
        fullOllamaList = []
        claudeList = ['claude-3-opus-20240229']
        headers = {
            "Authorization": f"Bearer {self.client.api_key}"
        }

        response = requests.get(
            "https://api.openai.com/v1/models", headers=headers)

        if response.status_code == 200:
            models = response.json().get("data", [])
            # Filter only gpt models
            gpt_models = [model for model in models if model.get(
                "id", "").startswith(("gpt"))]
            # Sort the models alphabetically by their ID
            sorted_gpt_models = sorted(
                gpt_models, key=lambda x: x.get("id"))

            for model in sorted_gpt_models:
                gptlist.append(model.get("id"))
        else:
            print(f"Failed to fetch models: HTTP {response.status_code}")
            sys.exit()
        import ollama
        try:
            default_modelollamaList = ollama.list()['models']
            for model in default_modelollamaList:
                fullOllamaList.append(model['name'].rstrip(":latest"))
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
        home_directory = os.path.expanduser("~")
        self.patterns = os.path.join(home_directory, ".config/fabric/patterns")
        if os.path.exists(os.path.join(home_directory, ".bashrc")):
            self.config_files.append(os.path.join(home_directory, ".bashrc"))
        if os.path.exists(os.path.join(home_directory, ".zshrc")):
            self.config_files.append(os.path.join(home_directory, ".zshrc"))
        if os.path.exists(os.path.join(home_directory, ".bash_profile")):
            self.config_files.append(os.path.join(
                home_directory, ".bash_profile"))
        self.remove_all_patterns()
        self.add_patterns()
        print('Aliases added successfully. Please restart your terminal to use them.')

    def add(self, name, alias):
        for file in self.config_files:
            with open(file, "a") as f:
                f.write(f"alias {name}='{alias}'\n")

    def remove(self, pattern):
        for file in self.config_files:
            # Read the whole file first
            with open(file, "r") as f:
                wholeFile = f.read()

            # Determine if the line to be removed is in the file
            target_line = f"alias {pattern}='fabric --pattern {pattern}'\n"
            if target_line in wholeFile:
                # If the line exists, replace it with nothing (remove it)
                wholeFile = wholeFile.replace(target_line, "")

                # Write the modified content back to the file
                with open(file, "w") as f:
                    f.write(wholeFile)

    def remove_all_patterns(self):
        allPatterns = os.listdir(self.patterns)
        for pattern in allPatterns:
            self.remove(pattern)

    def find_line(self, name):
        for file in self.config_files:
            with open(file, "r") as f:
                lines = f.readlines()
            for line in lines:
                if line.strip("\n") == f"alias ${name}='{alias}'":
                    return line

    def add_patterns(self):
        allPatterns = os.listdir(self.patterns)
        for pattern in allPatterns:
            self.add(pattern, f"fabric --pattern {pattern}")


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
        try:
            self.fetch_available_models()
        except:
            pass

    def fetch_available_models(self):
        headers = {
            "Authorization": f"Bearer {self.openaiapi_key}"
        }

        response = requests.get(
            "https://api.openai.com/v1/models", headers=headers)

        if response.status_code == 200:
            models = response.json().get("data", [])
            # Filter only gpt models
            gpt_models = [model for model in models if model.get(
                "id", "").startswith(("gpt"))]
            # Sort the models alphabetically by their ID
            sorted_gpt_models = sorted(
                gpt_models, key=lambda x: x.get("id"))

            for model in sorted_gpt_models:
                self.gptlist.append(model.get("id"))
        else:
            print(f"Failed to fetch models: HTTP {response.status_code}")
            sys.exit()
        import ollama
        try:
            default_modelollamaList = ollama.list()['models']
            for model in default_modelollamaList:
                self.fullOllamaList.append(model['name'].rstrip(":latest"))
        except:
            self.fullOllamaList = []
        allmodels = self.gptlist + self.fullOllamaList + self.claudeList
        return allmodels

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
                f.write(f"OPENAI_API_KEY={api_key}")
            print(f"OpenAI API key set to {api_key}")
        elif api_key:
            # erase the line OPENAI_API_KEY=key and write the new key
            with open(self.env_file, "r") as f:
                lines = f.readlines()
            with open(self.env_file, "w") as f:
                for line in lines:
                    if "OPENAI_API_KEY" not in line:
                        f.write(line)
                f.write(f"OPENAI_API_KEY={api_key}")

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
                f.write(f"CLAUDE_API_KEY={claude_key}")
        elif claude_key:
            with open(self.env_file, "w") as f:
                f.write(f"CLAUDE_API_KEY={claude_key}")

    def update_fabric_command(self, line, model):
        fabric_command_regex = re.compile(
            r"(alias.*fabric --pattern\s+\S+.*?)( --model.*)?'")
        match = fabric_command_regex.search(line)
        if match:
            base_command = match.group(1)
            # Provide a default value for current_flag
            current_flag = match.group(2) if match.group(2) else ""
            new_flag = ""
            new_flag = f" --model {model}"
            # Update the command if the new flag is different or to remove an existing flag.
            # Ensure to add the closing quote that was part of the original regex
            return f"{base_command}{new_flag}'\n"
        else:
            return line  # Return the line unmodified if no match is found.

    def update_fabric_alias(self, line, model):
        fabric_alias_regex = re.compile(
            r"(alias fabric='[^']+?)( --model.*)?'")
        match = fabric_alias_regex.search(line)
        if match:
            base_command, current_flag = match.groups()
            new_flag = f" --model {model}"
            # Update the alias if the new flag is different or to remove an existing flag.
            return f"{base_command}{new_flag}'\n"
        else:
            return line  # Return the line unmodified if no match is found.

    def clear_alias(self, line):
        fabric_command_regex = re.compile(
            r"(alias fabric='[^']+?)( --model.*)?'")
        match = fabric_command_regex.search(line)
        if match:
            base_command = match.group(1)
            return f"{base_command}'\n"
        else:
            return line  # Return the line unmodified if no match is found.

    def clear_env_line(self, line):
        fabric_command_regex = re.compile(
            r"(alias.*fabric --pattern\s+\S+.*?)( --model.*)?'")
        match = fabric_command_regex.search(line)
        if match:
            base_command = match.group(1)
            return f"{base_command}'\n"
        else:
            return line  # Return the line unmodified if no match is found.

    def pattern(self, line):
        fabric_command_regex = re.compile(
            r"(alias fabric='[^']+?)( --model.*)?'")
        match = fabric_command_regex.search(line)
        if match:
            base_command = match.group(1)
            return f"{base_command}'\n"
        else:
            return line  # Return the line unmodified if no match is found.

    def clean_env(self):
        """Clear the DEFAULT_MODEL from the environment file.

        Returns:
            None
        """
        user_home = os.path.expanduser("~")
        sh_config = None
        # Check for shell configuration files
        if os.path.exists(os.path.join(user_home, ".bashrc")):
            sh_config = os.path.join(user_home, ".bashrc")
        elif os.path.exists(os.path.join(user_home, ".zshrc")):
            sh_config = os.path.join(user_home, ".zshrc")
        else:
            print("No environment file found.")
        if sh_config:
            with open(sh_config, "r") as f:
                lines = f.readlines()
            with open(sh_config, "w") as f:
                for line in lines:
                    modified_line = line
                    # Update existing fabric commands
                    if "fabric --pattern" in line:
                        modified_line = self.clear_env_line(
                            modified_line)
                    elif "fabric=" in line:
                        modified_line = self.clear_alias(
                            modified_line)
                    f.write(modified_line)
            self.remove_duplicates(env_file)
        else:
            print("No shell configuration file found.")

    def default_model(self, model):
        """Set the default model in the environment file.

        Args:
            model (str): The model to be set.
        """
        model = model.strip()
        if model:
            # Write or update the DEFAULT_MODEL in env_file
            allModels = self.claudeList + self.fullOllamaList + self.gptlist
            if model not in allModels:
                print(
                    f"Error: {model} is not a valid model. Please run fabric --listmodels to see the available models.")
                sys.exit()

        # Compile regular expressions outside of the loop for efficiency

        user_home = os.path.expanduser("~")
        sh_config = None
        # Check for shell configuration files
        if os.path.exists(os.path.join(user_home, ".bashrc")):
            sh_config = os.path.join(user_home, ".bashrc")
        elif os.path.exists(os.path.join(user_home, ".zshrc")):
            sh_config = os.path.join(user_home, ".zshrc")

        if sh_config:
            with open(sh_config, "r") as f:
                lines = f.readlines()
            with open(sh_config, "w") as f:
                for line in lines:
                    modified_line = line
                    # Update existing fabric commands
                    if "fabric --pattern" in line:
                        modified_line = self.update_fabric_command(
                            modified_line, model)
                    elif "fabric=" in line:
                        modified_line = self.update_fabric_alias(
                            modified_line, model)
                    f.write(modified_line)
            print(f"""Default model changed to {
                  model}. Please restart your terminal to use it.""")
        else:
            print("No shell configuration file found.")

    def remove_duplicates(self, filename):
        unique_lines = set()
        with open(filename, 'r') as file:
            lines = file.readlines()

        with open(filename, 'w') as file:
            for line in lines:
                if line not in unique_lines:
                    file.write(line)
                    unique_lines.add(line)

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
        self.api_key(apikey.strip())
        print("Please enter your claude API key. If you do not have one, or if you have already entered it, press enter.\n")
        claudekey = input()
        self.claude_key(claudekey.strip())
        self.patterns()


class Transcribe:
    def youtube(video_id):
        """ 
        This method gets the transciption
        of a YouTube video designated with the video_id

        Input:
            the video id specifing a YouTube video
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
        browserless = input("Please enter your Browserless API key\n")
        serper = input("Please enter your Serper API key\n")

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
