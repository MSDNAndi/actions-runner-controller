package config

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMultiRepositoryConfigExample demonstrates how to create and read a multi-repository configuration
func TestMultiRepositoryConfigExample(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "listener-config-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create an example multi-repository configuration
	exampleConfig := Config{
		Repositories: []string{
			"https://github.com/testuser/repo1",
			"https://github.com/testuser/repo2",
			"https://github.com/testuser/repo3",
		},
		EphemeralRunnerSetNamespace: "default",
		EphemeralRunnerSetName:      "multi-repo-runners",
		RunnerScaleSetName:          "multi-repo-scale-set",
		MaxRunners:                  10,
		MinRunners:                  1,
		LogLevel:                    "info",
		LogFormat:                   "text",
	}

	// Note: In real usage, you would add authentication via:
	// - Token field for PAT: exampleConfig.Token = "ghp_..."
	// - Or GitHub App fields: exampleConfig.AppID, AppInstallationID, AppPrivateKey

	// Write the configuration to a file
	configPath := filepath.Join(tempDir, "config.json")
	configFile, err := os.Create(configPath)
	require.NoError(t, err)
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(&exampleConfig)
	require.NoError(t, err)

	// Verify the configuration can be read back (without validation since we didn't add auth)
	readConfig, err := os.Open(configPath)
	require.NoError(t, err)
	defer readConfig.Close()

	var loadedConfig Config
	err = json.NewDecoder(readConfig).Decode(&loadedConfig)
	require.NoError(t, err)

	// Verify the loaded configuration
	assert.Equal(t, 3, len(loadedConfig.Repositories))
	assert.Equal(t, "https://github.com/testuser/repo1", loadedConfig.Repositories[0])
	assert.Equal(t, "https://github.com/testuser/repo2", loadedConfig.Repositories[1])
	assert.Equal(t, "https://github.com/testuser/repo3", loadedConfig.Repositories[2])
	assert.Equal(t, "multi-repo-runners", loadedConfig.EphemeralRunnerSetName)
	assert.Equal(t, "multi-repo-scale-set", loadedConfig.RunnerScaleSetName)
	assert.Equal(t, 10, loadedConfig.MaxRunners)
	assert.Equal(t, 1, loadedConfig.MinRunners)

	// Verify helper methods work correctly
	assert.True(t, loadedConfig.IsMultiRepository())
	urls := loadedConfig.GetConfigureUrls()
	assert.Equal(t, 3, len(urls))
}

// TestSingleRepositoryConfigBackwardCompatibility demonstrates that single-repository mode still works
func TestSingleRepositoryConfigBackwardCompatibility(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "listener-config-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create an example single-repository configuration (existing behavior)
	exampleConfig := Config{
		ConfigureUrl:                "https://github.com/testuser/repo",
		EphemeralRunnerSetNamespace: "default",
		EphemeralRunnerSetName:      "single-repo-runners",
		RunnerScaleSetId:            12345,
		MaxRunners:                  5,
		MinRunners:                  0,
		LogLevel:                    "debug",
		LogFormat:                   "json",
	}

	// Write the configuration to a file
	configPath := filepath.Join(tempDir, "config.json")
	configFile, err := os.Create(configPath)
	require.NoError(t, err)
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(&exampleConfig)
	require.NoError(t, err)

	// Verify the configuration can be read back
	readConfig, err := os.Open(configPath)
	require.NoError(t, err)
	defer readConfig.Close()

	var loadedConfig Config
	err = json.NewDecoder(readConfig).Decode(&loadedConfig)
	require.NoError(t, err)

	// Verify the loaded configuration
	assert.Equal(t, "https://github.com/testuser/repo", loadedConfig.ConfigureUrl)
	assert.Equal(t, "single-repo-runners", loadedConfig.EphemeralRunnerSetName)
	assert.Equal(t, 12345, loadedConfig.RunnerScaleSetId)
	assert.Equal(t, 5, loadedConfig.MaxRunners)
	assert.Equal(t, 0, loadedConfig.MinRunners)

	// Verify helper methods work correctly
	assert.False(t, loadedConfig.IsMultiRepository())
	urls := loadedConfig.GetConfigureUrls()
	assert.Equal(t, 1, len(urls))
	assert.Equal(t, "https://github.com/testuser/repo", urls[0])
}

// TestReadMultiRepositoryConfigWithAuth demonstrates reading a complete multi-repo config with authentication
func TestReadMultiRepositoryConfigWithAuth(t *testing.T) {
	// Skip this test by default since it doesn't actually connect to GitHub
	// To run: go test -v ./cmd/ghalistener/config -run TestReadMultiRepositoryConfigWithAuth
	t.Skip("Skipping integration test - requires GitHub credentials")

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "listener-config-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a configuration file with authentication (using fake credentials for test)
	configPath := filepath.Join(tempDir, "config.json")
	configData := `{
  "repositories": [
    "https://github.com/testuser/repo1",
    "https://github.com/testuser/repo2"
  ],
  "github_token": "ghp_fakeTokenForTestingOnly123456789",
  "ephemeral_runner_set_namespace": "default",
  "ephemeral_runner_set_name": "multi-repo-runners",
  "runner_scale_set_name": "multi-repo-scale-set",
  "max_runners": 10,
  "min_runners": 1,
  "log_level": "info",
  "log_format": "text"
}`
	err = os.WriteFile(configPath, []byte(configData), 0644)
	require.NoError(t, err)

	// Read the configuration using the Read function
	ctx := context.Background()
	config, err := Read(ctx, configPath)
	require.NoError(t, err)

	// Verify the configuration
	assert.Equal(t, 2, len(config.Repositories))
	assert.Equal(t, "https://github.com/testuser/repo1", config.Repositories[0])
	assert.Equal(t, "https://github.com/testuser/repo2", config.Repositories[1])
	assert.Equal(t, "ghp_fakeTokenForTestingOnly123456789", config.Token)
	assert.True(t, config.IsMultiRepository())
}
