# CannaNote Development Roadmap

## Core Philosophy: Kina'ole Development

The primary metric for success is reducing cannabis session logging to under 30 seconds on mobile devices. All features and infrastructure decisions must align with our privacy promises—kina'ole development where our actions match our values.

## PRIORITY PHASE: Privacy Foundation

**Status**: Our privacy promises exceed our technical implementation. This must be corrected before other development.

### Data Encryption & Local-First
- **Client-Side Encryption** - Encrypt sensitive data before database storage
- **Local Storage Primary** - Web: IndexedDB, Mobile: Drift as primary storage
- **Encrypted Cloud Sync** - Sync only encrypted blobs and metadata
- **User-Controlled Keys** - Keys from user credentials, never server-side
- **Data Export** - Complete portable data export

### True Consent Management
- **Granular Permissions** - Specific, revocable consent per data type
- **Local-Only Default** - Cloud sync as explicit opt-in
- **Transparency Dashboard** - Real-time view of data storage locations
- **Complete Deletion** - Full data purging including backups
- **Consent History** - Audit trail of permission changes

### Privacy Validation
- **Automated Privacy Tests** - Verify encryption and access controls
- **Data Flow Mapping** - Document all personal data paths
- **Encryption Test Suite** - Comprehensive client-side crypto testing
- **Cross-Platform Consistency** - Identical privacy on web and mobile

**Goal**: Technical delivery of "radical data transparency" promise.

## Phase 1: Mobile Foundation & Core User Value

### Priority 1A: Offline-First Mobile Architecture
- **Flutter + Drift Implementation** - Replace basic sqflite with Drift for robust offline database
- **30-Second Logging UX** - Core journaling flow optimized for speed and simplicity
- **Local-First Data Storage** - Full functionality without network dependency
- **Background Sync** - Intelligent synchronization when network available
- **Cross-Platform Consistency** - iOS and Android parity for core features

### Priority 1B: Essential Harm Reduction Features
- **Dosage Calculator Service** - `core/application/dosage_service.go` with evidence-based calculations
- **Session Timing System** - Built-in reminders and spacing recommendations
- **Tolerance Break Suggestions** - Pattern-based recommendations for mindful usage
- **Safety Guardrails** - Dosage warnings and consumption frequency alerts
- **Educational Integration** - Contextual harm reduction information during logging

### Priority 1C: Core Backend Services
- **Mobile API Optimization** - Endpoints designed specifically for mobile offline-sync patterns
- **Authentication Flow** - Streamlined mobile auth with Supabase integration
- **Data Synchronization** - Conflict resolution and delta sync capabilities
- **Health Check Enhancement** - Mobile-specific connectivity and sync status monitoring

## Phase 2: User Experience & Documentation Polish

### Priority 2A: Documentation Enhancement
- **Root README Overhaul** - Quick-start guide with "make dev" prominence
- **Architecture Diagram** - Visual representation of hexagonal architecture and mobile-first design
- **Values Integration** - Explicit connections to privacy and harm reduction throughout docs
- **CONTRIBUTING.md Creation** - Values-based contribution guidelines emphasizing evidence-based development
- **Mobile Development Guide** - Comprehensive Flutter development documentation

### Priority 2B: Testing Foundation
- **Domain Logic Coverage** - High test coverage for `core/domain/` and `core/application/`
- **Mobile Integration Testing** - Flutter widget and integration test suites
- **Auth Flow Testing** - Comprehensive authentication and authorization test scenarios
- **Database Migration Testing** - Supabase RLS and schema validation testing
- **API Contract Testing** - Mobile-backend contract validation

### Priority 2C: Core Features Expansion
- **Pattern Recognition Engine** - Basic consumption pattern analysis and insights
- **Entry Enhancement** - Rich consumption method tracking and effects correlation
- **Search and Filtering** - Efficient local search with server backup
- **Data Export** - User-owned data export in multiple formats
- **Privacy Dashboard** - Granular privacy controls and data transparency

## Phase 3: Intelligence & Environmental Integration

### Priority 3A: Advanced Analytics
- **Machine Learning Pipeline** - Local ML models for personalized insights
- **Correlation Analysis** - Environmental factors and consumption outcome analysis
- **Recommendation Engine** - Strain and dosage recommendations based on patterns
- **Trend Visualization** - Charts and graphs for long-term pattern recognition
- **Health Integration** - Optional integration with health tracking platforms

### Priority 3B: Environmental Responsibility
- **Carbon Footprint Integration** - Simple endpoint for consumption environmental impact
- **Sustainable Cannabis Tracking** - Cultivation method and packaging impact awareness
- **Environmental Education** - Content about regenerative cannabis practices
- **Eco-Conscious Features** - Features that promote sustainable consumption practices

### Priority 3C: Community Features (Privacy-Optional)
- **Anonymous Insights Sharing** - Community pattern sharing without personal identification
- **Research Participation** - Opt-in anonymous data contribution for cannabis research
- **Educational Content Platform** - Peer-reviewed cannabis research integration
- **Expert Content Curation** - Integration with cannabis researchers and medical professionals

## Phase 4: Platform Maturation & Scaling

### Priority 4A: Advanced Technical Infrastructure
- **Microservice Architecture** - Domain extraction for independent scaling
- **GraphQL API** - Flexible query interface for advanced client applications
- **Real-Time Features** - WebSocket integration for live synchronization
- **Performance Optimization** - Advanced caching and rendering optimizations
- **Monitoring & Observability** - Comprehensive application performance monitoring

