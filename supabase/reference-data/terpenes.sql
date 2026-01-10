-- ============================================================================
-- TERPENES REFERENCE DATA - EDUCATIONAL PURPOSES ONLY
-- 
-- IMPORTANT LEGAL DISCLAIMERS:
-- * This data is for educational and research purposes only
-- * NOT intended as medical advice, diagnosis, or treatment recommendations  
-- * Consult qualified healthcare professionals for medical guidance
-- * Individual experiences may vary significantly
-- * No health claims are made or implied
-- * Based on 2025 cannabis research and user reports
-- ============================================================================

-- Clear existing data for fresh update
TRUNCATE TABLE terpenes RESTART IDENTITY CASCADE;

-- Insert comprehensive terpene compound data for educational reference
INSERT INTO terpenes (id, name, aroma_profile, reported_effects, boiling_point_celsius, common_sources, research_notes) VALUES

-- Myrcene - Most abundant terpene in cannabis
(
    gen_random_uuid(),
    'Myrcene',
    '{"primary_aroma": "earthy", "secondary_notes": ["musky", "clove-like", "herbal"], "intensity": "moderate_to_strong", "descriptors": ["fruity undertones", "hoppy", "spicy"]}',
    '{
        "commonly_reported": ["relaxation", "sedation", "muscle_relief"],
        "user_experiences": ["nighttime_use", "stress_relief", "body_relaxation"],
        "anecdotal_effects": [
            "enhanced_thc_absorption", "couch_lock_sensation", 
            "sleep_support", "muscle_tension_relief"
        ],
        "timing_preference": "evening_nighttime"
    }'::jsonb,
    166.0,
    '{"botanical_sources": ["mangoes", "hops", "thyme", "lemongrass", "bay_leaves"], "prevalence": "most_abundant_in_cannabis"}'::jsonb,
    '{
        "significance_threshold": "0.5_percent_indicates_sedating_effects",
        "research_status": "extensively_studied",
        "entourage_effect": "enhances_thc_psychoactivity",
        "mechanism": "reported_to_affect_muscle_relaxation"
    }'::jsonb
),

-- Limonene - Second most common, citrus profile
(
    gen_random_uuid(),
    'Limonene',
    '{"primary_aroma": "citrus", "secondary_notes": ["lemon_zest", "orange_peel", "fresh"], "intensity": "strong", "descriptors": ["bright", "clean", "energizing"]}',
    '{
        "commonly_reported": ["mood_elevation", "stress_relief", "anxiety_reduction"],
        "user_experiences": ["daytime_use", "social_activities", "creative_tasks"],
        "anecdotal_effects": [
            "counteracts_thc_anxiety", "uplifting_mood", 
            "enhanced_focus", "anti_stress"
        ],
        "timing_preference": "daytime_morning"
    }'::jsonb,
    176.0,
    '{"botanical_sources": ["citrus_fruits", "juniper", "peppermint", "rosemary"], "prevalence": "second_most_common_in_cannabis"}'::jsonb,
    '{
        "therapeutic_research": "anti_anxiety_anti_inflammatory",
        "mechanism": "affects_serotonin_dopamine_signaling", 
        "absorption": "well_absorbed_through_skin",
        "safety_profile": "generally_recognized_as_safe"
    }'::jsonb
),

-- Beta-Caryophyllene - Unique CB2 receptor interaction
(
    gen_random_uuid(),
    'Beta-Caryophyllene',
    '{"primary_aroma": "spicy", "secondary_notes": ["peppery", "woody", "cloves"], "intensity": "strong", "descriptors": ["sharp", "warm", "pungent"]}',
    '{
        "commonly_reported": ["anti_inflammatory", "pain_relief", "stress_reduction"],
        "user_experiences": ["wellness_routines", "physical_discomfort", "research_interest"],
        "anecdotal_effects": [
            "muscle_soreness_relief", "inflammation_support", 
            "tension_reduction", "physical_wellness"
        ],
        "timing_preference": "any_time"
    }'::jsonb,
    199.0,
    '{"botanical_sources": ["black_pepper", "cloves", "cinnamon", "oregano", "basil"], "prevalence": "common_in_many_strains"}'::jsonb,
    '{
        "unique_property": "only_terpene_that_directly_activates_cb2_receptors",
        "classification": "dietary_cannabinoid",
        "research_focus": "anti_inflammatory_analgesic_properties",
        "mechanism": "cb2_receptor_agonist"
    }'::jsonb
),

