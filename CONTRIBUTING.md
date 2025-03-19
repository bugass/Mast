# Contributing to Mast

Thank you for your interest in contributing to Mast! This document provides guidelines and steps for contributing.

## Development Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/your-username/mast.git
   cd mast
   ```
3. Create a new branch for your feature:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Code Style

- Follow the standard Go code style
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and small
- Write tests for new features

## Testing

Before submitting a pull request:
1. Run the tests:
   ```bash
   go test ./...
   ```
2. Ensure all tests pass
3. Add tests for new features

## Pull Request Process

1. Update the README.md if needed
2. Add tests for new features
3. Update documentation
4. Create a pull request with a clear description
5. Wait for review and address any feedback

## Feature Requests

When suggesting new features:
1. Check if the feature already exists
2. Provide a clear use case
3. Explain the benefits
4. Consider implementation complexity

## Bug Reports

When reporting bugs:
1. Use the issue template
2. Provide steps to reproduce
3. Include system information
4. Share relevant logs

## Code Review

- Review code for:
  - Correctness
  - Performance
  - Security
  - Maintainability
  - Test coverage

## Questions?

Feel free to open an issue for any questions or concerns. 