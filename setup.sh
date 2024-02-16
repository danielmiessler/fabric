#!/bin/bash

# Installs poetry-based python dependencies
echo "Installing python dependencies"
poetry install

# List of commands to check and add or update alias for
commands=("fabric" "fabric-api" "fabric-webui")

# List of shell configuration files to update
config_files=(~/.bashrc ~/.zshrc)

for config_file in "${config_files[@]}"; do
  # Check if the configuration file exists
  if [ -f "$config_file" ]; then
    echo "Updating $config_file"
    for cmd in "${commands[@]}"; do
      # Get the path of the command
      CMD_PATH=$(poetry run which $cmd)

      # Check if the config file contains an alias for the command
      if grep -q "alias $cmd=" "$config_file"; then
        # Replace the existing alias with the new one
        sed -i "/alias $cmd=/c\alias $cmd='$CMD_PATH'" "$config_file"
        echo "Updated alias for $cmd in $config_file."
      else
        # If not, add the alias to the config file
        echo "alias $cmd='$CMD_PATH'" >> "$config_file"
        echo "Added alias for $cmd to $config_file."
      fi
    done
  else
    echo "$config_file does not exist."
  fi
done

echo "Please close this terminal window to have new aliases work."