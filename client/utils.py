import requests
import os
from openai import OpenAI
import pyperclip
import sys

current_directory = os.path.dirname(os.path.realpath(__file__))
config_directory = os.path.expanduser("~/.config/fabric")
env_file = os.path.join(config_directory, ".env")


class Standalone:
    def __init__(self, args, pattern=""):
        try:
            with open(env_file, "r") as f:
                apikey = f.read().split("=")[1]
                self.client = OpenAI(api_key=apikey)
        except FileNotFoundError:
            print("No API key found. Use the --apikey option to set the key")
            sys.exit()
        self.config_pattern_directory = config_directory
        self.pattern = pattern
        self.args = args

    def streamMessage(self, input_data: str):
        wisdomFilePath = os.path.join(
            config_directory, f"patterns/{self.pattern}/system.md"
        )
        user_message = {"role": "user", "content": f"{input_data}"}
        wisdom_File = os.path.join(current_directory, wisdomFilePath)
        buffer = ""
        if self.pattern:
            try:
                with open(wisdom_File, "r") as f:
                    system = f.read()
                    system_message = {"role": "system", "content": system}
                messages = [system_message, user_message]
            except FileNotFoundError:
                print("pattern not found")
                return
        else:
            messages = [user_message]
        try:
            stream = self.client.chat.completions.create(
                model="gpt-4-turbo-preview",
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
        if self.args.copy:
            pyperclip.copy(buffer)
        if self.args.output:
            with open(self.args.output, "w") as f:
                f.write(buffer)

    def sendMessage(self, input_data: str):
        wisdomFilePath = os.path.join(
            config_directory, f"patterns/{self.pattern}/system.md"
        )
        user_message = {"role": "user", "content": f"{input_data}"}
        wisdom_File = os.path.join(current_directory, wisdomFilePath)
        if self.pattern:
            try:
                with open(wisdom_File, "r") as f:
                    system = f.read()
                    system_message = {"role": "system", "content": system}
                messages = [system_message, user_message]
            except FileNotFoundError:
                print("pattern not found")
                return
        else:
            messages = [user_message]
        try:
            response = self.client.chat.completions.create(
                model="gpt-4-turbo-preview",
                messages=messages,
                temperature=0.0,
                top_p=1,
                frequency_penalty=0.1,
                presence_penalty=0.1,
            )
            print(response.choices[0].message.content)
        except Exception as e:
            print(f"Error: {e}")
        if self.args.copy:
            pyperclip.copy(response.choices[0].message.content)
        if self.args.output:
            with open(self.args.output, "w") as f:
                f.write(response.choices[0].message.content)


class Update:
    def __init__(self):
        # Initialize with the root API URL
        self.root_api_url = "https://api.github.com/repos/danielmiessler/fabric/contents/patterns?ref=main"
        self.config_directory = os.path.expanduser("~/.config/fabric")
        self.pattern_directory = os.path.join(self.config_directory, "patterns")
        # Ensure local directory exists
        os.makedirs(self.pattern_directory, exist_ok=True)
        self.get_github_directory_contents(self.root_api_url, self.pattern_directory)

    def download_file(self, url, local_path):
        """
        Download a file from a URL to a local path.
        """
        response = requests.get(url)
        response.raise_for_status()  # This will raise an exception for HTTP error codes
        with open(local_path, "wb") as f:
            f.write(response.content)

    def process_item(self, item, local_dir):
        """
        Process an individual item, downloading if it's a file, or processing further if it's a directory.
        """
        if item["type"] == "file":
            print(f"Downloading file: {item['name']} to {local_dir}")
            self.download_file(
                item["download_url"], os.path.join(local_dir, item["name"])
            )
        elif item["type"] == "dir":
            new_dir = os.path.join(local_dir, item["name"])
            os.makedirs(new_dir, exist_ok=True)
            self.get_github_directory_contents(item["url"], new_dir)

    def get_github_directory_contents(self, api_url, local_dir):
        """
        Fetches the contents of a directory in a GitHub repository and downloads files, recursively handling directories.
        """
        response = requests.get(api_url)
        response.raise_for_status()  # This will raise an exception for HTTP error codes
        jsonList = response.json()
        for item in jsonList:
            self.process_item(item, local_dir)
