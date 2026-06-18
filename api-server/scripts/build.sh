#!/usr/bin/env bash
set -euo pipefail

SETTINGS_FILE="settings/settings.conf.jsonc"
PROJECT_ROOT="$(cd "$("${BASH_SOURCE[0]}")" && pwd)"

# ── 1. Read env_file_name from settings ──────────────────────────────────────
if [[ ! -f "${PROJECT_ROOT}/${SETTINGS_FILE}" ]]; then
  echo "ERROR: Settings file not found at ${PROJECT_ROOT}/${SETTINGS_FILE}" >&2
  exit 1
fi

# Extract the value of env_file_name (handles both 'key: value' and 'key:value')
ENV_FILE_NAME=$(grep -E '^\s*env_file_name\s*:' "${PROJECT_ROOT}/${SETTINGS_FILE}" \
  | sed 's/.*:\s*//' \
  | tr -d "\"'" \
  | tr -d '[:space:]')

if [[ -z "${ENV_FILE_NAME}" ]]; then
  echo "ERROR: 'env_file_name' field not found in ${SETTINGS_FILE}" >&2
  exit 1
fi

if [[ ! -f "${PROJECT_ROOT}/${ENV_FILE_NAME}" ]]; then
  echo "ERROR: env file '${ENV_FILE_NAME}' not found in project root" >&2
  exit 1
fi

# echo "→ Copying '${ENV_FILE_NAME}' to '.env.example' ..."
# cp "${PROJECT_ROOT}/${ENV_FILE_NAME}" "${PROJECT_ROOT}/.env.example"
# echo "  Done."

echo "→ Copying 'settings.conf.jsonc' to 'settings.conf.example.yaml' ..."
cp "${PROJECT_ROOT}/settings/settings.conf.jsonc" "${PROJECT_ROOT}/settings/settings.conf.example.yaml"
echo "  Done."

# ── 2. Build ──────────────────────────────────────────────────────────────────
mkdir -p "${PROJECT_ROOT}/bin"

case "$(uname -s)" in
  MINGW*|MSYS*|CYGWIN*|Windows_NT)
    echo "→ Building secveil-asm (windows/amd64) ..."
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o "${PROJECT_ROOT}/bin/secveil-asm.exe"
    echo "  Build complete → bin/secveil-asm.exe"
    ;;
  *)
    echo "→ Building secveil-asm (linux/amd64) ..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "${PROJECT_ROOT}/bin/secveil-asm"
    echo "  Build complete → bin/secveil-asm"
    ;;
esac