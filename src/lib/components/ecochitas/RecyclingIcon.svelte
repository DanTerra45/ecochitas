<script lang="ts">
	export type RecyclingIconName =
		// Category tabs
		| 'cat_papel'
		| 'cat_plasticos'
		| 'cat_vidrio'
		| 'cat_metales'
		| 'cat_ewaste'
		// How-it-works steps
		| 'step_separar'
		| 'step_limpiar'
		| 'step_entregar'
		| 'step_ganar'
		// Papel & Cartón items
		| 'periodico'
		| 'revista'
		| 'caja_carton'
		| 'hoja_papel'
		| 'folleto'
		| 'tubo_carton'
		// Plásticos items
		| 'botella_pet'
		| 'envase_limpieza'
		| 'vaso_plastico'
		| 'bolsa'
		| 'bidon_plastico'
		| 'envase_bebida'
		// Vidrio items
		| 'botella_vidrio'
		| 'frasco'
		| 'lata_conserva'
		| 'botella_vino'
		| 'frasco_medicina'
		| 'tarro_vidrio'
		// Metales items
		| 'lata_aluminio'
		| 'chatarra'
		| 'pieza_metalica'
		| 'herramienta'
		| 'bidon_metal'
		| 'tuerca'
		// E-Waste items
		| 'celular'
		| 'laptop'
		| 'monitor'
		| 'bateria'
		| 'impresora'
		| 'consola';

	let { name, size = 24 }: { name: RecyclingIconName; size?: number } = $props();
</script>

<svg
	viewBox="0 0 24 24"
	fill="none"
	stroke="currentColor"
	stroke-linecap="round"
	stroke-linejoin="round"
	stroke-width="1.75"
	style="width:{size}px;height:{size}px;flex-shrink:0"
	aria-hidden="true"
