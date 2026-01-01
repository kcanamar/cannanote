# CannaNote Database Architecture

Comprehensive database infrastructure for CannaNote personal cannabis journaling platform, built on Supabase PostgreSQL with privacy-first design, Row Level Security, and optimized JSONB storage for cannabis chemical profile data.

## Database Architecture Overview

The database implements **privacy-by-design architecture** with granular Row Level Security policies, efficient JSONB storage for complex cannabis data, and scalable schema design optimized for personal wellness tracking and optional community features.

### Architectural Principles

- **Privacy by Design** — User data isolation through Row Level Security policies
- **Data Ownership** — Users maintain complete control over personal consumption data
- **Efficient Storage** — JSONB optimization for complex cannabis chemical profiles
- **Scalable Design** — Schema architecture supporting personal to community scale
- **Reference Data Management** — Centralized scientific cannabis data with version control
- **Performance Optimization** — Strategic indexing for fast pattern recognition queries

### Technology Stack

- **Database Platform**: Supabase PostgreSQL 15+ with extensions
- **Security Model**: Row Level Security (RLS) with policy-based access control
- **Authentication**: Supabase Auth with JWT token validation
- **Storage Format**: JSONB for complex data structures with GIN indexing
- **Real-time Sync**: Supabase Realtime for optional multi-device synchronization
- **Development Tools**: Supabase CLI with local development environment
- **Migration Management**: SQL-based schema migrations with version control

## Project Structure

```
supabase/
├── config.toml                      # Supabase project configuration
├── seed.sql                         # Initial database seeding and test data
├── reference-data/                  # Scientific cannabis reference data
│   ├── cannabinoids.sql            # Cannabinoid reference data (THC, CBD, CBG, etc.)
│   └── terpenes.sql                # Terpene reference data (Myrcene, Limonene, etc.)
├── migrations/                      # Database schema migrations (generated)
│   ├── 20241201000001_initial_schema.sql     # Initial schema creation
│   ├── 20241201000002_rls_policies.sql       # Row Level Security policies
│   ├── 20241201000003_reference_data.sql     # Reference data population
│   └── 20241201000004_indexes.sql            # Performance optimization indexes
├── functions/                       # Supabase Edge Functions (planned)
│   ├── sync-conflict-resolution/    # Data synchronization conflict resolution
│   ├── pattern-recognition/         # Cannabis consumption pattern analysis
│   └── privacy-export/              # User data export functionality
├── storage/                         # Supabase Storage configuration (planned)
│   └── policies.sql                # Storage bucket policies for user uploads
├── types/                          # Generated TypeScript types
│   └── database.types.ts          # Auto-generated from schema
└── README.md                       # This documentation file
```

## Database Schema Design

### Core Entity Model

The database schema is optimized for personal cannabis tracking with efficient chemical profile storage and privacy-focused user data management.

#### User Management

