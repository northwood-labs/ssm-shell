#!/usr/bin/env bash

# Remove on the runner.
RUNNER_TEMP="/tmp/terraform-makefile"

# Clone repo into TMP directory.
rm -Rf "${RUNNER_TEMP}"
git clone \
    --depth 1 \
    --branch main \
    --single-branch \
    https://github.com/northwood-labs/.github.git \
    "${RUNNER_TEMP}" \
    ;

# Copy all "full-copy" files from the root into the repository.
find "${RUNNER_TEMP}/full-copy/" -maxdepth 1 -type f -print0 |
    xargs -0 -I% cp -Rfv "%" "${PWD}" ||
    true

# Folders to copy
FOLDERS=(
    ".githooks"
    ".github"
    "scripts"
)

for FOLDER in "${FOLDERS[@]}"; do
    # Copy all files from this directory into the root of the repository.
    mkdir -p "${PWD}/${FOLDER}"
    find "${RUNNER_TEMP}/full-copy/${FOLDER}/" -maxdepth 1 -type f -print0 |
        xargs -0 -I% cp -Rfv "%" "${PWD}/${FOLDER}" ||
        true
done

TYPES=()

# Pass GO=true when calling the script.
# shellcheck disable=2154
if [[ "${GO}" == "true" ]]; then
    TYPES+=("go")
fi

# Pass TF=true when calling the script.
# shellcheck disable=2154
if [[ "${TF}" == "true" ]]; then
    TYPES+=("tf")
fi

for TYPE in "${TYPES[@]}"; do
    # Copy all files from this directory into the root of the repository.
    mkdir -p "${PWD}"
    find "${RUNNER_TEMP}/full-copy/${TYPE}/" -maxdepth 1 -type f -not \( -name "*tmpl*" \) -print0 |
        xargs -0 -I% cp -Rfv "%" "${PWD}" ||
        true
done

# Copy all "updates" files from the root into the repository.
find "${RUNNER_TEMP}/updates/" -maxdepth 1 -type f -not \( -name "*tmpl*" \) -print0 |
    xargs -0 -I% cp -Rfv "%" "${PWD}" ||
    true

# Folders to copy
FOLDERS=(
    ".github"
    ".vscode"
)

for FOLDER in "${FOLDERS[@]}"; do
    # Copy all files from this directory into the root of the repository.
    mkdir -p "${PWD}/${FOLDER}"
    find "${RUNNER_TEMP}/updates/${FOLDER}/" -maxdepth 1 -type f -not \( -name "*tmpl*" \) -print0 |
        xargs -0 -I% cp -Rfv "%" "${PWD}/${FOLDER}" ||
        true
done

TYPES=()

# Pass GO=true when calling the script.
# shellcheck disable=2154
if [[ "${GO}" == "true" ]]; then
    TYPES+=("go")
fi

# Pass TF=true when calling the script.
# shellcheck disable=2154
if [[ "${TF}" == "true" ]]; then
    TYPES+=("tf")
fi

for TYPE in "${TYPES[@]}"; do
    # Copy all files from this directory into the root of the repository.
    mkdir -p "${PWD}"
    find "${RUNNER_TEMP}/updates/${TYPE}/" -maxdepth 1 -type f -not \( -name "*tmpl*" \) -print0 |
        xargs -0 -I% cp -Rfv "%" "${PWD}" ||
        true
done

#  Run Goplicate
goplicate run --allow-dirty --confirm --stash-changes

# Generate .ecrc
tomljson ecrc.toml >.ecrc

# Make shell scripts executable
find "${PWD}" -type f -name "*.sh" -print0 |
    xargs -0 -I% chmod +x "%" ||
    true
