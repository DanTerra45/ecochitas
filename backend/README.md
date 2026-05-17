# EcoChitas Backend (Go)

Backend del PoC para:

- tracking GPS de camiones en tiempo real,
- stream SSE para mapa,
- modelo MVP por zonas para bonificación de reciclaje.

## Requisitos

- Go 1.22+
- PostgreSQL (18.x recomendado)
- NATS Server 2.x
- PowerShell (scripts incluidos)

## Variables de entorno clave

- `NATS_BIN_SENSOR_EVENTS_SUBJECT` (default: `operations.bins.sensor_events`)
- `NATS_DRIVER_COLLECTION_EVENTS_SUBJECT` (default: `operations.driver.collection_events`)
- `NATS_ROUTE_BLOCKAGE_EVENTS_SUBJECT` (default: `operations.routes.blockage_events`)
- `NATS_ROUTE_DEVIATION_ALERTS_SUBJECT` (default: `operations.routes.deviation_alerts`)
- `AUTH_JWT_SIGNING_KEY` (default: `dev_only_change_me`)
- `AUTH_JWT_ISSUER` (default: `ecochitas_backend`)
- `AUTH_JWT_AUDIENCE` (default: `ecochitas_api`)
- `AUTH_ACCESS_TOKEN_TTL_MINUTES` (default: `120`)
- `AUTH_ENABLE_DEV_TOKEN_ISSUE` (default: `true`)

## Arranque rápido

```powershell
cd D:\academico_1\ecochitas\backend
.\scripts\start_poc.ps1
```

El script:

1. inicializa clúster PostgreSQL si no existe,
2. crea rol/db de PoC si falta,
3. aplica migraciones de `db/migrations/*.sql`,
4. levanta NATS,
5. levanta API en `http://127.0.0.1:8080`.

Parar servicios:

```powershell
.\scripts\stop_poc.ps1 -force_ports
```

## Comandos útiles

### API

```powershell
go run ./cmd/api
```

### Publicador de un evento GPS

```powershell
go run ./cmd/gps_event_publisher
```

### Simulador continuo de ruta GPS

```powershell
go run ./cmd/gps_route_simulator -truck_identifier TRUCK-001 -publish_interval 1s
```

### Seed de demo para zonas y hogares

```powershell
.\scripts\seed_zone_rewards_demo.ps1
```

### Seed de demo para rutas y operaciones

```powershell
.\scripts\seed_route_operations_demo.ps1
```

## Endpoints

### Salud

- `GET /healthz`

### Auth (MVP)

- `GET /v1/auth/me` (requiere `Authorization: Bearer <token>`)
- `POST /v1/auth/dev-token` (solo para desarrollo si `AUTH_ENABLE_DEV_TOKEN_ISSUE=true`)

### GPS

- `GET /v1/trucks/latest-position?truck_identifier=TRUCK-001`
- `GET /v1/trucks/latest-positions`
- `GET /v1/trucks/stream` (SSE)

### Rutas de recoleccion (mapa)

- `GET /v1/collection-routes?zone_name=Cochabamba%20Centro&collection_weekday=2&is_active=true`
- `GET /v1/collection-routes/{route_identifier}/stops`
- `GET /v1/driver/route-deviation?truck_identifier=TRUCK-001&route_identifier=<uuid>&deviation_threshold_meters=200` (roles: `driver`, `admin`)

### Administracion de rutas (admin)

- `POST /v1/admin/collection-routes` (role: `admin`)
- `PATCH /v1/admin/collection-routes/{route_identifier}` (role: `admin`)
- `PUT /v1/admin/collection-routes/{route_identifier}/stops` (role: `admin`)
- `GET /v1/admin/collection-routes/{route_identifier}/revisions?limit=20` (role: `admin`)
- `POST /v1/admin/truck-route-assignments` (role: `admin`)
- `GET /v1/admin/truck-route-assignments?truck_identifier=TRUCK-001&is_active=true` (role: `admin`)
- `GET /v1/admin/route-deviation-alerts?status=open&limit=50` (role: `admin`)

