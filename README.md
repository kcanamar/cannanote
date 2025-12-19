# CannaNote - Cannabis Experience Journaling Platform

> A lifelong opus project: Building the first HIPAA-compliant cannabis tracking platform with hexagonal architecture for maximum maintainability and scalability.

## Vision

CannaNote evolves from recreational cannabis tracking to wellness optimization, culminating in a full HIPAA-compliant medical cannabis care platform for patients and providers. Built with hexagonal architecture principles for decades of maintainability.

## Core Opinions

- **"Core domain is pure—no I/O"** - Business logic isolated from external dependencies
- **Humans, not users** - Respectful, professional terminology throughout
- **Hexagonal architecture first** - Ports & adapters for maximum flexibility
- **No formal DI** - Security best practice, explicit dependency passing
- **Test coverage 80%+** - Quality enforced, maintainability guaranteed

## Project Structure

```
cannanote/                  # Monorepo root
├── README.md               # This file - vision, setup, architecture decisions  
├── CLAUDE.md               # AI development guidance and strategic direction
├── docs/                   # Comprehensive technical documentation
├── supabase/               # Database schema/migrations (Supabase CLI)
├── backend/                # Go backend with hexagonal architecture
│   ├── cmd/                # Binaries: server (API + HTMX), future microservices
│   ├── internal/           # Core domain + adapters (Supabase, Gin, Resend)
│   ├── pkg/                # Shared utilities and domain models
│   ├── web/                # HTMX templates + Tailwind CSS static assets
│   └── fly.toml            # Deployment configuration for Fly.io
├── mobile/                 # Flutter app (iOS/Android) with offline capabilities
│   ├── lib/                # Dart code: screens, services, state management
│   └── pubspec.yaml        # Dependencies: Supabase Flutter, Provider
└── .github/workflows/      # CI/CD: Lint/test/build/deploy for all platforms
```

## Architecture Philosophy

### Hexagonal (Ports & Adapters) Design
- **Pure Domain Logic**: No I/O operations in `internal/core/domain/`
- **Port Interfaces**: Define contracts in `internal/core/ports/`  
- **Adapter Implementations**: External services in `internal/adapters/`
- **Progressive Extraction**: Start monolithic, evolve to microservices

### Technology Stack
- **Backend**: Go 1.22+ with Gin framework (go-blueprint scaffold)
- **Database**: Supabase PostgreSQL with Row Level Security (HIPAA-ready)
- **Frontend**: HTMX + templ templates + Tailwind CSS (server-rendered)
- **Mobile**: Flutter with Supabase client and offline synchronization
- **Deployment**: Docker containers on Fly.io with HIPAA compliance add-ons
- **Integration**: Microsoft Graph API for calendar correlation

### Domain-Driven Design
```
Core Domains:
├── Human Management: Auth, profiles, consent (foundation domain)
├── Experience Tracking: Entries, strains, effects (core business value)
├── Analytics: Insights, correlations, B2B dashboards
└── Partnership: Dispensary integrations, affiliate management
```

## Development Phases

### Phase 1: Recreational MVP (Weeks 1-12)
- **Goal**: $5K MRR through freemium SaaS and dispensary partnerships
- **Tech**: Go hexagonal architecture, Supabase free tier
- **Legal**: No medical claims, no PHI collection, recreational tracking only

### Phase 2: Wellness Platform (Months 3-9)  
- **Goal**: $20K MRR with wellness correlation features
- **Tech**: Mobile app launch, Microsoft Graph calendar integration
- **Legal**: Wellness tracking without medical claims

### Phase 3: Medical Compliance (Months 9+)
- **Goal**: $75K+ MRR through healthcare provider partnerships
- **Tech**: Full HIPAA compliance, enterprise B2B features
- **Legal**: Medical cannabis care platform with provider integrations

## Getting Started

### Prerequisites
- Go 1.22+
- Node.js 18+ (for Tailwind CSS compilation)
- Flutter 3.10+ (for mobile development)
- Supabase CLI
- Docker

### Backend Setup
```bash
cd backend
# Use go-blueprint for initial scaffold
go-blueprint create --framework gin --driver postgres --htmx --tailwind --docker --github-actions

# Install dependencies
go mod tidy

# Start development server
make dev
```

### Supabase Setup
```bash
# Initialize Supabase project
supabase init

# Link to remote project
supabase link --project-ref your-project-ref

# Push schema changes
supabase db push
```

### Mobile Setup
```bash
cd mobile
flutter create .
flutter pub add supabase_flutter http provider
flutter run
```

## Key Features

### Human-Centric Design
- **Cannabis Journey Tracking**: Strains, effects, mood correlation
- **Privacy-First**: Row Level Security, consent management
- **Wellness Insights**: Calendar integration for lifestyle correlation
- **Social Features**: Community sharing, dispensary partnerships

### Technical Excellence
- **Security**: No formal DI, Trivy scans, encrypted PHI storage
- **Performance**: Server-rendered HTMX, efficient mobile offline sync
- **Scalability**: Microservice extraction path, event sourcing ready
- **Maintainability**: 80%+ test coverage, clean architecture principles

### Compliance Ready
- **HIPAA Preparation**: Field-level encryption, audit logging, BAAs
- **SOC2 Controls**: Access management, incident response procedures  
- **Accessibility**: WCAG 2.1 AAA compliance for inclusive design
- **International**: GDPR compliance for global expansion

## Contributing

This is a lifelong opus project built with hexagonal architecture for maximum maintainability. All contributions should:

1. **Preserve Core Purity**: No I/O operations in domain logic
2. **Follow Interface Contracts**: Use existing ports, create adapters
3. **Maintain Test Coverage**: 80%+ coverage enforced in CI
4. **Respect Human Terminology**: Use "humans" not "users"
5. **Document Architectural Decisions**: Update ADRs for major changes

## License

Proprietary - All rights reserved. This is a commercial cannabis platform with plans for open-source components in the future.

---

*Built with ♥️ and hexagonal architecture for the cannabis community*