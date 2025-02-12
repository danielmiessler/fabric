#!/bin/bash

# Kill all frontend servers
pkill -f "npm run dev" || true

# Kill all backend servers
pkill -f "fabric --serve" || true

# Close any Chrome tabs pointing to any of the dev servers (optional)
osascript -e '
tell application "Google Chrome"
    close (tabs of windows whose URL contains "localhost:5173")
    close (tabs of windows whose URL contains "localhost:3001")
    close (tabs of windows whose URL contains "localhost:8080")
end tell'
