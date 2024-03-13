#!/bin/bash

# Check if pyproject.toml exists in the current directory
if [ ! -f "pyproject.toml" ]; then
  echo "Please navigate to the project directory where pyproject.toml is located and rerun this script."
  exit 1
fi

echo "Installing fabric via pipx"
pipx install .
pipx ensurepath

# Path to the bootstrap file
fabric_path="$HOME/.config/fabric/"
bootstrap="fabric-bootstrap.inc"
context="context.md"

bootstrap_file="$fabric_path$bootstrap"
context_file="$fabric_path$context"
# Ensure the directory for the bootstrap file exists
mkdir -p "$(dirname "$bootstrap_file")"

if [ -e "$bootstrap_file" ]; then
  echo "$bootstrap_file exists. Will only append new aliases"
else
  touch $bootstrap_file
  echo "created $bootstrap_file"
fi
if [ -e "$context_file" ]; then
  echo "$context_file exists. Doing nothing"
else
  touch $bootstrap_file
  echo "created $bootstrap_file"
fi

# List of shell configuration files to update
config_files=("$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.bash_profile")

for config_file in "${config_files[@]}"
do
  # Check if the configuration file exists
  if [ -e "$config_file" ]; then
    echo "Checking $config_file"
    # Ensure the bootstrap script is sourced from the shell configuration file
    source_line="if [ -f \"$bootstrap_file\" ]; then . \"$bootstrap_file\"; fi"
    if ! grep -qF -- "$source_line" "$config_file"; then

      echo -e "\n# Load custom aliases for fabric\n$source_line" >> "$config_file"
      echo "Added source command for $bootstrap_file in $config_file."
    fi
  fi
done

# set aliases for all known patterns
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
echo "Setup completed."
