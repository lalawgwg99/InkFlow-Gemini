#!/bin/bash
#
# Sync script - Keep skill/ and plugins/ directories in sync
# Usage: ./scripts/sync.sh
#

# Don't exit on error, we want to sync all files
set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Source and destination directories
SRC_DIR="$PROJECT_ROOT/skill/md2wechat"
DST_DIR="$PROJECT_ROOT/plugins/md2wechat/skills/md2wechat"

echo "========================================"
echo "   Syncing Skill Directories"
echo "========================================"
echo ""
echo "Source: $SRC_DIR"
echo "Target: $DST_DIR"
echo ""

# Check if source directory exists
if [[ ! -d "$SRC_DIR" ]]; then
    echo "❌ Error: Source directory not found: $SRC_DIR"
    exit 1
fi

# Create destination directory structure
mkdir -p "$DST_DIR/scripts"
mkdir -p "$DST_DIR/references"

# Files to sync
FILES=(
    "SKILL.md"
    "scripts/run.sh"
    "references/themes.md"
    "references/html-guide.md"
    "references/image-syntax.md"
    "references/wechat-api.md"
    "references/writing-guide.md"
)

# Sync each file
sync_count=0
for file in "${FILES[@]}"; do
    src_path="$SRC_DIR/$file"
    dst_path="$DST_DIR/$file"

    if [[ -f "$src_path" ]]; then
        # Copy file if source exists
        cp "$src_path" "$dst_path"
        echo "✓ Synced: $file"
        ((sync_count++))
    else
        echo "⚠ Skipped: $file (not found in source)"
    fi
done

echo ""
echo "========================================"
echo "   Sync Complete: $sync_count files"
echo "========================================"
echo ""
echo "Both directories are now in sync."
echo "  - skill/md2wechat/"
echo "  - plugins/md2wechat/skills/md2wechat/"