-- Linalool - Floral relaxation compound
(
    gen_random_uuid(),
    'Linalool',
    '{"primary_aroma": "floral", "secondary_notes": ["lavender", "soft", "perfumed"], "intensity": "gentle", "descriptors": ["soothing", "sweet", "delicate"]}',
    '{
        "commonly_reported": ["calming", "anxiety_relief", "sleep_support"],
        "user_experiences": ["bedtime_routines", "relaxation", "stress_management"],
        "anecdotal_effects": [
            "improved_sleep_quality", "reduced_anxiety", 
            "mood_stabilization", "relaxation_response"
        ],
        "timing_preference": "evening_nighttime"
    }'::jsonb,
    198.0,
    '{"botanical_sources": ["lavender", "mint", "cinnamon", "rosewood"], "prevalence": "moderate_in_cannabis"}'::jsonb,
    '{
        "mechanism": "modulates_gaba_receptors",
        "research_applications": "sleep_disorders_anxiety_studies",
        "aromatherapy_use": "widely_used_for_relaxation",
        "safety_profile": "well_tolerated"
    }'::jsonb
),

-- Alpha-Pinene - Focus and alertness compound
(
    gen_random_uuid(),
    'Alpha-Pinene',
    '{"primary_aroma": "pine", "secondary_notes": ["forest", "fresh", "woody"], "intensity": "moderate", "descriptors": ["crisp", "natural", "invigorating"]}',
    '{
        "commonly_reported": ["alertness", "focus", "memory_support"],
        "user_experiences": ["daytime_productivity", "mental_clarity", "outdoor_activities"],
        "anecdotal_effects": [
            "enhanced_alertness", "memory_retention", 
            "reduced_thc_memory_impairment", "mental_clarity"
        ],
        "timing_preference": "daytime"
    }'::jsonb,
    156.0,
    '{"botanical_sources": ["pine_trees", "rosemary", "basil", "dill"], "prevalence": "common_in_many_strains"}'::jsonb,
    '{
        "cognitive_effects": "may_counteract_thc_memory_effects",
        "research_focus": "bronchodilator_anti_inflammatory",
        "mechanism": "acetylcholinesterase_inhibition",
        "therapeutic_potential": "respiratory_cognitive_support"
    }'::jsonb
),

-- Terpinolene - Complex uplifting profile
(
    gen_random_uuid(),
    'Terpinolene',
    '{"primary_aroma": "complex", "secondary_notes": ["floral", "herbal", "citrus"], "intensity": "moderate", "descriptors": ["fresh", "piney", "sweet"]}',
    '{
        "commonly_reported": ["uplifting", "creative", "energetic"],
        "user_experiences": ["creative_projects", "social_situations", "daytime_activities"],
        "anecdotal_effects": [
            "enhanced_creativity", "mood_elevation", 
            "social_confidence", "mental_stimulation"
        ],
        "timing_preference": "daytime"
    }'::jsonb,
    186.0,
    '{"botanical_sources": ["nutmeg", "tea_tree", "conifers", "cumin"], "prevalence": "rare_as_dominant_terpene"}'::jsonb,
    '{
        "rarity": "uncommon_as_primary_terpene",
        "research_status": "limited_studies",
        "reported_properties": "antioxidant_sedative_paradox",
        "user_reports": "uplifting_despite_sedative_classification"
    }'::jsonb
),

-- Humulene - Appetite and inflammation compound
(
    gen_random_uuid(),
    'Humulene',
    '{"primary_aroma": "woody", "secondary_notes": ["earthy", "spicy", "hoppy"], "intensity": "moderate", "descriptors": ["subtle", "herbal", "dry"]}',
    '{
        "commonly_reported": ["appetite_suppression", "anti_inflammatory", "alertness"],
        "user_experiences": ["wellness_routines", "research_interest", "herbal_products"],
        "anecdotal_effects": [
            "reduced_appetite", "inflammation_support", 
            "energy_without_stimulation", "focus_enhancement"
        ],
        "timing_preference": "any_time"
    }'::jsonb,
    106.0,
    '{"botanical_sources": ["hops", "coriander", "cloves", "sage"], "prevalence": "moderate_in_cannabis"}'::jsonb,
    '{
        "unique_property": "potential_appetite_suppressant",
        "research_applications": "weight_management_studies",
        "anti_inflammatory": "significant_research_interest",
        "cannabimimetic": "reported_cannabinoid_like_activity"
    }'::jsonb
),

