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
commands=("fabric" "fabric-api" "fabric-webui" "ts", "yt")

# List of shell configuration files to update
config_files=("$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.bash_profile")

# Initialize an array to hold the paths of the sourced files
source_commands=()

for config_file in "${config_files[@]}"; do
  # Check if the configuration file exists
  if [ -f "$config_file" ]; then
    echo "Updating $config_file"
    for cmd in "${commands[@]}"; do
      # Get the path of the command
      CMD_PATH=$(poetry run which $cmd 2>/dev/null)

      # Check if CMD_PATH is empty
      if [ -z "$CMD_PATH" ]; then
        echo "Command $cmd not found in the current Poetry environment."
        continue
      fi

      # Check if the config file contains an alias for the command
      if grep -qE "alias $cmd=|alias $cmd =" "$config_file"; then
        # Compatibility with GNU and BSD sed: Check for operating system and apply appropriate sed syntax
        if [[ "$OSTYPE" == "darwin"* ]]; then
          # BSD sed (macOS)
          sed -i '' "/alias $cmd=/c\\
alias $cmd='$CMD_PATH'" "$config_file"
        else
          # GNU sed (Linux and others)
          sed -i "/alias $cmd=/c\alias $cmd='$CMD_PATH'" "$config_file"
        fi
        echo "Updated alias for $cmd in $config_file."
      else
        # If not, add the alias to the config file
        echo -e "\nalias $cmd='$CMD_PATH'" >>"$config_file"
        echo "Added alias for $cmd to $config_file."
      fi
    done
    # Add to source_commands array
    source_commands+=("$config_file")
  else
    echo "$config_file does not exist."
  fi
done

# Provide instruction to source the updated files
if [ ${#source_commands[@]} -ne 0 ]; then
  echo "To apply the changes, please run the following command(s) in your terminal:"
  for file in "${source_commands[@]}"; do
    echo "source $file"
  done
else
  echo "No configuration files were updated. No need to source."
fi
