# GitHub Copilot Coding Agents Configuration

This directory contains agent specifications for GitHub Copilot coding agents working in the Actions Runner Controller repository. These specifications enable AI-powered assistance for various development tasks through specialized agents.

## Overview

The agent system uses a hierarchical structure with a main coding agent that delegates to specialized sub-agents based on the task at hand. This approach ensures that domain-specific expertise is applied to each problem.

## Agent Architecture

```
coding-agent (Main Coordinator)
├── agent-manager (Lifecycle Management)
│   └── agent-recruiter (Specification Creator)
├── go-expert-agent (Go Development)
├── kubernetes-expert-agent (K8s/Helm)
├── documentation-agent (Documentation)
└── test-specialist-agent (Testing)
```

## Available Agents

### Core Agents

#### `coding-agent.md`
Main coordinating agent for general coding tasks. Delegates to specialized agents based on the task domain.

**When to Use**: General development tasks, coordinating multi-domain changes, overall project guidance.

### Management Agents

#### `agent-manager.md`
Manages the agent lifecycle - identifies when agents need to be created, updated, or retired.

**When to Use**: Analyzing agent effectiveness, identifying gaps in coverage, deciding on agent updates.

#### `agent-recruiter.md`
Creates new agent specifications when new domains or patterns emerge.

**When to Use**: Need for a new specialized agent, creating agent documentation.

### Specialist Agents

#### `go-expert-agent.md`
Go language specialist for all Go code changes, refactoring, and Go-specific tooling.

**When to Use**: 
- Modifying Go source code
- Adding new features in Go
- Refactoring Go code
- Fixing Go-related bugs
- Managing Go dependencies

**Expertise**: Go 1.25+, controller-runtime, Ginkgo/Gomega, error handling, concurrency.

#### `kubernetes-expert-agent.md`
Kubernetes and Helm specialist for CRDs, manifests, charts, and K8s configurations.

**When to Use**:
- Modifying CRDs or API types
- Updating Helm charts
- Changing RBAC policies
- Kubernetes resource configurations
- Webhook configurations

**Expertise**: CRDs, Helm charts, RBAC, Kubebuilder, K8s resources, admission webhooks.

#### `documentation-agent.md`
Documentation specialist for markdown files, README updates, and user guides.

**When to Use**:
- Updating documentation
- Creating new guides
- Maintaining README files
- Writing release notes
- Helm chart documentation

**Expertise**: Technical writing, markdown, API docs, user guides, troubleshooting docs.

#### `test-specialist-agent.md`
Testing expert for unit tests, integration tests, E2E tests, and test infrastructure.

**When to Use**:
- Writing new tests
- Fixing broken tests
- Improving test coverage
- Creating test fixtures
- Test infrastructure changes

**Expertise**: Ginkgo/Gomega, Go testing, envtest, E2E testing, mocking, test patterns.

## How to Use Agents

### For Copilot Coding Agent
When working on a task:

1. **Identify the Domain**: Determine which domain(s) the task touches
2. **Select Specialist**: Choose the most appropriate specialized agent
3. **Delegate**: Pass the task to the specialist with full context
4. **Coordinate**: If multiple domains are involved, coordinate between specialists
5. **Validate**: Ensure all changes work together correctly

### For Human Developers
These agent specifications also serve as:

- **Onboarding Material**: Understanding project conventions and structure
- **Quick Reference**: Finding commands, patterns, and best practices
- **Decision Guide**: Knowing which approach to use for common tasks
- **Quality Standards**: Understanding expectations for different types of work

## Agent Communication Protocol

### Delegation Pattern
```
Main Agent → "I need to modify CRD validation rules"
          → Delegates to kubernetes-expert-agent with full context
          → Kubernetes agent updates CRDs and manifests
          → Returns to main agent for integration
```

### Multi-Agent Coordination
```
Main Agent → "Add new autoscaling feature"
          → go-expert-agent: Implement controller logic
          → kubernetes-expert-agent: Update CRDs and Helm chart
          → test-specialist-agent: Add unit and integration tests
          → documentation-agent: Document the feature
          → Main agent validates integration
```