### Operaciones de campo (driver/admin)

- `POST /v1/bins/sensor-events` (roles: `driver`, `admin`)
- `POST /v1/driver/collection-events` (roles: `driver`, `admin`)
- `POST /v1/driver/route-blockages` (roles: `driver`, `admin`)
- `GET /v1/driver/route-blockages?status=open&limit=20` (roles: `driver`, `admin`)
- `PATCH /v1/driver/route-blockages/{blockage_identifier}` (roles: `driver`, `admin`)
- `POST /v1/driver/route-deviation-alerts` (roles: `driver`, `admin`)

### Zonas y puntos de reciclaje (MVP v1)

- `GET /v1/recycling/containers?zone_code=ZONA-NORTE`
- `POST /v1/recycling/cycles/start` (roles: `driver`, `admin`)
- `POST /v1/recycling/evidence-submissions` (roles: `citizen`, `driver`, `admin`, `condominium_admin`)
- `POST /v1/recycling/cycles/close` (roles: `driver`, `admin`)
- `GET /v1/recycling/cycles/summary?cycle_identifier=<uuid>`

## Flujo rápido de token en desarrollo

### 1) Emitir token

```powershell
$token_response = Invoke-RestMethod `
	-Method POST `
	-Uri "http://127.0.0.1:8080/v1/auth/dev-token" `
	-ContentType "application/json" `
	-Body '{"user_identifier":"driver-demo-001","role_name":"driver","full_name":"Driver Demo"}'

$access_token = $token_response.access_token
```

### 2) Probar endpoint protegido

```powershell
Invoke-RestMethod `
	-Method GET `
	-Uri "http://127.0.0.1:8080/v1/auth/me" `
	-Headers @{ Authorization = "Bearer $access_token" }