### Priority 4B: Platform Integration
- **Health App Integration** - Apple Health and Google Fit data sharing
- **Wearable Device Support** - Heart rate and activity data correlation
- **Calendar Integration** - Lifestyle correlation with calendar events
- **Third-Party APIs** - Integration with other wellness tracking platforms
- **Dispensary Integrations** - Product information and availability APIs

### Priority 4C: Enterprise & Research Capabilities
- **Healthcare Provider Integration** - Secure data sharing with medical professionals
- **Research Platform** - Anonymized aggregate data for scientific research
- **Compliance Framework** - Legal compliance monitoring for medical cannabis patients
- **Multi-Device Sync** - Seamless experience across multiple devices and platforms
- **Family Sharing** - Secure sharing between trusted family members and caregivers

## Risk Mitigation Strategies

### Technical Dependencies
- **Supabase Migration Path** - Documented PostgreSQL self-hosting migration strategy
- **Database Abstraction** - Maintain vendor-agnostic domain layer for database independence
- **Mobile Platform Changes** - Flexible architecture to adapt to iOS/Android policy changes
- **Network Dependency** - Full offline functionality as primary design constraint

### Product-Market Fit
- **Mobile vs Web Balance** - 80% of development decisions driven by mobile user experience
- **Feature Complexity** - Prioritize "breath-like" simplicity over feature richness
- **App Store Compliance** - "Wellness journaling" positioning to avoid cannabis-specific restrictions
- **Privacy Regulations** - Proactive compliance with evolving privacy legislation

### Business Sustainability
- **Open Source Strategy** - Core functionality remains open with optional premium features
- **Data Ownership** - User data portability and ownership as competitive advantage
- **Community Building** - Evidence-based community contributions over marketing-driven growth
- **Values Alignment** - Consistent harm reduction and privacy focus in all decisions

## Success Metrics

### Primary Metrics (Mobile-First Focus)
- **Logging Time** - Average time from app open to session logged (target: <30 seconds)
- **Offline Functionality** - Percentage of features available without network (target: >95%)
- **User Retention** - Monthly active users with consistent logging behavior
- **Data Quality** - Completeness and accuracy of user-entered consumption data

### Secondary Metrics (Supporting Infrastructure)
- **Sync Reliability** - Successful background synchronization rate
- **App Performance** - Load times, battery usage, and crash rates
- **Privacy Compliance** - User privacy preference adherence rate
- **Educational Impact** - User engagement with harm reduction content

### Long-Term Metrics (Platform Value)
- **Pattern Recognition Accuracy** - Quality of personalized insights and recommendations
- **Health Outcome Correlation** - User-reported wellness improvements
- **Community Contribution** - Anonymous research data contribution rate
- **Environmental Impact** - Carbon footprint reduction through informed consumption

## Implementation Notes

### Development Principles
- **Evidence-Based Features** - All features backed by research or user data
- **Privacy by Design** - Privacy considerations integrated at architecture level
- **Harm Reduction First** - User safety prioritized over engagement metrics
- **Sustainable Technology** - Environmental impact considered in technical decisions

### Quality Standards
- **Test Coverage** - Minimum 85% coverage for domain logic and critical paths
- **Performance Benchmarks** - Mobile app startup time <2 seconds, API response <500ms
- **Accessibility Compliance** - Full screen reader support and keyboard navigation
- **Security Standards** - Regular penetration testing and dependency security scanning

## Security & Privacy Roadmap (Integrated)

### Phase S1: Foundation Security (With Priority Phase)
- **RLS Policy Audit** - Strengthen Row Level Security on all tables
- **User Isolation** - Zero cross-user data access capability
- **Auth Hardening** - Rate limiting, email confirmation, CAPTCHA
- **Service Key Restriction** - Server-side only access controls
- **Transparency Documentation** - User-facing privacy technical docs

### Phase S2: Code Integrity
- **Security Linting** - gosec with cannabis data sensitivity rules
- **Dependency Scanning** - govulncheck and OSV Scanner
- **PII Exposure Prevention** - Automated detection of data leaks
- **Privacy-Safe CI/CD** - Pre-commit hooks for encryption validation
- **Audit Preparation** - Code structure for independent privacy reviews

### Phase S3: Secrets & Configuration
- **Secret Management** - Infisical or Doppler for cannabis data classification
- **Environment Isolation** - Production key separation
- **Key Rotation** - Automated rotation without user disruption
- **Configuration Transparency** - User documentation of third-party sharing
- **History Cleaning** - Remove any historical secret exposure

### Phase S4: Runtime Protection
- **Privacy-Safe Monitoring** - Sentry with PII filtering
- **Data Access Logging** - Audit trail without storing personal data
- **Encryption Monitoring** - Alerts for unencrypted personal data
- **User Privacy Dashboard** - Real-time data handling visibility
- **Privacy Incident Response** - Procedures aligned with harm reduction

### Phase S5: Advanced Validation
- **Regular Privacy Audits** - Quarterly technical reviews
- **User Privacy Education** - In-app technical protection explanations
- **Privacy Innovation** - Research emerging protection technologies
- **Community Advocacy** - Share learnings with cannabis tech community
- **Values-Aligned Evolution** - Privacy enhancements supporting harm reduction

## Privacy Validation Metrics

- **Encryption Coverage** - 100% of personal consumption data encrypted
- **Local-First Functionality** - 95%+ features available offline
- **User Data Control** - Complete export/deletion <24hr fulfillment
- **Access Isolation** - Zero cross-user data access incidents
- **Encryption Reliability** - 99.9%+ success rate for crypto operations