import requests
import yaml
import os
from openai import OpenAI
from dotenv import load_dotenv
import pyperclip
import socketio
import sys
import subprocess
import shlex

current_directory = os.path.dirname(os.path.realpath(__file__))
config_file = os.path.join(current_directory, "config.yaml")
config_directory = os.path.expanduser("~/.config/fabric")
env_file = os.path.join(config_directory, '.env')
try:
    config = yaml.safe_load(open(config_file))['server']
except FileNotFoundError:
    config = {}
gunicorn_directory = os.path.join(current_directory, ".venv/bin")


class Utilities:
    def __init__(self):
        domain = config['domain']
        port = config['port']
        baseurl = ''
        if config['port'] != 443:
            baseurl = f'http://{domain}:{port}'
        else:
            baseurl = f'https://{domain}'
        self.summarizestream = f"ws://{domain}:{port}"


class Standalone:
    def __init__(self, args, pattern=''):
        with open(env_file, "r") as f:
            apikey = f.read().split("=")[1]
            self.client = OpenAI(api_key=apikey)
        self.pattern = pattern
        self.args = args

    def streamMessage(self, input_data: str):
        wisdomfileDirectory = os.path.join(
            current_directory, "server/app/chatgpt/patterns")
        wisdomFilePath = os.path.join(
            wisdomfileDirectory, f"{self.pattern}/system.md")
        user_message = {"role": "user", "content": f"{input_data}"}
        wisdom_File = os.path.join(
            current_directory, wisdomFilePath)
        buffer = ''
        if self.pattern:
            try:
                with open(wisdom_File, "r") as f:
                    system = f.read()
                    system_message = {"role": "system", "content": system}
                messages = [system_message, user_message]
            except FileNotFoundError:
                print('pattern not found')
                return
        else:
            messages = [user_message]
        try:
            stream = self.client.chat.completions.create(
                model="gpt-4-1106-preview",
                messages=messages,
                temperature=0.0,
                top_p=1,
                frequency_penalty=0.1,
                presence_penalty=0.1,
                stream=True
            )
            for chunk in stream:
                if chunk.choices[0].delta.content is not None:
                    char = chunk.choices[0].delta.content
                    buffer += char
                    if char not in ['\n', ' ']:
                        print(char, end='')
                    elif char == ' ':
                        print(' ', end='')  # Explicitly handle spaces
                    elif char == '\n':
                        print()  # Handle newlines
                sys.stdout.flush()
        except Exception as e:
            print(f"Error: {e}")
        if self.args.copy:
            pyperclip.copy(buffer)
        if self.args.output:
            with open(self.args.output, 'w') as f:
                f.write(buffer)

    def sendMessage(self, input_data: str):
        wisdomfileDirectory = os.path.join(
            current_directory, "server/app/chatgpt/patterns")
        wisdomFilePath = os.path.join(
            wisdomfileDirectory, f"{self.pattern}/system.md")
        user_message = {"role": "user", "content": f"{input_data}"}
        wisdom_File = os.path.join(
            current_directory, wisdomFilePath)
        if self.pattern:
            try:
                with open(wisdom_File, "r") as f:
                    system = f.read()
                    system_message = {"role": "system", "content": system}
                messages = [system_message, user_message]
            except FileNotFoundError:
                print('pattern not found')
                return
        else:
            messages = [user_message]
        try:
            response = self.client.chat.completions.create(
                model="gpt-4-1106-preview",
                messages=messages,
                temperature=0.0,
                top_p=1,
                frequency_penalty=0.1,
                presence_penalty=0.1
            )
            print(response.choices[0].message.content)
        except Exception as e:
            print(f"Error: {e}")
        if self.args.copy:
            pyperclip.copy(response.choices[0].message.content)
        if self.args.output:
            with open(self.args.output, 'w') as f:
                f.write(response.choices[0].message.content)


class Remote:
    def __init__(self, module, args):
        self.module = module
        self.buffer = ''
        self.utils = Utilities()
        self.sio = socketio.Client()
        self.setup_handlers()
        self.args = args

    def setup_handlers(self):
        @self.sio.event
        def message(data):
            global buffer
            self.buffer += data
            if self.args.stream:
                for char in data:
                    if char not in ['\n', ' ']:  # If the character is not a newline or space
                        print(char, end='')
                    elif char == ' ':
                        print(' ', end='')  # Explicitly handle spaces
                    elif char == '\n':
                        print()  # Handle newlines
                sys.stdout.flush()

        @self.sio.event
        def error(data):
            print(data)

        @self.sio.event
        def disconnect():
            if self.args.copy:
                pyperclip.copy(self.buffer)
            if self.args.output:
                with open(self.args.output, 'w') as f:
                    f.write(self.buffer)
            if not self.args.stream:
                print(self.buffer, end='')

    def analyze(self, text, copy_to_clipboard=False, save_to_file=False):
        url = self.utils.summarizestream
        self.sio.connect(url)
        self.sio.emit('fabric', {'input_data': text, 'module': self.module})
        self.sio.wait()
        if copy_to_clipboard:
            pyperclip.copy(self.buffer)
        if save_to_file:
            with open(save_to_file, 'w') as f:
                f.write(self.buffer)

    def disconnect_handler(self):
        self.sio.disconnect()


class Server:
    def __init__(self):
        server_directory = os.path.join(
            current_directory, "server")
        os.chdir(server_directory)

    def run_server(self, domain: str, port: str):
        print(f"please visit http://{domain}:{port} to view the web frontend")
        command = f"{gunicorn_directory}/gunicorn -k gevent -w 1 --timeout 240 -b {domain}:{port} 'app:create_app()'"
        subprocess.run(shlex.split(command))
