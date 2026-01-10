# CannaNote Development Roadmap

## Core Philosophy: Kina'ole Development

The primary metric for success is reducing cannabis session logging to under 30 seconds on mobile devices. All features and infrastructure decisions must align with our privacy promises—kina'ole development where our actions match our values.

## 🌱 GERMINATION CYCLE: Privacy Foundation

**Status**: Our privacy promises exceed our technical implementation. This foundation must be established before other cultivation can begin.

**Historical Context**: After performance optimization work (Alpine.js → vanilla JS, Tailwind v4 migration), we recognized that true privacy requires architectural changes, not just framework choices.

### Data Encryption & Local-First
- **Local Storage Primary** - Web: IndexedDB, Mobile: Drift as primary storage
- **Encrypted Cloud Sync** - Sync only encrypted blobs and metadata (paid feature)
- **Human-Controlled Keys** - Keys from human credentials, never server-side
- **Data Export** - Complete portable data export

### True Consent Management
- **Granular Permissions** - Specific, revocable consent per data type
- **Local-Only Default** - Cloud sync as explicit opt-in (premium feature)
- **Transparency Dashboard** - Real-time view of data storage locations
- **Complete Deletion** - Full data purging including backups
- **Consent History** - Audit trail of permission changes

### Privacy Validation Through Beta
- **Beta Launch as Privacy Test** - Real humans testing local-first architecture
- **Automated Privacy Tests** - Verify encryption and access controls
- **Data Flow Mapping** - Document all personal data paths
- **Encryption Test Suite** - Comprehensive client-side crypto testing
- **Cross-Platform Consistency** - Identical privacy on web and mobile

**Goal**: Technical delivery of "radical data transparency" promise validated by real cannabis community members.

## 🌿 SEEDLING CYCLE: PWA-First Beta Launch

**Strategic Decision**: PWA-first approach validates core experience and API design before Flutter investment, getting cannabis community feedback faster.

**Business Model Foundation**: Local-first with optional premium sync establishes sustainable revenue while respecting human autonomy.

### Authentication Foundation

#### Supabase Integration (Auth Only - No Cannabis Data)
- **Email Signup Flow** - Email verification → profile creation → beta status
- **Database Schema**:
  ```sql
  profiles (
    id UUID REFERENCES auth.users(id),
    email TEXT,
    beta_joined_at TIMESTAMP,
    subscription_tier TEXT DEFAULT 'beta_grandfathered',
    beta_status TEXT DEFAULT 'waitlist', -- waitlist, invited, active
    referral_code TEXT UNIQUE,
    invited_at TIMESTAMP,
    activated_at TIMESTAMP
  )
  ```
- **Row Level Security** - Humans can only access their own profile data
- **Beta Grandfathering** - Lifetime sync access for early community members

#### Email System Setup (Resend + cannanote.org)
- **Welcome Email** - Immediate signup confirmation with privacy commitment
- **Beta Invitation** - Personal access link when spots available  
- **Onboarding Sequence** - PWA installation and first session guidance
- **Privacy-Focused Tone** - Cannabis wellness aligned, educational approach

### PWA Foundation

#### Progressive Web App Implementation
- **Web App Manifest** - Home screen installation capability
- **Service Worker** - Full offline functionality for cannabis sessions
- **Install Prompts** - Native-like installation experience
- **Offline-First Architecture** - All core features work without network

#### Local Cannabis Data Storage
- **IndexedDB Schema** - Complete cannabis session storage
  ```javascript
  sessions: {
    id: uuid,
    timestamp: datetime,
    strain: { name, type, thc_percent, cbd_percent },
    consumption: { method, amount, unit },
    setting: { location_type, mood_before, social_context },
    effects: { onset_minutes, duration_hours, intensity_1_10, feelings: [] },
    notes: encrypted_string,
    rating: 1_to_5,
    privacy_level: local_only_or_sync_eligible
  }
  ```
- **Client-Side Encryption** - All personal cannabis data encrypted before storage
- **Data Export Tools** - Complete human data portability (CSV/JSON formats)

### Core Cannabis Experience

