param(
	[switch]$force_ports
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function write_info_message {
	param([string]$message_text)
	Write-Host "[poc] $message_text" -ForegroundColor Yellow
}

function write_success_message {
	param([string]$message_text)
	Write-Host "[poc] $message_text" -ForegroundColor Green
}

function stop_process_if_running {
	param(
		[string]$service_name,
		[object]$service_state
	)

	if (-not $service_state) {
		write_info_message "$service_name has no state entry."
		return
	}

	$process_id = [int]$service_state.pid
	$started_by_script = [bool]$service_state.started_by_script
	if (-not $started_by_script -and -not $force_ports) {
		write_info_message "$service_name is running but was not started by script (pid=$process_id). Skipping."
		return
	}

	if ($process_id -le 0) {
		write_info_message "$service_name has invalid pid in state."
		return
	}

	$running_process = Get-Process -Id $process_id -ErrorAction SilentlyContinue
	if (-not $running_process) {
		write_info_message "$service_name process already stopped (pid=$process_id)."
		return
	}

	Stop-Process -Id $process_id -Force
	write_success_message "$service_name stopped (pid=$process_id)."
}

$script_directory = Split-Path -Parent $MyInvocation.MyCommand.Path
$backend_root_directory = (Resolve-Path (Join-Path $script_directory '..')).Path
$processes_state_file_path = Join-Path (Join-Path $backend_root_directory '.poc') 'processes.json'

if (-not (Test-Path $processes_state_file_path)) {
	write_info_message "State file not found: $processes_state_file_path"
	write_info_message 'Nothing to stop from script-managed services.'
	exit 0
}

$processes_state_payload = Get-Content $processes_state_file_path -Raw | ConvertFrom-Json

# Stop in reverse dependency order.
stop_process_if_running -service_name 'api' -service_state $processes_state_payload.api
stop_process_if_running -service_name 'nats' -service_state $processes_state_payload.nats
stop_process_if_running -service_name 'postgres' -service_state $processes_state_payload.postgres

Remove-Item $processes_state_file_path -Force
write_success_message 'State file removed.'
