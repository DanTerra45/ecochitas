<script lang="ts">
	// ─── User data ────────────────────────────────────────────────────────────
	const user = {
		name: 'María González', zone: 'Zona Norte', address: 'Casa #234',
		points: 1245, cash_equiv: 12.45, monthly_goal: 2000,
		deliveries: 18, recycled_kg: 42, reports: 6, co2_kg: 84, trees: 4,
		level_num: 4, days_inactive: 0
	};

	// ─── Sistema de niveles ───────────────────────────────────────────────────
	type Nivel = { num: number; nombre: string; min: number; bonus: number };
	const niveles: Nivel[] = [
		{ num: 1,  nombre: 'Curioso Verde',         min: 0,     bonus: 0  },
		{ num: 2,  nombre: 'Eco-Iniciado',           min: 200,   bonus: 2  },
		{ num: 3,  nombre: 'Eco-Ciudadano',          min: 500,   bonus: 4  },
		{ num: 4,  nombre: 'Eco-Ciudadano Pro',      min: 1000,  bonus: 6  },
		{ num: 5,  nombre: 'Guardián Verde',         min: 2000,  bonus: 8  },
		{ num: 6,  nombre: 'Maestro del Reciclaje',  min: 3500,  bonus: 10 },
		{ num: 7,  nombre: 'Defensor Ambiental',     min: 5500,  bonus: 11 },
		{ num: 8,  nombre: 'Eco-Héroe',              min: 8000,  bonus: 12 },
		{ num: 9,  nombre: 'Campeón Verde',          min: 12000, bonus: 13 },
		{ num: 10, nombre: 'Leyenda Ambiental',      min: 18000, bonus: 15 }
	];
	const current_nivel = $derived(niveles.filter((n) => user.points >= n.min).at(-1)!);
	const next_nivel = $derived(niveles.find((n) => n.min > user.points));
	const nivel_pct = $derived(
		next_nivel
			? Math.round(((user.points - current_nivel.min) / (next_nivel.min - current_nivel.min)) * 100)
			: 100
	);
	const monthly_pct = $derived(Math.min((user.points / user.monthly_goal) * 100, 100));

	// ─── Top stats ───────────────────────────────────────────────────────────
	const top_stats = $derived([
		{ label: 'EcoPoints', value: user.points.toLocaleString(), icon: 'M12 3l2.7 5.5 6.1.9-4.4 4.3 1 6.1L12 17l-5.4 2.8 1-6.1L3.2 9.4l6.1-.9z', color: '#d97706', bg: '#fffbeb', border: '#fde68a' },
		{ label: 'Equiv. Bs', value: `Bs. ${user.cash_equiv}`, icon: 'M12 1v22M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6', color: '#16a34a', bg: '#f0fdf4', border: '#bbf7d0' },
		{ label: 'Reciclado', value: `${user.recycled_kg}kg`, icon: 'M3 12a9 9 0 0 1 15.2-6.4M18.2 5.6V2.4h-3.2M21 12a9 9 0 0 1-15.2 6.4M5.8 18.4v3.2H9', color: '#0284c7', bg: '#eff6ff', border: '#bfdbfe' },
		{ label: 'Reportes', value: user.reports.toString(), icon: 'M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0zM12 9v4M12 17h.01', color: '#dc2626', bg: '#fef2f2', border: '#fecaca' },
		{ label: 'Nivel', value: `Nv. ${current_nivel.num}`, icon: 'M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z', color: '#7c3aed', bg: '#f5f3ff', border: '#ddd6fe' }
	]);

	// ─── Gamificación ─────────────────────────────────────────────────────────
	const streak = { days: 7, best: 14 };

	type Insignia = { id: string; title: string; desc: string; emoji: string; color: string; bg: string; border: string; earned: boolean; date?: string };
	const insignias: Insignia[] = [
		{ id: 'i1', title: 'Primera Entrega',       desc: 'Realizaste tu primera entrega de reciclaje', emoji: '🌱', color: '#16a34a', bg: '#f0fdf4', border: '#bbf7d0', earned: true,  date: 'Ene 2026' },
		{ id: 'i2', title: '10 Entregas',            desc: 'Completaste 10 entregas de materiales',      emoji: '♻️', color: '#0284c7', bg: '#eff6ff', border: '#bfdbfe', earned: true,  date: 'Mar 2026' },
		{ id: 'i3', title: 'Maestro del Vidrio',     desc: 'Entregaste más de 10 kg de vidrio',          emoji: '🏆', color: '#d97706', bg: '#fffbeb', border: '#fde68a', earned: true,  date: 'Abr 2026' },
		{ id: 'i4', title: 'Eco-Héroe del Mes',      desc: 'Fuiste el mayor contribuyente de tu zona',   emoji: '⭐', color: '#7c3aed', bg: '#f5f3ff', border: '#ddd6fe', earned: false },
		{ id: 'i5', title: 'Reporte Procesado',      desc: 'Tu reporte fue atendido por la GAMC',        emoji: '📋', color: '#dc2626', bg: '#fef2f2', border: '#fecaca', earned: true,  date: 'May 2026' }
	];

	type Actividad = { id: number; tipo: 'entrega' | 'reporte' | 'canje' | 'nivel'; desc: string; puntos: number; fecha: string };
	const actividad_reciente: Actividad[] = [
		{ id: 1, tipo: 'entrega', desc: 'Entrega de 3.2 kg de plástico PET',        puntos: 128,  fecha: '15 May' },
		{ id: 2, tipo: 'reporte', desc: 'Reporte en Av. América procesado',          puntos: 50,   fecha: '13 May' },
		{ id: 3, tipo: 'entrega', desc: 'Entrega de 2 kg de papel y cartón',         puntos: 70,   fecha: '11 May' },
		{ id: 4, tipo: 'nivel',   desc: '¡Subiste al Nivel 4 — Eco-Ciudadano Pro!', puntos: 0,    fecha: '08 May' },
		{ id: 5, tipo: 'canje',   desc: 'Canjeaste: Café Orgánico Boliviano',        puntos: -150, fecha: '05 May' }
	];
	const actividad_icons: Record<Actividad['tipo'], string> = {
		entrega: 'M3 12a9 9 0 0 1 15.2-6.4M18.2 5.6V2.4h-3.2M21 12a9 9 0 0 1-15.2 6.4M5.8 18.4v3.2H9',
		reporte: 'M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0zM12 9v4M12 17h.01',
		canje: 'M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4zM3 6h18M16 10a4 4 0 0 1-8 0',
		nivel: 'M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z'
	};
	const actividad_colors: Record<Actividad['tipo'], string> = {
		entrega: '#16a34a', reporte: '#dc2626', canje: '#7c3aed', nivel: '#d97706'
	};

	// ─── Validador de reciclado ───────────────────────────────────────────────
	type ValidacionEstado = 'pendiente' | 'validado' | 'rechazado';
	type Validacion = { id: number; material: string; peso: string; estado: ValidacionEstado; fecha: string; puntos_obtenidos?: number };
	const validaciones: Validacion[] = [
		{ id: 1, material: 'Plástico PET',    peso: '3.2 kg', estado: 'validado',  fecha: '15 May 2026', puntos_obtenidos: 128 },
		{ id: 2, material: 'Papel y Cartón',  peso: '2.0 kg', estado: 'pendiente', fecha: '17 May 2026' }
	];
	const val_label: Record<ValidacionEstado, string> = { pendiente: 'Pendiente', validado: 'Validado', rechazado: 'Rechazado' };

	// ─── Condominio ───────────────────────────────────────────────────────────
	const building = { name: 'Torres Sofer', address: 'Av. América', units: 24, ranking: 3, reward_pct: 15, community_points: 8450, community_goal: 12000 };
	const community_pct = $derived(Math.min((building.community_points / building.community_goal) * 100, 100));
	const top_contributors = [
		{ unit: 'Depto 301', points: 850 }, { unit: 'Depto 205', points: 720 },
		{ unit: 'Depto 102', points: 680 }, { unit: 'Depto 401', points: 540 },
		{ unit: 'Depto 104', points: 480 }
	];

	// ─── Tienda ───────────────────────────────────────────────────────────────
	type ProductCategory = 'cafeteria' | 'ecostore' | 'restaurante' | 'bolivia' | 'transporte' | 'educacion';
	type StoreItem = { id: string; name: string; partner: string; category: ProductCategory; points: number; discount_pct: number; type: 'descuento' | 'producto'; description: string };

	const store_items: StoreItem[] = [
		{ id: 's1', name: 'Café Orgánico Boliviano',  partner: 'Café La Casona',       category: 'cafeteria',   points: 150, discount_pct: 20,  type: 'descuento', description: 'Cualquier bebida del menú' },
		{ id: 's2', name: 'Bolsa Eco-Reutilizable',   partner: 'EcoBolivia',            category: 'bolivia',     points: 80,  discount_pct: 100, type: 'producto',  description: 'Bolsa de tela certificada hecha en Bolivia' },
		{ id: 's3', name: 'Menú Eco-Lunch',           partner: 'Restaurante Verde',     category: 'restaurante', points: 200, discount_pct: 15,  type: 'descuento', description: 'Menú del día con ingredientes ecológicos' },
		{ id: 's4', name: 'Kit de Reciclaje',         partner: 'EcoTienda Naturaleza',  category: 'ecostore',    points: 320, discount_pct: 100, type: 'producto',  description: 'Contenedores clasificadores x 4 colores' },
		{ id: 's5', name: 'Almuerzo Familiar',        partner: 'Quinoa & Co.',          category: 'restaurante', points: 400, discount_pct: 25,  type: 'descuento', description: 'Para 2 personas, incluye postres orgánicos' },
		{ id: 's6', name: 'Planta de Interior',       partner: 'Vivero EcoVerde',       category: 'ecostore',    points: 250, discount_pct: 100, type: 'producto',  description: 'Planta nativa de Bolivia, certificada' },
		{ id: 's7', name: 'Bicicleta Compartida',     partner: 'BiciVerde Cbba',        category: 'transporte',  points: 300, discount_pct: 50,  type: 'descuento', description: '1 mes de acceso a bicicletas compartidas' },
		{ id: 's8', name: 'Taller de Compostaje',     partner: 'EcoAprender',           category: 'educacion',   points: 180, discount_pct: 100, type: 'producto',  description: 'Curso presencial de compostaje doméstico' }
	];

	const category_labels: Record<ProductCategory, { label: string; color: string; bg: string }> = {
		cafeteria:   { label: 'Cafetería',        color: '#92400e', bg: '#fef3c7' },
		ecostore:    { label: 'Eco-Tienda',        color: '#065f46', bg: '#d1fae5' },
		restaurante: { label: 'Restaurante',       color: '#7c3aed', bg: '#ede9fe' },
		bolivia:     { label: 'Hecho en Bolivia',  color: '#1e40af', bg: '#dbeafe' },
		transporte:  { label: 'Transporte',        color: '#0c4a6e', bg: '#e0f2fe' },
		educacion:   { label: 'Educación',         color: '#065f46', bg: '#ecfdf5' }
	};

	let store_filter = $state<ProductCategory | 'todos'>('todos');
	const visible_store = $derived(
		store_filter === 'todos' ? store_items : store_items.filter((s) => s.category === store_filter)
	);

	// ─── Beneficios municipales ───────────────────────────────────────────────
	type BenefitStatus = 'disponible' | 'por_alcanzar' | 'obtenido';
	type MunicipalBenefit = { id: string; title: string; desc: string; points_required: number; icon: string; color: string; status: BenefitStatus };

	const municipal_benefits: MunicipalBenefit[] = [
		{ id: 'b1', title: 'Descuento Predial 5%',       desc: 'Descuento simbólico sobre impuesto a la propiedad inmueble para el siguiente semestre.',  points_required: 1000, icon: 'M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z', color: '#16a34a', status: 'disponible' },
		{ id: 'b2', title: 'Certificado Eco-Responsable', desc: 'Reconocimiento público emitido por la GAMC que acredita tu compromiso ambiental.',         points_required: 800,  icon: 'M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z',    color: '#0284c7', status: 'disponible' },
		{ id: 'b3', title: 'Prioridad en Trámites',      desc: 'Acceso preferente en ventanillas municipales de servicios ciudadanos de la GAMC.',          points_required: 1500, icon: 'M9 11l3 3L22 4M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11', color: '#7c3aed', status: 'por_alcanzar' },
		{ id: 'b4', title: 'Fondo de Mejora Barrial',    desc: 'Tu zona acumula puntos comunitarios para financiar parques e iluminación pública.',          points_required: 2000, icon: 'M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75', color: '#d97706', status: 'por_alcanzar' },
		{ id: 'b5', title: 'Descuento en Tasas 10%',     desc: 'Descuento aplicable sobre tasas de servicios de agua y recolección de residuos.',           points_required: 3000, icon: 'M12 1v22M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6', color: '#dc2626', status: 'por_alcanzar' }
	];
	const benefit_status_labels: Record<BenefitStatus, string> = {
		disponible: 'Disponible', por_alcanzar: 'Por alcanzar', obtenido: 'Obtenido'
	};

	// ─── Tab state ────────────────────────────────────────────────────────────
	let active_tab: 'inicio' | 'individual' | 'condominio' | 'tienda' | 'municipal' = $state('inicio');

	// ─── Impacto comunitario mock (pestaña Inicio) ─────────────────────────────
	const impact_stats = [
		{ label: 'Usuarios activos',    value: '3,840', icon: 'M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75', color: '#16a34a' },
		{ label: 'EcoPoints generados', value: '2.1M',  icon: 'M12 3l2.7 5.5 6.1.9-4.4 4.3 1 6.1L12 17l-5.4 2.8 1-6.1L3.2 9.4l6.1-.9z',                       color: '#d97706' },
		{ label: 'Kg reciclados',       value: '98 T',  icon: 'M3 12a9 9 0 0 1 15.2-6.4M18.2 5.6V2.4h-3.2M21 12a9 9 0 0 1-15.2 6.4M5.8 18.4v3.2H9',          color: '#0284c7' },
		{ label: 'Reportes atendidos',  value: '1,204', icon: 'M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0zM12 9v4M12 17h.01', color: '#dc2626' },
		{ label: 'Puntos de acopio',    value: '24',    icon: 'M12 21s-8-7.5-8-12a8 8 0 0 1 16 0c0 4.5-8 12-8 12z',                                             color: '#7c3aed' },
		{ label: 'Zonas beneficiadas',  value: '18',    icon: 'M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z',                                                 color: '#059669' }
	];
