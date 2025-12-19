# CannaNote CI/CD Workflows

GitHub Actions for automated testing, building, and deployment across all platforms.

## Workflows

### Backend (Go)
- **Lint**: golangci-lint with opinionated rules
- **Test**: 80%+ coverage requirement enforced
- **Security**: Trivy vulnerability scanning
- **Deploy**: Automated deployment to Fly.io

### Mobile (Flutter)
- **Test**: Unit and widget tests
- **Build**: iOS and Android builds
- **Deploy**: App store deployment (future)

### Supabase
- **Schema**: Automated migration testing
- **Security**: RLS policy validation

## Planned Features

- Automated dependency updates
- Performance testing and benchmarks
- Compliance verification checks
- Multi-environment deployment strategies