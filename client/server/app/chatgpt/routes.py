from app.chatgpt import bp, sendMessage
from flask import request, jsonify
import os
from app.chatgpt import streamMessage
from app import sockio
from flask_socketio import emit, disconnect

cwd = os.path.dirname(os.path.abspath(__file__))
config_directory = os.path.expanduser("~/.config/fabric/patterns")
modules = os.listdir(config_directory)


@bp.route("/static/<path_name>", methods=['POST'])
# @jwt_required()
def extractwisdom(path_name):
    wisdomFile = os.path.join(cwd, f"patterns/{path_name}/system.md")
    if os.path.exists(wisdomFile):
        with open(wisdomFile, "r") as f:
            systemPrompt = f.read()
        data = request.get_json()
        input_data = data.get('input')
        return sendMessage(systemPrompt, input_data), 200
    else:
        return {'error': 'module not found'}, 404


@bp.route("/patterns", methods=['GET'])
# @jwt_required()
def get_patterns():
    return jsonify(modules), 200


@sockio.on('connect')
def handle_connect():
    # token = request.headers.get('Authorization')
    # if token:
    #     token = token.split(' ')[1]
    #     if not is_token_valid(token):
    #         emit('error', {'data': 'invalid token'})
    #         disconnect()

    # else:
    #     emit('error', {'data': 'no token'})
    #     disconnect()
    pass


@sockio.on("fabric")
def fabric(message):
    module = message['module']
    input_data = message['input_data']
    if module not in modules:
        available_modules = '\n'.join(modules)
        emit(
            'error', f"module {module} not found. Available modules are:\n{available_modules}")
        disconnect()
    else:
        streamMessage(input_data, f"patterns/{module}/system.md")
    disconnect()
