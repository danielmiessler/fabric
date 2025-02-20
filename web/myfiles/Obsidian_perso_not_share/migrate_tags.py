import os
import json
from typing import Dict, Any

def load_pattern_tags() -> Dict[str, Any]:
    """
    Load tags and descriptions from pattern_tags_review.md
    Returns a dictionary with pattern names as keys for exact matching
    """
    script_dir = os.path.dirname(os.path.abspath(__file__))
    tags_path = os.path.join(script_dir, "pattern_tags_review.md")
    tags_dict = {}
    
    with open(tags_path, 'r', encoding='utf-8') as f:
        current_pattern = None
        for line in f:
            line = line.strip()
            if line.startswith("## Pattern:"):
                current_pattern = line.split("Pattern:")[1].strip()
                tags_dict[current_pattern] = {"description": "", "tags": []}
            elif line.startswith("Description:") and current_pattern:
                tags_dict[current_pattern]["description"] = line.split("Description:")[1].strip()
            elif line.startswith("Tags:") and current_pattern:
                # Clean and normalize tag parsing
                tags_part = line.split("Tags:")[1].strip()
                tags = [tag.strip() for tag in tags_part.split(",")]
                tags = [tag for tag in tags if tag]  # Remove empty tags
                tags_dict[current_pattern]["tags"] = tags
                
    return tags_dict

def migrate_tags(dry_run: bool = False) -> None:
    """
    Migrate tags to pattern_descriptions.json with detailed verification
    dry_run: If True, only prints changes without writing files
    """
    script_dir = os.path.dirname(os.path.abspath(__file__))
    descriptions_path = os.path.join(script_dir, "pattern_descriptions.json")
    
    # Load existing descriptions
    with open(descriptions_path, 'r', encoding='utf-8') as f:
        descriptions = json.load(f)
    
    # Load and verify tags
    pattern_tags = load_pattern_tags()
    
    # Track changes for reporting
    patterns_updated = []
    patterns_not_found = []
    patterns_no_tags = []
    
    # Update descriptions with tags
    for pattern in descriptions["patterns"]:
        pattern_name = pattern["patternName"]
        if pattern_name in pattern_tags:
            new_tags = pattern_tags[pattern_name]["tags"]
            if new_tags:
                old_tags = pattern.get("tags", [])
                pattern["tags"] = new_tags
                patterns_updated.append({
                    "name": pattern_name,
                    "old_tags": old_tags,
                    "new_tags": new_tags
                })
            else:
                patterns_no_tags.append(pattern_name)
        else:
            patterns_not_found.append(pattern_name)
    
    # Print detailed migration report
    print("\nMigration Report:")
    print("-----------------")
    print(f"Total patterns in descriptions: {len(descriptions['patterns'])}")
    print(f"Total patterns with tag data: {len(pattern_tags)}")
    print(f"Patterns updated: {len(patterns_updated)}")
    
    print("\nDetailed Updates:")
    for p in patterns_updated:
        print(f"\nPattern: {p['name']}")
        print(f"  Old tags: {p['old_tags']}")
        print(f"  New tags: {p['new_tags']}")
    
    if patterns_not_found:
        print("\nPatterns without tag data:")
        for p in patterns_not_found:
            print(f"  - {p}")
    
    if patterns_no_tags:
        print("\nPatterns with empty tags:")
        for p in patterns_no_tags:
            print(f"  - {p}")
    
    if not dry_run:
        # Save updated descriptions
        with open(descriptions_path, 'w', encoding='utf-8') as f:
            json.dump(descriptions, f, indent=2, ensure_ascii=False)
        
        # Update web static copy
        web_static_path = os.path.join(script_dir, "..", "web", "static", "data", "pattern_descriptions.json")
        with open(web_static_path, 'w', encoding='utf-8') as f:
            json.dump(descriptions, f, indent=2, ensure_ascii=False)
        
        print("\nFiles updated:")
        print(f"  - {descriptions_path}")
        print(f"  - {web_static_path}")
    else:
        print("\nDRY RUN - No files were modified")

if __name__ == "__main__":
    # First run in dry-run mode to verify changes
    print("Performing dry run to verify changes...")
    migrate_tags(dry_run=True)
    
    # Prompt for actual migration
    response = input("\nProceed with actual migration? (yes/no): ")
    if response.lower() == 'yes':
        migrate_tags(dry_run=False)
        print("\nMigration completed successfully")
    else:
        print("\nMigration cancelled")

