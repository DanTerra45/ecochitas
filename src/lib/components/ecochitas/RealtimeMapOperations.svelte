<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { SvelteMap } from 'svelte/reactivity';
	import 'maplibre-gl/dist/maplibre-gl.css';
	import EcochitasIcon from '$lib/components/ecochitas/EcochitasIcon.svelte';

	type Truck_latest_position = {
		truck_identifier: string;
		latitude: number;
		longitude: number;
		speed_kmh?: number | null;
		heading_degrees?: number | null;
		captured_at: string;
		received_at: string;
	};

	type Snapshot_response_payload = {
		items: Truck_latest_position[];
		total: number;
	};

	type Stream_ready_payload = {
		subscriber_identifier: string;
		truck_identifier_filter: string;
		connected_at: string;
	};

	type Stream_status_label = 'disconnected' | 'connecting' | 'connected' | 'error';

	type Truck_feature_properties = {
		truck_identifier: string;
		latitude: number;
		longitude: number;
		speed_kmh: number | null;
		captured_at: string;
		popup_html: string;
	};

	type Truck_feature = {
		type: 'Feature';
		geometry: {
			type: 'Point';
			coordinates: [number, number];
		};
		properties: Truck_feature_properties;
	};

	type Truck_feature_collection = {
		type: 'FeatureCollection';
		features: Truck_feature[];
	};

	const cochabamba_center_latitude = -17.3935;
	const cochabamba_center_longitude = -66.157;
	const default_backend_api_url = 'http://127.0.0.1:8080';
	const truck_positions_source_identifier = 'truck_positions_source';
	const truck_positions_circle_layer_identifier = 'truck_positions_circle_layer';

	const map_style_specification: import('maplibre-gl').StyleSpecification = {
		version: 8,
		sources: {
			openstreetmap_tiles: {
				type: 'raster',
				tiles: ['https://tile.openstreetmap.org/{z}/{x}/{y}.png'],
				tileSize: 256,
				attribution: '&copy; OpenStreetMap contributors'
			}
		},
		layers: [
			{
				id: 'openstreetmap_base_layer',
				type: 'raster',
				source: 'openstreetmap_tiles',
				minzoom: 0,
				maxzoom: 22
			}
		]
	};

	let backend_api_url = $state(default_backend_api_url);
	let truck_identifier_filter = $state('');
	let stream_status = $state<Stream_status_label>('disconnected');
	let stream_ready_message = $state('');
	let latest_error_message = $state('');
	let snapshot_total = $state(0);
	let stream_event_count = $state(0);
	let is_loading_snapshot = $state(false);
	let latest_positions_list = $state<Truck_latest_position[]>([]);

	let map_container_element = $state<HTMLDivElement | null>(null);

	let maplibre_library: typeof import('maplibre-gl') | null = null;
	let map_instance: import('maplibre-gl').Map | null = null;
	let map_resize_observer: ResizeObserver | null = null;
	let truck_stream_connection: EventSource | null = null;
	let truck_popup_instance: import('maplibre-gl').Popup | null = null;
	let has_registered_layer_interactions = false;
	let latest_positions_by_truck_identifier = new SvelteMap<string, Truck_latest_position>();

	onMount(async () => {
		await initialize_map();
		setup_map_resize_observer();
		await load_latest_positions_snapshot();
		connect_truck_stream();
	});

	onDestroy(() => {
		disconnect_truck_stream();
		destroy_map();
	});

	async function initialize_map() {
		if (!map_container_element || map_instance) {
			return;
		}

		maplibre_library = await import('maplibre-gl');

		map_instance = new maplibre_library.Map({
			container: map_container_element,
			style: map_style_specification,
			center: [cochabamba_center_longitude, cochabamba_center_latitude],
			zoom: 13
		});

		map_instance.addControl(new maplibre_library.NavigationControl(), 'top-left');

		map_instance.on('load', () => {
			ensure_truck_layers_ready();
			sync_truck_source_data();
		});
	}

	function setup_map_resize_observer() {
		if (!map_container_element || !map_instance) {
			return;
		}

		map_resize_observer = new ResizeObserver(() => {
			if (!map_instance) {
				return;
			}
			map_instance.resize();
		});
		map_resize_observer.observe(map_container_element);
	}

	function destroy_map() {
		map_resize_observer?.disconnect();
		map_resize_observer = null;

		truck_popup_instance?.remove();
		truck_popup_instance = null;
		has_registered_layer_interactions = false;

		if (map_instance) {
			map_instance.remove();
			map_instance = null;
		}
	}

	function get_normalized_backend_api_url(): string {
		return backend_api_url.trim().replace(/\/+$/, '');
	}

	function build_snapshot_endpoint_url(): string {
		return `${get_normalized_backend_api_url()}/v1/trucks/latest-positions`;
	}

	function build_stream_endpoint_url(): string {
		const base_stream_endpoint_url = `${get_normalized_backend_api_url()}/v1/trucks/stream`;
		const normalized_truck_identifier_filter = truck_identifier_filter.trim();
		if (!normalized_truck_identifier_filter) {
			return base_stream_endpoint_url;
		}

		const encoded_truck_identifier = encodeURIComponent(normalized_truck_identifier_filter);
		return `${base_stream_endpoint_url}?truck_identifier=${encoded_truck_identifier}`;
	}

	async function load_latest_positions_snapshot() {
		is_loading_snapshot = true;
		latest_error_message = '';

		try {
			const snapshot_response = await fetch(build_snapshot_endpoint_url());
			if (!snapshot_response.ok) {
				throw new Error(`snapshot_request_failed_status_${snapshot_response.status}`);
			}

			const snapshot_payload = (await snapshot_response.json()) as Snapshot_response_payload;
			snapshot_total = snapshot_payload.total;
			replace_positions_from_snapshot(snapshot_payload.items);
		} catch (snapshot_error) {
			latest_error_message =
				snapshot_error instanceof Error ? snapshot_error.message : 'snapshot_request_failed';
		} finally {
			is_loading_snapshot = false;
		}
	}

	function replace_positions_from_snapshot(snapshot_items: Truck_latest_position[]) {
		latest_positions_by_truck_identifier.clear();

		for (const truck_position_item of snapshot_items) {
			latest_positions_by_truck_identifier.set(
				truck_position_item.truck_identifier,
				truck_position_item
			);
		}

		sync_positions_list_from_map();
		sync_truck_source_data();
	}

	function upsert_truck_position(truck_position_item: Truck_latest_position) {
		latest_positions_by_truck_identifier.set(
			truck_position_item.truck_identifier,
			truck_position_item
		);
		sync_positions_list_from_map();
		sync_truck_source_data();
	}

	function sync_positions_list_from_map() {
		latest_positions_list = Array.from(latest_positions_by_truck_identifier.values()).sort(
			(left_item, right_item) =>
				left_item.truck_identifier.localeCompare(right_item.truck_identifier)
		);
	}

	function ensure_truck_layers_ready() {
		if (!map_instance) {
			return;
		}

		if (!map_instance.getSource(truck_positions_source_identifier)) {
			map_instance.addSource(truck_positions_source_identifier, {
				type: 'geojson',
				data: build_truck_feature_collection() as unknown as GeoJSON.FeatureCollection
			});
		}

		if (!map_instance.getLayer(truck_positions_circle_layer_identifier)) {
			map_instance.addLayer({
				id: truck_positions_circle_layer_identifier,
				type: 'circle',
				source: truck_positions_source_identifier,
				paint: {
					'circle-color': '#18a26e',
					'circle-radius': 7,
					'circle-stroke-color': '#0a5e4a',
					'circle-stroke-width': 2.2,
					'circle-opacity': 0.95
				}
			});
		}

		if (!has_registered_layer_interactions) {
			register_layer_interactions();
			has_registered_layer_interactions = true;
		}
	}

	function sync_truck_source_data() {
		if (!map_instance || !map_instance.isStyleLoaded()) {
			return;
		}

		ensure_truck_layers_ready();

		const truck_positions_source = map_instance.getSource(truck_positions_source_identifier) as
			| import('maplibre-gl').GeoJSONSource
			| undefined;
		if (!truck_positions_source) {
			return;
		}

		truck_positions_source.setData(
			build_truck_feature_collection() as unknown as GeoJSON.FeatureCollection
		);
	}

	function register_layer_interactions() {
		if (!map_instance || !maplibre_library) {
			return;
		}

		map_instance.on('mouseenter', truck_positions_circle_layer_identifier, () => {
			if (!map_instance) {
				return;
			}
			map_instance.getCanvas().style.cursor = 'pointer';
		});

		map_instance.on('mouseleave', truck_positions_circle_layer_identifier, () => {
			if (!map_instance) {
				return;
			}
			map_instance.getCanvas().style.cursor = '';
		});

		map_instance.on('click', truck_positions_circle_layer_identifier, (layer_mouse_event) => {
			const clicked_feature = layer_mouse_event.features?.[0];
			if (!clicked_feature || !map_instance || !maplibre_library) {
				return;
			}

			const clicked_feature_geometry = clicked_feature.geometry;
			if (clicked_feature_geometry.type !== 'Point') {
				return;
			}

			const clicked_coordinates = clicked_feature_geometry.coordinates as [number, number];
			const popup_html = String(clicked_feature.properties?.['popup_html'] ?? '');

			truck_popup_instance?.remove();
			truck_popup_instance = new maplibre_library.Popup({
				closeButton: true,
				closeOnClick: true,
				offset: 16
			})
				.setLngLat(clicked_coordinates)
				.setHTML(popup_html)
				.addTo(map_instance);
		});
	}

	function build_truck_feature_collection(): Truck_feature_collection {
		return {
			type: 'FeatureCollection',
			features: latest_positions_list.map((truck_position_item) => ({
				type: 'Feature',
				geometry: {
					type: 'Point',
					coordinates: [truck_position_item.longitude, truck_position_item.latitude]
				},
				properties: {
					truck_identifier: truck_position_item.truck_identifier,
					latitude: truck_position_item.latitude,
					longitude: truck_position_item.longitude,
					speed_kmh: truck_position_item.speed_kmh ?? null,
					captured_at: truck_position_item.captured_at,
					popup_html: build_marker_popup_html(truck_position_item)
				}
			}))
		};
	}

	function build_marker_popup_html(truck_position_item: Truck_latest_position): string {
		const speed_label =
			truck_position_item.speed_kmh == null ? 'N/A' : `${truck_position_item.speed_kmh} km/h`;
		return `
			<strong>${truck_position_item.truck_identifier}</strong><br/>
			Lat: ${truck_position_item.latitude}<br/>
			Lng: ${truck_position_item.longitude}<br/>
			Speed: ${speed_label}<br/>
			Captured: ${truck_position_item.captured_at}
		`;
	}

	function disconnect_truck_stream() {
		if (truck_stream_connection) {
			truck_stream_connection.close();
			truck_stream_connection = null;
		}
		stream_status = 'disconnected';
		stream_ready_message = '';
		latest_error_message = '';
	}

	function connect_truck_stream() {
		disconnect_truck_stream();
		latest_error_message = '';
		stream_status = 'connecting';

		const stream_endpoint_url = build_stream_endpoint_url();
		const stream_connection = new EventSource(stream_endpoint_url);
		truck_stream_connection = stream_connection;

		stream_connection.onopen = () => {
			if (truck_stream_connection !== stream_connection) {
				return;
			}

			stream_status = 'connected';
			latest_error_message = '';
		};

		stream_connection.addEventListener('ready', (ready_event) => {
			if (truck_stream_connection !== stream_connection) {
				return;
			}

			try {
				const ready_payload = JSON.parse(
					(ready_event as MessageEvent).data
				) as Stream_ready_payload;
				stream_status = 'connected';
				stream_ready_message = `subscriber=${ready_payload.subscriber_identifier}`;
				latest_error_message = '';
			} catch {
				stream_status = 'connected';
				stream_ready_message = 'stream_ready';
				latest_error_message = '';
			}
		});

		stream_connection.addEventListener('truck_position', (position_event) => {
			if (truck_stream_connection !== stream_connection) {
				return;
			}

			try {
				const truck_position_payload = JSON.parse(
					(position_event as MessageEvent).data
				) as Truck_latest_position;
				upsert_truck_position(truck_position_payload);
				stream_event_count += 1;
				stream_status = 'connected';
				latest_error_message = '';
			} catch (parse_error) {
				latest_error_message =
					parse_error instanceof Error ? parse_error.message : 'stream_payload_parse_failed';
			}
		});

		stream_connection.onerror = () => {
			if (truck_stream_connection !== stream_connection) {
				return;
			}

			if (stream_connection.readyState === EventSource.CONNECTING) {
				stream_status = 'connecting';
				return;
			}

			stream_status = 'error';
			if (stream_connection.readyState === EventSource.CLOSED) {
				latest_error_message = 'stream_connection_closed';
				return;
			}

			latest_error_message = 'stream_connection_error';
		};
	}

	async function reload_snapshot_and_reconnect_stream() {
		await load_latest_positions_snapshot();
		connect_truck_stream();
	}

	function format_capture_time(captured_at_value: string): string {
		const captured_at_date = new Date(captured_at_value);
		if (Number.isNaN(captured_at_date.getTime())) {
			return captured_at_value;
		}
		return captured_at_date.toLocaleString();
	}
