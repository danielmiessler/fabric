function __fabric_models
    # Get models from fabric --listmodels
    fabric --listmodels 2>/dev/null | grep '\\[' | cut -f 3
end

function __fabric_patterns
    # Get patterns from fabric --listpatterns
    fabric --listpatterns 2>/dev/null
end

function __fabric_contexts
    # Get contexts from fabric --listcontexts
    fabric --listcontexts 2>/dev/null | grep -v 'No Contexts'
end

function __fabric_sessions
    # Get sessions from fabric --listsessions
    fabric --listsessions 2>/dev/null
end

function __fabric_strategies
    # Get strategies from fabric --liststrategies, add descriptions
    fabric --liststrategies 2>/dev/null | tail -n +2 | sed 's/ \+/\t/' 
end

function __fabric_extensions
    # Get extensions from fabric --listextensions
    fabric --listextensions 2>/dev/null
end

# Main fabric command
complete -c fabric -f

# Pattern selection
complete -c fabric -s p -l pattern -x -a "(__fabric_patterns)" -d "Choose a pattern from the available patterns"

# Variables for patterns
complete -c fabric -s v -l variable -x -d "Values for pattern variables, e.g. -v=#role:expert -v=#points:30"

# Context selection
complete -c fabric -s C -l context -x -a "(__fabric_contexts)" -d "Choose a context from the available contexts"

# Session selection
complete -c fabric -l session -x -a "(__fabric_sessions)" -d "Choose a session from the available sessions"

# Attachments
complete -c fabric -s a -l attachment -r -d "Attachment path or URL (e.g. for OpenAI image recognition messages)"

# Setup
complete -c fabric -s S -l setup -d "Run setup for all reconfigurable parts of fabric"

# Temperature
complete -c fabric -s t -l temperature -x -d "Set temperature (default: 0.7)"

# Top P
complete -c fabric -s T -l topp -x -d "Set top P (default: 0.9)"

# Stream output
complete -c fabric -s s -l stream -d "Stream output"

# Presence penalty
complete -c fabric -s P -l presencepenalty -x -d "Set presence penalty (default: 0.0)"

# Raw mode
complete -c fabric -s r -l raw -d "Use the defaults of the model without sending chat options"

# Frequency penalty
complete -c fabric -s F -l frequencypenalty -x -d "Set frequency penalty (default: 0.0)"

# List patterns
complete -c fabric -s l -l listpatterns -d "List all patterns"

# List models
complete -c fabric -s L -l listmodels -d "List all available models"

# List contexts
complete -c fabric -s x -l listcontexts -d "List all contexts"

# List sessions
complete -c fabric -s X -l listsessions -d "List all sessions"

# Update patterns
complete -c fabric -s U -l updatepatterns -d "Update patterns"

# Copy to clipboard
complete -c fabric -s c -l copy -d "Copy to clipboard"

# Model selection
complete -c fabric -s m -l model -x -a "(__fabric_models)" -d "Choose model"

# Model context length
complete -c fabric -l modelContextLength -x -d "Model context length (only affects ollama)"

# Output to file
complete -c fabric -s o -l output -r -F -d "Output to file"

# Output session
complete -c fabric -l output-session -d "Output the entire session to the output file"

# Latest patterns
complete -c fabric -s n -l latest -x -d "Number of latest patterns to list"

# Change default model
complete -c fabric -s d -l changeDefaultModel -d "Change default model"

# YouTube operations
complete -c fabric -s y -l youtube -x -d "YouTube video or playlist URL to grab content from"
complete -c fabric -l playlist -d "Prefer playlist over video if both ids are present in the URL"
complete -c fabric -l transcript -d "Grab transcript from YouTube video (default)"
complete -c fabric -l transcript-with-timestamps -d "Grab transcript from YouTube video with timestamps"
complete -c fabric -l comments -d "Grab comments from YouTube video"
complete -c fabric -l metadata -d "Output video metadata"

# Language specification
complete -c fabric -s g -l language -x -d "Specify the Language Code for the chat, e.g. -g=en -g=zh"

# Jina AI operations
complete -c fabric -s u -l scrape_url -x -d "Scrape website URL to markdown using Jina AI"
complete -c fabric -s q -l scrape_question -x -d "Search question using Jina AI"

# Seed
complete -c fabric -s e -l seed -x -d "Seed to be used for LMM generation"

# Context and session operations
complete -c fabric -s w -l wipecontext -x -a "(__fabric_contexts)" -d "Wipe context"
complete -c fabric -s W -l wipesession -x -a "(__fabric_sessions)" -d "Wipe session"
complete -c fabric -l printcontext -x -a "(__fabric_contexts)" -d "Print context"
complete -c fabric -l printsession -x -a "(__fabric_sessions)" -d "Print session"

# HTML readability
complete -c fabric -l readability -d "Convert HTML input into a clean, readable view"

# Variables in input
complete -c fabric -l input-has-vars -d "Apply variables to user input"

# Dry run
complete -c fabric -l dry-run -d "Show what would be sent to the model without actually sending it"

# Server options
complete -c fabric -l serve -d "Serve the Fabric Rest API"
complete -c fabric -l serveOllama -d "Serve the Fabric Rest API with ollama endpoints"
complete -c fabric -l address -x -d "The address to bind the REST API (default: :8080)"
complete -c fabric -l api-key -x -d "API key used to secure server routes"

# Config file
complete -c fabric -l config -r -F -d "Path to YAML config file"

# Version
complete -c fabric -l version -d "Print current version"

# Extensions
complete -c fabric -l listextensions -d "List all registered extensions"
complete -c fabric -l addextension -r -F -d "Register a new extension from config file path"
complete -c fabric -l rmextension -x -a "(__fabric_extensions)" -d "Remove a registered extension by name"

# Strategy
complete -c fabric -l strategy -x -a "(__fabric_strategies)" -d "Choose a strategy from the available strategies"
complete -c fabric -l liststrategies -d "List all strategies"