</script>

<!-- ── Top stats cards ── -->
<div class="stats_bar">
	{#each top_stats as s (s.label)}
		<div class="top_stat_card" style="background:{s.bg};border-color:{s.border}">
			<div class="top_stat_icon" style="color:{s.color}">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:15px;height:15px" aria-hidden="true"><path d={s.icon}/></svg>
			</div>
			<strong class="top_stat_val" style="color:{s.color}">{s.value}</strong>
			<span class="top_stat_lbl">{s.label}</span>
		</div>
	{/each}
</div>

<!-- ── Profile header ── -->
<div class="profile_header panel">
	<div class="avatar_wrap">
		<span class="avatar_initials">MG</span>
		<div class="avatar_level_ring" title="Nivel {current_nivel.num}"></div>
	</div>
	<div class="header_info">
		<h1 class="user_name">{user.name}</h1>
		<p class="user_sub">{user.zone} · {user.address}</p>
		<div class="header_badges">
			<span class="user_level">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:11px;height:11px" aria-hidden="true"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
				Nv. {current_nivel.num} — {current_nivel.nombre}
			</span>
			<span class="streak_chip">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:11px;height:11px" aria-hidden="true"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>
				{streak.days} días de racha
			</span>
		</div>
	</div>
</div>

<!-- ── Tab switcher ── -->
<div class="tab_switcher">
	{#each ([
		['inicio',     'Inicio',     'M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2zM9 22V12h6v10'],
		['individual', 'Yo',         'M4 11.5L12 5l8 6.5M6 10.8V20h12v-9.2M10 20v-4h4v4'],
		['condominio', 'Edificio',   'M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2zM9 22V12h6v10'],
		['tienda',     'Tienda',     'M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4zM3 6h18M16 10a4 4 0 0 1-8 0'],
		['municipal',  'GAMC',       'M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z']
	]) as [id, label, icon]}
		<button class="tab_btn" class:tab_active={active_tab === id} onclick={() => (active_tab = id as typeof active_tab)}>
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:13px;height:13px" aria-hidden="true"><path d={icon}/></svg>
			{label}
		</button>
	{/each}
</div>

<!-- ══════════ INICIO ══════════ -->
{#if active_tab === 'inicio'}

<div class="inicio_hero">
	<div class="inicio_hero_eyebrow">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:12px;height:12px" aria-hidden="true"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
		EcoChitas · GAMC Cochabamba
	</div>
	<h2 class="inicio_hero_title">Bienvenida,<br /><span class="inicio_name">{user.name.split(' ')[0]}</span> 👋</h2>
	<p class="inicio_hero_sub">Sos parte del movimiento ambiental más grande de Cochabamba. Cada entrega suma, cada reporte importa.</p>
</div>

<!-- Impacto stats -->
<div class="impact_grid">
	{#each impact_stats as s (s.label)}
		<div class="impact_stat_card">
			<div class="impact_stat_icon" style="color:{s.color}">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:16px;height:16px" aria-hidden="true"><path d={s.icon}/></svg>
			</div>
			<strong class="impact_stat_val" style="color:{s.color}">{s.value}</strong>
			<span class="impact_stat_lbl">{s.label}</span>
		</div>
	{/each}
</div>

<!-- Quiénes somos -->
<div class="panel inicio_section">
	<div class="inicio_section_head">
		<div class="inicio_section_icon" style="background:#f0fdf4;border-color:#bbf7d0;color:#16a34a">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:16px;height:16px" aria-hidden="true"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"/><circle cx="9" cy="7" r="4"/></svg>
		</div>
		<h3 class="inicio_section_title">Quiénes Somos</h3>
	</div>
	<p class="inicio_section_text">EcoChitas es una iniciativa tecnológica nacida en Cochabamba con la misión de modernizar la gestión de residuos sólidos. Combinamos sensores IoT, GPS y software avanzado para brindar herramientas a operadores municipales y motivar a los ciudadanos a reciclar.</p>
</div>

<!-- Misión -->
<div class="panel inicio_section">
	<div class="inicio_section_head">
		<div class="inicio_section_icon" style="background:#eff6ff;border-color:#bfdbfe;color:#0284c7">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:16px;height:16px" aria-hidden="true"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
		</div>
		<h3 class="inicio_section_title">Nuestra Misión</h3>
	</div>
	<p class="inicio_section_text">Digitalizamos el ciclo completo de recolección: desde el monitoreo en tiempo real de camiones y sensores de llenado, hasta la validación de materiales reciclables entregados por vecinos, otorgándoles puntos canjeables por beneficios locales.</p>
</div>

<!-- Cómo te beneficia -->
<div class="inicio_benefits_card">
	<h3 class="inicio_benefits_title">¿Cómo te beneficia?</h3>
	<div class="inicio_benefits_list">
		{#each [
			{ icon: 'M12 3l2.7 5.5 6.1.9-4.4 4.3 1 6.1L12 17l-5.4 2.8 1-6.1L3.2 9.4l6.1-.9z', text: 'Ganás EcoPoints por cada entrega de reciclaje verificada' },
			{ icon: 'M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4zM3 6h18M16 10a4 4 0 0 1-8 0', text: 'Canjeás puntos por descuentos en tiendas y restaurantes locales' },
			{ icon: 'M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z', text: 'Accedés a beneficios municipales de la GAMC por participación activa' },
			{ icon: 'M3 12a9 9 0 0 1 15.2-6.4M18.2 5.6V2.4h-3.2M21 12a9 9 0 0 1-15.2 6.4M5.8 18.4v3.2H9', text: 'Contribuís directamente a reducir la contaminación en tu ciudad' }
		] as b}
			<div class="inicio_benefit_item">
				<div class="inicio_benefit_dot">
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:13px;height:13px" aria-hidden="true"><path d={b.icon}/></svg>
				</div>
				<span>{b.text}</span>
			</div>
		{/each}
	</div>
</div>

<!-- Educación ambiental -->
<div class="panel">
	<h3 class="section_title" style="margin-bottom:0.85rem">💡 Educación Ambiental</h3>
	<div class="edu_list">
		{#each [
			{ tip: 'Enjuagá los envases plásticos antes de llevarlos al punto de acopio para evitar el rechazo.', color: '#0284c7', bg: '#eff6ff' },
			{ tip: 'El papel húmedo no se recicla. Guardá tu papel y cartón en lugar seco y entregalos limpio.', color: '#92400e', bg: '#fef3c7' },
			{ tip: 'Las pilas y baterías son residuos especiales — solo en puntos de acopio autorizados por la GAMC.', color: '#dc2626', bg: '#fef2f2' }
		] as e, i}
			<div class="edu_item" style="background:{e.bg};border-left:3px solid {e.color}">
				<span class="edu_num" style="color:{e.color}">{i + 1}</span>
				<p class="edu_text" style="color:{e.color}">{e.tip}</p>
			</div>
		{/each}
	</div>
