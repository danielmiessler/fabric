#!/usr/bin/env python3
import subprocess
import json
import os
import sys

def get_tree_json(directory):
    """Runs `tree -J` and returns the parsed JSON output."""
    try:
        result = subprocess.run(["tree", "-J", "-L", "3", directory], capture_output=True, text=True)
        return json.loads(result.stdout)
    except Exception as e:
        # Fallback to manual directory traversal if tree command fails
        return manual_directory_scan(directory)

def manual_directory_scan(directory):
    """Manually scan directory structure if tree command is unavailable."""
    tree_json = []
    root_entry = {
        "type": "directory",
        "name": ".",
        "contents": []
    }

    for root, dirs, files in os.walk(directory):
        current_dir = root_entry if root == directory else None

        # Find or create the correct directory entry
        if root != directory:
            relative_path = os.path.relpath(root, directory)
            path_parts = relative_path.split(os.path.sep)

            current_dir = root_entry
            for part in path_parts:
                matching_dirs = [d for d in current_dir.get('contents', []) if d['name'] == part and d['type'] == 'directory']
                if matching_dirs:
                    current_dir = matching_dirs[0]
                else:
                    new_dir = {"type": "directory", "name": part, "contents": []}
                    current_dir.setdefault('contents', []).append(new_dir)
                    current_dir = new_dir

        # Add files to the current directory
        for file in files:
            file_path = os.path.join(root, file)
            file_entry = {
                "type": "file",
                "name": file
            }
            current_dir.setdefault('contents', []).append(file_entry)

    tree_json.append(root_entry)
    return tree_json

def add_file_content(node, base_path):
    """Recursively adds file content to JSON tree."""
    if "type" in node and node["type"] == "file":
        file_path = os.path.abspath(os.path.join(base_path, node["name"]))
        try:
            with open(file_path, "r", encoding="utf-8", errors="ignore") as f:
                node["content"] = f.read()
        except Exception as e:
            node["content"] = f"Error reading file: {str(e)}"
    elif "contents" in node:
        for child in node["contents"]:
            add_file_content(child, os.path.join(base_path, node["name"]))

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 script.py <directory> [instruction]")
        sys.exit(1)

    directory = sys.argv[1]
    instruction_text = sys.argv[2] if len(sys.argv) > 2 else "No instructions provided."

    # Ensure directory exists
    if not os.path.isdir(directory):
        print(f"Error: {directory} is not a valid directory")
        sys.exit(1)

    # Get directory tree JSON
    tree_json = get_tree_json(directory)

    # Add file contents to JSON
    for item in tree_json:
        add_file_content(item, directory)

    # Add instructions to JSON output
    instructions = {
        "type": "instructions",
        "name": "code_change_instructions",
        "details": instruction_text
    }
    tree_json.append(instructions)

    # Add a report section
    report = {
        "type": "report",
        "directories": sum(1 for item in tree_json if item.get('type') == 'directory'),
        "files": sum(1 for item in tree_json for content in item.get('contents', []) if content.get('type') == 'file')
    }
    tree_json.append(report)

    # Output JSON to stdout for piping
    json.dump(tree_json, sys.stdout, indent=4)

if __name__ == "__main__":
    main()
