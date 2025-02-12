#!/bin/bash

# Change to the web directory
cd "$(dirname "$0")"

# Start backend server in a new terminal window
osascript -e 'tell application "Terminal" to do script "cd \"'$(pwd)'\" && fabric --serve"'

# Wait a moment for backend to start
sleep 3

# Start frontend server in a new terminal window with port 3001
osascript -e 'tell application "Terminal" to do script "cd \"'$(pwd)'\" && PORT=3001 npm run dev"'

# Wait a moment for frontend to start
sleep 3

# Open in default browser
open http://localhost:3001
