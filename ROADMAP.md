# CannaNote 2.0 Roadmap - Medical Cannabis Tracking Platform

## Project Vision

CannaNote 2.0 is a comprehensive medical cannabis tracking platform designed specifically for medicinal cannabis patients. The platform enables granular tracking of the complete consumption lifecycle, allowing users to self-identify their conditions, monitor treatment effects, and connect with other patients sharing similar experiences.

**Mission**: Empower medical cannabis patients with data-driven insights to optimize their treatment outcomes through community-driven knowledge sharing and AI-powered analysis.

**Target Market**: Medical cannabis patients seeking to optimize their treatment regimens through detailed tracking and community insights.

## Business Model

### Subscription Tiers with Premium Benefits

| Feature | Free | Basic ($4.20/mo, $42/yr) | Premium ($7.10/mo, $71/yr) |
|---------|------|---------------------------|-----------------------------|
| **Core Tracking** | ✅ | ✅ | ✅ |
| **Analytics History** | 30 days | Unlimited | Unlimited |
| **Data Export** | ❌ | ✅ CSV | ✅ CSV/PDF/Print |
| **Photo Uploads** | ❌ | 5/month | Unlimited |
| **Advanced Charts** | ❌ | ✅ | ✅ Enhanced AI insights |
| **Email Reports** | ❌ | Monthly summary | Weekly + custom reports |
| **AI Analysis** | ❌ | ❌ | ✅ Pattern recognition |
| **Priority Support** | ❌ | ❌ | ✅ Fast response |
| **POD Discounts** | ❌ | 10% off | 20% off |
| **Beta Features** | ❌ | ❌ | ✅ Early access |

### **Annual Subscription Exclusive Benefits**

**"Cannabis Journey Year in Review" (Annual subscribers only):**
- **Basic Annual ($42/year)**: Custom data visualization poster ($25 value) featuring:
  - Personal consumption patterns and favorite strains
  - "Spotify Wrapped" style infographic of cannabis journey
  - Achievement milestones and progress tracking
  - High-quality print delivered annually

- **Premium Annual ($71/year)**: Professional medical progress report ($50 value) featuring:
  - AI-powered treatment effectiveness analysis
  - Healthcare provider-ready symptom improvement report
  - Custom infographic of medical cannabis outcomes
  - Shareable HIPAA-compliant format
  - $20 credit for additional POD products

### Revenue Streams
1. **Subscription Revenue**: Monthly/annual premium tiers with immediate value
2. **Print-on-Demand Products**: Personalized cannabis tracking journals, data visualizations, medical reports
3. **Dispensary Partnerships**: Product placement and referral fees from Washington state dispensaries
4. **Premium Services**: Custom health reports, healthcare provider integrations
5. **Community Features**: Premium user matching and consultation services

## Complete Technical Rewrite

### Revised Stack for Lean Startup + Future Scale
| Component | Current | New (MVP) | Future Scale | Rationale |
|-----------|---------|-----------|--------------|-----------|
| **Backend** | Node.js/Express | **PocketBase + Go extensions** | Custom Go microservices | Rapid development + learning Go |
| **Database** | MongoDB | **NeonDB (PostgreSQL)** | Multi-region NeonDB | HIPAA compliance + scalability |
| **Frontend** | EJS Templates | **HTMX + Tailwind** | Same | Modern UX, reduced complexity |
| **Deployment** | Heroku | **Railway.app free tier** | Kubernetes/Docker | $0-5/month vs $200+/month |
| **Authentication** | Express-session | **Supabase Auth** | Custom JWT | 50K MAU free, OAuth built-in |
| **Payments** | None | **Stripe + Buy Me Coffee** | Multi-provider system | Subscriptions + tips/donations |
| **Storage** | Local | **Cloudflare R2** | Same | Free tier, global CDN |
| **Print-on-Demand** | None | **Shopify integration** | Custom fulfillment API | Revenue diversification |

### Implementation Steps

## Phase 1: Foundation & MVP (Months 1-3)

### 1.1 MVP Architecture: PocketBase + Go Extensions

