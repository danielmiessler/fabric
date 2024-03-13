#!/bin/bash

# Check if pyproject.toml exists in the current directory
if [ ! -f "pyproject.toml" ]; then
  echo "Please navigate to the project directory where pyproject.toml is located and rerun this script."
  exit 1
fi

echo "Installing fabric via pipx"
pipx install .
pipx ensurepath

PATH=$HOME/.local/bin:$PATH

# Path to the bootstrap file
fabric_path="$HOME/.config/fabric/"
bootstrap="fabric-bootstrap.inc"
context="context.md"
env=".env"

bootstrap_file="$fabric_path$bootstrap"
context_file="$fabric_path$context"
env_file=$fabric_path$env

# Ensure the directory for the bootstrap file exists
mkdir -p "$fabric_path"

# check for the alias file- create if it doesn't exist
if [ -e "$bootstrap_file" ]; then
  echo "$bootstrap_file exists. Will only append new aliases"
else
  touch $bootstrap_file
  echo "created $bootstrap_file"
fi

# check for the context file- create if it doesn't exist
if [ -e "$context_file" ]; then
  echo "$context_file exists. Doing nothing"
else
  touch $context_file
  echo "created $contex_file"
fi

# List of shell configuration files to update
config_files=("$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.bash_profile")

for config_file in "${config_files[@]}"
do
  # Check if the configuration file exists
  if [ -e "$config_file" ]; then
    # we could optionally 'break' after configuring a single file
    echo "Checking $config_file"

    # Ensure the bootstrap script is sourced from the shell configuration file
    source_line="if [ -f \"$bootstrap_file\" ]; then . \"$bootstrap_file\"; fi"
    if ! grep -qF -- "$source_line" "$config_file"; then
      echo -e "\n# Load custom aliases for fabric\n$source_line" >> "$config_file"
      echo "Added source command for $bootstrap_file in $config_file."
    fi
  fi
done

# if you don't have a .env for fabric, fabric --list will tell you to run --setup
# which will break alias creation below
if [ ! -e "$env_file" ]; then
  fabric --setup
fi

# set aliases for all known patterns, skip for existing patterns
echo "Checking alises:"
for i in $(fabric --list)
do
  new="alias $i='fabric -cp $i'"
  if ! grep -qF -- "$new" "$bootstrap_file"; then
    echo $new >> $bootstrap_file
    echo " - $i added to aliases"
  else
    echo " - $i alias exists- skipping"
  fi

done

echo "Please restart your terminal or source $bootstrap_file with '. $bootstrap_file'."
echo "Your personal context file is located at $context_file."
echo "Setup completed."
