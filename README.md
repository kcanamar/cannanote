# CannaNote - Your Personal Cannabis Educator and Tracking Companion

> **Your patterns, not their profits.** A personal cannabis wellness journaling app that transforms consumption data into clear, actionable insights through seamless, private, and mindful tracking.

## Mission

CannaNote empowers individuals to understand their unique relationship with cannabis by transforming personal consumption data into clear, actionable insights—delivered through a seamless, private, and mindful journaling experience.

## Core Values

- **Radical Data Transparency** - Users own and fully understand their data. We never hide, sell, or share without explicit consent.
- **User Safety & Harm Reduction First** - Every feature prioritizes well-being over engagement.
- **Environmental Responsibility** - We champion regenerative cannabis practices and sustainable technology.
- **Evidence-Based Insights** - No hype, no exaggeration. Only what science and personal data support.
- **Inclusive & Equitable Design** - The app works for everyone, regardless of background or ability.

## Why CannaNote Exists

Cannabis users deserve better than expensive trial and error. Built from real industry experience and consumer frustration, CannaNote bridges the gap between cannabis complexity and personal understanding. We provide the tools to track what works, learn the science, and make confident decisions based on your patterns—not marketing hype.

## Technology Stack

**Backend:**
- **Go 1.24+** - Fast, reliable, simple deployment
- **Gin Router** - HTTP routing and middleware
- **templ** - Type-safe HTML templates
- **HTMX** - Dynamic interactions without complex JavaScript

**Database:**
- **Supabase PostgreSQL** - Database with built-in authentication
- **Row Level Security** - Privacy protection at the database level

**Frontend:**
- **Server-side rendering** - Fast loading, accessible by default
- **Tailwind CSS** - Utility-first styling
- **Alpine.js** (minimal) - Lightweight JavaScript when needed

**Infrastructure:**
- **Fly.io** - Simple, reliable deployment
- **Docker** - Consistent containerized deployments

## Project Structure

```
cannanote/
├── README.md                # Project overview and setup
├── backend/                 # Go application
│   ├── cmd/
│   │   ├── api/             # API server entry point  
│   │   └── web/             # Web templates and assets
│   │       ├── assets/      # Static assets (CSS, JS, images)
│   │       │   └── images/logos/  # Logo assets
│   │       ├── *.templ      # HTML templates
│   │       └── *.go         # Generated template code
│   ├── internal/            # Private application code
│   │   ├── adapters/        # External service integrations
│   │   │   ├── http/        # HTTP handlers
│   │   │   └── repository/  # Data access layer
│   │   ├── core/            # Business logic
│   │   │   ├── domain/      # Entities and business rules
│   │   │   ├── application/ # Use cases and services
│   │   │   └── ports/       # Interface definitions
│   │   ├── database/        # Database connection setup
│   │   └── server/          # Server configuration and routing
│   ├── tests/               # Test files
│   ├── Dockerfile           # Container definition
│   ├── fly.toml            # Deployment configuration
│   ├── Makefile            # Build and development commands
│   └── go.mod              # Go dependencies
├── mobile/                  # Flutter mobile application
├── supabase/               # Database schema and configuration
│   ├── config.toml         # Supabase configuration
│   ├── seed.sql            # Initial data
│   └── reference-data/     # Cannabinoids and terpenes data
└── docs/                   # Documentation
    ├── style-guide.md      # Brand guidelines and design system
    ├── brand-strategy.md   # Business strategy and positioning
    ├── engineering.md      # Development guidelines
    └── archive/            # Archived documentation
```

## Architecture Philosophy

### Simple, Maintainable, Fast

We prioritize:
1. **Developer velocity** - Fast to develop and deploy
2. **User experience** - Quick loading, responsive interactions  
3. **Maintenance** - Easy to understand and modify

### Hexagonal Architecture

- **Core Domain** (`internal/core/`) - Pure business logic, no external dependencies
- **Ports** (`internal/core/ports/`) - Interfaces defining external needs
- **Adapters** (`internal/adapters/`) - Concrete implementations for external services

This pattern keeps business logic separate from infrastructure, making the code easier to test, modify, and maintain.

## Getting Started

### Prerequisites

