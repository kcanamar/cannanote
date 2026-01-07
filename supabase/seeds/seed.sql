-- CannaNote Database Schema & Seed Data
-- This file defines our application schema and initial data
-- Applied via: supabase db reset (local) or supabase db push (production)

-- ============================================================================
-- HUMANS TABLE - Core entity for our hexagonal architecture
-- ============================================================================

-- Drop existing table if exists (for reset scenarios)
DROP TABLE IF EXISTS humans CASCADE;

-- Create humans table matching our Go hexagonal architecture
CREATE TABLE humans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL, 
    profile JSONB NOT NULL DEFAULT '{}',
    consent JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- PERFORMANCE INDEXES
-- ============================================================================

-- Index for authentication lookups
CREATE INDEX idx_humans_email ON humans(email);

-- Index for user searches
CREATE INDEX idx_humans_username ON humans(username);

-- Index for recent users
CREATE INDEX idx_humans_created_at ON humans(created_at);

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

-- Enable RLS for security
ALTER TABLE humans ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only view their own records
-- Note: This assumes Supabase auth.uid() matches the human.id
CREATE POLICY "humans_select_own" ON humans
    FOR SELECT USING (auth.uid()::text = id::text);

-- Policy: Users can only insert their own records  
CREATE POLICY "humans_insert_own" ON humans
    FOR INSERT WITH CHECK (auth.uid()::text = id::text);

-- Policy: Users can only update their own records
CREATE POLICY "humans_update_own" ON humans
    FOR UPDATE USING (auth.uid()::text = id::text);

-- Policy: Users can only delete their own records
CREATE POLICY "humans_delete_own" ON humans
    FOR DELETE USING (auth.uid()::text = id::text);

-- ============================================================================
-- FUNCTIONS & TRIGGERS
-- ============================================================================

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to auto-update updated_at on humans table
CREATE TRIGGER update_humans_updated_at 
    BEFORE UPDATE ON humans 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- SAMPLE DATA (Optional - for development/testing)
-- ============================================================================

-- Uncomment for development seed data
-- INSERT INTO humans (id, username, email, profile, consent) VALUES 
-- (
--     gen_random_uuid(),
--     'demo_user',
--     'demo@cannanote.com',
--     '{
--         "preferred_strains": ["Blue Dream", "OG Kush"],
--         "preferred_consumption": ["vape", "edible"],
--         "experience_level": "intermediate",
--         "public_profile": false,
--         "share_experiences": true
--     }'::jsonb,
--     '{
--         "data_collection": true,
--         "marketing_emails": false,
--         "data_sharing": false,
--         "medical_data_sharing": false,
--         "consent_date": "2025-12-20T00:00:00Z",
--         "consent_version": "v1.0"
--     }'::jsonb
-- );

-- ============================================================================
-- CANNABINOIDS TABLE - Reference data for cannabinoid compounds
-- ============================================================================

DROP TABLE IF EXISTS cannabinoids CASCADE;

