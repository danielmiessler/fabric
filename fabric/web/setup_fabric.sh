#!/bin/bash

# Make the scripts executable
chmod +x start_fabric.sh stop_fabric.sh

# Define paths
DESKTOP_PATH="/Users/jmdb/Desktop"
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Create Start FABRIC script
cat > "$DESKTOP_PATH/Start FABRIC.command" << EOL
#!/bin/bash
osascript "$SCRIPT_DIR/start_fabric.applescript"
EOL
chmod +x "$DESKTOP_PATH/Start FABRIC.command"

# Create Stop FABRIC script
cat > "$DESKTOP_PATH/Stop FABRIC.command" << EOL
#!/bin/bash
osascript "$SCRIPT_DIR/stop_fabric.applescript"
EOL
chmod +x "$DESKTOP_PATH/Stop FABRIC.command"

echo "FABRIC launcher scripts have been created on your Desktop!"
echo "You can now use:"
echo "- Start FABRIC.command (on Desktop) to launch the development server and browser"
echo "- Stop FABRIC.command (on Desktop) to shut down the server and close the browser"
