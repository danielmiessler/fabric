import os
import json
import shutil

def load_existing_file(filepath):
    """Load existing JSON file or return default structure"""
    if os.path.exists(filepath):
        with open(filepath, 'r', encoding='utf-8') as f:
            return json.load(f)
    return {"patterns": []}

def extract_pattern_info():
    """Extract pattern information and manage both extract and description files"""
    script_dir = os.path.dirname(os.path.abspath(__file__))
    fabric_dir = os.path.dirname(script_dir)
    patterns_dir = os.path.join(fabric_dir, "patterns")
    extracts_path = os.path.join(script_dir, "pattern_extracts.json")
    descriptions_path = os.path.join(script_dir, "pattern_descriptions.json")
    
    # Load existing data
    existing_extracts = load_existing_file(extracts_path)
    existing_descriptions = load_existing_file(descriptions_path)
    
    # Create lookup sets
    existing_extract_names = {p["patternName"] for p in existing_extracts["patterns"]}
    existing_description_names = {p["patternName"] for p in existing_descriptions["patterns"]}
    
    # Track new patterns
    new_extracts = []
    new_descriptions = []
    
    # Process patterns directory
    for dirname in sorted(os.listdir(patterns_dir)):
        if dirname in ['.DS_Store', 'raycast']:
            continue
            
        pattern_path = os.path.join(patterns_dir, dirname)
        system_md_path = os.path.join(pattern_path, "system.md")
        
        if os.path.isdir(pattern_path) and os.path.exists(system_md_path):
            try:
                # Process pattern extracts
                if dirname not in existing_extract_names:
                    with open(system_md_path, 'r', encoding='utf-8') as f:
                        lines = []
                        for i, line in enumerate(f):
                            if i >= 25:
                                break
                            lines.append(line.rstrip())
                        
                        pattern_extract = "\n".join(lines)
                        new_extracts.append({
                            "patternName": dirname,
                            "pattern_extract": pattern_extract
                        })
                        print(f"Added new pattern extract: {dirname}")
                
                # Add placeholder for pattern descriptions
                if dirname not in existing_description_names:
                    new_descriptions.append({
                        "patternName": dirname,
                        "description": "[Description pending - Requires AI generation]"
                    })
                    print(f"Added description placeholder for: {dirname}")
                    
            except Exception as e:
                print(f"Error processing {dirname}: {str(e)}")
    
    # Merge new data with existing
    existing_extracts["patterns"].extend(new_extracts)
    existing_descriptions["patterns"].extend(new_descriptions)
    
    return existing_extracts, existing_descriptions

def update_web_static(descriptions_path):
    """Copy pattern descriptions to web static directory"""
    try:
        script_dir = os.path.dirname(os.path.abspath(__file__))
        fabric_dir = os.path.dirname(script_dir)
        static_dir = os.path.join(fabric_dir, "web", "static", "data")
        
        # Create static/data directory if it doesn't exist
        os.makedirs(static_dir, exist_ok=True)
        
        # Copy pattern descriptions to web static
        static_path = os.path.join(static_dir, "pattern_descriptions.json")
        shutil.copy2(descriptions_path, static_path)
        print(f"Updated web static file: {static_path}")
    except Exception as e:
        print(f"Error updating web static: {str(e)}")

def save_pattern_files():
    """Save both pattern files"""
    script_dir = os.path.dirname(os.path.abspath(__file__))
    extracts_path = os.path.join(script_dir, "pattern_extracts.json")
    descriptions_path = os.path.join(script_dir, "pattern_descriptions.json")
    
    # Load existing descriptions to calculate new patterns
    existing_descriptions = load_existing_file(descriptions_path)
    existing_description_count = len(existing_descriptions["patterns"])
    
    pattern_extracts, pattern_descriptions = extract_pattern_info()
    
    try:
        # Save pattern extracts
        with open(extracts_path, 'w', encoding='utf-8') as f:
            json.dump(pattern_extracts, f, indent=2, ensure_ascii=False)
            
        # Save pattern descriptions
        with open(descriptions_path, 'w', encoding='utf-8') as f:
            json.dump(pattern_descriptions, f, indent=2, ensure_ascii=False)
            
        # Update web static directory
        update_web_static(descriptions_path)
            
        print(f"\nProcessing complete:")
        print(f"Total patterns: {len(pattern_extracts['patterns'])}")
        print(f"New patterns added: {len(pattern_descriptions['patterns']) - existing_description_count}")
        
    except Exception as e:
        print(f"Error saving JSON files: {str(e)}")

if __name__ == "__main__":
    save_pattern_files()
