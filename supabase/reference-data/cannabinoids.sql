-- ============================================================================
-- CANNABINOIDS REFERENCE DATA - EDUCATIONAL PURPOSES ONLY
-- 
-- IMPORTANT LEGAL DISCLAIMERS:
-- * This data is for educational and research purposes only
-- * NOT intended as medical advice, diagnosis, or treatment recommendations  
-- * Consult qualified healthcare professionals for medical guidance
-- * Individual experiences may vary significantly
-- * No health claims are made or implied
-- ============================================================================

-- Clear existing data for fresh update
TRUNCATE TABLE cannabinoids RESTART IDENTITY CASCADE;

-- Insert cannabinoid compound data for educational reference
INSERT INTO cannabinoids (id, name, full_name, description, psychoactive, reported_experiences, compound_notes) VALUES

-- THC - Most well-known cannabis compound
(
    gen_random_uuid(),
    'THC',
    'Tetrahydrocannabinol',
    'The most widely recognized compound in cannabis. Commonly associated with psychoactive effects. Subject of ongoing scientific research.',
    true,
    '{
        "commonly_reported": ["euphoria", "relaxation", "altered_perception"],
        "user_reported_uses": ["recreational", "wellness", "research"],
        "anecdotal_experiences": [
            "appetite_changes", "relaxation", "mood_effects", 
            "sensory_changes", "time_perception_changes"
        ]
    }'::jsonb,
    '{
        "abbreviation": "THC",
        "research_status": "extensively_studied",
        "legal_status": "varies_by_jurisdiction"
    }'::jsonb
),

-- CBD - Second most studied compound
(
    gen_random_uuid(),
    'CBD',
    'Cannabidiol',
    'Widely studied non-psychoactive compound. Subject of ongoing research and commercial interest. Found in many consumer products.',
    false,
    '{
        "commonly_reported": ["relaxation", "wellness_support", "daily_use"],
        "user_reported_uses": ["wellness", "dietary_supplement", "research"],
        "anecdotal_experiences": [
            "relaxation", "daily_wellness_routine", "research_interest"
        ]
    }'::jsonb,
    '{
        "abbreviation": "CBD",
        "research_status": "heavily_researched",
        "availability": "widely_available_in_products"
    }'::jsonb
),

-- CBN - Sleep-associated compound
(
    gen_random_uuid(),
    'CBN',
    'Cannabinol',
    'Compound often found in aged cannabis. Commonly associated with evening use and relaxation in user reports. Typically found in small amounts.',
    false,
    '{
        "commonly_reported": ["evening_use", "relaxation", "sleep_routine"],
        "user_reported_uses": ["evening_wellness", "research"],
        "anecdotal_experiences": [
            "bedtime_routine", "relaxation", "evening_wind_down"
        ]
    }'::jsonb,
    '{
        "abbreviation": "CBN",
        "formation": "results_from_thc_degradation_over_time",
        "typical_amounts": "usually_less_than_1_percent"
    }'::jsonb
),

-- THCA - Raw plant compound
(
    gen_random_uuid(),
    'THCA',
    'Tetrahydrocannabinolic Acid',
    'Compound found in fresh cannabis plants. Converts to THC when heated. Subject of ongoing research into raw cannabis consumption.',
    false,
    '{
        "commonly_reported": ["raw_consumption", "juicing", "research_interest"],
        "user_reported_uses": ["raw_cannabis_products", "research"],
        "anecdotal_experiences": [
            "wellness_routines", "raw_plant_consumption", "research_participation"
        ]
    }'::jsonb,
    '{
        "abbreviation": "THCA",
        "conversion": "becomes_thc_when_heated",
        "consumption": "often_consumed_raw"
    }'::jsonb
),

-- THCV - Appetite-related compound
(
    gen_random_uuid(),
    'THCV',
    'Tetrahydrocannabivarin',
    'Compound with reported appetite effects different from THC. Subject of research into metabolic processes and appetite regulation.',
    true,
    '{
        "commonly_reported": ["appetite_effects", "energy", "focus"],
        "user_reported_uses": ["wellness_research", "appetite_management"],
        "anecdotal_experiences": [
            "different_appetite_effects", "energy_effects", "research_participation"
        ]
    }'::jsonb,
    '{
        "abbreviation": "THCV",
        "unique_property": "different_appetite_effects_than_thc",
        "research_focus": "metabolic_studies"
    }'::jsonb
),

