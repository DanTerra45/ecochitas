<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import type { PageData } from '../../../../routes/(authenticated)/$types';
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
	import SimulatorControlPanel from '$lib/components/SimulatorControlPanel.svelte';
	import AppModal from './AppModal.svelte';
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

	type Simulator_status = 'idle' | 'running' | 'paused';

	type Simulator_payload = {
		routeId: string;
		zone: string;
		truckCount: number;
		speedMultiplier: number;
		updateIntervalMs: number;
		snapToRoute: boolean;
	};

	type Simulator_route_option = {
		value: string;
		label: string;
	};

	type Simulator_truck_runtime_state = {
		truck_identifier: string;
		route_coordinates: [number, number][];
		cursor: number;
		speed_kmh: number;
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
	const stream_reconnect_base_delay_milliseconds = 2000;
	const stream_reconnect_max_delay_milliseconds = 30000;
	const stream_reconnect_jitter_ratio = 0.2;
	const stream_status_message_cooldown_milliseconds = 5000;
	const backend_health_poll_interval_milliseconds = 15000;
	const backend_data_refresh_interval_milliseconds = 30000;
	const backend_health_request_timeout_milliseconds = 3500;
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
	const simulated_truck_identifier_prefix = 'SIM-';
	const simulation_backend_start_endpoint_path = '/v1/admin/simulations/trucks/start';

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
	let stream_reconnect_attempt_count = 0;
	let stream_last_status_message_at_milliseconds = 0;
	let backend_health_poll_interval_identifier: number | null = null;
	let backend_data_refresh_interval_identifier: number | null = null;
	let active_map_popup_state = $state<Active_map_popup_state | null>(null);
	let is_simulator_panel_open = $state(false);
	let simulation_status = $state<Simulator_status>('idle');
	let simulation_active_trucks = $state(0);
	let simulation_messages_sent = $state(0);
	let simulation_route_id = $state('route-1');
	let simulation_zone = $state('north-zone');
	let simulation_truck_count = $state(5);
	let simulation_speed_multiplier = $state(1);
	let simulation_update_interval_milliseconds = $state(1000);
	let simulation_snap_to_route = $state(true);
	let simulation_loop_interval_identifier: number | null = null;
	let simulation_runtime_by_truck_identifier = new Map<string, Simulator_truck_runtime_state>();
	let simulated_truck_identifiers = new Set<string>();
	let simulation_backend_control_enabled = $state<null | boolean>(null);
	let simulation_last_message = $state('Local simulator is ready.');
	let backend_connection_status = $state<'unknown' | 'online' | 'offline'>('unknown');
	let backend_last_health_check_iso_timestamp = $state<string | null>(null);
	let backend_last_data_sync_iso_timestamp = $state<string | null>(null);
	let is_manual_reconnect_in_progress = $state(false);

	// --- UI Panel States ---
	let is_legend_open = $state(false);
	let is_status_open = $state(false);

	// --- Global Modal States ---
	let global_modal_open = $state(false);
	let global_modal_title = $state('');
	let global_modal_message = $state('');
	let global_modal_type = $state<'success' | 'error' | 'info'>('info');

	function show_modal(title: string, message: string, type: 'success' | 'error' | 'info' = 'info') {
		global_modal_title = title;
		global_modal_message = message;
		global_modal_type = type;
		global_modal_open = true;
	}

	// --- Route Builder States ---
	let is_route_builder_active = $state(false);
	let route_builder_selected_points = $state<{ id: string; lng: number; lat: number }[]>([]);
	let route_builder_route_code = $state('');
	let route_builder_route_name = $state('');
	let route_builder_zone = $state('Zona Norte');
	let route_builder_weekday = $state(1);
	let is_saving_route = $state(false);

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
	const simulation_route_options = $derived.by(() => {
		if (collection_routes_list.length === 0) {
			return [
				{ value: 'route-1', label: 'Route 1 (Downtown)' },
				{ value: 'route-2', label: 'Route 2 (Suburbs)' },
				{ value: 'route-3', label: 'Route 3 (Highway)' }
			] satisfies Simulator_route_option[];
		}

		return collection_routes_list.map((collection_route_item) => ({
			value: collection_route_item.route_identifier,
			label: `${collection_route_item.route_code} · ${collection_route_item.route_name}`
		}));
	});
	const simulation_zone_options = $derived.by(() => {
		if (collection_routes_list.length === 0) {
			return [
				{ value: 'north-zone', label: 'North Zone' },
				{ value: 'south-zone', label: 'South Zone' },
				{ value: 'east-zone', label: 'East Zone' },
				{ value: 'west-zone', label: 'West Zone' }
			] satisfies Simulator_route_option[];
		}

		const unique_zone_name_set = new Set<string>();
		for (const collection_route_item of collection_routes_list) {
			if (collection_route_item.zone_name.trim().length > 0) {
				unique_zone_name_set.add(collection_route_item.zone_name);
			}
		}

		return Array.from(unique_zone_name_set).map((zone_name) => ({
			value: zone_name,
			label: zone_name
		}));
	});
	const backend_connection_status_label = $derived(
		backend_connection_status === 'online'
			? 'online'
			: backend_connection_status === 'offline'
				? 'offline'
				: 'checking'
	);
	const backend_last_health_check_label = $derived(
		format_status_time_label(backend_last_health_check_iso_timestamp)
	);
	const backend_last_data_sync_label = $derived(
		format_status_time_label(backend_last_data_sync_iso_timestamp)
	);

	onMount(async () => {
		resolve_backend_api_base_url();
		is_dark_map_theme =
			typeof document !== 'undefined' && document.documentElement.dataset.theme === 'dark';
		setup_theme_observer();
		recompute_static_feature_collections();
		await refresh_backend_snapshots({ quiet: true });
		should_keep_stream_connected = true;
		connect_truck_stream();
		start_backend_live_refresh_loops();
		void check_backend_health({ quiet: true });
	});

	onDestroy(() => {
		should_keep_stream_connected = false;
		stop_local_simulation_loop();
		clear_simulated_trucks_from_map();
		stop_truck_motion_loop();
		clear_stream_reconnect_timeout();
		stop_backend_live_refresh_loops();
		disconnect_truck_stream();
		theme_observer?.disconnect();
	});

	$effect(() => {
		if (simulation_route_options.length === 0) return;
		const has_selected_route = simulation_route_options.some(
			(route_option_item) => route_option_item.value === simulation_route_id
		);
		if (!has_selected_route) {
			simulation_route_id = simulation_route_options[0]!.value;
		}
	});

	$effect(() => {
		if (simulation_zone_options.length === 0) return;
		const has_selected_zone = simulation_zone_options.some(
			(zone_option_item) => zone_option_item.value === simulation_zone
		);
		if (!has_selected_zone) {
			simulation_zone = simulation_zone_options[0]!.value;
		}
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

	function format_status_time_label(raw_iso_timestamp: string | null): string {
		if (!raw_iso_timestamp) {
			return 'n/a';
		}
		const parsed_unix_timestamp = Date.parse(raw_iso_timestamp);
		if (Number.isNaN(parsed_unix_timestamp)) {
			return 'n/a';
		}
		return new Date(parsed_unix_timestamp).toLocaleTimeString();
	}

	function set_backend_connection_status(
		next_status: 'unknown' | 'online' | 'offline',
		options: { quiet?: boolean } = {}
	) {
		const previous_status = backend_connection_status;
		backend_connection_status = next_status;
		if (options.quiet || previous_status === next_status) {
			return;
		}

		if (next_status === 'online') {
			maybe_set_stream_status_message('Backend is online.');
			return;
		}

		if (next_status === 'offline') {
			maybe_set_stream_status_message('Backend is offline.');
		}
	}

	function build_health_endpoint_url(): string {
		return `${backend_api_base_url}/healthz`;
	}

	async function check_backend_health(options: { quiet?: boolean } = {}): Promise<boolean> {
		if (typeof window === 'undefined') {
			return false;
		}

		const abort_controller = new AbortController();
		const timeout_identifier = window.setTimeout(() => {
			abort_controller.abort();
		}, backend_health_request_timeout_milliseconds);

		let is_backend_online = false;
		try {
			const health_response = await fetch(build_health_endpoint_url(), {
				cache: 'no-store',
				signal: abort_controller.signal
			});
			is_backend_online = health_response.ok;
		} catch {
			is_backend_online = false;
		} finally {
			window.clearTimeout(timeout_identifier);
			backend_last_health_check_iso_timestamp = new Date().toISOString();
		}

		set_backend_connection_status(is_backend_online ? 'online' : 'offline', {
			quiet: options.quiet
		});
		return is_backend_online;
	}

	async function refresh_backend_snapshots(options: { quiet?: boolean } = {}): Promise<boolean> {
		const [has_latest_positions_snapshot, has_routes_snapshot] = await Promise.all([
			load_latest_positions_snapshot(),
			load_collection_routes_snapshot()
		]);
		const has_any_snapshot = has_latest_positions_snapshot || has_routes_snapshot;
		if (has_any_snapshot) {
			backend_last_data_sync_iso_timestamp = new Date().toISOString();
		}
		if (!options.quiet && !has_any_snapshot) {
			maybe_set_stream_status_message('Snapshot refresh failed.');
		}
		return has_any_snapshot;
	}

	function stop_backend_live_refresh_loops() {
		if (typeof window === 'undefined') {
			return;
		}

		if (backend_health_poll_interval_identifier != null) {
			window.clearInterval(backend_health_poll_interval_identifier);
			backend_health_poll_interval_identifier = null;
		}
		if (backend_data_refresh_interval_identifier != null) {
			window.clearInterval(backend_data_refresh_interval_identifier);
			backend_data_refresh_interval_identifier = null;
		}
	}

	function start_backend_live_refresh_loops() {
		if (typeof window === 'undefined') {
			return;
		}
		stop_backend_live_refresh_loops();

		backend_health_poll_interval_identifier = window.setInterval(() => {
			void check_backend_health({ quiet: true });
		}, backend_health_poll_interval_milliseconds);

		backend_data_refresh_interval_identifier = window.setInterval(() => {
			if (backend_connection_status !== 'online') {
				return;
			}
			void refresh_backend_snapshots({ quiet: true });
		}, backend_data_refresh_interval_milliseconds);
	}

	function normalize_simulator_payload(raw_payload: Simulator_payload): Simulator_payload {
		const normalized_truck_count = Math.max(1, Math.min(20, Math.round(raw_payload.truckCount)));
		const normalized_speed_multiplier = Math.max(
			0.25,
			Math.min(4, Number(raw_payload.speedMultiplier) || 1)
		);
		const normalized_update_interval_milliseconds = Math.max(
			300,
			Math.min(5000, Math.round(Number(raw_payload.updateIntervalMs) || 1000))
		);
		return {
			routeId: String(raw_payload.routeId || '').trim(),
			zone: String(raw_payload.zone || '').trim(),
			truckCount: normalized_truck_count,
			speedMultiplier: normalized_speed_multiplier,
			updateIntervalMs: normalized_update_interval_milliseconds,
			snapToRoute: Boolean(raw_payload.snapToRoute)
		};
	}

	function apply_simulator_payload_to_state(normalized_payload: Simulator_payload) {
		simulation_route_id = normalized_payload.routeId;
		simulation_zone = normalized_payload.zone;
		simulation_truck_count = normalized_payload.truckCount;
		simulation_speed_multiplier = normalized_payload.speedMultiplier;
		simulation_update_interval_milliseconds = normalized_payload.updateIntervalMs;
		simulation_snap_to_route = normalized_payload.snapToRoute;
	}

	function build_fallback_simulation_loop_coordinates(): [number, number][] {
		const fallback_coordinates: [number, number][] = [];
		const step_total = 48;
		for (let step_index = 0; step_index < step_total; step_index += 1) {
			const angle_radians = (Math.PI * 2 * step_index) / step_total;
			const latitude_offset = 0.0115 * Math.sin(angle_radians);
			const longitude_offset = 0.0142 * Math.cos(angle_radians);
			fallback_coordinates.push([
				cochabamba_center_longitude + longitude_offset,
				cochabamba_center_latitude + latitude_offset
			]);
		}
		fallback_coordinates.push(fallback_coordinates[0]!);
		return fallback_coordinates;
	}

	function ensure_loop_route_coordinates(
		route_coordinates: [number, number][]
	): [number, number][] {
		if (route_coordinates.length < 2) {
			return build_fallback_simulation_loop_coordinates();
		}
		const first_coordinates = route_coordinates[0]!;
		const last_coordinates = route_coordinates[route_coordinates.length - 1]!;
		const is_closed_loop = is_same_coordinates(first_coordinates, last_coordinates);
		if (is_closed_loop) {
			return route_coordinates;
		}
		return [...route_coordinates, first_coordinates];
	}

	function resolve_simulation_route_coordinates(
		normalized_payload: Simulator_payload
	): [number, number][] {
		const selected_route_item = collection_routes_list.find(
			(collection_route_item) =>
				collection_route_item.route_identifier === normalized_payload.routeId
		);
		if (selected_route_item) {
			return ensure_loop_route_coordinates(
				resolve_collection_route_line_coordinates(selected_route_item)
			);
		}

		const zone_route_item = collection_routes_list.find(
			(collection_route_item) => collection_route_item.zone_name === normalized_payload.zone
		);
		if (zone_route_item) {
			return ensure_loop_route_coordinates(
				resolve_collection_route_line_coordinates(zone_route_item)
			);
		}

		if (collection_routes_list.length > 0) {
			return ensure_loop_route_coordinates(
				resolve_collection_route_line_coordinates(collection_routes_list[0]!)
			);
		}

		return build_fallback_simulation_loop_coordinates();
	}

	function resolve_coordinates_at_cursor(
		route_coordinates: [number, number][],
		cursor: number
	): [number, number] {
		if (route_coordinates.length === 0) {
			return [cochabamba_center_longitude, cochabamba_center_latitude];
		}
		if (route_coordinates.length === 1) {
			return route_coordinates[0]!;
		}

		const segment_total = route_coordinates.length - 1;
		if (segment_total <= 0) {
			return route_coordinates[0]!;
		}

		const normalized_cursor = ((cursor % segment_total) + segment_total) % segment_total;
		const segment_index = Math.floor(normalized_cursor);
		const segment_progress = normalized_cursor - segment_index;
		const segment_start_coordinates = route_coordinates[segment_index]!;
		const segment_end_coordinates = route_coordinates[segment_index + 1]!;
		return lerp_coordinates(segment_start_coordinates, segment_end_coordinates, segment_progress);
	}

	function calculate_heading_degrees_from_coordinates(
		from_coordinates: [number, number],
		to_coordinates: [number, number]
	): number {
		const delta_longitude = to_coordinates[0] - from_coordinates[0];
		const delta_latitude = to_coordinates[1] - from_coordinates[1];
		if (Math.abs(delta_longitude) < 1e-9 && Math.abs(delta_latitude) < 1e-9) {
			return 0;
		}
		const heading_radians = Math.atan2(delta_longitude, delta_latitude);
		const heading_degrees = (heading_radians * 180) / Math.PI;
		return (heading_degrees + 360) % 360;
	}

	function clear_simulated_trucks_from_map() {
		if (simulated_truck_identifiers.size === 0) {
			return;
		}

		for (const simulated_truck_identifier of simulated_truck_identifiers) {
			latest_positions_by_truck_identifier.delete(simulated_truck_identifier);
			latest_snap_state_by_truck_identifier.delete(simulated_truck_identifier);
			rendered_coordinates_by_truck_identifier.delete(simulated_truck_identifier);
			active_motion_by_truck_identifier.delete(simulated_truck_identifier);
		}
		simulation_runtime_by_truck_identifier.clear();
		simulated_truck_identifiers.clear();
		sync_positions_list_from_map();
		recompute_truck_feature_collection({ animate: false });
	}

	function stop_local_simulation_loop() {
		if (simulation_loop_interval_identifier == null || typeof window === 'undefined') {
			simulation_loop_interval_identifier = null;
			return;
		}
		window.clearInterval(simulation_loop_interval_identifier);
		simulation_loop_interval_identifier = null;
	}

	function run_local_simulation_tick() {
		if (simulation_runtime_by_truck_identifier.size === 0) {
			return;
		}

		const now_iso_timestamp = new Date().toISOString();
		const step_segments =
			Math.max(
				0.06,
				0.85 * simulation_speed_multiplier * (simulation_update_interval_milliseconds / 1000)
			) * (simulation_snap_to_route ? 1 : 1.05);

		for (const simulation_runtime_item of simulation_runtime_by_truck_identifier.values()) {
			const previous_cursor = simulation_runtime_item.cursor;
			simulation_runtime_item.cursor += step_segments;
			const previous_coordinates = resolve_coordinates_at_cursor(
				simulation_runtime_item.route_coordinates,
				previous_cursor
			);
			const next_coordinates = resolve_coordinates_at_cursor(
				simulation_runtime_item.route_coordinates,
				simulation_runtime_item.cursor
			);
			const heading_degrees = calculate_heading_degrees_from_coordinates(
				previous_coordinates,
				next_coordinates
			);

			upsert_truck_position(
				{
					truck_identifier: simulation_runtime_item.truck_identifier,
					latitude: next_coordinates[1],
					longitude: next_coordinates[0],
					speed_kmh: Math.max(
						6,
						simulation_runtime_item.speed_kmh + Math.sin(simulation_runtime_item.cursor) * 1.6
					),
					heading_degrees,
					captured_at: now_iso_timestamp,
					received_at: now_iso_timestamp
				},
				{ defer_sync: true }
			);
		}

		sync_positions_list_from_map();
		recompute_truck_feature_collection({ animate: true });
		simulation_messages_sent += simulation_runtime_by_truck_identifier.size;
	}

	function start_local_simulation_loop() {
		stop_local_simulation_loop();
		if (typeof window === 'undefined') {
			return;
		}

		simulation_loop_interval_identifier = window.setInterval(() => {
			if (simulation_status !== 'running') {
				return;
			}
			run_local_simulation_tick();
		}, simulation_update_interval_milliseconds);
	}

	function start_local_simulation(normalized_payload: Simulator_payload) {
		apply_simulator_payload_to_state(normalized_payload);
		stop_local_simulation_loop();
		clear_simulated_trucks_from_map();

		const resolved_route_coordinates = resolve_simulation_route_coordinates(normalized_payload);
		const effective_truck_total = Math.max(1, normalized_payload.truckCount);
		const route_segment_total = Math.max(1, resolved_route_coordinates.length - 1);

		simulation_runtime_by_truck_identifier.clear();
		simulated_truck_identifiers.clear();
		for (let truck_index = 0; truck_index < effective_truck_total; truck_index += 1) {
			const simulated_truck_identifier = `${simulated_truck_identifier_prefix}${String(truck_index + 1).padStart(3, '0')}`;
			const cursor_offset = (route_segment_total * truck_index) / effective_truck_total;
			const speed_kmh = 18 + normalized_payload.speedMultiplier * 8 + (truck_index % 3) * 1.5;
			simulation_runtime_by_truck_identifier.set(simulated_truck_identifier, {
				truck_identifier: simulated_truck_identifier,
				route_coordinates: resolved_route_coordinates,
				cursor: cursor_offset,
				speed_kmh
			});
			simulated_truck_identifiers.add(simulated_truck_identifier);
		}

		simulation_messages_sent = 0;
		simulation_active_trucks = simulation_runtime_by_truck_identifier.size;
		simulation_status = 'running';
		simulation_last_message = 'Local simulation running.';

		run_local_simulation_tick();
		start_local_simulation_loop();
	}

	function pause_local_simulation() {
		if (simulation_status !== 'running') {
			return;
		}
		stop_local_simulation_loop();
		simulation_status = 'paused';
		simulation_last_message = 'Simulation paused.';
	}

	function resume_local_simulation() {
		if (simulation_status !== 'paused') {
			return;
		}
		simulation_status = 'running';
		simulation_last_message = 'Simulation resumed.';
		start_local_simulation_loop();
	}

	function stop_local_simulation_and_reset_map() {
		stop_local_simulation_loop();
		clear_simulated_trucks_from_map();
		simulation_status = 'idle';
		simulation_active_trucks = 0;
		simulation_messages_sent = 0;
		simulation_last_message = 'Simulation stopped.';
	}

	async function try_start_backend_simulation(
		normalized_payload: Simulator_payload
	): Promise<boolean> {
		try {
			const response = await fetch(
				`${backend_api_base_url}${simulation_backend_start_endpoint_path}`,
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json'
					},
					body: JSON.stringify(normalized_payload)
				}
			);

			if (!response.ok) {
				if (response.status === 404 || response.status === 405) {
					simulation_last_message =
						'Backend simulation API not available. Running local simulator.';
					return false;
				}
				simulation_last_message = `Backend simulation rejected (${response.status}). Running local simulator.`;
				return false;
			}

			simulation_last_message = 'Backend simulation started.';
			return true;
		} catch {
			simulation_last_message = 'Backend simulation unreachable. Running local simulator.';
			return false;
		}
	}

	async function handle_simulation_start(simulation_start_event: CustomEvent<Simulator_payload>) {
		const normalized_payload = normalize_simulator_payload(simulation_start_event.detail);
		apply_simulator_payload_to_state(normalized_payload);

		const backend_started = await try_start_backend_simulation(normalized_payload);
		if (backend_started) {
			stop_local_simulation_loop();
			clear_simulated_trucks_from_map();
			simulation_backend_control_enabled = true;
			simulation_status = 'running';
			simulation_active_trucks = normalized_payload.truckCount;
			return;
		}

		simulation_backend_control_enabled = false;
		start_local_simulation(normalized_payload);
	}

	function handle_simulation_pause() {
		if (simulation_backend_control_enabled) {
			simulation_status = 'paused';
			simulation_last_message = 'Backend simulation paused in UI.';
			return;
		}
		pause_local_simulation();
	}

	function handle_simulation_resume() {
		if (simulation_backend_control_enabled) {
			simulation_status = 'running';
			simulation_last_message = 'Backend simulation resumed in UI.';
			return;
		}
		resume_local_simulation();
	}

	function handle_simulation_stop() {
		stop_local_simulation_and_reset_map();
		simulation_backend_control_enabled = false;
	}

	function handle_simulation_change(simulation_change_event: CustomEvent<Simulator_payload>) {
		const normalized_payload = normalize_simulator_payload(simulation_change_event.detail);
		apply_simulator_payload_to_state(normalized_payload);

		if (simulation_status === 'running' && !simulation_backend_control_enabled) {
			start_local_simulation(normalized_payload);
			simulation_last_message = 'Local simulation settings updated.';
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

	async function load_latest_positions_snapshot(): Promise<boolean> {
		try {
			const snapshot_response = await fetch(build_snapshot_endpoint_url());
			if (!snapshot_response.ok) {
				return false;
			}
			const snapshot_payload = (await snapshot_response.json()) as Snapshot_response_payload;
			replace_positions_from_snapshot(snapshot_payload.items);
			return true;
		} catch {
			// silently ignore — map will populate via stream
			return false;
		}
	}

	async function load_collection_routes_snapshot(): Promise<boolean> {
		try {
			const collection_routes_response = await fetch(build_collection_routes_endpoint_url());
			if (!collection_routes_response.ok) {
				return false;
			}
			const collection_routes_payload =
				(await collection_routes_response.json()) as Collection_routes_response_payload;
			replace_routes_from_snapshot(collection_routes_payload.items);
			return true;
		} catch {
			// silently ignore route load errors in PoC mode
			return false;
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

	function upsert_truck_position(
		truck_position_item: Truck_latest_position,
		options: { defer_sync?: boolean } = {}
	) {
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
		if (options.defer_sync) {
			return;
		}
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

		if (is_route_builder_active) {
			return; // No mostramos popup normal en modo edición
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

	$effect(() => {
		if (map_instance && is_route_builder_active) {
			const active_map_instance = map_instance;
			const click_handler = (e: maplibregl.MapMouseEvent) => {
				route_builder_selected_points.push({
					id: `point_${Date.now()}_${Math.random()}`,
					lng: e.lngLat.lng,
					lat: e.lngLat.lat
				});
			};
			active_map_instance.on('click', click_handler);
			return () => {
				active_map_instance.off('click', click_handler);
			};
		}
	});

	async function save_custom_route() {
		if (route_builder_selected_points.length < 2) {
			show_modal('Alerta', 'Haz clic en el mapa para añadir al menos 2 puntos para crear una ruta.', 'info');
			return;
		}
		if (!route_builder_route_code.trim() || !route_builder_route_name.trim()) {
			show_modal('Alerta', 'Por favor llena el código y nombre de la ruta.', 'info');
			return;
		}
		is_saving_route = true;
		try {
			const res_route = await fetch(`${backend_api_base_url}/v1/admin/demo-routes`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					route_code: route_builder_route_code,
					route_name: route_builder_route_name,
					zone_name: route_builder_zone,
					collection_weekday: route_builder_weekday,
					is_active: true,
					points: route_builder_selected_points.map((p) => ({
						longitude: p.lng,
						latitude: p.lat
					}))
				})
			});
			if (!res_route.ok) {
				const errBody = await res_route.text();
				throw new Error('Error al crear la ruta: ' + errBody);
			}

			show_modal('Éxito', '¡Ruta creada exitosamente!', 'success');
			is_route_builder_active = false;
			route_builder_selected_points = [];
			route_builder_route_code = '';
			route_builder_route_name = '';
			await load_collection_routes_snapshot();
		} catch (error) {
			console.error(error);
			show_modal('Error', error instanceof Error ? error.message : 'Error desconocido', 'error');
		} finally {
			is_saving_route = false;
		}
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
		const is_simulated_truck = truck_position_item.truck_identifier.startsWith(
			simulated_truck_identifier_prefix
		);
		if (is_simulated_truck && !simulation_snap_to_route) {
			latest_snap_state_by_truck_identifier.delete(truck_position_item.truck_identifier);
			return null;
		}

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

	function popup_icon_svg(path_definition: string): string {
		return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round" style="width:18px;height:18px;display:block"><path d="${path_definition}"/></svg>`;
	}

	function popup_card(icon_markup: string, title: string, rows: string, cta?: string): string {
		return `<div style="font-family:system-ui,sans-serif;padding:2px 0">
			<div style="display:flex;align-items:center;gap:8px;margin-bottom:10px">
				<span style="display:inline-flex;color:#22c55e">${icon_markup}</span>
				<strong style="font-size:0.9rem;color:#ffffff">${title}</strong>
			</div>
			<div style="display:grid;gap:4px;font-size:0.8rem;color:#cbd5e1">${rows}</div>
			${cta ? `<div style="margin-top:10px">${cta}</div>` : ''}
		</div>`;
	}

	function popup_row(label: string, value: string): string {
		return `<div style="display:flex;justify-content:space-between;gap:12px">
			<span style="color:#94a3b8">${label}</span>
			<span style="font-weight:600;color:#f8fafc">${value}</span>
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
		return popup_card(
			popup_icon_svg(
				'M10 17h4M1 9h14l4 5v3h-2M1 9v8h2M5 17a2 2 0 1 0 0.001 0M17 17a2 2 0 1 0 0.001 0M7 17h8'
			),
			'Carro recolector',
			rows
		);
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
		return popup_card(
			popup_icon_svg('M3 6l6-2 6 2 6-2v14l-6 2-6-2-6 2zM9 4v14M15 6v14'),
			collection_route_item.route_name,
			rows
		);
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
		return popup_card(
			popup_icon_svg(
				'M12 21s-8-7.5-8-12a8 8 0 0 1 16 0c0 4.5-8 12-8 12zM12 9a2.5 2.5 0 1 0 0.001 0'
			),
			`Parada de ruta`,
			rows
		);
	}

	function build_smart_bin_popup_html(bin: SmartBin): string {
		const cap = bin.capacity_pct;
		const status_color = cap >= 80 ? '#ef4444' : cap >= 60 ? '#f59e0b' : '#22c55e';
		const status_bg =
			cap >= 80
				? 'rgba(239, 68, 68, 0.15)'
				: cap >= 60
					? 'rgba(245, 158, 11, 0.15)'
					: 'rgba(34, 197, 94, 0.15)';
		const status_label = cap >= 80 ? 'Saturado' : cap >= 60 ? 'Moderado' : 'Disponible';
		const bar = `<div style="margin:8px 0 4px">
			<div style="display:flex;justify-content:space-between;font-size:0.72rem;margin-bottom:3px">
				<span style="color:#94a3b8">Capacidad</span>
				<span style="font-weight:700;color:${status_color}">${cap}%</span>
			</div>
			<div style="height:6px;border-radius:999px;background:#262626;overflow:hidden">
				<div style="height:100%;width:${cap}%;background:${status_color};border-radius:999px"></div>
			</div>
		</div>`;
		const rows = [popup_row('Zona', bin.zone), popup_row('Actualizado', bin.last_updated)].join('');
		const badge = popup_badge(status_label, status_color, status_bg);
		return popup_card(
			popup_icon_svg('M3 6h18M8 6V4h8v2M6 6l1 14h10l1-14M10 10v7M14 10v7'),
			bin.label,
			`${badge}${bar}${rows}`
		);
	}

	function build_acopio_popup_html(pt: AcopioMarker): string {
		const cap = pt.capacity_pct;
		const cap_color = cap >= 85 ? '#ef4444' : cap >= 60 ? '#f59e0b' : '#22c55e';
		const bar = `<div style="margin:8px 0 4px">
			<div style="display:flex;justify-content:space-between;font-size:0.72rem;margin-bottom:3px">
				<span style="color:#94a3b8">Capacidad</span>
				<span style="font-weight:700;color:${cap_color}">${cap}%</span>
			</div>
			<div style="height:6px;border-radius:999px;background:#262626;overflow:hidden">
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
		return popup_card(
			popup_icon_svg(
				'M7 19H4.815a1.83 1.83 0 0 1-1.57-.881 1.785 1.785 0 0 1-.004-1.784L7.196 9.5M11 19h8.203a1.83 1.83 0 0 0 1.556-.89 1.784 1.784 0 0 0 0-1.775l-1.226-2.12M14 16l-3 3 3 3M8.293 13.596 7.196 9.5 3.1 10.598M9.344 5.811l1.093-1.892A1.83 1.83 0 0 1 11.985 3a1.784 1.784 0 0 1 1.546.888l3.943 6.843M13.378 9.633l4.096 1.098 1.097-4.096'
			),
			pt.name,
			`${bar}${rows}`,
			cta
		);
	}

	function build_user_popup_html(): string {
		return popup_card(
			popup_icon_svg(
				'M12 21s-8-7.5-8-12a8 8 0 0 1 16 0c0 4.5-8 12-8 12zM12 9a2.5 2.5 0 1 0 0.001 0'
			),
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

	function maybe_set_stream_status_message(next_message: string) {
		const now_milliseconds = Date.now();
		if (
			now_milliseconds - stream_last_status_message_at_milliseconds <
			stream_status_message_cooldown_milliseconds
		) {
			return;
		}
		stream_last_status_message_at_milliseconds = now_milliseconds;
		simulation_last_message = next_message;
	}

	function schedule_stream_reconnect() {
		if (!should_keep_stream_connected || typeof window === 'undefined') {
			return;
		}
		if (stream_reconnect_timeout_identifier != null) {
			return;
		}

		stream_reconnect_attempt_count += 1;
		const exponential_delay_milliseconds = Math.min(
			stream_reconnect_max_delay_milliseconds,
			stream_reconnect_base_delay_milliseconds * 2 ** (stream_reconnect_attempt_count - 1)
		);
		const jitter_milliseconds =
			exponential_delay_milliseconds * stream_reconnect_jitter_ratio * (Math.random() - 0.5) * 2;
		const reconnect_delay_milliseconds = Math.max(
			1000,
			Math.round(exponential_delay_milliseconds + jitter_milliseconds)
		);
		maybe_set_stream_status_message(
			`Backend stream unavailable. Retrying in ${Math.round(reconnect_delay_milliseconds / 1000)}s.`
		);

		stream_reconnect_timeout_identifier = window.setTimeout(() => {
			stream_reconnect_timeout_identifier = null;
			if (!should_keep_stream_connected) {
				return;
			}
			connect_truck_stream();
		}, reconnect_delay_milliseconds);
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
			stream_reconnect_attempt_count = 0;
			clear_stream_reconnect_timeout();
			set_backend_connection_status('online', { quiet: true });
			maybe_set_stream_status_message('Backend stream connected.');
			void refresh_backend_snapshots({ quiet: true });
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
			set_backend_connection_status('offline', { quiet: true });
			schedule_stream_reconnect();
		});
	}

	async function reconnect_backend_stream_now() {
		if (is_manual_reconnect_in_progress) {
			return;
		}

		is_manual_reconnect_in_progress = true;
		try {
			clear_stream_reconnect_timeout();
			stream_reconnect_attempt_count = 0;
			await check_backend_health({ quiet: true });
			await refresh_backend_snapshots({ quiet: true });
			if (should_keep_stream_connected) {
				connect_truck_stream();
			}
			maybe_set_stream_status_message('Manual reconnect requested.');
		} finally {
			is_manual_reconnect_in_progress = false;
		}
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

	<GeoJSONSource
		id="route_builder_line_source"
		data={{
			type: 'FeatureCollection',
			features:
				route_builder_selected_points.length > 1
					? [
							{
								type: 'Feature',
								geometry: {
									type: 'LineString',
									coordinates: route_builder_selected_points.map((b) => [b.lng, b.lat])
								},
								properties: {}
							}
						]
					: []
		}}
	>
		<LineLayer
			id="route_builder_line_layer"
			paint={{
				'line-color': '#22c55e',
				'line-width': 4,
				'line-dasharray': [2, 2]
			}}
		/>
	</GeoJSONSource>

	{#if is_route_builder_active}
		{#each route_builder_selected_points as pt, i}
			<Marker lnglat={[pt.lng, pt.lat]}>
				<div class="route_builder_marker_badge">{i + 1}</div>
			</Marker>
		{/each}
	{/if}

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

<div class="map_tools_bar">
	<button
		type="button"
		class="map_tool_btn {is_route_builder_active ? 'map_tool_btn_active' : ''}"
		onclick={() => {
			is_route_builder_active = !is_route_builder_active;
			if (!is_route_builder_active) route_builder_selected_points = [];
		}}
		title="Constructor de Ruta Demo"
	>
		<svg
			width="20"
			height="20"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"><path d="M3 12h4l3-9 5 18 3-9h3" /></svg
		>
	</button>
	<button
		type="button"
		class="map_tool_btn {is_simulator_panel_open ? 'map_tool_btn_active' : ''}"
		onclick={() => (is_simulator_panel_open = !is_simulator_panel_open)}
		title="Panel de Simulador"
	>
		<svg
			width="20"
			height="20"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
			><circle cx="12" cy="12" r="10" /><polygon points="10 8 16 12 10 16 10 8" /></svg
		>
	</button>
	<button
		type="button"
		class="map_tool_btn {is_status_open ? 'map_tool_btn_active' : ''}"
		onclick={() => (is_status_open = !is_status_open)}
		title="Estado de API"
	>
		<svg
			width="20"
			height="20"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"><path d="M22 12h-4l-3 9L9 3l-3 9H2" /></svg
		>
	</button>
	<button
		type="button"
		class="map_tool_btn {is_legend_open ? 'map_tool_btn_active' : ''}"
		onclick={() => (is_legend_open = !is_legend_open)}
		title="Leyenda del Mapa"
	>
		<svg
			width="20"
			height="20"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
			><line x1="8" y1="6" x2="21" y2="6" /><line x1="8" y1="12" x2="21" y2="12" /><line
				x1="8"
				y1="18"
				x2="21"
				y2="18"
			/><line x1="3" y1="6" x2="3.01" y2="6" /><line x1="3" y1="12" x2="3.01" y2="12" /><line
				x1="3"
				y1="18"
				x2="3.01"
				y2="18"
			/></svg
		>
	</button>
</div>

{#if is_route_builder_active}
	<div class="route_builder_panel_wrap">
		<h3 style="margin:0 0 12px 0; font-size:1.1rem; color:#fff">Constructor de Ruta</h3>
		<p style="margin:0 0 16px 0; font-size:0.85rem; color:#a3a3a3">
			Haz clic en cualquier punto del mapa para trazar tu ruta. Puntos actuales: {route_builder_selected_points.length}
		</p>

		<div style="display:flex; flex-direction:column; gap:12px; margin-bottom: 16px;">
			<input
				type="text"
				placeholder="Código (ej. D-01)"
				bind:value={route_builder_route_code}
				class="builder_input"
			/>
			<input
				type="text"
				placeholder="Nombre (ej. Demo Ruta Libre)"
				bind:value={route_builder_route_name}
				class="builder_input"
			/>
			<select bind:value={route_builder_zone} class="builder_input">
				<option value="Zona Norte">Zona Norte</option>
				<option value="Zona Sur">Zona Sur</option>
				<option value="Zona Central">Zona Central</option>
				<option value="Zona Este">Zona Este</option>
				<option value="Zona Oeste">Zona Oeste</option>
			</select>
			<select bind:value={route_builder_weekday} class="builder_input">
				<option value={1}>Lunes</option>
				<option value={2}>Martes</option>
				<option value={3}>Miércoles</option>
				<option value={4}>Jueves</option>
				<option value={5}>Viernes</option>
				<option value={6}>Sábado</option>
				<option value={0}>Domingo</option>
			</select>
		</div>

		<button
			class="simulator_toggle_btn save_route_btn"
			onclick={save_custom_route}
			disabled={is_saving_route}
		>
			{is_saving_route ? 'Guardando...' : 'Guardar Ruta'}
		</button>
	</div>
{/if}

{#if is_simulator_panel_open}
	<div id="simulator-panel" class="simulator_panel_wrap">
		<SimulatorControlPanel
			bind:routeId={simulation_route_id}
			bind:zone={simulation_zone}
			bind:truckCount={simulation_truck_count}
			bind:speedMultiplier={simulation_speed_multiplier}
			bind:updateIntervalMs={simulation_update_interval_milliseconds}
			bind:snapToRoute={simulation_snap_to_route}
			routeOptions={simulation_route_options}
			zoneOptions={simulation_zone_options}
			status={simulation_status}
			activeTrucks={simulation_active_trucks}
			messagesSent={simulation_messages_sent}
			on:start={handle_simulation_start}
			on:pause={handle_simulation_pause}
			on:resume={handle_simulation_resume}
			on:stop={handle_simulation_stop}
			on:change={handle_simulation_change}
		/>
	</div>
{/if}

{#if is_status_open}
	<div class="map_status_panel">
		<div>API: {backend_api_base_url}</div>
		<div>
			Backend:
			<span
				class="backend_status_badge"
				class:backend_status_badge_online={backend_connection_status === 'online'}
				class:backend_status_badge_offline={backend_connection_status === 'offline'}
			>
				{backend_connection_status_label}
			</span>
		</div>
		<div>Last health check: {backend_last_health_check_label}</div>
		<div>Last data sync: {backend_last_data_sync_label}</div>
		<div>Rutas activas: {collection_routes_list.length}</div>
		<div>Camiones en vivo: {latest_positions_list.length}</div>
		<div>
			Sim: {simulation_status}
			{#if simulation_backend_control_enabled === true}(backend){:else if simulation_backend_control_enabled === false}
				(local){/if}
		</div>
		<div>{simulation_last_message}</div>
		<button
			type="button"
			class="status_reconnect_btn"
			onclick={reconnect_backend_stream_now}
			disabled={is_manual_reconnect_in_progress}
		>
			{is_manual_reconnect_in_progress ? 'Reconnecting...' : 'Reconnect now'}
		</button>
	</div>
{/if}

{#if is_legend_open}
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
{/if}

<AppModal
	bind:is_open={global_modal_open}
	title={global_modal_title}
	message={global_modal_message}
	type={global_modal_type}
/>

<style>
	.map_tools_bar {
		position: fixed;
		left: 1rem;
		top: 5rem;
		display: flex;
		flex-direction: column;
		gap: 0.6rem;
		z-index: 10;
		background: #ffffff;
		padding: 0.6rem;
		border-radius: 1rem;
		border: 1px solid var(--ecochitas-border);
		box-shadow: 0 4px 20px rgba(0,0,0,0.05);
	}
	:global([data-theme='dark']) .map_tools_bar {
		background: rgba(10, 10, 10, 0.85);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		border-color: #262626;
		box-shadow: 0 4px 30px rgba(0, 0, 0, 0.5);
	}
	@media (max-width: 768px) {
		.map_tools_bar { display: none; }
	}

	.map_tool_btn {
		width: 44px;
		height: 44px;
		border-radius: 0.75rem;
		background: #f3f4f6;
		border: 1px solid transparent;
		color: var(--ecochitas-ink);
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition: all 0.2s ease;
	}
	.map_tool_btn:hover {
		background: #e5e7eb;
	}
	:global([data-theme='dark']) .map_tool_btn {
		background: #171717;
		border-color: #262626;
		color: #ffffff;
	}
	:global([data-theme='dark']) .map_tool_btn:hover {
		background: #262626;
		border-color: #404040;
		color: #22c55e;
	}

	.map_tool_btn_active {
		background: rgba(34, 197, 94, 0.15) !important;
		border-color: #22c55e !important;
		color: #22c55e !important;
		box-shadow: 0 0 15px rgba(34, 197, 94, 0.2);
	}

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

	.route_builder_marker_badge {
		background: #22c55e;
		color: #000;
		font-weight: bold;
		font-size: 0.8rem;
		width: 22px;
		height: 22px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		border: 2px solid #000;
		box-shadow: 0 2px 6px rgba(0, 0, 0, 0.5);
		transform: translate(-50%, -50%);
	}

	.route_builder_panel_wrap {
		position: absolute;
		top: 70px;
		right: 20px;
		z-index: 10;
		width: 320px;
		background: rgba(10, 10, 10, 0.95);
		border: 1px solid #262626;
		border-radius: 12px;
		padding: 20px;
		box-shadow: 0 10px 25px rgba(0, 0, 0, 0.8);
		backdrop-filter: blur(8px);
	}

	.builder_input {
		background: #000;
		border: 1px solid #262626;
		color: #fff;
		border-radius: 6px;
		padding: 8px 12px;
		font-size: 0.9rem;
		width: 100%;
		box-sizing: border-box;
		outline: none;
		transition: border-color 0.2s;
	}
	.builder_input:focus {
		border-color: #22c55e;
	}

	.save_route_btn {
		width: 100%;
		justify-content: center;
		background: #22c55e;
		border-color: #16a34a;
		color: #000;
		transition: all 0.2s;
	}
	.save_route_btn:hover:not(:disabled) {
		background: #16a34a;
	}
	.save_route_btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	/* Force pure dark aesthetic on map panels */
	:global(.maplibregl-popup-content) {
		background: #0a0a0a !important;
		color: #fff !important;
		border: 1px solid #262626;
		border-radius: 8px;
	}
	:global(.maplibregl-popup-anchor-bottom .maplibregl-popup-tip) {
		border-top-color: #0a0a0a !important;
	}
	:global(.maplibregl-popup-anchor-top .maplibregl-popup-tip) {
		border-bottom-color: #0a0a0a !important;
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
		background: rgba(10, 10, 10, 0.95);
		border: 1px solid #262626;
		color: #f8fafc;
		font-size: 0.75rem;
		line-height: 1.25;
		backdrop-filter: blur(8px);
		-webkit-backdrop-filter: blur(8px);
		max-width: min(90vw, 22rem);
		word-break: break-all;
	}

	.backend_status_badge {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		margin-left: 0.35rem;
		padding: 0.08rem 0.42rem;
		border-radius: 999px;
		border: 1px solid #525252;
		background: rgba(115, 115, 115, 0.2);
		color: #e5e5e5;
		font-size: 0.68rem;
		font-weight: 700;
		text-transform: uppercase;
	}

	.backend_status_badge_online {
		border-color: rgba(34, 197, 94, 0.55);
		background: rgba(34, 197, 94, 0.18);
		color: #86efac;
	}

	.backend_status_badge_offline {
		border-color: rgba(239, 68, 68, 0.55);
		background: rgba(239, 68, 68, 0.18);
		color: #fca5a5;
	}

	.status_reconnect_btn {
		margin-top: 0.2rem;
		border: 1px solid #3f3f46;
		background: #18181b;
		color: #fafafa;
		border-radius: 0.55rem;
		padding: 0.34rem 0.52rem;
		font-size: 0.7rem;
		font-weight: 700;
		cursor: pointer;
	}

	.status_reconnect_btn:hover:not(:disabled) {
		border-color: #22c55e;
		color: #bbf7d0;
	}

	.status_reconnect_btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.simulator_panel_wrap {
		position: fixed;
		left: 5rem;
		top: 5.2rem;
		z-index: 45;
		width: min(430px, calc(100vw - 6rem));
		max-height: calc(100vh - 6rem);
		overflow: auto;
		border-radius: 1rem;
	}

	.map_legend {
		position: fixed;
		left: 5rem;
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
			left: 4.5rem;
			max-width: none;
		}

		.map_tools_bar {
			left: 0.7rem;
			top: 5.4rem;
		}

		.simulator_panel_wrap {
			left: 4.5rem;
			right: 0.7rem;
			top: 5.4rem;
			width: auto;
		}

		.map_legend {
			left: 4.5rem;
			bottom: 6rem;
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