</div>

<!-- Mapa preview CTA -->
<a href="/map" class="mapa_preview_cta">
	<div class="mapa_preview_left">
		<div class="mapa_preview_icon">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:22px;height:22px" aria-hidden="true"><polygon points="3 11 22 2 13 21 11 13 3 11"/></svg>
		</div>
		<div>
			<strong class="mapa_preview_title">Ver Mapa Operativo</strong>
			<p class="mapa_preview_sub">Camiones, basureros y puntos de acopio en tiempo real</p>
		</div>
	</div>
	<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" style="width:16px;height:16px;flex-shrink:0" aria-hidden="true"><polyline points="9 18 15 12 9 6"/></svg>
</a>

<!-- ══════════ INDIVIDUAL ══════════ -->
{:else if active_tab === 'individual'}

<!-- Eco-Wallet -->
<div class="wallet_card">
	<div class="wallet_header">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:17px;height:17px;opacity:0.75" aria-hidden="true"><rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/></svg>
		<span class="wallet_label">Eco-Wallet</span>
		<span class="wallet_mvp_tag">Gamificación</span>
	</div>
	<p class="ecopoints_sublabel">EcoPoints acumulados</p>
	<div class="ecopoints_row">
		<strong class="points_num">{user.points.toLocaleString()}</strong>
		<span class="ecopoints_unit">EP</span>
	</div>
	<div class="cash_disclaimer">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:11px;height:11px;flex-shrink:0;margin-top:1px;opacity:0.7" aria-hidden="true"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
		<span>≈ <strong>Bs. {user.cash_equiv}</strong> equivalencia referencial. Esto NO es dinero real — es una estimación si el material hubiese sido vendido directamente por vos.</span>
	</div>
	<button class="redeem_btn">Canjear EcoPoints</button>
</div>

