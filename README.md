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
- **Resend** - Transactional email delivery

**Database:**
- **NeonDB** - Managed PostgreSQL (serverless)

**Frontend:**
- **Server-side rendering** - Fast loading, accessible by default
- **Tailwind CSS v4** - Utility-first styling with local build
- **Vanilla JavaScript** - Zero framework overhead, minimal bundle size

**Infrastructure:**
- **Fly.io** - Simple, reliable deployment
- **Docker** - Consistent containerized deployments

## Project Structure

```
cannanote/
├── README.md                # Project overview and setup
├── CLAUDE.md                # Development roadmap and future plans
├── backend/                 # Go application
│   ├── cmd/
│   │   ├── api/             # API server entry point
│   │   └── web/             # Web templates and assets
│   │       ├── assets/      # Static assets (CSS, JS, images)
│   │       ├── components/  # Reusable templ components
│   │       └── *.templ      # Page templates
│   ├── internal/            # Private application code
│   │   ├── adapters/        # External service integrations
│   │   │   ├── external/    # Third-party services (email, etc.)
│   │   │   ├── http/        # HTTP handlers and middleware
│   │   │   └── repository/  # Data access layer
│   │   ├── core/            # Business logic
│   │   │   ├── domain/      # Entities and business rules
│   │   │   ├── application/ # Use cases and services
│   │   │   └── ports/       # Interface definitions
│   │   ├── database/        # Database connection setup
│   │   └── server/          # Server configuration and routing
│   ├── neondb/              # Database migrations
│   │   └── migrations/      # SQL migration files
│   ├── Dockerfile           # Container definition
│   ├── fly.toml             # Deployment configuration
│   ├── Makefile             # Build and development commands
│   └── go.mod               # Go dependencies
├── mobile/                  # Flutter mobile application (future)
└── docs/                    # Documentation
    ├── content/             # Markdown documentation content
    ├── style-guide.md       # Brand guidelines and design system
    └── engineering.md       # Development guidelines
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

## Current Features

### Beta Access System
- **Email signup flow** - Request beta access with email verification
- **Account activation** - Secure email verification and password setup
- **Session management** - Persistent authentication with secure cookies
- **Beta grandfathering** - Early supporters receive lifetime premium benefits

### Application Interface
- **Landing page** - Clear value proposition and beta signup
- **Authenticated dashboard** - Protected app area for verified users
- **Collapsible sidebar** - Slack/Discord-style navigation (desktop + mobile)
- **Dark mode support** - System-aware theme with manual toggle
- **Mobile-first design** - Responsive layouts optimized for touch

### Documentation System
- **Markdown-based docs** - Easy to maintain documentation
- **Searchable content** - Find information quickly
- **Cannabis education** - Harm reduction and science-based content

### Privacy Protection
- **Local-first architecture** - Cannabis data stays on your device
- **No third-party tracking** - Zero analytics or advertising SDKs
- **Minimal data collection** - Only what's necessary for core functionality
- **User data ownership** - Full control over your information

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
# Database (NeonDB)
DB_HOST=your-neondb-host
DB_PORT=5432
DB_DATABASE=your-database
DB_USERNAME=your-username
DB_PASSWORD=your-password
DB_SCHEMA=public

# Application
PORT=8080
APP_ENV=production
APP_URL=https://your-domain.com

# Email (Resend)
RESEND_API_KEY=your-resend-api-key
```

## Privacy & Security

### Data Protection
- **Local-first architecture** - Cannabis data stays on user devices
- **Minimal server data** - Only authentication and sync metadata on server
- **No third-party tracking** - Zero analytics or advertising SDKs
- **Transparent practices** - Clear documentation of all data handling

### Security Practices
- **Input validation** - All user inputs sanitized
- **HTTPS enforcement** - Encrypted connections required
- **Regular dependency updates** - Automated security scanning
- **Rate limiting** - Protection against abuse

## Roadmap

See [CLAUDE.md](CLAUDE.md) for the detailed development roadmap.

### Current Focus
- **Offline-first PWA** - Full functionality without network connection
- **Session logging** - 30-second cannabis session tracking
- **Local data storage** - IndexedDB-based client-side persistence

### Future Plans
- **Optional cloud sync** - Paid feature for cross-device access
- **Mobile application** - Native apps via Flutter
- **Pattern insights** - Personal analytics and recommendations

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