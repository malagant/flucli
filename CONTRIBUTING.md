# Contributing to FluxCLI

Thank you for your interest in contributing to FluxCLI! This document provides guidelines and instructions for contributing to the project.

## 🤝 Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or later
- Kubernetes cluster with FluxCD v2.0+
- kubectl configured
- Git

### Development Environment

We recommend using Nix for a reproducible development environment:

```bash
# Clone the repository
git clone https://github.com/malagant/fluxcli.git
cd fluxcli

# Enter development shell (provides all tools)
./dev.sh shell

# Or build directly
./dev.sh build
```

Alternative without Nix:

```bash
# Ensure you have Go 1.21+ installed
go version

# Build the project
make build

# Run tests
make test
```

## 🛠️ Development Workflow

### 1. Fork and Clone

```bash
# Fork the repository on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/fluxcli.git
cd fluxcli

# Add the original repository as upstream
git remote add upstream https://github.com/malagant/fluxcli.git
```

### 2. Create a Feature Branch

```bash
# Create and switch to a new branch
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/issue-description
```

### 3. Make Your Changes

- Follow the existing code style and conventions
- Add tests for new functionality
- Update documentation as needed
- Ensure your changes don't break existing functionality

### 4. Test Your Changes

```bash
# Run all tests
make test

# Run linter
make lint

# Build to ensure no compilation errors
make build

# Test the binary manually
./fluxcli --help
```

### 5. Commit Your Changes

Follow conventional commit format:

```bash
# Examples of good commit messages
git commit -m "feat: add search functionality to resource view"
git commit -m "fix: resolve cluster switching bug"
git commit -m "docs: update installation instructions"
git commit -m "test: add unit tests for event view"
```

### 6. Push and Create Pull Request

```bash
# Push to your fork
git push origin feature/your-feature-name

# Create a pull request on GitHub
```

## 📝 Coding Standards

### Go Style Guide

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` to format code
- Follow Go naming conventions
- Write clear, descriptive variable and function names
- Add comments for exported functions and types

### Code Organization

```
pkg/
├── core/       # Core business logic, no UI dependencies
├── k8s/        # Kubernetes client and resource management
└── ui/         # Terminal UI components using Bubble Tea

internal/
└── config/     # Configuration management (internal use only)

cmd/
└── root.go     # CLI command definitions
```

### Testing

- Write unit tests for new functionality
- Use table-driven tests where appropriate
- Mock external dependencies (Kubernetes API, etc.)
- Aim for meaningful test coverage, not just high percentages

Example test structure:

```go
func TestResourceManager_ListResources(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected []Resource
        wantErr  bool
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 🎨 UI Guidelines

### Bubble Tea Components

- Keep components focused and single-purpose
- Implement proper keyboard navigation
- Handle window resizing gracefully
- Use consistent styling with lipgloss

### Color and Styling

- Use the existing color scheme consistently
- Ensure accessibility (readable color contrasts)
- Test in both light and dark terminal themes
- Use semantic colors (red for errors, green for success, etc.)

### Keyboard Shortcuts

- Follow K9s conventions where applicable
- Use vim-like navigation (j/k, g/G, etc.)
- Ensure shortcuts are discoverable (help screen)
- Avoid conflicts with standard terminal shortcuts

## 📚 Documentation

### Code Documentation

- Add godoc comments for all exported functions and types
- Include examples in documentation where helpful
- Keep comments up-to-date with code changes

### User Documentation

When adding new features:

- Update the user guide (`docs/user-guide.md`)
- Add keyboard shortcuts to help screens
- Update configuration examples if needed
- Consider adding screenshots or demos

## 🐛 Bug Reports

When reporting bugs, please include:

- FluxCLI version (`fluxcli --version`)
- Operating system and terminal
- Kubernetes version and FluxCD version
- Steps to reproduce the issue
- Expected vs actual behavior
- Relevant log output (use `--debug` flag)

## 💡 Feature Requests

For new features:

- Check if the feature aligns with the project goals
- Consider if it fits the MVP scope or should be a future enhancement
- Provide clear use cases and benefits
- Consider implementation complexity
- Be open to feedback and alternative solutions

## 🔄 Pull Request Process

### Before Submitting

- [ ] Tests pass (`make test`)
- [ ] Linting passes (`make lint`)
- [ ] Code builds successfully (`make build`)
- [ ] Documentation updated if needed
- [ ] Commit messages follow conventional format
- [ ] Changes are focused and atomic

### Pull Request Template

Your PR should include:

1. **Description**: Clear description of changes
2. **Type**: Bug fix, new feature, documentation, etc.
3. **Testing**: How you tested the changes
4. **Screenshots**: For UI changes (if applicable)
5. **Breaking Changes**: Any breaking changes and migration notes

### Review Process

1. Automated checks must pass (tests, linting)
2. At least one maintainer review required
3. Address review feedback promptly
4. Squash commits before merging (if requested)

## 🏗️ Architecture Guidelines

### Adding New Resource Types

When adding support for new FluxCD resources:

1. Update the Kubernetes client (`pkg/k8s/`)
2. Add resource type to the core manager (`pkg/core/`)
3. Update UI components to display the new resource type
4. Add tests for all layers
5. Update documentation

### Adding New Views

For new UI views:

1. Create a new component in `pkg/ui/`
2. Implement the `tea.Model` interface
3. Add keyboard navigation and help text
4. Integrate with the main app model
5. Add tests for the component

### Configuration Changes

When modifying configuration:

1. Update the config struct (`internal/config/`)
2. Handle backward compatibility
3. Update default configuration
4. Add validation for new fields
5. Update documentation

## 🚀 Release Process

Releases are handled by maintainers, but contributors should:

- Keep CHANGELOG.md updated (if exists)
- Tag significant changes in commit messages
- Consider backward compatibility
- Update version numbers in documentation

## ❓ Getting Help

- Check existing [issues](https://github.com/malagant/fluxcli/issues)
- Read the [documentation](docs/)
- Join discussions in GitHub Discussions
- Ask questions in issues with the "question" label

## 🙏 Recognition

Contributors will be recognized in:

- Release notes for significant contributions
- README contributors section
- GitHub contributor graphs

Thank you for contributing to FluxCLI! Your efforts help make GitOps management more accessible and enjoyable for everyone. 🌟