<!-- Nivel -->
<div class="panel nivel_panel">
	<div class="nivel_head_row">
		<div class="nivel_badge_icon">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:15px;height:15px" aria-hidden="true"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
		</div>
		<div>
			<h2 class="section_title">Sistema de Niveles</h2>
			<p class="nivel_sub">Nivel {current_nivel.num} — {current_nivel.nombre} · +{current_nivel.bonus}% bonus</p>
		</div>
	</div>
	<div class="nivel_progress_row">
		<span class="nivel_pts_text">{user.points.toLocaleString()} pts</span>
		{#if next_nivel}
			<span class="nivel_next_text">→ {next_nivel.min.toLocaleString()} pts para Nv. {next_nivel.num}</span>
		{:else}
			<span class="nivel_next_text" style="color:var(--ecochitas-leaf)">¡Nivel máximo alcanzado!</span>
		{/if}
	</div>
	<div class="nivel_track"><div class="nivel_fill" style="width:{nivel_pct}%"></div></div>
	{#if next_nivel}
		<p class="nivel_hint">Faltan <strong>{(next_nivel.min - user.points).toLocaleString()} pts</strong> para {next_nivel.nombre}</p>
	{/if}
	<!-- Tabla bonus compacta -->
	<div class="nivel_table_scroll">
		<div class="nivel_table">
			{#each niveles as n (n.num)}
				<div class="nivel_row" class:nivel_row_active={n.num === current_nivel.num} class:nivel_row_locked={user.points < n.min}>
					<span class="nivel_row_num">Nv.{n.num}</span>
					<span class="nivel_row_name">{n.nombre}</span>
					<span class="nivel_row_bonus">+{n.bonus}%</span>
				</div>
			{/each}
		</div>
	</div>
	<div class="nivel_inactivity_warn">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:12px;height:12px;flex-shrink:0;margin-top:1px" aria-hidden="true"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0zM12 9v4M12 17h.01"/></svg>
		<span>La inactividad prolongada (más de 30 días sin entregas) puede reducir tu nivel. ¡Mantenete activo!</span>
	</div>
</div>

<!-- Progreso mensual -->
<div class="panel">
	<div class="section_row">
		<svg viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:18px;height:18px;flex-shrink:0" aria-hidden="true"><path d="M12 3l2.7 5.5 6.1.9-4.4 4.3 1 6.1L12 17l-5.4 2.8 1-6.1L3.2 9.4l6.1-.9z"/></svg>
		<h2 class="section_title">Progreso Este Mes</h2>
	</div>
	<div class="progress_row">
		<span class="progress_label">Meta Mensual</span>
		<span class="progress_val" style="color:var(--ecochitas-sky)">{user.points.toLocaleString()} / {user.monthly_goal.toLocaleString()}</span>
	</div>
	<div class="progress_track"><div class="progress_fill" style="width:{monthly_pct}%"></div></div>
	<div class="stats_duo">
		<div class="stat_pill stat_pill_green"><strong class="stat_big">{user.deliveries}</strong><span class="stat_small">Entregas</span></div>
		<div class="stat_pill stat_pill_blue"><strong class="stat_big stat_blue">{user.recycled_kg}kg</strong><span class="stat_small">Reciclado</span></div>
	</div>

	<!-- Validador de reciclado -->
	<div class="validator_section">
		<div class="validator_head">
			<div class="validator_title_row">
				<svg viewBox="0 0 24 24" fill="none" stroke="var(--ecochitas-leaf)" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:15px;height:15px;flex-shrink:0" aria-hidden="true"><rect x="3" y="3" width="18" height="18" rx="2.5"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
				<strong class="validator_title">Validador de Reciclado</strong>
			</div>
			<span class="validator_quota">2 fotos/semana máx.</span>
		</div>
		<p class="validator_sub">Subí fotos de tus entregas como evidencia. Un operador las revisará y acreditará los puntos.</p>

		<div class="validaciones_list">
			{#each validaciones as v (v.id)}
				<div class="validacion_item">
					<div class="val_icon val_icon_{v.estado}">
						{#if v.estado === 'validado'}
							<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.2" style="width:13px;height:13px" aria-hidden="true"><polyline points="20 6 9 17 4 12"/></svg>
						{:else if v.estado === 'rechazado'}
							<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.2" style="width:13px;height:13px" aria-hidden="true"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
						{:else}
							<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:13px;height:13px" aria-hidden="true"><circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/></svg>
						{/if}
					</div>
					<div class="val_info">
						<strong class="val_material">{v.material}</strong>
						<span class="val_meta">{v.peso} · {v.fecha}</span>
					</div>
					<div class="val_right">
						<span class="val_status val_status_{v.estado}">{val_label[v.estado]}</span>
						{#if v.puntos_obtenidos}
							<span class="val_pts">+{v.puntos_obtenidos} EP</span>
						{/if}
					</div>
				</div>
			{/each}
		</div>

		<label class="validator_upload_btn" for="val_photo_input">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:14px;height:14px;flex-shrink:0" aria-hidden="true"><rect x="3" y="3" width="18" height="18" rx="2.5"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
			Subir evidencia fotográfica
		</label>
		<input type="file" id="val_photo_input" accept="image/*" class="upload_input_hidden" />
	</div>
</div>

<!-- Impacto ambiental -->
<div class="impact_card">
	<div class="impact_header"><span class="impact_emoji" aria-hidden="true">🌱</span><h2 class="impact_title">Impacto Ambiental</h2></div>
	<p class="impact_desc">Tu contribución este mes evitó la emisión de</p>
	<strong class="impact_co2">{user.co2_kg}kg CO<sub>2</sub></strong>
	<p class="impact_equiv">Equivalente a plantar {user.trees} árboles nativos</p>
</div>

<!-- Insignias -->
<div class="panel">
	<div class="section_row">
		<svg viewBox="0 0 24 24" fill="none" stroke="#d97706" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:18px;height:18px;flex-shrink:0" aria-hidden="true"><circle cx="12" cy="8" r="6"/><path d="M15.477 12.89L17 22l-5-3-5 3 1.523-9.11"/></svg>
		<h2 class="section_title">Mis Insignias</h2>
	</div>
	<div class="insignias_grid">
		{#each insignias as ins (ins.id)}
			<div class="insignia_card" class:insignia_locked={!ins.earned} style="background:{ins.earned ? ins.bg : 'var(--ecochitas-surface)'};border-color:{ins.earned ? ins.border : 'var(--ecochitas-border)'}">
				<span class="insignia_emoji" class:insignia_emoji_locked={!ins.earned}>{ins.emoji}</span>
				<strong class="insignia_title" style="color:{ins.earned ? ins.color : 'var(--ecochitas-muted)'}">{ins.title}</strong>
				<span class="insignia_date">{ins.earned ? (ins.date ?? '') : 'Bloqueada'}</span>
			</div>
		{/each}
	</div>
</div>

<!-- Racha + actividad -->
<div class="panel">
	<div class="streak_row">
		<div class="streak_card">
			<svg viewBox="0 0 24 24" fill="none" stroke="#ef4444" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:18px;height:18px" aria-hidden="true"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>
			<div>
				<strong class="streak_val">{streak.days} días</strong>
				<span class="streak_lbl">Racha actual</span>
			</div>
		</div>
		<div class="streak_card streak_card_muted">
			<svg viewBox="0 0 24 24" fill="none" stroke="var(--ecochitas-muted)" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:16px;height:16px" aria-hidden="true"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
			<div>
				<strong class="streak_val streak_val_muted">{streak.best} días</strong>
				<span class="streak_lbl">Mejor racha</span>
			</div>
		</div>
	</div>
	<h2 class="section_title" style="margin-top:1.1rem;margin-bottom:0.75rem">Actividad Reciente</h2>
	<div class="activity_list">
		{#each actividad_reciente as act (act.id)}
			<div class="activity_item">
				<div class="activity_icon" style="background:color-mix(in srgb,{actividad_colors[act.tipo]} 12%,transparent);border-color:color-mix(in srgb,{actividad_colors[act.tipo]} 25%,transparent);color:{actividad_colors[act.tipo]}">
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:14px;height:14px" aria-hidden="true"><path d={actividad_icons[act.tipo]}/></svg>
				</div>
				<div class="activity_info">
					<span class="activity_desc">{act.desc}</span>
					<span class="activity_date">{act.fecha}</span>
				</div>
				<span class="activity_pts" class:activity_pts_neg={act.puntos < 0} class:activity_pts_zero={act.puntos === 0}>
					{act.puntos > 0 ? `+${act.puntos} EP` : act.puntos < 0 ? `${act.puntos} EP` : '—'}
				</span>
			</div>
		{/each}
	</div>
</div>

<!-- ══════════ CONDOMINIO ══════════ -->
{:else if active_tab === 'condominio'}

<div class="panel building_header_card">
	<div class="building_icon_wrap">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:20px;height:20px" aria-hidden="true"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M9 3v18M15 3v18M3 9h18M3 15h18"/></svg>
	</div>
	<div><h2 class="building_name">{building.name}</h2><p class="building_sub">{building.address} · {building.units} departamentos</p></div>
</div>

<div class="reward_card">
	<div class="reward_header">
		<div class="reward_icon">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:16px;height:16px" aria-hidden="true"><circle cx="12" cy="8" r="6"/><path d="M15.477 12.89L17 22l-5-3-5 3 1.523-9.11"/></svg>
		</div>
		<h3 class="reward_title">Recompensa de Condominio Top</h3>
	</div>
	<p class="reward_desc">{building.reward_pct}% de descuento en impuesto predial para todos los residentes este mes</p>
	<div class="reward_rank">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.2" style="width:12px;height:12px;flex-shrink:0" aria-hidden="true"><polyline points="22 7 13.5 15.5 8.5 10.5 2 17"/><polyline points="16 7 22 7 22 13"/></svg>
		Ranking Actual: #{building.ranking}
	</div>
</div>

<div class="panel">
	<div class="progress_row">
		<span class="progress_label">Meta Comunitaria Mensual</span>
		<span class="progress_val" style="color:var(--ecochitas-sky)">{building.community_points.toLocaleString()} / {building.community_goal.toLocaleString()}</span>
	</div>
	<div class="progress_track"><div class="progress_fill" style="width:{community_pct}%"></div></div>
	<p class="progress_hint">Faltan {(building.community_goal - building.community_points).toLocaleString()} puntos para alcanzar la meta</p>
</div>

<div class="panel">
	<h2 class="section_title" style="margin-bottom:0.85rem">Top Contribuyentes del Edificio</h2>
	<div class="rank_list">
		{#each top_contributors as contrib, i (contrib.unit)}
			<div class="rank_row" class:rank_row_gold={i === 0}>
				<div class="rank_badge" class:rank_badge_gold={i === 0} class:rank_badge_silver={i === 1} class:rank_badge_bronze={i === 2}>{i + 1}</div>
				<span class="rank_unit">{contrib.unit}</span>
				<div class="rank_right">
					<strong class="rank_pts" class:rank_pts_gold={i === 0}>{contrib.points.toLocaleString()}</strong>
					<span class="rank_pts_lbl">EcoPoints</span>
				</div>
			</div>
		{/each}
	</div>
</div>

<div class="benefits_card">
	<div class="benefits_header"><span class="benefits_emoji" aria-hidden="true">🏆</span><h2 class="benefits_title">Beneficios Comunitarios</h2></div>
	<ul class="benefits_list">
		{#each ['Descuentos grupales en impuestos municipales', 'Reconocimiento como Edificio Eco-Responsable GAMC', 'Acceso prioritario a programas de mejora urbana'] as b}
			<li class="benefit_item">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.2" style="width:11px;height:11px;flex-shrink:0;margin-top:2px;opacity:0.8" aria-hidden="true"><polyline points="20 6 9 17 4 12"/></svg>
				{b}
			</li>
		{/each}
	</ul>
</div>

<!-- ══════════ TIENDA ══════════ -->
{:else if active_tab === 'tienda'}

<div class="store_hero">
	<div class="store_hero_icon" aria-hidden="true">🛍️</div>
	<div>
		<h2 class="store_hero_title">Tienda Ecológica</h2>
		<p class="store_hero_sub">Canjeá tus EcoPoints por descuentos y productos en negocios locales.</p>
	</div>
	<div class="store_pts_chip">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:12px;height:12px" aria-hidden="true"><path d="M12 3l2.7 5.5 6.1.9-4.4 4.3 1 6.1L12 17l-5.4 2.8 1-6.1L3.2 9.4l6.1-.9z"/></svg>
		{user.points.toLocaleString()} EP
	</div>
</div>

<div class="filter_scroll">
	<div class="filter_pills">
		{#each ([['todos','Todos'],['cafeteria','Cafeterías'],['restaurante','Restaurantes'],['ecostore','Eco-Tiendas'],['bolivia','Bolivia'],['transporte','Transporte'],['educacion','Educación']]) as [id, label]}
			<button class="filter_pill" class:filter_pill_active={store_filter === id} onclick={() => (store_filter = id as typeof store_filter)}>
				{label}
			</button>
		{/each}
	</div>
</div>

<div class="store_grid">
	{#each visible_store as item (item.id)}
		<div class="store_card">
			<div class="store_logo_area">
				<div class="store_logo_placeholder">
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" style="width:20px;height:20px;opacity:0.35" aria-hidden="true"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4zM3 6h18"/></svg>
				</div>
				<div class="store_card_badges">
					<span class="store_cat_tag" style="background:{category_labels[item.category].bg};color:{category_labels[item.category].color}">
						{category_labels[item.category].label}
					</span>
					<span class="store_type_badge" class:store_type_product={item.type === 'producto'}>
						{item.type === 'descuento' ? `${item.discount_pct}% off` : 'Gratis'}
					</span>
				</div>
			</div>
			<h3 class="store_item_name">{item.name}</h3>
			<p class="store_partner">{item.partner}</p>
			<p class="store_desc">{item.description}</p>
			<div class="store_card_foot">
				<span class="store_pts_needed">
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:11px;height:11px;flex-shrink:0" aria-hidden="true"><path d="M12 3l2.7 5.5 6.1.9-4.4 4.3 1 6.1L12 17l-5.4 2.8 1-6.1L3.2 9.4l6.1-.9z"/></svg>
					{item.points} EP
				</span>
				<button class="store_redeem_btn" disabled={user.points < item.points}>
					{user.points >= item.points ? 'Canjear' : 'Insuficiente'}
				</button>
			</div>
		</div>
	{/each}
</div>

<!-- ══════════ MUNICIPAL ══════════ -->
{:else if active_tab === 'municipal'}

<div class="muni_hero">
	<div class="muni_hero_icon" aria-hidden="true">🏛️</div>
	<div>
		<h2 class="muni_hero_title">Beneficios Municipales</h2>
		<p class="muni_hero_sub">La GAMC recompensa a ciudadanos activos con descuentos y reconocimientos oficiales.</p>
	</div>
</div>

<div class="muni_mvp_notice">
	<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:13px;height:13px;flex-shrink:0;margin-top:1px" aria-hidden="true"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
	<span><strong>Propuesta MVP:</strong> Los beneficios municipales mostrados son propuestas/recompensas modelo para demostración del sistema. Su implementación real depende de la GAMC.</span>
</div>

<div class="panel muni_progress_card">
	<div class="muni_prog_row">
		<div>
			<p class="muni_prog_label">Tus EcoPoints actuales</p>
			<strong class="muni_prog_pts">{user.points.toLocaleString()} EP</strong>
		</div>
		<div class="muni_level_badge">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:13px;height:13px" aria-hidden="true"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
			{current_nivel.nombre}
		</div>
	</div>
</div>

<div class="muni_list">
	{#each municipal_benefits as b (b.id)}
		<div class="muni_card" class:muni_card_available={user.points >= b.points_required}>
			<div class="muni_card_icon" style="background:color-mix(in srgb,{b.color} 12%,transparent);border-color:color-mix(in srgb,{b.color} 25%,transparent);color:{b.color}">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:18px;height:18px" aria-hidden="true"><path d={b.icon}/></svg>
			</div>
			<div class="muni_card_body">
				<div class="muni_card_top">
					<h3 class="muni_card_title">{b.title}</h3>
					<span class="muni_status muni_status_{b.status}">{benefit_status_labels[b.status]}</span>
				</div>
				<p class="muni_card_desc">{b.desc}</p>
				<div class="muni_card_foot">
					<span class="muni_pts_req" style="color:{b.color}">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:11px;height:11px;flex-shrink:0" aria-hidden="true"><path d="M12 3l2.7 5.5 6.1.9-4.4 4.3 1 6.1L12 17l-5.4 2.8 1-6.1L3.2 9.4l6.1-.9z"/></svg>
						{b.points_required.toLocaleString()} EP
					</span>
					{#if user.points >= b.points_required}
						<button class="muni_claim_btn" style="background:{b.color}">Reclamar</button>
					{:else}
						<span class="muni_pts_falt">Faltan {(b.points_required - user.points).toLocaleString()} EP</span>
					{/if}
				</div>
			</div>
		</div>
	{/each}
</div>

{/if}

<style>
	h1, h2, h3 { margin: 0; font-family: 'Sora', 'Plus Jakarta Sans', sans-serif; letter-spacing: -0.02em; }
	p { margin: 0; }
	sub { font-size: 0.6em; }
	a { text-decoration: none; }

	/* ── Top stats bar ── */
	.stats_bar { display: grid; grid-template-columns: repeat(5, 1fr); gap: 0.45rem; }
	.top_stat_card { border: 1.5px solid; border-radius: 1rem; padding: 0.6rem 0.45rem; display: flex; flex-direction: column; align-items: center; gap: 0.22rem; text-align: center; }
	.top_stat_icon { display: flex; }
	.top_stat_val { font-size: 0.78rem; font-weight: 900; font-family: 'Sora', sans-serif; line-height: 1; }
	.top_stat_lbl { font-size: 0.57rem; color: var(--ecochitas-muted); font-weight: 600; }

	/* ── Profile header ── */
	.profile_header { display: flex; align-items: center; gap: 1rem; padding: 1.1rem 1.25rem; }
	.avatar_wrap { flex-shrink: 0; position: relative; width: 54px; height: 54px; }
	.avatar_initials {
		position: absolute; inset: 0;
		background: linear-gradient(135deg, #0d9488, #16a34a);
		border-radius: 50%; display: flex; align-items: center; justify-content: center;
		font-size: 1rem; font-weight: 900; color: white; font-family: 'Sora', sans-serif;
		box-shadow: 0 4px 14px rgba(13,148,136,0.4);
	}
	.avatar_level_ring {
		position: absolute; inset: -3px;
		border-radius: 50%; border: 2px solid var(--ecochitas-leaf);
		opacity: 0.6; pointer-events: none;
	}
	.user_name { font-size: 1.05rem; font-weight: 900; color: var(--ecochitas-ink); line-height: 1.2; }
	.user_sub { font-size: 0.74rem; color: var(--ecochitas-muted); margin-top: 0.12rem; }
	.header_info { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 0.15rem; }
	.header_badges { display: flex; flex-wrap: wrap; gap: 0.35rem; margin-top: 0.3rem; }
	.user_level {
		display: inline-flex; align-items: center; gap: 0.3rem;
		font-size: 0.65rem; font-weight: 800;
		background: color-mix(in srgb, var(--ecochitas-leaf) 12%, transparent);
		color: var(--ecochitas-leaf); padding: 0.18rem 0.55rem;
		border-radius: 999px; font-family: 'Sora', sans-serif;
	}
	.streak_chip {
		display: inline-flex; align-items: center; gap: 0.3rem;
		font-size: 0.65rem; font-weight: 800;
		background: rgba(239,68,68,0.1); color: #dc2626;
		padding: 0.18rem 0.55rem; border-radius: 999px; font-family: 'Sora', sans-serif;
	}
	:root[data-theme='dark'] .streak_chip { background: rgba(239,68,68,0.15); color: #f87171; }

	/* ── Tab switcher ── */
	.tab_switcher { display: grid; grid-template-columns: repeat(5, 1fr); background: var(--ecochitas-surface); border: 1px solid var(--ecochitas-border); border-radius: 1.2rem; padding: 0.3rem; gap: 0.2rem; }
	.tab_btn { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 0.28rem; padding: 0.52rem 0.2rem; border-radius: 0.82rem; border: none; background: transparent; color: var(--ecochitas-muted); font-size: 0.6rem; font-weight: 700; font-family: 'Sora', 'Plus Jakarta Sans', sans-serif; cursor: pointer; transition: background 0.15s, color 0.15s; }
	.tab_active { background: var(--ecochitas-leaf); color: white; box-shadow: 0 3px 10px color-mix(in srgb, var(--ecochitas-leaf) 40%, transparent); }

	/* ── Inicio tab ── */
	.inicio_hero { background: linear-gradient(145deg, oklch(0.22 0.09 175), oklch(0.30 0.11 163)); border-radius: 1.5rem; padding: 1.5rem 1.4rem; display: flex; flex-direction: column; gap: 0.65rem; }
	.inicio_hero_eyebrow { display: inline-flex; align-items: center; gap: 0.4rem; font-size: 0.68rem; font-weight: 700; letter-spacing: 0.06em; text-transform: uppercase; color: oklch(1 0 0 / 0.6); background: oklch(1 0 0 / 0.1); border: 1px solid oklch(1 0 0 / 0.18); border-radius: 999px; padding: 0.22rem 0.7rem; width: fit-content; }
	.inicio_hero_title { font-size: 1.55rem; font-weight: 900; color: white; line-height: 1.12; }
	.inicio_name { color: var(--ecochitas-leaf); }
	.inicio_hero_sub { font-size: 0.82rem; color: oklch(1 0 0 / 0.72); line-height: 1.55; }
	.impact_grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 0.5rem; }
	.impact_stat_card { border: 1.5px solid var(--ecochitas-border); border-radius: 1rem; padding: 0.75rem 0.5rem; display: flex; flex-direction: column; align-items: center; gap: 0.22rem; text-align: center; background: var(--ecochitas-surface); }
	.impact_stat_icon { display: flex; }
	.impact_stat_val { font-size: 0.95rem; font-weight: 900; font-family: 'Sora', sans-serif; line-height: 1; }
	.impact_stat_lbl { font-size: 0.56rem; color: var(--ecochitas-muted); font-weight: 600; }
	.inicio_section { display: flex; flex-direction: column; gap: 0.55rem; }
	.inicio_section_head { display: flex; align-items: center; gap: 0.6rem; }
	.inicio_section_icon { flex-shrink: 0; width: 34px; height: 34px; border-radius: 0.75rem; border: 1.5px solid; display: flex; align-items: center; justify-content: center; }
	.inicio_section_title { font-size: 0.95rem; font-weight: 800; color: var(--ecochitas-ink); }
	.inicio_section_text { font-size: 0.82rem; color: var(--ecochitas-muted); line-height: 1.6; }
	.inicio_benefits_card { background: linear-gradient(145deg, oklch(0.97 0.02 163), oklch(0.94 0.04 170)); border: 1.5px solid #bbf7d0; border-radius: 1.5rem; padding: 1.25rem 1.35rem; }
	:root[data-theme='dark'] .inicio_benefits_card { background: rgba(34,197,94,0.07); border-color: rgba(34,197,94,0.2); }
	.inicio_benefits_title { font-size: 0.95rem; font-weight: 800; color: #166534; margin-bottom: 0.85rem; font-family: 'Sora', sans-serif; }
	:root[data-theme='dark'] .inicio_benefits_title { color: #4ade80; }
	.inicio_benefits_list { display: flex; flex-direction: column; gap: 0.65rem; }
	.inicio_benefit_item { display: flex; align-items: flex-start; gap: 0.55rem; font-size: 0.8rem; color: #166534; line-height: 1.5; }
	:root[data-theme='dark'] .inicio_benefit_item { color: #4ade80; }
	.inicio_benefit_dot { flex-shrink: 0; width: 26px; height: 26px; border-radius: 0.55rem; background: rgba(22,163,74,0.15); display: flex; align-items: center; justify-content: center; color: #16a34a; margin-top: 1px; }
	.edu_list { display: flex; flex-direction: column; gap: 0.55rem; }
	.edu_item { display: flex; align-items: flex-start; gap: 0.55rem; padding: 0.75rem 0.9rem; border-radius: 0.85rem; }
	.edu_num { flex-shrink: 0; font-size: 0.75rem; font-weight: 900; font-family: 'Sora', sans-serif; width: 18px; }
	.edu_text { font-size: 0.78rem; line-height: 1.5; font-weight: 500; }
	.mapa_preview_cta { display: flex; align-items: center; justify-content: space-between; gap: 0.85rem; padding: 1rem 1.25rem; border-radius: 1.25rem; background: var(--ecochitas-surface); border: 1.5px solid var(--ecochitas-border); transition: box-shadow 0.15s, border-color 0.15s; color: var(--ecochitas-ink); }
	.mapa_preview_cta:hover { border-color: var(--ecochitas-leaf); box-shadow: 0 4px 16px rgba(22,163,74,0.1); }
	.mapa_preview_left { display: flex; align-items: center; gap: 0.8rem; }
	.mapa_preview_icon { flex-shrink: 0; width: 42px; height: 42px; border-radius: 0.85rem; background: color-mix(in srgb, var(--ecochitas-leaf) 12%, transparent); border: 1.5px solid color-mix(in srgb, var(--ecochitas-leaf) 25%, transparent); display: flex; align-items: center; justify-content: center; color: var(--ecochitas-leaf); }
	.mapa_preview_title { font-size: 0.88rem; font-weight: 800; color: var(--ecochitas-ink); font-family: 'Sora', sans-serif; }
	.mapa_preview_sub { font-size: 0.72rem; color: var(--ecochitas-muted); margin-top: 0.1rem; }

	/* ── EcoWallet ── */
	.wallet_card { background: linear-gradient(145deg, oklch(0.30 0.11 163), oklch(0.22 0.09 175)); border-radius: 1.5rem; padding: 1.35rem 1.5rem 1.2rem; display: flex; flex-direction: column; gap: 0.25rem; }
	.wallet_header { display: flex; align-items: center; gap: 0.5rem; color: oklch(1 0 0 / 0.7); margin-bottom: 0.4rem; }
	.wallet_label { font-size: 0.78rem; font-weight: 800; font-family: 'Sora', sans-serif; letter-spacing: 0.05em; text-transform: uppercase; flex: 1; }
	.wallet_mvp_tag { font-size: 0.6rem; font-weight: 800; background: oklch(1 0 0 / 0.12); color: oklch(1 0 0 / 0.7); padding: 0.15rem 0.5rem; border-radius: 999px; border: 1px solid oklch(1 0 0 / 0.2); font-family: 'Sora', sans-serif; }
	.ecopoints_sublabel { font-size: 0.72rem; color: oklch(1 0 0 / 0.55); }
	.ecopoints_row { display: flex; align-items: baseline; gap: 0.4rem; margin: 0.1rem 0 0.2rem; }
	.points_num { font-size: 2.6rem; font-weight: 900; color: white; line-height: 1; letter-spacing: -0.04em; font-family: 'Sora', sans-serif; }
	.ecopoints_unit { font-size: 1rem; font-weight: 800; color: oklch(1 0 0 / 0.55); font-family: 'Sora', sans-serif; }
	.cash_disclaimer { display: flex; align-items: flex-start; gap: 0.4rem; background: oklch(1 0 0 / 0.07); border: 1px solid oklch(1 0 0 / 0.12); border-radius: 0.75rem; padding: 0.5rem 0.75rem; font-size: 0.71rem; color: oklch(1 0 0 / 0.6); line-height: 1.5; margin-bottom: 0.25rem; }
	.cash_disclaimer strong { color: oklch(1 0 0 / 0.85); font-weight: 700; }
	.redeem_btn { width: 100%; padding: 0.75rem 1rem; border-radius: 0.9rem; border: 1.5px solid oklch(1 0 0 / 0.3); background: oklch(1 0 0 / 0.08); color: white; font-size: 0.85rem; font-weight: 800; font-family: 'Sora', sans-serif; cursor: pointer; transition: background 0.15s; }
	.redeem_btn:hover { background: oklch(1 0 0 / 0.16); }

	/* ── Sistema de niveles ── */
	.nivel_panel { display: flex; flex-direction: column; gap: 0.5rem; }
	.nivel_head_row { display: flex; align-items: flex-start; gap: 0.7rem; }
	.nivel_badge_icon { flex-shrink: 0; width: 36px; height: 36px; border-radius: 0.75rem; background: color-mix(in srgb, #d97706 12%, transparent); border: 1.5px solid color-mix(in srgb, #d97706 25%, transparent); display: flex; align-items: center; justify-content: center; color: #d97706; }
	.nivel_sub { font-size: 0.74rem; color: var(--ecochitas-muted); margin-top: 0.12rem; }
	.nivel_progress_row { display: flex; align-items: center; justify-content: space-between; gap: 0.5rem; }
	.nivel_pts_text { font-size: 0.82rem; font-weight: 800; color: var(--ecochitas-ink); font-family: 'Sora', sans-serif; }
	.nivel_next_text { font-size: 0.72rem; color: var(--ecochitas-muted); text-align: right; }
	.nivel_track { height: 8px; background: var(--ecochitas-border); border-radius: 999px; overflow: hidden; }
	.nivel_fill { height: 100%; background: linear-gradient(90deg, #d97706, #f59e0b); border-radius: 999px; transition: width 0.5s ease; }
	.nivel_hint { font-size: 0.72rem; color: var(--ecochitas-muted); }
	.nivel_hint strong { color: var(--ecochitas-ink); }
	.nivel_table_scroll { overflow-x: auto; margin: 0.25rem 0; }
	.nivel_table { display: flex; gap: 0.3rem; width: max-content; }
	.nivel_row { display: flex; flex-direction: column; align-items: center; gap: 0.18rem; padding: 0.5rem 0.6rem; border-radius: 0.75rem; border: 1.5px solid var(--ecochitas-border); background: var(--ecochitas-surface); min-width: 72px; text-align: center; }
	.nivel_row_active { border-color: #d97706; background: #fffbeb; }
	:root[data-theme='dark'] .nivel_row_active { background: rgba(217,119,6,0.1); }
	.nivel_row_locked { opacity: 0.45; }
	.nivel_row_num { font-size: 0.62rem; font-weight: 800; color: var(--ecochitas-muted); font-family: 'Sora', sans-serif; }
	.nivel_row_active .nivel_row_num { color: #d97706; }
	.nivel_row_name { font-size: 0.58rem; font-weight: 700; color: var(--ecochitas-ink); line-height: 1.2; }
	.nivel_row_bonus { font-size: 0.65rem; font-weight: 900; color: #16a34a; font-family: 'Sora', sans-serif; }
	.nivel_inactivity_warn { display: flex; align-items: flex-start; gap: 0.45rem; padding: 0.6rem 0.8rem; border-radius: 0.75rem; background: #fffbeb; border: 1px solid #fde68a; font-size: 0.72rem; color: #92400e; line-height: 1.5; }
	:root[data-theme='dark'] .nivel_inactivity_warn { background: rgba(245,158,11,0.08); border-color: rgba(245,158,11,0.25); color: #fcd34d; }

	/* ── Section commons ── */
	.section_row { display: flex; align-items: center; gap: 0.45rem; margin-bottom: 0.85rem; }
	.section_title { font-size: 1rem; font-weight: 800; color: var(--ecochitas-ink); }
	.progress_row { display: flex; align-items: center; justify-content: space-between; gap: 0.5rem; margin-bottom: 0.55rem; }
	.progress_label { font-size: 0.82rem; color: var(--ecochitas-muted); font-weight: 600; }
	.progress_val { font-size: 0.82rem; font-weight: 800; font-family: 'Sora', sans-serif; white-space: nowrap; }
	.progress_track { height: 10px; background: var(--ecochitas-border); border-radius: 999px; overflow: hidden; }
	.progress_fill { height: 100%; background: linear-gradient(90deg, var(--ecochitas-leaf), #22c55e); border-radius: 999px; transition: width 0.5s ease; }
	.progress_hint { font-size: 0.74rem; color: var(--ecochitas-muted); margin-top: 0.45rem; }
	.stats_duo { display: grid; grid-template-columns: 1fr 1fr; gap: 0.6rem; margin-top: 0.85rem; }
	.stat_pill { border-radius: 1rem; padding: 0.85rem 0.75rem; display: flex; flex-direction: column; align-items: center; gap: 0.2rem; text-align: center; }
	.stat_pill_green { background: #f0fdf4; border: 1.5px solid #bbf7d0; }
	.stat_pill_blue { background: #eff6ff; border: 1.5px solid #bfdbfe; }
	:root[data-theme='dark'] .stat_pill_green { background: rgba(34,197,94,0.08); border-color: rgba(34,197,94,0.2); }
	:root[data-theme='dark'] .stat_pill_blue { background: rgba(59,130,246,0.08); border-color: rgba(59,130,246,0.2); }
	.stat_big { font-size: 1.5rem; font-weight: 900; color: var(--ecochitas-leaf); line-height: 1; font-family: 'Sora', sans-serif; }
	.stat_blue { color: var(--ecochitas-sky) !important; }
	.stat_small { font-size: 0.72rem; color: var(--ecochitas-muted); font-weight: 600; }

	/* ── Validador de reciclado ── */
	.validator_section { margin-top: 1.1rem; padding-top: 1rem; border-top: 1.5px solid var(--ecochitas-border); }
	.validator_head { display: flex; align-items: center; justify-content: space-between; gap: 0.5rem; margin-bottom: 0.35rem; }
	.validator_title_row { display: flex; align-items: center; gap: 0.4rem; }
	.validator_title { font-size: 0.88rem; font-weight: 800; color: var(--ecochitas-ink); font-family: 'Sora', sans-serif; }
	.validator_quota { font-size: 0.65rem; font-weight: 700; color: var(--ecochitas-muted); background: var(--ecochitas-border); padding: 0.15rem 0.5rem; border-radius: 999px; white-space: nowrap; }
	.validator_sub { font-size: 0.76rem; color: var(--ecochitas-muted); line-height: 1.5; margin-bottom: 0.65rem; }
	.validaciones_list { display: flex; flex-direction: column; gap: 0.45rem; margin-bottom: 0.65rem; }
	.validacion_item { display: flex; align-items: center; gap: 0.65rem; padding: 0.65rem 0.8rem; border-radius: 0.9rem; border: 1.5px solid var(--ecochitas-border); background: var(--ecochitas-surface); }
	.val_icon { flex-shrink: 0; width: 30px; height: 30px; border-radius: 50%; display: flex; align-items: center; justify-content: center; }
	.val_icon_validado { background: #dcfce7; color: #16a34a; }
	.val_icon_pendiente { background: #fef9c3; color: #d97706; }
	.val_icon_rechazado { background: #fef2f2; color: #ef4444; }
	:root[data-theme='dark'] .val_icon_validado { background: rgba(34,197,94,0.15); }
	:root[data-theme='dark'] .val_icon_pendiente { background: rgba(245,158,11,0.15); }
	:root[data-theme='dark'] .val_icon_rechazado { background: rgba(239,68,68,0.15); }
	.val_info { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 0.1rem; }
	.val_material { font-size: 0.8rem; font-weight: 800; color: var(--ecochitas-ink); font-family: 'Sora', sans-serif; }
	.val_meta { font-size: 0.68rem; color: var(--ecochitas-muted); }
	.val_right { display: flex; flex-direction: column; align-items: flex-end; gap: 0.1rem; }
	.val_status { font-size: 0.62rem; font-weight: 800; padding: 0.15rem 0.5rem; border-radius: 999px; white-space: nowrap; font-family: 'Sora', sans-serif; }
	.val_status_validado { background: #dcfce7; color: #166534; }
	:root[data-theme='dark'] .val_status_validado { background: rgba(34,197,94,0.15); color: #4ade80; }
	.val_status_pendiente { background: #fef9c3; color: #854d0e; }
	:root[data-theme='dark'] .val_status_pendiente { background: rgba(245,158,11,0.15); color: #fbbf24; }
	.val_status_rechazado { background: #fef2f2; color: #b91c1c; }
	:root[data-theme='dark'] .val_status_rechazado { background: rgba(239,68,68,0.15); color: #fca5a5; }
	.val_pts { font-size: 0.68rem; font-weight: 800; color: var(--ecochitas-leaf); font-family: 'Sora', sans-serif; }
	.validator_upload_btn {
		display: flex; align-items: center; justify-content: center; gap: 0.45rem;
		width: 100%; padding: 0.68rem 1rem; border-radius: 0.9rem;
		border: 1.5px dashed var(--ecochitas-leaf);
		background: color-mix(in srgb, var(--ecochitas-leaf) 6%, transparent);
		color: var(--ecochitas-leaf); font-size: 0.8rem; font-weight: 700;
		cursor: pointer; font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		transition: background 0.15s;
	}
	.validator_upload_btn:hover { background: color-mix(in srgb, var(--ecochitas-leaf) 12%, transparent); }
	.upload_input_hidden { display: none; }

	/* ── Impact ── */
	.impact_card { background: linear-gradient(145deg, #16a34a, #15803d); border-radius: 1.5rem; padding: 1.35rem 1.5rem; display: flex; flex-direction: column; gap: 0.35rem; }
	.impact_header { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.1rem; }
	.impact_emoji { font-size: 1.1rem; }
	.impact_title { font-size: 1rem; font-weight: 900; color: white; }
	.impact_desc { font-size: 0.8rem; color: oklch(1 0 0 / 0.7); }
	.impact_co2 { font-size: 2.3rem; font-weight: 900; color: white; line-height: 1.1; letter-spacing: -0.04em; font-family: 'Sora', sans-serif; margin: 0.1rem 0; }
	.impact_equiv { font-size: 0.78rem; color: oklch(1 0 0 / 0.65); }

	/* ── Insignias ── */
	.insignias_grid { display: grid; grid-template-columns: repeat(5, 1fr); gap: 0.4rem; margin-top: 0.65rem; }
	.insignia_card { border: 1.5px solid; border-radius: 0.85rem; padding: 0.6rem 0.35rem; display: flex; flex-direction: column; align-items: center; gap: 0.2rem; text-align: center; transition: transform 0.15s; }
	.insignia_card:hover { transform: translateY(-2px); }
	.insignia_locked { opacity: 0.5; }
	.insignia_emoji { font-size: 1.3rem; line-height: 1; }
	.insignia_emoji_locked { filter: grayscale(1); opacity: 0.5; }
	.insignia_title { font-size: 0.55rem; font-weight: 800; line-height: 1.2; font-family: 'Sora', sans-serif; }
	.insignia_date { font-size: 0.52rem; color: var(--ecochitas-muted); }

	/* ── Streak & actividad ── */
	.streak_row { display: grid; grid-template-columns: 1fr 1fr; gap: 0.6rem; }
	.streak_card { display: flex; align-items: center; gap: 0.65rem; padding: 0.85rem 0.9rem; border-radius: 1rem; background: rgba(239,68,68,0.06); border: 1.5px solid rgba(239,68,68,0.18); }
	:root[data-theme='dark'] .streak_card { background: rgba(239,68,68,0.1); }
	.streak_card_muted { background: var(--ecochitas-surface) !important; border-color: var(--ecochitas-border) !important; }
	.streak_val { font-size: 1rem; font-weight: 900; color: #dc2626; font-family: 'Sora', sans-serif; display: block; }
	:root[data-theme='dark'] .streak_val { color: #f87171; }
	.streak_val_muted { color: var(--ecochitas-muted) !important; }
	.streak_lbl { font-size: 0.65rem; color: var(--ecochitas-muted); font-weight: 600; }
	.activity_list { display: flex; flex-direction: column; gap: 0.45rem; }
	.activity_item { display: flex; align-items: center; gap: 0.65rem; padding: 0.65rem 0.75rem; border-radius: 0.9rem; border: 1.5px solid var(--ecochitas-border); background: var(--ecochitas-surface); }
	.activity_icon { flex-shrink: 0; width: 30px; height: 30px; border-radius: 50%; border: 1.5px solid; display: flex; align-items: center; justify-content: center; }
	.activity_info { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 0.08rem; }
	.activity_desc { font-size: 0.76rem; font-weight: 600; color: var(--ecochitas-ink); line-height: 1.2; }
	.activity_date { font-size: 0.62rem; color: var(--ecochitas-muted); }
	.activity_pts { font-size: 0.72rem; font-weight: 900; color: var(--ecochitas-leaf); font-family: 'Sora', sans-serif; white-space: nowrap; }
	.activity_pts_neg { color: #dc2626 !important; }
	:root[data-theme='dark'] .activity_pts_neg { color: #f87171 !important; }
	.activity_pts_zero { color: var(--ecochitas-muted) !important; }

	/* ── Condominio ── */
	.building_header_card { display: flex; align-items: center; gap: 0.9rem; padding: 1rem 1.25rem; }
	.building_icon_wrap { flex-shrink: 0; width: 46px; height: 46px; border-radius: 0.85rem; background: color-mix(in srgb, var(--ecochitas-leaf) 12%, transparent); border: 1.5px solid color-mix(in srgb, var(--ecochitas-leaf) 25%, transparent); display: flex; align-items: center; justify-content: center; color: var(--ecochitas-leaf); }
	.building_name { font-size: 1rem; font-weight: 900; color: var(--ecochitas-ink); }
	.building_sub { font-size: 0.76rem; color: var(--ecochitas-muted); margin-top: 0.12rem; }
	.reward_card { background: linear-gradient(145deg, #fef9c3, #fef3c7); border: 1.5px solid #fde68a; border-radius: 1.5rem; padding: 1.1rem 1.25rem; display: flex; flex-direction: column; gap: 0.45rem; }
	:root[data-theme='dark'] .reward_card { background: rgba(251,191,36,0.08); border-color: rgba(251,191,36,0.22); }
	.reward_header { display: flex; align-items: center; gap: 0.55rem; }
	.reward_icon { width: 32px; height: 32px; border-radius: 0.6rem; background: #fde68a; border: 1px solid #fcd34d; display: flex; align-items: center; justify-content: center; color: #92400e; flex-shrink: 0; }
	:root[data-theme='dark'] .reward_icon { background: rgba(251,191,36,0.18); border-color: rgba(251,191,36,0.32); color: #fcd34d; }
	.reward_title { font-size: 0.84rem; font-weight: 800; color: #92400e; font-family: 'Sora', sans-serif; }
	:root[data-theme='dark'] .reward_title { color: #fcd34d; }
	.reward_desc { font-size: 0.78rem; color: #78350f; line-height: 1.5; }
	:root[data-theme='dark'] .reward_desc { color: #fde68a; }
	.reward_rank { display: inline-flex; align-items: center; gap: 0.35rem; font-size: 0.76rem; font-weight: 800; color: #b45309; font-family: 'Sora', sans-serif; }
	:root[data-theme='dark'] .reward_rank { color: #fbbf24; }
	.rank_list { display: flex; flex-direction: column; gap: 0.5rem; }
	.rank_row { display: flex; align-items: center; gap: 0.75rem; padding: 0.72rem 0.85rem; border-radius: 1rem; border: 1.5px solid var(--ecochitas-border); background: var(--ecochitas-surface); }
	.rank_row_gold { background: #fefce8; border-color: #fde68a; }
	:root[data-theme='dark'] .rank_row_gold { background: rgba(251,191,36,0.07); border-color: rgba(251,191,36,0.22); }
	.rank_badge { flex-shrink: 0; width: 28px; height: 28px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 0.75rem; font-weight: 900; font-family: 'Sora', sans-serif; background: var(--ecochitas-border); color: var(--ecochitas-muted); }
	.rank_badge_gold { background: #f59e0b; color: white; }
	.rank_badge_silver { background: #94a3b8; color: white; }
	.rank_badge_bronze { background: #c2855a; color: white; }
	.rank_unit { flex: 1; font-size: 0.83rem; font-weight: 700; color: var(--ecochitas-ink); }
	.rank_right { display: flex; flex-direction: column; align-items: flex-end; gap: 0.05rem; }
	.rank_pts { font-size: 0.95rem; font-weight: 900; color: var(--ecochitas-sky); font-family: 'Sora', sans-serif; }
	.rank_pts_gold { color: #d97706; }
	.rank_pts_lbl { font-size: 0.6rem; color: var(--ecochitas-muted); }
	.benefits_card { background: linear-gradient(145deg, #7c3aed, #6d28d9); border-radius: 1.5rem; padding: 1.25rem 1.35rem; display: flex; flex-direction: column; gap: 0.8rem; }
	.benefits_header { display: flex; align-items: center; gap: 0.55rem; }
	.benefits_emoji { font-size: 1.1rem; }
	.benefits_title { font-size: 0.95rem; font-weight: 900; color: white; }
	.benefits_list { margin: 0; padding: 0; list-style: none; display: flex; flex-direction: column; gap: 0.5rem; }
	.benefit_item { display: flex; align-items: flex-start; gap: 0.45rem; font-size: 0.8rem; color: oklch(1 0 0 / 0.85); line-height: 1.5; }

	/* ── Tienda ── */
	.store_hero { display: flex; align-items: flex-start; gap: 0.85rem; padding: 1.25rem 1.4rem; border-radius: 1.5rem; background: linear-gradient(145deg, oklch(0.95 0.03 290), oklch(0.92 0.04 280)); border: 1.5px solid oklch(0.87 0.06 285); flex-wrap: wrap; }
	:root[data-theme='dark'] .store_hero { background: rgba(124,58,237,0.1); border-color: rgba(124,58,237,0.22); }
	.store_hero_icon { font-size: 1.8rem; line-height: 1; }
	.store_hero_title { font-size: 1.05rem; font-weight: 900; color: #4c1d95; }
	:root[data-theme='dark'] .store_hero_title { color: #c4b5fd; }
	.store_hero_sub { font-size: 0.78rem; color: #6d28d9; margin-top: 0.12rem; }
	:root[data-theme='dark'] .store_hero_sub { color: #a78bfa; }
	.store_pts_chip { margin-left: auto; display: inline-flex; align-items: center; gap: 0.35rem; font-size: 0.72rem; font-weight: 800; font-family: 'Sora', sans-serif; background: #7c3aed; color: white; padding: 0.3rem 0.7rem; border-radius: 999px; }
	.filter_scroll { overflow-x: auto; padding-bottom: 0.2rem; }
	.filter_pills { display: flex; gap: 0.4rem; width: max-content; }
	.filter_pill { display: inline-flex; align-items: center; gap: 0.4rem; padding: 0.4rem 0.8rem; border-radius: 999px; border: 1.5px solid var(--ecochitas-border); background: var(--ecochitas-surface); color: var(--ecochitas-muted); font-size: 0.74rem; font-weight: 700; cursor: pointer; white-space: nowrap; transition: border-color 0.15s, background 0.15s, color 0.15s; font-family: 'Sora', 'Plus Jakarta Sans', sans-serif; }
	.filter_pill:hover { border-color: var(--ecochitas-leaf); color: var(--ecochitas-leaf); }
	.filter_pill_active { border-color: var(--ecochitas-leaf); background: color-mix(in srgb, var(--ecochitas-leaf) 12%, transparent); color: var(--ecochitas-leaf); }
	.store_grid { display: grid; grid-template-columns: 1fr; gap: 0.7rem; }
	.store_card { border: 1.5px solid var(--ecochitas-border); border-radius: 1.25rem; padding: 1rem 1.1rem; background: var(--ecochitas-surface); display: flex; flex-direction: column; gap: 0.42rem; transition: box-shadow 0.15s, transform 0.12s; }
	.store_card:hover { box-shadow: 0 4px 18px rgba(0,0,0,0.07); transform: translateY(-2px); }
	.store_logo_area { display: flex; align-items: center; justify-content: space-between; gap: 0.5rem; margin-bottom: 0.1rem; }
	.store_logo_placeholder { width: 40px; height: 40px; border-radius: 0.75rem; background: var(--ecochitas-border); display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
	.store_card_badges { display: flex; align-items: center; gap: 0.35rem; flex-wrap: wrap; justify-content: flex-end; }
	.store_cat_tag { font-size: 0.63rem; font-weight: 800; padding: 0.18rem 0.5rem; border-radius: 999px; font-family: 'Sora', sans-serif; white-space: nowrap; }
	.store_type_badge { font-size: 0.68rem; font-weight: 900; color: #16a34a; background: #dcfce7; padding: 0.18rem 0.5rem; border-radius: 999px; font-family: 'Sora', sans-serif; white-space: nowrap; }
	.store_type_product { color: #7c3aed; background: #ede9fe; }
	.store_item_name { font-size: 0.9rem; font-weight: 800; color: var(--ecochitas-ink); line-height: 1.2; }
	.store_partner { font-size: 0.7rem; color: var(--ecochitas-muted); font-weight: 600; }
	.store_desc { font-size: 0.74rem; color: var(--ecochitas-muted); line-height: 1.45; }
	.store_card_foot { display: flex; align-items: center; justify-content: space-between; gap: 0.5rem; margin-top: 0.2rem; padding-top: 0.65rem; border-top: 1px solid var(--ecochitas-border); }
	.store_pts_needed { display: inline-flex; align-items: center; gap: 0.28rem; font-size: 0.8rem; font-weight: 800; color: #d97706; font-family: 'Sora', sans-serif; }
	.store_redeem_btn { padding: 0.42rem 0.85rem; border-radius: 0.7rem; border: none; background: var(--ecochitas-leaf); color: white; font-size: 0.74rem; font-weight: 800; cursor: pointer; font-family: 'Sora', sans-serif; transition: background 0.15s; }
	.store_redeem_btn:disabled { background: var(--ecochitas-border); color: var(--ecochitas-muted); cursor: not-allowed; }

	/* ── Municipal ── */
	.muni_hero { display: flex; align-items: flex-start; gap: 0.85rem; padding: 1.2rem 1.35rem; border-radius: 1.5rem; background: linear-gradient(145deg, #1e3a5f, #1e40af); }
	.muni_hero_icon { font-size: 1.8rem; line-height: 1; }
	.muni_hero_title { font-size: 1.05rem; font-weight: 900; color: white; }
	.muni_hero_sub { font-size: 0.78rem; color: oklch(1 0 0 / 0.7); margin-top: 0.12rem; line-height: 1.5; }
	.muni_mvp_notice { display: flex; align-items: flex-start; gap: 0.5rem; padding: 0.7rem 0.9rem; border-radius: 0.9rem; background: #fffbeb; border: 1px solid #fde68a; font-size: 0.74rem; color: #92400e; line-height: 1.5; }
	:root[data-theme='dark'] .muni_mvp_notice { background: rgba(245,158,11,0.08); border-color: rgba(245,158,11,0.25); color: #fcd34d; }
	.muni_mvp_notice strong { font-weight: 800; }
	.muni_progress_card { padding: 1rem 1.1rem; }
	.muni_prog_row { display: flex; align-items: center; justify-content: space-between; gap: 0.75rem; }
	.muni_prog_label { font-size: 0.76rem; color: var(--ecochitas-muted); }
	.muni_prog_pts { font-size: 1.4rem; font-weight: 900; color: var(--ecochitas-ink); font-family: 'Sora', sans-serif; }
	.muni_level_badge { display: inline-flex; align-items: center; gap: 0.3rem; font-size: 0.68rem; font-weight: 800; background: color-mix(in srgb, var(--ecochitas-leaf) 12%, transparent); color: var(--ecochitas-leaf); padding: 0.28rem 0.7rem; border-radius: 999px; font-family: 'Sora', sans-serif; }
	.muni_list { display: flex; flex-direction: column; gap: 0.75rem; }
	.muni_card { display: flex; align-items: flex-start; gap: 0.85rem; padding: 1rem 1.1rem; border: 1.5px solid var(--ecochitas-border); border-radius: 1.2rem; background: var(--ecochitas-surface); transition: box-shadow 0.15s; }
	.muni_card_available { border-color: #bbf7d0; background: #f0fdf4; }
	:root[data-theme='dark'] .muni_card_available { border-color: rgba(34,197,94,0.22); background: rgba(34,197,94,0.05); }
	.muni_card_icon { flex-shrink: 0; width: 44px; height: 44px; border-radius: 0.85rem; border: 1.5px solid; display: flex; align-items: center; justify-content: center; }
	.muni_card_body { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 0.3rem; }
	.muni_card_top { display: flex; align-items: flex-start; justify-content: space-between; gap: 0.5rem; }
	.muni_card_title { font-size: 0.88rem; font-weight: 800; color: var(--ecochitas-ink); line-height: 1.2; font-family: 'Sora', sans-serif; }
	.muni_card_desc { font-size: 0.74rem; color: var(--ecochitas-muted); line-height: 1.45; }
	.muni_card_foot { display: flex; align-items: center; justify-content: space-between; gap: 0.5rem; flex-wrap: wrap; margin-top: 0.15rem; padding-top: 0.55rem; border-top: 1px solid var(--ecochitas-border); }
	.muni_pts_req { display: inline-flex; align-items: center; gap: 0.28rem; font-size: 0.72rem; font-weight: 700; font-family: 'Sora', sans-serif; }
	.muni_status { font-size: 0.6rem; font-weight: 800; padding: 0.18rem 0.5rem; border-radius: 999px; white-space: nowrap; font-family: 'Sora', sans-serif; }
	.muni_status_disponible { background: #dcfce7; color: #166534; }
	:root[data-theme='dark'] .muni_status_disponible { background: rgba(34,197,94,0.15); color: #4ade80; }
	.muni_status_por_alcanzar { background: #fef9c3; color: #854d0e; }
	:root[data-theme='dark'] .muni_status_por_alcanzar { background: rgba(234,179,8,0.15); color: #fbbf24; }
	.muni_status_obtenido { background: #ede9fe; color: #4c1d95; }
	.muni_claim_btn { padding: 0.38rem 0.8rem; border-radius: 0.65rem; border: none; color: white; font-size: 0.72rem; font-weight: 800; cursor: pointer; font-family: 'Sora', sans-serif; }
	.muni_pts_falt { font-size: 0.68rem; font-weight: 700; color: var(--ecochitas-muted); }

	/* ── Responsive ── */
	@media (min-width: 380px) {
		.store_grid { grid-template-columns: repeat(2, 1fr); }
	}
	@media (max-width: 340px) {
		.stats_bar { grid-template-columns: repeat(3, 1fr); }
		.insignias_grid { grid-template-columns: repeat(4, 1fr); }
	}
</style>