#### 30-Second Session Logging
- **Quick Entry Interface** - Touch-optimized for immediate post-consumption logging
- **Strain Database** - Local searchable cannabis strain library
- **Consumption Methods** - Vaping, smoking, edibles, tinctures, topicals
- **Effects Tracking** - Standardized feelings taxonomy + personal notes
- **Pattern Recognition** - Client-side analytics for personal insights

#### Beta Community Management
- **Waitlist Strategy** - Creates scarcity and quality control for early adoption
- **Manual Invitation System** - Curated onboarding for better community experience
- **Admin Dashboard** - Manage beta invitations, track community growth metrics
- **Community Feedback Collection** - In-app feedback and usage pattern analysis

### Growth Metrics (Seedling Health)
- **PWA Install Rate** - Percentage of humans adding to home screen
- **Session Logging Speed** - Time from app open to session saved (<30s target)
- **Offline Usage** - Sessions logged without network connection
- **Beta Community Satisfaction** - NPS and retention rates among early adopters

## 🍃 VEGETATIVE CYCLE: Mobile App Development & API Maturation

**Transition Point**: After PWA beta validates core experience and API design, begin Flutter development informed by real human usage patterns.

### Flutter + Drift Implementation
- **Native Mobile Experience** - iOS and Android apps with offline-first Drift database
- **API Optimization** - Backend endpoints refined based on PWA beta feedback
- **Cross-Platform Sync** - Local data synchronization between PWA and mobile
- **Performance Optimization** - Sub-2-second startup, immediate session logging

### Essential Harm Reduction Features
- **Dosage Calculator Service** - `core/application/dosage_service.go` with evidence-based calculations
- **Session Timing System** - Built-in reminders and spacing recommendations
- **Tolerance Break Suggestions** - Pattern-based recommendations for mindful usage
- **Safety Guardrails** - Dosage warnings and consumption frequency alerts
- **Educational Integration** - Contextual harm reduction information during logging

### Backend Services Expansion
- **Mobile API Endpoints** - Designed for offline-sync patterns based on PWA learnings
- **Authentication Flow** - Streamlined mobile auth with existing Supabase integration
- **Data Synchronization** - Conflict resolution and delta sync capabilities
- **Health Check Enhancement** - Mobile-specific connectivity and sync status monitoring

## 🌸 PRE-FLOWER CYCLE: Experience & Documentation Polish

### Documentation Enhancement
- **Root README Overhaul** - Quick-start guide reflecting PWA-first then mobile approach
- **Architecture Diagram** - Visual representation of local-first hexagonal architecture
- **Values Integration** - Cannabis wellness culture throughout technical documentation
- **CONTRIBUTING.md Creation** - Community contribution guidelines emphasizing evidence-based development
- **Mobile Development Guide** - Flutter development documentation informed by PWA experience

### Testing Foundation
- **Domain Logic Coverage** - High test coverage for `core/domain/` and `core/application/`
- **Mobile Integration Testing** - Flutter widget and integration test suites
- **Auth Flow Testing** - Comprehensive authentication and authorization scenarios
- **Database Migration Testing** - Supabase RLS and schema validation testing
- **API Contract Testing** - Mobile-backend contract validation based on PWA API usage

### Core Features Expansion
- **Pattern Recognition Engine** - Advanced consumption pattern analysis and insights
- **Entry Enhancement** - Rich consumption method tracking and effects correlation
- **Search and Filtering** - Efficient local search with optional cloud backup
- **Data Export Enhancement** - Multiple format support for human data ownership
- **Privacy Dashboard** - Granular privacy controls and data transparency

## 💐 FLOWERING CYCLE: Intelligence & Environmental Integration

### Advanced Analytics (Local-First)
- **Machine Learning Pipeline** - Local ML models for personalized insights
- **Correlation Analysis** - Environmental factors and consumption outcome analysis
- **Recommendation Engine** - Strain and dosage recommendations based on patterns
- **Trend Visualization** - Charts and graphs for long-term pattern recognition
- **Health Integration** - Optional integration with health tracking platforms