**Hybrid Architecture for Rapid Development:**
```
cannanote-v2/
├── pocketbase/           # PocketBase binary + config
│   ├── pb_data/         # Database and files
│   ├── pb_hooks/        # Custom Go hooks
│   └── pb_migrations/   # Schema migrations
├── extensions/          # Custom Go services
│   ├── ai-analysis/     # ML pattern recognition
│   ├── pdf-generator/   # Year in review reports
│   ├── scraper/         # Dispensary data collection
│   └── payments/        # Stripe integration
├── web/
│   ├── static/          # Tailwind CSS, HTMX
│   ├── templates/       # Go templates
│   └── components/      # Reusable HTMX components
├── shopify/             # Print-on-demand integration
│   ├── products/        # POD product templates
│   └── webhooks/        # Order processing
└── deploy/
    ├── railway.toml     # Railway deployment config
    └── docker/          # Containerization for scale
```

**Benefits of PocketBase MVP Approach:**
- **30x faster development**: Built-in admin, auth, real-time, file storage
- **Single binary deployment**: No complex orchestration needed
- **Learn Go gradually**: Start with PocketBase hooks, add custom services
- **Future-proof**: Can extract to microservices as you scale

**MVP Implementation Steps (Weeks, not Months):**
1. **Week 1**: Download PocketBase, set up basic schema, deploy to Railway
2. **Week 2**: Configure Supabase Auth, create consumption tracking forms with HTMX
3. **Week 3**: Implement Stripe subscriptions with immediate premium features
4. **Week 4**: Build PDF generation for basic analytics reports

**Immediate Premium Value Features (Month 1):**
```go
// PocketBase hooks for premium features
func (app *pocketbase.PocketBase) bindPremiumHooks() {
    // Data export for premium users
    app.OnRecordAfterCreateRequest("consumption_entries").Add(func(e *core.RecordCreateEvent) error {
        if user.IsPremium() {
            generateWeeklyReport(user)
        }
        return nil
    })
    
    // Year in review generation
    app.OnRecordAfterUpdateRequest("users").Add(func(e *core.RecordUpdateEvent) error {
        if isAnnualRenewal(e.Record) {
            triggerYearInReviewGeneration(e.Record)
        }
        return nil
    })
}
```

**Go Learning Path Through Extensions:**
- **Month 1**: PocketBase hooks and simple Go templates
- **Month 3**: Custom PDF generation with Go libraries
- **Month 6**: AI analysis microservice in Go
- **Month 9**: Full custom backend migration option

### 1.2 Database Design: NeonDB with HIPAA Compliance

**Why NeonDB over TursoDB for Medical Data:**
- ✅ **HIPAA Compliance**: BAA available, SOC2 Type II certified
- ✅ **PostgreSQL**: More mature for complex medical queries
- ✅ **Branching**: Database branches for safe migrations
- ✅ **Auto-scaling**: Scales to zero when not in use
- **Cost**: $19/month Pro plan vs TursoDB $8/month (worthwhile for compliance)

