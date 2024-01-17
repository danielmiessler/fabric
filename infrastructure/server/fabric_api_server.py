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


# Set your OpenAI API key
with open("openai.key", "r") as key_file:
    openai.api_key = key_file.read().strip()

## APIs


# /extwis
@app.route("/extwis", methods=["POST"])
@auth_required  # Require authentication
def extwis():
    data = request.get_json()

    # Warn if there's no input
    if "input" not in data:
        return jsonify({"error": "Missing input parameter"}), 400

    # Get data from client
    input_data = data["input"]

    # Set the system and user URLs
    system_url = "https://raw.githubusercontent.com/danielmiessler/fabric/main/patterns/extract_wisdom/system.md"
    user_url = "https://raw.githubusercontent.com/danielmiessler/fabric/main/patterns/extract_wisdom/user.md"

    # Fetch the prompt content
    system_content = fetch_content_from_url(system_url)
    user_file_content = fetch_content_from_url(user_url)

    # Build the API call
    system_message = {"role": "system", "content": system_content}
    user_message = {"role": "user", "content": user_file_content + "\n" + input_data}
    messages = [system_message, user_message]
    try:
        response = openai.ChatCompletion.create(
            model="gpt-4-1106-preview",
            messages=messages,
            temperature=0.0,
            top_p=1,
            frequency_penalty=0.1,
            presence_penalty=0.1,
        )
        assistant_message = response["choices"][0]["message"]["content"]
        return jsonify({"response": assistant_message})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


# /labelandrate
@app.route("/labelandrate", methods=["POST"])
def labelandrate():
    data = request.get_json()

    if "input" not in data:
        return jsonify({"error": "Missing input parameter"}), 400

    input_data = data["input"]
    system_message = {
        "role": "system",
        "content": """

You are an ultra-wise and brilliant classifier and judge of content. You label content with a a comma-separated list of single-word labels and then give it a quality rating.

Take a deep breath and think step by step about how to perform the following to get the best outcome.

STEPS:

1. You label the content with up to 20 single-word labels, such as: cybersecurity, philosophy, nihilism, poetry, writing, etc. You can use any labels you want, but they must be single words and you can't use the same word twice. This goes in a section called LABELS:.

2. You then rate the content based on the number of ideas in the input (below ten is bad, between 11 and 20 is good, and above 25 is excellent) combined with how well it matches the THEMES of: human meaning, the future of AI, mental models, abstract thinking, unconvential thinking, meaning in a post-ai world, continuous improvement, reading, art, books, and related topics.

You use the following rating levels:

S Tier (Must Consume Original Content Immediately): 18+ ideas and/or STRONG theme matching with the themes in STEP #2.
A Tier (Should Consume Original Content): 15+ ideas and/or GOOD theme matching with the THEMES in STEP #2.
B Tier (Consume Original When Time Allows): 12+ ideas and/or DECENT theme matching with the THEMES in STEP #2.
C Tier (Maybe Skip It): 10+ ideas and/or SOME theme matching with the THEMES in STEP #2.
D Tier (Definitely Skip It): Few quality ideas and/or little theme matching with the THEMES in STEP #2.

Also provide a score between 1 and 100 for the overall quality ranking, where 100 is a perfect match with the highest number of high quality ideas, and 1 is the worst match with a low number of the worst ideas.

The output should look like the following:

LABELS:

Cybersecurity, Writing, Running, Copywriting 

RATING:

S Tier: (Must Consume Original Content Immediately)

Explanation: $$Explanation in 5 short bullets for why you gave that rating.$$

QUALITY SCORE:

$$The 1-100 quality score$$

Explanation: $$Explanation in 5 short bullets for why you gave that score.$$

""",
    }
    user_message = {
        "role": "user",
        "content": """

CONTENT:

        """,
    }

    messages = [system_message, {"role": "user", "content": input_data}]

    try:
        response = openai.ChatCompletion.create(
            model="gpt-4-1106-preview",
            messages=messages,
            temperature=0.0,
            top_p=1,
            frequency_penalty=0.1,
            presence_penalty=0.1,
        )

        assistant_message = response["choices"][0]["message"]["content"]
        return jsonify({"response": assistant_message})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


# Run the application
if __name__ == "__main__":
    app.run(host="1.1.1.1", port=13337, debug=True)
