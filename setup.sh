#!/bin/bash

# Check if pyproject.toml exists in the current directory
if [ ! -f "pyproject.toml" ]; then
  echo "Poetry could not find a pyproject.toml file in the current directory or its parents."
  echo "Please navigate to the project directory where pyproject.toml is located and rerun this script."
  exit 1
fi

# Installs poetry-based python dependencies
echo "Installing python dependencies"
poetry install

# List of commands to check and add or update alias for
# Add 'yt' and 'ts' to the list of commands
commands=("fabric" "fabric-api" "fabric-webui" "ts" "yt")

# Path to the bootstrap file
bootstrap_file="$HOME/.config/fabric/fabric-bootstrap.inc"

# Ensure the directory for the bootstrap file exists
mkdir -p "$(dirname "$bootstrap_file")"

# Start the bootstrap file with a shebang if it doesn't already exist
if [ ! -f "$bootstrap_file" ]; then
  echo "#!/bin/bash" > "$bootstrap_file"
fi

# List of shell configuration files to update
config_files=("$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.bash_profile")

for config_file in "${config_files[@]}"; do
  # Check if the configuration file exists
  if [ -f "$config_file" ]; then
    echo "Checking $config_file"

    # Ensure the bootstrap script is sourced from the shell configuration file
    source_line="if [ -f \"$bootstrap_file\" ]; then . \"$bootstrap_file\"; fi"
    if ! grep -qF -- "$source_line" "$config_file"; then
      echo -e "\n# Load custom aliases\n$source_line" >> "$config_file"
      echo "Added source command for $bootstrap_file in $config_file."
    fi
    sed -i '/alias fabric=/d' "$config_file"
    sed -i '/fabric --pattern/d' "$config_file"


  else
    echo "$config_file does not exist."
  fi
done

# Add aliases to the bootstrap file
for cmd in "${commands[@]}"; do
  CMD_PATH=$(poetry run which $cmd 2>/dev/null)
  if [ -z "$CMD_PATH" ]; then
    echo "Command $cmd not found in the current Poetry environment."
    continue
  fi
  
  # Check if the alias already exists in the bootstrap file
  if ! grep -qF "alias $cmd=" "$bootstrap_file"; then
    echo "alias $cmd='$CMD_PATH'" >> "$bootstrap_file"
    echo "Added alias for $cmd to $bootstrap_file."
  fi
done

echo "Setup completed. Please restart your terminal or source your configuration files to apply changes."
