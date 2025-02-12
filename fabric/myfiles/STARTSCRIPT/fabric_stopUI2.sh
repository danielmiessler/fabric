#!/bin/bash

# Kill the frontend server
pkill -f "npm run dev" || true

# Kill the backend server
pkill -f "fabric --serve" || true

# Close any Chrome tabs pointing to the dev server (optional)
osascript -e 'tell application "Google Chrome" to close (tabs of windows whose URL contains "localhost:3001")'
