import requests
import os
from openai import OpenAI
import pyperclip
import sys
import platform
from dotenv import load_dotenv
from requests.exceptions import HTTPError
from tqdm import tqdm
import zipfile
import tempfile
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
        except KeyError:
            print("OPENAI_API_KEY not found in environment variables.")

        except FileNotFoundError:
            print("No API key found. Use the --apikey option to set the key")
            sys.exit()
        self.config_pattern_directory = config_directory
        self.pattern = pattern
        self.args = args
        self.model = args.model

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
            print(f"Error: {e}")
            print(e)
        if self.args.copy:
            pyperclip.copy(response.choices[0].message.content)
        if self.args.output:
            with open(self.args.output, "w") as f:
                f.write(response.choices[0].message.content)

    def fetch_available_models(self):
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
            sorted_gpt_models = sorted(gpt_models, key=lambda x: x.get("id"))

            for model in sorted_gpt_models:
                print(model.get("id"))
        else:
            print(f"Failed to fetch models: HTTP {response.status_code}")

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

    def api_key(self, api_key):
        """        Set the OpenAI API key in the environment file.

        Args:
            api_key (str): The API key to be set.

        Returns:
            None

        Raises:
            OSError: If the environment file does not exist or cannot be accessed.
        """

        if not os.path.exists(self.env_file):
            with open(self.env_file, "w") as f:
                f.write(f"OPENAI_API_KEY={api_key}")
            print(f"OpenAI API key set to {api_key}")

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
        apikey = input("Please enter your OpenAI API key\n")
        self.api_key(apikey.strip())
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
