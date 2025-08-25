package cmd

import "testing"

func TestEnabledApiToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		enabled  bool
		expected string
	}{
		{
			name:     "enabled true returns True",
			enabled:  true,
			expected: "True",
		},
		{
			name:     "enabled false returns False",
			enabled:  false,
			expected: "False",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := enabledApiToString(tt.enabled)
			if result != tt.expected {
				t.Errorf("enabledApiToString(%v) = %q, want %q", tt.enabled, result, tt.expected)
			}
		})
	}
}

func TestCommandApiToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		command  []string
		expected string
	}{
		{
			name:     "empty slice returns empty brackets",
			command:  []string{},
			expected: "[]",
		},
		{
			name:     "single command",
			command:  []string{"echo"},
			expected: `["echo"]`,
		},
		{
			name:     "multiple commands",
			command:  []string{"ls", "-la", "/tmp"},
			expected: `["ls","-la","/tmp"]`,
		},
		{
			name:     "commands with spaces and special characters",
			command:  []string{"echo", "hello world", "test & more"},
			expected: `["echo","hello world","test & more"]`,
		},
		{
			name:     "commands with quotes",
			command:  []string{`echo "quoted"`, "test's apostrophe"},
			expected: `["echo \"quoted\"","test's apostrophe"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := commandApiToString(tt.command)
			if result != tt.expected {
				t.Errorf("commandApiToString(%v) = %q, want %q", tt.command, result, tt.expected)
			}
		})
	}
}