-- CBG - Precursor compound
(
    gen_random_uuid(),
    'CBG',
    'Cannabigerol',
    'Often called the "mother cannabinoid" as other compounds develop from it. Subject of increasing research and commercial interest.',
    false,
    '{
        "commonly_reported": ["wellness_interest", "research_participation", "daily_use"],
        "user_reported_uses": ["wellness_routines", "research"],
        "anecdotal_experiences": [
            "daily_wellness", "research_interest", "product_exploration"
        ]
    }'::jsonb,
    '{
        "abbreviation": "CBG",
        "nickname": "mother_cannabinoid",
        "research_status": "increasing_commercial_interest"
    }'::jsonb
),

-- CBC - Emerging research compound
(
    gen_random_uuid(),
    'CBC',
    'Cannabichromene',
    'Less well-known compound that is gaining research attention. Often found alongside other cannabinoids in full-spectrum products.',
    false,
    '{
        "commonly_reported": ["wellness_interest", "research_curiosity", "product_component"],
        "user_reported_uses": ["full_spectrum_products", "research"],
        "anecdotal_experiences": [
            "wellness_exploration", "research_participation", "product_interest"
        ]
    }'::jsonb,
    '{
        "abbreviation": "CBC",
        "research_status": "emerging_interest",
        "occurrence": "often_found_with_other_compounds"
    }'::jsonb
),

-- CBDA - Raw CBD compound
(
    gen_random_uuid(),
    'CBDA',
    'Cannabidiolic Acid',
    'Precursor to CBD found in raw cannabis plants. Converts to CBD when heated. Interest in raw cannabis consumption methods.',
    false,
    '{
        "commonly_reported": ["raw_consumption", "juicing", "wellness_exploration"],
        "user_reported_uses": ["raw_products", "research"],
        "anecdotal_experiences": [
            "raw_plant_consumption", "wellness_routines", "product_exploration"
        ]
    }'::jsonb,
    '{
        "abbreviation": "CBDA",
        "conversion": "becomes_cbd_when_heated",
        "consumption": "raw_cannabis_methods"
    }'::jsonb
),

-- CBDV - Specialized research compound
(
    gen_random_uuid(),
    'CBDV',
    'Cannabidivarin',
    'Compound that has gained attention in specialized research. Similar to CBD but with structural differences that interest researchers.',
    false,
    '{
        "commonly_reported": ["research_interest", "specialized_studies", "wellness_exploration"],
        "user_reported_uses": ["research_participation", "specialized_products"],
        "anecdotal_experiences": [
            "research_participation", "wellness_interest", "product_curiosity"
        ]
    }'::jsonb,
    '{
        "abbreviation": "CBDV",
        "research_focus": "specialized_studies",
        "relationship": "similar_to_cbd_with_differences"
    }'::jsonb
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_cannabinoids_name ON cannabinoids(name);
CREATE INDEX IF NOT EXISTS idx_cannabinoids_psychoactive ON cannabinoids(psychoactive);
CREATE INDEX IF NOT EXISTS idx_cannabinoids_experiences_gin ON cannabinoids USING gin(reported_experiences);

-- ============================================================================
-- LEGAL COMPLIANCE DOCUMENTATION
-- ============================================================================

COMMENT ON TABLE cannabinoids IS 'Educational reference data only - not medical advice';
COMMENT ON COLUMN cannabinoids.reported_experiences IS 'Anecdotal user reports - not medical claims';
COMMENT ON COLUMN cannabinoids.compound_notes IS 'Educational information for research purposes';

-- Track data provenance for legal compliance
INSERT INTO data_sources (table_name, source, last_updated, legal_notes) VALUES 
('cannabinoids', 'Educational compilation for research purposes', NOW(), 
 'No medical claims made. Educational use only. Consult healthcare professionals.')
ON CONFLICT (table_name) DO UPDATE SET 
    source = EXCLUDED.source,
    last_updated = EXCLUDED.last_updated,
    legal_notes = EXCLUDED.legal_notes;