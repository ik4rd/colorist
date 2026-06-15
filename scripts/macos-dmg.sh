#!/usr/bin/env bash

set -euo pipefail

app="${1:?usage: macos-dmg.sh <app-path> <version> [out-dir]}"
version="${2:?usage: macos-dmg.sh <app-path> <version> [out-dir]}"
out_dir="${3:-dist}"

if [[ ! -d "$app" ]]; then
  echo "macos-dmg: app not found: $app" >&2
  exit 1
fi

vol_name="colorist"
dmg_path="$out_dir/colorist-${version}-arm64.dmg"

mkdir -p "$out_dir"
rm -f "$dmg_path"

staging="$(mktemp -d)"
trap 'rm -rf "$staging"' EXIT

cp -R "$app" "$staging/"
ln -s /Applications "$staging/Applications"

hdiutil create \
  -volname "$vol_name" \
  -srcfolder "$staging" \
  -ov \
  -format UDZO \
  "$dmg_path" >/dev/null

echo "$dmg_path"
