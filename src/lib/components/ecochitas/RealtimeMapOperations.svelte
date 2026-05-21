<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import {
		MapLibre,
		RasterTileSource,
		RasterLayer,
		GeoJSONSource,
		LineLayer,
		CircleLayer,
		Marker,
		Popup
	} from 'svelte-maplibre-gl';
	import { SvelteMap } from 'svelte/reactivity';
	import type * as maplibregl from 'maplibre-gl';
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

	type Route_road_path_coordinate = [number, number];

	type Collection_route_view = {
		route_identifier: string;
		route_code: string;
		route_name: string;
		zone_name: string;
		collection_weekday: number;
		is_active: boolean;
		stop_total: number;
		path_coordinates: Route_path_coordinate[];
		road_path_coordinates?: Route_road_path_coordinate[];
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

	type Route_line_candidate = {
		route_identifier: string;
		line_coordinates: [number, number][];
	};

	type Route_projection_match = {
		route_identifier: string;
		segment_index: number;
		snapped_coordinates: [number, number];
		distance_meters: number;
	};

	type Truck_route_snap_state = {
		route_identifier: string;
		segment_index: number;
		snapped_coordinates: [number, number];
		raw_coordinates: [number, number];
		captured_at_unix_milliseconds: number;
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

	type SmartBin = {
		id: string;
		label: string;
		lat: number;
		lng: number;
		capacity_pct: number;
		zone: string;
		last_updated: string;
	};

	type AcopioMarker = {
		id: string;
		name: string;
		address: string;
		zone: string;
		schedule: string;
		lat: number;
		lng: number;
		capacity_pct: number;
		materials: string;
	};

	type Active_map_popup_state = {
		lnglat: [number, number];
		popup_html: string;
		associated_truck_identifier: string | null;
		open: boolean;
	};

	type Truck_motion_animation_state = {
		start_coordinates: [number, number];
		end_coordinates: [number, number];
		started_at_unix_milliseconds: number;
		duration_milliseconds: number;
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
	const smart_bins_source_identifier = 'smart_bins_source';
	const smart_bins_circle_layer_identifier = 'smart_bins_circle_layer';
	const acopio_points_source_identifier = 'acopio_points_source';
	const acopio_points_circle_layer_identifier = 'acopio_points_circle_layer';
	const truck_route_snap_distance_threshold_meters = 250;
	const truck_route_snap_segment_window_radius = 6;
	const truck_route_snap_max_jump_meters = 90;
	const truck_motion_instant_jump_threshold_meters = 400;
	const truck_motion_min_duration_milliseconds = 240;
	const truck_motion_max_duration_milliseconds = 1400;
	const truck_motion_milliseconds_per_meter = 12;
	const meters_per_degree_latitude = 110540;
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

	const mock_smart_bins: SmartBin[] = [
		{
			id: 'b01',
			label: 'Contenedor A-01',
			lat: -17.3835,
			lng: -66.157,
			capacity_pct: 45,
			zone: 'Zona Norte',
			last_updated: 'hace 12 min'
		},
		{
			id: 'b02',
			label: 'Contenedor A-02',
			lat: -17.3855,
			lng: -66.161,
			capacity_pct: 88,
			zone: 'Zona Norte',
			last_updated: 'hace 8 min'
		},
		{
			id: 'b03',
			label: 'Contenedor B-01',
			lat: -17.391,
			lng: -66.164,
			capacity_pct: 32,
			zone: 'Zona Central',
			last_updated: 'hace 5 min'
		},
		{
			id: 'b04',
			label: 'Contenedor B-02',
			lat: -17.393,
			lng: -66.151,
			capacity_pct: 91,
			zone: 'Zona Central',
			last_updated: 'hace 3 min'
		},
		{
			id: 'b05',
			label: 'Contenedor C-01',
			lat: -17.3985,
			lng: -66.158,
			capacity_pct: 60,
			zone: 'Zona Sur',
			last_updated: 'hace 20 min'
		},
		{
			id: 'b06',
			label: 'Contenedor C-02',
			lat: -17.401,
			lng: -66.1555,
			capacity_pct: 78,
			zone: 'Zona Sur',
			last_updated: 'hace 15 min'
		},
		{
			id: 'b07',
			label: 'Contenedor D-01',
			lat: -17.3872,
			lng: -66.149,
			capacity_pct: 95,
			zone: 'Zona Este',
			last_updated: 'hace 2 min'
		},
		{
			id: 'b08',
			label: 'Contenedor D-02',
			lat: -17.39,
			lng: -66.146,
			capacity_pct: 22,
			zone: 'Zona Este',
			last_updated: 'hace 30 min'
		},
		{
			id: 'b09',
			label: 'Contenedor E-01',
			lat: -17.396,
			lng: -66.169,
			capacity_pct: 55,
			zone: 'Zona Oeste',
			last_updated: 'hace 18 min'
		},
		{
			id: 'b10',
			label: 'Contenedor E-02',
			lat: -17.4005,
			lng: -66.172,
			capacity_pct: 83,
			zone: 'Zona Oeste',
			last_updated: 'hace 10 min'
		}
	];

	const mock_acopio_points: AcopioMarker[] = [
		{
			id: 'a01',
			name: 'Centro de Acopio Norte',
			address: 'Av. América 2345',
			zone: 'Zona Norte',
			schedule: 'Lun–Sáb 8:00–18:00',
			lat: -17.384,
			lng: -66.16,
			capacity_pct: 45,
			materials: 'Plástico · Papel · Cartón · Vidrio'
		},
		{
			id: 'a02',
			name: 'Acopio Villa Pagador',
			address: 'Calle Sucre 1100',
			zone: 'Zona Sur',
			schedule: 'Mar–Dom 9:00–17:00',
			lat: -17.402,
			lng: -66.156,
			capacity_pct: 70,
			materials: 'Metal · Electrónico · Plástico'
		},
		{
			id: 'a03',
			name: 'Acopio Queru Queru',
			address: 'Av. Blanco Galindo km 3',
			zone: 'Zona Oeste',
			schedule: 'Lun–Vie 7:00–16:00',
			lat: -17.395,
			lng: -66.17,
			capacity_pct: 30,
			materials: 'Orgánico · Papel · Cartón'
		},
		{
			id: 'a04',
			name: 'Acopio Colón',
			address: 'Plaza Colón 500',
			zone: 'Zona Central',
			schedule: 'Lun–Dom 8:00–20:00',
			lat: -17.3925,
			lng: -66.1565,
			capacity_pct: 85,
			materials: 'Plástico · Vidrio · Metal · Aceite'
		},
		{
			id: 'a05',
			name: 'Acopio Muyurina',
			address: 'Av. Aniceto Arce 800',
			zone: 'Zona Este',
			schedule: 'Mié–Dom 10:00–18:00',
			lat: -17.3875,
			lng: -66.148,
			capacity_pct: 55,
			materials: 'Pilas · Electrónico · Plástico'
		}
	];

	function get_base_tile_url(is_dark: boolean): string {
		return is_dark
			? 'https://basemaps.cartocdn.com/rastertiles/dark_all/{z}/{x}/{y}.png'
			: 'https://basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}.png';
	}

	const empty_map_style: maplibregl.StyleSpecification = {
		version: 8,
		sources: {},
		layers: []
	};

	const user_marker_coordinates: [number, number] = [
		cochabamba_center_longitude - 0.0015,
		cochabamba_center_latitude + 0.005
	];

	let latest_positions_list = $state<Truck_latest_position[]>([]);
	let collection_routes_list = $state<Collection_route_view[]>([]);
	let backend_api_base_url = $state(default_backend_api_base_url);
	let is_dark_map_theme = $state(false);
	let map_canvas_cursor = $state('');
	let map_instance = $state<maplibregl.Map | undefined>(undefined);

	let truck_stream_connection: EventSource | null = null;
	let should_keep_stream_connected = false;
	let stream_reconnect_timeout_identifier: number | null = null;
	let active_map_popup_state = $state<Active_map_popup_state | null>(null);

	let latest_positions_by_truck_identifier = new SvelteMap<string, Truck_latest_position>();
	let latest_snap_state_by_truck_identifier = new SvelteMap<string, Truck_route_snap_state>();
	let rendered_coordinates_by_truck_identifier = new SvelteMap<string, [number, number]>();
	let active_motion_by_truck_identifier = new Map<string, Truck_motion_animation_state>();
	let truck_motion_frame_request_identifier: number | null = null;

	let truck_feature_collection = $state<Truck_feature_collection>({
		type: 'FeatureCollection',
		features: []
	});
	let collection_route_line_feature_collection = $state<Collection_route_line_feature_collection>({
		type: 'FeatureCollection',
		features: []
	});
	let collection_route_stop_feature_collection = $state<Collection_route_stop_feature_collection>({
		type: 'FeatureCollection',
		features: []
	});
	let smart_bin_feature_collection = $state<GeoJSON.FeatureCollection>({
		type: 'FeatureCollection',
		features: []
	});
	let acopio_feature_collection = $state<GeoJSON.FeatureCollection>({
		type: 'FeatureCollection',
		features: []
	});

	let theme_observer: MutationObserver | null = null;

	const base_tile_url = $derived(get_base_tile_url(is_dark_map_theme));

	onMount(async () => {
		resolve_backend_api_base_url();
		is_dark_map_theme =
			typeof document !== 'undefined' && document.documentElement.dataset.theme === 'dark';
		setup_theme_observer();
		recompute_static_feature_collections();
		await Promise.all([load_latest_positions_snapshot(), load_collection_routes_snapshot()]);
		should_keep_stream_connected = true;
		connect_truck_stream();
	});

	onDestroy(() => {
		should_keep_stream_connected = false;
		stop_truck_motion_loop();
		clear_stream_reconnect_timeout();
		disconnect_truck_stream();
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
				is_local_runtime_host ? `${runtime_protocol}//${runtime_hostname}:8080` : runtime_url.origin
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
					is_dark_map_theme = document.documentElement.dataset.theme === 'dark';
				}
			}
		});
		theme_observer.observe(document.documentElement, { attributes: true });
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
		latest_snap_state_by_truck_identifier.clear();
		rendered_coordinates_by_truck_identifier.clear();
		active_motion_by_truck_identifier.clear();
		stop_truck_motion_loop();

		for (const truck_position_item of snapshot_items) {
			latest_positions_by_truck_identifier.set(
				truck_position_item.truck_identifier,
				truck_position_item
			);
		}

		sync_positions_list_from_map();
		recompute_truck_feature_collection({ animate: false });
	}

	function replace_routes_from_snapshot(route_items: Collection_route_view[]) {
		collection_routes_list = route_items;
		recompute_route_feature_collections();
		recompute_truck_feature_collection({ animate: false });
	}

	function upsert_truck_position(truck_position_item: Truck_latest_position) {
		const existing_truck_position = latest_positions_by_truck_identifier.get(
			truck_position_item.truck_identifier
		);
		if (
			existing_truck_position &&
			!is_truck_position_captured_after(existing_truck_position, truck_position_item)
		) {
			return;
		}

		latest_positions_by_truck_identifier.set(
			truck_position_item.truck_identifier,
			truck_position_item
		);
		sync_positions_list_from_map();
		recompute_truck_feature_collection({ animate: true });
	}

	function is_truck_position_captured_after(
		existing_truck_position: Truck_latest_position,
		incoming_truck_position: Truck_latest_position
	): boolean {
		const existing_captured_at_unix_milliseconds = Date.parse(existing_truck_position.captured_at);
		const incoming_captured_at_unix_milliseconds = Date.parse(incoming_truck_position.captured_at);

		if (
			Number.isNaN(existing_captured_at_unix_milliseconds) ||
			Number.isNaN(incoming_captured_at_unix_milliseconds)
		) {
			return true;
		}

		return incoming_captured_at_unix_milliseconds > existing_captured_at_unix_milliseconds;
	}

	function sync_positions_list_from_map() {
		latest_positions_list = Array.from(latest_positions_by_truck_identifier.values()).sort(
			(left_item, right_item) =>
				left_item.truck_identifier.localeCompare(right_item.truck_identifier)
		);
	}

	function recompute_static_feature_collections() {
		smart_bin_feature_collection = build_smart_bin_feature_collection();
		acopio_feature_collection = build_acopio_feature_collection();
	}

	function recompute_route_feature_collections() {
		collection_route_line_feature_collection = build_collection_route_line_feature_collection();
		collection_route_stop_feature_collection = build_collection_route_stop_feature_collection();
	}

	function resolve_motion_duration_milliseconds(distance_meters: number): number {
		const projected_duration = distance_meters * truck_motion_milliseconds_per_meter;
		return Math.max(
			truck_motion_min_duration_milliseconds,
			Math.min(truck_motion_max_duration_milliseconds, projected_duration)
		);
	}

	function smoothstep_unit_interval(raw_value: number): number {
		const clamped_value = Math.max(0, Math.min(1, raw_value));
		return clamped_value * clamped_value * (3 - 2 * clamped_value);
	}

	function lerp_coordinates(
		start_coordinates: [number, number],
		end_coordinates: [number, number],
		t: number
	): [number, number] {
		return [
			start_coordinates[0] + (end_coordinates[0] - start_coordinates[0]) * t,
			start_coordinates[1] + (end_coordinates[1] - start_coordinates[1]) * t
		];
	}

	function is_same_coordinates(
		left_coordinates: [number, number],
		right_coordinates: [number, number]
	): boolean {
		return (
			Math.abs(left_coordinates[0] - right_coordinates[0]) < 1e-7 &&
			Math.abs(left_coordinates[1] - right_coordinates[1]) < 1e-7
		);
	}

	function stop_truck_motion_loop() {
		if (truck_motion_frame_request_identifier == null || typeof window === 'undefined') {
			truck_motion_frame_request_identifier = null;
			return;
		}
		window.cancelAnimationFrame(truck_motion_frame_request_identifier);
		truck_motion_frame_request_identifier = null;
	}

	function ensure_truck_motion_loop() {
		if (typeof window === 'undefined') {
			return;
		}
		if (truck_motion_frame_request_identifier != null) {
			return;
		}
		truck_motion_frame_request_identifier = window.requestAnimationFrame(
			run_truck_motion_animation_frame
		);
	}

	function run_truck_motion_animation_frame(frame_timestamp: number) {
		truck_motion_frame_request_identifier = null;
		if (active_motion_by_truck_identifier.size === 0) {
			return;
		}

		let has_pending_motions = false;
		for (const [truck_identifier, motion_state] of active_motion_by_truck_identifier) {
			const elapsed_milliseconds = frame_timestamp - motion_state.started_at_unix_milliseconds;
			const raw_progress =
				motion_state.duration_milliseconds <= 0
					? 1
					: elapsed_milliseconds / motion_state.duration_milliseconds;
			const eased_progress = smoothstep_unit_interval(raw_progress);
			const next_coordinates = lerp_coordinates(
				motion_state.start_coordinates,
				motion_state.end_coordinates,
				eased_progress
			);
			rendered_coordinates_by_truck_identifier.set(truck_identifier, next_coordinates);

			if (raw_progress >= 1) {
				rendered_coordinates_by_truck_identifier.set(
					truck_identifier,
					motion_state.end_coordinates
				);
				active_motion_by_truck_identifier.delete(truck_identifier);
			} else {
				has_pending_motions = true;
			}
		}

		truck_feature_collection = build_truck_feature_collection();
		refresh_active_truck_popup_if_needed();

		if (has_pending_motions && typeof window !== 'undefined') {
			truck_motion_frame_request_identifier = window.requestAnimationFrame(
				run_truck_motion_animation_frame
			);
		}
	}

	function get_target_truck_coordinates(
		truck_position_item: Truck_latest_position,
		route_line_candidate_list: Route_line_candidate[]
	): [number, number] {
		const snapped_truck_coordinates = resolve_truck_coordinates_snapped_to_route_candidates(
			truck_position_item,
			route_line_candidate_list
		);
		return (
			snapped_truck_coordinates ?? [truck_position_item.longitude, truck_position_item.latitude]
		);
	}

	function update_rendered_truck_coordinates(animate: boolean) {
		const route_line_candidate_list = build_route_line_candidate_list();
		const seen_truck_identifiers = new Set<string>();

		const now_milliseconds = typeof performance !== 'undefined' ? performance.now() : Date.now();

		for (const truck_position_item of latest_positions_list) {
			const truck_identifier = truck_position_item.truck_identifier;
			seen_truck_identifiers.add(truck_identifier);

			const next_target_coordinates = get_target_truck_coordinates(
				truck_position_item,
				route_line_candidate_list
			);
			const current_rendered_coordinates =
				rendered_coordinates_by_truck_identifier.get(truck_identifier);

			if (!current_rendered_coordinates || !animate) {
				rendered_coordinates_by_truck_identifier.set(truck_identifier, next_target_coordinates);
				active_motion_by_truck_identifier.delete(truck_identifier);
				continue;
			}

			const distance_to_target_meters = calculate_distance_meters_between_coordinates(
				current_rendered_coordinates,
				next_target_coordinates
			);
			if (
				distance_to_target_meters <= 0.7 ||
				distance_to_target_meters >= truck_motion_instant_jump_threshold_meters
			) {
				rendered_coordinates_by_truck_identifier.set(truck_identifier, next_target_coordinates);
				active_motion_by_truck_identifier.delete(truck_identifier);
				continue;
			}

			const existing_motion = active_motion_by_truck_identifier.get(truck_identifier);
			if (
				existing_motion &&
				is_same_coordinates(existing_motion.end_coordinates, next_target_coordinates)
			) {
				continue;
			}

			active_motion_by_truck_identifier.set(truck_identifier, {
				start_coordinates: current_rendered_coordinates,
				end_coordinates: next_target_coordinates,
				started_at_unix_milliseconds: now_milliseconds,
				duration_milliseconds: resolve_motion_duration_milliseconds(distance_to_target_meters)
			});
		}

		const known_identifiers = Array.from(rendered_coordinates_by_truck_identifier.keys());
		for (const truck_identifier of known_identifiers) {
			if (seen_truck_identifiers.has(truck_identifier)) {
				continue;
			}
			rendered_coordinates_by_truck_identifier.delete(truck_identifier);
			active_motion_by_truck_identifier.delete(truck_identifier);
			latest_snap_state_by_truck_identifier.delete(truck_identifier);
		}

		if (active_motion_by_truck_identifier.size > 0) {
			ensure_truck_motion_loop();
		} else {
			stop_truck_motion_loop();
		}
	}

	function recompute_truck_feature_collection(options: { animate: boolean }) {
		update_rendered_truck_coordinates(options.animate);
		truck_feature_collection = build_truck_feature_collection();
		refresh_active_truck_popup_if_needed();
	}

	function refresh_active_truck_popup_if_needed() {
		if (!active_map_popup_state?.associated_truck_identifier) {
			return;
		}

		const tracked_truck_identifier = active_map_popup_state.associated_truck_identifier;
		const latest_truck_position =
			latest_positions_by_truck_identifier.get(tracked_truck_identifier);
		if (!latest_truck_position) {
			active_map_popup_state = null;
			return;
		}

		const rendered_coordinates = rendered_coordinates_by_truck_identifier.get(
			tracked_truck_identifier
		) ?? [latest_truck_position.longitude, latest_truck_position.latitude];
		active_map_popup_state = {
			...active_map_popup_state,
			lnglat: rendered_coordinates,
			popup_html: build_truck_marker_popup_html(latest_truck_position)
		};
	}

	function set_active_popup_state(
		lnglat: [number, number],
		popup_html: string,
		associated_truck_identifier: string | null = null
	) {
		active_map_popup_state = {
			lnglat,
			popup_html,
			associated_truck_identifier,
			open: true
		};
	}

	function clear_active_popup_state() {
		active_map_popup_state = null;
	}

	function on_layer_pointer_enter() {
		map_canvas_cursor = 'pointer';
	}

	function on_layer_pointer_leave() {
		map_canvas_cursor = '';
	}

	function on_truck_layer_click(layer_mouse_event: maplibregl.MapLayerMouseEvent) {
		const clicked_feature = layer_mouse_event.features?.[0];
		if (!clicked_feature || clicked_feature.geometry.type !== 'Point') {
			return;
		}
		const clicked_coordinates = clicked_feature.geometry.coordinates as [number, number];
		const popup_html = String(clicked_feature.properties?.['popup_html'] ?? '');
		const truck_identifier_value = String(clicked_feature.properties?.['truck_identifier'] ?? '');
		set_active_popup_state(
			clicked_coordinates,
			popup_html,
			truck_identifier_value.length > 0 ? truck_identifier_value : null
		);
	}

	function on_collection_route_layer_click(layer_mouse_event: maplibregl.MapLayerMouseEvent) {
		const clicked_feature = layer_mouse_event.features?.[0];
		if (!clicked_feature || clicked_feature.geometry.type !== 'LineString') {
			return;
		}
		const clicked_coordinates = layer_mouse_event.lngLat.toArray() as [number, number];
		const popup_html = String(clicked_feature.properties?.['popup_html'] ?? '');
		set_active_popup_state(clicked_coordinates, popup_html, null);
	}

	function on_route_stop_layer_click(layer_mouse_event: maplibregl.MapLayerMouseEvent) {
		const clicked_feature = layer_mouse_event.features?.[0];
		if (!clicked_feature || clicked_feature.geometry.type !== 'Point') {
			return;
		}
		const clicked_coordinates = clicked_feature.geometry.coordinates as [number, number];
		const popup_html = String(clicked_feature.properties?.['popup_html'] ?? '');
		set_active_popup_state(clicked_coordinates, popup_html, null);
	}

	function on_smart_bin_layer_click(layer_mouse_event: maplibregl.MapLayerMouseEvent) {
		const clicked_feature = layer_mouse_event.features?.[0];
		if (!clicked_feature || clicked_feature.geometry.type !== 'Point') {
			return;
		}
		const clicked_coordinates = clicked_feature.geometry.coordinates as [number, number];
		const popup_html = String(clicked_feature.properties?.['popup_html'] ?? '');
		set_active_popup_state(clicked_coordinates, popup_html, null);
	}

	function on_acopio_layer_click(layer_mouse_event: maplibregl.MapLayerMouseEvent) {
		const clicked_feature = layer_mouse_event.features?.[0];
		if (!clicked_feature || clicked_feature.geometry.type !== 'Point') {
			return;
		}
		const clicked_coordinates = clicked_feature.geometry.coordinates as [number, number];
		const popup_html = String(clicked_feature.properties?.['popup_html'] ?? '');
		set_active_popup_state(clicked_coordinates, popup_html, null);
	}

	function build_truck_feature_collection(): Truck_feature_collection {
		return {
			type: 'FeatureCollection',
			features: latest_positions_list.map((truck_position_item) => {
				const display_truck_coordinates = rendered_coordinates_by_truck_identifier.get(
					truck_position_item.truck_identifier
				) ?? [truck_position_item.longitude, truck_position_item.latitude];

				return {
					type: 'Feature',
					geometry: {
						type: 'Point',
						coordinates: display_truck_coordinates as [number, number]
					},
					properties: {
						truck_identifier: truck_position_item.truck_identifier,
						latitude: display_truck_coordinates[1],
						longitude: display_truck_coordinates[0],
						speed_kmh: truck_position_item.speed_kmh ?? null,
						captured_at: truck_position_item.captured_at,
						popup_html: build_truck_marker_popup_html(truck_position_item)
					}
				};
			})
		};
	}

	function build_collection_route_line_feature_collection(): Collection_route_line_feature_collection {
		const line_feature_list: Collection_route_line_feature[] = [];

		for (const collection_route_item of collection_routes_list) {
			const resolved_line_coordinates =
				resolve_collection_route_line_coordinates(collection_route_item);
			if (resolved_line_coordinates.length < 2) {
				continue;
			}

			line_feature_list.push({
				type: 'Feature',
				geometry: {
					type: 'LineString',
					coordinates: resolved_line_coordinates
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

	function resolve_collection_route_line_coordinates(
		collection_route_item: Collection_route_view
	): [number, number][] {
		const sorted_path_coordinates = [...collection_route_item.path_coordinates].sort(
			(left_coordinate, right_coordinate) =>
				left_coordinate.stop_order - right_coordinate.stop_order
		);
		const fallback_line_coordinates = sorted_path_coordinates.map((route_coordinate) => [
			route_coordinate.longitude,
			route_coordinate.latitude
		]) as [number, number][];
		const persisted_road_path_coordinates = collection_route_item.road_path_coordinates ?? [];
		return persisted_road_path_coordinates.length >= 2
			? persisted_road_path_coordinates
			: fallback_line_coordinates;
	}

	function build_route_line_candidate_list(): Route_line_candidate[] {
		const route_line_candidate_list: Route_line_candidate[] = [];
		for (const collection_route_item of collection_routes_list) {
			const resolved_line_coordinates =
				resolve_collection_route_line_coordinates(collection_route_item);
			if (resolved_line_coordinates.length < 2) {
				continue;
			}

			route_line_candidate_list.push({
				route_identifier: collection_route_item.route_identifier,
				line_coordinates: resolved_line_coordinates
			});
		}

		return route_line_candidate_list;
	}

	function resolve_truck_coordinates_snapped_to_route_candidates(
		truck_position_item: Truck_latest_position,
		route_line_candidate_list: Route_line_candidate[]
	): [number, number] | null {
		if (route_line_candidate_list.length === 0) {
			latest_snap_state_by_truck_identifier.delete(truck_position_item.truck_identifier);
			return null;
		}

		const truck_coordinates: [number, number] = [
			truck_position_item.longitude,
			truck_position_item.latitude
		];
		const best_global_projection_match = find_best_route_projection_match(
			truck_coordinates,
			route_line_candidate_list
		);
		if (!best_global_projection_match) {
			latest_snap_state_by_truck_identifier.delete(truck_position_item.truck_identifier);
			return null;
		}
		if (best_global_projection_match.distance_meters > truck_route_snap_distance_threshold_meters) {
			latest_snap_state_by_truck_identifier.delete(truck_position_item.truck_identifier);
			return null;
		}

		let selected_projection_match = best_global_projection_match;
		const previous_snap_state = latest_snap_state_by_truck_identifier.get(
			truck_position_item.truck_identifier
		);
		if (previous_snap_state) {
			const same_route_candidate = route_line_candidate_list.find(
				(route_line_candidate) =>
					route_line_candidate.route_identifier === previous_snap_state.route_identifier
			);
			if (same_route_candidate) {
				const local_projection_match = find_best_route_projection_match(
					truck_coordinates,
					[same_route_candidate],
					previous_snap_state.segment_index - truck_route_snap_segment_window_radius,
					previous_snap_state.segment_index + truck_route_snap_segment_window_radius
				);
				if (
					local_projection_match &&
					local_projection_match.distance_meters <= truck_route_snap_distance_threshold_meters
				) {
					selected_projection_match = local_projection_match;
				}
			}

			const snapped_jump_distance_meters = calculate_distance_meters_between_coordinates(
				previous_snap_state.snapped_coordinates,
				selected_projection_match.snapped_coordinates
			);
			if (snapped_jump_distance_meters > truck_route_snap_max_jump_meters) {
				latest_snap_state_by_truck_identifier.delete(truck_position_item.truck_identifier);
				return null;
			}
		}

		const captured_at_unix_milliseconds = Date.parse(truck_position_item.captured_at);
		latest_snap_state_by_truck_identifier.set(truck_position_item.truck_identifier, {
			route_identifier: selected_projection_match.route_identifier,
			segment_index: selected_projection_match.segment_index,
			snapped_coordinates: selected_projection_match.snapped_coordinates,
			raw_coordinates: truck_coordinates,
			captured_at_unix_milliseconds: Number.isNaN(captured_at_unix_milliseconds)
				? Date.now()
				: captured_at_unix_milliseconds
		});

		return selected_projection_match.snapped_coordinates;
	}

	function find_best_route_projection_match(
		point_coordinates: [number, number],
		route_line_candidate_list: Route_line_candidate[],
		start_segment_index_inclusive: number = 0,
		end_segment_index_inclusive: number = Number.POSITIVE_INFINITY
	): Route_projection_match | null {
		let best_projection_match: Route_projection_match | null = null;

		for (const route_line_candidate of route_line_candidate_list) {
			const route_segment_count = route_line_candidate.line_coordinates.length - 1;
			if (route_segment_count <= 0) {
				continue;
			}

			const bounded_segment_start_index = Math.max(0, start_segment_index_inclusive);
			const bounded_segment_end_index = Math.min(
				route_segment_count - 1,
				end_segment_index_inclusive
			);
			for (
				let segment_index = bounded_segment_start_index;
				segment_index <= bounded_segment_end_index;
				segment_index += 1
			) {
				const segment_start_coordinates = route_line_candidate.line_coordinates[segment_index];
				const segment_end_coordinates = route_line_candidate.line_coordinates[segment_index + 1];
				const segment_projection_result = project_point_to_route_segment(
					point_coordinates,
					segment_start_coordinates,
					segment_end_coordinates
				);

				if (
					!best_projection_match ||
					segment_projection_result.distance_meters < best_projection_match.distance_meters
				) {
					best_projection_match = {
						route_identifier: route_line_candidate.route_identifier,
						segment_index,
						snapped_coordinates: segment_projection_result.snapped_coordinates,
						distance_meters: segment_projection_result.distance_meters
					};
				}
			}
		}

		return best_projection_match;
	}

	function project_point_to_route_segment(
		point_coordinates: [number, number],
		segment_start_coordinates: [number, number],
		segment_end_coordinates: [number, number]
	): { snapped_coordinates: [number, number]; distance_meters: number } {
		const point_longitude = point_coordinates[0];
		const point_latitude = point_coordinates[1];
		const segment_start_longitude = segment_start_coordinates[0];
		const segment_start_latitude = segment_start_coordinates[1];
		const segment_end_longitude = segment_end_coordinates[0];
		const segment_end_latitude = segment_end_coordinates[1];

		const reference_latitude_radians =
			((point_latitude + segment_start_latitude + segment_end_latitude) / 3) * (Math.PI / 180);
		const meters_per_degree_longitude = 111320 * Math.cos(reference_latitude_radians);

		const segment_delta_x_meters =
			(segment_end_longitude - segment_start_longitude) * meters_per_degree_longitude;
		const segment_delta_y_meters =
			(segment_end_latitude - segment_start_latitude) * meters_per_degree_latitude;
		const point_delta_x_meters =
			(point_longitude - segment_start_longitude) * meters_per_degree_longitude;
		const point_delta_y_meters =
			(point_latitude - segment_start_latitude) * meters_per_degree_latitude;

		const segment_length_squared =
			segment_delta_x_meters * segment_delta_x_meters +
			segment_delta_y_meters * segment_delta_y_meters;
		const raw_projection_ratio =
			segment_length_squared <= 0
				? 0
				: (point_delta_x_meters * segment_delta_x_meters +
						point_delta_y_meters * segment_delta_y_meters) /
					segment_length_squared;
		const normalized_projection_ratio = Math.max(0, Math.min(1, raw_projection_ratio));

		const snapped_longitude =
			segment_start_longitude +
			(segment_end_longitude - segment_start_longitude) * normalized_projection_ratio;
		const snapped_latitude =
			segment_start_latitude +
			(segment_end_latitude - segment_start_latitude) * normalized_projection_ratio;

		const snapped_delta_x_meters =
			(point_longitude - snapped_longitude) * meters_per_degree_longitude;
		const snapped_delta_y_meters = (point_latitude - snapped_latitude) * meters_per_degree_latitude;
		const snapped_distance_meters = Math.sqrt(
			snapped_delta_x_meters * snapped_delta_x_meters +
				snapped_delta_y_meters * snapped_delta_y_meters
		);

		return {
			snapped_coordinates: [snapped_longitude, snapped_latitude],
			distance_meters: snapped_distance_meters
		};
	}

	function calculate_distance_meters_between_coordinates(
		left_coordinates: [number, number],
		right_coordinates: [number, number]
	): number {
		const left_longitude = left_coordinates[0];
		const left_latitude = left_coordinates[1];
		const right_longitude = right_coordinates[0];
		const right_latitude = right_coordinates[1];

		const average_latitude_radians = ((left_latitude + right_latitude) / 2) * (Math.PI / 180);
		const meters_per_degree_longitude = 111320 * Math.cos(average_latitude_radians);
		const delta_longitude_meters = (right_longitude - left_longitude) * meters_per_degree_longitude;
		const delta_latitude_meters = (right_latitude - left_latitude) * meters_per_degree_latitude;
		return Math.sqrt(
			delta_longitude_meters * delta_longitude_meters +
				delta_latitude_meters * delta_latitude_meters
		);
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

	function build_smart_bin_feature_collection(): GeoJSON.FeatureCollection {
		return {
			type: 'FeatureCollection',
			features: mock_smart_bins.map((bin) => ({
				type: 'Feature',
				geometry: { type: 'Point', coordinates: [bin.lng, bin.lat] },
				properties: {
					id: bin.id,
					label: bin.label,
					zone: bin.zone,
					capacity_pct: bin.capacity_pct,
					last_updated: bin.last_updated,
					popup_html: build_smart_bin_popup_html(bin)
				}
			}))
		};
	}

	function build_acopio_feature_collection(): GeoJSON.FeatureCollection {
		return {
			type: 'FeatureCollection',
			features: mock_acopio_points.map((pt) => ({
				type: 'Feature',
				geometry: { type: 'Point', coordinates: [pt.lng, pt.lat] },
				properties: {
					id: pt.id,
					name: pt.name,
					zone: pt.zone,
					capacity_pct: pt.capacity_pct,
					popup_html: build_acopio_popup_html(pt)
				}
			}))
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
			hash_accumulator = (hash_accumulator * 31 + raw_text.charCodeAt(character_index)) >>> 0;
		}
		return hash_accumulator;
	}

	function popup_card(icon: string, title: string, rows: string, cta?: string): string {
		return `<div style="font-family:system-ui,sans-serif;padding:2px 0">
			<div style="display:flex;align-items:center;gap:8px;margin-bottom:10px">
				<span style="font-size:1.15rem">${icon}</span>
				<strong style="font-size:0.9rem;color:#0f172a">${title}</strong>
			</div>
			<div style="display:grid;gap:4px;font-size:0.8rem;color:#475569">${rows}</div>
			${cta ? `<div style="margin-top:10px">${cta}</div>` : ''}
		</div>`;
	}

	function popup_row(label: string, value: string): string {
		return `<div style="display:flex;justify-content:space-between;gap:12px">
			<span style="color:#94a3b8">${label}</span>
			<span style="font-weight:600;color:#1e293b">${value}</span>
		</div>`;
	}

	function popup_badge(text: string, color: string, bg: string): string {
		return `<span style="display:inline-block;padding:2px 8px;border-radius:999px;font-size:0.72rem;font-weight:700;color:${color};background:${bg}">${text}</span>`;
	}

	function build_truck_marker_popup_html(truck_position_item: Truck_latest_position): string {
		const speed_label =
			truck_position_item.speed_kmh == null ? 'N/A' : `${truck_position_item.speed_kmh} km/h`;
		const rows = [
			popup_row('ID', truck_position_item.truck_identifier),
			popup_row('Velocidad', speed_label),
			popup_row(
				'Lat / Lng',
				`${truck_position_item.latitude.toFixed(4)}, ${truck_position_item.longitude.toFixed(4)}`
			),
			popup_row('Actualizado', truck_position_item.captured_at)
		].join('');
		return popup_card('🚛', 'Carro recolector', rows);
	}

	function build_collection_route_popup_html(collection_route_item: Collection_route_view): string {
		const days = ['Dom', 'Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb'];
		const day_label =
			days[collection_route_item.collection_weekday] ??
			String(collection_route_item.collection_weekday);
		const rows = [
			popup_row('Código', collection_route_item.route_code),
			popup_row('Zona', collection_route_item.zone_name),
			popup_row('Día de recolección', day_label),
			popup_row('Paradas', String(collection_route_item.stop_total))
		].join('');
		return popup_card('🗺️', collection_route_item.route_name, rows);
	}

	function build_collection_route_stop_popup_html(
		collection_route_item: Collection_route_view,
		route_stop_coordinate: Route_path_coordinate
	): string {
		const rows = [
			popup_row('Ruta', collection_route_item.route_code),
			popup_row('Parada', `#${route_stop_coordinate.stop_order}`),
			popup_row('Contenedor', route_stop_coordinate.bin_code)
		].join('');
		return popup_card('📍', `Parada de ruta`, rows);
	}

	function build_smart_bin_popup_html(bin: SmartBin): string {
		const cap = bin.capacity_pct;
		const status_color = cap >= 80 ? '#ef4444' : cap >= 60 ? '#f59e0b' : '#16a34a';
		const status_bg = cap >= 80 ? '#fef2f2' : cap >= 60 ? '#fffbeb' : '#f0fdf4';
		const status_label = cap >= 80 ? 'Saturado' : cap >= 60 ? 'Moderado' : 'Disponible';
		const bar = `<div style="margin:8px 0 4px">
			<div style="display:flex;justify-content:space-between;font-size:0.72rem;margin-bottom:3px">
				<span style="color:#94a3b8">Capacidad</span>
				<span style="font-weight:700;color:${status_color}">${cap}%</span>
			</div>
			<div style="height:6px;border-radius:999px;background:#e2e8f0;overflow:hidden">
				<div style="height:100%;width:${cap}%;background:${status_color};border-radius:999px"></div>
			</div>
		</div>`;
		const rows = [popup_row('Zona', bin.zone), popup_row('Actualizado', bin.last_updated)].join('');
		const badge = popup_badge(status_label, status_color, status_bg);
		return popup_card('🗑️', bin.label, `${badge}${bar}${rows}`);
	}

	function build_acopio_popup_html(pt: AcopioMarker): string {
		const cap = pt.capacity_pct;
		const cap_color = cap >= 85 ? '#ef4444' : cap >= 60 ? '#f59e0b' : '#16a34a';
		const bar = `<div style="margin:8px 0 4px">
			<div style="display:flex;justify-content:space-between;font-size:0.72rem;margin-bottom:3px">
				<span style="color:#94a3b8">Capacidad</span>
				<span style="font-weight:700;color:${cap_color}">${cap}%</span>
			</div>
			<div style="height:6px;border-radius:999px;background:#e2e8f0;overflow:hidden">
				<div style="height:100%;width:${cap}%;background:${cap_color};border-radius:999px"></div>
			</div>
		</div>`;
		const rows = [
			popup_row('Dirección', pt.address),
			popup_row('Zona', pt.zone),
			popup_row('Horario', pt.schedule),
			popup_row('Materiales', pt.materials)
		].join('');
		const cta = `<a href="/recycling" style="display:block;text-align:center;padding:6px;border-radius:8px;background:#0284c7;color:#fff;font-size:0.78rem;font-weight:600;text-decoration:none">Ver más detalles →</a>`;
		return popup_card('♻️', pt.name, `${bar}${rows}`, cta);
	}

	function build_user_popup_html(): string {
		return popup_card(
			'📍',
			'Estás aquí',
			[popup_row('Zona', 'Zona Norte'), popup_row('Acopio más cercano', 'Av. América 2345')].join(
				''
			)
		);
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

{#snippet user_marker_content()}
	<div class="user_marker_shell" aria-hidden="true">
		<span class="user_marker_pulse"></span>
		<span class="user_marker_dot"></span>
	</div>
{/snippet}

<MapLibre
	bind:map={map_instance}
	class="map_fullscreen"
	style={empty_map_style}
	center={[cochabamba_center_longitude, cochabamba_center_latitude]}
	zoom={13}
	attributionControl={false}
	cursor={map_canvas_cursor}
	autoloadGlobalCss={false}
>
	<RasterTileSource
		id="base_tiles"
		tiles={[base_tile_url]}
		tileSize={256}
		attribution="&copy; OpenStreetMap contributors &copy; CARTO"
	>
		<RasterLayer id="base_layer" minzoom={0} maxzoom={22} />
	</RasterTileSource>

	<GeoJSONSource
		id={collection_routes_source_identifier}
		data={collection_route_line_feature_collection}
	>
		<LineLayer
			id={collection_routes_line_layer_identifier}
			paint={{
				'line-color': ['coalesce', ['get', 'line_color'], '#0284c7'],
				'line-width': 4,
				'line-opacity': 0.85
			}}
			layout={{
				'line-join': 'round',
				'line-cap': 'round'
			}}
			onmouseenter={on_layer_pointer_enter}
			onmouseleave={on_layer_pointer_leave}
			onclick={on_collection_route_layer_click}
		/>
	</GeoJSONSource>

	<GeoJSONSource
		id={collection_route_stops_source_identifier}
		data={collection_route_stop_feature_collection}
	>
		<CircleLayer
			id={collection_route_stops_circle_layer_identifier}
			paint={{
				'circle-color': '#f8fafc',
				'circle-radius': 4.2,
				'circle-stroke-color': '#0f172a',
				'circle-stroke-width': 1.2,
				'circle-opacity': 0.95
			}}
			onmouseenter={on_layer_pointer_enter}
			onmouseleave={on_layer_pointer_leave}
			onclick={on_route_stop_layer_click}
		/>
	</GeoJSONSource>

	<GeoJSONSource id={truck_positions_source_identifier} data={truck_feature_collection}>
		<CircleLayer
			id={truck_positions_circle_layer_identifier}
			paint={{
				'circle-color': '#18a26e',
				'circle-radius': 7,
				'circle-stroke-color': '#0a5e4a',
				'circle-stroke-width': 2.2,
				'circle-opacity': 0.95
			}}
			onmouseenter={on_layer_pointer_enter}
			onmouseleave={on_layer_pointer_leave}
			onclick={on_truck_layer_click}
		/>
	</GeoJSONSource>

	<GeoJSONSource id={smart_bins_source_identifier} data={smart_bin_feature_collection}>
		<CircleLayer
			id={smart_bins_circle_layer_identifier}
			paint={{
				'circle-color': [
					'case',
					['>=', ['get', 'capacity_pct'], 80],
					'#ef4444',
					['>=', ['get', 'capacity_pct'], 60],
					'#f59e0b',
					'#16a34a'
				],
				'circle-radius': 6,
				'circle-stroke-color': [
					'case',
					['>=', ['get', 'capacity_pct'], 80],
					'#991b1b',
					['>=', ['get', 'capacity_pct'], 60],
					'#92400e',
					'#14532d'
				],
				'circle-stroke-width': 1.8,
				'circle-opacity': 0.92
			}}
			onmouseenter={on_layer_pointer_enter}
			onmouseleave={on_layer_pointer_leave}
			onclick={on_smart_bin_layer_click}
		/>
	</GeoJSONSource>

	<GeoJSONSource id={acopio_points_source_identifier} data={acopio_feature_collection}>
		<CircleLayer
			id={acopio_points_circle_layer_identifier}
			paint={{
				'circle-color': '#0284c7',
				'circle-radius': 8,
				'circle-stroke-color': '#075985',
				'circle-stroke-width': 2,
				'circle-opacity': 0.9
			}}
			onmouseenter={on_layer_pointer_enter}
			onmouseleave={on_layer_pointer_leave}
			onclick={on_acopio_layer_click}
		/>
	</GeoJSONSource>

	<Marker lnglat={user_marker_coordinates} content={user_marker_content}>
		<Popup offset={20} maxWidth="220px">
			{@html build_user_popup_html()}
		</Popup>
	</Marker>

	{#if active_map_popup_state}
		<Popup
			lnglat={active_map_popup_state.lnglat}
			bind:open={active_map_popup_state.open}
			closeButton={true}
			closeOnClick={true}
			offset={16}
			maxWidth="280px"
			onclose={clear_active_popup_state}
		>
			{@html active_map_popup_state.popup_html}
		</Popup>
	{/if}
</MapLibre>

<div class="map_status_panel">
	<div>API: {backend_api_base_url}</div>
	<div>Rutas activas: {collection_routes_list.length}</div>
	<div>Camiones en vivo: {latest_positions_list.length}</div>
</div>

<div class="map_legend">
	<p class="legend_title">Leyenda</p>
	<div class="legend_items">
		<div class="legend_row">
			<span class="legend_dot" style="background:#18a26e;border-color:#0a5e4a"></span>
			<span>Carro recolector</span>
		</div>
		<div class="legend_row">
			<span class="legend_dot" style="background:#16a34a;border-color:#14532d"></span>
			<span>Contenedor OK</span>
		</div>
		<div class="legend_row">
			<span class="legend_dot" style="background:#f59e0b;border-color:#92400e"></span>
			<span>Contenedor moderado</span>
		</div>
		<div class="legend_row">
			<span class="legend_dot" style="background:#ef4444;border-color:#991b1b"></span>
			<span>Contenedor saturado</span>
		</div>
		<div class="legend_row">
			<span class="legend_dot" style="background:#0284c7;border-color:#075985"></span>
			<span>Punto de acopio</span>
		</div>
		<div class="legend_row">
			<span class="legend_dot" style="background:#3b82f6;border-color:#fff"></span>
			<span>Estás aquí</span>
		</div>
		<div class="legend_row">
			<span class="legend_line"></span>
			<span>Ruta de recolección</span>
		</div>
	</div>
</div>

<style>
	:global(.map_fullscreen) {
		position: fixed;
		top: 0;
		left: 0;
		width: 100vw;
		height: 100vh;
		z-index: 0;
	}

	.user_marker_shell {
		position: relative;
		width: 20px;
		height: 20px;
		cursor: default;
	}

	.user_marker_pulse {
		position: absolute;
		top: 50%;
		left: 50%;
		width: 36px;
		height: 36px;
		margin: -18px 0 0 -18px;
		border-radius: 50%;
		background: rgba(59, 130, 246, 0.25);
		animation: eco_pulse 2s ease-out infinite;
	}

	.user_marker_dot {
		position: absolute;
		top: 50%;
		left: 50%;
		width: 14px;
		height: 14px;
		margin: -7px 0 0 -7px;
		border-radius: 50%;
		background: #3b82f6;
		border: 3px solid #fff;
		box-shadow: 0 2px 8px rgba(59, 130, 246, 0.6);
	}

	@keyframes eco_pulse {
		0% {
			transform: scale(0.6);
			opacity: 0.8;
		}

		70% {
			transform: scale(1.8);
			opacity: 0;
		}

		100% {
			transform: scale(2.2);
			opacity: 0;
		}
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

	.map_legend {
		position: fixed;
		left: 1rem;
		bottom: 5rem;
		z-index: 40;
		padding: 0.75rem 1rem;
		border-radius: 0.9rem;
		background: oklch(0.98 0 0 / 0.88);
		border: 1px solid oklch(0.85 0 0 / 0.6);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
		min-width: 170px;
	}

	:global([data-theme='dark']) .map_legend {
		background: oklch(0.18 0 0 / 0.88);
		border-color: oklch(1 0 0 / 0.12);
	}

	.legend_title {
		margin: 0 0 0.55rem;
		font-size: 0.72rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: #64748b;
	}

	.legend_items {
		display: flex;
		flex-direction: column;
		gap: 0.45rem;
	}

	.legend_row {
		display: flex;
		align-items: center;
		gap: 0.55rem;
		font-size: 0.75rem;
		color: #334155;
	}

	:global([data-theme='dark']) .legend_row {
		color: #cbd5e1;
	}

	.legend_dot {
		display: inline-block;
		width: 11px;
		height: 11px;
		border-radius: 50%;
		border: 2px solid;
		flex-shrink: 0;
	}

	.legend_line {
		display: inline-block;
		width: 18px;
		height: 3px;
		border-radius: 999px;
		background: linear-gradient(to right, #16a34a, #0284c7);
		flex-shrink: 0;
	}

	@media (max-width: 1023px) {
		.map_status_panel {
			top: 5.4rem;
			right: 0.7rem;
			left: 0.7rem;
			max-width: none;
		}

		.map_legend {
			left: 0.7rem;
			bottom: 5.5rem;
		}
	}

	:global(.maplibregl-ctrl-attrib) {
		font-size: 10px;
		background: oklch(1 0 0 / 0.84);
	}

	:global(.maplibregl-popup-content) {
		border-radius: 12px !important;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.18) !important;
		padding: 14px 16px !important;
		font-family: system-ui, sans-serif;
		min-width: 200px;
	}

	:global(.maplibregl-popup-close-button) {
		font-size: 1rem;
		padding: 4px 8px;
		color: #94a3b8;
	}

	:global(.maplibregl-popup-tip) {
		border-top-color: #fff !important;
	}
</style>
