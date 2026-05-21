<script lang="ts">
	import EcochitasIcon from '$lib/components/ecochitas/EcochitasIcon.svelte';
	import RecyclingIcon, {
		type RecyclingIconName
	} from '$lib/components/ecochitas/RecyclingIcon.svelte';

	// ─── Recycling category types ─────────────────────────────────────────────
	type Recyclable_item = { icon: RecyclingIconName; name: string; desc: string };
	type Recycling_category = {
		id: string; name: string; icon: RecyclingIconName;
		accent: string; accent_light: string; accent_border: string;
		price: string; price_numeric: number; unit: string; tip: string;
		items: Recyclable_item[];
	};
	type MaterialId = 'plastico' | 'vidrio' | 'papel' | 'metal' | 'electronico' | 'organico' | 'pilas' | 'carton' | 'aceite';
	type AcopioStatus = 'activo' | 'lleno' | 'cerrado';
	type AcopioPoint = {
		id: string; name: string; address: string; zone: string;
		schedule: string; materials: MaterialId[]; status: AcopioStatus;
		capacity_pct: number; phone?: string;
	};

	// ─── Ruta del reciclaje (6 pasos) ────────────────────────────────────────
	const ruta_steps = [
		{ num: 1, color: '#4ade80', bg: 'rgba(74, 222, 128, 0.1)', border: 'rgba(74, 222, 128, 0.2)', title: 'Entrega', desc: 'Separás y llevás residuos limpios al punto de acopio más cercano', icon: 'M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2M9 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75' },
		{ num: 2, color: '#38bdf8', bg: 'rgba(56, 189, 248, 0.1)', border: 'rgba(56, 189, 248, 0.2)', title: 'Clasificación', desc: 'El operador verifica y clasifica cada material por tipo', icon: 'M22 3H2l8 9.46V19l4 2v-8.54L22 3z' },
		{ num: 3, color: '#a78bfa', bg: 'rgba(167, 139, 250, 0.1)', border: 'rgba(167, 139, 250, 0.2)', title: 'Recolección', desc: 'El camión de la GAMC transporta los materiales al centro', icon: 'M3.5 8.2h9.7v6.1H3.5zM13.2 10h3.7l2 2.5v1.8h-5.7zM7.3 14.8a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3zM15.6 14.8a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3z' },
		{ num: 4, color: '#f87171', bg: 'rgba(248, 113, 113, 0.1)', border: 'rgba(248, 113, 113, 0.2)', title: 'Reciclaje', desc: 'Las plantas especializadas procesan y transforman los materiales', icon: 'M3 12a9 9 0 0 1 15.2-6.4M18.2 5.6V2.4h-3.2M21 12a9 9 0 0 1-15.2 6.4M5.8 18.4v3.2H9' },
		{ num: 5, color: '#fbbf24', bg: 'rgba(251, 191, 36, 0.1)', border: 'rgba(251, 191, 36, 0.2)', title: 'Puntos', desc: 'Recibís puntos ecológicos acreditados en tu perfil automáticamente', icon: 'M12 3l2.7 5.5 6.1.9-4.4 4.3 1 6.1L12 17l-5.4 2.8 1-6.1L3.2 9.4l6.1-.9z' },
		{ num: 6, color: '#22c55e', bg: 'rgba(34, 197, 94, 0.1)', border: 'rgba(34, 197, 94, 0.2)', title: 'Recompensa', desc: 'Canjeás puntos por descuentos, productos y beneficios municipales', icon: 'M20 12v10H4V12M2 7h20v5H2zM12 22V7M12 7H7.5a2.5 2.5 0 0 1 0-5C11 2 12 7 12 7zM12 7h4.5a2.5 2.5 0 0 0 0-5C13 2 12 7 12 7z' }
	];

	// ─── Material filters ─────────────────────────────────────────────────────
	const material_filters = [
		{ id: 'todos', label: 'Todos', icon: 'M4 6h16M4 12h16M4 18h16' },
		{ id: 'plastico', label: 'Plástico', icon: 'M10 2h4M10 2v3l-2 3.5V20a1.5 1.5 0 0 0 1.5 1.5h5A1.5 1.5 0 0 0 16 20V8.5L14 5V2M9.5 13h5' },
		{ id: 'vidrio', label: 'Vidrio', icon: 'M8.5 2h7M8.5 2v2.5L6 9v11a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V9l-2.5-4.5V2M5.8 11.5h12.4' },
		{ id: 'papel', label: 'Papel', icon: 'M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8zM14 2v6h6M8 13h8M8 17h5' },
		{ id: 'electronico', label: 'Electrónicos', icon: 'M9 3H5a2 2 0 0 0-2 2v4m6-6h10a2 2 0 1 2 2v4M9 3v18m0 0h10a2 2 0 0 0 2-2V9M9 21H5a2 2 0 1 1-2-2V9m0 0h18' },
		{ id: 'organico', label: 'Orgánicos', icon: 'M12 22a7 7 0 0 0 7-7c0-2-1-3.9-3-5.5s-3.5-4-4-6.5c-.5 2.5-2 4.9-4 6.5C6 11.1 5 13 5 15a7 7 0 0 0 7 7z' },
		{ id: 'pilas', label: 'Pilas', icon: 'M7 7h10a2 2 0 0 1 2 2v6a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V9a2 2 0 0 1 2-2zM9 4h6M10 11v2M14 11v2' },
		{ id: 'carton', label: 'Cartón', icon: 'M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16zM3.27 6.96L12 12.01l8.73-5.05M12 22.08V12' },
		{ id: 'aceite', label: 'Aceite', icon: 'M12 2.69l5.66 5.66a8 8 0 1 1-11.31 0z' }
	];

	const material_meta: Record<string, { label: string; color: string; light: string }> = {
		plastico:   { label: 'Plástico',    color: '#38bdf8', light: 'rgba(56, 189, 248, 0.15)' },
		vidrio:     { label: 'Vidrio',      color: '#34d399', light: 'rgba(52, 211, 153, 0.15)' },
		papel:      { label: 'Papel',       color: '#fbbf24', light: 'rgba(251, 191, 36, 0.15)' },
		metal:      { label: 'Metal',       color: '#a78bfa', light: 'rgba(167, 139, 250, 0.15)' },
		electronico:{ label: 'E-Waste',     color: '#f87171', light: 'rgba(248, 113, 113, 0.15)' },
		organico:   { label: 'Orgánico',    color: '#84cc16', light: 'rgba(132, 204, 22, 0.15)' },
		pilas:      { label: 'Pilas',       color: '#fb923c', light: 'rgba(251, 146, 60, 0.15)' },
		carton:     { label: 'Cartón',      color: '#d97706', light: 'rgba(217, 119, 6, 0.15)' },
		aceite:     { label: 'Aceite',      color: '#65a30d', light: 'rgba(101, 163, 13, 0.15)' }
	};

	// ─── Puntos de acopio mock data ───────────────────────────────────────────
	const acopio_points: AcopioPoint[] = [
		{ id: 'p1', name: 'Centro de Acopio Norte', address: 'Av. América 2345', zone: 'Zona Norte', schedule: 'Lun–Sáb 8:00–18:00', materials: ['plastico', 'papel', 'carton', 'vidrio'], status: 'activo', capacity_pct: 45, phone: '4-452-1890' },
		{ id: 'p2', name: 'Punto Verde Sarco', address: 'Calle Heroínas 456', zone: 'Sarco', schedule: 'Lun–Vie 9:00–17:00', materials: ['metal', 'electronico', 'pilas', 'aceite'], status: 'activo', capacity_pct: 60, phone: '4-441-2233' },
		{ id: 'p3', name: 'Eco-Punto Central', address: 'Plaza Principal 100', zone: 'Centro', schedule: 'Mar–Dom 7:00–19:00', materials: ['plastico', 'vidrio', 'organico', 'carton'], status: 'activo', capacity_pct: 30, phone: '4-450-7799' },
		{ id: 'p4', name: 'Centro Queru Queru', address: 'Av. Pando 789', zone: 'Queru Queru', schedule: 'Lun–Sáb 8:30–17:30', materials: ['papel', 'plastico', 'carton'], status: 'lleno', capacity_pct: 93, phone: '4-443-6612' },
		{ id: 'p5', name: 'Punto Tiquipaya Sur', address: 'Calle Tahuapalca 33', zone: 'Tiquipaya Sur', schedule: 'Mié–Dom 10:00–18:00', materials: ['organico', 'electronico', 'aceite'], status: 'activo', capacity_pct: 55, phone: '4-478-3344' },
		{ id: 'p6', name: 'Acopio Villa Pagador', address: 'Av. Blanco Galindo Km 5', zone: 'Villa Pagador', schedule: 'Lun–Vie 8:00–16:00', materials: ['metal', 'pilas', 'papel', 'carton'], status: 'activo', capacity_pct: 40, phone: '4-462-8800' },
		{ id: 'p7', name: 'Eco-Centro Cala Cala', address: 'Calle Lanza 200', zone: 'Cala Cala', schedule: 'Mar–Sáb 9:00–18:00', materials: ['plastico', 'vidrio', 'papel', 'pilas', 'carton'], status: 'activo', capacity_pct: 70, phone: '4-455-1177' },
		{ id: 'p8', name: 'Punto Verde Mayorazgo', address: 'Av. Costanera 512', zone: 'Mayorazgo', schedule: 'Lun–Dom 7:00–20:00', materials: ['organico', 'plastico', 'aceite'], status: 'cerrado', capacity_pct: 85, phone: '4-489-9900' }
	];

	let active_filter = $state('todos');
	let show_all_acopio = $state(false);

	const filtered_points = $derived(
		active_filter === 'todos'
			? acopio_points
			: acopio_points.filter((p) => p.materials.includes(active_filter as MaterialId))
	);
	const visible_points = $derived(
		show_all_acopio ? filtered_points : filtered_points.slice(0, 6)
	);
	const has_more = $derived(!show_all_acopio && filtered_points.length > 6);

	// ─── Recycling categories (unchanged) ────────────────────────────────────
	const recycling_categories: Recycling_category[] = [
		{
			id: 'papel', name: 'Papel & Cartón', icon: 'cat_papel',
			accent: '#fbbf24', accent_light: 'rgba(251, 191, 36, 0.08)', accent_border: 'rgba(251, 191, 36, 0.25)',
			price: '3.50 Bs/kg', price_numeric: 3.5, unit: 'por kilogramo',
			tip: 'Evitar papel húmedo o mezclado con comida',
			items: [
				{ icon: 'periodico', name: 'Periódicos', desc: 'Diarios y semanales' },
				{ icon: 'revista', name: 'Revistas', desc: 'Cualquier publicación' },
				{ icon: 'caja_carton', name: 'Cajas', desc: 'Cajas de cartón limpio' },
				{ icon: 'hoja_papel', name: 'Hojas', desc: 'Papel de oficina' },
				{ icon: 'folleto', name: 'Folletos', desc: 'Trípticos y volantes' },
				{ icon: 'tubo_carton', name: 'Tubos', desc: 'Rollos de cartón' }
			]
		},
		{
			id: 'plasticos', name: 'Plásticos', icon: 'cat_plasticos',
			accent: '#38bdf8', accent_light: 'rgba(56, 189, 248, 0.08)', accent_border: 'rgba(56, 189, 248, 0.25)',
			price: '4.00 Bs/kg', price_numeric: 4.0, unit: 'por kilogramo',
			tip: 'Enjuagado y compactado, sin tapa',
			items: [
				{ icon: 'botella_pet', name: 'Botellas PET', desc: 'Transparentes de bebida' },
				{ icon: 'envase_limpieza', name: 'Limpieza', desc: 'Dispensadores y frascos' },
				{ icon: 'vaso_plastico', name: 'Vasos', desc: 'Desechables y reutilizables' },
				{ icon: 'bolsa', name: 'Bolsas', desc: 'Bolsas reutilizables' },
				{ icon: 'bidon_plastico', name: 'Bidones', desc: 'Contenedores grandes' },
				{ icon: 'envase_bebida', name: 'Envases', desc: 'Cajas de jugo y leche' }
			]
		},
		{
			id: 'vidrio', name: 'Vidrio', icon: 'cat_vidrio',
			accent: '#34d399', accent_light: 'rgba(52, 211, 153, 0.08)', accent_border: 'rgba(52, 211, 153, 0.25)',
			price: '1.60 Bs/kg', price_numeric: 1.6, unit: 'por kilogramo',
			tip: 'Separar por color reduce los rechazos',
			items: [
				{ icon: 'botella_vidrio', name: 'Botellas', desc: 'De agua y gaseosas' },
				{ icon: 'frasco', name: 'Frascos', desc: 'Con boca ancha' },
				{ icon: 'lata_conserva', name: 'Conservas', desc: 'Envases de alimentos' },
				{ icon: 'botella_vino', name: 'Vino / cerveza', desc: 'Botellas de bebidas' },
				{ icon: 'frasco_medicina', name: 'Medicinas', desc: 'Frascos de farmacia' },
				{ icon: 'tarro_vidrio', name: 'Tarros', desc: 'Tarros y potes' }
			]
		},
		{
			id: 'metales', name: 'Metales', icon: 'cat_metales',
			accent: '#a78bfa', accent_light: 'rgba(167, 139, 250, 0.08)', accent_border: 'rgba(167, 139, 250, 0.25)',
			price: '6.20 Bs/kg', price_numeric: 6.2, unit: 'por kilogramo',
			tip: 'Sin residuos orgánicos adheridos',
			items: [
				{ icon: 'lata_aluminio', name: 'Latas', desc: 'Aluminio de bebidas' },
				{ icon: 'chatarra', name: 'Chatarra', desc: 'Piezas metálicas varias' },
				{ icon: 'pieza_metalica', name: 'Repuestos', desc: 'Engranajes y piezas' },
				{ icon: 'herramienta', name: 'Herramientas', desc: 'En desuso o rotas' },
				{ icon: 'bidon_metal', name: 'Bidones', desc: 'Tambores metálicos' },
				{ icon: 'tuerca', name: 'Tornillería', desc: 'Tuercas, pernos, clavos' }
			]
		},
		{
			id: 'ewaste', name: 'E-Waste', icon: 'cat_ewaste',
			accent: '#f87171', accent_light: 'rgba(248, 113, 113, 0.08)', accent_border: 'rgba(248, 113, 113, 0.25)',
			price: 'Punto especial', price_numeric: 0, unit: 'acopio GAMC',
			tip: 'Solo en puntos de acopio autorizados por la GAMC',
			items: [
				{ icon: 'celular', name: 'Celulares', desc: 'Teléfonos en desuso' },
				{ icon: 'laptop', name: 'Laptops', desc: 'Computadoras portátiles' },
				{ icon: 'monitor', name: 'Monitores', desc: 'Pantallas y TVs' },
				{ icon: 'bateria', name: 'Baterías', desc: 'Pilas y acumuladores' },
				{ icon: 'impresora', name: 'Impresoras', desc: 'Equipos de oficina' },
				{ icon: 'consola', name: 'Consolas', desc: 'Videojuegos y controles' }
			]
		}
	];

	const max_price = Math.max(...recycling_categories.map((c) => c.price_numeric));
	let selected_id = $state(recycling_categories[0].id);
	const selected = $derived(recycling_categories.find((c) => c.id === selected_id) ?? recycling_categories[0]);
	const sorted_by_price = $derived([...recycling_categories].sort((a, b) => b.price_numeric - a.price_numeric));
