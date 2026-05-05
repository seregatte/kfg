// Package urlresolve provides URL detection and resolution utilities.
package urlresolve

import "strings"

// IsGitHubURL checks if the given argument is a GitHub repository URL.
// It supports:
//   - https://github.com/owner/repo
//   - https://github.com/owner/repo//path
//   - https://github.com/owner/repo//path?ref=branch
//   - http://github.com/... (same patterns)
//   - github.com/owner/repo (without protocol)
func IsGitHubURL(arg string) bool {
	return strings.Contains(arg, "github.com") ||
		strings.HasPrefix(arg, "https://github.com") ||
		strings.HasPrefix(arg, "http://github.com")
}