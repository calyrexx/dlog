package git

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
)

func Branch() (string, error) {
	out, err := run("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("get branch: %w", err)
	}

	return out, nil
}

func CommitHash() (string, error) {
	out, err := run("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "", fmt.Errorf("get commit hash: %w", err)
	}

	return out, nil
}

func RepoName() (string, error) {
	out, err := run("git", "remote", "get-url", "origin")
	if err != nil {
		return "", fmt.Errorf("get remote url: %w", err)
	}

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

func run(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", fmt.Errorf("exec command %s: %w", name, err)
	}

	return strings.TrimSpace(string(out)), nil
}