CREATE TABLE cannabinoids (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    full_name TEXT NOT NULL,
    description TEXT NOT NULL,
    psychoactive BOOLEAN NOT NULL DEFAULT false,
    reported_experiences JSONB NOT NULL DEFAULT '{}',
    compound_notes JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- TERPENES TABLE - Reference data for terpene compounds  
-- ============================================================================

DROP TABLE IF EXISTS terpenes CASCADE;

CREATE TABLE terpenes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    aroma_profile JSONB NOT NULL DEFAULT '{}',
    reported_effects JSONB NOT NULL DEFAULT '{}',
    boiling_point_celsius DECIMAL(5,2),
    common_sources JSONB NOT NULL DEFAULT '{}',
    research_notes JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- CONSUMPTION_METHODS TABLE - Reference data for consumption methods
-- ============================================================================

DROP TABLE IF EXISTS consumption_methods CASCADE;

CREATE TABLE consumption_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    category TEXT NOT NULL, -- 'inhalation', 'oral', 'topical', 'sublingual'
    description TEXT NOT NULL,
    onset_time_minutes INTEGER,
    duration_hours INTEGER,
    bioavailability_percent INTEGER,
    notes JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- CHEMICAL_PROFILES TABLE - Lab test data with JSONB efficiency
-- ============================================================================

DROP TABLE IF EXISTS chemical_profiles CASCADE;

CREATE TABLE chemical_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_name TEXT,
    brand TEXT,
    lab_name TEXT,
    batch_id TEXT,
    garden_name TEXT,
    harvest_date DATE,
    test_date DATE,
    cannabinoids JSONB NOT NULL DEFAULT '{}', -- {'THC': 23.5, 'CBD': 0.8, 'CBG': 1.2}
    terpenes JSONB NOT NULL DEFAULT '{}', -- {'Myrcene': 1.2, 'Limonene': 0.8}
    raw_data JSONB NOT NULL DEFAULT '{}', -- Full lab report if available
    total_cannabinoids DECIMAL(5,2),
    total_terpenes DECIMAL(5,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- ENTRIES TABLE - Cannabis consumption experiences
-- ============================================================================

DROP TABLE IF EXISTS entries CASCADE;

CREATE TABLE entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    human_id UUID NOT NULL REFERENCES humans(id) ON DELETE CASCADE,
    chemical_profile_id UUID REFERENCES chemical_profiles(id) ON DELETE SET NULL,
    consumption_method_id UUID NOT NULL REFERENCES consumption_methods(id),
    product_name TEXT NOT NULL,
    brand TEXT,
    dose_amount DECIMAL(8,3) NOT NULL,
    dose_unit TEXT NOT NULL, -- 'mg', 'g', 'ml', 'puffs', etc.
    rating INTEGER CHECK (rating >= 1 AND rating <= 10),
    effects TEXT, -- Free text description of effects experienced
    notes TEXT, -- Additional notes
    mood_before INTEGER CHECK (mood_before >= 1 AND mood_before <= 10),
    mood_after INTEGER CHECK (mood_after >= 1 AND mood_after <= 10),
    consumed_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- DATA_SOURCES TABLE - Track reference data provenance for legal compliance
-- ============================================================================

DROP TABLE IF EXISTS data_sources CASCADE;

CREATE TABLE data_sources (
    table_name TEXT PRIMARY KEY,
    source TEXT NOT NULL,
    last_updated TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    legal_notes TEXT NOT NULL,
    version TEXT DEFAULT 'v1.0'
);

-- ============================================================================
-- PERFORMANCE INDEXES FOR ALL TABLES
-- ============================================================================

-- Humans indexes (already created above)

-- Cannabinoids indexes
CREATE INDEX idx_cannabinoids_name ON cannabinoids(name);
CREATE INDEX idx_cannabinoids_psychoactive ON cannabinoids(psychoactive);
CREATE INDEX idx_cannabinoids_experiences_gin ON cannabinoids USING gin(reported_experiences);

-- Terpenes indexes  
CREATE INDEX idx_terpenes_name ON terpenes(name);
CREATE INDEX idx_terpenes_aroma_gin ON terpenes USING gin(aroma_profile);
CREATE INDEX idx_terpenes_effects_gin ON terpenes USING gin(reported_effects);
CREATE INDEX idx_terpenes_boiling_point ON terpenes(boiling_point_celsius);

-- Consumption methods indexes
CREATE INDEX idx_consumption_methods_name ON consumption_methods(name);
CREATE INDEX idx_consumption_methods_category ON consumption_methods(category);

-- Chemical profiles indexes
CREATE INDEX idx_chemical_profiles_product_name ON chemical_profiles(product_name);
CREATE INDEX idx_chemical_profiles_brand ON chemical_profiles(brand);
CREATE INDEX idx_chemical_profiles_lab_name ON chemical_profiles(lab_name);
CREATE INDEX idx_chemical_profiles_cannabinoids_gin ON chemical_profiles USING gin(cannabinoids);
CREATE INDEX idx_chemical_profiles_terpenes_gin ON chemical_profiles USING gin(terpenes);
CREATE INDEX idx_chemical_profiles_total_cannabinoids ON chemical_profiles(total_cannabinoids);
CREATE INDEX idx_chemical_profiles_total_terpenes ON chemical_profiles(total_terpenes);

-- Entries indexes
CREATE INDEX idx_entries_human_id ON entries(human_id);
CREATE INDEX idx_entries_chemical_profile_id ON entries(chemical_profile_id);
CREATE INDEX idx_entries_consumption_method_id ON entries(consumption_method_id);
CREATE INDEX idx_entries_consumed_at ON entries(consumed_at);
CREATE INDEX idx_entries_rating ON entries(rating);
CREATE INDEX idx_entries_created_at ON entries(created_at);

-- ============================================================================
-- ROW LEVEL SECURITY POLICIES
-- ============================================================================

-- Chemical profiles: Public reference data (no RLS needed)
-- Cannabinoids: Public reference data (no RLS needed)  
-- Terpenes: Public reference data (no RLS needed)
-- Consumption methods: Public reference data (no RLS needed)

-- Entries: Private user data with RLS
ALTER TABLE entries ENABLE ROW LEVEL SECURITY;

CREATE POLICY "entries_select_own" ON entries
    FOR SELECT USING (auth.uid()::text = human_id::text);

CREATE POLICY "entries_insert_own" ON entries
    FOR INSERT WITH CHECK (auth.uid()::text = human_id::text);

CREATE POLICY "entries_update_own" ON entries  
    FOR UPDATE USING (auth.uid()::text = human_id::text);

CREATE POLICY "entries_delete_own" ON entries
    FOR DELETE USING (auth.uid()::text = human_id::text);

-- ============================================================================
-- TRIGGERS FOR UPDATED_AT TIMESTAMPS
-- ============================================================================

-- Function to automatically update updated_at timestamp (already created above)

-- Apply trigger to all tables with updated_at column
CREATE TRIGGER update_cannabinoids_updated_at 
    BEFORE UPDATE ON cannabinoids 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_terpenes_updated_at 
    BEFORE UPDATE ON terpenes 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_consumption_methods_updated_at 
    BEFORE UPDATE ON consumption_methods 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_chemical_profiles_updated_at 
    BEFORE UPDATE ON chemical_profiles 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_entries_updated_at 
    BEFORE UPDATE ON entries 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- COMMENTS & DOCUMENTATION
-- ============================================================================

COMMENT ON TABLE humans IS 'Core human entities for CannaNote cannabis journaling platform';
COMMENT ON COLUMN humans.id IS 'Primary key UUID, matches Supabase auth.uid() for RLS';
COMMENT ON COLUMN humans.username IS 'Unique human-friendly identifier';
COMMENT ON COLUMN humans.email IS 'Unique email for authentication';
COMMENT ON COLUMN humans.profile IS 'JSON object containing cannabis preferences and profile data';
COMMENT ON COLUMN humans.consent IS 'JSON object containing HIPAA-ready consent settings and history';

COMMENT ON TABLE cannabinoids IS 'Reference data for cannabinoid compounds - educational purposes only';
COMMENT ON TABLE terpenes IS 'Reference data for terpene compounds - educational purposes only';
COMMENT ON TABLE consumption_methods IS 'Reference data for cannabis consumption methods';
COMMENT ON TABLE chemical_profiles IS 'Lab test data using efficient JSONB storage for chemical analysis';
COMMENT ON TABLE entries IS 'Personal cannabis consumption experiences and effects tracking';
COMMENT ON TABLE data_sources IS 'Legal compliance tracking for reference data provenance';