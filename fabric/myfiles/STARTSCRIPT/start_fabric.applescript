 so tell application "Terminal"
	activate
	
	-- Start backend server
	do script "cd /Users/jmdb/Documents/GitHub/FABRIC2/fabric/web && ../fabric --serve"
	
	-- Start frontend server in new window
	do script "cd /Users/jmdb/Documents/GitHub/FABRIC2/fabric/web && npm run dev"
end tell

delay 3

tell application "Google Chrome"
	activate
	open location "http://localhost:5173"
end tell
