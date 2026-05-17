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

require_command -command_name psql

$postgres_port_number = if ([string]::IsNullOrWhiteSpace($env:PGPORT)) { 5432 } else { [int]$env:PGPORT }
$postgres_connection_string = "postgresql://ecochitas_user:ecochitas_password@localhost:$postgres_port_number/ecochitas_db?sslmode=disable"

$zone_identifier = execute_scalar_query -postgres_connection_string $postgres_connection_string -sql_text @"
INSERT INTO recycling_zones (zone_code, zone_name)
VALUES ('ZONA-NORTE', 'Zona Norte')
ON CONFLICT (zone_code)
DO UPDATE SET zone_name = EXCLUDED.zone_name
RETURNING id;
"@

$container_identifier = execute_scalar_query -postgres_connection_string $postgres_connection_string -sql_text @"
INSERT INTO zone_recycling_containers (zone_id, container_code, latitude, longitude)
VALUES ('$zone_identifier'::uuid, 'RC-NORTE-001', -17.3920, -66.1560)
ON CONFLICT (container_code)
DO UPDATE SET
	zone_id = EXCLUDED.zone_id,
	latitude = EXCLUDED.latitude,
	longitude = EXCLUDED.longitude
RETURNING id;
"@

$household_identifier_list = @()
$household_seed_list = @(
	@{ household_code = 'HOGAR-NORTE-001'; street_address = 'Av. Blanco Galindo #100'; house_reference = 'Casa 1' },
	@{ household_code = 'HOGAR-NORTE-002'; street_address = 'Av. Blanco Galindo #102'; house_reference = 'Casa 2' }
)

foreach ($household_seed_item in $household_seed_list) {
	$household_identifier = execute_scalar_query -postgres_connection_string $postgres_connection_string -sql_text @"
INSERT INTO zone_households (zone_id, household_code, street_address, house_reference)
VALUES (
	'$zone_identifier'::uuid,
	'$($household_seed_item.household_code)',
	'$($household_seed_item.street_address)',
	'$($household_seed_item.house_reference)'
)
ON CONFLICT (household_code)
DO UPDATE SET
	zone_id = EXCLUDED.zone_id,
	street_address = EXCLUDED.street_address,
	house_reference = EXCLUDED.house_reference
RETURNING id;
"@

	$household_identifier_list += $household_identifier
	psql $postgres_connection_string -qAt -v ON_ERROR_STOP=1 -c @"
INSERT INTO zone_recycling_container_household_assignments (container_id, household_id)
VALUES ('$container_identifier'::uuid, '$household_identifier'::uuid)
ON CONFLICT (household_id)
DO UPDATE SET
	container_id = EXCLUDED.container_id,
	is_active = TRUE;
"@ | Out-Null
}

$seed_output_payload = [ordered]@{
	zone_identifier = $zone_identifier
	container_identifier = $container_identifier
	household_identifiers = $household_identifier_list
}

$seed_output_payload | ConvertTo-Json -Depth 4
