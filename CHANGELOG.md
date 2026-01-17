# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.5.0] - 2025-01-17

### Added
- **Writer Style Assistant**: AI-powered writing assistance with customizable creator styles
  - New `write` command for assisted article generation
  - Support for multiple input types: idea, fragment, outline, title
  - Automatic cover prompt generation matching writing style
- **Dan Koe Style**: First built-in creator writing style
  - Profound, sharp, grounded tone for personal growth and opinion pieces
  - Complete writing framework with hooks, structures, and quote extraction
  - Victorian Woodcut/Etching style cover generation
- **Custom Style System**: YAML-based style definitions in `writers/` directory
  - Easy to add custom creator styles
  - Configurable writing prompts and cover styles
- **Image Size Control**: New `--size` parameter for `generate_image` command
  - Support for 16:9 ratio (2560x1440) for WeChat cover images
  - Multiple preset sizes: 2048x2048, 1920x1920, 2560x1440, 1440x2560, etc.
- **Documentation**:
  - New `writers/README.md` with custom style guide
  - New `docs/WRITING_FAQ.md` for writing beginners
  - New `references/writing-guide.md` with complete write command reference
  - Updated `references/image-syntax.md` with size parameter documentation
  - Enhanced README.md with write command workflows and diagrams
  - Enhanced SKILL.md with writing assistance natural language examples

### Changed
- README: Added API service notice at top with md2wechat.cn contact information
- README: Updated author section with donation information and QR codes
- README: Added writer style comparison table and workflow diagrams
- SKILL.md: Improved LLM instruction following with explicit trigger conditions
- Sync script: Added `writing-guide.md` to file synchronization list

### Technical Details
- **New Files**:
  - `internal/writer/` - Writer style assistant module
    - `types.go` - Core data structures
    - `style.go` - Style management system
    - `generator.go` - Article generation logic
    - `assistant.go` - High-level assistant API
    - `cover_generator.go` - Cover prompt generation
  - `cmd/md2wechat/write.go` - Write command implementation
  - `writers/dan-koe.yaml` - Dan Koe style configuration
  - `writers/README.md` - Custom style guide
  - `docs/WRITING_FAQ.md` - Writing functionality FAQ
  - `skill/md2wechat/references/writing-guide.md` - Writing command reference

### Breaking Changes
- None

### Migration Guide
No migration required. The write command is a new feature and doesn't affect existing functionality.

---

## [1.4.0] - 2025-01-14

### Added
- **Provider Pattern**: Extensible image generation service architecture
- **TuZi Integration**: Support for TuZi (tu-zi.com) image generation service
  - Models: `doubao-seedream-4-5-251128` (default), `gemini-3-pro-image-preview`
  - Sizes: 2048x2048 (default), 1920x1920, 2560x1440, 1440x2560, 3072x2048, 2048x3072, 3840x2160, 2160x3840
- **Natural Language Image Generation**: Generate images via conversational interface
  - Method 1: Insert into article at specific position
  - Method 2: Standalone image generation
  - Method 3: Manual Markdown syntax `![alt](__generate:prompt__)`
- **Configuration Fields**: `image_provider`, `image_api_base`, `image_model`, `image_size`
- **Documentation**: New `docs/IMAGE_PROVISIONERS.md` with complete provider configuration guide
- **Sync Script**: `scripts/sync.sh` to keep `skill/` and `plugins/` directories synchronized
- **Makefile Target**: `make sync` for easy directory synchronization

### Changed
- README: Enhanced platform-specific download instructions (Mac Intel vs Apple Silicon)
- README: Added AI image generation section with natural language examples
- `scripts/install.sh`: Added Linux ARM64 detection and support
- Default TuZi image size: `1024x1024` → `2048x2048` (TuZi requires minimum 3.7M pixels)

### Fixed
- URL download with query parameters: Fixed file name too long error when downloading generated images
- Platform detection in install script for ARM64 systems
- API base URL documentation: Updated to correct `https://api.tu-zi.com/v1`

### Technical Details
- **New Files**:
  - `internal/image/provider.go` - Provider interface and factory
  - `internal/image/openai.go` - OpenAI DALL-E provider
  - `internal/image/tuzi.go` - TuZi image generation provider
  - `docs/IMAGE_PROVISIONERS.md` - Provider configuration guide
  - `scripts/sync.sh` - Directory synchronization script

---

## [1.3.1] - 2025-01-12

### Added
- Auto-download binary from GitHub releases on first run
- User-friendly error messages with actionable guidance
- System dependencies declaration in `manifest.json`
- XDG-compliant cache directory (`~/.cache/claude/`) for binary storage
- Fallback to local `bin/` directory for development/offline usage
- Automatic version checking and update prompt when binary is outdated

### Changed
- Binary naming format now uses hyphens (`md2wechat-linux-amd64`) to match GitHub releases
- Error messages now show human-readable platform names (e.g., "macOS (Apple Silicon)")
- Download progress displays concise status information

