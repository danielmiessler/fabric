#!/bin/zsh

# Set the base directory
BASE_DIR="${HOME}/.config/fabric/patterns"
OUTPUT_DIR="${BASE_DIR}/suggest_pattern"

# Create the output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# If the file doesn't exist, create the user.md file
if [[ ! -f "$OUTPUT_DIR/user.md" ]]; then
		echo "Creating user.md file..."
		echo "CONTENT: \n\n" > $OUTPUT_DIR/user.md
fi

# Process README.md for command formation info
README_FILE="$BASE_DIR/../README.md"
# Skip if user.md has "# OVERVIEW" section
if grep -q "# OVERVIEW" $OUTPUT_DIR/user.md; then
		echo "Skipping README.md..."
elif [[ -f "$README_FILE" ]]; then
    echo "Processing README.md for command formation info..."
    
    # Extract relevant sections from README.md
    fabric -p explain_docs < "$README_FILE" >> $OUTPUT_DIR/user.md
fi

# If the file doesn't include "# PATTERNS" section, add it
if ! grep -q "# PATTERNS" $OUTPUT_DIR/user.md; then
		echo "\n\n# PATTERNS\n" >> $OUTPUT_DIR/user.md
fi

# Function to process each pattern directory
process_pattern() {
    local pattern_name=$(basename "$1")
    local system_file="$1/system.md"
		local readme_file="$1/README.md"
    
		# Skip if pattern_name is already in the user.md file
		if grep -q "# $pattern_name" $OUTPUT_DIR/user.md; then
				echo "Skipping $pattern_name..."
				return
		fi

		# Use the README file if it exists
		if [[ -f "$readme_file" ]]; then
				echo "Processing $pattern_name (README.md)..."
				
				# Run fabric summarize on system.md and append to user.md
				local summary=$(fabric -p summarize_prompt < "$readme_file")
				echo "## $pattern_name\n$summary\n" >> $OUTPUT_DIR/user.md
		elif [[ -f "$system_file" ]]; then
				echo "Processing $pattern_name (system.md)..."
				
				# Run fabric summarize on system.md and append to user.md
				local summary=$(fabric -p summarize_prompt < "$system_file")
				echo "## $pattern_name\n$summary\n" >> $OUTPUT_DIR/user.md
		fi
}

# Process all pattern directories
for dir in $BASE_DIR/*/; do
    process_pattern "$dir"
done

echo "Processing complete. Context file updated at $OUTPUT_DIR/user.md"