**PocketBase Schema Configuration:**
```javascript
// PocketBase collection definitions (replaces SQL migrations)
const collections = [
  {
    name: "users",
    type: "auth",
    schema: [
      { name: "medical_card_number", type: "text", options: { encrypted: true }},
      { name: "state", type: "text", required: true },
      { name: "conditions", type: "json" }, // Array of medical conditions
      { name: "privacy_level", type: "number", default: 1 },
      { name: "subscription_tier", type: "select", options: ["free", "basic", "premium"] },
      { name: "annual_subscriber", type: "bool", default: false },
      { name: "subscription_expires", type: "date" }
    ]
  },

  {
    name: "consumption_entries", 
    schema: [
      { name: "user", type: "relation", options: { collection: "users" }},
      { name: "product_name", type: "text", required: true },
      { name: "product_type", type: "select", options: ["flower", "edible", "concentrate", "topical"] },
      { name: "strain", type: "text" },
      { name: "thc_percentage", type: "number" },
      { name: "cbd_percentage", type: "number" },
      { name: "consumption_method", type: "select", options: ["vape", "smoke", "edible", "sublingual", "topical"] },
      { name: "amount_mg", type: "number" },
      { name: "amount_description", type: "text" },
      { name: "consumed_at", type: "date", required: true },
      { name: "location", type: "text" },
      { name: "purpose", type: "text" }, // pain relief, sleep, anxiety, etc.
      { name: "photos", type: "file", options: { maxSelect: 3, maxSize: 5242880 }} // 5MB limit
    ],
    indexes: ["user", "consumed_at", "strain"]
  },

  {
    name: "consumption_effects",
    schema: [
      { name: "consumption_entry", type: "relation", options: { collection: "consumption_entries" }},
      { name: "effect_type", type: "select", options: ["pain_relief", "euphoria", "drowsiness", "anxiety_relief", "focus", "creativity", "appetite", "nausea_relief"] },
      { name: "intensity", type: "number", min: 1, max: 10 },
      { name: "onset_time_minutes", type: "number" },
      { name: "duration_minutes", type: "number" },
      { name: "notes", type: "editor" } // Rich text for detailed notes
    ]
  },
  {
    name: "year_in_review_orders",
    schema: [
      { name: "user", type: "relation", options: { collection: "users" }},
      { name: "year", type: "number", required: true },
      { name: "product_type", type: "select", options: ["poster", "journal", "medical_report"] },
      { name: "generated_data", type: "json" }, // Cached analytics for printing
      { name: "shopify_order_id", type: "text" },
      { name: "fulfillment_status", type: "select", options: ["pending", "generated", "ordered", "shipped", "delivered"] },
      { name: "shipping_address", type: "json" }
    ]
  }
];
```

**Automatic HIPAA Compliance Features:**
- **Encryption at rest**: NeonDB handles automatically
- **Access logging**: PocketBase logs all data access
- **Data retention**: Configurable auto-deletion policies
- **Backup encryption**: Automated encrypted backups

**Implementation Steps:**
1. Create database migration system using golang-migrate
2. Implement repository pattern for data access
3. Set up database connection pooling and health checks
4. Create seed data for testing
5. Implement soft deletes for user data compliance

### 1.3 Authentication & Authorization

**JWT-based Authentication Implementation:**
```go
// JWT service with proper security
type AuthService struct {
    jwtSecret []byte
    tokenTTL  time.Duration
    db        *sql.DB
}

// Implement OAuth2 for social login
// Add rate limiting for auth endpoints
// Implement proper password hashing with bcrypt
// Add email verification flow
// Medical card verification process
```

**Implementation Steps:**
1. Create JWT service with proper token validation
2. Implement middleware for protected routes
3. Add rate limiting to prevent brute force attacks
4. Create email verification system
5. Build medical card verification workflow
6. Implement proper session management

### 1.4 Core API Development

**RESTful API Endpoints:**
```
POST   /api/auth/register        # User registration
POST   /api/auth/login           # User login
POST   /api/auth/verify-email    # Email verification
GET    /api/users/profile        # User profile
PUT    /api/users/profile        # Update profile
POST   /api/consumption          # Log consumption
GET    /api/consumption          # Get consumption history
PUT    /api/consumption/:id      # Update consumption entry
POST   /api/consumption/:id/effects # Log effects
GET    /api/products/search      # Search products
GET    /api/users/public         # Find other users (filtered by privacy)
```

**Implementation Steps:**
1. Create handlers with proper input validation
2. Implement pagination for list endpoints
3. Add comprehensive error handling
4. Create OpenAPI/Swagger documentation
5. Implement request/response logging
6. Add health check endpoints for monitoring

### 1.5 Frontend with HTMX

**HTMX Integration Pattern:**
```html
<!-- Real-time consumption tracking form -->
<form hx-post="/api/consumption" 
      hx-target="#consumption-list" 
      hx-swap="afterbegin">
    <!-- Form fields for consumption entry -->
</form>

<!-- Live effects tracking -->
<div hx-trigger="every 30s" 
     hx-get="/api/consumption/recent-effects"
     hx-target="#effects-timeline">
    <!-- Effects timeline updates -->
</div>
```

