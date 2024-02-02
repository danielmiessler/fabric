#!/usr/bin/python3

import sys
import json
import requests


def send_request(prompt):
    url = "http://hostorip.tld:13337/extwis"
    headers = {
        "Content-Type": "application/json",
        "Authorization": "eJ4f1e0b-25wO-47f9-97ec-6b5335b2",
    }
    data = json.dumps({"input": prompt})
    response = requests.post(url, headers=headers, data=data)

    try:
        print(response.json()["response"])
    except KeyError:
        print("Error: The API response does not contain a 'response' key.")
        print("Received response:", response.json())


if __name__ == "__main__":
    if len(sys.argv) > 1:
        prompt = " ".join(sys.argv[1:])
        send_request(prompt)
    else:
        prompt = sys.stdin.read()
        send_request(prompt)
