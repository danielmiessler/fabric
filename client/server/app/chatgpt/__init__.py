from openai import OpenAI, chat
from flask import jsonify, Blueprint
import os
from flask_socketio import send, emit

current_directory = os.path.dirname(os.path.abspath(__file__))
config_directory = os.path.expanduser("~/.config/fabric")
api_key_file = os.path.join(config_directory, ".env")
with open(api_key_file, "r") as f:
    apiKey = f.read().split("=")[1]
    client = OpenAI(api_key=apiKey)
bp = Blueprint('chatgpt', __name__)


def sendMessage(system: str, input_data: str, user=''):
    system_message = {"role": "system", "content": system}
    user_message = {"role": "user", "content": f"{user}\n{input_data}"}
    messages = [system_message, user_message]
    try:
        response = client.chat.completions.create(
            model="gpt-4-1106-preview",
            messages=messages,
            temperature=0.0,
            top_p=1,
            frequency_penalty=0.1,
            presence_penalty=0.1
        )
        assistant_message = response.choices[0].message.content
        return jsonify({"response": assistant_message})
    except Exception as e:
        return jsonify({"error": str(e)})


def streamMessage(input_data: str, wisdomFile: str):
    # Similar logic as sendMessage but adapted for streaming
    user_message = {"role": "user", "content": f"{input_data}"}
    wisdom_File = os.path.join(
        current_directory, wisdomFile)
    with open(wisdom_File, "r") as f:
        system = f.read()
        system_message = {"role": "system", "content": system}
    messages = [system_message, user_message]
    try:
        # Note: You need to modify the API call to support streaming
        stream = client.chat.completions.create(
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
                send(chunk.choices[0].delta.content, end="")
    except Exception as e:
        emit('error', {'data': str(e)})


if 1 == 1:
    from app.chatgpt import routes
