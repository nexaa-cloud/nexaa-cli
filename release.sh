#!/bin/bash

set -e

# Check if version is provided as an argument
if [[ -z "$1" ]]; then
  echo "Usage: $0 <version>"
  exit 1
fi

VERSION="$1" # Version argument
PROJECT_ID="65059736" # Set your GitLab project ID
API_BASE_URL="https://gitlab.com/api/v4/projects"
BUILD_DIR="./build"
START_DIR="${PWD}"

# Check if GITLAB_TOKEN is set
if [[ -z "$GITLAB_TOKEN" ]]; then
  echo "Error: Environment variable GITLAB_TOKEN is not set."
  exit 1
fi

# Upload function
upload_binary() {
  local platform=${1%/}
  local binary_name=$2

  if [[ "$platform" == windows_* ]]; then
    binary_name="tilaa.exe"
  fi

  local file_path="${platform}/${binary_name}"

  if [[ -f "$file_path" ]]; then
    echo "Uploading $file_path to GitLab package repository..."

    curl --location \
      --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
      --upload-file "$file_path" \
      "${API_BASE_URL}/${PROJECT_ID}/packages/generic/${platform}/${VERSION}/${binary_name}"

    if [[ $? -ne 0 ]]; then
      echo "Failed to upload $file_path."
    else
      echo "$file_path uploaded successfully."
    fi
  else
    echo "File $file_path does not exist. Skipping."
  fi
}

cd "${BUILD_DIR}"
# Main script
for platform in */; do
  upload_binary "$platform" "tilaa"
done

echo "All uploads completed."
