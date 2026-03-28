package git

// Branch returns the name of the currently checked-out git branch.
// Runs: git rev-parse --abbrev-ref HEAD
func Branch() (string, error) {
	// TODO: exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	return "", nil
}

// CommitHash returns the short hash of the latest commit on HEAD.
// Runs: git rev-parse --short HEAD
func CommitHash() (string, error) {
	// TODO: exec.Command("git", "rev-parse", "--short", "HEAD")
	return "", nil
}

// RepoName returns the repository name derived from the remote origin URL
// or, as a fallback, from the base name of the working directory.
func RepoName() (string, error) {
	// TODO: exec.Command("git", "remote", "get-url", "origin"), then parse URL
	// Fallback: filepath.Base(workdir)
	return "", nil
}
