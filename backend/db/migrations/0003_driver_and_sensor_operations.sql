ALTER TABLE driver_collection_events
	ALTER COLUMN driver_user_id DROP NOT NULL;

ALTER TABLE driver_collection_events
	ADD COLUMN IF NOT EXISTS driver_user_identifier TEXT;

CREATE INDEX IF NOT EXISTS driver_collection_events_bin_action_at_idx
	ON driver_collection_events (bin_id, action_at DESC);

CREATE TABLE IF NOT EXISTS bin_sensor_events (
	id BIGSERIAL PRIMARY KEY,
	bin_id UUID NOT NULL REFERENCES bins(id) ON DELETE CASCADE,
	fill_percentage INTEGER NOT NULL CHECK (fill_percentage >= 0 AND fill_percentage <= 100),
	sensor_status TEXT NOT NULL CHECK (sensor_status IN ('online', 'offline')),
	measured_at TIMESTAMPTZ NOT NULL,
	source_identifier TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS bin_sensor_events_bin_measured_at_idx
	ON bin_sensor_events (bin_id, measured_at DESC);

CREATE TABLE IF NOT EXISTS collection_route_blockage_events (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	route_id UUID REFERENCES collection_routes(id) ON DELETE SET NULL,
	route_stop_id UUID REFERENCES route_stops(id) ON DELETE SET NULL,
	bin_id UUID REFERENCES bins(id) ON DELETE SET NULL,
	reported_by_user_identifier TEXT NOT NULL,
	blockage_reason TEXT NOT NULL,
	evidence_photo_url TEXT,
	severity_level TEXT NOT NULL DEFAULT 'medium' CHECK (severity_level IN ('low', 'medium', 'high')),
	reported_at TIMESTAMPTZ NOT NULL,
	status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'resolved', 'dismissed')),
	resolved_at TIMESTAMPTZ,
	resolution_notes TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS collection_route_blockage_events_route_status_idx
	ON collection_route_blockage_events (route_id, status, reported_at DESC);
