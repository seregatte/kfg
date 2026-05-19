package urlresolve

import "testing"

func TestIsGitHubURL(t *testing.T) {
	tests := []struct {
		name     string
		arg      string
		expected bool
	}{
		// HTTPS GitHub URLs
		{
			name:     "HTTPS GitHub URL",
			arg:      "https://github.com/owner/repo",
			expected: true,
		},
		{
			name:     "HTTPS GitHub URL with path separator",
			arg:      "https://github.com/owner/repo//manifests",
			expected: true,
		},
		{
			name:     "HTTPS GitHub URL with ref parameter",
			arg:      "https://github.com/owner/repo//manifests?ref=v1.0.0",
			expected: true,
		},
		{
			name:     "HTTPS GitHub URL with branch ref",
			arg:      "https://github.com/owner/repo//path?ref=main",
			expected: true,
		},
		// HTTP GitHub URLs
		{
			name:     "HTTP GitHub URL",
			arg:      "http://github.com/owner/repo",
			expected: true,
		},
		{
			name:     "HTTP GitHub URL with path",
			arg:      "http://github.com/owner/repo//path",
			expected: true,
		},
		// Short form (without protocol)
		{
			name:     "GitHub short form with github.com",
			arg:      "github.com/owner/repo",
			expected: true,
		},
		// Non-GitHub URLs
		{
			name:     "Non-GitHub URL",
			arg:      "https://example.com/path",
			expected: false,
		},
		{
			name:     "GitLab URL",
			arg:      "https://gitlab.com/owner/repo",
			expected: false,
		},
		// Local paths
		{
			name:     "Local relative path",
			arg:      "./manifests",
			expected: false,
		},
		{
			name:     "Local absolute path",
			arg:      "/home/user/manifests",
			expected: false,
		},
		{
			name:     "Local path with subdirectory",
			arg:      ".manifests/overlay/dev",
			expected: false,
		},
		{
			name:     "Current directory",
			arg:      ".",
			expected: false,
		},
		{
			name:     "Parent directory",
			arg:      "..",
			expected: false,
		},
		// Edge cases
		{
			name:     "Empty string",
			arg:      "",
			expected: false,
		},
		{
			name:     "URL containing github.com in path",
			arg:      "https://example.com/github.com/repo",
			expected: true, // Contains "github.com" - matches current design
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsGitHubURL(tt.arg)
			if result != tt.expected {
				t.Errorf("IsGitHubURL(%q) = %v, expected %v", tt.arg, result, tt.expected)
			}
		})
	}
}
