# Multi-Repository Listener Support

## Overview

The GitHub Actions Listener now supports listening to multiple repositories from a single listener pod. This is useful for individual users or small teams with multiple repositories who want to share a single runner pool across all their repositories.

## Features

- **Multiple Repository Support**: Configure a single listener to monitor jobs from multiple GitHub repositories
- **Automatic Scale Set Management**: The listener automatically creates or retrieves scale sets for each repository
- **Shared Runner Pool**: All repositories share the same Kubernetes runner pool
- **Robust Error Handling**: Each repository listener runs independently with retry and reconnect logic
- **Minimal Configuration**: Simple JSON configuration file with list of repository URLs

## Configuration

### Prerequisites

- A Personal Access Token (PAT) or GitHub App with permissions across all target repositories
- Kubernetes cluster with actions-runner-controller installed
- Sufficient permissions to create runner scale sets in each repository

### Configuration File Format

Create a JSON configuration file with the following structure:

```json
{
  "repositories": [
    "https://github.com/myuser/repo1",
    "https://github.com/myuser/repo2",
    "https://github.com/myuser/repo3"
  ],
  "github_token": "ghp_yourPersonalAccessToken",
  "ephemeral_runner_set_namespace": "default",
  "ephemeral_runner_set_name": "multi-repo-runners",
  "runner_scale_set_name": "multi-repo-scale-set",
  "max_runners": 10,
  "min_runners": 1,
  "log_level": "info",
  "log_format": "text"
}
```

### Configuration Fields

- **repositories** (required): Array of GitHub repository URLs to listen to
- **github_token** (required): Personal Access Token with repo and workflow permissions
- **ephemeral_runner_set_namespace** (required): Kubernetes namespace for the runner pods
- **ephemeral_runner_set_name** (required): Name of the ephemeral runner set
- **runner_scale_set_name** (required): Base name for the scale sets (one will be created per repository)
- **max_runners** (optional): Maximum number of concurrent runners (default: 5)
- **min_runners** (optional): Minimum number of idle runners (default: 0)
- **log_level** (optional): Logging level - debug, info, warn, error (default: debug)
- **log_format** (optional): Log format - text or json (default: text)

### Alternative: GitHub App Configuration

You can also use a GitHub App instead of a PAT:

```json
{
  "repositories": [
    "https://github.com/myuser/repo1",
    "https://github.com/myuser/repo2"
  ],
  "github_app_id": "123456",
  "github_app_installation_id": "12345678",
  "github_app_private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
  "ephemeral_runner_set_namespace": "default",
  "ephemeral_runner_set_name": "multi-repo-runners",
  "runner_scale_set_name": "multi-repo-scale-set",
  "max_runners": 10,
  "min_runners": 1
}
```

## Usage

### Running the Listener

1. Create your configuration file (e.g., `multi-repo-config.json`)

2. Set the configuration path environment variable:
   ```bash
   export LISTENER_CONFIG_PATH=/path/to/multi-repo-config.json
   ```

3. Run the listener:
   ```bash
   ./bin/github-runnerscaleset-listener
   ```

### Kubernetes Deployment

Create a ConfigMap with your configuration:

```bash
kubectl create configmap listener-config \
  --from-file=config.json=multi-repo-config.json \
  --namespace=default
```

Update your listener deployment to mount the ConfigMap and set the environment variable:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: multi-repo-listener
spec:
  replicas: 1
  selector:
    matchLabels:
      app: multi-repo-listener
  template:
    metadata:
      labels:
        app: multi-repo-listener
    spec:
      containers:
      - name: listener
        image: ghcr.io/actions/actions-runner-controller:latest
        command: ["/github-runnerscaleset-listener"]
        env:
        - name: LISTENER_CONFIG_PATH
          value: /etc/listener/config.json
        volumeMounts:
        - name: config
          mountPath: /etc/listener
      volumes:
      - name: config
        configMap:
          name: listener-config
```

## How It Works

1. **Initialization**: On startup, the listener reads the configuration file and validates all repository URLs

2. **Scale Set Creation**: For each repository:
   - The listener connects to the GitHub API
   - Attempts to retrieve an existing scale set with the configured name
   - If not found, creates a new scale set for that repository

3. **Listener Setup**: A separate listener instance is created for each repository's scale set

4. **Concurrent Operation**: All listeners run concurrently, each maintaining its own:
   - Message session with GitHub Actions
   - Job queue monitoring
   - Retry and reconnect logic

5. **Runner Allocation**: When jobs are triggered in any repository:
   - The corresponding listener receives the job notification
   - Jobs are acquired and runners are scaled accordingly
   - All runners are created in the same Kubernetes namespace

## Retry and Reconnect Logic

The listener includes robust error handling:

- **Session Creation**: Up to 10 retries with 30-second intervals if session creation fails
- **Token Refresh**: Automatic token refresh when message queue tokens expire
- **Network Resilience**: Retryable HTTP client with exponential backoff for transient failures
- **Independent Listeners**: Each repository listener operates independently; if one fails, others continue

## Limitations

- All repositories must be accessible with the same authentication credentials (PAT or GitHub App)
- Scale sets are created with default runner group (ID: 1)
- All runners share the same resource limits and pod specifications
- Metrics are aggregated across all repositories

## Migration from Single Repository

To migrate from single repository mode:

1. Keep your existing configuration structure
2. Add the `repositories` field with your repository URLs
3. Remove the `configure_url` field (it will be ignored if `repositories` is present)
4. Restart the listener

The system will automatically create scale sets for each repository while maintaining backward compatibility with single-repository configurations.

## Troubleshooting

### Listener fails to start

- Verify all repository URLs are valid and accessible
- Ensure your PAT or GitHub App has permissions for all repositories
- Check Kubernetes namespace exists and has sufficient resources

### Scale set creation fails

- Verify authentication credentials have admin permissions
- Check GitHub API rate limits
- Ensure runner group ID (1) exists in your GitHub organization/repository

### Jobs not being picked up

- Verify scale sets were created successfully in each repository
- Check listener logs for connection errors
- Ensure runners have proper labels configured in workflows

### High resource usage

- Reduce `max_runners` if too many concurrent runners are created
- Adjust `min_runners` to 0 to avoid idle runners
- Monitor Kubernetes resource usage and adjust pod limits

## Example Logs

Successful multi-repository initialization:

```
INFO  listener-app  Initializing multi-repository mode  repositories=3
INFO  listener-app  Getting or creating scale set  repository=https://github.com/myuser/repo1 scaleSetName=multi-repo-scale-set
INFO  listener-app  Using scale set  repository=https://github.com/myuser/repo1 scaleSetId=12345 scaleSetName=multi-repo-scale-set
INFO  listener-app  Getting or creating scale set  repository=https://github.com/myuser/repo2 scaleSetName=multi-repo-scale-set
INFO  listener-app  Using scale set  repository=https://github.com/myuser/repo2 scaleSetId=12346 scaleSetName=multi-repo-scale-set
INFO  listener-app  Multi-repository mode initialized  listeners=3
INFO  listener-app  Starting multi-repo listener  listenerIndex=0
INFO  listener-app  Starting multi-repo listener  listenerIndex=1
INFO  listener-app  Starting multi-repo listener  listenerIndex=2
```
