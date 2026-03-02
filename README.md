# skill-copy

A CLI tool that copies a single [Agent Skill](https://agentskills.io) file from a public GitHub repository into a local agent's skills directory.

Makes it easy to initialize a project with the skills you need.

[![Go Reference](https://pkg.go.dev/badge/github.com/briangershon/skill-copy.svg)](https://pkg.go.dev/github.com/briangershon/skill-copy)

## Usage

    skill-copy <github-tree-url> <destination>

The tool:

1. Validates the folder is a skill (must contain `SKILL.md`)
2. Creates `<destination>/<skill-name>/`
3. Copies all files and subdirectories into it

## Install

**Option 1: Download a pre-built binary** (no Go required)

Download the latest release for your platform from the [GitHub Releases page](https://github.com/briangershon/skill-copy/releases/latest), extract the archive, and move the `skill-copy` binary to a directory on your `PATH`.

**Option 2: Install with Go**

    # first make sure your GOPATH is set
    go install github.com/briangershon/skill-copy@latest

**Option 3: Installation into a Dockerfile**

For an example of `skill-copy` being installed and used in a `Dockerfile`, see [Agent Workspace](https://github.com/briangershon/agent-workspace/blob/main/Dockerfile).

## Example Installing a Skill with `skill-copy`

Each agent has a different skills directory. For example, Claude's skills directory is `~/.claude/skills` for the current user.

So if you want to install `skill-creator` Skill for Claude, you would run:

```bash
skill-copy https://github.com/anthropics/skills/tree/main/skills/skill-creator ~/.claude/skills
```

This will copy it into `~/.claude/skills`.

## Development

### Build for local development

    go build -o skill-copy .

### Publishing a new release

Merge changes to `main`, then tag and push:

    make tag TAG=v1.2.3

This creates the git tag and pushes it to origin. Once pushed, `go install github.com/briangershon/skill-copy@latest` will resolve to the new tag, or visit Releases on GitHub to download the pre-built binaries.
