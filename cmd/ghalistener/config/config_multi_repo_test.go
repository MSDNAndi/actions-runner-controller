package config

import (
	"testing"

	"github.com/actions/actions-runner-controller/apis/actions.github.com/v1alpha1/appconfig"
	"github.com/stretchr/testify/assert"
)

func TestConfig_MultiRepository(t *testing.T) {
	t.Run("GetConfigureUrls returns repositories when set", func(t *testing.T) {
		config := &Config{
			Repositories: []string{
				"https://github.com/user/repo1",
				"https://github.com/user/repo2",
				"https://github.com/user/repo3",
			},
			AppConfig: &appconfig.AppConfig{
				Token: "test-token",
			},
			EphemeralRunnerSetNamespace: "default",
			EphemeralRunnerSetName:      "test",
			RunnerScaleSetName:          "test-scale-set",
		}

		urls := config.GetConfigureUrls()
		assert.Equal(t, 3, len(urls))
		assert.Equal(t, "https://github.com/user/repo1", urls[0])
		assert.Equal(t, "https://github.com/user/repo2", urls[1])
		assert.Equal(t, "https://github.com/user/repo3", urls[2])
	})

	t.Run("GetConfigureUrls returns single URL when ConfigureUrl is set", func(t *testing.T) {
		config := &Config{
			ConfigureUrl: "https://github.com/user/repo",
			AppConfig: &appconfig.AppConfig{
				Token: "test-token",
			},
			EphemeralRunnerSetNamespace: "default",
			EphemeralRunnerSetName:      "test",
			RunnerScaleSetId:            123,
		}

		urls := config.GetConfigureUrls()
		assert.Equal(t, 1, len(urls))
		assert.Equal(t, "https://github.com/user/repo", urls[0])
	})

	t.Run("IsMultiRepository returns true when Repositories is set", func(t *testing.T) {
		config := &Config{
			Repositories: []string{
				"https://github.com/user/repo1",
				"https://github.com/user/repo2",
			},
		}

		assert.True(t, config.IsMultiRepository())
	})

	t.Run("IsMultiRepository returns false when only ConfigureUrl is set", func(t *testing.T) {
		config := &Config{
			ConfigureUrl: "https://github.com/user/repo",
		}

		assert.False(t, config.IsMultiRepository())
	})

	t.Run("Validate fails when neither ConfigureUrl nor Repositories is set", func(t *testing.T) {
		config := &Config{
			AppConfig: &appconfig.AppConfig{
				Token: "test-token",
			},
			EphemeralRunnerSetNamespace: "default",
			EphemeralRunnerSetName:      "test",
			RunnerScaleSetId:            123,
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "either GitHubConfigUrl or Repositories list must be provided")
	})

	t.Run("Validate succeeds with Repositories set", func(t *testing.T) {
		config := &Config{
			Repositories: []string{
				"https://github.com/user/repo1",
				"https://github.com/user/repo2",
			},
			AppConfig: &appconfig.AppConfig{
				Token: "test-token",
			},
			EphemeralRunnerSetNamespace: "default",
			EphemeralRunnerSetName:      "test",
			RunnerScaleSetName:          "test-scale-set",
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Validate fails with empty repository URL", func(t *testing.T) {
		config := &Config{
			Repositories: []string{
				"https://github.com/user/repo1",
				"",
				"https://github.com/user/repo2",
			},
			AppConfig: &appconfig.AppConfig{
				Token: "test-token",
			},
			EphemeralRunnerSetNamespace: "default",
			EphemeralRunnerSetName:      "test",
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty repository URL")
	})

	t.Run("Validate does not require RunnerScaleSetId in multi-repo mode", func(t *testing.T) {
		config := &Config{
			Repositories: []string{
				"https://github.com/user/repo1",
				"https://github.com/user/repo2",
			},
			AppConfig: &appconfig.AppConfig{
				Token: "test-token",
			},
			EphemeralRunnerSetNamespace: "default",
			EphemeralRunnerSetName:      "test",
			RunnerScaleSetName:          "test-scale-set",
			// RunnerScaleSetId intentionally not set
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Validate requires RunnerScaleSetId in single-repo mode", func(t *testing.T) {
		config := &Config{
			ConfigureUrl: "https://github.com/user/repo",
			AppConfig: &appconfig.AppConfig{
				Token: "test-token",
			},
			EphemeralRunnerSetNamespace: "default",
			EphemeralRunnerSetName:      "test",
			// RunnerScaleSetId intentionally not set (0)
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "RunnerScaleSetId")
	})

	t.Run("Validate requires RunnerScaleSetName in multi-repo mode", func(t *testing.T) {
		config := &Config{
			Repositories: []string{
				"https://github.com/user/repo1",
				"https://github.com/user/repo2",
			},
			AppConfig: &appconfig.AppConfig{
				Token: "test-token",
			},
			EphemeralRunnerSetNamespace: "default",
			EphemeralRunnerSetName:      "test",
			// RunnerScaleSetName intentionally not set
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "RunnerScaleSetName is required")
	})
}
