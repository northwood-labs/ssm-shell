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
FILES="$(find "${RUNNER_TEMP}/full-copy/" -maxdepth 1 -type f)"

# Files that should only be copied the first time. Do not overwrite on
# subsequent copies.
ONE_TIME_ONLY=(
    ".markdownlint.jsonc"
)

# shellcheck disable=2068
for FILE in ${FILES[@]}; do
    for IGNORE in "${ONE_TIME_ONLY[@]}"; do
        # If the file does not exist, go ahead and copy it (first time)
        if [[ ! -f "${PWD}/${IGNORE}" ]]; then
            cp -Rfv "${FILE}" "${PWD}"

        # Otherwise, as long as the copied file is not the ignored file, go
        # ahead and copy it (no restricton)
        elif [[ "${FILE}" != "${RUNNER_TEMP}/full-copy/${IGNORE}" ]]; then
            cp -Rfv "${FILE}" "${PWD}"
        fi
    done
done

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
    find "${RUNNER_TEMP}/full-copy/${TYPE}/" -maxdepth 1 -type f -print0 |
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
