Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function require_command {
	param([string]$command_name)
	$resolved_command = Get-Command $command_name -ErrorAction SilentlyContinue
	if (-not $resolved_command) {
		throw "required_command_not_found: $command_name"
	}
}

function execute_scalar_query {
	param(
		[string]$postgres_connection_string,
		[string]$sql_text
	)
	return ((psql $postgres_connection_string -qAt -v ON_ERROR_STOP=1 -c $sql_text | Out-String).Trim())
}

function execute_non_query {
	param(
		[string]$postgres_connection_string,
		[string]$sql_text
	)
	psql $postgres_connection_string -qAt -v ON_ERROR_STOP=1 -c $sql_text | Out-Null
}

require_command -command_name psql

$postgres_port_number = if ([string]::IsNullOrWhiteSpace($env:PGPORT)) { 5432 } else { [int]$env:PGPORT }
$postgres_connection_string = "postgresql://ecochitas_user:ecochitas_password@localhost:$postgres_port_number/ecochitas_db?sslmode=disable"

$bin_demo_001_identifier = execute_scalar_query -postgres_connection_string $postgres_connection_string -sql_text @"
INSERT INTO bins (bin_code, zone_name, latitude, longitude, fill_percentage, sensor_status, bin_status)
VALUES ('BIN-DEMO-001', 'Cochabamba Centro', -17.3935000, -66.1571000, 20, 'online', 'available')
ON CONFLICT (bin_code)
DO UPDATE SET
	zone_name = EXCLUDED.zone_name,
	latitude = EXCLUDED.latitude,
	longitude = EXCLUDED.longitude,
	fill_percentage = EXCLUDED.fill_percentage,
	sensor_status = EXCLUDED.sensor_status,
	bin_status = EXCLUDED.bin_status,
	updated_at = NOW()
RETURNING id;
"@

$bin_demo_002_identifier = execute_scalar_query -postgres_connection_string $postgres_connection_string -sql_text @"
INSERT INTO bins (bin_code, zone_name, latitude, longitude, fill_percentage, sensor_status, bin_status)
VALUES ('BIN-DEMO-002', 'Cochabamba Centro', -17.3917000, -66.1548000, 40, 'online', 'available')
ON CONFLICT (bin_code)
DO UPDATE SET
	zone_name = EXCLUDED.zone_name,
	latitude = EXCLUDED.latitude,
	longitude = EXCLUDED.longitude,
	fill_percentage = EXCLUDED.fill_percentage,
	sensor_status = EXCLUDED.sensor_status,
	bin_status = EXCLUDED.bin_status,
	updated_at = NOW()
RETURNING id;
"@

$bin_demo_003_identifier = execute_scalar_query -postgres_connection_string $postgres_connection_string -sql_text @"
INSERT INTO bins (bin_code, zone_name, latitude, longitude, fill_percentage, sensor_status, bin_status)
VALUES ('BIN-DEMO-003', 'Cochabamba Centro', -17.3899000, -66.1527000, 60, 'online', 'warning')
ON CONFLICT (bin_code)
DO UPDATE SET
	zone_name = EXCLUDED.zone_name,
	latitude = EXCLUDED.latitude,
	longitude = EXCLUDED.longitude,
	fill_percentage = EXCLUDED.fill_percentage,
	sensor_status = EXCLUDED.sensor_status,
	bin_status = EXCLUDED.bin_status,
	updated_at = NOW()
RETURNING id;
"@

$route_demo_001_identifier = execute_scalar_query -postgres_connection_string $postgres_connection_string -sql_text @"
INSERT INTO collection_routes (route_code, route_name, zone_name, collection_weekday, is_active)
VALUES ('RUTA-DEMO-001', 'Ruta Centro Demo', 'Cochabamba Centro', 2, TRUE)
ON CONFLICT (route_code)
DO UPDATE SET
	route_name = EXCLUDED.route_name,
	zone_name = EXCLUDED.zone_name,
	collection_weekday = EXCLUDED.collection_weekday,
	is_active = EXCLUDED.is_active,
	updated_at = NOW()
RETURNING id;
"@

execute_non_query -postgres_connection_string $postgres_connection_string -sql_text @"
DELETE FROM route_stops
WHERE route_id = '$route_demo_001_identifier'::uuid;
"@

execute_non_query -postgres_connection_string $postgres_connection_string -sql_text @"
INSERT INTO route_stops (route_id, bin_id, stop_order, planned_time)
VALUES
	('$route_demo_001_identifier'::uuid, '$bin_demo_001_identifier'::uuid, 1, '08:00'),
	('$route_demo_001_identifier'::uuid, '$bin_demo_002_identifier'::uuid, 2, '08:20'),
	('$route_demo_001_identifier'::uuid, '$bin_demo_003_identifier'::uuid, 3, '08:40');
"@

execute_non_query -postgres_connection_string $postgres_connection_string -sql_text @"
UPDATE truck_route_assignments
SET
	is_active = FALSE,
	unassigned_at = NOW(),
	updated_at = NOW()
WHERE truck_identifier = 'TRUCK-001'
	AND is_active = TRUE;
"@

$truck_assignment_identifier = execute_scalar_query -postgres_connection_string $postgres_connection_string -sql_text @"
INSERT INTO truck_route_assignments (
	truck_identifier,
	route_id,
	is_active,
	assigned_by_user_identifier,
	assignment_notes,
	assigned_at
)
VALUES (
	'TRUCK-001',
	'$route_demo_001_identifier'::uuid,
	TRUE,
	'admin-demo-001',
	'seed assignment for integration tests',
	NOW()
)
RETURNING id;
"@

execute_non_query -postgres_connection_string $postgres_connection_string -sql_text @"
INSERT INTO truck_positions (truck_identifier, latitude, longitude, speed_kmh, heading_degrees, captured_at)
VALUES
	('TRUCK-001', -17.3934000, -66.1570000, 18.50, 90.0, NOW() - INTERVAL '2 minutes'),
	('TRUCK-001', -17.3812000, -66.1400000, 24.00, 120.0, NOW() - INTERVAL '1 minutes')
ON CONFLICT (truck_identifier, captured_at)
DO NOTHING;
"@

$seed_output_payload = [ordered]@{
	route_identifier = $route_demo_001_identifier
	truck_identifier = 'TRUCK-001'
	truck_assignment_identifier = $truck_assignment_identifier
	bin_identifiers = @(
		$bin_demo_001_identifier,
		$bin_demo_002_identifier,
		$bin_demo_003_identifier
	)
	next_tests = @(
		"GET /v1/collection-routes?zone_name=Cochabamba%20Centro",
		"GET /v1/driver/route-deviation?truck_identifier=TRUCK-001",
		"POST /v1/driver/route-deviation-alerts"
	)
}

$seed_output_payload | ConvertTo-Json -Depth 5