>
	<!-- ── CATEGORY ICONS ── -->
	{#if name === 'cat_papel'}
		<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
		<path d="M14 2v6h6" />
		<path d="M8 13h8M8 17h5" />

	{:else if name === 'cat_plasticos'}
		<path d="M10 2h4" />
		<path d="M10 2v3l-2 3.5V20a1.5 1.5 0 0 0 1.5 1.5h5A1.5 1.5 0 0 0 16 20V8.5L14 5V2" />
		<path d="M9.5 13h5" />

	{:else if name === 'cat_vidrio'}
		<path d="M8.5 2h7" />
		<path d="M8.5 2v2.5L6 9v11a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V9l-2.5-4.5V2" />
		<path d="M5.8 11.5h12.4" />

	{:else if name === 'cat_metales'}
		<circle cx="12" cy="12" r="3" />
		<path
			d="M12 2v2.5M12 19.5V22M4.22 4.22l1.77 1.77M18.01 18.01l1.77 1.77M2 12h2.5M19.5 12H22M4.22 19.78l1.77-1.77M18.01 5.99l1.77-1.77"
		/>

	{:else if name === 'cat_ewaste'}
		<rect x="7" y="7" width="10" height="10" rx="1.5" />
		<path d="M7 10H4M7 14H4M17 10h3M17 14h3M10 7V4M14 7V4M10 17v3M14 17v3" />

	<!-- ── STEP ICONS ── -->
	{:else if name === 'step_separar'}
		<rect x="2" y="13" width="6" height="8" rx="1.5" />
		<rect x="9" y="10" width="6" height="11" rx="1.5" />
		<rect x="16" y="13" width="6" height="8" rx="1.5" />
		<path d="M12 2v6M9.5 5l2.5-3 2.5 3" />

	{:else if name === 'step_limpiar'}
		<path d="M12 2C12 2 5 10 5 15a7 7 0 0 0 14 0C19 10 12 2 12 2z" />
		<path d="M9.5 15.5a3 3 0 0 0 5 0" />

	{:else if name === 'step_entregar'}
		<path d="M3.5 8.2h9.7v6.1H3.5z" />
		<path d="M13.2 10h3.7l2 2.5v1.8h-5.7z" />
		<circle cx="7.3" cy="16.3" r="1.5" />
		<circle cx="15.6" cy="16.3" r="1.5" />

	{:else if name === 'step_ganar'}
		<path
			d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"
		/>

	<!-- ── PAPEL & CARTÓN ITEMS ── -->
	{:else if name === 'periodico'}
		<rect x="3" y="5" width="18" height="14" rx="1.5" />
		<path d="M3 10h18M7 10v9" />
		<path d="M10 13h8M10 16h6" />

	{:else if name === 'revista'}
		<path d="M4 3h14a1 1 0 0 1 1 1v16a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1z" />
		<path d="M8 3v18" />
		<path d="M11 8h5M11 12h5M11 16h3" />

	{:else if name === 'caja_carton'}
		<path
			d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"
		/>
		<path d="M3.27 6.96L12 12.01l8.73-5.05M12 22.08V12" />

	{:else if name === 'hoja_papel'}
		<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6z" />
		<path d="M14 2v6h6" />
		<path d="M8 13h8M8 17h5" />

	{:else if name === 'folleto'}
		<rect x="2" y="4" width="20" height="16" rx="1.5" />
		<path d="M9 4v16M15 4v16" />
		<path d="M4 10h3M11 10h2M17 10h3M4 14h3M11 14h2M17 14h3" />

	{:else if name === 'tubo_carton'}
		<ellipse cx="12" cy="7" rx="6" ry="2.5" />
		<path d="M6 7v10M18 7v10" />
		<ellipse cx="12" cy="17" rx="6" ry="2.5" />

	<!-- ── PLÁSTICOS ITEMS ── -->
	{:else if name === 'botella_pet'}
		<path d="M10 2h4v3l-1.5 2.5V20a1.5 1.5 0 0 1-1.5 1.5h-1A1.5 1.5 0 0 1 8.5 20V7.5L10 5V2z" />
		<path d="M8.5 12h7" />
		<path d="M10 2h4" />

	{:else if name === 'envase_limpieza'}
		<rect x="7" y="9" width="10" height="13" rx="1.5" />
		<path d="M11 9V6" />
		<path d="M8.5 6h3a2.5 2.5 0 0 1 2.5 2.5" />
		<path d="M14 5h3v3.5" />
		<path d="M9 13.5h6" />

	{:else if name === 'vaso_plastico'}
		<path d="M6 4h12l-2 17H8L6 4z" />
		<path d="M6.8 9h10.4" />

	{:else if name === 'bolsa'}
		<path d="M6 7h12l-1.5 14H7.5L6 7z" />
		<path d="M9 7C9 5.34 10.34 4 12 4s3 1.34 3 3" />

	{:else if name === 'bidon_plastico'}
		<rect x="6" y="7" width="12" height="14" rx="2" />
		<path d="M9 7V5a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2" />
		<path d="M18 11h2a1 1 0 0 1 0 4h-2" />
		<path d="M7.5 12h9" />

	{:else if name === 'envase_bebida'}
		<rect x="5" y="6" width="14" height="15" rx="1.5" />
		<path d="M9 3h6v3H9z" />
		<path d="M12 6V3" />
		<path d="M8 12h8M8 15.5h5" />

	<!-- ── VIDRIO ITEMS ── -->
	{:else if name === 'botella_vidrio'}
		<path d="M10.5 2h3v4.5l2 4.5V20a1.5 1.5 0 0 1-1.5 1.5h-4A1.5 1.5 0 0 1 8.5 20V11l2-4.5V2z" />
		<path d="M10.5 2h3" />
		<path d="M9 14.5h6" />

	{:else if name === 'frasco'}
		<path d="M7 9h10v10a2 2 0 0 1-2 2H9a2 2 0 0 1-2-2V9z" />
		<path d="M5.5 6h13a1 1 0 0 1 1 1v2H5.5V7a1 1 0 0 1 1-1z" />
		<path d="M9 6V4.5h6V6" />
		<path d="M9 13.5h6" />

	{:else if name === 'lata_conserva'}
		<ellipse cx="12" cy="6" rx="7" ry="2.5" />
		<path d="M5 6v12M19 6v12" />
		<ellipse cx="12" cy="18" rx="7" ry="2.5" />
		<path d="M5 10h14" />

	{:else if name === 'botella_vino'}
		<path d="M10 2h4v4l2.5 5v9a1.5 1.5 0 0 1-1.5 1.5h-6A1.5 1.5 0 0 1 7.5 20v-9L10 6V2z" />
		<path d="M10 2h4" />
		<path d="M8.5 16h7" />

	{:else if name === 'frasco_medicina'}
		<rect x="8" y="8" width="8" height="13" rx="1.5" />
		<rect x="9.5" y="5" width="5" height="3" rx="1" />
		<path d="M10 13h4M12 11v4" />

	{:else if name === 'tarro_vidrio'}
		<rect x="5" y="9" width="14" height="12" rx="2" />
		<path d="M4 6.5h16a1 1 0 0 1 1 1V9H3V7.5a1 1 0 0 1 1-1z" />
		<path d="M8 9V7M16 9V7" />
		<path d="M8 14.5h8M8 17.5h5" />

	<!-- ── METALES ITEMS ── -->
	{:else if name === 'lata_aluminio'}
		<ellipse cx="12" cy="6" rx="6" ry="2" />
		<path d="M6 6v12M18 6v12" />
		<ellipse cx="12" cy="18" rx="6" ry="2" />
		<path d="M11 4.5h2v2.5" />
		<path d="M13 4.5l2-1.5" />

	{:else if name === 'chatarra'}
		<path
			d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"
		/>

	{:else if name === 'pieza_metalica'}
		<circle cx="12" cy="12" r="3" />
		<path
			d="M12 2v2.5M12 19.5V22M4.22 4.22l1.77 1.77M18.01 18.01l1.77 1.77M2 12h2.5M19.5 12H22M4.22 19.78l1.77-1.77M18.01 5.99l1.77-1.77"
		/>

	{:else if name === 'herramienta'}
		<path
			d="M3 21l9.5-9.5M21 7a4 4 0 0 1-5.66 5.66l-9.9 9.9a1.41 1.41 0 0 1-2-2l9.9-9.9A4 4 0 0 1 21 7z"
		/>

	{:else if name === 'bidon_metal'}
		<ellipse cx="12" cy="5" rx="8" ry="3" />
		<path d="M4 5v14M20 5v14" />
		<ellipse cx="12" cy="19" rx="8" ry="3" />
		<path d="M4 10h16M4 14.5h16" />

	{:else if name === 'tuerca'}
		<path d="M12 2l9 5.2v9.6L12 22l-9-5.2V7.2z" />
		<circle cx="12" cy="12" r="3.5" />

	<!-- ── E-WASTE ITEMS ── -->
	{:else if name === 'celular'}
		<rect x="7" y="2" width="10" height="20" rx="2.5" />
		<path d="M10 5h4" />
		<circle cx="12" cy="18" r="1" fill="currentColor" stroke="none" />

	{:else if name === 'laptop'}
		<rect x="3" y="4" width="18" height="12" rx="2" />
		<path d="M1 20h22" />
		<path d="M7 20h10" />
		<path d="M8 10h8M8 13h5" />

	{:else if name === 'monitor'}
		<rect x="2" y="3" width="20" height="14" rx="2" />
		<path d="M8 21h8M12 17v4" />
		<path d="M7 9h10M7 13h6" />

	{:else if name === 'bateria'}
		<rect x="2" y="7" width="18" height="10" rx="2" />
		<path d="M20 11h2v2h-2" />
		<path d="M10 10l-2 2h4l-2 2" />

	{:else if name === 'impresora'}
		<path d="M6 9V2.5h12V9" />
		<rect x="2" y="9" width="20" height="9" rx="2" />
		<path d="M6 18v3.5h12V18" />
		<path d="M6 13.5h.01M10 13.5h4" />

	{:else if name === 'consola'}
		<path d="M6 10a6 6 0 0 1 12 0v5a3 3 0 0 1-3 3H9a3 3 0 0 1-3-3v-5z" />
		<path d="M9 12h2M10 11v2" />
		<circle cx="15" cy="12" r="0.8" fill="currentColor" stroke="none" />
		<circle cx="15" cy="14.2" r="0.8" fill="currentColor" stroke="none" />
	{/if}
</svg>
