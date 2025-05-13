package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigLoad(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pushmonitor-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testConfig := `
push_tests:
  - name: "test1"
    url: "http://example.com/test1"
    interval: 60
    retries: 3
    command: "echo test"
  - name: "test2"
    url: "http://example.com/test2"
    interval: 120
    retries: 2
logging:
  file: "/var/log/pushmonitor.log"
  level: "info"
timeout: 30
`
	testConfigPath := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(testConfigPath, []byte(testConfig), 0644)
	assert.NoError(t, err)

	// Test cases
	tests := []struct {
		name            string
		setup           func() error
		cleanup         func()
		configLocations []string
		expectError     bool
	}{
		{
			name: "successful load from custom path",
			setup: func() error {
				return nil
			},
			cleanup:         func() {},
			configLocations: []string{testConfigPath},
			expectError:     false,
		},
		{
			name: "file not found",
			setup: func() error {
				return nil
			},
			cleanup:         func() {},
			configLocations: []string{"/nonexistent/path/config.yaml"},
			expectError:     true,
		},
		{
			name: "invalid yaml",
			setup: func() error {
				// Write invalid YAML
				invalidConfig := `invalid: yaml: content:`
				return os.WriteFile(testConfigPath, []byte(invalidConfig), 0644)
			},
			cleanup:         func() {},
			configLocations: []string{testConfigPath},
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				err := tt.setup()
				assert.NoError(t, err)
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			cfg := &Config{}
			err := cfg.Load(tt.configLocations...)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, cfg.PushTests)
				assert.Equal(t, 2, len(cfg.PushTests))
				assert.Equal(t, "test1", cfg.PushTests[0].Name)
				assert.Equal(t, "test2", cfg.PushTests[1].Name)
				assert.Equal(t, "/var/log/pushmonitor.log", cfg.Logging.File)
				assert.Equal(t, "info", cfg.Logging.Level)
				assert.Equal(t, 30, cfg.Timeout)
			}
		})
	}
}
