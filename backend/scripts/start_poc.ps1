param(
	[switch]$skip_migration
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $false

function write_info_message {
	param([string]$message_text)
	Write-Host "[poc] $message_text" -ForegroundColor Cyan
}

function write_success_message {
	param([string]$message_text)
	Write-Host "[poc] $message_text" -ForegroundColor Green
}

function ensure_command_available {
	param([string]$command_name)
	$resolved_command = Get-Command $command_name -ErrorAction SilentlyContinue
	if (-not $resolved_command) {
		throw "required_command_not_found: $command_name"
	}
}

function get_listener_process_id {
	param([int]$port_number)
	$tcp_listener = Get-NetTCPConnection -State Listen -LocalPort $port_number -ErrorAction SilentlyContinue | Select-Object -First 1
	if ($tcp_listener) {
		return [int]$tcp_listener.OwningProcess
	}
	return 0
}

function wait_for_port_listener {
	param(
		[int]$port_number,
		[int]$timeout_seconds
	)

	$timeout_deadline = (Get-Date).AddSeconds($timeout_seconds)
	while ((Get-Date) -lt $timeout_deadline) {
		$listener_process_id = get_listener_process_id -port_number $port_number
		if ($listener_process_id -gt 0) {
			return $true
		}
		Start-Sleep -Milliseconds 300
	}

	return $false
}

function wait_for_postgres_sql_ready {
	param(
		[int]$port_number,
		[int]$timeout_seconds
	)

	$timeout_deadline = (Get-Date).AddSeconds($timeout_seconds)
	while ((Get-Date) -lt $timeout_deadline) {
		$previous_error_action_preference = $ErrorActionPreference
		$ErrorActionPreference = 'Continue'
		try {
			$null = & psql -h localhost -p $port_number -U postgres -d postgres -tAc "SELECT 1;" 2>$null
		} catch {
		} finally {
			$ErrorActionPreference = $previous_error_action_preference
		}
		if ($LASTEXITCODE -eq 0) {
			return $true
		}
		Start-Sleep -Milliseconds 600
	}

	return $false
}

function get_nats_server_binary_path {
	$nats_command = Get-Command nats-server -ErrorAction SilentlyContinue
	if ($nats_command) {
		return $nats_command.Source
	}

	$go_workspace_path = (& go env GOPATH).Trim()
	if ([string]::IsNullOrWhiteSpace($go_workspace_path)) {
		throw 'go_env_gopath_is_empty'
	}

	$nats_binary_path = Join-Path $go_workspace_path 'bin\nats-server.exe'
	if (-not (Test-Path $nats_binary_path)) {
		write_info_message "Installing nats-server into $go_workspace_path\bin ..."
		& go install github.com/nats-io/nats-server/v2@v2.11.11
		if ($LASTEXITCODE -ne 0) {
			throw 'failed_to_install_nats_server'
		}
	}

	if (-not (Test-Path $nats_binary_path)) {
		throw "nats_server_binary_not_found_after_install: $nats_binary_path"
	}

	return $nats_binary_path
}

function ensure_database_role_and_database {
	param([int]$postgres_port_number)

	$role_exists_query = "SELECT 1 FROM pg_roles WHERE rolname = 'ecochitas_user';"
	$role_exists_raw = & psql -h localhost -p $postgres_port_number -U postgres -d postgres -tAc $role_exists_query 2>$null
	$role_exists_value = ($role_exists_raw | Out-String).Trim()
	if ($role_exists_value -ne '1') {
		write_info_message 'Creating role ecochitas_user ...'
		& psql -h localhost -p $postgres_port_number -U postgres -d postgres -v ON_ERROR_STOP=1 -c "CREATE ROLE ecochitas_user LOGIN PASSWORD 'ecochitas_password';" | Out-Null
		if ($LASTEXITCODE -ne 0) {
			throw 'failed_to_create_role_ecochitas_user'
		}
	}

	$database_exists_query = "SELECT 1 FROM pg_database WHERE datname = 'ecochitas_db';"
	$database_exists_raw = & psql -h localhost -p $postgres_port_number -U postgres -d postgres -tAc $database_exists_query 2>$null
	$database_exists_value = ($database_exists_raw | Out-String).Trim()
	if ($database_exists_value -ne '1') {
		write_info_message 'Creating database ecochitas_db ...'
		& psql -h localhost -p $postgres_port_number -U postgres -d postgres -v ON_ERROR_STOP=1 -c "CREATE DATABASE ecochitas_db OWNER ecochitas_user;" | Out-Null
		if ($LASTEXITCODE -ne 0) {
			throw 'failed_to_create_database_ecochitas_db'
		}
	}
}

function run_database_migration {
	param(
		[int]$postgres_port_number,
		[string]$backend_root_directory
	)

	$migrations_directory_path = Join-Path $backend_root_directory 'db\migrations'
	if (-not (Test-Path $migrations_directory_path)) {
		throw "migration_directory_not_found: $migrations_directory_path"
	}
	$migration_file_list = Get-ChildItem -Path $migrations_directory_path -Filter '*.sql' | Sort-Object Name
	if ($migration_file_list.Count -eq 0) {
		throw "migration_files_not_found_in_directory: $migrations_directory_path"
	}

	$postgres_connection_string = "postgresql://ecochitas_user:ecochitas_password@localhost:$postgres_port_number/ecochitas_db?sslmode=disable"
	foreach ($migration_file in $migration_file_list) {
		& psql $postgres_connection_string -q -v ON_ERROR_STOP=1 -f $migration_file.FullName | Out-Null
		if ($LASTEXITCODE -ne 0) {
			throw "failed_to_apply_database_migration_file: $($migration_file.Name)"
		}
	}
}

function wait_for_backend_health_endpoint {
	param(
		[string]$health_endpoint_url,
		[int]$timeout_seconds
	)

	$timeout_deadline = (Get-Date).AddSeconds($timeout_seconds)
	while ((Get-Date) -lt $timeout_deadline) {
		try {
			$health_response = Invoke-RestMethod -Uri $health_endpoint_url -TimeoutSec 2
			if ($health_response.status -eq 'ok') {
				return $true
			}
		} catch {
		}
		Start-Sleep -Milliseconds 500
	}

	return $false
}

$script_directory = Split-Path -Parent $MyInvocation.MyCommand.Path
$backend_root_directory = (Resolve-Path (Join-Path $script_directory '..')).Path
$processes_state_directory = Join-Path $backend_root_directory '.poc'
$processes_state_file_path = Join-Path $processes_state_directory 'processes.json'

ensure_command_available -command_name go
ensure_command_available -command_name initdb
ensure_command_available -command_name postgres
ensure_command_available -command_name psql

if ([string]::IsNullOrWhiteSpace($env:PGDATA)) {
	throw 'environment_variable_pgdata_is_required'
}

$postgres_port_number = if ([string]::IsNullOrWhiteSpace($env:PGPORT)) { 5432 } else { [int]$env:PGPORT }
$postgres_data_directory = $env:PGDATA
$nats_port_number = 4222
$backend_port_number = 8080

New-Item -ItemType Directory -Force $processes_state_directory | Out-Null

Push-Location $backend_root_directory
try {
	if (-not (Test-Path (Join-Path $postgres_data_directory 'PG_VERSION'))) {
		write_info_message "Initializing PostgreSQL cluster at $postgres_data_directory ..."
		& initdb -D $postgres_data_directory -U postgres -A trust -E UTF8 | Out-Null
		if ($LASTEXITCODE -ne 0) {
			throw 'failed_to_initialize_postgres_cluster'
		}
	}

	$postgres_process_id = get_listener_process_id -port_number $postgres_port_number
	$postgres_started_by_script = $false
	if ($postgres_process_id -le 0) {
		write_info_message "Starting PostgreSQL on port $postgres_port_number ..."
		$postgres_process = Start-Process -FilePath 'postgres' -ArgumentList @('-D', $postgres_data_directory, '-p', "$postgres_port_number") -WindowStyle Hidden -PassThru
		$postgres_process_id = $postgres_process.Id
		$postgres_started_by_script = $true
	}

	if (-not (wait_for_port_listener -port_number $postgres_port_number -timeout_seconds 20)) {
		throw "postgres_listener_not_ready_on_port_$postgres_port_number"
	}
	if (-not (wait_for_postgres_sql_ready -port_number $postgres_port_number -timeout_seconds 60)) {
		throw 'postgres_sql_not_ready'
	}
	write_success_message "PostgreSQL ready (pid=$postgres_process_id, port=$postgres_port_number)"

	ensure_database_role_and_database -postgres_port_number $postgres_port_number
	if (-not $skip_migration) {
		run_database_migration -postgres_port_number $postgres_port_number -backend_root_directory $backend_root_directory
		write_success_message 'Database migration applied'
	}

	$nats_server_binary_path = get_nats_server_binary_path
	$nats_process_id = get_listener_process_id -port_number $nats_port_number
	$nats_started_by_script = $false
	if ($nats_process_id -le 0) {
		write_info_message "Starting NATS on port $nats_port_number ..."
		$nats_process = Start-Process -FilePath $nats_server_binary_path -ArgumentList @('-p', "$nats_port_number") -WindowStyle Hidden -PassThru
		$nats_process_id = $nats_process.Id
		$nats_started_by_script = $true
	}
	if (-not (wait_for_port_listener -port_number $nats_port_number -timeout_seconds 20)) {
		throw "nats_listener_not_ready_on_port_$nats_port_number"
	}
	write_success_message "NATS ready (pid=$nats_process_id, port=$nats_port_number)"

	$backend_process_id = get_listener_process_id -port_number $backend_port_number
	$backend_started_by_script = $false
	if ($backend_process_id -le 0) {
		write_info_message "Starting backend API on port $backend_port_number ..."
		$env:POSTGRES_URL = "postgresql://ecochitas_user:ecochitas_password@localhost:$postgres_port_number/ecochitas_db?sslmode=disable"
		$env:NATS_URL = "nats://localhost:$nats_port_number"
		$env:NATS_GPS_SUBJECT = 'gps.trucks.location'
		$env:NATS_GPS_QUEUE_GROUP = 'ecochitas_gps_workers'
		$env:NATS_ROUTE_DEVIATION_ALERTS_SUBJECT = 'operations.routes.deviation_alerts'
		$env:HTTP_WRITE_TIMEOUT_SECONDS = '0'
		$backend_process = Start-Process -FilePath 'go' -ArgumentList @('run', './cmd/api') -WorkingDirectory $backend_root_directory -WindowStyle Hidden -PassThru
		$backend_process_id = $backend_process.Id
		$backend_started_by_script = $true
	}

	if (-not (wait_for_backend_health_endpoint -health_endpoint_url "http://127.0.0.1:$backend_port_number/healthz" -timeout_seconds 45)) {
		throw 'backend_health_endpoint_not_ready'
	}
	write_success_message "Backend API ready (pid=$backend_process_id, port=$backend_port_number)"

	$processes_state_payload = @{
		generated_at = (Get-Date).ToString('o')
		backend_root = $backend_root_directory
		postgres = @{
			pid = $postgres_process_id
			started_by_script = $postgres_started_by_script
			port = $postgres_port_number
			data_directory = $postgres_data_directory
		}
		nats = @{
			pid = $nats_process_id
			started_by_script = $nats_started_by_script
			port = $nats_port_number
			binary_path = $nats_server_binary_path
		}
		api = @{
			pid = $backend_process_id
			started_by_script = $backend_started_by_script
			port = $backend_port_number
		}
	}
	$processes_state_payload | ConvertTo-Json -Depth 8 | Set-Content -Path $processes_state_file_path -Encoding UTF8

	write_success_message "PoC services are running. State file: $processes_state_file_path"
	write_info_message "Quick test: go run ./cmd/gps_event_publisher"
	write_info_message "API URL: http://127.0.0.1:$backend_port_number"
}
finally {
	Pop-Location
}
