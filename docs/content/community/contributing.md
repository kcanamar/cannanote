---
title: "Contributing to CannaNote"
description: "Learn how to contribute to CannaNote's mission of privacy-first cannabis wellness tracking."
sidebar_label: "Contributing"
sidebar_order: 2
section: "community"
keywords: ["contributing", "development", "open source", "community"]
related_pages: ["community/content-workflow"]
---

# Contributing to CannaNote

Welcome to the CannaNote community! We're building the future of privacy-first cannabis wellness tracking, and we'd love your help.

## Our Mission

CannaNote exists to help people develop a healthier relationship with cannabis through:

- **Privacy by design** - Your data belongs to you
- **Evidence-based features** - Backed by research and user needs
- **Harm reduction focus** - Supporting mindful consumption
- **Radical transparency** - Open source and honest practices

## How to Contribute

### Code Contributions

#### Prerequisites

- **Go 1.21+** for backend development
- **Flutter/Dart** for mobile development
- **Git** for version control
- **Make** for development commands

#### Getting Started

1. **Fork the repository** on [GitHub](https://github.com/kcanamar/cannanote)
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/yourusername/cannanote.git
   cd cannanote
   ```
3. **Set up development environment**:
   ```bash
   make dev
   ```
4. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

#### Development Workflow

1. **Write tests first** - We practice test-driven development
2. **Follow existing patterns** - Check similar features for consistency
3. **Keep privacy first** - Never log or expose personal data
4. **Document your changes** - Update relevant docs and comments
5. **Test thoroughly** - Run the full test suite before submitting

#### Code Standards

- **Go**: Follow `gofmt` and `golint` conventions
- **Flutter**: Use `dart format` and follow Flutter best practices
- **Commit messages**: Use conventional commits format
- **Privacy**: Never commit secrets or personal data

#### Submitting Changes

1. **Push your branch** to your fork
2. **Open a pull request** against the main branch
3. **Describe your changes** clearly in the PR description
4. **Link any related issues** using GitHub keywords
5. **Wait for review** - We'll provide feedback within 48 hours

### Documentation Contributions

#### Content Guidelines

- **Clear and concise** - Help users accomplish their goals quickly
- **Privacy-focused** - Emphasize user control and data ownership
- **Evidence-based** - Include research links where relevant
- **Inclusive language** - Welcoming to all cannabis users

#### How to Add Documentation

1. **Follow the [Content Workflow](content-workflow)** guide
2. **Use proper frontmatter** with all required fields
3. **Test locally** before submitting
4. **Create related page links** to improve navigation

### Feature Requests

#### Before Submitting

- **Search existing issues** to avoid duplicates
- **Consider privacy implications** of your request
- **Think about harm reduction** impact
- **Provide use case examples** from real experience

#### Creating Good Feature Requests

1. **Use the issue template** on GitHub
2. **Describe the problem** you're trying to solve
3. **Explain your proposed solution** with specific details
4. **Consider alternative approaches** and their tradeoffs
5. **Estimate impact** on user privacy and data security

### Bug Reports

#### Information to Include

- **Steps to reproduce** the issue
- **Expected behavior** vs actual behavior
- **Platform details** (web, iOS, Android)
- **App version** and device information
- **Screenshots or screen recordings** if helpful

#### Privacy in Bug Reports

- **Never include** personal consumption data
- **Redact any** personal information from screenshots
- **Use example data** instead of your actual entries
- **Report security issues** privately to security@cannanote.com

## Community Guidelines

### Our Values

- **Respect** - Treat all community members with kindness
- **Privacy** - Respect others' personal information and choices
- **Evidence** - Base discussions on research and facts
- **Harm Reduction** - Support safer cannabis use practices
- **Inclusivity** - Welcome all responsible cannabis users

### Code of Conduct

We follow a zero-tolerance policy for:
- Harassment or discrimination
- Sharing others' personal information
- Promoting unsafe cannabis practices
- Commercial spam or self-promotion

## Recognition

Contributors are recognized in several ways:

### Hall of Fame

Outstanding contributors are featured in our documentation and app credits.

### Community Roles

Active contributors may be invited to become:
- **Maintainers** - Help review and merge contributions
- **Moderators** - Assist with community management
- **Advisors** - Provide guidance on product direction

### Beta Access

All contributors receive early access to new features and beta releases.

## Development Resources

### Architecture Overview

CannaNote uses hexagonal architecture with these key components:

- **Domain Layer** - Core business logic (Go)
- **Application Layer** - Use cases and services (Go)
- **Adapters** - External interfaces (HTTP, database, etc.)
- **Mobile App** - Flutter with Drift for local storage
- **Web Frontend** - Server-side rendering with HTMX

### Key Principles

1. **Privacy by Design** - Encrypt sensitive data client-side
2. **Offline First** - Full functionality without network
3. **Performance** - Sub-30-second logging experience
4. **Accessibility** - Support screen readers and keyboard navigation

### Testing Strategy

- **Unit tests** for domain logic (85%+ coverage required)
- **Integration tests** for API endpoints
- **Widget tests** for Flutter UI components
- **End-to-end tests** for critical user journeys

## Getting Help

### Documentation

- **[Content Workflow](content-workflow)** - Adding/editing documentation
- **README.md** - Quick start development guide
- **Architecture docs** - Technical implementation details

### Community Support

- **GitHub Discussions** - Ask questions and share ideas
- **Issue tracker** - Report bugs and request features
- **Email** - Contact maintainers at hello@cannanote.com

### Office Hours

Core maintainers host virtual office hours:
- **When**: Every other Wednesday at 7 PM EST
- **Where**: Announced in GitHub Discussions
- **Format**: Casual Q&A and pair programming

## License and Legal

### Open Source License

CannaNote is released under the [MIT License](https://opensource.org/licenses/MIT), which means:

- ✅ Use for any purpose (personal, commercial, etc.)
- ✅ Modify and distribute freely
- ✅ Include in proprietary software
- ❗ Must include original license notice
- ❗ No warranty provided

### Contributor License Agreement

By contributing, you agree that:

1. **You have the right** to license your contribution
2. **Your contribution** will be licensed under MIT
3. **You grant permission** for your code to be used in CannaNote
4. **You understand** this is a cannabis-related project

### Privacy Commitment

All contributors must commit to:

- Never accessing or storing user consumption data
- Reporting security vulnerabilities responsibly
- Following privacy-by-design principles in all code
- Respecting user autonomy and data ownership

---

Thank you for contributing to CannaNote! Together, we're building a better future for cannabis wellness tracking.

**Questions?** Reach out to us at hello@cannanote.com or start a discussion on GitHub.