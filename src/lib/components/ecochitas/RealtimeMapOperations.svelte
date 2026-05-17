<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { SvelteMap } from 'svelte/reactivity';
	import 'maplibre-gl/dist/maplibre-gl.css';

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

	type Route_path_coordinate = {
		stop_order: number;
		bin_identifier: string;
		bin_code: string;
		latitude: number;
		longitude: number;
	};

	type Collection_route_view = {
		route_identifier: string;
		route_code: string;
		route_name: string;
		zone_name: string;
		collection_weekday: number;
		is_active: boolean;
		stop_total: number;
		path_coordinates: Route_path_coordinate[];
	};

	type Collection_routes_response_payload = {
		items: Collection_route_view[];
		total: number;
	};

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

	type Collection_route_line_feature_properties = {
		route_identifier: string;
		route_code: string;
		route_name: string;
		zone_name: string;
		collection_weekday: number;
		stop_total: number;
		line_color: string;
		popup_html: string;
	};

	type Collection_route_line_feature = {
		type: 'Feature';
		geometry: {
			type: 'LineString';
			coordinates: [number, number][];
		};
		properties: Collection_route_line_feature_properties;
	};

	type Collection_route_line_feature_collection = {
		type: 'FeatureCollection';
		features: Collection_route_line_feature[];
	};

	type Collection_route_stop_feature_properties = {
		route_identifier: string;
		route_code: string;
		route_name: string;
		bin_identifier: string;
		bin_code: string;
		stop_order: number;
		popup_html: string;
	};

	type Collection_route_stop_feature = {
		type: 'Feature';
		geometry: {
			type: 'Point';
			coordinates: [number, number];
		};
		properties: Collection_route_stop_feature_properties;
	};

	type Collection_route_stop_feature_collection = {
		type: 'FeatureCollection';
		features: Collection_route_stop_feature[];
	};

	const cochabamba_center_latitude = -17.3935;
	const cochabamba_center_longitude = -66.157;
	const default_backend_api_base_url = 'http://127.0.0.1:8080';
	const backend_api_url_query_parameter_name = 'backend_api_url';
	const backend_api_url_storage_key = 'ecochitas_backend_api_url';
	const truck_positions_source_identifier = 'truck_positions_source';
	const truck_positions_circle_layer_identifier = 'truck_positions_circle_layer';
	const collection_routes_source_identifier = 'collection_routes_source';
	const collection_routes_line_layer_identifier = 'collection_routes_line_layer';
	const collection_route_stops_source_identifier = 'collection_route_stops_source';
	const collection_route_stops_circle_layer_identifier = 'collection_route_stops_circle_layer';
	const route_line_color_palette = [
		'#16a34a',
		'#0284c7',
		'#dc2626',
		'#7c3aed',
		'#0f766e',
		'#ea580c',
		'#2563eb',
		'#e11d48'
	];

	function get_map_style(is_dark: boolean): import('maplibre-gl').StyleSpecification {
		const tile_url = is_dark
			? 'https://basemaps.cartocdn.com/rastertiles/dark_all/{z}/{x}/{y}.png'
			: 'https://basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}.png';

		return {
			version: 8,
			sources: {
				base_tiles: {
					type: 'raster',
					tiles: [tile_url],
					tileSize: 256,
					attribution: '&copy; OpenStreetMap contributors &copy; CARTO'
				}
			},
			layers: [
				{
					id: 'base_layer',
					type: 'raster',
					source: 'base_tiles',
					minzoom: 0,
					maxzoom: 22
				}
			]
		};
	}

	let latest_positions_list = $state<Truck_latest_position[]>([]);
	let collection_routes_list = $state<Collection_route_view[]>([]);
	let map_container_element = $state<HTMLDivElement | null>(null);
	let backend_api_base_url = $state(default_backend_api_base_url);

	let maplibre_library: typeof import('maplibre-gl') | null = null;
	let map_instance: import('maplibre-gl').Map | null = null;
	let map_resize_observer: ResizeObserver | null = null;
	let truck_stream_connection: EventSource | null = null;
	let map_popup_instance: import('maplibre-gl').Popup | null = null;
	let has_registered_layer_interactions = false;
	let latest_positions_by_truck_identifier = new SvelteMap<string, Truck_latest_position>();
	let theme_observer: MutationObserver | null = null;
	let should_keep_stream_connected = false;
	let stream_reconnect_timeout_identifier: number | null = null;

	onMount(async () => {
		resolve_backend_api_base_url();
		await initialize_map();
		setup_map_resize_observer();
		setup_theme_observer();
		await Promise.all([load_latest_positions_snapshot(), load_collection_routes_snapshot()]);
		should_keep_stream_connected = true;
		connect_truck_stream();
	});

	onDestroy(() => {
		should_keep_stream_connected = false;
		clear_stream_reconnect_timeout();
		disconnect_truck_stream();
		destroy_map();
		theme_observer?.disconnect();
	});

	function resolve_backend_api_base_url() {
		if (typeof window === 'undefined') {
			return;
		}

		const runtime_url = new URL(window.location.href);
		const backend_api_url_query_value = runtime_url.searchParams.get(
			backend_api_url_query_parameter_name
		);
		if (backend_api_url_query_value) {
			backend_api_base_url = normalize_backend_api_base_url(backend_api_url_query_value);
			localStorage.setItem(backend_api_url_storage_key, backend_api_base_url);
			return;
		}

		const saved_backend_api_url = localStorage.getItem(backend_api_url_storage_key);
		if (saved_backend_api_url) {
			backend_api_base_url = normalize_backend_api_base_url(saved_backend_api_url);
			return;
		}

		const runtime_protocol = runtime_url.protocol === 'https:' ? 'https:' : 'http:';
		const runtime_hostname = runtime_url.hostname;
		const is_local_runtime_host =
			runtime_hostname === 'localhost' ||
			runtime_hostname === '127.0.0.1' ||
			runtime_hostname === '0.0.0.0' ||
			runtime_hostname === '[::1]';

		if (runtime_hostname) {
			backend_api_base_url = normalize_backend_api_base_url(
				is_local_runtime_host
					? `${runtime_protocol}//${runtime_hostname}:8080`
					: runtime_url.origin
			);
		}
	}

	function normalize_backend_api_base_url(raw_backend_api_url: string): string {
		return raw_backend_api_url.trim().replace(/\/+$/, '');
	}

	function setup_theme_observer() {
		if (typeof window === 'undefined') return;
		theme_observer = new MutationObserver((mutations) => {
			for (const mutation of mutations) {
				if (mutation.attributeName === 'data-theme') {
					apply_map_theme(document.documentElement.dataset.theme === 'dark');
				}
			}
		});
		theme_observer.observe(document.documentElement, { attributes: true });
	}

	function apply_map_theme(is_dark: boolean) {
		if (!map_instance) return;

		if (!map_instance.isStyleLoaded()) {
			map_instance.once('styledata', () => apply_map_theme(is_dark));
			return;
		}

		map_instance.setStyle(get_map_style(is_dark));

		map_instance.once('style.load', () => {
			has_registered_layer_interactions = false;
			ensure_layers_ready();
			sync_truck_source_data();
			sync_collection_route_source_data();
		});
	}

	async function initialize_map() {
		if (!map_container_element || map_instance) {
			return;
		}

		maplibre_library = await import('maplibre-gl');
		const is_dark_initially =
			typeof document !== 'undefined' && document.documentElement.dataset.theme === 'dark';

		map_instance = new maplibre_library.Map({
			container: map_container_element,
			style: get_map_style(is_dark_initially),
			center: [cochabamba_center_longitude, cochabamba_center_latitude],
			zoom: 13,
			attributionControl: false
		});

		map_instance.on('load', () => {
			ensure_layers_ready();
			sync_truck_source_data();
			sync_collection_route_source_data();
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

		map_popup_instance?.remove();
		map_popup_instance = null;
		has_registered_layer_interactions = false;

		if (map_instance) {
			map_instance.remove();
			map_instance = null;
		}
	}

	function build_snapshot_endpoint_url(): string {
		return `${backend_api_base_url}/v1/trucks/latest-positions`;
	}

	function build_stream_endpoint_url(): string {
		return `${backend_api_base_url}/v1/trucks/stream`;
	}

	function build_collection_routes_endpoint_url(): string {
		return `${backend_api_base_url}/v1/collection-routes?is_active=true`;
	}

	async function load_latest_positions_snapshot() {
		try {
			const snapshot_response = await fetch(build_snapshot_endpoint_url());
			if (!snapshot_response.ok) {
				return;
			}
			const snapshot_payload = (await snapshot_response.json()) as Snapshot_response_payload;
			replace_positions_from_snapshot(snapshot_payload.items);
		} catch {
			// silently ignore — map will populate via stream
		}
	}

	async function load_collection_routes_snapshot() {
		try {
			const collection_routes_response = await fetch(build_collection_routes_endpoint_url());
			if (!collection_routes_response.ok) {
				return;
			}
			const collection_routes_payload =
				(await collection_routes_response.json()) as Collection_routes_response_payload;
			replace_routes_from_snapshot(collection_routes_payload.items);
		} catch {
			// silently ignore route load errors in PoC mode
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

	function replace_routes_from_snapshot(route_items: Collection_route_view[]) {
		collection_routes_list = route_items;
		sync_collection_route_source_data();
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

	function ensure_layers_ready() {
		ensure_truck_layers_ready();
		ensure_collection_route_layers_ready();
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

	function ensure_collection_route_layers_ready() {
		if (!map_instance) {
			return;
		}

		if (!map_instance.getSource(collection_routes_source_identifier)) {
			map_instance.addSource(collection_routes_source_identifier, {
				type: 'geojson',
				data: build_collection_route_line_feature_collection() as unknown as GeoJSON.FeatureCollection
			});
		}

		if (!map_instance.getLayer(collection_routes_line_layer_identifier)) {
			map_instance.addLayer({
				id: collection_routes_line_layer_identifier,
				type: 'line',
				source: collection_routes_source_identifier,
				paint: {
					'line-color': ['coalesce', ['get', 'line_color'], '#0284c7'],
					'line-width': 4,
					'line-opacity': 0.85
				},
				layout: {
					'line-join': 'round',
					'line-cap': 'round'
				}
			});
		}

		if (!map_instance.getSource(collection_route_stops_source_identifier)) {
			map_instance.addSource(collection_route_stops_source_identifier, {
				type: 'geojson',
				data: build_collection_route_stop_feature_collection() as unknown as GeoJSON.FeatureCollection
			});
		}

		if (!map_instance.getLayer(collection_route_stops_circle_layer_identifier)) {
			map_instance.addLayer({
				id: collection_route_stops_circle_layer_identifier,
				type: 'circle',
				source: collection_route_stops_source_identifier,
				paint: {
					'circle-color': '#f8fafc',
					'circle-radius': 4.2,
					'circle-stroke-color': '#0f172a',
					'circle-stroke-width': 1.2,
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

	function sync_collection_route_source_data() {
		if (!map_instance || !map_instance.isStyleLoaded()) {
			return;
		}

		ensure_collection_route_layers_ready();

		const collection_routes_source = map_instance.getSource(collection_routes_source_identifier) as
			| import('maplibre-gl').GeoJSONSource
			| undefined;
		if (collection_routes_source) {
			collection_routes_source.setData(
				build_collection_route_line_feature_collection() as unknown as GeoJSON.FeatureCollection
			);
		}

		const collection_route_stops_source = map_instance.getSource(
			collection_route_stops_source_identifier
		) as import('maplibre-gl').GeoJSONSource | undefined;
		if (collection_route_stops_source) {
			collection_route_stops_source.setData(
				build_collection_route_stop_feature_collection() as unknown as GeoJSON.FeatureCollection
			);
		}
	}

	function register_layer_interactions() {
		if (!map_instance || !maplibre_library) {
			return;
		}

		map_instance.on('mouseenter', truck_positions_circle_layer_identifier, () => {
			if (!map_instance) return;
			map_instance.getCanvas().style.cursor = 'pointer';
		});

		map_instance.on('mouseleave', truck_positions_circle_layer_identifier, () => {
			if (!map_instance) return;
			map_instance.getCanvas().style.cursor = '';
		});

		map_instance.on('click', truck_positions_circle_layer_identifier, (layer_mouse_event) => {
			const clicked_feature = layer_mouse_event.features?.[0];
			if (!clicked_feature || !map_instance || !maplibre_library) return;
			if (clicked_feature.geometry.type !== 'Point') return;

			const clicked_coordinates = clicked_feature.geometry.coordinates as [number, number];
			const popup_html = String(clicked_feature.properties?.['popup_html'] ?? '');
			open_popup(clicked_coordinates, popup_html);
		});

		map_instance.on('mouseenter', collection_routes_line_layer_identifier, () => {
			if (!map_instance) return;
			map_instance.getCanvas().style.cursor = 'pointer';
		});

		map_instance.on('mouseleave', collection_routes_line_layer_identifier, () => {
			if (!map_instance) return;
			map_instance.getCanvas().style.cursor = '';
		});

		map_instance.on('click', collection_routes_line_layer_identifier, (layer_mouse_event) => {
			const clicked_feature = layer_mouse_event.features?.[0];
			if (!clicked_feature || !map_instance || !maplibre_library) return;
			if (clicked_feature.geometry.type !== 'LineString') return;

			const clicked_coordinates = layer_mouse_event.lngLat.toArray() as [number, number];
			const popup_html = String(clicked_feature.properties?.['popup_html'] ?? '');
			open_popup(clicked_coordinates, popup_html);
		});

		map_instance.on('mouseenter', collection_route_stops_circle_layer_identifier, () => {
			if (!map_instance) return;
			map_instance.getCanvas().style.cursor = 'pointer';
		});

		map_instance.on('mouseleave', collection_route_stops_circle_layer_identifier, () => {
			if (!map_instance) return;
			map_instance.getCanvas().style.cursor = '';
		});

		map_instance.on('click', collection_route_stops_circle_layer_identifier, (layer_mouse_event) => {
			const clicked_feature = layer_mouse_event.features?.[0];
			if (!clicked_feature || !map_instance || !maplibre_library) return;
			if (clicked_feature.geometry.type !== 'Point') return;

			const clicked_coordinates = clicked_feature.geometry.coordinates as [number, number];
			const popup_html = String(clicked_feature.properties?.['popup_html'] ?? '');
			open_popup(clicked_coordinates, popup_html);
		});
	}

	function open_popup(clicked_coordinates: [number, number], popup_html: string) {
		if (!map_instance || !maplibre_library) return;

		map_popup_instance?.remove();
		map_popup_instance = new maplibre_library.Popup({
			closeButton: true,
			closeOnClick: true,
			offset: 16
		})
			.setLngLat(clicked_coordinates)
			.setHTML(popup_html)
			.addTo(map_instance);
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
					popup_html: build_truck_marker_popup_html(truck_position_item)
				}
			}))
		};
	}

	function build_collection_route_line_feature_collection(): Collection_route_line_feature_collection {
		const line_feature_list: Collection_route_line_feature[] = [];

		for (const collection_route_item of collection_routes_list) {
			const sorted_coordinates = [...collection_route_item.path_coordinates].sort(
				(left_coordinate, right_coordinate) => left_coordinate.stop_order - right_coordinate.stop_order
			);
			if (sorted_coordinates.length < 2) {
				continue;
			}

			line_feature_list.push({
				type: 'Feature',
				geometry: {
					type: 'LineString',
					coordinates: sorted_coordinates.map((route_coordinate) => [
						route_coordinate.longitude,
						route_coordinate.latitude
					])
				},
				properties: {
					route_identifier: collection_route_item.route_identifier,
					route_code: collection_route_item.route_code,
					route_name: collection_route_item.route_name,
					zone_name: collection_route_item.zone_name,
					collection_weekday: collection_route_item.collection_weekday,
					stop_total: collection_route_item.stop_total,
					line_color: resolve_route_line_color(collection_route_item.route_identifier),
					popup_html: build_collection_route_popup_html(collection_route_item)
				}
			});
		}

		return {
			type: 'FeatureCollection',
			features: line_feature_list
		};
	}

	function build_collection_route_stop_feature_collection(): Collection_route_stop_feature_collection {
		const stop_feature_list: Collection_route_stop_feature[] = [];

		for (const collection_route_item of collection_routes_list) {
			for (const route_stop_coordinate of collection_route_item.path_coordinates) {
				stop_feature_list.push({
					type: 'Feature',
					geometry: {
						type: 'Point',
						coordinates: [route_stop_coordinate.longitude, route_stop_coordinate.latitude]
					},
					properties: {
						route_identifier: collection_route_item.route_identifier,
						route_code: collection_route_item.route_code,
						route_name: collection_route_item.route_name,
						bin_identifier: route_stop_coordinate.bin_identifier,
						bin_code: route_stop_coordinate.bin_code,
						stop_order: route_stop_coordinate.stop_order,
						popup_html: build_collection_route_stop_popup_html(
							collection_route_item,
							route_stop_coordinate
						)
					}
				});
			}
		}

		return {
			type: 'FeatureCollection',
			features: stop_feature_list
		};
	}

	function resolve_route_line_color(route_identifier: string): string {
		const route_identifier_hash = hash_string(route_identifier);
		const color_palette_index = route_identifier_hash % route_line_color_palette.length;
		return route_line_color_palette[color_palette_index] ?? '#0284c7';
	}

	function hash_string(raw_text: string): number {
		let hash_accumulator = 0;
		for (let character_index = 0; character_index < raw_text.length; character_index += 1) {
			hash_accumulator =
				(hash_accumulator * 31 + raw_text.charCodeAt(character_index)) >>> 0;
		}
		return hash_accumulator;
	}

	function build_truck_marker_popup_html(truck_position_item: Truck_latest_position): string {
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

	function build_collection_route_popup_html(collection_route_item: Collection_route_view): string {
		return `
			<strong>${collection_route_item.route_code}</strong><br/>
			${collection_route_item.route_name}<br/>
			Zona: ${collection_route_item.zone_name}<br/>
			Dia: ${collection_route_item.collection_weekday}<br/>
			Paradas: ${collection_route_item.stop_total}
		`;
	}

	function build_collection_route_stop_popup_html(
		collection_route_item: Collection_route_view,
		route_stop_coordinate: Route_path_coordinate
	): string {
		return `
			<strong>${collection_route_item.route_code}</strong><br/>
			Parada #${route_stop_coordinate.stop_order}<br/>
			Contenedor: ${route_stop_coordinate.bin_code}<br/>
			Lat: ${route_stop_coordinate.latitude}<br/>
			Lng: ${route_stop_coordinate.longitude}
		`;
	}

	function clear_stream_reconnect_timeout() {
		if (stream_reconnect_timeout_identifier == null) {
			return;
		}
		window.clearTimeout(stream_reconnect_timeout_identifier);
		stream_reconnect_timeout_identifier = null;
	}

	function schedule_stream_reconnect() {
		if (!should_keep_stream_connected || typeof window === 'undefined') {
			return;
		}
		if (stream_reconnect_timeout_identifier != null) {
			return;
		}

		stream_reconnect_timeout_identifier = window.setTimeout(() => {
			stream_reconnect_timeout_identifier = null;
			if (!should_keep_stream_connected) {
				return;
			}
			connect_truck_stream();
		}, 2000);
	}

	function disconnect_truck_stream() {
		if (truck_stream_connection) {
			truck_stream_connection.close();
			truck_stream_connection = null;
		}
	}

	function connect_truck_stream() {
		disconnect_truck_stream();

		const stream_endpoint_url = build_stream_endpoint_url();
		const stream_connection = new EventSource(stream_endpoint_url);
		truck_stream_connection = stream_connection;

		stream_connection.addEventListener('open', () => {
			clear_stream_reconnect_timeout();
		});

		stream_connection.addEventListener('ready', (ready_event) => {
			if (truck_stream_connection !== stream_connection) {
				return;
			}
			try {
				JSON.parse((ready_event as MessageEvent).data) as Stream_ready_payload;
			} catch {
				// ignore malformed ready payload
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
			} catch {
				// ignore malformed position payload
			}
		});

		stream_connection.addEventListener('error', () => {
			if (truck_stream_connection !== stream_connection) {
				return;
			}
			stream_connection.close();
			truck_stream_connection = null;
			schedule_stream_reconnect();
		});
	}
</script>

<div bind:this={map_container_element} class="map_fullscreen"></div>

<div class="map_status_panel">
	<div>API: {backend_api_base_url}</div>
	<div>Rutas activas: {collection_routes_list.length}</div>
	<div>Camiones en vivo: {latest_positions_list.length}</div>
</div>

<style>
	.map_fullscreen {
		position: fixed;
		top: 0;
		left: 0;
		width: 100vw;
		height: 100vh;
		z-index: 0;
	}

	.map_status_panel {
		position: fixed;
		right: 1rem;
		top: 5.2rem;
		z-index: 40;
		display: grid;
		gap: 0.25rem;
		padding: 0.65rem 0.8rem;
		border-radius: 0.75rem;
		background: oklch(0.2 0 0 / 0.75);
		border: 1px solid oklch(1 0 0 / 0.2);
		color: #f8fafc;
		font-size: 0.75rem;
		line-height: 1.25;
		backdrop-filter: blur(8px);
		-webkit-backdrop-filter: blur(8px);
		max-width: min(90vw, 22rem);
		word-break: break-all;
	}

	@media (max-width: 1023px) {
		.map_status_panel {
			top: 5.4rem;
			right: 0.7rem;
			left: 0.7rem;
			max-width: none;
		}
	}

	:global(.maplibregl-ctrl-attrib) {
		font-size: 10px;
		background: oklch(1 0 0 / 0.84);
	}
</style>