```

## Ejemplos payload

### Iniciar ciclo

```json
{
	"container_identifier": "e9960178-c1a8-4de6-ac89-b30f2d0356b9",
	"scheduled_collection_date": "2026-05-16",
	"collection_operator_name": "Mamut"
}
```

### Evento de sensor

```json
{
	"bin_identifier": "697d3564-8662-4f1e-9dc4-2fa5c6855a55",
	"source_identifier": "sensor-lab-001",
	"fill_percentage": 92,
	"sensor_status": "online",
	"measured_at": "2026-05-17T10:30:00Z"
}
```

### Evento de recoleccion del chofer

```json
{
	"route_stop_identifier": "",
	"bin_identifier": "697d3564-8662-4f1e-9dc4-2fa5c6855a55",
	"action_type": "emptied",
	"evidence_photo_url": "https://example.com/empty-bin-photo.jpg",
	"action_notes": "recolectado en prueba",
	"action_at": "2026-05-17T10:35:00Z"
}
```

### Reporte de bloqueo de ruta

```json
{
	"route_identifier": "",
	"route_stop_identifier": "",
	"bin_identifier": "697d3564-8662-4f1e-9dc4-2fa5c6855a55",
	"blockage_reason": "vehiculo estacionado bloqueando acceso",
	"evidence_photo_url": "https://example.com/blockage-photo.jpg",
	"severity_level": "medium",
	"reported_at": "2026-05-17T10:40:00Z"
}
```

### Consulta de bloqueos de ruta

```text
GET /v1/driver/route-blockages?status=open&limit=20
```

Filtros disponibles:

- `status`: `open`, `resolved`, `dismissed`
- `route_identifier`: `uuid` opcional
- `route_stop_identifier`: `uuid` opcional
- `bin_identifier`: `uuid` opcional
- `limit`: por defecto `50`, max `200`

### Consulta de rutas para mapa

```text
GET /v1/collection-routes?zone_name=Cochabamba%20Centro&collection_weekday=2&is_active=true
GET /v1/collection-routes/{route_identifier}/stops
```

Filtros disponibles en `collection-routes`:

- `zone_name`: exacto, opcional
- `collection_weekday`: `1..7` (lunes..domingo), opcional
- `is_active`: `true|false`, opcional

### Crear ruta de recoleccion

```json
{
	"route_code": "RUTA-CENTRO-01",
	"route_name": "Ruta Centro Martes/Jueves",
	"zone_name": "Cochabamba Centro",
	"collection_weekday": 2,
	"is_active": true
}
```

### Actualizar metadatos de ruta

```json
{
	"route_name": "Ruta Centro Ajustada",
	"collection_weekday": 4,
	"is_active": true
}
```

### Sincronizar paradas de ruta

```json
{
	"stops": [
		{
			"bin_identifier": "697d3564-8662-4f1e-9dc4-2fa5c6855a55",
			"stop_order": 1,
			"planned_time": "08:00"
		},
		{
			"bin_identifier": "f40e2a53-eb4c-4647-bef6-77533817a4a4",
			"stop_order": 2,
			"planned_time": "08:20"
		}
	]
}
```

### Consultar revisiones de ruta

```text
GET /v1/admin/collection-routes/{route_identifier}/revisions?limit=20
```

### Consultar desvio de camion contra ruta

```text
GET /v1/driver/route-deviation?truck_identifier=TRUCK-001&route_identifier=<uuid>&deviation_threshold_meters=200
```

Reglas:

- `deviation_threshold_meters` es opcional (default `200`)
- `route_identifier` es opcional si existe asignacion activa en `truck_route_assignments`
- si el camion excede el umbral se responde `is_off_route=true`

### Crear alerta de desvio

```json
{
	"truck_identifier": "TRUCK-001",
	"route_identifier": "f33403f5-89c8-4f2f-a26f-c4958ac6e577",
	"deviation_threshold_meters": 200,
	"alert_notes": "desvio detectado en avenida bloqueada"
}
```

Notas:

- si `route_identifier` se omite, se usa la ruta activa asignada al camion
- si el camion esta dentro del umbral se retorna `409 truck_is_within_route_threshold`

### Actualizar estado de bloqueo de ruta

```json
{
	"status": "resolved",
	"resolution_notes": "acceso liberado",
	"resolved_at": "2026-05-17T13:35:00Z"
}
```

Estados permitidos:

- `resolved`
- `dismissed`

### Enviar evidencia

```json
{
	"cycle_identifier": "6e7abf86-b9d5-4fb9-a745-d1f6ca064207",
	"household_identifier": "997aa825-d99d-4d54-b2ba-4f2d138bcf6f",
	"evidence_photo_url": "https://example.com/evidence-001.jpg",
	"evidence_captured_at": "2026-05-16T22:30:00Z",
	"evidence_latitude": -17.3921,
	"evidence_longitude": -66.1561,
	"validation_status": "accepted"
}
```

### Cerrar ciclo

```json
{
	"cycle_identifier": "6e7abf86-b9d5-4fb9-a745-d1f6ca064207",
	"raw_points_total": 100,
	"contamination_level": "medium",
	"contamination_notes": "Se encontro algo de comida",
	"collection_operator_name": "Mamut"
}
```

## Regla actual de descuento por contaminación

- `low`: 0%
- `medium`: 10%
- `high`: 25%

Si envías `contamination_discount_percentage`, ese valor tiene prioridad.

## Eventos NATS emitidos por operaciones

- `operations.bins.sensor_events`: al registrar `POST /v1/bins/sensor-events`
- `operations.driver.collection_events`: al registrar `POST /v1/driver/collection-events`
- `operations.routes.blockage_events`: al registrar `POST /v1/driver/route-blockages`
- `operations.routes.blockage_events`: al registrar `PATCH /v1/driver/route-blockages/{blockage_identifier}` con `event_name=route_blockage_report_status_updated`
- `operations.routes.deviation_alerts`: al registrar `POST /v1/driver/route-deviation-alerts`