### Environmental Responsibility
- **Carbon Footprint Integration** - Simple endpoint for consumption environmental impact
- **Sustainable Cannabis Tracking** - Cultivation method and packaging impact awareness
- **Environmental Education** - Content about regenerative cannabis practices
- **Eco-Conscious Features** - Features that promote sustainable consumption practices

### Privacy-Optional Community Features
- **Anonymous Insights Sharing** - Community pattern sharing without personal identification
- **Research Participation** - Opt-in anonymous data contribution for cannabis research
- **Educational Content Platform** - Peer-reviewed cannabis research integration
- **Expert Content Curation** - Integration with cannabis researchers and medical professionals

## 🔥 HARVEST CYCLE: Platform Maturation & Scaling

### Advanced Technical Infrastructure
- **Microservice Architecture** - Domain extraction for independent scaling
- **GraphQL API** - Flexible query interface for advanced client applications
- **Real-Time Features** - WebSocket integration for live synchronization
- **Performance Optimization** - Advanced caching and rendering optimizations
- **Monitoring & Observability** - Comprehensive application performance monitoring

### Platform Integration
- **Health App Integration** - Apple Health and Google Fit data sharing
- **Wearable Device Support** - Heart rate and activity data correlation
- **Calendar Integration** - Lifestyle correlation with calendar events
- **Third-Party APIs** - Integration with other wellness tracking platforms
- **Dispensary Integrations** - Product information and availability APIs

### Enterprise & Research Capabilities
- **Healthcare Provider Integration** - Secure data sharing with medical professionals
- **Research Platform** - Anonymized aggregate data for scientific research
- **Compliance Framework** - Legal compliance monitoring for medical cannabis patients
- **Multi-Device Sync** - Seamless experience across multiple devices and platforms
- **Family Sharing** - Secure sharing between trusted family members and caregivers

## 🏺 CURING CYCLE: Optimization & Community Scaling

### Performance Refinement
- **Advanced Caching Strategies** - Multi-layer caching for optimal performance
- **Database Optimization** - Query optimization and indexing strategies
- **CDN Implementation** - Global content delivery for educational resources
- **Load Testing** - Community scale stress testing and optimization

### Community Platform Evolution
- **Community Guidelines** - Cannabis culture aligned community standards
- **Peer Support Systems** - Human-to-human wellness support networks
- **Educational Content Curation** - Community-contributed harm reduction resources
- **Research Collaboration** - Academic and medical research partnerships

## Risk Mitigation Strategies

### Technical Dependencies
- **Supabase Migration Path** - Documented PostgreSQL self-hosting migration strategy
- **Database Abstraction** - Maintain vendor-agnostic domain layer for database independence
- **Mobile Platform Changes** - Flexible architecture to adapt to iOS/Android policy changes
- **Network Dependency** - Full offline functionality as primary design constraint

### Product-Market Fit
- **Mobile vs Web Balance** - 80% of development decisions driven by mobile human experience
- **Feature Complexity** - Prioritize "breath-like" simplicity over feature richness
- **App Store Compliance** - "Wellness journaling" positioning to avoid cannabis-specific restrictions
- **Privacy Regulations** - Proactive compliance with evolving privacy legislation

### Business Sustainability
- **Open Source Strategy** - Core functionality remains open with optional premium features
- **Data Ownership** - Human data portability and ownership as competitive advantage
- **Community Building** - Evidence-based community contributions over marketing-driven growth
- **Values Alignment** - Consistent harm reduction and privacy focus in all decisions

## Growth Metrics

### Primary Metrics (Mobile-First Focus)
- **Logging Time** - Average time from app open to session logged (target: <30 seconds)
- **Offline Functionality** - Percentage of features available without network (target: >95%)
- **Community Retention** - Monthly active humans with consistent logging behavior
- **Data Quality** - Completeness and accuracy of human-entered consumption data

### Secondary Metrics (Supporting Infrastructure)
- **Sync Reliability** - Successful background synchronization rate
- **App Performance** - Load times, battery usage, and crash rates
- **Privacy Compliance** - Human privacy preference adherence rate
- **Educational Impact** - Community engagement with harm reduction content

