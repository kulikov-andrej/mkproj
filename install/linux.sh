#!/usr/bin/env bash

set -euo pipefail

REPO="kulikov-andrej/mkproj"
ASSET_NAME="mkproj-linux-amd64"

INSTALL_DIR="${HOME}/.local/bin"
INSTALL_PATH="${INSTALL_DIR}/mkproj"

CONFIG_HOME="${XDG_CONFIG_HOME:-${HOME}/.config}"
TEMPLATES_DIR="${CONFIG_HOME}/mkproj/templates"

API_URL="https://api.github.com/repos/${REPO}/releases/latest"

echo "Installing mkproj..."

for command in curl sha256sum; do
    if ! command -v "$command" >/dev/null 2>&1; then
        echo "Required command not found: $command" >&2
        exit 1
    fi
done

release_json="$(curl -fsSL \
    -H "Accept: application/vnd.github+json" \
    "$API_URL" \
    --connect-timeout 10 \
    --max-time 30
)"

if command -v jq >/dev/null 2>&1; then
    tag="$(
        printf '%s' "$release_json" |
            jq -r '.tag_name'
    )"

    download_url="$(
        printf '%s' "$release_json" |
            jq -r \
                --arg name "$ASSET_NAME" \
                '.assets[] | select(.name == $name) | .browser_download_url'
    )"

    digest="$(
        printf '%s' "$release_json" |
            jq -r \
                --arg name "$ASSET_NAME" \
                '.assets[] | select(.name == $name) | .digest'
    )"
elif command -v python3 >/dev/null 2>&1; then
    readarray -t release_info < <(
        printf '%s' "$release_json" |
            python3 -c '
import json
import sys

asset_name = sys.argv[1]
release = json.load(sys.stdin)

for asset in release.get("assets", []):
    if asset.get("name") == asset_name:
        print(release.get("tag_name", ""))
        print(asset.get("browser_download_url", ""))
        print(asset.get("digest", ""))
        break
' "$ASSET_NAME"
    )

    tag="${release_info[0]:-}"
    download_url="${release_info[1]:-}"
    digest="${release_info[2]:-}"
else
    echo "jq or python3 is required to read release metadata." >&2
    exit 1
fi

if [[ -z "$download_url" ]]; then
    echo "Release asset not found: $ASSET_NAME" >&2
    exit 1
fi

if [[ "$digest" != sha256:* ]]; then
    echo "Release asset has no supported SHA-256 digest." >&2
    exit 1
fi

expected_hash="${digest#sha256:}"

temp_file="$(mktemp)"
trap 'rm -f "$temp_file"' EXIT

echo "Downloading ${tag}..."

curl -fL \
    "$download_url" \
    -o "$temp_file"

actual_hash="$(
    sha256sum "$temp_file" |
        awk '{ print $1 }'
)"

if [[ "${actual_hash,,}" != "${expected_hash,,}" ]]; then
    echo "SHA-256 verification failed." >&2
    exit 1
fi

echo
echo "Installed:"
echo "  $INSTALL_PATH"

mkdir -p "$INSTALL_DIR"
install -m 0755 "$temp_file" "$INSTALL_PATH"

if [[ ! -d "$TEMPLATES_DIR" ]]; then
    mkdir -p "$TEMPLATES_DIR"

    echo
    echo "Created template directory:"
    echo "  $TEMPLATES_DIR"
fi

"$INSTALL_PATH" --version

case ":${PATH}:" in
    *":${INSTALL_DIR}:"*)
        ;;
    *)
        echo
        echo "'$INSTALL_DIR' is not in PATH."
        echo "Add it to your shell configuration to use 'mkproj' directly."
        ;;
esac