- **Go 1.24+** - [Install Go](https://golang.org/doc/install)
- **Git** - Version control
- **Docker** (optional) - For containerized development
- **Make** (optional) - For build commands

### Quick Start

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/cannanote.git
   cd cannanote/backend
   ```

2. **Set up environment:**
   ```bash
   # Copy example environment file
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Run the application:**
   ```bash
   # Install dependencies and start development server
   make dev
   ```

4. **Access the application:**
   - Local app: http://localhost:3001
   - Health check: http://localhost:3001/health

### Development Commands

```bash
# Start development server with hot reload
make dev

# Run tests
make test

# Run linting
make lint

# Build for production
make build

# Deploy to production
make deploy
```

## Features

### Personal Cannabis Tracking
- **Seamless logging** - Track sessions in under 30 seconds
- **Pattern recognition** - See what actually works for your body
- **Multiple consumption methods** - Flower, concentrates, edibles, topicals
- **Comprehensive data** - Strains, effects, dosage, timing, environment

### Evidence-Based Education
- **Cannabis science** explained in accessible language
- **Harm reduction** guidance and safety information
- **Myth-busting** content based on peer-reviewed research
- **Personal insights** derived from your own consumption patterns

### Privacy Protection
- **Local-first data storage** - Your data stays on your device
- **No data selling** - We never monetize your personal information
- **Radical transparency** - Clear documentation of all data practices
- **User ownership** - You control your data completely

### Mindful Design
- **Harm reduction first** - Features prioritize your well-being
- **Non-judgmental approach** - "You're not wrong" philosophy
- **Accessible interface** - Works for everyone, regardless of experience level
- **Environmental consciousness** - Sustainable technology choices

## Contributing

CannaNote is built with the cannabis community in mind. We welcome contributions that align with our values of privacy, harm reduction, and evidence-based education.

### Development Guidelines

1. **Preserve architectural patterns** - Follow existing hexagonal structure
2. **Maintain test coverage** - Write tests for new functionality
3. **Follow code standards** - Use `gofmt` and follow Go conventions
4. **Document changes** - Update documentation for significant modifications
5. **Prioritize privacy** - Never compromise user data protection

### Code of Conduct

- **Respectful communication** - Professional, inclusive interactions
- **Evidence-based discussions** - Support claims with research when possible
- **Harm reduction focus** - Features should prioritize user safety and well-being
- **Privacy conscious** - Consider data protection in all decisions

## Deployment

### Production Deployment

The application is designed for simple deployment to modern cloud platforms:

```bash
# Deploy to Fly.io (configured)
make deploy
```

This will:
1. Run tests and linting
2. Build Docker image
3. Deploy to production
4. Run health checks

### Environment Variables

Required configuration in `.env`:

```bash
# Database
DB_HOST=your-supabase-host
DB_PORT=5432
DB_DATABASE=postgres
DB_USERNAME=postgres
DB_PASSWORD=your-password

# Application
PORT=3001
APP_ENV=production
```

## Privacy & Security

### Data Protection
- **Local-first architecture** - Data stays on user devices when possible
- **Minimal data collection** - Only collect what's necessary for functionality
- **No third-party tracking** - No analytics or advertising SDKs
- **Transparent practices** - Clear documentation of all data handling

### Security Measures
- **Parameterized queries** - Prevent SQL injection
- **Input validation** - Sanitize all user inputs
- **HTTPS enforcement** - Encrypted connections in production
- **Regular updates** - Keep dependencies current for security

## Roadmap

### Current Focus
- **Core tracking functionality** - Reliable, fast journaling experience
- **Basic pattern recognition** - Help users understand their data
- **Mobile application** - Flutter app for iOS and Android

### Future Considerations
- **Advanced analytics** - Deeper insights from consumption patterns
- **Community features** - Optional sharing with privacy controls
- **Integration options** - Connect with other health tracking apps
- **Offline capabilities** - Full functionality without internet

## Support

### Documentation
- **Engineering Guide** - See `docs/engineering.md` for development details
- **Style Guide** - See `docs/style-guide.md` for brand and design guidelines
- **API Documentation** - Generated from code comments

### Getting Help
- **Issues** - Report bugs or request features via GitHub Issues
- **Discussions** - General questions and community discussions
- **Security** - Email security@cannanote.app for security concerns

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Individual business logic components may be licensed separately for commercial use.

---

**CannaNote** - Transforming personal cannabis data into clear, actionable insights through seamless, private, and mindful journaling. Because your patterns matter more than their profits.