### Fixed
- Binary name mismatch between workflow and run.sh that caused download failures

---

## [1.3.0] - 2025-01-11

### Added
- Plugin Marketplace support with `.claude-plugin/marketplace.json`
- One-command Claude Code installation via `/plugin marketplace add`
- Prominent Claude Code installation section at top of README
- Claude Code badge for quick identification
- Detailed binary installation instructions with location guidance
- Installation steps for Windows, Mac, and Linux users

### Changed
- Updated QUICKSTART.md with Claude Code section at the beginning
- Enhanced docs/USAGE.md with Claude Code integration guide
- Improved download table with installation locations
- Added collapsible installation steps for each platform

### Installation
```bash
# Claude Code users (simplest)
/plugin marketplace add geekjourneyx/md2wechat-skill
/plugin install md2wechat@geekjourneyx-md2wechat-skill
```

---

## [1.2.0] - 2025-01-11

### Added
- Claude Code Skill support with `.claude-plugin/plugin.json`
- Claude Code Skill in `skill/md2wechat/` directory for distribution
- New API themes: `bytedance`, `apple`, `sports`, `chinese`, `cyber`
- Comprehensive troubleshooting guide in SKILL.md
- API theme selection section in README.md
- CHANGELOG.md with version history and upgrade guide

### Changed
- Updated themes.md with complete API and AI theme documentation
- Enhanced HTML guide with AI theme specific requirements
- Improved SKILL.md with detailed error handling and FAQ
- Updated command help text to reflect all available themes
- Enhanced FAQ.md with IP whitelist configuration guide

### Removed
- `leo` theme (deprecated)

---

## [1.1.0] - 2025-01-11

### Added
- YAML-based theme system for AI mode
- AI themes: `autumn-warm`, `spring-fresh`, `ocean-calm`
- Custom theme support with `custom.yaml`
- Theme configuration with color schemes and style info
- Reference documentation for themes, HTML guide, image syntax, and WeChat API

### Changed
- Refactored theme system to use YAML files instead of hardcoded prompts
- Updated SKILL.md with new AI theme workflow
- Enhanced themes.md with detailed style specifications

---

## [1.0.1] - 2025-01-11

### Fixed
- Mermaid diagrams rendering for GitHub documentation

---

## [1.0.0] - 2025-01-11

### Added
- Initial release of md2wechat
- API mode conversion using md2wechat.cn API
- AI mode conversion with Claude AI
- WeChat draft upload functionality
- Image upload to WeChat material library
- Configuration management with YAML support
- Command-line interface with cobra
- Multi-platform binary support (Windows, macOS, Linux)
- Comprehensive documentation (README, FAQ, USAGE)
- Test draft command for HTML validation

### Features
- Convert Markdown to WeChat Official Account formatted HTML
- Support for local images, online images, and AI-generated images
- Automatic image compression and optimization
- Draft creation with cover image support
- Environment variable and config file support
- Claude Code Skill integration

---

## Version History Summary

| Version | Date | Description |
|---------|------|-------------|
| [1.5.0] | 2025-01-17 | Writer style assistant, Dan Koe style, image size control |
| [1.4.0] | 2025-01-14 | TuZi image provider, natural language image generation |
| [1.3.1] | 2025-01-12 | Auto binary download, user-friendly errors, system dependencies |
| [1.3.0] | 2025-01-11 | Plugin Marketplace support, enhanced installation docs |
| [1.2.0] | 2025-01-11 | Claude Code plugin support, new API themes |
| [1.1.0] | 2025-01-11 | YAML theme system, AI themes (autumn-warm, spring-fresh, ocean-calm) |
| [1.0.1] | 2025-01-11 | Fixed Mermaid diagrams for GitHub rendering |
| [1.0.0] | 2025-01-11 | Initial release with full md2wechat functionality |

---

## Upgrade Guide

### From v1.1.0 to Unreleased

**New API Themes Available:**
```bash
# New themes available
md2wechat convert article.md --theme bytedance
md2wechat convert article.md --theme apple
md2wechat convert article.md --theme sports
md2wechat convert article.md --theme chinese
md2wechat convert article.md --theme cyber
```

**Claude Code Skill Integration:**
```bash
# Install as Claude Code Skill
cp -r skill/md2wechat ~/.claude/skills/
```

### From v1.0.0 to v1.1.0

**Theme System Migration:**
- Old theme names have been updated
- AI themes now use YAML configuration files
- Update your commands to use new theme names:
  - `autumn-warm` instead of `elegant`
  - `spring-fresh` instead of `minimal`
  - `ocean-calm` instead of `tech`

---

## Links

- [GitHub Repository](https://github.com/geekjourneyx/md2wechat-skill)
- [Documentation](README.md)
- [Issues](https://github.com/geekjourneyx/md2wechat-skill/issues)
