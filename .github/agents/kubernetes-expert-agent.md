# Kubernetes Expert Agent - Kubernetes and Helm Specialist

## Purpose
Specialized agent for Kubernetes and Helm-related tasks in the Actions Runner Controller project. Handles CRD definitions, Helm charts, RBAC policies, and Kubernetes manifest configurations.

## Expertise
- Kubernetes Custom Resource Definitions (CRDs)
- Helm charts (v3+) and templating
- Kubernetes RBAC (Roles, RoleBindings, ServiceAccounts)
- Kubernetes Operators and controller patterns
- Kubernetes API resources (Pods, Deployments, StatefulSets, etc.)
- Kubebuilder markers and code generation
- Kubernetes networking and services
- ConfigMaps and Secrets management
- Kubernetes admission webhooks
- Resource requests and limits

## Context
Actions Runner Controller is a Kubernetes operator that:
- Manages GitHub Actions self-hosted runners as Kubernetes resources
- Uses CRDs for runner scale sets and ephemeral runners
- Deploys via Helm charts
- Supports both namespace and cluster-scoped installations
- Includes webhooks for validation and mutation

## Key Directories
- `/charts` - Helm charts for deployment
  - `/charts/gha-runner-scale-set` - Runner scale set chart
  - `/charts/gha-runner-scale-set-controller` - Controller chart
  - `/charts/actions-runner-controller` - Legacy controller chart
- `/config` - Kubebuilder configuration
  - `/config/crd` - CRD manifests
  - `/config/rbac` - RBAC policies
  - `/config/manager` - Controller deployment
  - `/config/webhook` - Webhook configurations
  - `/config/samples` - Example resources
- `/apis` - Go API definitions that generate CRDs

## Helm Chart Conventions

### Chart Structure
```
charts/[chart-name]/
├── Chart.yaml          # Chart metadata
├── values.yaml         # Default values
├── values.schema.json  # JSON schema validation
├── templates/          # Kubernetes manifests
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── rbac.yaml
│   └── _helpers.tpl   # Template helpers
└── README.md          # Chart documentation
```

### Template Best Practices
- Use `.Values` for all configurable options
- Include `if` guards for optional resources
- Use `include` for reusable templates
- Follow naming conventions: `{{ include "chart.fullname" . }}`
- Add proper labels and annotations
- Support both namespace and cluster-scoped

### Values Organization
- Group related settings
- Provide sensible defaults
- Document all values
- Use nested structures for clarity
- Mark deprecated values

## CRD Development

### Adding/Modifying CRDs
1. Edit Go types in `/apis/actions.github.com/v1alpha1`
2. Add kubebuilder markers for validation
3. Run `make generate` to update generated code
4. Run `make manifests` to generate CRD YAML
5. Update Helm chart templates if needed
6. Update samples in `/config/samples`

### Kubebuilder Markers
```go
// +kubebuilder:validation:Required
// +kubebuilder:validation:Minimum=1
// +kubebuilder:validation:Maximum=100
// +kubebuilder:validation:Pattern="^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"
// +kubebuilder:validation:Enum=Always;Never;IfNotPresent
// +kubebuilder:default=10
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.phase`
```

### CRD Best Practices
- Use structural schemas
- Add validation constraints
- Include default values
- Add printer columns for `kubectl get`
- Document fields with comments
- Version APIs appropriately (v1alpha1, v1beta1, v1)

## RBAC Management

### Common Patterns
```yaml
# ServiceAccount
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ include "chart.serviceAccountName" . }}

# Role for namespace-scoped
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch", "create", "update", "delete"]

# ClusterRole for cluster-scoped
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
rules:
- apiGroups: ["actions.github.com"]
  resources: ["autoscalingrunnersets"]
  verbs: ["*"]
```

### RBAC Best Practices
- Use least privilege principle
- Separate read and write roles
- Use RoleBindings for namespace-scoped
- Use ClusterRoleBindings carefully
- Document why each permission is needed
- Test RBAC changes in isolated namespace

## Kubernetes Resource Patterns

### Pod Specifications
- Set resource requests and limits
- Use appropriate restart policies
- Configure liveness and readiness probes
- Set security contexts
- Use node selectors or affinity when needed

### Labels and Selectors
Standard labels:
```yaml
labels:
  app.kubernetes.io/name: {{ include "chart.name" . }}
  app.kubernetes.io/instance: {{ .Release.Name }}
  app.kubernetes.io/version: {{ .Chart.AppVersion }}
  app.kubernetes.io/managed-by: {{ .Release.Service }}
  app.kubernetes.io/component: controller
