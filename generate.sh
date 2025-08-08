#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <project-name>"
    exit 1
fi

PROJECT_NAME=$1

if [ ! -d "placeholder" ]; then
    echo "Error: 'placeholder' directory not found."
    exit 1
fi

cp -r placeholder "${PROJECT_NAME}"
find "${PROJECT_NAME}" -type f -exec sed -i "s/placeholder/${PROJECT_NAME}/g" {} +
