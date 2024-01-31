#!/usr/bin/python3

import argparse
import json
import requests
import sys
import os
from pathlib import Path


def list_patterns():
    # Gets the location of this script and ensures we resolve the right pattern path
    pattern_directory = f"{Path(__file__).resolve().parent}/../../patterns"  
    try:
        return sorted(os.listdir(pattern_directory))
    except FileNotFoundError:
        print("Pattern directory not found.")
        return []


def send_request(prompt, pattern):
    """sends a pattern and the associated prompt to our general api

    Args:
        prompt (string): stdin string of context we will send to openai
        pattern (string): list a pattern you would like the general client to run
    """
    
    url = "http://localhost:13337/general"
    headers = {
        "Content-Type": "application/json",
        "Authorization": "b246f5c2-6b45-492a-a230-52f2d04b3dc0",
    }
    data = json.dumps({"input": prompt, "pattern": pattern})
    response = requests.post(url, headers=headers, data=data)

    try:
        print(response.json()["response"])
    except KeyError:
        print("Error: The API response does not contain a 'response' key.")
        print("Received response:", response.json())

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Send a request to the API.')
    parser.add_argument('-p', '--pattern', default='extract_wisdom', help='Specify the pattern to use')
    parser.add_argument('-l', '--list', action='store_true', help='List all available patterns')
    
    args = parser.parse_args()

    if args.list:
        patterns = list_patterns()
        print("Available patterns:")
        for pattern in patterns:
            print(pattern)
    else:
        prompt = sys.stdin.read()
        send_request(prompt, args.pattern)