</script>

<!-- ── Hero ── -->
<div class="hero">
	<div class="hero_eyebrow">
		<EcochitasIcon name="recycling" size={14} />
		<span>EcoHub · Cochabamba · GAMC</span>
	</div>
	<h1 class="hero_title">Recicla mejor,<br />gana más puntos</h1>
	<p class="hero_subtitle">
		Encuentra puntos de acopio, aprende la ruta del reciclaje y descubrí cuánto valen tus
		materiales en el mercado.
	</p>
	<div class="hero_stats">
		<div class="hero_stat"><span class="stat_num">8</span><span class="stat_lbl">Puntos activos</span></div>
		<div class="stat_div"></div>
		<div class="hero_stat"><span class="stat_num">6.20 Bs</span><span class="stat_lbl">Precio máx/kg</span></div>
		<div class="stat_div"></div>
		<div class="hero_stat"><span class="stat_num">Viernes</span><span class="stat_lbl">Turno recolección</span></div>
	</div>
</div>

<!-- ── Ruta del Reciclaje ── -->
<section class="panel" id="como-funciona">
	<div class="section_head_row">
		<EcochitasIcon name="recycling" size={18} />
		<h2 class="section_title">Ruta del Reciclaje</h2>
	</div>
	<p class="section_sub">Así viaja tu residuo desde tu mano hasta convertirse en un nuevo producto.</p>

	<div class="ruta_grid">
		{#each ruta_steps as step, i (step.num)}
			<div class="ruta_step">
				<div class="ruta_icon_wrap" style="background:{step.bg};border-color:{step.border};color:{step.color}">
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:22px;height:22px" aria-hidden="true">
						<path d={step.icon} />
					</svg>
					<span class="ruta_num" style="background:{step.color}">{step.num}</span>
				</div>
				<strong class="ruta_title" style="color:{step.color}">{step.title}</strong>
				<p class="ruta_desc">{step.desc}</p>
				{#if i < ruta_steps.length - 1}
					<div class="ruta_arrow" aria-hidden="true">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" style="width:14px;height:14px;color:var(--ecochitas-muted)" aria-hidden="true">
							<polyline points="9 18 15 12 9 6"/>
						</svg>
					</div>
				{/if}
			</div>
		{/each}
	</div>
</section>

<!-- ── Puntos de Acopio ── -->
<section class="panel">
	<div class="section_head_row">
		<svg viewBox="0 0 24 24" fill="none" stroke="var(--ecochitas-leaf)" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:18px;height:18px;flex-shrink:0" aria-hidden="true">
			<path d="M12 21s-8-7.5-8-12a8 8 0 0 1 16 0c0 4.5-8 12-8 12z"/>
			<circle cx="12" cy="9" r="2.5"/>
		</svg>
		<h2 class="section_title">Puntos de Acopio</h2>
		<span class="acopio_count_badge">{filtered_points.length} activos</span>
	</div>
	<p class="section_sub">Encontrá el punto más cercano según el material que querés reciclar.</p>

	<!-- Pill filters -->
	<div class="filter_scroll">
		<div class="filter_pills">
			{#each material_filters as f (f.id)}
				<button
					class="filter_pill"
					class:filter_pill_active={active_filter === f.id}
					onclick={() => { active_filter = f.id; show_all_acopio = false; }}
				>
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:13px;height:13px;flex-shrink:0" aria-hidden="true">
						<path d={f.icon}/>
					</svg>
					{f.label}
				</button>
			{/each}
		</div>
	</div>

	<!-- Cards grid -->
	<div class="acopio_grid">
		{#each visible_points as pt (pt.id)}
			<div class="acopio_card">
				<!-- Card header -->
				<div class="acopio_card_head">
					<div class="acopio_card_icon">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:18px;height:18px" aria-hidden="true">
							<path d="M12 21s-8-7.5-8-12a8 8 0 0 1 16 0c0 4.5-8 12-8 12z"/>
							<circle cx="12" cy="9" r="2.5"/>
						</svg>
					</div>
					<div class="acopio_card_head_text">
						<strong class="acopio_name">{pt.name}</strong>
						<span class="acopio_zone">{pt.zone}</span>
					</div>
					<span class="acopio_status acopio_status_{pt.status}">{pt.status === 'activo' ? 'Disponible' : pt.status === 'lleno' ? 'Saturado' : 'Cerrado temp.'}</span>
				</div>

				<!-- Card body -->
				<div class="acopio_card_body">
					<div class="acopio_info_row">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:13px;height:13px;flex-shrink:0;color:var(--ecochitas-muted)" aria-hidden="true">
							<path d="M21 10c0 7-9 13-9 13S3 17 3 10a9 9 0 0 1 18 0z"/>
							<circle cx="12" cy="10" r="3"/>
						</svg>
						<span>{pt.address}</span>
					</div>
					<div class="acopio_info_row">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:13px;height:13px;flex-shrink:0;color:var(--ecochitas-muted)" aria-hidden="true">
							<circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/>
						</svg>
						<span>{pt.schedule}</span>
					</div>

					<!-- Capacity -->
					<div class="acopio_capacity_row">
						<span class="capacity_label">Capacidad</span>
						<div class="capacity_bar_wrap">
							<div class="capacity_bar" style="width:{pt.capacity_pct}%;background:{pt.capacity_pct >= 85 ? '#ef4444' : pt.capacity_pct >= 60 ? '#f59e0b' : '#16a34a'}"></div>
						</div>
						<span class="capacity_pct_text" style="color:{pt.capacity_pct >= 85 ? '#ef4444' : pt.capacity_pct >= 60 ? '#d97706' : '#16a34a'}">{pt.capacity_pct}%</span>
					</div>

					<!-- Materials -->
					<div class="acopio_materials">
						{#each pt.materials as mat (mat)}
							<span class="mat_tag" style="background:{material_meta[mat]?.light};color:{material_meta[mat]?.color}">
								{material_meta[mat]?.label}
							</span>
						{/each}
					</div>
				</div>

				<!-- Card footer -->
				<div class="acopio_card_foot">
					<button class="acopio_loc_btn">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:13px;height:13px;flex-shrink:0" aria-hidden="true">
							<path d="M12 21s-8-7.5-8-12a8 8 0 0 1 16 0c0 4.5-8 12-8 12z"/>
							<circle cx="12" cy="9" r="2.5"/>
						</svg>
						Ver mapa
					</button>
					<button class="acopio_detail_btn">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:13px;height:13px;flex-shrink:0" aria-hidden="true">
							<circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/>
						</svg>
						Ver detalle
					</button>
				</div>
			</div>
		{/each}
	</div>

	{#if has_more}
		<button class="ver_todos_btn" onclick={() => (show_all_acopio = true)}>
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" style="width:15px;height:15px" aria-hidden="true">
				<polyline points="6 9 12 15 18 9"/>
			</svg>
			Ver todos ({filtered_points.length})
		</button>
	{:else if show_all_acopio && filtered_points.length > 6}
		<button class="ver_todos_btn ver_menos_btn" onclick={() => (show_all_acopio = false)}>
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" style="width:15px;height:15px" aria-hidden="true">
				<polyline points="18 15 12 9 6 15"/>
			</svg>
			Ver menos
		</button>
	{/if}
</section>

<!-- ── ¿Qué puedo reciclar? ── -->
<section class="panel">
	<h2 class="section_title">¿Qué puedo reciclar?</h2>
	<p class="section_sub">Seleccioná una categoría para ver sus artículos aceptados.</p>

	<div class="cat_picker" role="tablist">
		{#each recycling_categories as cat (cat.id)}
			<button role="tab" aria-selected={selected_id === cat.id} class="cat_btn" class:cat_btn_active={selected_id === cat.id} style="--accent:{cat.accent};--accent-light:{cat.accent_light};--accent-border:{cat.accent_border}" onclick={() => (selected_id = cat.id)}>
				<span class="cat_btn_icon"><RecyclingIcon name={cat.icon} size={22} /></span>
				<span class="cat_btn_name">{cat.name}</span>
				<span class="cat_btn_price">{cat.price}</span>
			</button>
		{/each}
	</div>

	<div class="cat_detail" style="--accent:{selected.accent};--accent-light:{selected.accent_light};--accent-border:{selected.accent_border}">
		<div class="cat_detail_head">
			<div class="cat_detail_icon"><RecyclingIcon name={selected.icon} size={28} /></div>
			<div class="cat_detail_meta">
				<h3 class="cat_detail_name">{selected.name}</h3>
				<p class="cat_detail_tip"><EcochitasIcon name="alert" size={12} />{selected.tip}</p>
			</div>
			<div class="cat_price_box">
				<span class="price_box_lbl">Valor de mercado</span>
				<strong class="price_box_val">{selected.price}</strong>
				<span class="price_box_unit">{selected.unit}</span>
			</div>
		</div>
		<div class="items_grid">
			{#each selected.items as item (item.name)}
				<div class="item_card">
					<div class="item_icon"><RecyclingIcon name={item.icon} size={28} /></div>
					<div class="item_text">
						<strong class="item_name">{item.name}</strong>
						<span class="item_desc">{item.desc}</span>
					</div>
				</div>
			{/each}
		</div>
	</div>
</section>

<!-- ── Precios ── -->
<section class="panel">
	<div class="section_head_row">
		<EcochitasIcon name="rewards" size={18} />
		<h2 class="section_title">Precios de Mercado</h2>
	</div>
	<p class="section_sub">Valores de referencia para el turno de reciclaje en Cochabamba.</p>
	<div class="ranking_list">
		{#each sorted_by_price as cat, i (cat.id)}
			<div class="rank_row" style="--accent:{cat.accent};--accent-light:{cat.accent_light};--accent-border:{cat.accent_border}">
				<span class="rank_pos">#{i + 1}</span>
				<span class="rank_icon"><RecyclingIcon name={cat.icon} size={20} /></span>
				<span class="rank_name">{cat.name}</span>
				<div class="rank_bar_track">
					{#if cat.price_numeric > 0}
						<div class="rank_bar_fill" style="width:{(cat.price_numeric / max_price) * 100}%"></div>
					{:else}
						<span class="rank_special_label">Punto de acopio GAMC</span>
					{/if}
				</div>
				<strong class="rank_price">{cat.price}</strong>
			</div>
		{/each}
	</div>
</section>

<!-- ── Reglas ── -->
<section class="panel">
	<div class="section_head_row">
		<EcochitasIcon name="alert" size={18} />
		<h2 class="section_title">Reglas del Programa</h2>
	</div>
	<div class="rules_cols">
		<div class="rules_group rules_do">
			<div class="rules_group_head">
				<span class="rules_head_icon do_icon"><EcochitasIcon name="upload" size={14} /></span>
				<span class="rules_head_label">Qué debes hacer</span>
			</div>
			<ul class="rules_list">
				<li>Toma una foto del aporte en bolsa limpia antes de entregar.</li>
				<li>Selecciona la categoría del material al registrar el aporte.</li>
				<li>Activa la geolocalización del hogar o contenedor asignado.</li>
				<li>Espera la confirmación del recolector al pesar el material.</li>
			</ul>
		</div>
		<div class="rules_group rules_dont">
			<div class="rules_group_head">
				<span class="rules_head_icon dont_icon"><EcochitasIcon name="alert" size={14} /></span>
				<span class="rules_head_label">Qué debes evitar</span>
			</div>
			<ul class="rules_list">
				<li>Mezclar residuos orgánicos con materiales reciclables.</li>
				<li>Entregar envases mojados o con restos de comida adheridos.</li>
				<li>Superar el 20 % de contaminación — congela tu bonificación.</li>
				<li>Llevar e-waste al recolector común; usa solo los puntos GAMC.</li>
			</ul>
		</div>
	</div>
</section>

<style>
	h1, h2, h3 { margin: 0; font-family: 'Sora', 'Plus Jakarta Sans', sans-serif; letter-spacing: -0.02em; }
	p { margin: 0; }

	/* ── Panel (shared section card) ── */
	.panel {
		background: var(--ecochitas-surface);
		border: 1px solid var(--ecochitas-border);
		border-radius: 1.25rem;
		padding: 1.25rem 1.15rem;
		overflow: hidden;
		margin-bottom: 0.75rem;
	}

	/* ── Hero ── */
	.hero {
		background: var(--ecochitas-surface);
		border: 1px solid var(--ecochitas-border);
		box-shadow: 0 0 40px rgba(34, 197, 94, 0.05);
		border-radius: 1.5rem; padding: 1.75rem 1.5rem 1.4rem;
		display: flex; flex-direction: column; gap: 0.9rem;
	}
	:root[data-theme='dark'] .hero {
		background: rgba(255, 255, 255, 0.03);
		border-color: rgba(255, 255, 255, 0.1);
	}
	.hero_eyebrow {
		display: inline-flex; align-items: center; gap: 0.4rem;
		font-size: 0.72rem; font-weight: 700; letter-spacing: 0.06em; text-transform: uppercase;
		color: var(--ecochitas-leaf); background: rgba(34, 197, 94, 0.1);
		border: 1px solid rgba(34, 197, 94, 0.2); border-radius: 999px;
		padding: 0.25rem 0.7rem; width: fit-content;
	}
	:root[data-theme='dark'] .hero_eyebrow { color: #4ade80; }
	.hero_title { font-size: 1.7rem; font-weight: 900; line-height: 1.12; color: var(--ecochitas-ink); }
	.hero_subtitle { font-size: 0.88rem; line-height: 1.6; color: var(--ecochitas-muted); max-width: 36ch; }
	.hero_stats {
		display: flex; align-items: center; gap: 0.75rem; flex-wrap: wrap;
		padding-top: 0.85rem; border-top: 1px solid var(--ecochitas-border);
	}
	:root[data-theme='dark'] .hero_stats { border-top-color: var(--ecochitas-border); }
	
	.hero_stat { display: flex; flex-direction: column; gap: 0.08rem; }
	.stat_num { font-size: 1rem; font-weight: 800; color: var(--ecochitas-ink); line-height: 1; font-family: 'Sora', sans-serif; }
	.stat_lbl { font-size: 0.68rem; color: var(--ecochitas-muted); line-height: 1; }
	.stat_div { width: 1px; height: 30px; background: var(--ecochitas-border); flex-shrink: 0; }
	:root[data-theme='dark'] .stat_div { background: #262626; }

	/* ── Section commons ── */
	.section_title { font-size: 1.05rem; font-weight: 800; color: var(--ecochitas-ink); }
	.section_sub { font-size: 0.82rem; color: var(--ecochitas-muted); margin-top: 0.22rem; line-height: 1.5; }
	.section_head_row { display: flex; align-items: center; gap: 0.45rem; color: var(--ecochitas-ink); }

	/* ── Ruta del reciclaje ── */
	.ruta_grid {
		display: grid;
		grid-template-columns: 1fr;
		gap: 0.75rem;
		margin-top: 1.1rem;
		position: relative;
	}
	.ruta_step {
		display: flex; flex-direction: column; align-items: center;
		text-align: center; gap: 0.45rem; position: relative;
		padding: 0.85rem 0.6rem;
		background: var(--ecochitas-surface);
		border: 1px solid var(--ecochitas-border);
		border-radius: 1.1rem;
	}
	:root[data-theme='dark'] .ruta_step { background: rgba(255,255,255,0.03); border-color: var(--ecochitas-border); }
	.ruta_icon_wrap {
		position: relative; width: 52px; height: 52px;
		border: 1px solid; border-radius: 1rem;
		display: flex; align-items: center; justify-content: center;
	}
	.ruta_num {
		position: absolute; top: -7px; right: -7px;
		width: 18px; height: 18px; color: #ffffff; background: var(--ecochitas-ink);
		font-size: 0.6rem; font-weight: 900; border-radius: 999px;
		display: flex; align-items: center; justify-content: center;
		font-family: 'Sora', sans-serif;
	}
	:root[data-theme='dark'] .ruta_num { color: var(--ecochitas-background); background: var(--ecochitas-ink); }
	.ruta_title { font-size: 0.8rem; font-weight: 800; line-height: 1.2; color: var(--ecochitas-ink); }
	.ruta_desc { font-size: 0.7rem; color: var(--ecochitas-muted); line-height: 1.4; }
	.ruta_arrow { display: none; }

	/* ── Filter pills ── */
	.filter_scroll { overflow-x: auto; padding-bottom: 0.25rem; margin-top: 1rem; }
	.filter_pills { display: flex; gap: 0.45rem; width: max-content; }
	.filter_pill {
		display: inline-flex; align-items: center; gap: 0.4rem;
		padding: 0.45rem 0.85rem; border-radius: 999px;
		border: 1px solid var(--ecochitas-border);
		background: var(--ecochitas-surface);
		color: var(--ecochitas-muted);
		font-size: 0.78rem; font-weight: 700;
		cursor: pointer; white-space: nowrap;
		transition: border-color 0.15s, background 0.15s, color 0.15s;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
	}
	:root[data-theme='dark'] .filter_pill { background: rgba(255,255,255,0.03); border-color: var(--ecochitas-border); }
	.filter_pill:hover { border-color: var(--ecochitas-leaf); color: var(--ecochitas-leaf); }
	.filter_pill_active {
		border-color: var(--ecochitas-leaf);
		background: rgba(34, 197, 94, 0.1);
		color: var(--ecochitas-leaf);
	}
	:root[data-theme='dark'] .filter_pill_active { border-color: #22c55e; color: #22c55e; }

	/* ── Acopio section ── */
	.acopio_count_badge {
		margin-left: auto;
		font-size: 0.68rem; font-weight: 800;
		background: rgba(34, 197, 94, 0.15);
		color: var(--ecochitas-leaf);
		padding: 0.2rem 0.55rem; border-radius: 999px;
		font-family: 'Sora', sans-serif;
	}
	:root[data-theme='dark'] .acopio_count_badge { color: #4ade80; }
	.acopio_grid {
		display: grid;
		grid-template-columns: 1fr;
		gap: 0.75rem;
		margin-top: 1rem;
	}
	.acopio_card {
		border: 1px solid var(--ecochitas-border);
		border-radius: 1.25rem;
		overflow: hidden;
		background: var(--ecochitas-surface);
		transition: box-shadow 0.15s, transform 0.12s;
	}
	:root[data-theme='dark'] .acopio_card { border-color: var(--ecochitas-border); background: rgba(255,255,255,0.03); }
	.acopio_card:hover {
		box-shadow: 0 6px 24px rgba(0,0,0,0.08);
		transform: translateY(-2px);
	}
	:root[data-theme='dark'] .acopio_card:hover { box-shadow: 0 6px 24px rgba(0,0,0,0.4); }
	.acopio_card_head {
		display: flex; align-items: center; gap: 0.7rem;
		padding: 0.85rem 1rem;
		background: rgba(0,0,0,0.02);
		border-bottom: 1px solid var(--ecochitas-border);
	}
	:root[data-theme='dark'] .acopio_card_head { background: rgba(0,0,0,0.3); border-bottom-color: var(--ecochitas-border); }
	.acopio_card_icon {
		width: 36px; height: 36px; flex-shrink: 0;
		background: var(--ecochitas-surface);
		border: 1px solid var(--ecochitas-border);
		border-radius: 0.7rem; display: flex; align-items: center; justify-content: center;
		color: var(--ecochitas-leaf);
	}
	:root[data-theme='dark'] .acopio_card_icon { background: rgba(255,255,255,0.03); border-color: var(--ecochitas-border); color: var(--ecochitas-leaf); }
	.acopio_card_head_text { flex: 1; min-width: 0; }
	.acopio_name { display: block; font-size: 0.85rem; font-weight: 800; color: var(--ecochitas-ink); font-family: 'Sora', sans-serif; }
	.acopio_zone { font-size: 0.72rem; color: var(--ecochitas-muted); }
	.acopio_status {
		flex-shrink: 0; font-size: 0.65rem; font-weight: 800;
		padding: 0.2rem 0.6rem; border-radius: 999px;
		font-family: 'Sora', sans-serif;
	}
	.acopio_status_activo { background: rgba(34, 197, 94, 0.15); color: var(--ecochitas-leaf); }
	.acopio_status_lleno  { background: rgba(234, 179, 8, 0.15); color: #b45309; }
	.acopio_status_cerrado{ background: rgba(239, 68, 68, 0.15); color: #b91c1c; }
	:root[data-theme='dark'] .acopio_status_activo { color: #4ade80; }
	:root[data-theme='dark'] .acopio_status_lleno { color: #fbbf24; }
	:root[data-theme='dark'] .acopio_status_cerrado { color: #f87171; }
	
	.acopio_card_body { padding: 0.85rem 1rem; display: flex; flex-direction: column; gap: 0.45rem; }
	.acopio_info_row {
		display: flex; align-items: center; gap: 0.4rem;
		font-size: 0.78rem; color: var(--ecochitas-muted);
	}
	.acopio_materials { display: flex; flex-wrap: wrap; gap: 0.35rem; margin-top: 0.2rem; }
	.mat_tag {
		font-size: 0.65rem; font-weight: 700;
		padding: 0.18rem 0.55rem; border-radius: 999px;
		font-family: 'Sora', sans-serif;
	}
	.acopio_capacity_row {
		display: flex; align-items: center; gap: 0.5rem;
		margin-top: 0.5rem;
	}
	.capacity_label { font-size: 0.68rem; color: var(--ecochitas-muted); font-weight: 600; white-space: nowrap; }
	.capacity_bar_wrap { flex: 1; height: 5px; background: var(--ecochitas-border); border-radius: 999px; overflow: hidden; }
	:root[data-theme='dark'] .capacity_bar_wrap { background: #262626; }
	.capacity_bar { height: 100%; border-radius: 999px; transition: width 0.3s ease; }
	.capacity_pct_text { font-size: 0.68rem; font-weight: 800; font-family: 'Sora', sans-serif; white-space: nowrap; min-width: 30px; text-align: right; color: var(--ecochitas-ink); }

	.acopio_card_foot {
		padding: 0.7rem 1rem;
		border-top: 1px solid var(--ecochitas-border);
		display: flex; gap: 0.5rem;
	}
	:root[data-theme='dark'] .acopio_card_foot { border-top-color: #262626; }
	.acopio_loc_btn {
		flex: 1; display: inline-flex; align-items: center; justify-content: center; gap: 0.4rem;
		padding: 0.5rem 0.6rem; border-radius: 0.75rem;
		border: 1px solid var(--ecochitas-border);
		background: var(--ecochitas-surface);
		color: var(--ecochitas-muted);
		font-size: 0.74rem; font-weight: 700; cursor: pointer;
		transition: background 0.15s, border-color 0.15s, color 0.15s;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
	}
	:root[data-theme='dark'] .acopio_loc_btn { border-color: #262626; background: #111111; }
	.acopio_loc_btn:hover { border-color: var(--ecochitas-ink); color: var(--ecochitas-ink); }
	:root[data-theme='dark'] .acopio_loc_btn:hover { border-color: #a3a3a3; color: #ffffff; }
	.acopio_detail_btn {
		flex: 1; display: inline-flex; align-items: center; justify-content: center; gap: 0.4rem;
		padding: 0.5rem 0.6rem; border-radius: 0.75rem;
		border: 1px solid var(--ecochitas-leaf);
		background: rgba(34, 197, 94, 0.1);
		color: var(--ecochitas-leaf);
		font-size: 0.74rem; font-weight: 700; cursor: pointer;
		transition: background 0.15s;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
	}
	:root[data-theme='dark'] .acopio_detail_btn { border-color: #22c55e; color: #22c55e; }
	.acopio_detail_btn:hover { background: rgba(34, 197, 94, 0.15); }
	.ver_todos_btn {
		display: flex; align-items: center; justify-content: center; gap: 0.5rem;
		width: 100%; margin-top: 0.85rem;
		padding: 0.75rem 1rem; border-radius: 1rem;
		border: 1px dashed var(--ecochitas-border);
		background: transparent; color: var(--ecochitas-muted);
		font-size: 0.82rem; font-weight: 700; cursor: pointer;
		transition: border-color 0.15s, color 0.15s;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
	}
	:root[data-theme='dark'] .ver_todos_btn { border-color: #262626; }
	.ver_todos_btn:hover { border-color: var(--ecochitas-leaf); color: var(--ecochitas-leaf); }
	:root[data-theme='dark'] .ver_todos_btn:hover { border-color: #22c55e; color: #22c55e; }
	.ver_menos_btn { border-style: solid; }

	/* ── Category picker ── */
	.cat_picker { display: grid; grid-template-columns: 1fr 1fr; gap: 0.55rem; margin-top: 1.1rem; }
	.cat_btn {
		display: flex; flex-direction: column; align-items: center; text-align: center;
		gap: 0.4rem; padding: 0.85rem 0.5rem 0.7rem;
		border-radius: 1.1rem; border: 1px solid var(--ecochitas-border);
		background: var(--ecochitas-surface); cursor: pointer;
		transition: border-color 0.15s, background 0.15s, transform 0.12s, box-shadow 0.15s;
	}
	:root[data-theme='dark'] .cat_btn { border-color: var(--ecochitas-border); background: rgba(255,255,255,0.03); }
	.cat_btn:hover { border-color: var(--accent-border); background: var(--accent-light); transform: translateY(-1px); }
	.cat_btn_active { border-color: var(--accent); background: var(--accent-light); transform: translateY(-2px); box-shadow: 0 6px 20px color-mix(in srgb, var(--accent) 15%, transparent); }
	.cat_btn_active .cat_btn_name, .cat_btn_active .cat_btn_price, .cat_btn_active .cat_btn_icon { color: var(--accent); }
	.cat_btn_icon { color: var(--accent); display: flex; }
	.cat_btn_name { font-size: 0.72rem; font-weight: 800; color: var(--ecochitas-ink); line-height: 1.2; font-family: 'Sora', sans-serif; }
	.cat_btn_price { font-size: 0.62rem; font-weight: 700; color: var(--accent); background: color-mix(in srgb, var(--accent) 15%, transparent); padding: 0.15rem 0.45rem; border-radius: 999px; }

	/* ── Category detail ── */
	.cat_detail { margin-top: 0.8rem; border: 1px solid var(--accent-border); border-radius: 1.25rem; overflow: hidden; background: #ffffff; }
	:root[data-theme='dark'] .cat_detail { background: #0a0a0a; }
	.cat_detail_head { display: flex; align-items: center; gap: 0.75rem; padding: 0.95rem 1rem; border-bottom: 1px solid var(--accent-border); background: #f9fafb; }
	:root[data-theme='dark'] .cat_detail_head { background: #111111; }
	.cat_detail_icon { flex-shrink: 0; width: 46px; height: 46px; background: #ffffff; border: 1px solid var(--accent-border); border-radius: 0.9rem; display: flex; align-items: center; justify-content: center; color: var(--accent); }
	:root[data-theme='dark'] .cat_detail_icon { background: #0a0a0a; }
	.cat_detail_meta { flex: 1; min-width: 0; }
	.cat_detail_name { font-size: 0.92rem; font-weight: 800; color: var(--accent); line-height: 1.2; }
	.cat_detail_tip { font-size: 0.74rem; color: var(--ecochitas-muted); display: flex; align-items: center; gap: 0.28rem; margin-top: 0.25rem; line-height: 1.4; }
	.cat_price_box { flex-shrink: 0; display: flex; flex-direction: column; align-items: flex-end; gap: 0.05rem; background: var(--ecochitas-surface); border: 1px solid var(--accent-border); border-radius: 0.85rem; padding: 0.45rem 0.75rem; }
	:root[data-theme='dark'] .cat_price_box { background: rgba(255,255,255,0.03); }
	.price_box_lbl { font-size: 0.62rem; color: var(--ecochitas-muted); text-align: right; }
	.price_box_val { font-size: 0.95rem; font-weight: 900; color: var(--accent); line-height: 1.1; font-family: 'Sora', sans-serif; }
	.price_box_unit { font-size: 0.6rem; color: var(--ecochitas-muted); text-align: right; }
	.items_grid { display: grid; grid-template-columns: 1fr; gap: 1px; background: var(--ecochitas-border); }
	:root[data-theme='dark'] .items_grid { background: var(--ecochitas-border); }
	.item_card { display: flex; align-items: center; gap: 0.7rem; padding: 0.8rem 0.9rem; background: var(--ecochitas-surface); transition: background 0.12s; }
	:root[data-theme='dark'] .item_card { background: rgba(255,255,255,0.03); }
	.item_card:hover { background: rgba(0,0,0,0.02); }
	:root[data-theme='dark'] .item_card:hover { background: rgba(255,255,255,0.06); }
	.item_icon { flex-shrink: 0; width: 40px; height: 40px; background: var(--accent-light); border: 1px solid var(--accent-border); border-radius: 0.75rem; display: flex; align-items: center; justify-content: center; color: var(--accent); }
	.item_text { display: flex; flex-direction: column; gap: 0.12rem; min-width: 0; }
	.item_name { font-size: 0.8rem; font-weight: 800; color: var(--ecochitas-ink); line-height: 1.2; font-family: 'Sora', sans-serif; }
	.item_desc { font-size: 0.7rem; color: var(--ecochitas-muted); line-height: 1.3; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

	/* ── Ranking ── */
	.ranking_list { display: flex; flex-direction: column; gap: 0.5rem; margin-top: 0.9rem; }
	.rank_row { display: grid; grid-template-columns: 26px 26px 1fr auto auto; align-items: center; gap: 0.55rem; padding: 0.65rem 0.85rem; background: var(--ecochitas-surface); border: 1px solid var(--ecochitas-border); border-radius: 1rem; }
	:root[data-theme='dark'] .rank_row { background: rgba(255,255,255,0.03); border-color: var(--ecochitas-border); }
	.rank_pos { font-size: 0.72rem; font-weight: 900; color: var(--accent); text-align: center; font-family: 'Sora', sans-serif; }
	.rank_icon { display: flex; color: var(--accent); }
	.rank_name { font-size: 0.82rem; font-weight: 700; color: var(--ecochitas-ink); white-space: nowrap; }
	.rank_bar_track { height: 8px; background: var(--ecochitas-border); border-radius: 999px; overflow: hidden; display: flex; align-items: center; }
	:root[data-theme='dark'] .rank_bar_track { background: rgba(255,255,255,0.1); }
	.rank_bar_fill { height: 100%; background: linear-gradient(90deg, var(--accent), color-mix(in srgb, var(--accent) 60%, white)); border-radius: 999px; transition: width 0.35s ease; }
	.rank_special_label { font-size: 0.66rem; color: var(--ecochitas-muted); font-style: italic; white-space: nowrap; padding-left: 0.3rem; }
	.rank_price { font-size: 0.8rem; font-weight: 800; color: var(--accent); white-space: nowrap; text-align: right; min-width: 76px; font-family: 'Sora', sans-serif; }

	/* ── Rules ── */
	.rules_cols { display: grid; grid-template-columns: 1fr; gap: 0.75rem; margin-top: 0.9rem; }
	.rules_group { border-radius: 1.1rem; overflow: hidden; border: 1px solid var(--ecochitas-border); }
	:root[data-theme='dark'] .rules_group { border-color: var(--ecochitas-border); }
	.rules_group_head { display: flex; align-items: center; gap: 0.5rem; padding: 0.65rem 1rem; font-size: 0.82rem; font-weight: 800; font-family: 'Sora', sans-serif; }
	.rules_head_icon { width: 24px; height: 24px; border-radius: 6px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
	.rules_head_label { font-size: 0.82rem; font-weight: 800; }
	.rules_do .rules_group_head { background: rgba(34, 197, 94, 0.1); border-bottom: 1px solid rgba(34, 197, 94, 0.2); color: var(--ecochitas-leaf); }
	:root[data-theme='dark'] .rules_do .rules_group_head { color: #4ade80; }
	.do_icon { background: rgba(34, 197, 94, 0.2); color: var(--ecochitas-leaf); }
	:root[data-theme='dark'] .do_icon { color: #4ade80; }
	.rules_dont .rules_group_head { background: rgba(239, 68, 68, 0.1); border-bottom: 1px solid rgba(239, 68, 68, 0.2); color: #b91c1c; }
	:root[data-theme='dark'] .rules_dont .rules_group_head { color: #f87171; }
	.dont_icon { background: rgba(239, 68, 68, 0.2); color: #b91c1c; }
	:root[data-theme='dark'] .dont_icon { color: #f87171; }
	.rules_list { margin: 0; padding: 0.75rem 1rem 0.75rem 1.85rem; display: flex; flex-direction: column; gap: 0.5rem; background: var(--ecochitas-surface); }
	:root[data-theme='dark'] .rules_list { background: rgba(255,255,255,0.03); }
	.rules_list li { font-size: 0.82rem; color: var(--ecochitas-ink); line-height: 1.5; }

	/* ── Responsive ── */
	@media screen and (min-width: 640px) {
		.ruta_grid { grid-template-columns: repeat(2, 1fr) !important; }
		.cat_picker { grid-template-columns: repeat(3, 1fr) !important; }
		.items_grid { grid-template-columns: repeat(2, 1fr) !important; }
		.rules_cols { grid-template-columns: 1fr 1fr !important; }
		.acopio_grid { grid-template-columns: repeat(2, 1fr) !important; }
	}
	@media screen and (min-width: 1024px) {
		.hero_title { font-size: 2.1rem !important; }
		.ruta_grid { grid-template-columns: repeat(3, 1fr) !important; }
		.cat_picker { grid-template-columns: repeat(5, 1fr) !important; }
		.items_grid { grid-template-columns: repeat(3, 1fr) !important; }
		.acopio_grid { grid-template-columns: repeat(3, 1fr) !important; }
	}
</style>
