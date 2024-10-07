import os
import subprocess
import logging
from flask import Flask, request, jsonify
from flask_cors import CORS

app = Flask(__name__)
CORS(app, resources={r"/*": {"origins": "http://192.168.4.21:3000"}})

logging.basicConfig(level=logging.INFO)

def get_patterns():
    patterns_dir = 'patterns'
    patterns = []
    for item in os.listdir(patterns_dir):
        if os.path.isdir(os.path.join(patterns_dir, item)):
            patterns.append(item)
    return patterns

@app.route('/patterns', methods=['GET'])
def list_patterns():
    patterns = get_patterns()
    return jsonify(patterns)

@app.route('/run-pattern', methods=['POST'])
def run_pattern():
    data = request.json
    pattern = data.get('pattern')
    params = data.get('params', [])

    fabric_command = ['fabric']
    
    if pattern:
        fabric_command.extend(['--pattern', pattern])

    youtube_url = next((param for param in params if param.startswith('-y')), None)
    if youtube_url:
        youtube_index = params.index(youtube_url)
        fabric_command.extend(['-y', params[youtube_index + 1]])
        params = params[:youtube_index] + params[youtube_index + 2:]

    fabric_command.extend(params)

    logging.info(f"Executing command: {' '.join(fabric_command)}")

    try:
        result = subprocess.run(fabric_command, capture_output=True, text=True)
        if result.returncode != 0:
            logging.error(f"Command failed with error: {result.stderr}")
            return jsonify({'error': result.stderr}), 500
        return jsonify({'output': result.stdout})
    except Exception as e:
        logging.exception("An error occurred while running the fabric command")
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
