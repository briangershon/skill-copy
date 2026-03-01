# skill-copy

A CLI tool that copies a single Agent skill file from a GitHub repository folder into a local agent's skills directory.

Makes it easy to initialize a project with the skills you need, such as `skill-creator` for creating new skills.`

## Install

    go install github.com/briangershon/skill-copy@latest

Or copy the binary to a directory on your PATH:

    go build -o skill-copy . && mv skill-copy /usr/local/bin/

## Usage

    skill-copy <github-tree-url> <destination>

The tool:

1. Validates the folder is a skill (must contain `SKILL.md`)
2. Creates `<destination>/<skill-name>/`
3. Copies all files and subdirectories into it

## Build for local development

    go build -o skill-copy .

## Example: Install Anthropic's skill-creator

### Install skills for a project

    skill-copy https://github.com/anthropics/skills/tree/main/skills/skill-creator ./.claude/skills

### Install skills for a user

    skill-copy https://github.com/anthropics/skills/tree/main/skills/skill-creator ~/.claude/skills

This copies the `skill-creator` skill into `~/.claude/skills/skill-creator/`.

## Requirements

- Go 1.16+
- Public GitHub repositories only (no authentication required)
