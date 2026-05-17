CREATE TABLE IF NOT EXISTS truck_route_assignments (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	truck_identifier TEXT NOT NULL,
	route_id UUID NOT NULL REFERENCES collection_routes(id) ON DELETE CASCADE,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	assigned_by_user_identifier TEXT,
	assignment_notes TEXT,
	assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	unassigned_at TIMESTAMPTZ,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS truck_route_assignments_truck_active_idx
	ON truck_route_assignments (truck_identifier)
	WHERE is_active = TRUE;

CREATE INDEX IF NOT EXISTS truck_route_assignments_route_active_idx
	ON truck_route_assignments (route_id, is_active);

CREATE TABLE IF NOT EXISTS route_deviation_alert_events (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	truck_identifier TEXT NOT NULL,
	route_id UUID NOT NULL REFERENCES collection_routes(id) ON DELETE CASCADE,
	route_stop_id UUID REFERENCES route_stops(id) ON DELETE SET NULL,
	distance_to_route_meters NUMERIC(10, 2) NOT NULL CHECK (distance_to_route_meters >= 0),
	deviation_threshold_meters NUMERIC(10, 2) NOT NULL CHECK (deviation_threshold_meters > 0),
	severity_level TEXT NOT NULL CHECK (severity_level IN ('low', 'medium', 'high')),
	alert_status TEXT NOT NULL DEFAULT 'open' CHECK (alert_status IN ('open', 'resolved', 'dismissed')),
	alert_notes TEXT,
	triggered_by_user_identifier TEXT,
	detected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	resolved_at TIMESTAMPTZ,
	metadata_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS route_deviation_alert_events_route_detected_idx
	ON route_deviation_alert_events (route_id, detected_at DESC);

CREATE INDEX IF NOT EXISTS route_deviation_alert_events_truck_status_idx
	ON route_deviation_alert_events (truck_identifier, alert_status, detected_at DESC);

CREATE TABLE IF NOT EXISTS notification_outbox_events (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	event_type TEXT NOT NULL,
	channel_type TEXT NOT NULL CHECK (channel_type IN ('push', 'email', 'sms', 'webhook')),
	target_reference TEXT NOT NULL,
	title_text TEXT NOT NULL,
	body_text TEXT NOT NULL,
	payload JSONB NOT NULL DEFAULT '{}'::jsonb,
	delivery_status TEXT NOT NULL DEFAULT 'pending' CHECK (delivery_status IN ('pending', 'sent', 'failed', 'cancelled')),
	failure_reason TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	sent_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS notification_outbox_events_status_created_idx
	ON notification_outbox_events (delivery_status, created_at ASC);
