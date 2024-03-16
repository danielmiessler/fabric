import argparse
import sys
import requests
from bs4 import BeautifulSoup

# Argument parsing setup
parser = argparse.ArgumentParser(description='Extract text from specified elements within the body of a given webpage, with options for post-processing. Example: --element "p,code" --process "div:decompose=script, p:extract=span"')
parser.add_argument('-u', '--url', help='URL of the webpage to scrape', default=None)
parser.add_argument('-e', '--element', help='Comma-separated list of elements to match within the body.', default='')
parser.add_argument('-p', '--process', help='Comma-separated list of post-process actions by tag and action, e.g., "div:decompose=script, p:extract=span".', default='')
parser.add_argument('-d', '--debug', action='store_true', help='Enable debug mode to print matched element names.')
args = parser.parse_args()

# Function definitions
def get_input_url():
    if not sys.stdin.isatty():
        return sys.stdin.readline().strip()
    return None

def apply_processing(element, actions):
    for action in actions:
        target, command = action.split(':')
        action, tags = command.split('=')
        if element.name == target:
            for tag in tags.split(','):
                for sub_element in element.find_all(tag.strip()):
                    getattr(sub_element, action)()

def process_element_text(element):
    texts = []
    for child in element.children:
        if hasattr(child, 'name') and child.name:
            child_text = child.get_text(strip=True)
            if child_text:
                texts.append(child_text + '\n')
        elif child.string:
            texts.append(child.string.strip())
    return ' '.join(texts)

# Main code logic
if args.url is None:
    args.url = get_input_url()

if not args.url:
    sys.stderr.write("Error: No URL provided. Please provide a URL as an argument or pipe one in.\n")
    sys.exit(1)

tags = set()
if args.element:
    tags.update([tag.strip() for tag in args.element.split(',')])

process_actions = [action.strip() for action in args.process.split(',')] if args.process else []

response = requests.get(args.url)
if response.status_code != 200:
    sys.stderr.write(f"Error: Failed to fetch the URL with status code {response.status_code}.\n")
    sys.exit(1)

soup = BeautifulSoup(response.text, 'html.parser')
body_content = soup.body

for element in body_content.find_all(tags):
    if args.debug:
        print(f"Matched element: {element.name}")
    apply_processing(element, process_actions)
    processed_text = process_element_text(element)
    print(processed_text)
    