```sql
-- Human profiles with cannabis preferences
CREATE TABLE humans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    profile JSONB DEFAULT '{}',
    privacy_settings JSONB DEFAULT '{
        "data_sync_enabled": false,
        "community_visible": false,
        "analytics_participation": false
    }',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### Cannabis Experience Tracking

```sql
-- Cannabis consumption entries with privacy controls
CREATE TABLE entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    human_id UUID NOT NULL REFERENCES humans(id) ON DELETE CASCADE,
    product_name VARCHAR(255) NOT NULL,
    brand VARCHAR(255),
    consumption_method_id UUID REFERENCES consumption_methods(id),
    chemical_profile_id UUID REFERENCES chemical_profiles(id),
    dose_amount DECIMAL(10,3),
    dose_unit VARCHAR(50),
    effects JSONB DEFAULT '{}',
    mood_before INTEGER CHECK (mood_before BETWEEN 1 AND 10),
    mood_after INTEGER CHECK (mood_after BETWEEN 1 AND 10),
    environment JSONB DEFAULT '{}',
    notes TEXT,
    rating INTEGER CHECK (rating BETWEEN 1 AND 10),
    consumed_at TIMESTAMP WITH TIME ZONE NOT NULL,
    visibility VARCHAR(20) DEFAULT 'private' CHECK (visibility IN ('private', 'community', 'followers')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### Cannabis Chemical Profiles

```sql
-- Chemical profiles with JSONB optimization for lab data
CREATE TABLE chemical_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_name VARCHAR(255),
    brand VARCHAR(255),
    batch_id VARCHAR(255),
    lab_name VARCHAR(255),
    test_date DATE,
    harvest_date DATE,
    cannabinoids JSONB DEFAULT '{}',
    terpenes JSONB DEFAULT '{}',
    other_compounds JSONB DEFAULT '{}',
    total_cannabinoids DECIMAL(5,2),
    total_terpenes DECIMAL(5,2),
    lab_report_url TEXT,
    verified_lab BOOLEAN DEFAULT FALSE,
    data_source VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Efficient JSONB indexing for chemical profile queries
CREATE INDEX idx_chemical_profiles_cannabinoids ON chemical_profiles USING GIN (cannabinoids);
CREATE INDEX idx_chemical_profiles_terpenes ON chemical_profiles USING GIN (terpenes);
CREATE INDEX idx_chemical_profiles_total_cannabinoids ON chemical_profiles (total_cannabinoids);
```

### Reference Data Tables

#### Cannabis Cannabinoids Reference

```sql
-- Scientific cannabinoid reference data
CREATE TABLE cannabinoids (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(255),
    description TEXT,
    molecular_formula VARCHAR(50),
    common_effects JSONB DEFAULT '[]',
    therapeutic_potential JSONB DEFAULT '[]',
    research_status VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### Cannabis Terpenes Reference

```sql
-- Scientific terpene reference data
CREATE TABLE terpenes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    common_name VARCHAR(100),
    aroma_profile JSONB DEFAULT '[]',
    flavor_notes JSONB DEFAULT '[]',
    common_effects JSONB DEFAULT '[]',
    boiling_point DECIMAL(6,2),
    also_found_in JSONB DEFAULT '[]',
    research_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### Consumption Methods

```sql
-- Consumption method reference data
CREATE TABLE consumption_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    category VARCHAR(50),
    description TEXT,
    onset_time_minutes INTEGER,
    duration_hours DECIMAL(3,1),
    bioavailability_percent INTEGER,
    harm_reduction_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Row Level Security Implementation

### Privacy-First Security Policies

CannaNote implements comprehensive Row Level Security to ensure user data privacy and granular access control.

```sql
-- Enable Row Level Security on personal data tables
ALTER TABLE humans ENABLE ROW LEVEL SECURITY;
ALTER TABLE entries ENABLE ROW LEVEL SECURITY;

-- Users can only access their own profile data
CREATE POLICY "Users can view own profile" ON humans
    FOR ALL USING (auth.uid() = id);

-- Entry visibility policies based on user preferences
CREATE POLICY "Users can view own entries" ON entries
    FOR ALL USING (auth.uid() = human_id);

CREATE POLICY "Community can view public entries" ON entries
    FOR SELECT USING (
        visibility = 'community' AND 
        EXISTS (
            SELECT 1 FROM humans 
            WHERE id = entries.human_id 
            AND (profile->>'community_participation')::boolean = true
        )
    );

-- Reference data is publicly readable
CREATE POLICY "Reference data is public" ON cannabinoids FOR SELECT TO authenticated USING (true);
CREATE POLICY "Reference data is public" ON terpenes FOR SELECT TO authenticated USING (true);
CREATE POLICY "Reference data is public" ON consumption_methods FOR SELECT TO authenticated USING (true);
```

### Data Synchronization Policies

```sql
-- Optional community features with privacy controls
CREATE TABLE user_follows (
    follower_id UUID REFERENCES humans(id) ON DELETE CASCADE,
    following_id UUID REFERENCES humans(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id),
    CHECK (follower_id != following_id)
);

-- Followers can view followed users' community-visible entries
CREATE POLICY "Followers can view followed entries" ON entries
    FOR SELECT USING (
        visibility IN ('community', 'followers') AND
        EXISTS (
            SELECT 1 FROM user_follows 
            WHERE follower_id = auth.uid() 
            AND following_id = entries.human_id
        )
    );
```

## Development Environment Setup

### Prerequisites

- **Supabase CLI** — Latest version for local development and deployment
- **PostgreSQL 15+** — Compatible with Supabase cloud platform
- **Node.js 18+** — Required for Supabase CLI and type generation
- **Git** — Version control for schema migrations

### Local Development Setup

#### 1. Supabase CLI Installation

```bash
# Install Supabase CLI globally
npm install -g supabase

# Verify installation
supabase --version

# Login to Supabase account
supabase login
```

#### 2. Project Initialization

```bash
# Navigate to supabase directory
cd cannanote/supabase

# Initialize local Supabase project (if needed)
supabase init

# Link to remote Supabase project
supabase link --project-ref citdskdmralncvjyybin
```

#### 3. Local Database Setup

```bash
# Start local Supabase development environment
supabase start

# This will start:
# - PostgreSQL database on localhost:54322
# - Supabase Studio on http://localhost:54323
# - API Gateway on http://localhost:54321
# - Auth server on http://localhost:54324

# Apply database schema and seed data
supabase db reset

# Generate TypeScript types from schema
supabase gen types typescript --local > types/database.types.ts
```

### Environment Configuration

```bash
# Local development environment variables
SUPABASE_URL=http://localhost:54321
SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
SUPABASE_SERVICE_ROLE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Production environment variables
SUPABASE_URL=https://citdskdmralncvjyybin.supabase.co
SUPABASE_ANON_KEY=your_production_anon_key
SUPABASE_SERVICE_ROLE_KEY=your_production_service_role_key
```

### Development Commands

```bash
# Database Management
supabase start                       # Start local development environment
supabase stop                        # Stop local environment
supabase status                      # Check service status
supabase db reset                    # Reset database with fresh schema/data
supabase db push                     # Push schema changes to remote
supabase db pull                     # Pull schema changes from remote

# Migration Management
supabase migration new migration_name # Create new migration file
supabase migration up                # Apply pending migrations
supabase migration repair             # Repair migration state

# Type Generation
supabase gen types typescript --local > types/database.types.ts
supabase gen types typescript --project-id citdskdmralncvjyybin > types/database.types.ts

# Function Deployment
supabase functions new function_name # Create new edge function
supabase functions deploy function_name # Deploy function to remote

# Project Management
supabase projects list              # List accessible projects
supabase link --project-ref <ref>  # Link to different project
supabase unlink                     # Unlink current project
```

## Reference Data Management

### Scientific Cannabis Data

The database includes comprehensive reference data for cannabinoids and terpenes based on peer-reviewed research and scientific literature.

#### Cannabinoid Data Structure

```sql
-- Example cannabinoid reference data
INSERT INTO cannabinoids (name, full_name, description, molecular_formula, common_effects, therapeutic_potential) VALUES
('THC', 'Δ9-Tetrahydrocannabinol', 'Primary psychoactive compound in cannabis', 'C21H30O2', 
 '["euphoria", "relaxation", "altered_perception", "increased_appetite"]',
 '["pain_relief", "appetite_stimulation", "anti_nausea", "muscle_relaxation"]'),
('CBD', 'Cannabidiol', 'Non-psychoactive cannabinoid with therapeutic properties', 'C21H30O2',
 '["calm", "clarity", "reduced_anxiety", "anti_inflammatory"]',
 '["anxiety_relief", "anti_inflammatory", "seizure_control", "neuroprotective"]'),
('CBG', 'Cannabigerol', 'Non-psychoactive cannabinoid, precursor to other cannabinoids', 'C21H32O2',
 '["alertness", "focus", "mild_euphoria"]',
 '["antibacterial", "anti_inflammatory", "appetite_stimulation"]');
```

#### Terpene Data Structure

```sql
-- Example terpene reference data
INSERT INTO terpenes (name, common_name, aroma_profile, flavor_notes, common_effects, boiling_point, also_found_in) VALUES
('Myrcene', 'β-Myrcene', '["earthy", "musky", "herbal"]', '["clove", "cardamom", "balsamic"]',
 '["sedating", "muscle_relaxing", "enhanced_cannabinoid_absorption"]', 166.7, 
 '["mangoes", "hops", "bay_leaves", "eucalyptus"]'),
('Limonene', 'D-Limonene', '["citrus", "lemon", "orange"]', '["lemon_zest", "orange_peel"]',
 '["mood_elevating", "stress_relief", "alertness", "anti_anxiety"]', 176.0,
 '["citrus_peels", "juniper", "peppermint", "fennel"]'),
('Pinene', 'α-Pinene', '["pine", "fresh", "sharp"]', '["pine_needles", "rosemary"]',
 '["alertness", "memory_retention", "counteracts_thc_memory_effects"]', 156.0,
 '["pine_trees", "rosemary", "basil", "dill"]');
```

### Reference Data Updates

```bash
# Update cannabinoid reference data
psql --host=localhost --port=54322 --username=postgres --dbname=postgres \
     --file=reference-data/cannabinoids.sql

# Update terpene reference data
psql --host=localhost --port=54322 --username=postgres --dbname=postgres \
     --file=reference-data/terpenes.sql

# Production reference data updates
supabase db push --include-all
```

## Performance Optimization

### Indexing Strategy

```sql
-- Primary performance indexes for cannabis tracking queries
CREATE INDEX idx_entries_human_id_consumed_at ON entries (human_id, consumed_at DESC);
CREATE INDEX idx_entries_rating ON entries (rating) WHERE rating IS NOT NULL;
CREATE INDEX idx_entries_consumption_method ON entries (consumption_method_id);
CREATE INDEX idx_entries_chemical_profile ON entries (chemical_profile_id) WHERE chemical_profile_id IS NOT NULL;

-- Chemical profile optimization for pattern recognition
CREATE INDEX idx_chemical_profiles_cannabinoids_thc ON chemical_profiles ((cannabinoids->>'THC')) WHERE cannabinoids ? 'THC';
CREATE INDEX idx_chemical_profiles_cannabinoids_cbd ON chemical_profiles ((cannabinoids->>'CBD')) WHERE cannabinoids ? 'CBD';
CREATE INDEX idx_chemical_profiles_terpenes_myrcene ON chemical_profiles ((terpenes->>'Myrcene')) WHERE terpenes ? 'Myrcene';

-- Full-text search optimization
CREATE INDEX idx_entries_fulltext ON entries USING GIN (to_tsvector('english', product_name || ' ' || COALESCE(brand, '') || ' ' || COALESCE(notes, '')));
CREATE INDEX idx_chemical_profiles_fulltext ON chemical_profiles USING GIN (to_tsvector('english', COALESCE(product_name, '') || ' ' || COALESCE(brand, '')));
```

### Query Optimization Examples

```sql
-- Efficient cannabinoid profile queries
SELECT * FROM chemical_profiles 
WHERE cannabinoids ? 'CBD' 
AND (cannabinoids->>'CBD')::numeric > 15.0;

-- Pattern recognition for user consumption preferences
SELECT 
    cm.name as consumption_method,
    AVG(e.rating) as avg_rating,
    COUNT(*) as usage_count
FROM entries e
JOIN consumption_methods cm ON e.consumption_method_id = cm.id
WHERE e.human_id = $1 
AND e.consumed_at >= NOW() - INTERVAL '90 days'
GROUP BY cm.name
ORDER BY avg_rating DESC, usage_count DESC;

-- Chemical profile effectiveness analysis
WITH user_ratings AS (
    SELECT 
        cp.id,
        cp.cannabinoids,
        cp.terpenes,
        AVG(e.rating) as avg_rating,
        COUNT(e.id) as entry_count
    FROM chemical_profiles cp
    JOIN entries e ON cp.id = e.chemical_profile_id
    WHERE e.human_id = $1
    GROUP BY cp.id, cp.cannabinoids, cp.terpenes
    HAVING COUNT(e.id) >= 3
)
SELECT * FROM user_ratings ORDER BY avg_rating DESC LIMIT 10;
```

## Planned Feature Development

### Phase 1: Core Cannabis Tracking Infrastructure

#### Enhanced Entry Management
- **Consumption Timeline** — Detailed tracking of onset, peak, and duration
- **Environment Correlation** — Location, weather, social context tracking
- **Mood Tracking** — Before/after mood and energy level measurements
- **Tolerance Management** — T-break tracking and tolerance reset indicators
- **Dosage Precision** — Micro-dosing and precise measurement tracking

#### Chemical Profile Enhancement
- **Lab Data Integration** — Direct API connections with testing laboratories
- **QR Code Support** — Product verification through manufacturer QR codes
- **Batch Tracking** — Complete supply chain visibility from seed to sale
- **Quality Indicators** — Cultivation method, curing process, storage conditions

### Phase 2: Pattern Recognition and Analytics

#### Personal Analytics Engine
- **Consumption Pattern Analysis** — Statistical analysis of usage patterns and effects
- **Chemical Profile Optimization** — Recommendation engine based on historical ratings
- **Environment Correlation** — Analysis of external factors affecting experiences
- **Tolerance Modeling** — Predictive modeling for tolerance breaks and dosage adjustments
- **Health Integration** — Optional integration with health tracking platforms

#### Community Features (Privacy-Optional)
- **Anonymous Insights** — Aggregated pattern sharing without personal identification
- **Strain Effectiveness** — Community-driven effectiveness ratings by condition
- **Research Participation** — Opt-in anonymous data contribution for cannabis research
- **Expert Content** — Integration with cannabis researchers and medical professionals

### Phase 3: Advanced Platform Integration

#### External Data Integration
- **Dispensary APIs** — Real-time product availability and pricing
- **Cultivation Data** — Integration with grow tracking platforms
- **Research Databases** — Connection to academic cannabis research platforms
- **Health Records** — Secure integration with electronic health record systems

#### Enterprise and Research Capabilities
- **Research Platform** — Anonymized aggregate data for scientific research
- **Clinical Integration** — Healthcare provider access with patient consent
- **Compliance Tracking** — Legal compliance monitoring for medical cannabis patients
- **Data Export** — Comprehensive data portability for patient records

## Security and Compliance

### Data Protection Implementation

#### Privacy Controls
- **Granular Permissions** — Fine-grained control over data sharing and visibility
- **Data Minimization** — Collect only necessary data for core functionality
- **Consent Management** — Clear, revocable consent for all data collection
- **Right to Deletion** — Complete data removal capabilities with audit trail

#### Database Security
- **Row Level Security** — Database-level access control for all personal data
- **Encryption at Rest** — Database encryption for sensitive user information
- **Audit Logging** — Comprehensive audit trail for all data access and modifications
- **Access Controls** — Role-based access control for administrative functions

#### Compliance Considerations
- **GDPR Compliance** — European data protection regulation compliance
- **CCPA Compliance** — California Consumer Privacy Act compliance
- **HIPAA Considerations** — Health information privacy for medical cannabis patients
- **Cannabis Regulations** — Compliance with state and federal cannabis regulations

## Monitoring and Maintenance

### Database Health Monitoring

```sql
-- Database performance monitoring queries
SELECT 
    schemaname,
    tablename,
    n_live_tup as row_count,
    n_dead_tup as dead_rows,
    last_vacuum,
    last_autovacuum
FROM pg_stat_user_tables
WHERE schemaname = 'public'
ORDER BY n_live_tup DESC;

-- Index usage analysis
SELECT 
    indexrelname as index_name,
    relname as table_name,
    idx_tup_read,
    idx_tup_fetch,
    idx_scan
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;

-- Query performance analysis
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    stddev_time
FROM pg_stat_statements
WHERE query LIKE '%entries%' OR query LIKE '%chemical_profiles%'
ORDER BY total_time DESC;
```

### Maintenance Procedures

```bash
# Regular database maintenance
supabase db vacuum                   # Database cleanup and optimization
supabase db analyze                  # Update table statistics
supabase db reindex                  # Rebuild indexes for optimal performance

# Backup and recovery procedures
supabase db dump --file backup.sql  # Create database backup
supabase db restore --file backup.sql # Restore from backup

# Migration management
supabase migration repair            # Repair migration state if needed
supabase migration squash            # Squash multiple migrations into one
```

## Testing and Validation

### Database Testing Strategy

#### Schema Testing
```sql
-- Test Row Level Security policies
BEGIN;
SET ROLE authenticated;
SET request.jwt.claim.sub = '11111111-1111-1111-1111-111111111111';

-- Should return user's own entries only
SELECT COUNT(*) FROM entries; 

-- Should fail for other user's private entries
SELECT COUNT(*) FROM entries WHERE human_id = '22222222-2222-2222-2222-222222222222';
ROLLBACK;
```

#### Performance Testing
```sql
-- Test JSONB query performance
EXPLAIN ANALYZE 
SELECT * FROM chemical_profiles 
WHERE cannabinoids ? 'THC' 
AND (cannabinoids->>'THC')::numeric > 20;

-- Test join query performance
EXPLAIN ANALYZE
SELECT e.*, cp.cannabinoids 
FROM entries e 
JOIN chemical_profiles cp ON e.chemical_profile_id = cp.id
WHERE e.human_id = $1 
AND e.consumed_at >= NOW() - INTERVAL '30 days';
```

#### Data Integrity Testing
```sql
-- Test referential integrity constraints
INSERT INTO entries (human_id, product_name, consumed_at) 
VALUES ('non-existent-uuid', 'Test Product', NOW()); -- Should fail

-- Test check constraints
INSERT INTO entries (human_id, product_name, rating, consumed_at) 
VALUES (gen_random_uuid(), 'Test Product', 11, NOW()); -- Should fail (rating > 10)
```

## Contributing Guidelines

### Database Development Standards

#### Schema Changes
- **Migration-Based** — All schema changes through versioned migration files
- **Backward Compatible** — Ensure changes don't break existing functionality
- **Performance Tested** — Analyze impact of indexes and constraints
- **Security Reviewed** — Validate Row Level Security policy updates

#### Data Quality
- **Referential Integrity** — Maintain proper foreign key relationships
- **Constraint Validation** — Implement appropriate check constraints
- **Data Types** — Use appropriate PostgreSQL data types for performance
- **JSONB Standards** — Follow consistent structure for JSONB columns

### Development Workflow

#### Schema Development
```bash
# Create new migration for schema changes
supabase migration new add_new_feature

# Edit migration file with schema changes
# supabase/migrations/timestamp_add_new_feature.sql

# Test migration locally
supabase db reset

# Generate updated TypeScript types
supabase gen types typescript --local > types/database.types.ts

# Commit changes with descriptive message
git add . && git commit -m "feat: add new feature schema"
```

#### Code Review Process
- **Schema Review** — Validate migration files and RLS policies
- **Performance Review** — Analyze query performance and index usage
- **Security Review** — Verify data access controls and privacy protection
- **Documentation Review** — Ensure comprehensive documentation updates

This database architecture provides the foundation for CannaNote's privacy-first cannabis journaling platform, emphasizing user data ownership, efficient storage, and scalable design for personal wellness tracking with optional community features.