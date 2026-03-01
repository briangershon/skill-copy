package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GitHubEntry struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
	URL         string `json:"url"`
}

func parseGitHubURL(rawURL string) (owner, repo, branch, path string, err error) {
	// Expected: https://github.com/{owner}/{repo}/tree/{branch}/{path...}
	if !strings.HasPrefix(rawURL, "https://github.com/") {
		return "", "", "", "", fmt.Errorf("invalid URL: must start with https://github.com/")
	}

	trimmed := strings.TrimPrefix(rawURL, "https://github.com/")
	parts := strings.SplitN(trimmed, "/tree/", 2)
	if len(parts) != 2 {
		return "", "", "", "", fmt.Errorf("invalid URL: missing /tree/ segment")
	}

	repoParts := strings.SplitN(parts[0], "/", 2)
	if len(repoParts) != 2 {
		return "", "", "", "", fmt.Errorf("invalid URL: missing owner/repo")
	}
	owner = repoParts[0]
	repo = repoParts[1]

	afterTree := parts[1]
	slashIdx := strings.Index(afterTree, "/")
	if slashIdx < 0 {
		return "", "", "", "", fmt.Errorf("invalid URL: missing path after branch")
	}
	branch = afterTree[:slashIdx]
	path = afterTree[slashIdx+1:]

	if owner == "" || repo == "" || branch == "" || path == "" {
		return "", "", "", "", fmt.Errorf("invalid URL: incomplete components")
	}

	return owner, repo, branch, path, nil
}

func listContents(owner, repo, branch, path string) ([]GitHubEntry, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", owner, repo, path, branch)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "skill-copy/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", apiURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var entries []GitHubEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return entries, nil
}

func downloadFile(url, destPath string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "skill-copy/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("downloading %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %d for %s", resp.StatusCode, url)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", destPath, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("writing file %s: %w", destPath, err)
	}

	return nil
}

func copyDir(owner, repo, branch, remotePath, localDir string) error {
	entries, err := listContents(owner, repo, branch, remotePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		destPath := filepath.Join(localDir, entry.Name)

		switch entry.Type {
		case "file":
			fmt.Printf("  copying %s\n", entry.Name)
			if err := downloadFile(entry.DownloadURL, destPath); err != nil {
				return fmt.Errorf("copying file %s: %w", entry.Name, err)
			}
		case "dir":
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("creating directory %s: %w", destPath, err)
			}
			subPath := remotePath + "/" + entry.Name
			if err := copyDir(owner, repo, branch, subPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: skill-copy <github-tree-url> <destination>\n")
		fmt.Fprintf(os.Stderr, "Example: skill-copy https://github.com/anthropics/skills/tree/main/skills/skill-creator ./claude/skills\n")
		os.Exit(1)
	}

	rawURL := os.Args[1]
	dest := os.Args[2]

	owner, repo, branch, path, err := parseGitHubURL(rawURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Validate target is a real skill by checking for SKILL.md
	entries, err := listContents(owner, repo, branch, path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	hasSkillMD := false
	for _, e := range entries {
		if e.Type == "file" && e.Name == "SKILL.md" {
			hasSkillMD = true
			break
		}
	}
	if !hasSkillMD {
		fmt.Fprintf(os.Stderr, "Error: no SKILL.md found at %s — not a valid skill\n", rawURL)
		os.Exit(1)
	}

	// Skill name is the last segment of the path
	skillName := filepath.Base(path)
	destDir := filepath.Join(dest, skillName)

	fmt.Printf("Copying skill '%s' to %s\n", skillName, destDir)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating destination: %v\n", err)
		os.Exit(1)
	}

	if err := copyDir(owner, repo, branch, path, destDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done.\n")
}
