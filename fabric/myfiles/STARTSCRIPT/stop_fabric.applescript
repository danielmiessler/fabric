-- Kill the frontend server
do shell script "pkill -f 'npm run dev' || true"

-- Kill the backend server
do shell script "pkill -f 'fabric --serve' || true"

-- Close Chrome tabs
tell application "Google Chrome"
	close (tabs of windows whose URL contains "localhost:5173")
end tell