-- Ocimene - Energy and immunity compound
(
    gen_random_uuid(),
    'Ocimene',
    '{"primary_aroma": "sweet", "secondary_notes": ["fruity", "citrus", "herbal"], "intensity": "light", "descriptors": ["fresh", "tropical", "bright"]}',
    '{
        "commonly_reported": ["energizing", "uplifting", "immune_support"],
        "user_experiences": ["daytime_wellness", "outdoor_activities", "social_energy"],
        "anecdotal_effects": [
            "natural_energy_boost", "mood_enhancement", 
            "immune_system_support", "antiviral_properties"
        ],
        "timing_preference": "daytime"
    }'::jsonb,
    50.0,
    '{"botanical_sources": ["mint", "parsley", "pepper", "basil"], "prevalence": "less_common_in_cannabis"}'::jsonb,
    '{
        "therapeutic_research": "antiviral_antifungal_properties",
        "immune_support": "traditional_use_in_aromatherapy",
        "volatility": "highly_volatile_low_boiling_point",
        "preservation": "degrades_quickly_requires_proper_storage"
    }'::jsonb
),

-- Eucalyptol (1,8-Cineole) - Clarity and respiratory support
(
    gen_random_uuid(),
    'Eucalyptol',
    '{"primary_aroma": "eucalyptus", "secondary_notes": ["medicinal", "cooling", "fresh"], "intensity": "strong", "descriptors": ["penetrating", "clear", "camphor-like"]}',
    '{
        "commonly_reported": ["mental_clarity", "respiratory_support", "anti_inflammatory"],
        "user_experiences": ["focus_enhancement", "breathing_wellness", "mental_alertness"],
        "anecdotal_effects": [
            "enhanced_concentration", "respiratory_comfort", 
            "reduced_inflammation", "cognitive_clarity"
        ],
        "timing_preference": "daytime"
    }'::jsonb,
    176.0,
    '{"botanical_sources": ["eucalyptus", "rosemary", "sage", "tea_tree"], "prevalence": "uncommon_in_cannabis"}'::jsonb,
    '{
        "respiratory_benefits": "traditional_use_for_breathing_support",
        "cognitive_effects": "reported_memory_enhancement",
        "anti_inflammatory": "significant_research_support",
        "mechanism": "reported_to_increase_cerebral_blood_flow"
    }'::jsonb
),

-- Bisabolol - Gentle healing compound
(
    gen_random_uuid(),
    'Bisabolol',
    '{"primary_aroma": "floral", "secondary_notes": ["chamomile", "honey", "sweet"], "intensity": "gentle", "descriptors": ["soothing", "soft", "delicate"]}',
    '{
        "commonly_reported": ["skin_soothing", "anti_inflammatory", "calming"],
        "user_experiences": ["topical_wellness", "relaxation", "skin_care"],
        "anecdotal_effects": [
            "skin_comfort", "reduced_irritation", 
            "gentle_relaxation", "healing_support"
        ],
        "timing_preference": "any_time"
    }'::jsonb,
    153.0,
    '{"botanical_sources": ["chamomile", "candeia_tree"], "prevalence": "rare_in_cannabis"}'::jsonb,
    '{
        "skin_benefits": "widely_used_in_cosmetics",
        "anti_inflammatory": "gentle_non_irritating",
        "research_applications": "dermatological_studies",
        "safety_profile": "excellent_tolerance"
    }'::jsonb
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_terpenes_name ON terpenes(name);
CREATE INDEX IF NOT EXISTS idx_terpenes_aroma_gin ON terpenes USING gin(aroma_profile);
CREATE INDEX IF NOT EXISTS idx_terpenes_effects_gin ON terpenes USING gin(reported_effects);
CREATE INDEX IF NOT EXISTS idx_terpenes_boiling_point ON terpenes(boiling_point_celsius);

-- ============================================================================
-- LEGAL COMPLIANCE DOCUMENTATION
-- ============================================================================

COMMENT ON TABLE terpenes IS 'Educational terpene reference data - not medical advice or health claims';
COMMENT ON COLUMN terpenes.aroma_profile IS 'Descriptive aroma characteristics for identification purposes';
COMMENT ON COLUMN terpenes.reported_effects IS 'User-reported experiences - not medical claims';
COMMENT ON COLUMN terpenes.research_notes IS 'Educational research information - consult professionals for medical guidance';

-- Track data provenance for legal compliance
INSERT INTO data_sources (table_name, source, last_updated, legal_notes) VALUES 
('terpenes', 'Comprehensive 2025 cannabis research compilation for educational purposes', NOW(), 
 'No medical claims made. Educational and research use only. Individual experiences vary. Consult healthcare professionals for medical guidance.')
ON CONFLICT (table_name) DO UPDATE SET 
    source = EXCLUDED.source,
    last_updated = EXCLUDED.last_updated,
    legal_notes = EXCLUDED.legal_notes;