package git

import (
	"context"
	"fmt"
	"os/exec"
	"path"
	"strings"
)

// Branch returns the name of the currently checked-out git branch.
func Branch(ctx context.Context) (string, error) {
	out, err := run(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("get branch: %w", err)
	}

	return out, nil
}

// CommitHash returns the short hash of the latest commit on HEAD.
func CommitHash(ctx context.Context) (string, error) {
	out, err := run(ctx, "git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "", fmt.Errorf("get commit hash: %w", err)
	}

	return out, nil
}

// RepoName returns the repository name (e.g. "dlog") derived from the remote
// origin URL. Handles both HTTPS and SSH remote formats:
//
//	https://github.com/user/repo.git
//	git@github.com:user/repo.git
func RepoName(ctx context.Context) (string, error) {
	out, err := run(ctx, "git", "remote", "get-url", "origin")
	if err != nil {
		return "", fmt.Errorf("get remote url: %w", err)
	}

	// SSH: git@github.com:user/repo.git → normalize colon to slash
	if strings.HasPrefix(out, "git@") {
		parts := strings.SplitN(out, ":", 2)
		if len(parts) == 2 {
			out = parts[1]
		}
	}

	name := path.Base(out)
	name = strings.TrimSuffix(name, ".git")

	return name, nil
}

func run(ctx context.Context, name string, args ...string) (string, error) {
	out, err := exec.CommandContext(ctx, name, args...).Output()
	if err != nil {
		return "", fmt.Errorf("exec command %s: %w", name, err)
	}

	return strings.TrimSpace(string(out)), nil
}
