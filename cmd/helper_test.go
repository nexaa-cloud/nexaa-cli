package cmd

import (
	"testing"

	"github.com/nexaa-cloud/nexaa-cli/api"
)

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

func TestEnvsToApi(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		envs         []string
		secret       bool
		state        api.State
		expectedEnvs []api.EnvironmentVariableInput
	}{
		{
			name:   "single environment variable, not secret",
			envs:   []string{"KEY=value"},
			secret: false,
			state:  api.StatePresent,
			expectedEnvs: []api.EnvironmentVariableInput{
				{
					Name:   "KEY",
					Value:  "value",
					Secret: false,
					State:  api.StatePresent,
				},
			},
		},
		{
			name:   "single environment variable, secret",
			envs:   []string{"PASSWORD=secret123"},
			secret: true,
			state:  api.StatePresent,
			expectedEnvs: []api.EnvironmentVariableInput{
				{
					Name:   "PASSWORD",
					Value:  "secret123",
					Secret: true,
					State:  api.StatePresent,
				},
			},
		},
		{
			name:   "multiple environment variables",
			envs:   []string{"KEY1=value1", "KEY2=value2"},
			secret: false,
			state:  api.StatePresent,
			expectedEnvs: []api.EnvironmentVariableInput{
				{
					Name:   "KEY1",
					Value:  "value1",
					Secret: false,
					State:  api.StatePresent,
				},
				{
					Name:   "KEY2",
					Value:  "value2",
					Secret: false,
					State:  api.StatePresent,
				},
			},
		},
		{
			name:         "empty environment variables",
			envs:         []string{},
			secret:       false,
			state:        api.StatePresent,
			expectedEnvs: []api.EnvironmentVariableInput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := envsToApi(tt.envs, tt.secret, tt.state)

			if len(result) != len(tt.expectedEnvs) {
				t.Errorf("envsToApi() returned %d items, want %d", len(result), len(tt.expectedEnvs))
			}

			for _, got := range result {
				var want api.EnvironmentVariableInput
				for _, expectedEnv := range tt.expectedEnvs {
					if expectedEnv.Name == got.Name {
						want = expectedEnv
						break
					}
				}

				if got.Name != want.Name {
					t.Errorf("envsToApi() Name = %v, want %v", got.Name, want.Name)
				}
				if got.Value != want.Value {
					t.Errorf("envsToApi() Value = %v, want %v", got.Value, want.Value)
				}
				if got.Secret != want.Secret {
					t.Errorf("envsToApi() Secret = %v, want %v", got.Secret, want.Secret)
				}
				if got.State != want.State {
					t.Errorf("envsToApi() State = %v, want %v", got.State, want.State)
				}
			}
		})
	}
}
