#!/usr/bin/env bash
set -euo pipefail

APP_NAME="Mado-Tray"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP_PATH="${ROOT_DIR}/build/bin/${APP_NAME}.app"

log() {
  printf "\n==> %s\n" "$1"
}

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    printf "Missing required command: %s\n" "$1" >&2
    exit 1
  fi
}

cd "${ROOT_DIR}"

require_command go
require_command npm
require_command wails

log "Installing frontend dependencies"
npm install

log "Tidying Go modules"
go mod tidy

log "Building frontend"
npm run build

log "Building ${APP_NAME}.app with Wails"
wails build

if [[ ! -d "${APP_PATH}" ]]; then
  printf "Build finished, but app bundle was not found at: %s\n" "${APP_PATH}" >&2
  exit 1
fi

log "Build complete"
printf "App bundle: %s\n" "${APP_PATH}"
printf "Tip: move it to /Applications before enabling startup from Options.\n"
