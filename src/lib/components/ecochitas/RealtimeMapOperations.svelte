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
	const backend_api_url = 'http://127.0.0.1:8080';
	const truck_positions_source_identifier = 'truck_positions_source';
	const truck_positions_circle_layer_identifier = 'truck_positions_circle_layer';

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
	let map_container_element = $state<HTMLDivElement | null>(null);

	let maplibre_library: typeof import('maplibre-gl') | null = null;
	let map_instance: import('maplibre-gl').Map | null = null;
	let map_resize_observer: ResizeObserver | null = null;
	let truck_stream_connection: EventSource | null = null;
	let truck_popup_instance: import('maplibre-gl').Popup | null = null;
	let has_registered_layer_interactions = false;
	let latest_positions_by_truck_identifier = new SvelteMap<string, Truck_latest_position>();
	let theme_observer: MutationObserver | null = null;

	onMount(async () => {
		await initialize_map();
		setup_map_resize_observer();
		setup_theme_observer();
		await load_latest_positions_snapshot();
		connect_truck_stream();
	});

	onDestroy(() => {
		disconnect_truck_stream();
		destroy_map();
		theme_observer?.disconnect();
	});

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
			ensure_truck_layers_ready();
			sync_truck_source_data();
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

	function build_snapshot_endpoint_url(): string {
		return `${backend_api_url}/v1/trucks/latest-positions`;
	}

	function build_stream_endpoint_url(): string {
		return `${backend_api_url}/v1/trucks/stream`;
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
	}

	function connect_truck_stream() {
		disconnect_truck_stream();

		const stream_endpoint_url = build_stream_endpoint_url();
		const stream_connection = new EventSource(stream_endpoint_url);
		truck_stream_connection = stream_connection;

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
	}
</script>

<div bind:this={map_container_element} class="map_fullscreen"></div>

<style>
	.map_fullscreen {
		position: fixed;
		top: 0;
		left: 0;
		width: 100vw;
		height: 100vh;
		z-index: 0;
	}

	:global(.maplibregl-ctrl-attrib) {
		font-size: 10px;
		background: oklch(1 0 0 / 0.84);
	}
</style>
