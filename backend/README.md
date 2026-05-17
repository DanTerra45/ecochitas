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

## Endpoints

### Salud

- `GET /healthz`

### GPS

- `GET /v1/trucks/latest-position?truck_identifier=TRUCK-001`
- `GET /v1/trucks/latest-positions`
- `GET /v1/trucks/stream` (SSE)

### Zonas y puntos de reciclaje (MVP v1)

- `GET /v1/recycling/containers?zone_code=ZONA-NORTE`
- `POST /v1/recycling/cycles/start`
- `POST /v1/recycling/evidence-submissions`
- `POST /v1/recycling/cycles/close`
- `GET /v1/recycling/cycles/summary?cycle_identifier=<uuid>`

## Ejemplos payload

### Iniciar ciclo

```json
{
	"container_identifier": "e9960178-c1a8-4de6-ac89-b30f2d0356b9",
	"scheduled_collection_date": "2026-05-16",
	"collection_operator_name": "Mamut"
}
```

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
