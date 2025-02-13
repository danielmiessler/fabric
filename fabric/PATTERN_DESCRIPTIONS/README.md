# Pattern Descriptions Management

This document explains how to update pattern descriptions and maintain synchronization between the source files and the web interface.

## Overview

The pattern system follows this hierarchy:
1. `patterns/` directory: The source of truth for available patterns
2. `pattern_descriptions.json`: Generated descriptions for each pattern
3. `web/static/data/pattern_descriptions.json`: Web-accessible copy for the modal interface

The system is managed through:
- `extract_patterns.py`: Python script that scans patterns directory and updates descriptions

## Updating Pattern Descriptions

### Using the Python Script

The `extract_patterns.py` script maintains synchronization between all components:

```bash
# From the fabric/PATTERN_DESCRIPTIONS directory
python extract_patterns.py
```

This will automatically:
1. Scan all pattern directories
2. Extract pattern names and descriptions
3. Update pattern_descriptions.json
4. Copy the file to web/static/data/ for the modal interface

No manual file copying is needed - the script handles the entire synchronization process.

## File Structure

```
fabric/
├── patterns/                    # Source of truth - pattern directories
├── PATTERN_DESCRIPTIONS/
│   ├── extract_patterns.py      # Script to update & sync descriptions
│   ├── pattern_descriptions.json # Generated descriptions
│   └── README.md               # This documentation
└── web/
    └── static/
        └── data/
            └── pattern_descriptions.json # Auto-synced web copy
```

## How It Works

1. Pattern Description Flow:
   ```
   patterns/ directory (source of truth)
          ↓
   extract_patterns.py (scans & processes)
          ↓
   PATTERN_DESCRIPTIONS/pattern_descriptions.json
          ↓ (automatic sync)
   web/static/data/pattern_descriptions.json
          ↓
   Web Modal Display
   ```

2. The web interface:
   - Loads descriptions from `/static/data/pattern_descriptions.json`
   - Uses PatternDescription interface to type the data:
     ```typescript
     interface PatternDescription {
       patternName: string;
       description: string;
     }
     ```
   - Displays patterns in the modal with search and sort capabilities

## Best Practices

1. Always run `extract_patterns.py` when:
   - Adding new patterns to the patterns/ directory
   - Modifying pattern descriptions
   - Removing patterns

2. The script automatically maintains synchronization:
   - Scans patterns/ directory for the complete list
   - Updates pattern_descriptions.json with any changes
   - Syncs to web/static/data/ for the modal interface

3. Test the web interface after updates to ensure:
   - All patterns from patterns/ directory are listed
   - Descriptions are correct
   - Search functionality works
   - Pattern selection works

## Troubleshooting

If patterns are not showing in the web modal:

1. Verify JSON files:
   - Check both JSON files exist
   - Ensure JSON format is valid
   - Compare contents of both files

2. Check web paths:
   - Verify `/static/data/pattern_descriptions.json` is accessible
   - Check browser console for loading errors

3. Clear cache:
   - Clear browser cache
   - Rebuild web application
