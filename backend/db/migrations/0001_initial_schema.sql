CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL,
    role_name TEXT NOT NULL CHECK (role_name IN ('citizen', 'driver', 'admin', 'condominium_admin')),
    zone_name TEXT,
    home_reference TEXT,
    accumulated_points NUMERIC(12, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS condominiums (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    condominium_name TEXT NOT NULL UNIQUE,
    zone_name TEXT NOT NULL,
    monthly_target_points NUMERIC(12, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS condominium_memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    condominium_id UUID NOT NULL REFERENCES condominiums(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (condominium_id, user_id)
);

CREATE TABLE IF NOT EXISTS bins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bin_code TEXT NOT NULL UNIQUE,
    zone_name TEXT NOT NULL,
    latitude NUMERIC(10, 7) NOT NULL,
    longitude NUMERIC(10, 7) NOT NULL,
    fill_percentage INTEGER NOT NULL DEFAULT 0 CHECK (fill_percentage >= 0 AND fill_percentage <= 100),
    sensor_status TEXT NOT NULL DEFAULT 'online' CHECK (sensor_status IN ('online', 'offline')),
    bin_status TEXT NOT NULL DEFAULT 'unknown' CHECK (bin_status IN ('available', 'warning', 'full', 'unknown')),
    last_emptied_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS collection_routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_code TEXT NOT NULL UNIQUE,
    route_name TEXT NOT NULL,
    zone_name TEXT NOT NULL,
    collection_weekday SMALLINT NOT NULL CHECK (collection_weekday >= 1 AND collection_weekday <= 7),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS route_stops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id UUID NOT NULL REFERENCES collection_routes(id) ON DELETE CASCADE,
    bin_id UUID NOT NULL REFERENCES bins(id) ON DELETE CASCADE,
    stop_order INTEGER NOT NULL CHECK (stop_order > 0),
    planned_time TIME,
    UNIQUE (route_id, stop_order),
    UNIQUE (route_id, bin_id)
);

CREATE TABLE IF NOT EXISTS driver_collection_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_stop_id UUID REFERENCES route_stops(id) ON DELETE SET NULL,
    bin_id UUID NOT NULL REFERENCES bins(id) ON DELETE CASCADE,
    driver_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    action_type TEXT NOT NULL CHECK (action_type IN ('emptied', 'not_accessible', 'contaminated')),
    evidence_photo_url TEXT,
    action_notes TEXT,
    action_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS truck_positions (
    id BIGSERIAL PRIMARY KEY,
    truck_identifier TEXT NOT NULL,
    latitude NUMERIC(10, 7) NOT NULL,
    longitude NUMERIC(10, 7) NOT NULL,
    speed_kmh NUMERIC(6, 2),
    heading_degrees NUMERIC(6, 2),
    captured_at TIMESTAMPTZ NOT NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (truck_identifier, captured_at)
);

CREATE INDEX IF NOT EXISTS truck_positions_truck_identifier_captured_at_idx
    ON truck_positions (truck_identifier, captured_at DESC);

CREATE TABLE IF NOT EXISTS recycling_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_name TEXT NOT NULL UNIQUE,
    category_description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS recycling_materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES recycling_categories(id) ON DELETE CASCADE,
    material_name TEXT NOT NULL,
    material_description TEXT,
    points_per_kg NUMERIC(8, 2) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (category_id, material_name)
);

CREATE TABLE IF NOT EXISTS recycling_price_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    material_id UUID NOT NULL REFERENCES recycling_materials(id) ON DELETE CASCADE,
    price_bs_per_kg NUMERIC(8, 2) NOT NULL CHECK (price_bs_per_kg >= 0),
    effective_from TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS recycling_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    material_id UUID NOT NULL REFERENCES recycling_materials(id) ON DELETE RESTRICT,
    submitted_weight_kg NUMERIC(10, 3) NOT NULL CHECK (submitted_weight_kg > 0),
    awarded_points NUMERIC(10, 2) NOT NULL CHECK (awarded_points >= 0),
    estimated_sale_bs NUMERIC(10, 2) NOT NULL CHECK (estimated_sale_bs >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS illegal_dumping_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reporter_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    zone_name TEXT NOT NULL,
    latitude NUMERIC(10, 7) NOT NULL,
    longitude NUMERIC(10, 7) NOT NULL,
    media_url TEXT NOT NULL,
    report_description TEXT,
    report_status TEXT NOT NULL DEFAULT 'received' CHECK (report_status IN ('received', 'under_review', 'validated', 'discarded', 'resolved')),
    suspected_person_fingerprint TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS weekly_offender_highlights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES illegal_dumping_reports(id) ON DELETE CASCADE,
    week_start_date DATE NOT NULL,
    publication_status TEXT NOT NULL DEFAULT 'pending' CHECK (publication_status IN ('pending', 'published', 'rejected')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (report_id, week_start_date)
);