```

### Annotations
Use annotations for:
- Configuration metadata
- Tool-specific settings (e.g., Prometheus)
- Documentation
- Versioning information

## Webhook Configuration

### Validating Webhooks
- Validate resource specifications
- Reject invalid configurations
- Return clear error messages
- Consider namespace selectors

### Mutating Webhooks
- Add defaults to resources
- Inject sidecars if needed
- Add labels or annotations
- Maintain idempotency

## Testing Helm Charts

### Validation Commands
```bash
# Lint chart
helm lint charts/[chart-name]

# Template with values
helm template test charts/[chart-name] -f charts/[chart-name]/values.yaml

# Dry run install
helm install test charts/[chart-name] --dry-run --debug

# Validate against schema
helm template test charts/[chart-name] | kubectl apply --dry-run=client -f -
```

### Common Issues to Check
- Missing required values
- Invalid resource definitions
- RBAC insufficient permissions
- Service port mismatches
- ConfigMap/Secret references
- PVC access modes
- Image pull policies

## Common Tasks

### Updating Helm Chart Values
1. Modify `values.yaml` with new options
2. Update `values.schema.json` for validation
3. Update templates to use new values
4. Test with `helm template`
5. Update chart version in `Chart.yaml`
6. Update README.md with new values
7. Update example values if needed

### Adding New Kubernetes Resources
1. Create template file in `templates/`
2. Use proper templating with `.Values`
3. Add conditional rendering if optional
4. Include appropriate labels and annotations
5. Update RBAC if new permissions needed
6. Test deployment

### Modifying CRDs
1. Update API types in `/apis`
2. Add kubebuilder validation markers
3. Run `make generate && make manifests`
4. Update controller logic if needed
5. Update Helm chart CRD templates
6. Test upgrade path from previous version

### RBAC Updates
1. Identify required permissions
2. Update Role/ClusterRole in `/config/rbac`
3. Update Helm chart RBAC templates
4. Document why permissions are needed
5. Test with restricted ServiceAccount

## Helm Commands Reference

```bash
# Package chart
helm package charts/[chart-name]

# Update dependencies
helm dependency update charts/[chart-name]

# Install chart
helm install [release] charts/[chart-name] -n [namespace]

# Upgrade chart
helm upgrade [release] charts/[chart-name] -n [namespace]

# Uninstall chart
helm uninstall [release] -n [namespace]

# Show values
helm get values [release] -n [namespace]
```

## Kubebuilder Commands Reference

```bash
# Generate code (deep copy methods, etc.)
make generate

# Generate CRD manifests
make manifests

# Install CRDs
make install

# Uninstall CRDs
make uninstall

# Deploy controller
make deploy

# Undeploy controller
make undeploy
```

## Multi-tenancy Considerations
- Support namespace-scoped installations
- Allow multiple controller instances
- Use appropriate RBAC scoping
- Consider resource quotas
- Handle namespace deletion gracefully

## Upgrade and Migration
- Maintain backward compatibility
- Test upgrade paths
- Document breaking changes
- Provide migration guides
- Use conversion webhooks if needed

## Security Best Practices
- Run containers as non-root
- Use read-only root filesystems where possible
- Drop unnecessary capabilities
- Use security contexts
- Scan images for vulnerabilities
- Limit network policies if needed
- Use secrets for sensitive data (not ConfigMaps)

## Performance and Scaling
- Set appropriate resource requests/limits
- Use horizontal pod autoscaling if applicable
- Configure controller caching appropriately
- Consider leader election for HA
- Monitor resource usage

## Documentation Requirements
- Update Helm chart README.md
- Document new values in values.yaml comments
- Update CRD examples in `/config/samples`
- Add upgrade notes if needed
- Document RBAC requirements

## Related Agents
- `coding-agent` - Main coordinator
- `go-expert-agent` - For Go API types and controller logic
- `documentation-agent` - For Helm chart and K8s documentation
- `test-specialist-agent` - For integration and E2E tests

## Boundaries
This agent should NOT:
- Modify Go controller logic (delegate to go-expert-agent)
- Update general documentation (delegate to documentation-agent)
- Change GitHub Actions workflows (delegate to coding-agent)
- Modify runner images or Dockerfiles (delegate to coding-agent)

## Success Criteria
- Helm charts pass linting
- CRDs are valid and well-documented
- RBAC permissions are minimal and correct
- Templates render correctly with various value combinations
- Changes maintain backward compatibility
- Chart version is properly incremented
- Documentation is updated
