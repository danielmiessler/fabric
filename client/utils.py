import requests
import os
from openai import OpenAI
import pyperclip
import sys
from dotenv import load_dotenv
from requests.exceptions import HTTPError
from tqdm import tqdm

current_directory = os.path.dirname(os.path.realpath(__file__))
config_directory = os.path.expanduser("~/.config/fabric")
env_file = os.path.join(config_directory, ".env")


class Standalone:
    def __init__(self, args, pattern="", env_file="~/.config/fabric/.env"):
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
            print(e)
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
            print(e)
        if self.args.copy:
            pyperclip.copy(response.choices[0].message.content)
        if self.args.output:
            with open(self.args.output, "w") as f:
                f.write(response.choices[0].message.content)


class Update:
    def __init__(self):
        self.root_api_url = "https://api.github.com/repos/danielmiessler/fabric/contents/patterns?ref=main"
        self.config_directory = os.path.expanduser("~/.config/fabric")
        self.pattern_directory = os.path.join(self.config_directory, "patterns")
        os.makedirs(self.pattern_directory, exist_ok=True)
        self.update_patterns()  # Call the update process from a method.

    def update_patterns(self):
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
        if item["type"] == "file":
            self.download_file(
                item["download_url"], os.path.join(local_dir, item["name"])
            )
        elif item["type"] == "dir":
            new_dir = os.path.join(local_dir, item["name"])
            os.makedirs(new_dir, exist_ok=True)
            self.get_github_directory_contents(item["url"], new_dir)

    def get_github_directory_contents(self, api_url, local_dir):
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
                print(f"Failed to fetch directory contents due to an HTTP error: {e}")


class Setup:
    def __init__(self):
        self.config_directory = os.path.expanduser("~/.config/fabric")
        self.pattern_directory = os.path.join(self.config_directory, "patterns")
        os.makedirs(self.pattern_directory, exist_ok=True)
        self.env_file = os.path.join(self.config_directory, ".env")

    def api_key(self, api_key):
        if not os.path.exists(self.env_file):
            with open(self.env_file, "w") as f:
                f.write(f"OPENAI_API_KEY={api_key}")
            print(f"OpenAI API key set to {api_key}")

    def patterns(self):
        Update()
        sys.exit()

    def run(self):
        print("Welcome to Fabric. Let's get started.")
        apikey = input("Please enter your OpenAI API key\n")
        self.api_key(apikey.strip())
        self.patterns()
