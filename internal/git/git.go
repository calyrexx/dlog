package git

import (
	"context"
	"fmt"
	"os/exec"
	"path"
	"strings"
)

// Branch returns the currently checked-out git branch, or "" outside a git repo.
func Branch(ctx context.Context) string {
	out, err := run(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return ""
	}

	return out
}

// CommitHash returns the short hash of HEAD, or "" outside a git repo.
func CommitHash(ctx context.Context) string {
	out, err := run(ctx, "git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return ""
	}

	return out
}

// RepoName returns the repository name derived from the remote origin URL,
// or "" if unavailable. Handles both HTTPS and SSH remote formats:
//
//	https://github.com/user/repo.git → "repo"
//	git@github.com:user/repo.git    → "repo"
func RepoName(ctx context.Context) string {
	out, err := run(ctx, "git", "remote", "get-url", "origin")
	if err != nil {
		return ""
	}

	// SSH: git@github.com:user/repo.git → normalize colon to slash.
	if strings.HasPrefix(out, "git@") {
		parts := strings.SplitN(out, ":", 2)
		if len(parts) == 2 {
			out = parts[1]
		}
	}

	name := path.Base(out)
	name = strings.TrimSuffix(name, ".git")

	return name
}

func run(ctx context.Context, name string, args ...string) (string, error) {
	out, err := exec.CommandContext(ctx, name, args...).Output()
	if err != nil {
		return "", fmt.Errorf("exec %s: %w", name, err)
	}

	return strings.TrimSpace(string(out)), nil
}
