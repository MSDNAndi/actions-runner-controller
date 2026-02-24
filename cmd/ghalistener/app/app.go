package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/actions/actions-runner-controller/cmd/ghalistener/config"
	"github.com/actions/actions-runner-controller/cmd/ghalistener/listener"
	"github.com/actions/actions-runner-controller/cmd/ghalistener/metrics"
	"github.com/actions/actions-runner-controller/cmd/ghalistener/worker"
	"github.com/actions/actions-runner-controller/github/actions"
	"github.com/go-logr/logr"
	"golang.org/x/sync/errgroup"
)

// App is responsible for initializing required components and running the app.
type App struct {
	// configured fields
	config *config.Config
	logger logr.Logger

	// initialized fields
	listener  Listener   // Single listener (used in single-repo mode)
	listeners []Listener // Multiple listeners (used in multi-repo mode)
	worker    Worker
	metrics   metrics.ServerExporter
}

//go:generate mockery --name Listener --output ./mocks --outpkg mocks --case underscore
type Listener interface {
	Listen(ctx context.Context, handler listener.Handler) error
}

//go:generate mockery --name Worker --output ./mocks --outpkg mocks --case underscore
type Worker interface {
	HandleJobStarted(ctx context.Context, jobInfo *actions.JobStarted) error
	HandleDesiredRunnerCount(ctx context.Context, count int, jobsCompleted int) (int, error)
}