### Long-Term Metrics (Platform Value)
- **Pattern Recognition Accuracy** - Quality of personalized insights and recommendations
- **Health Outcome Correlation** - Human-reported wellness improvements
- **Community Contribution** - Anonymous research data contribution rate
- **Environmental Impact** - Carbon footprint reduction through informed consumption

## Implementation Principles

### Development Philosophy
- **Evidence-Based Features** - All features backed by research or human data
- **Privacy by Design** - Privacy considerations integrated at architecture level
- **Harm Reduction First** - Human safety prioritized over engagement metrics
- **Sustainable Technology** - Environmental impact considered in technical decisions
- **Cannabis Culture Alignment** - Language and features reflect cannabis wellness values

### Quality Standards
- **Test Coverage** - Minimum 85% coverage for domain logic and critical paths
- **Performance Benchmarks** - Mobile app startup time <2 seconds, API response <500ms
- **Accessibility Compliance** - Full screen reader support and keyboard navigation
- **Security Standards** - Regular penetration testing and dependency security scanning

## Security & Privacy Cultivation

### Germination Security (With Privacy Foundation)
- **RLS Policy Audit** - Strengthen Row Level Security on all tables
- **Human Isolation** - Zero cross-human data access capability
- **Auth Hardening** - Rate limiting, email confirmation, CAPTCHA
- **Service Key Restriction** - Server-side only access controls
- **Transparency Documentation** - Human-facing privacy technical documentation

### Seedling Security (Beta Launch Security)
- **Security Linting** - gosec with cannabis data sensitivity rules
- **Dependency Scanning** - govulncheck and OSV Scanner
- **PII Exposure Prevention** - Automated detection of data leaks
- **Privacy-Safe CI/CD** - Pre-commit hooks for encryption validation
- **Beta Privacy Auditing** - Real-world privacy validation with beta community

### Vegetative Security (Mobile Security)
- **Secret Management** - Infisical or Doppler for cannabis data classification
- **Environment Isolation** - Production key separation
- **Key Rotation** - Automated rotation without human disruption
- **Configuration Transparency** - Human documentation of third-party sharing
- **Mobile Security Standards** - Platform-specific security implementations

### Flowering Security (Runtime Protection)
- **Privacy-Safe Monitoring** - Sentry with PII filtering
- **Data Access Logging** - Audit trail without storing personal data
- **Encryption Monitoring** - Alerts for unencrypted personal data
- **Human Privacy Dashboard** - Real-time data handling visibility
- **Privacy Incident Response** - Procedures aligned with harm reduction

### Harvest Security (Advanced Validation)
- **Regular Privacy Audits** - Quarterly technical reviews
- **Human Privacy Education** - In-app technical protection explanations
- **Privacy Innovation** - Research emerging protection technologies
- **Cannabis Community Advocacy** - Share learnings with cannabis tech community
- **Values-Aligned Evolution** - Privacy enhancements supporting harm reduction

## Privacy Validation Metrics

- **Encryption Coverage** - 100% of personal consumption data encrypted
- **Local-First Functionality** - 95%+ features available offline
- **Human Data Control** - Complete export/deletion <24hr fulfillment
- **Access Isolation** - Zero cross-human data access incidents
- **Encryption Reliability** - 99.9%+ success rate for crypto operations
- **Community Trust** - Human satisfaction with privacy transparency
- **Beta Grandfathering Success** - Lifetime access retention and satisfaction rates

## Historical Context & Decision Points

### Performance Optimization Era (Completed)
- **Framework Migration** - Alpine.js → Vanilla JS for reliability and performance
- **CSS Architecture** - Tailwind v3 → v4 with local build process
- **Bundle Optimization** - 28KB → 3KB JavaScript bundle with minification
- **Cache Strategy** - Static asset caching and CDN optimization

### PWA-First Strategic Decision
- **Rationale** - Faster community validation than waiting for Flutter development
- **Local-First Commitment** - Cannabis data never leaves device unless explicitly paid for sync
- **Beta Grandfathering** - Lifetime sync access rewards early community supporters
- **Cannabis Culture Language** - Cultivation terminology aligns with our community values

This roadmap represents our evolution from technical optimization to community-focused cannabis wellness platform, maintaining privacy promises while building sustainable business model through optional premium features.