**Implementation Steps:**
1. Create responsive layouts with Tailwind CSS
2. Implement progressive enhancement patterns
3. Build consumption tracking forms with real-time validation
4. Create effects timeline with live updates
5. Implement search and filtering with HTMX
6. Add offline-first capabilities with service workers

## Phase 2: AI Analysis & Recommendations (Months 4-6)

### 2.1 Data Collection & Web Scraping

**Web Scraping Infrastructure:**
```go
// Washington state dispensary scraping
type ScraperService struct {
    client   *http.Client
    db       *sql.DB
    logger   *zerolog.Logger
    rateLimiter *rate.Limiter
}

// Target sites: Leafly, Weedmaps, dispensary websites
// Respect robots.txt and implement proper rate limiting
// Store product data with source attribution
```

**Implementation Steps:**
1. Build web scraper for major dispensary websites
2. Create product normalization and deduplication system
3. Implement rate limiting to respect target sites
4. Set up automated scraping schedules
5. Create data validation and quality checks
6. Build admin dashboard for monitoring scraped data

### 2.2 AI Analysis Engine

**Consumption Pattern Analysis:**
```go
type AnalysisService struct {
    db          *sql.DB
    mlModel     MLPredictor
    insights    InsightGenerator
}

// Features to analyze:
// - Consumption timing patterns
// - Product effectiveness by condition
// - Dosage optimization
// - Onset/duration patterns
// - Cross-product recommendations
```

**Implementation Steps:**
1. Create data preprocessing pipeline for ML analysis
2. Implement basic pattern recognition algorithms
3. Build recommendation engine using collaborative filtering
4. Create personalized insights dashboard
5. Implement A/B testing framework for recommendations
6. Add machine learning model training pipeline

### 2.3 Premium Features Implementation

**Premium Subscription System:**
- Stripe integration for payments
- Usage tracking and billing
- Feature flagging system
- Advanced analytics dashboard
- Export functionality for user data

## Phase 3: Scaling & DevOps (Months 7-9)

### 3.1 Containerization & Orchestration

**Docker Configuration:**
```dockerfile
# Multi-stage build for optimal image size
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

**Kubernetes Deployment:**
- Create Kubernetes manifests for production deployment
- Implement horizontal pod autoscaling
- Set up ingress controllers and load balancing
- Create persistent volume claims for database storage

### 3.2 Monitoring & Observability

**Comprehensive Monitoring Stack:**
```go
// Prometheus metrics
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
)

// Structured logging with context
logger.Info().
    Str("user_id", userID).
    Str("action", "consumption_logged").
    Dur("duration", time.Since(start)).
    Msg("Successfully logged consumption entry")
