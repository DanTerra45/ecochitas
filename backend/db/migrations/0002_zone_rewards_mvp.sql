CREATE TABLE IF NOT EXISTS recycling_zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_code TEXT NOT NULL UNIQUE,
    zone_name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS zone_households (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_id UUID NOT NULL REFERENCES recycling_zones(id) ON DELETE CASCADE,
    household_code TEXT NOT NULL UNIQUE,
    street_address TEXT NOT NULL,
    house_reference TEXT,
    representative_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS zone_recycling_containers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_id UUID NOT NULL REFERENCES recycling_zones(id) ON DELETE CASCADE,
    container_code TEXT NOT NULL UNIQUE,
    latitude NUMERIC(10, 7) NOT NULL,
    longitude NUMERIC(10, 7) NOT NULL,
    container_status TEXT NOT NULL DEFAULT 'available' CHECK (container_status IN ('available', 'warning', 'full', 'maintenance')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS zone_recycling_container_household_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    container_id UUID NOT NULL REFERENCES zone_recycling_containers(id) ON DELETE CASCADE,
    household_id UUID NOT NULL REFERENCES zone_households(id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (container_id, household_id),
    UNIQUE (household_id)
);

CREATE TABLE IF NOT EXISTS recycling_collection_cycles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    container_id UUID NOT NULL REFERENCES zone_recycling_containers(id) ON DELETE CASCADE,
    scheduled_collection_date DATE NOT NULL,
    cycle_status TEXT NOT NULL DEFAULT 'scheduled' CHECK (cycle_status IN ('scheduled', 'in_progress', 'closed', 'reprogrammed', 'cancelled')),
    collection_operator_name TEXT,
    raw_points_total NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (raw_points_total >= 0),
    contamination_level TEXT NOT NULL DEFAULT 'low' CHECK (contamination_level IN ('low', 'medium', 'high')),
    contamination_discount_percentage NUMERIC(5, 2) NOT NULL DEFAULT 0 CHECK (contamination_discount_percentage >= 0 AND contamination_discount_percentage <= 100),
    discount_points_total NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (discount_points_total >= 0),
    final_points_total NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (final_points_total >= 0),
    contamination_notes TEXT,
    started_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS recycling_cycle_evidence_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cycle_id UUID NOT NULL REFERENCES recycling_collection_cycles(id) ON DELETE CASCADE,
    household_id UUID NOT NULL REFERENCES zone_households(id) ON DELETE CASCADE,
    evidence_photo_url TEXT NOT NULL,
    evidence_captured_at TIMESTAMPTZ NOT NULL,
    evidence_latitude NUMERIC(10, 7) NOT NULL,
    evidence_longitude NUMERIC(10, 7) NOT NULL,
    validation_status TEXT NOT NULL DEFAULT 'accepted' CHECK (validation_status IN ('pending', 'accepted', 'rejected')),
    rejection_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (cycle_id, household_id)
);

CREATE TABLE IF NOT EXISTS recycling_cycle_household_points (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cycle_id UUID NOT NULL REFERENCES recycling_collection_cycles(id) ON DELETE CASCADE,
    household_id UUID NOT NULL REFERENCES zone_households(id) ON DELETE CASCADE,
    awarded_points NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (awarded_points >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (cycle_id, household_id)
);

DROP INDEX IF EXISTS recycling_collection_cycles_container_id_scheduled_collection_d;
DROP INDEX IF EXISTS recycling_cycle_evidence_submissions_cycle_id_validation_status;

CREATE INDEX IF NOT EXISTS zone_households_zone_id_idx
    ON zone_households (zone_id);

CREATE INDEX IF NOT EXISTS zone_recycling_containers_zone_id_idx
    ON zone_recycling_containers (zone_id);

CREATE INDEX IF NOT EXISTS recycling_cycles_container_schedule_idx
    ON recycling_collection_cycles (container_id, scheduled_collection_date DESC);

CREATE INDEX IF NOT EXISTS recycling_evidence_cycle_status_idx
    ON recycling_cycle_evidence_submissions (cycle_id, validation_status);
