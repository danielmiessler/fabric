#!/bin/bash

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <YouTube URL> <Keyword>"
    exit 1
fi

# Assign input arguments to variables
youtube_url=$1
keyword=$2

# Generate file names based on the keyword
path_to_artefacts="../artifacts/"
transcript_file="${keyword}_transcript.md"
summary_file="${keyword}_Summary.md"
essay_file="${keyword}_essay.md"

# Execute the commands with the input arguments
echo "Geting transcript"
yt --transcript "$youtube_url" > "$path_to_artefacts$transcript_file"
echo "Transcript saved to $path_to_artefacts$transcript_file"

echo "Extracting wisdom"
cat "$path_to_artefacts$transcript_file" | fabric --model gpt-4o -sp extract_wisdom > "$path_to_artefacts$summary_file"
echo "Summary saved to $path_to_artefacts$summary_file"

echo "Writing essay"
cat "$path_to_artefacts$summary_file" | fabric --model gpt-4o -sp write_essay > "$path_to_artefacts$essay_file"

echo "Essay saved to $path_to_artefacts$essay_file"