func New(config config.Config) (*App, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	app := &App{
		config: &config,
	}

	{
		logger, err := config.Logger()
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %w", err)
		}
		app.logger = logger.WithName("listener-app")
	}

	worker, err := worker.New(
		worker.Config{
			EphemeralRunnerSetNamespace: config.EphemeralRunnerSetNamespace,
			EphemeralRunnerSetName:      config.EphemeralRunnerSetName,
			MaxRunners:                  config.MaxRunners,
			MinRunners:                  config.MinRunners,
		},
		worker.WithLogger(app.logger.WithName("worker")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new kubernetes worker: %w", err)
	}
	app.worker = worker

	if config.IsMultiRepository() {
		// Multi-repository mode
		app.logger.Info("Initializing multi-repository mode", "repositories", len(config.Repositories))

		for _, repoUrl := range config.Repositories {
			ghConfig, err := actions.ParseGitHubConfigFromURL(repoUrl)
			if err != nil {
				return nil, fmt.Errorf("failed to parse GitHub config from URL %s: %w", repoUrl, err)
			}

			actionsClient, err := config.ActionsClientForURL(repoUrl, 0, app.logger)
			if err != nil {
				return nil, fmt.Errorf("failed to create actions client for %s: %w", repoUrl, err)
			}

			// Get or create scale set for this repository
			scaleSetName := config.RunnerScaleSetName
			if scaleSetName == "" {
				scaleSetName = config.EphemeralRunnerSetName
			}

			app.logger.Info("Getting or creating scale set", "repository", repoUrl, "scaleSetName", scaleSetName)

			scaleSet, err := actionsClient.GetRunnerScaleSet(context.Background(), 1, scaleSetName)
			if err != nil {
				// Try to create the scale set
				app.logger.Info("Scale set not found, attempting to create", "repository", repoUrl, "scaleSetName", scaleSetName)
				newScaleSet := &actions.RunnerScaleSet{
					Name:          scaleSetName,
					RunnerGroupId: 1,
					Labels:        []actions.Label{{Type: "System", Name: "self-hosted"}},
					RunnerSetting: actions.RunnerSetting{},
				}
				scaleSet, err = actionsClient.CreateRunnerScaleSet(context.Background(), newScaleSet)
				if err != nil {
					return nil, fmt.Errorf("failed to create scale set for %s: %w", repoUrl, err)
				}
			}

			app.logger.Info("Using scale set", "repository", repoUrl, "scaleSetId", scaleSet.Id, "scaleSetName", scaleSet.Name)

			// Create metrics exporter if configured
			var metricsPublisher metrics.Publisher
			if config.MetricsAddr != "" && app.metrics == nil {
				// Only create one metrics server for all listeners
				app.metrics = metrics.NewExporter(metrics.ExporterConfig{
					ScaleSetName:      config.EphemeralRunnerSetName,
					ScaleSetNamespace: config.EphemeralRunnerSetNamespace,
					Enterprise:        ghConfig.Enterprise,
					Organization:      ghConfig.Organization,
					Repository:        ghConfig.Repository,
					ServerAddr:        config.MetricsAddr,
					ServerEndpoint:    config.MetricsEndpoint,
					Metrics:           config.Metrics,
					Logger:            app.logger.WithName("metrics exporter"),
				})
				metricsPublisher = app.metrics
			}

			// Create listener for this repository
			listener, err := listener.New(listener.Config{
				Client:     actionsClient,
				ScaleSetID: scaleSet.Id,
				MinRunners: app.config.MinRunners,
				MaxRunners: app.config.MaxRunners,
				Logger:     app.logger.WithName("listener").WithValues("repository", repoUrl, "scaleSetId", scaleSet.Id),
				Metrics:    metricsPublisher,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create listener for %s: %w", repoUrl, err)
			}

			app.listeners = append(app.listeners, listener)
		}

		app.logger.Info("Multi-repository mode initialized", "listeners", len(app.listeners))
	} else {
		// Single repository mode (existing behavior)
		ghConfig, err := actions.ParseGitHubConfigFromURL(config.ConfigureUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse GitHub config from URL: %w", err)
		}

		actionsClient, err := config.ActionsClient(app.logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create actions client: %w", err)
		}

		if config.MetricsAddr != "" {
			app.metrics = metrics.NewExporter(metrics.ExporterConfig{
				ScaleSetName:      config.EphemeralRunnerSetName,
				ScaleSetNamespace: config.EphemeralRunnerSetNamespace,
				Enterprise:        ghConfig.Enterprise,
				Organization:      ghConfig.Organization,
				Repository:        ghConfig.Repository,
				ServerAddr:        config.MetricsAddr,
				ServerEndpoint:    config.MetricsEndpoint,
				Metrics:           config.Metrics,
				Logger:            app.logger.WithName("metrics exporter"),
			})
		}

		listener, err := listener.New(listener.Config{
			Client:     actionsClient,
			ScaleSetID: app.config.RunnerScaleSetId,
			MinRunners: app.config.MinRunners,
			MaxRunners: app.config.MaxRunners,
			Logger:     app.logger.WithName("listener"),
			Metrics:    app.metrics,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create new listener: %w", err)
		}
		app.listener = listener
	}

	app.logger.Info("app initialized")

	return app, nil
}

func (app *App) Run(ctx context.Context) error {
	var errs []error
	if app.worker == nil {
		errs = append(errs, fmt.Errorf("worker not initialized"))
	}
	if app.listener == nil && len(app.listeners) == 0 {
		errs = append(errs, fmt.Errorf("listener not initialized"))
	}
	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("app not initialized: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	metricsCtx, cancelMetrics := context.WithCancelCause(ctx)

	if len(app.listeners) > 0 {
		// Multi-repository mode: run all listeners
		for i, listener := range app.listeners {
			listenerIndex := i
			listenerInstance := listener
			g.Go(func() error {
				app.logger.Info("Starting multi-repo listener", "listenerIndex", listenerIndex)
				err := listenerInstance.Listen(ctx, app.worker)
				if err != nil {
					app.logger.Error(err, "Multi-repo listener exited with error", "listenerIndex", listenerIndex)
				}
				// Don't cancel metrics or other listeners if one fails
				// Just log and return the error
				return fmt.Errorf("listener %d failed: %w", listenerIndex, err)
			})
		}
	} else {
		// Single repository mode
		g.Go(func() error {
			app.logger.Info("Starting listener")
			listenerErr := app.listener.Listen(ctx, app.worker)
			cancelMetrics(fmt.Errorf("Listener exited: %w", listenerErr))
			return listenerErr
		})
	}

	if app.metrics != nil {
		g.Go(func() error {
			app.logger.Info("Starting metrics server")
			return app.metrics.ListenAndServe(metricsCtx)
		})
	}

	return g.Wait()
}
