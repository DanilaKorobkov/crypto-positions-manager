#!/bin/bash
# Script name: setup.sh
# Description: Copies .env.example to .env if .env does not already exist

TEMPLATE_FILE=".env.example"
TARGET_FILE=".env"

# Check if the target file already exists
if [ ! -f "$TARGET_FILE" ]; then
    echo "File $TARGET_FILE not found. Copying $TEMPLATE_FILE..."
    cp "$TEMPLATE_FILE" "$TARGET_FILE"
    echo "Copy completed."
else
    echo "File $TARGET_FILE already exists. Doing nothing."
fi
