CREATE TABLE IF NOT EXISTS collection_route_revisions (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	route_id UUID NOT NULL REFERENCES collection_routes(id) ON DELETE CASCADE,
	revision_number INTEGER NOT NULL CHECK (revision_number > 0),
	change_type TEXT NOT NULL CHECK (change_type IN ('route_created', 'route_updated', 'route_stops_synced')),
	changed_by_user_identifier TEXT,
	change_payload JSONB NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE (route_id, revision_number)
);

CREATE INDEX IF NOT EXISTS collection_route_revisions_route_created_idx
	ON collection_route_revisions (route_id, created_at DESC);
