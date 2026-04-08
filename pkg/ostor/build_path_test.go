package ostor

import (
	"testing"
)

func TestBuildPath(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		query    map[string]string
		expected string
	}{
		{
			name:     "command only",
			cmd:      "ostor-users",
			query:    map[string]string{},
			expected: "/?ostor-users",
		},
		{
			name:     "with email param",
			cmd:      "ostor-users",
			query:    map[string]string{"emailAddress": "user@example.org"},
			expected: "/?ostor-users&emailAddress=user%40example.org",
		},
		{
			name:     "with empty value flag",
			cmd:      "ostor-users",
			query:    map[string]string{"emailAddress": "user@example.org", "disable": ""},
			expected: "/?ostor-users&disable&emailAddress=user%40example.org",
		},
		{
			name:     "multiple flags sorted",
			cmd:      "ostor-users",
			query:    map[string]string{"emailAddress": "user@example.org", "genKey": ""},
			expected: "/?ostor-users&emailAddress=user%40example.org&genKey",
		},
		{
			name:     "nil query",
			cmd:      "ostor-usage",
			query:    nil,
			expected: "/?ostor-usage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildPath(tt.cmd, tt.query)
			if got != tt.expected {
				t.Errorf("buildPath(%q, %v)\n  got:  %s\n  want: %s", tt.cmd, tt.query, got, tt.expected)
			}
		})
	}
}