</script>

<div class="map_page_grid">
	<aside class="panel controls_panel">
		<div class="panel_title_row">
			<EcochitasIcon name="schedule" size={18} />
			<h2>Conexion y control</h2>
		</div>
		<p class="muted_text">Configura backend y stream para pruebas de campo.</p>

		<label>
			Backend API URL
			<input bind:value={backend_api_url} placeholder="http://127.0.0.1:8080" />
		</label>

		<label>
			Filtro de camion (opcional)
			<input bind:value={truck_identifier_filter} placeholder="TRUCK-001" />
		</label>

		<div class="button_grid">
			<button onclick={reload_snapshot_and_reconnect_stream} disabled={is_loading_snapshot}>
				{is_loading_snapshot ? 'Cargando...' : 'Recargar + reconectar'}
			</button>
			<button onclick={load_latest_positions_snapshot} disabled={is_loading_snapshot}>
				Recargar snapshot
			</button>
			<button onclick={connect_truck_stream}>Reconectar stream</button>
			<button onclick={disconnect_truck_stream}>Desconectar stream</button>
		</div>

		<div class="status_grid">
			<div>
				<span>Estado</span>
				<strong>{stream_status}</strong>
			</div>
			<div>
				<span>Snapshot total</span>
				<strong>{snapshot_total}</strong>
			</div>
			<div>
				<span>Eventos stream</span>
				<strong>{stream_event_count}</strong>
			</div>
			<div>
				<span>Subscriber</span>
				<strong>{stream_ready_message || 'N/A'}</strong>
			</div>
		</div>

		{#if latest_error_message}
			<p class="error_message">Error: {latest_error_message}</p>
		{/if}
	</aside>

	<section class="panel map_panel">
		<div class="map_panel_header">
			<div class="panel_title_row">
				<EcochitasIcon name="map" size={18} />
				<h2>Mapa en vivo</h2>
			</div>
			<p>{latest_positions_list.length} camion(es) activos</p>
		</div>
		<div bind:this={map_container_element} class="maplibre_map"></div>
	</section>

	<section class="panel list_panel">
		<div class="map_panel_header">
			<div class="panel_title_row">
				<EcochitasIcon name="truck" size={18} />
				<h2>Ultimas posiciones</h2>
			</div>
			<p>Stream SSE + snapshot</p>
		</div>

		{#if latest_positions_list.length === 0}
			<p class="empty_state_text">
				No hay camiones aun. Publica eventos con
				<code>go run ./cmd/gps_route_simulator -publish_interval 1s</code>
			</p>
		{:else}
			<ul class="positions_list">
				{#each latest_positions_list as truck_position_item (truck_position_item.truck_identifier)}
					<li>
						<div class="position_row_top">
							<strong>{truck_position_item.truck_identifier}</strong>
							<span>{format_capture_time(truck_position_item.captured_at)}</span>
						</div>
						<div class="position_row_bottom">
							<span>{truck_position_item.latitude}, {truck_position_item.longitude}</span>
							<span>
								{truck_position_item.speed_kmh == null
									? 'N/A'
									: `${truck_position_item.speed_kmh} km/h`}
							</span>
						</div>
					</li>
				{/each}
			</ul>
		{/if}
	</section>
</div>

<style>
	h2 {
		margin: 0;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		font-size: 1rem;
		line-height: 1.2;
		letter-spacing: -0.015em;
	}

	label {
		display: block;
		margin-top: 0.65rem;
		font-size: 0.82rem;
		font-weight: 600;
	}

	input {
		width: 100%;
		margin-top: 0.3rem;
		padding: 0.63rem 0.72rem;
		font-size: 0.9rem;
		border-radius: 0.8rem;
		border: 1px solid oklch(0.84 0.02 140);
		background: oklch(0.99 0.005 120);
		color: var(--ecochitas-ink);
		box-sizing: border-box;
	}

	input:focus-visible {
		outline: none;
		border-color: var(--ecochitas-leaf);
		box-shadow: 0 0 0 3px oklch(0.88 0.07 155 / 0.42);
	}

	button {
		padding: 0.6rem 0.68rem;
		border-radius: 0.9rem;
		border: 1px solid oklch(0.86 0.03 145);
		background: linear-gradient(135deg, oklch(0.98 0.01 150), oklch(0.95 0.03 220));
		color: var(--ecochitas-ink);
		font-size: 0.8rem;
		font-weight: 700;
		cursor: pointer;
	}

	button:disabled {
		opacity: 0.55;
		cursor: not-allowed;
	}

	.map_page_grid {
		display: grid;
		grid-template-columns: 1fr;
		grid-template-areas:
			'controls'
			'map'
			'list';
		gap: 0.85rem;
	}

	.panel {
		background: var(--ecochitas-surface);
		backdrop-filter: blur(14px);
		-webkit-backdrop-filter: blur(14px);
		border: 1px solid var(--ecochitas-border);
		border-radius: 1.5rem;
		padding: 1rem;
		box-shadow: var(--ecochitas-shadow-soft);
	}

	.controls_panel {
		grid-area: controls;
	}

	.map_panel {
		grid-area: map;
		padding: 0.85rem;
	}

	.list_panel {
		grid-area: list;
	}

	.panel_title_row {
		display: flex;
		align-items: center;
		gap: 0.45rem;
		color: var(--ecochitas-ink);
	}

	.muted_text {
		margin: 0.3rem 0 0.9rem;
		font-size: 0.82rem;
		color: var(--ecochitas-muted);
	}

	.button_grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0.5rem;
		margin-top: 0.95rem;
	}

	.status_grid {
		margin-top: 0.9rem;
		padding: 0.7rem 0.8rem;
		border-radius: 1rem;
		background: oklch(1 0 0 / 0.72);
		border: 1px solid var(--ecochitas-border);
		display: grid;
		gap: 0.5rem;
	}

	.status_grid div {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
		gap: 0.5rem;
	}

	.status_grid span {
		font-size: 0.78rem;
		color: var(--ecochitas-muted);
	}

	.status_grid strong {
		font-size: 0.78rem;
		font-weight: 700;
		word-break: break-word;
		text-align: right;
	}

	.error_message {
		margin: 0.7rem 0 0;
		font-size: 0.8rem;
		color: var(--ecochitas-alert);
		font-weight: 700;
	}

	.map_panel_header {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
		padding: 0.2rem 0.2rem 0.6rem;
		gap: 0.6rem;
	}

	.map_panel_header p {
		margin: 0;
		font-size: 0.76rem;
		color: var(--ecochitas-muted);
	}

	.maplibre_map {
		width: 100%;
		height: 52vh;
		min-height: 300px;
		border-radius: 1.1rem;
		overflow: hidden;
		border: 1px solid oklch(0.88 0.03 165);
	}

	.empty_state_text {
		margin: 0;
		font-size: 0.84rem;
		color: var(--ecochitas-muted);
		line-height: 1.4;
	}

	.empty_state_text code {
		font-size: 0.74rem;
		background: oklch(0.96 0.02 150);
		padding: 0.16rem 0.3rem;
		border-radius: 0.4rem;
	}

	.positions_list {
		margin: 0;
		padding: 0;
		list-style: none;
		display: grid;
		gap: 0.5rem;
	}

	.positions_list li {
		padding: 0.66rem 0.75rem;
		border-radius: 0.9rem;
		background: oklch(1 0 0 / 0.72);
		border: 1px solid oklch(0.88 0.02 145);
	}

	.position_row_top,
	.position_row_bottom {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
		gap: 0.6rem;
	}

	.position_row_top strong {
		font-size: 0.82rem;
	}

	.position_row_top span,
	.position_row_bottom span {
		font-size: 0.74rem;
		color: var(--ecochitas-muted);
	}

	.position_row_bottom {
		margin-top: 0.22rem;
	}

	:global(.maplibregl-ctrl-attrib) {
		font-size: 10px;
		background: oklch(1 0 0 / 0.84);
	}

	@media (min-width: 1024px) {
		.map_page_grid {
			grid-template-columns: 360px 1fr;
			grid-template-areas:
				'controls map'
				'controls list';
			gap: 0.95rem;
		}

		.controls_panel {
			max-height: 100%;
		}

		.map_panel {
			min-height: 520px;
		}

		.maplibre_map {
			height: 57vh;
			min-height: 500px;
		}

		.list_panel {
			max-height: 280px;
			overflow: auto;
		}
	}

	@media (min-width: 1300px) {
		.maplibre_map {
			height: 60vh;
		}
	}
</style>
