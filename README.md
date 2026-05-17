# EcoChitas PoC

Prototipo de plataforma para gestión de recolección y reciclaje en Cochabamba.

## Arquitectura actual

- `backend` en Go (`net/http` + `pgx`).
- NATS para ingestión/eventos GPS de camiones.
- PostgreSQL para persistencia.
- Frontend en SvelteKit + Leaflet (mapa en tiempo real).

## Estado implementado

### 1) Tracking GPS en tiempo real

- Persistencia de posiciones de camiones.
- Snapshot de última posición por camión.
- SSE para stream en vivo de posiciones.
- Simulador de ruta en tiempo real (`gps_route_simulator`).

### 2) Modelo MVP por zonas para reciclaje

Se implementó base backend para el flujo acordado con el PO:

- Zonas de reciclaje.
- Hogares asignados a contenedor de reciclaje.
- Ciclos de recolección de reciclaje por contenedor.
- Evidencia por hogar (foto + ubicación + hora).
- Cierre de ciclo con puntos brutos, contaminación y descuento.
- Reparto de puntos a hogares elegibles (evidencia aceptada).

## Comandos de ejecución

### Backend PoC (PostgreSQL + NATS + API)

```powershell
cd D:\academico_1\ecochitas\backend
.\scripts\start_poc.ps1
```

Notas:

- `start_poc.ps1` ahora aplica **todas** las migraciones de `backend/db/migrations/*.sql`.
- Si ya migraste y quieres arrancar rápido:

```powershell
.\scripts\start_poc.ps1 -skip_migration
```

Para apagar servicios:

```powershell
.\scripts\stop_poc.ps1 -force_ports
```

### Simulación de movimiento GPS continuo

```powershell
cd D:\academico_1\ecochitas\backend
go run ./cmd/gps_route_simulator -truck_identifier TRUCK-001 -publish_interval 1s
```

### Seed demo de zona/contendor/hogares

```powershell
cd D:\academico_1\ecochitas\backend
.\scripts\seed_zone_rewards_demo.ps1
```

### Frontend

```powershell
cd D:\academico_1\ecochitas
npm run dev -- --open
```

## Endpoints principales

### GPS y mapa

- `GET /healthz`
- `GET /v1/trucks/latest-position?truck_identifier=TRUCK-001`
- `GET /v1/trucks/latest-positions`
- `GET /v1/trucks/stream` (SSE)

### Zonas y gamificación MVP

- `GET /v1/recycling/containers?zone_code=ZONA-NORTE`
- `POST /v1/recycling/cycles/start`
- `POST /v1/recycling/evidence-submissions`
- `POST /v1/recycling/cycles/close`
- `GET /v1/recycling/cycles/summary?cycle_identifier=<uuid>`

## Validaciones realizadas

- `go test ./...` en `backend` sin errores.
- Flujo validado en API:

1. Crear ciclo.
2. Registrar evidencias de hogares asignados.
3. Cerrar ciclo con contaminación `medium`.
4. Reparto correcto de puntos entre hogares elegibles.
