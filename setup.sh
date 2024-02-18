#!/bin/bash

# Installs poetry-based python dependencies
echo "Installing python dependencies"
poetry install

# List of commands to check and add or update alias for
commands=("fabric" "fabric-api" "fabric-webui")

# List of shell configuration files to update
config_files=(~/.bashrc ~/.zshrc ~/.bash_profile)

# Initialize an empty string to hold the path of the sourced file
source_command=""

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
        echo -e "\nalias $cmd='$CMD_PATH'" >> "$config_file"
        echo "Added alias for $cmd to $config_file."
      fi
    done
    # Set source_command to source the updated file
    source_command="source $config_file"
  else
    echo "$config_file does not exist."
  fi
done

# Provide instruction to source the updated file
if [ ! -z "$source_command" ]; then
  echo "To apply the changes, please run the following command in your terminal:"
  echo "$source_command"
else
  echo "No configuration files were updated. No need to source."
fi