## Agent Lifecycle

### Creating New Agents

When a new specialized domain emerges:

1. **agent-manager** identifies the need
2. **agent-recruiter** creates the specification
3. Specification is reviewed and approved
4. Agent is added to the `.github/agents` directory
5. This README is updated
6. Other agents are updated to reference the new agent

### Updating Agents

When agent specifications need changes:

1. **agent-manager** identifies outdated content
2. Specification is updated with new information
3. Version and update date are recorded
4. Dependent agents are updated if needed
5. Changes are documented

### Retiring Agents

When agents are no longer needed:

1. **agent-manager** confirms obsolescence
2. Agent is moved to `.github/agents/archived/`
3. Retirement reason is documented
4. References are removed from other agents
5. This README is updated

## Best Practices

### For Agent Users (Copilot)
- Always read the full agent specification before starting work
- Pass complete context when delegating
- Respect agent boundaries (don't ask them to work outside their domain)
- Coordinate between agents for cross-cutting concerns
- Validate that delegated work integrates correctly

### For Agent Maintainers (Developers)
- Keep specifications current with project evolution
- Document changes with dates and reasons
- Avoid overlap between agent domains
- Include concrete examples
- Reference actual code patterns from the project
- Update related agents when making changes

## Quick Reference

| Task | Primary Agent | Supporting Agents |
|------|---------------|-------------------|
| Add Go controller logic | go-expert-agent | test-specialist-agent |
| Update CRD schema | kubernetes-expert-agent | go-expert-agent |
| Modify Helm chart | kubernetes-expert-agent | documentation-agent |
| Write new tests | test-specialist-agent | go-expert-agent |
| Update documentation | documentation-agent | - |
| Add new feature | coding-agent | All specialists as needed |
| Fix bug | coding-agent | Relevant specialist |
| Create new agent | agent-recruiter | agent-manager |
| Update agent spec | agent-manager | agent-recruiter |

## Directory Structure

```
.github/agents/
├── README.md                      # This file
├── coding-agent.md               # Main coordinator
├── agent-manager.md              # Lifecycle management
├── agent-recruiter.md            # Specification creator
├── go-expert-agent.md            # Go specialist
├── kubernetes-expert-agent.md    # K8s/Helm specialist
├── documentation-agent.md        # Documentation specialist
├── test-specialist-agent.md      # Testing specialist
└── archived/                     # Retired agents (when needed)
```

## Getting Started

### For New Contributors
1. Read `coding-agent.md` for project overview
2. Review specialist agent that matches your interest area
3. Follow the guidelines and conventions in the agent specifications

### For Existing Contributors
1. Use agent specifications as quick reference
2. Suggest updates when you find outdated information
3. Propose new agents when new patterns emerge

## Feedback and Improvements

Agent specifications are living documents that should evolve with the project:

- **Found outdated info?** Update the agent spec or file an issue
- **Missing important context?** Add it to the appropriate agent
- **Need a new agent?** Discuss with maintainers
- **Agent overlap?** Suggest consolidation or boundary clarification

## Version History

- **v1.0** (2025-12-02): Initial agent system with 7 agents
  - Core: coding-agent
  - Management: agent-manager, agent-recruiter
  - Specialists: go-expert, kubernetes-expert, documentation, test-specialist

## Resources

### Project Documentation
- [Contributing Guide](../../CONTRIBUTING.md)
- [Main README](../../README.md)
- [Troubleshooting](../../TROUBLESHOOTING.md)

### External Resources
- [GitHub Copilot Agent Best Practices](https://gh.io/copilot-coding-agent-tips)
- [Kubebuilder Documentation](https://book.kubebuilder.io/)
- [Ginkgo Testing Framework](https://onsi.github.io/ginkgo/)
- [Helm Documentation](https://helm.sh/docs/)

## License

These agent specifications are part of the Actions Runner Controller project and follow the same Apache 2.0 license.
