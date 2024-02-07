# Imports
import openai
import json
from flask import Flask, request, jsonify
from functools import wraps
import re
import requests

## Define Flask app
app = Flask(__name__)

##################################################
##################################################
#
# ⚠️ CAUTION: This is an HTTP-only server!
#
# If you don't know what you're doing, don't run
#
##################################################
##################################################

## Setup

## Did I mention this is HTTP only? Don't run this on the public internet.

## Set authentication on your APIs
## Let's at least have some kind of auth

# Load your OpenAI API key from a file
with open("openai.key", "r") as key_file:
    openai.api_key = key_file.read().strip()

## Define our own client
client = openai.OpenAI(api_key = openai.api_key)


# Read API tokens from the apikeys.json file
with open("fabric_api_keys.json", "r") as tokens_file:
    valid_tokens = json.load(tokens_file)


# The function to check if the token is valid
def auth_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        # Get the authentication token from request header
        auth_token = request.headers.get("Authorization", "")

        # Remove any bearer token prefix if present
        if auth_token.lower().startswith("bearer "):
            auth_token = auth_token[7:]

        # Get API endpoint from request
        endpoint = request.path

        # Check if token is valid
        user = check_auth_token(auth_token, endpoint)
        if user == "Unauthorized: You are not authorized for this API":
            return jsonify({"error": user}), 401

        return f(*args, **kwargs)

    return decorated_function


# Check for a valid token/user for the given route
def check_auth_token(token, route):
    # Check if token is valid for the given route and return corresponding user
    if route in valid_tokens and token in valid_tokens[route]:
        return valid_tokens[route][token]
    else:
        return "Unauthorized: You are not authorized for this API"


# Define the allowlist of characters
ALLOWLIST_PATTERN = re.compile(r"^[a-zA-Z0-9\s.,;:!?\-]+$")


# Sanitize the content, sort of. Prompt injection is the main threat so this isn't a huge deal
def sanitize_content(content):
    return "".join(char for char in content if ALLOWLIST_PATTERN.match(char))


# Pull the URL content's from the GitHub repo
def fetch_content_from_url(url):
    try:
        response = requests.get(url)
        response.raise_for_status()
        sanitized_content = sanitize_content(response.text)
        return sanitized_content
    except requests.RequestException as e:
        return str(e)


## APIs
# Make path mapping flexible and scalable
pattern_path_mappings = {
    "extwis": {"system_url": "https://raw.githubusercontent.com/danielmiessler/fabric/main/patterns/extract_wisdom/system.md",
               "user_url": "https://raw.githubusercontent.com/danielmiessler/fabric/main/patterns/extract_wisdom/user.md"}
}

# /extwis
@app.route("/<pattern>", methods=["POST"])
@auth_required  # Require authentication
def mill(pattern):
    data = request.get_json()

    # Warn if there's no input
    if "input" not in data:
        return jsonify({"error": "Missing input parameter"}), 400

    # Get data from client
    input_data = data["input"]

    # Set the system and user URLs
    urls = pattern_path_mappings[pattern]
    system_url, user_url = urls["system_url"], urls["user_url"]

    # Fetch the prompt content
    system_content = fetch_content_from_url(system_url)
    user_file_content = fetch_content_from_url(user_url)

    # Build the API call
    system_message = {"role": "system", "content": system_content}
    user_message = {"role": "user", "content": user_file_content + "\n" + input_data}
    messages = [system_message, user_message]
    try:
        response = openai.chat.completions.create(
            model="gpt-4-1106-preview",
            messages=messages,
            temperature=0.0,
            top_p=1,
            frequency_penalty=0.1,
            presence_penalty=0.1,
        )
        assistant_message = response.choices[0].message.content
        return jsonify({"response": assistant_message})
    except Exception as e:
        return jsonify({"error": str(e)}), 500

# Run the application
if __name__ == "__main__":
    app.run(host="127.0.0.1", port=13337, debug=True)
