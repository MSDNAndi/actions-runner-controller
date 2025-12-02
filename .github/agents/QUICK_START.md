# Agent System Quick Start Guide

This guide helps you quickly understand and use the GitHub Copilot coding agent system in the Actions Runner Controller repository.

## What Are These Agents?

These are specialized AI assistant configurations that help with different aspects of development in this project. Think of them as expert teammates who know specific domains really well.

## The Three Types of Agents

### 1. Main Coordinator
- **coding-agent**: Your starting point for any development task

### 2. Management Agents (Meta-Agents)
- **agent-manager**: Keeps the agent system healthy and up-to-date
- **agent-recruiter**: Creates new agents when needed

### 3. Specialist Agents (The Experts)
- **go-expert-agent**: Everything Go-related
- **kubernetes-expert-agent**: Kubernetes, Helm, CRDs, RBAC
- **documentation-agent**: All documentation
- **test-specialist-agent**: All testing

## How To Use This System

### For Quick Tasks

```
Question: "How do I add a new field to the AutoscalingRunnerSet CRD?"
→ Check: kubernetes-expert-agent.md
→ Section: "Adding/Modifying CRDs"
```

```
Question: "How do I write a controller test?"
→ Check: test-specialist-agent.md
→ Section: "Writing Controller Tests"
```

```
Question: "How do I handle errors in Go code?"
→ Check: go-expert-agent.md
→ Section: "Error Handling"
```

### For Complex Tasks

```
Task: "Add a new autoscaling metric"
Steps:
1. coding-agent coordinates the work
2. go-expert-agent implements the logic
3. kubernetes-expert-agent updates CRDs
4. test-specialist-agent adds tests
5. documentation-agent documents it
```

## Quick Reference Card

| I Need To... | Look Here |
|--------------|-----------|
| Add Go code | go-expert-agent.md |
| Modify a CRD | kubernetes-expert-agent.md |
| Update Helm chart | kubernetes-expert-agent.md |
| Write tests | test-specialist-agent.md |
| Update docs | documentation-agent.md |
| Create new agent | agent-recruiter.md |
| Evaluate agents | agent-manager.md |
| General guidance | coding-agent.md |

## Common Commands By Agent

### Go Expert
```bash
make test          # Run tests
make lint          # Run linters
make build         # Build binary
make generate      # Generate code
```

### Kubernetes Expert
```bash
make manifests     # Generate CRDs
helm lint charts/* # Lint charts
make install       # Install CRDs
make deploy        # Deploy controller
```

### Test Specialist
```bash
ginkgo -r          # Run all Ginkgo tests
go test -cover     # Coverage report
ginkgo -focus      # Run specific tests
```

## File Location Patterns

```
Need to modify...        → Use agent...
/controllers/*.go        → go-expert-agent
/apis/*.go              → go-expert-agent + kubernetes-expert-agent
/charts/*               → kubernetes-expert-agent
/config/crd/*           → kubernetes-expert-agent
/config/rbac/*          → kubernetes-expert-agent
*.md files              → documentation-agent
*_test.go files         → test-specialist-agent
/test_e2e_arc/*         → test-specialist-agent
```

## Agent Philosophy

### Key Principles

1. **Specialization**: Each agent is an expert in their domain
2. **Delegation**: Main agent delegates to specialists
3. **Minimal Overlap**: Clear boundaries between agents
4. **Living Documents**: Agents evolve with the project
5. **Self-Managing**: Agents manage their own ecosystem

### Why This Matters

- **Faster Onboarding**: New contributors find answers quickly
- **Consistent Quality**: Best practices are documented
- **Less Context Switching**: Focus on one domain at a time
- **Better Collaboration**: Clear handoff protocols
- **Reduced Errors**: Domain experts catch domain-specific issues

## Examples of Agent Usage

### Example 1: Bug Fix in Controller

```
Problem: Runner pods not getting created

Steps:
1. Read go-expert-agent for controller patterns
2. Read test-specialist-agent for test approach
3. Write failing test (test-specialist-agent guides)
4. Fix bug (go-expert-agent guides)
5. Update documentation (documentation-agent guides)
```

### Example 2: New Helm Chart Value

```
Task: Add new chart value for resource limits

Steps:
1. Read kubernetes-expert-agent section on Helm
2. Update values.yaml with new option
3. Update values.schema.json
4. Update templates to use new value
5. Document in chart README (documentation-agent)
6. Test with helm template
```

### Example 3: Adding New Feature

```
Feature: Add webhook for custom validation

Steps:
1. coding-agent reviews overall approach
2. go-expert-agent implements webhook logic
3. kubernetes-expert-agent updates webhook config
4. kubernetes-expert-agent adds RBAC rules
5. test-specialist-agent adds webhook tests
6. documentation-agent documents the feature
```

## Troubleshooting Agent Usage

### "I don't know which agent to use"
Start with coding-agent.md - it will point you to the right specialist.

### "Multiple agents seem relevant"
That's normal for cross-cutting tasks. Use coding-agent to coordinate between them.

### "Agent information is outdated"
Update the agent spec and note what changed. Consider informing agent-manager.

### "I need an agent that doesn't exist"
Review agent-manager.md to understand if a new agent is warranted, then check agent-recruiter.md for how to create one.

### "Agent boundaries are unclear"
Check the "Boundaries" section in each agent - it explicitly states what they don't do.

## Best Practices

### DO
- ✓ Read the relevant agent spec before starting work
- ✓ Follow the guidelines in the agent specs
- ✓ Update agent specs when patterns change
- ✓ Use agents as onboarding material
- ✓ Suggest improvements to agent specs

### DON'T
- ✗ Ignore agent boundaries
- ✗ Let agent specs become outdated
- ✗ Create overlapping agents
- ✗ Skip testing because the agent didn't mention it
- ✗ Follow agent advice blindly - use judgment

## Measuring Success

The agent system is working when:

- [ ] New contributors find answers in agent specs
- [ ] PRs follow conventions documented in agents
- [ ] Less time spent on code review for style/conventions
- [ ] Fewer bugs in specialist areas
- [ ] Faster feature development
- [ ] More consistent code quality

## Evolution

This agent system will evolve as the project evolves:

- **New technology?** → Consider new specialist agent
- **Pattern changes?** → Update relevant agent specs
- **Agent unused?** → Consider retirement
- **Agent confusion?** → Clarify or split
- **Missing coverage?** → Add new agent

## Getting Help

- **About this project**: Read coding-agent.md
- **Specific domain**: Read relevant specialist agent
- **Agent system itself**: Read this guide and README.md
- **Creating agents**: Read agent-recruiter.md
- **Managing agents**: Read agent-manager.md

## Summary

The agent system is your AI-powered development guide. Each agent is a specialist who helps you work effectively in their domain. Start with the coding-agent, delegate to specialists, and keep the specs updated as the project grows.

---

**Remember**: These agents are tools to help you succeed. Use them as guides, but always apply your own judgment and expertise.
