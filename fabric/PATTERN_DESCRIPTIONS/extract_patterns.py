import os
import json

def extract_pattern_info():
    # Get the directory where the script is located
    script_dir = os.path.dirname(os.path.abspath(__file__))
    # Go up one level to fabric directory
    fabric_dir = os.path.dirname(script_dir)
    patterns_dir = os.path.join(fabric_dir, "patterns")
    
    pattern_info = {"patterns": []}
    
    # Walk through all directories in patterns folder
    for dirname in sorted(os.listdir(patterns_dir)):
        # Skip .DS_Store and raycast directory
        if dirname in ['.DS_Store', 'raycast']:
            continue
            
        pattern_path = os.path.join(patterns_dir, dirname)
        system_md_path = os.path.join(pattern_path, "system.md")
        
        # Check if it's a directory and contains system.md
        if os.path.isdir(pattern_path) and os.path.exists(system_md_path):
            try:
                with open(system_md_path, 'r', encoding='utf-8') as f:
                    # Read first 25 lines
                    lines = []
                    for i, line in enumerate(f):
                        if i >= 25:
                            break
                        lines.append(line.rstrip())
                    
                    # Join lines with newlines
                    pattern_extract = "\n".join(lines)
                    
                    # Add to pattern info
                    pattern_info["patterns"].append({
                        "patternName": dirname,
                        "pattern_extract": pattern_extract
                    })
                    print(f"Processed {dirname}")
            except Exception as e:
                print(f"Error processing {dirname}: {str(e)}")
    
    return pattern_info

def save_pattern_info():
    # Get the pattern information
    pattern_info = extract_pattern_info()
    
    # Get script directory for output path
    script_dir = os.path.dirname(os.path.abspath(__file__))
    output_path = os.path.join(script_dir, "pattern_extracts.json")
    
    try:
        with open(output_path, 'w', encoding='utf-8') as f:
            json.dump(pattern_info, f, indent=2, ensure_ascii=False)
        print(f"\nSuccessfully saved {len(pattern_info['patterns'])} pattern extracts to {output_path}")
    except Exception as e:
        print(f"Error saving JSON file: {str(e)}")

if __name__ == "__main__":
    save_pattern_info()