```

**Implementation Steps:**
1. Integrate Prometheus for metrics collection
2. Set up Grafana dashboards for monitoring
3. Implement distributed tracing with Jaeger
4. Create alerting rules for critical system events
5. Set up log aggregation with ELK stack
6. Implement health checks and readiness probes

### 3.3 Performance Optimization

**Scaling Strategies:**
- Database read replicas for improved read performance
- Redis caching layer for frequently accessed data
- CDN integration for static asset delivery
- Database connection pooling optimization
- Query optimization and indexing strategy

**Load Testing:**
- Create comprehensive load testing suite
- Test for 1M+ concurrent user scenarios
- Identify and optimize bottlenecks
- Implement circuit breaker patterns
- Add graceful degradation mechanisms

## Compliance & Legal Framework

### HIPAA Compliance Considerations
- **Data Encryption**: At rest and in transit using AES-256
- **Access Controls**: Role-based access with audit logging
- **Data Retention**: Configurable retention policies
- **Breach Notification**: Automated incident response procedures
- **Business Associate Agreements**: With any third-party integrations

### Cannabis Regulations
- **State-Specific Compliance**: Washington state medical cannabis laws
- **Age Verification**: Robust identity verification system
- **Data Sovereignty**: Ensure data residency requirements
- **Product Verification**: Integration with state tracking systems
- **Advertising Restrictions**: Compliant marketing materials

### Privacy Controls
- **Granular Privacy Settings**: Control over data sharing levels
- **Data Portability**: Export functionality for user data
- **Right to Deletion**: Complete data removal capabilities
- **Anonymization**: Options for anonymous participation
- **Consent Management**: Clear consent flows for data usage

## Go-to-Market Strategy

### Phase 1: MVP Launch with Premium Benefits (Months 1-4)
**Target**: 200 users, 20 annual subscribers, break-even

**Immediate Value Strategy:**
- Launch with annual subscriptions offering instant benefits
- Focus on "Year in Review" exclusive value proposition  
- Target early adopters willing to pay upfront for personalized products

**Strategy:**
1. **Partner with 3-5 dispensaries** for user acquisition
   - Offer free premium trials to dispensary customers
   - Create co-marketing materials focusing on treatment optimization
   - Implement referral tracking for partnership ROI

2. **Focus on High-Impact Conditions**
   - Chronic pain management (largest patient group)
   - Anxiety and PTSD treatment
   - Cancer treatment support
   - Epilepsy and seizure management

3. **Content Marketing**
   - Educational blog content about medical cannabis dosing
   - Patient success stories (with consent)
   - Healthcare provider information packets
   - SEO optimization for medical cannabis search terms

**Success Metrics:**
- User registration conversion rate from dispensary partnerships
- Daily active usage rates (target: 60% weekly retention)
- Premium conversion rate from free trial users

### Phase 2: Revenue Growth & Product Expansion (Months 4-8) 
**Target**: 1,000 users, 100 annual subscribers, $5k MRR

**Premium Product Strategy:**
- Launch print-on-demand integration with Shopify
- Implement AI-powered "Year in Review" generation
- Add healthcare provider report features
- Expand to quarterly personalized products for annual subscribers

**Strategy:**
1. **Premium Feature Launch**
   - AI-powered consumption pattern analysis
   - Personalized product recommendations
   - Advanced analytics and trend tracking
   - Integration with Washington state dispensary inventory

2. **Strategic Partnerships**
   - Major dispensary chain partnerships (contracts with 2-3 chains)
   - Healthcare provider pilot programs
   - Patient advocacy group collaborations

3. **Product Development**
   - Web scraping infrastructure for real-time product data
   - Machine learning model training and optimization
   - Advanced reporting and export capabilities

### Phase 3: Scale & Compliance (Months 8-12)
**Target**: 5,000 users, $15k MRR

**Enterprise & Compliance Strategy:**
- Full HIPAA compliance certification
- Healthcare provider dashboard and integrations
- Multi-state expansion (Oregon, California)
- B2B dispensary partnership program

**Strategy:**
1. **Geographic Expansion**
   - Oregon and California medical markets
   - Replicate Washington state partnership model
   - Adapt to state-specific regulations and product catalogs

2. **Platform Maturation**
   - Advanced AI recommendations with larger dataset
   - Community features and patient matching
   - Healthcare provider dashboard and integration APIs

## Technical Implementation Timeline with Premium Benefits

### Month 1: MVP Foundation + Immediate Premium Value
- [ ] PocketBase setup with NeonDB backend
- [ ] Supabase authentication integration  
- [ ] Basic consumption tracking with HTMX forms
- [ ] Stripe subscription integration (monthly + annual)
- [ ] **Premium Feature**: Data export (CSV/PDF generation)
- [ ] **Premium Feature**: Photo uploads for consumption entries
- [ ] Railway deployment with custom domain

### Month 2: Analytics & Premium Reports
- [ ] Advanced charts and visualization for premium users
- [ ] Email report system (weekly summaries for premium)
- [ ] **Annual Exclusive**: Basic "Year in Review" PDF generation
- [ ] User profile privacy controls
- [ ] Mobile-responsive HTMX components

### Month 3: Print-on-Demand Integration
- [ ] Shopify store setup for POD products
- [ ] **Annual Exclusive**: Custom poster/journal generation
- [ ] Automated order fulfillment webhooks
- [ ] Payment processing for physical products
- [ ] Customer shipping address management

### Month 4: AI Analysis & Medical Reports 
- [ ] Basic pattern recognition algorithms
- [ ] **Premium Annual**: Professional medical progress reports
- [ ] Healthcare provider-ready PDF formatting
- [ ] Treatment effectiveness analysis
- [ ] HIPAA-compliant data handling protocols

### **Zero-Cost Infrastructure Strategy (Months 1-4)**

**Phase 0: $0.83/month (0-100 users) - Pure Free Tier**

| Service | Provider | Cost | Limits |
|---------|----------|------|--------|
| **Backend** | Railway.app | $0 | 500GB transfer, $5 credit/month |
| **Database** | NeonDB | $0* | 0.5GB storage (*HIPAA requires Pro $19/month) |
| **Frontend/CDN** | Vercel + Cloudflare | $0 | 100GB bandwidth |
| **Authentication** | Supabase Auth | $0 | 50,000 MAU |
| **Email** | SendGrid | $0 | 100 emails/day |
| **Print-on-Demand** | Shopify Starter | $5/month | Only when POD launches |
| **Domain + SSL** | Namecheap + Let's Encrypt | $10/year | = $0.83/month |

**Total Infrastructure Cost: $19.83/month (with HIPAA compliance)**

**Break-even: 3-5 paid subscribers at $4.20-7.10/month**

### **Scaling Thresholds & Cost Management**

**Phase 1: $50/month (100-500 users)**
- Railway Pro: $20/month
- NeonDB Scale: $69/month  
- Enhanced email: $15/month
- **Total: $104/month = $0.21 per user at 500 users**

**Phase 2: $200/month (500-2000 users)**  
- Multiple Railway services: $100/month
- NeonDB Pro: $169/month
- Advanced monitoring: $30/month
- **Total: $299/month = $0.15 per user at 2000 users**

### **Financial Break-Even Analysis**

| Month | Users | Annual Subs | Monthly Subs | Revenue | Infrastructure | Profit |
|-------|-------|-------------|--------------|---------|----------------|---------|
| 1     | 50    | 3 ($126)    | 2 ($14)     | $140    | $20           | $120    |
| 2     | 100   | 8 ($336)    | 5 ($35)     | $371    | $20           | $351    |
| 3     | 200   | 15 ($630)   | 10 ($70)    | $700    | $50           | $650    |
| 4     | 400   | 25 ($1,050) | 20 ($140)   | $1,190  | $100          | $1,090  |

**Key Insight: Profitable from Month 1 with just 3 annual subscribers**

## Risk Mitigation

### Technical Risks
- **Performance**: Regular load testing and optimization
- **Security**: Automated security scanning and regular audits
- **Data Loss**: Comprehensive backup and disaster recovery plans
- **Scalability**: Cloud-native architecture with auto-scaling

### Business Risks
- **Regulatory Changes**: Legal counsel and compliance monitoring
- **Market Competition**: Unique value proposition through medical focus
- **User Acquisition**: Diversified marketing channels and partnerships
- **Revenue Model**: Multiple revenue streams and flexible pricing

### Operational Risks
- **Team Scaling**: Clear documentation and knowledge sharing
- **Technology Debt**: Regular refactoring and code quality maintenance
- **Vendor Dependencies**: Abstraction layers and fallback options

## Success Metrics & KPIs

### User Engagement
- Daily/Weekly/Monthly active users
- Session duration and frequency
- Feature usage patterns
- User retention cohort analysis

### Business Metrics
- Customer acquisition cost (CAC)
- Lifetime value (LTV)
- Monthly recurring revenue (MRR)
- Premium conversion rates

### Product Metrics
- Consumption entry frequency
- Effects logging completion rates
- Recommendation click-through rates
- User-generated content quality

### Technical Metrics
- API response times (target: <100ms p95)
- System uptime (target: 99.9%)
- Error rates and resolution times
- Database performance metrics

## Conclusion

This roadmap represents a complete transformation of CannaNote from a general cannabis tracking app to a specialized medical cannabis platform built with modern, scalable technology. The focus on medical patients, combined with robust technical architecture and clear go-to-market strategy, positions the platform for sustainable growth and meaningful impact in the medical cannabis community.

The phased approach allows for validation of key assumptions while building towards the ultimate vision of serving 1M+ medical cannabis patients with AI-powered treatment optimization tools.