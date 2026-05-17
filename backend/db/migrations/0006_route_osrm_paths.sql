ALTER TABLE collection_routes
	ADD COLUMN IF NOT EXISTS stops_hash TEXT,
	ADD COLUMN IF NOT EXISTS road_path_coordinates JSONB,
	ADD COLUMN IF NOT EXISTS routed_at TIMESTAMPTZ,
	ADD COLUMN IF NOT EXISTS routing_provider TEXT,
	ADD COLUMN IF NOT EXISTS routing_status TEXT,
	ADD COLUMN IF NOT EXISTS routing_error TEXT;

UPDATE collection_routes
SET road_path_coordinates = '[]'::jsonb
WHERE road_path_coordinates IS NULL;
