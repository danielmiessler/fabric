# Pattern Description Update Process

This document explains the process of generating and updating pattern descriptions for the Fabric project.

## Overview

The pattern description generation is a two-step process:
1. Extract pattern information using Python script
2. Generate descriptions using AI

## Step 1: Extract Pattern Information

First, run the extract_patterns.py script to gather information about all patterns:

```bash
cd fabric/PATTERN_DESCRIPTIONS
python extract_patterns.py
```

This script:
1. Scans the patterns directory
2. Reads each pattern's system.md file
3. Extracts key information like identity, purpose, and steps
4. Outputs the extracted data to pattern_extracts.json

## Step 2: Generate Descriptions

After extracting the pattern information:

1. Use an AI model (like GPT-4) to read pattern_extracts.json
2. For each pattern, ask the AI to:
   - Analyze the extracted information
   - Generate a concise, clear description
   - Focus on the pattern's primary purpose and capabilities
   - Keep descriptions consistent in style and length

Example prompt for the AI:
```
Please read the pattern information from pattern_extracts.json and generate concise, clear descriptions for each pattern. Each description should:
- Be 1-2 sentences long
- Focus on what the pattern does and its primary use case
- Use consistent style across all patterns
- Be clear and actionable
```

3. The AI will generate descriptions in the format required by pattern_descriptions.json:
```json
{
  "patterns": [
    {
      "patternName": "pattern_name",
      "description": "Generated description"
    }
  ]
}
```

## Updating Process for New Patterns

When new patterns are added:

1. Run the Extract Script:
```bash
python extract_patterns.py
```
This script will:
- Update pattern_extracts.json with information from all patterns, including new ones
- Automatically detect new patterns by comparing with existing pattern_descriptions.json
- Add placeholder entries in pattern_descriptions.json for new patterns with a description of "NEEDS_DESCRIPTION"

Example of placeholder entries:
```json
{
  "patterns": [
    {
      "patternName": "existing_pattern",
      "description": "Original description"
    },
    {
      "patternName": "new_pattern",
      "description": "NEEDS_DESCRIPTION"
    }
  ]
}
```

2. Generate Descriptions for New Patterns:
   - Look for patterns marked with "NEEDS_DESCRIPTION" in pattern_descriptions.json
   - Use AI to generate descriptions only for these marked patterns
   - The placeholders make it easy to identify which patterns need new descriptions
   - Replace "NEEDS_DESCRIPTION" with the AI-generated descriptions

3. Verify Updates:
   - Check that no "NEEDS_DESCRIPTION" placeholders remain
   - Ensure descriptions are consistent in style with existing ones
   - Verify the JSON format is valid

## Best Practices

1. Always run extract_patterns.py first to ensure you have the latest pattern information
2. Keep descriptions concise and focused on the pattern's primary purpose
3. Maintain consistent style across all descriptions
4. Verify JSON formatting after updates
5. Keep backups of pattern_descriptions.json before major updates

## File Locations

- Extract Script: `/PATTERN_DESCRIPTIONS/extract_patterns.py`
- Pattern Extracts: `/PATTERN_DESCRIPTIONS/pattern_extracts.json`
- Pattern Descriptions: `/PATTERN_DESCRIPTIONS/pattern_descriptions.json`

## Notes

- The extract_patterns.py script automatically handles new patterns
- Pattern descriptions are used by both the CLI and web interface
- Descriptions should be clear enough for users to understand the pattern's purpose without reading the full implementation
