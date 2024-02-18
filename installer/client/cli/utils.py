import requests
import os
from openai import OpenAI
import pyperclip
import sys
import platform
from dotenv import load_dotenv
from requests.exceptions import HTTPError
from tqdm import tqdm

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
                user_message += {role: "system", content: context}
                messages = [user_message]
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
                user_message += {'role': 'system', 'content': context}
                messages = [user_message]
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
        """        Initialize the object with default values and update patterns.

        This method initializes the object with default values for root_api_url, config_directory, and pattern_directory.
        It then creates the pattern_directory if it does not exist and calls the update_patterns method to update the patterns.

        Raises:
            OSError: If there is an issue creating the pattern_directory.
        """

        self.root_api_url = "https://api.github.com/repos/danielmiessler/fabric/contents/patterns?ref=main"
        self.config_directory = os.path.expanduser("~/.config/fabric")
        self.pattern_directory = os.path.join(
            self.config_directory, "patterns")
        os.makedirs(self.pattern_directory, exist_ok=True)
        self.update_patterns()  # Call the update process from a method.

    def update_patterns(self):
        """        Update the patterns by downloading from the GitHub directory.

        Raises:
            HTTPError: If there is an HTTP error while downloading patterns.
        """

        try:
            self.progress_bar = tqdm(desc="Downloading Patterns…", unit="file")
            self.get_github_directory_contents(
                self.root_api_url, self.pattern_directory
            )
            # Close progress bar on success before printing the message.
            self.progress_bar.close()
        except HTTPError as e:
            # Ensure progress bar is closed on HTTPError as well.
            self.progress_bar.close()
            if e.response.status_code == 403:
                print(
                    "GitHub API rate limit exceeded. Please wait before trying again."
                )
                sys.exit()
            else:
                print(f"Failed to download patterns due to an HTTP error: {e}")
            sys.exit()  # Exit after handling the error.

    def download_file(self, url, local_path):
        """        Download a file from the given URL and save it to the local path.

        Args:
            url (str): The URL of the file to be downloaded.
            local_path (str): The local path where the file will be saved.

        Raises:
            HTTPError: If an HTTP error occurs during the download process.
        """

        try:
            response = requests.get(url)
            response.raise_for_status()
            with open(local_path, "wb") as f:
                f.write(response.content)
            self.progress_bar.update(1)
        except HTTPError as e:
            print(f"Failed to download file {url}. HTTP error: {e}")
            sys.exit()

    def process_item(self, item, local_dir):
        """        Process the given item and save it to the local directory.

        Args:
            item (dict): The item to be processed, containing information about the type, download URL, name, and URL.
            local_dir (str): The local directory where the item will be saved.

        Returns:
            None

        Raises:
            OSError: If there is an issue creating the new directory using os.makedirs.
        """

        if item["type"] == "file":
            self.download_file(
                item["download_url"], os.path.join(local_dir, item["name"])
            )
        elif item["type"] == "dir":
            new_dir = os.path.join(local_dir, item["name"])
            os.makedirs(new_dir, exist_ok=True)
            self.get_github_directory_contents(item["url"], new_dir)

    def get_github_directory_contents(self, api_url, local_dir):
        """        Get the contents of a directory from GitHub API and process each item.

        Args:
            api_url (str): The URL of the GitHub API endpoint for the directory.
            local_dir (str): The local directory where the contents will be processed.

        Returns:
            None

        Raises:
            HTTPError: If an HTTP error occurs while fetching the directory contents.
                If the status code is 403, it prints a message about GitHub API rate limit exceeded
                and closes the progress bar. For any other status code, it prints a message
                about failing to fetch directory contents due to an HTTP error.
        """

        try:
            response = requests.get(api_url)
            response.raise_for_status()
            jsonList = response.json()
            for item in jsonList:
                self.process_item(item, local_dir)
        except HTTPError as e:
            if e.response.status_code == 403:
                print(
                    "GitHub API rate limit exceeded. Please wait before trying again."
                )
                self.progress_bar.close()  # Ensure the progress bar is cleaned up properly
            else:
                print(
                    f"Failed to fetch directory contents due to an HTTP error: {e}")


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
        sys.exit()

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
