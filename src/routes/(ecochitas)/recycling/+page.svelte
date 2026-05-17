<script lang="ts">
	import EcochitasIcon from '$lib/components/ecochitas/EcochitasIcon.svelte';
	import RecyclingIcon, {
		type RecyclingIconName
	} from '$lib/components/ecochitas/RecyclingIcon.svelte';

	type Recyclable_item = {
		icon: RecyclingIconName;
		name: string;
		desc: string;
	};

	type Recycling_category = {
		id: string;
		name: string;
		icon: RecyclingIconName;
		accent: string;
		accent_light: string;
		accent_border: string;
		price: string;
		price_numeric: number;
		unit: string;
		tip: string;
		items: Recyclable_item[];
	};

	const recycling_categories: Recycling_category[] = [
		{
			id: 'papel',
			name: 'Papel & Cartón',
			icon: 'cat_papel',
			accent: '#92400e',
			accent_light: '#fef3c7',
			accent_border: '#fcd34d',
			price: '3.50 Bs/kg',
			price_numeric: 3.5,
			unit: 'por kilogramo',
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
			id: 'plasticos',
			name: 'Plásticos',
			icon: 'cat_plasticos',
			accent: '#075985',
			accent_light: '#e0f2fe',
			accent_border: '#7dd3fc',
			price: '4.00 Bs/kg',
			price_numeric: 4.0,
			unit: 'por kilogramo',
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
			id: 'vidrio',
			name: 'Vidrio',
			icon: 'cat_vidrio',
			accent: '#065f46',
			accent_light: '#d1fae5',
			accent_border: '#6ee7b7',
			price: '1.60 Bs/kg',
			price_numeric: 1.6,
			unit: 'por kilogramo',
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
			id: 'metales',
			name: 'Metales',
			icon: 'cat_metales',
			accent: '#4c1d95',
			accent_light: '#ede9fe',
			accent_border: '#c4b5fd',
			price: '6.20 Bs/kg',
			price_numeric: 6.2,
			unit: 'por kilogramo',
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
			id: 'ewaste',
			name: 'E-Waste',
			icon: 'cat_ewaste',
			accent: '#991b1b',
			accent_light: '#fee2e2',
			accent_border: '#fca5a5',
			price: 'Punto especial',
			price_numeric: 0,
			unit: 'acopio GAMC',
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

	const how_it_works = [
		{
			icon: 'step_separar' as RecyclingIconName,
			title: 'Separa',
			desc: 'Clasifica los materiales por tipo en casa antes de la recolección'
		},
		{
			icon: 'step_limpiar' as RecyclingIconName,
			title: 'Limpia',
			desc: 'Enjuaga envases y retira tapas, etiquetas y restos de comida'
		},
		{
			icon: 'step_entregar' as RecyclingIconName,
			title: 'Entrega',
			desc: 'Lleva los materiales al recolector asignado el viernes por la mañana'
		},
		{
			icon: 'step_ganar' as RecyclingIconName,
			title: 'Gana puntos',
			desc: 'Acumula bonificaciones individuales y por zona con cada aporte válido'
		}
	];

	const max_price = Math.max(...recycling_categories.map((c) => c.price_numeric));

	let selected_id = $state(recycling_categories[0].id);

	const selected = $derived(
		recycling_categories.find((c) => c.id === selected_id) ?? recycling_categories[0]
	);

	const sorted_by_price = $derived(
		[...recycling_categories].sort((a, b) => b.price_numeric - a.price_numeric)
	);

	// ── Reporting form state ──
	type ViolationType = {
		id: string;
		label: string;
		icon_path: string;
	};

	const violation_types: ViolationType[] = [
		{
			id: 'contaminacion',
			label: 'Contaminación en bolsa',
			icon_path:
				'M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0zM12 9v4M12 17h.01'
		},
		{
			id: 'mezcla',
			label: 'Mezcla con orgánicos',
			icon_path: 'M3 3h18v18H3zM12 8v8M8 12h8'
		},
		{
			id: 'material_sucio',
			label: 'Material sucio / húmedo',
			icon_path:
				'M12 22a7 7 0 0 0 7-7c0-2-1-3.9-3-5.5s-3.5-4-4-6.5c-.5 2.5-2 4.9-4 6.5C6 11.1 5 13 5 15a7 7 0 0 0 7 7z'
		},
		{
			id: 'fuera_horario',
			label: 'Depósito fuera de horario',
			icon_path: 'M12 22c5.5 0 10-4.5 10-10S17.5 2 12 2 2 6.5 2 12s4.5 10 10 10zM12 6v6l4 2'
		},
		{
			id: 'ewaste',
			label: 'E-waste en recolector común',
			icon_path:
				'M9 3H5a2 2 0 0 0-2 2v4m6-6h10a2 2 0 0 1 2 2v4M9 3v18m0 0h10a2 2 0 0 0 2-2V9M9 21H5a2 2 0 0 1-2-2V9m0 0h18'
		},
		{ id: 'otro', label: 'Otro', icon_path: 'M12 22c5.5 0 10-4.5 10-10S17.5 2 12 2 2 6.5 2 12s4.5 10 10 10zM9 9h.01M15 9h.01M8 13s1.5 2 4 2 4-2 4-2' }
	];

	let selected_violation = $state('');
	let location_detected = $state(false);
	let report_sent = $state(false);

	function detect_location() {
		location_detected = true;
	}

	function send_report() {
		if (!selected_violation) return;
		report_sent = true;
		setTimeout(() => {
			report_sent = false;
			selected_violation = '';
			location_detected = false;
		}, 3000);
	}

	// ── Contaminadores data ──
	type Contaminador = {
		id: number;
		location: string;
		status: string;
		status_type: 'pending' | 'warning' | 'investigating';
		reports: number;
		last_report: string;
	};

	const contaminadores: Contaminador[] = [
		{
			id: 1,
			location: 'Av. América & Pando',
			status: 'Multa Municipal Pendiente',
			status_type: 'pending',
			reports: 7,
			last_report: '14 May 2026'
		},
		{
			id: 2,
			location: 'Plaza Colón',
			status: 'Aviso Emitido',
			status_type: 'warning',
			reports: 5,
			last_report: '12 May 2026'
		},
		{
			id: 3,
			location: 'Av. Heroínas',
			status: 'En Investigación',
			status_type: 'investigating',
			reports: 4,
			last_report: '10 May 2026'
		},
		{
			id: 4,
			location: 'Parque Tunari',
			status: 'Multa Municipal Pendiente',
			status_type: 'pending',
			reports: 3,
			last_report: '09 May 2026'
		}
	];
</script>

<!-- ── Hero ── -->
<div class="hero">
	<div class="hero_eyebrow">
		<EcochitasIcon name="recycling" size={14} />
		<span>EcoHub Educativo · Cochabamba</span>
	</div>
	<h1 class="hero_title">Recicla mejor,<br />gana más puntos</h1>
	<p class="hero_subtitle">
		Aprende qué materiales acepta Ecochita, cómo prepararlos correctamente y cuánto valen en el
		mercado.
	</p>
	<div class="hero_stats">
		<div class="hero_stat">
			<span class="stat_num">5</span>
			<span class="stat_lbl">Categorías</span>
		</div>
		<div class="stat_div"></div>
		<div class="hero_stat">
			<span class="stat_num">6.20 Bs</span>
			<span class="stat_lbl">Precio máx/kg</span>
		</div>
		<div class="stat_div"></div>
		<div class="hero_stat">
			<span class="stat_num">Viernes</span>
			<span class="stat_lbl">Turno de recolección</span>
		</div>
	</div>
</div>

<!-- ── Cómo funciona ── -->
<section class="panel">
	<h2 class="section_title">¿Cómo funciona?</h2>
	<p class="section_sub">Cuatro pasos para reciclar correctamente y ganar tu bonificación.</p>
	<div class="steps_row">
		{#each how_it_works as step, i (step.title)}
			<div class="step_item">
				<div class="step_icon_wrap">
					<RecyclingIcon name={step.icon} size={26} />
					<span class="step_num">{i + 1}</span>
				</div>
				<strong class="step_title">{step.title}</strong>
				<p class="step_desc">{step.desc}</p>
			</div>
			{#if i < how_it_works.length - 1}
				<div class="step_connector" aria-hidden="true">
					<div class="connector_line"></div>
				</div>
			{/if}
		{/each}
	</div>
</section>

<!-- ── Categorías ── -->
<section class="panel">
	<h2 class="section_title">¿Qué puedo reciclar?</h2>
	<p class="section_sub">Selecciona una categoría para ver sus artículos aceptados.</p>

	<!-- Category picker -->
	<div class="cat_picker" role="tablist">
		{#each recycling_categories as cat (cat.id)}
			<button
				role="tab"
				aria-selected={selected_id === cat.id}
				class="cat_btn"
				class:cat_btn_active={selected_id === cat.id}
				style="--accent:{cat.accent};--accent-light:{cat.accent_light};--accent-border:{cat.accent_border}"
				onclick={() => (selected_id = cat.id)}
			>
				<span class="cat_btn_icon">
					<RecyclingIcon name={cat.icon} size={22} />
				</span>
				<span class="cat_btn_name">{cat.name}</span>
				<span class="cat_btn_price">{cat.price}</span>
			</button>
		{/each}
	</div>

	<!-- Selected category detail -->
	<div
		class="cat_detail"
		style="--accent:{selected.accent};--accent-light:{selected.accent_light};--accent-border:{selected.accent_border}"
	>
		<!-- Header -->
		<div class="cat_detail_head">
			<div class="cat_detail_icon">
				<RecyclingIcon name={selected.icon} size={28} />
			</div>
			<div class="cat_detail_meta">
				<h3 class="cat_detail_name">{selected.name}</h3>
				<p class="cat_detail_tip">
					<EcochitasIcon name="alert" size={12} />
					{selected.tip}
				</p>
			</div>
			<div class="cat_price_box">
				<span class="price_box_lbl">Valor de mercado</span>
				<strong class="price_box_val">{selected.price}</strong>
				<span class="price_box_unit">{selected.unit}</span>
			</div>
		</div>

		<!-- Items grid -->
		<div class="items_grid">
			{#each selected.items as item (item.name)}
				<div class="item_card">
					<div class="item_icon">
						<RecyclingIcon name={item.icon} size={28} />
					</div>
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
	<p class="section_sub">
		Valores de referencia para el turno de reciclaje de viernes en Cochabamba.
	</p>

	<div class="ranking_list">
		{#each sorted_by_price as cat, i (cat.id)}
			<div
				class="rank_row"
				style="--accent:{cat.accent};--accent-light:{cat.accent_light};--accent-border:{cat.accent_border}"
			>
				<span class="rank_pos">#{i + 1}</span>
				<span class="rank_icon"><RecyclingIcon name={cat.icon} size={20} /></span>
				<span class="rank_name">{cat.name}</span>
				<div class="rank_bar_track">
					{#if cat.price_numeric > 0}
						<div
							class="rank_bar_fill"
							style="width:{(cat.price_numeric / max_price) * 100}%"
						></div>
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
				<span class="rules_head_icon do_icon">
					<EcochitasIcon name="upload" size={14} />
				</span>
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
				<span class="rules_head_icon dont_icon">
					<EcochitasIcon name="alert" size={14} />
				</span>
				<span class="rules_head_label">Qué debes evitar</span>
			</div>
			<ul class="rules_list">
				<li>Mezclar residuos orgánicos con materiales reciclables.</li>
				<li>Entregar envases mojados o con restos de comida adheridos.</li>
				<li>Superar el 20 % de contaminación — congela tu bonificación semanal.</li>
				<li>Llevar e-waste al recolector común; usa solo los puntos GAMC.</li>
			</ul>
		</div>
	</div>
</section>

<!-- ── Reportar Práctica Incorrecta ── -->
<section class="panel">
	<div class="section_head_row">
		<svg
			viewBox="0 0 24 24"
			fill="none"
			stroke="#ef4444"
			stroke-linecap="round"
			stroke-linejoin="round"
			stroke-width="1.9"
			style="width:18px;height:18px;flex-shrink:0"
			aria-hidden="true"
		>
			<path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z" />
			<line x1="12" y1="9" x2="12" y2="13" />
			<line x1="12" y1="17" x2="12.01" y2="17" />
		</svg>
		<h2 class="section_title">Reportar Práctica Incorrecta</h2>
	</div>
	<p class="section_sub">Denuncia anónimamente malas prácticas de reciclaje en tu zona.</p>

	<!-- GPS Button -->
	<button
		class="gps_btn"
		class:gps_btn_active={location_detected}
		onclick={detect_location}
		aria-pressed={location_detected}
	>
		{#if location_detected}
			<svg
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2.2"
				style="width:16px;height:16px;flex-shrink:0"
				aria-hidden="true"
			>
				<polyline points="20 6 9 17 4 12" />
			</svg>
			Ubicación Detectada
		{:else}
			<svg
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="1.9"
				style="width:16px;height:16px;flex-shrink:0"
				aria-hidden="true"
			>
				<path d="M12 21s-8-7.5-8-12a8 8 0 0 1 16 0c0 4.5-8 12-8 12z" />
				<circle cx="12" cy="9" r="2.5" />
			</svg>
			Detectar Mi Ubicación (GPS)
		{/if}
	</button>

	<!-- Violation type picker -->
	<p class="field_label">Tipo de Infracción <span class="field_required">*</span></p>
	<div class="violation_grid">
		{#each violation_types as vtype (vtype.id)}
			<button
				class="vtype_btn"
				class:vtype_btn_active={selected_violation === vtype.id}
				onclick={() => (selected_violation = vtype.id)}
				aria-pressed={selected_violation === vtype.id}
			>
				<svg
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="1.75"
					style="width:14px;height:14px;flex-shrink:0"
					aria-hidden="true"
				>
					<path d={vtype.icon_path} />
				</svg>
				{vtype.label}
			</button>
		{/each}
	</div>

	<!-- Submit button -->
	<button
		class="submit_btn"
		class:submit_btn_sent={report_sent}
		disabled={!selected_violation || report_sent}
		onclick={send_report}
	>
		{#if report_sent}
			<svg
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2.2"
				style="width:16px;height:16px;flex-shrink:0"
				aria-hidden="true"
			>
				<polyline points="20 6 9 17 4 12" />
			</svg>
			Reporte Enviado
		{:else}
			<svg
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="1.9"
				style="width:16px;height:16px;flex-shrink:0"
				aria-hidden="true"
			>
				<line x1="22" y1="2" x2="11" y2="13" />
				<polygon points="22 2 15 22 11 13 2 9 22 2" fill="currentColor" stroke="none" />
			</svg>
			Enviar Reporte Anónimo
		{/if}
	</button>

	<!-- Privacy notice -->
	<div class="privacy_notice">
		<svg
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-linecap="round"
			stroke-linejoin="round"
			stroke-width="1.9"
			style="width:13px;height:13px;flex-shrink:0;margin-top:1px"
			aria-hidden="true"
		>
			<rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
			<path d="M7 11V7a5 5 0 0 1 10 0v4" />
		</svg>
		<span>
			<strong>Tu privacidad está protegida:</strong> Los reportes son completamente anónimos. La
			información se usa solo para mejorar el programa de reciclaje.
		</span>
	</div>
</section>

<!-- ── Contaminadores de la Semana ── -->
<section class="panel contam_panel">
	<div class="section_head_row">
		<svg
			viewBox="0 0 24 24"
			fill="none"
			stroke="#ef4444"
			stroke-linecap="round"
			stroke-linejoin="round"
			stroke-width="1.9"
			style="width:18px;height:18px;flex-shrink:0"
			aria-hidden="true"
		>
			<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
			<circle cx="9" cy="7" r="4" />
			<line x1="17" y1="8" x2="23" y2="14" />
			<line x1="23" y1="8" x2="17" y2="14" />
		</svg>
		<h2 class="section_title contam_title">Contaminadores de la Semana</h2>
	</div>
	<p class="section_sub">Infractores recurrentes reportados 3+ veces en el mismo sector.</p>

	<div class="contam_list">
		{#each contaminadores as c (c.id)}
			<div class="contam_card">
				<!-- Avatar -->
				<div class="contam_avatar">
					<svg
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="1.75"
						style="width:20px;height:20px"
						aria-hidden="true"
					>
						<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
						<circle cx="9" cy="7" r="4" />
						<line x1="17" y1="8" x2="23" y2="14" />
						<line x1="23" y1="8" x2="17" y2="14" />
					</svg>
				</div>

				<!-- Info -->
				<div class="contam_info">
					<div class="contam_name_row">
						<strong class="contam_name">Infractor #{c.id}</strong>
						<span class="contam_badge">{c.reports} Reportes</span>
					</div>
					<div class="contam_location">
						<svg
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="1.9"
							style="width:11px;height:11px;flex-shrink:0"
							aria-hidden="true"
						>
							<path d="M12 21s-8-7.5-8-12a8 8 0 0 1 16 0c0 4.5-8 12-8 12z" />
							<circle cx="12" cy="9" r="2" />
						</svg>
						{c.location}
					</div>
					<div class="contam_meta_row">
						<span class="contam_status contam_status_{c.status_type}">
							{#if c.status_type === 'pending'}
								<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" style="width:10px;height:10px;flex-shrink:0" aria-hidden="true"><circle cx="12" cy="12" r="9"/><path d="M12 6v6l3 3"/></svg>
							{:else if c.status_type === 'warning'}
								<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" style="width:10px;height:10px;flex-shrink:0" aria-hidden="true"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
							{:else}
								<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" style="width:10px;height:10px;flex-shrink:0" aria-hidden="true"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
							{/if}
							{c.status}
						</span>
						<span class="contam_date">
							<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:10px;height:10px;flex-shrink:0" aria-hidden="true"><rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
							{c.last_report}
						</span>
					</div>
				</div>
			</div>
		{/each}
	</div>
</section>

<style>
	/* ── Reset & base ── */
	h1,
	h2,
	h3 {
		margin: 0;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		letter-spacing: -0.02em;
	}

	p {
		margin: 0;
	}

	/* ── Hero ── */
	.hero {
		background: linear-gradient(145deg, oklch(0.33 0.12 163), oklch(0.24 0.09 178));
		border-radius: 1.5rem;
		padding: 1.75rem 1.5rem 1.4rem;
		display: flex;
		flex-direction: column;
		gap: 0.9rem;
	}

	.hero_eyebrow {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		font-size: 0.72rem;
		font-weight: 700;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		color: oklch(1 0 0 / 0.6);
		background: oklch(1 0 0 / 0.1);
		border: 1px solid oklch(1 0 0 / 0.18);
		border-radius: 999px;
		padding: 0.25rem 0.7rem;
		width: fit-content;
	}

	.hero_title {
		font-size: 1.7rem;
		font-weight: 900;
		line-height: 1.12;
		color: #ffffff;
	}

	.hero_subtitle {
		font-size: 0.88rem;
		line-height: 1.6;
		color: oklch(1 0 0 / 0.68);
		max-width: 36ch;
	}

	.hero_stats {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		flex-wrap: wrap;
		padding-top: 0.85rem;
		border-top: 1px solid oklch(1 0 0 / 0.14);
	}

	.hero_stat {
		display: flex;
		flex-direction: column;
		gap: 0.08rem;
	}

	.stat_num {
		font-size: 1rem;
		font-weight: 800;
		color: #ffffff;
		line-height: 1;
		font-family: 'Sora', sans-serif;
	}

	.stat_lbl {
		font-size: 0.68rem;
		color: oklch(1 0 0 / 0.55);
		line-height: 1;
	}

	.stat_div {
		width: 1px;
		height: 30px;
		background: oklch(1 0 0 / 0.18);
		flex-shrink: 0;
	}

	/* ── Section commons ── */
	.section_title {
		font-size: 1.05rem;
		font-weight: 800;
		color: var(--ecochitas-ink);
	}

	.section_sub {
		font-size: 0.82rem;
		color: var(--ecochitas-muted);
		margin-top: 0.22rem;
		line-height: 1.5;
	}

	.section_head_row {
		display: flex;
		align-items: center;
		gap: 0.45rem;
	}

	/* ── Steps ── */
	.steps_row {
		display: flex;
		align-items: flex-start;
		gap: 0;
		margin-top: 1.1rem;
		overflow-x: auto;
		padding-bottom: 0.25rem;
	}

	.step_item {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		flex: 1;
		min-width: 80px;
		gap: 0.55rem;
	}

	.step_icon_wrap {
		position: relative;
		width: 54px;
		height: 54px;
		background: oklch(0.96 0.02 150);
		border: 1.5px solid var(--ecochitas-border);
		border-radius: 1.1rem;
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--ecochitas-leaf);
	}

	.step_num {
		position: absolute;
		top: -7px;
		right: -7px;
		width: 18px;
		height: 18px;
		background: var(--ecochitas-leaf);
		color: white;
		font-size: 0.62rem;
		font-weight: 900;
		border-radius: 999px;
		display: flex;
		align-items: center;
		justify-content: center;
		font-family: 'Sora', sans-serif;
	}

	.step_title {
		font-size: 0.82rem;
		font-weight: 800;
		color: var(--ecochitas-ink);
		line-height: 1.2;
	}

	.step_desc {
		font-size: 0.72rem;
		color: var(--ecochitas-muted);
		line-height: 1.4;
	}

	.step_connector {
		flex-shrink: 0;
		width: 24px;
		display: flex;
		align-items: center;
		justify-content: center;
		padding-top: 27px;
	}

	.connector_line {
		width: 100%;
		height: 1.5px;
		background: var(--ecochitas-border);
		border-radius: 999px;
	}

	/* ── Category picker ── */
	.cat_picker {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 0.55rem;
		margin-top: 1.1rem;
	}

	.cat_btn {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		gap: 0.4rem;
		padding: 0.85rem 0.5rem 0.7rem;
		border-radius: 1.1rem;
		border: 1.5px solid var(--ecochitas-border);
		background: var(--ecochitas-surface);
		cursor: pointer;
		transition:
			border-color 0.15s,
			background 0.15s,
			transform 0.12s,
			box-shadow 0.15s;
	}

	.cat_btn:hover {
		border-color: var(--accent-border);
		background: var(--accent-light);
		transform: translateY(-1px);
	}

	.cat_btn_active {
		border-color: var(--accent);
		background: var(--accent);
		transform: translateY(-2px);
		box-shadow: 0 6px 20px color-mix(in srgb, var(--accent) 30%, transparent);
	}

	.cat_btn_active .cat_btn_name,
	.cat_btn_active .cat_btn_price {
		color: white;
	}

	.cat_btn_active .cat_btn_icon {
		color: white;
	}

	.cat_btn_icon {
		color: var(--accent);
		display: flex;
	}

	.cat_btn_name {
		font-size: 0.72rem;
		font-weight: 800;
		color: var(--ecochitas-ink);
		line-height: 1.2;
		font-family: 'Sora', sans-serif;
	}

	.cat_btn_price {
		font-size: 0.62rem;
		font-weight: 700;
		color: var(--accent);
		background: color-mix(in srgb, var(--accent) 10%, transparent);
		padding: 0.15rem 0.45rem;
		border-radius: 999px;
		line-height: 1.4;
	}

	/* ── Category detail ── */
	.cat_detail {
		margin-top: 0.8rem;
		border: 1.5px solid var(--accent-border);
		border-radius: 1.25rem;
		overflow: hidden;
		background: var(--accent-light);
	}

	.cat_detail_head {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.95rem 1rem;
		border-bottom: 1.5px solid var(--accent-border);
		background: color-mix(in srgb, var(--accent-light) 60%, white);
	}

	.cat_detail_icon {
		flex-shrink: 0;
		width: 46px;
		height: 46px;
		background: white;
		border: 1.5px solid var(--accent-border);
		border-radius: 0.9rem;
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--accent);
	}

	.cat_detail_meta {
		flex: 1;
		min-width: 0;
	}

	.cat_detail_name {
		font-size: 0.92rem;
		font-weight: 800;
		color: var(--accent);
		line-height: 1.2;
	}

	.cat_detail_tip {
		font-size: 0.74rem;
		color: var(--ecochitas-muted);
		display: flex;
		align-items: center;
		gap: 0.28rem;
		margin-top: 0.25rem;
		line-height: 1.4;
	}

	.cat_price_box {
		flex-shrink: 0;
		display: flex;
		flex-direction: column;
		align-items: flex-end;
		gap: 0.05rem;
		background: white;
		border: 1.5px solid var(--accent-border);
		border-radius: 0.85rem;
		padding: 0.45rem 0.75rem;
	}

	.price_box_lbl {
		font-size: 0.62rem;
		color: var(--ecochitas-muted);
		text-align: right;
	}

	.price_box_val {
		font-size: 0.95rem;
		font-weight: 900;
		color: var(--accent);
		line-height: 1.1;
		font-family: 'Sora', sans-serif;
	}

	.price_box_unit {
		font-size: 0.6rem;
		color: var(--ecochitas-muted);
		text-align: right;
	}

	/* ── Items grid ── */
	.items_grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1px;
		background: var(--accent-border);
	}

	.item_card {
		display: flex;
		align-items: center;
		gap: 0.7rem;
		padding: 0.8rem 0.9rem;
		background: white;
		transition: background 0.12s;
	}

	.item_card:hover {
		background: var(--accent-light);
	}

	.item_icon {
		flex-shrink: 0;
		width: 40px;
		height: 40px;
		background: var(--accent-light);
		border: 1.5px solid var(--accent-border);
		border-radius: 0.75rem;
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--accent);
	}

	.item_text {
		display: flex;
		flex-direction: column;
		gap: 0.12rem;
		min-width: 0;
	}

	.item_name {
		font-size: 0.8rem;
		font-weight: 800;
		color: var(--ecochitas-ink);
		line-height: 1.2;
		font-family: 'Sora', sans-serif;
	}

	.item_desc {
		font-size: 0.7rem;
		color: var(--ecochitas-muted);
		line-height: 1.3;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	/* ── Ranking ── */
	.ranking_list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		margin-top: 0.9rem;
	}

	.rank_row {
		display: grid;
		grid-template-columns: 26px 26px 1fr auto auto;
		align-items: center;
		gap: 0.55rem;
		padding: 0.65rem 0.85rem;
		background: var(--accent-light);
		border: 1.5px solid var(--accent-border);
		border-radius: 1rem;
	}

	.rank_pos {
		font-size: 0.72rem;
		font-weight: 900;
		color: var(--accent);
		text-align: center;
		font-family: 'Sora', sans-serif;
	}

	.rank_icon {
		display: flex;
		color: var(--accent);
	}

	.rank_name {
		font-size: 0.82rem;
		font-weight: 700;
		color: var(--ecochitas-ink);
		white-space: nowrap;
	}

	.rank_bar_track {
		height: 8px;
		background: color-mix(in srgb, var(--accent) 14%, white);
		border-radius: 999px;
		overflow: hidden;
		display: flex;
		align-items: center;
	}

	.rank_bar_fill {
		height: 100%;
		background: linear-gradient(90deg, var(--accent), color-mix(in srgb, var(--accent) 60%, white));
		border-radius: 999px;
		transition: width 0.35s ease;
	}

	.rank_special_label {
		font-size: 0.66rem;
		color: var(--ecochitas-muted);
		font-style: italic;
		white-space: nowrap;
		padding-left: 0.3rem;
	}

	.rank_price {
		font-size: 0.8rem;
		font-weight: 800;
		color: var(--accent);
		white-space: nowrap;
		text-align: right;
		min-width: 76px;
		font-family: 'Sora', sans-serif;
	}

	/* ── Rules ── */
	.rules_cols {
		display: grid;
		gap: 0.75rem;
		margin-top: 0.9rem;
	}

	.rules_group {
		border-radius: 1.1rem;
		overflow: hidden;
		border: 1.5px solid var(--ecochitas-border);
	}

	.rules_group_head {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.65rem 1rem;
		font-size: 0.82rem;
		font-weight: 800;
		font-family: 'Sora', sans-serif;
	}

	.rules_head_icon {
		width: 24px;
		height: 24px;
		border-radius: 6px;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}

	.rules_head_label {
		font-size: 0.82rem;
		font-weight: 800;
	}

	.rules_do .rules_group_head {
		background: oklch(0.93 0.06 150);
		border-bottom: 1.5px solid oklch(0.84 0.08 150);
		color: oklch(0.28 0.1 155);
	}

	.do_icon {
		background: oklch(0.84 0.08 150);
		color: oklch(0.28 0.1 155);
	}

	.rules_dont .rules_group_head {
		background: oklch(0.94 0.05 25);
		border-bottom: 1.5px solid oklch(0.87 0.07 25);
		color: oklch(0.38 0.14 25);
	}

	.dont_icon {
		background: oklch(0.87 0.07 25);
		color: oklch(0.38 0.14 25);
	}

	.rules_list {
		margin: 0;
		padding: 0.75rem 1rem 0.75rem 1.85rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		background: white;
	}

	.rules_list li {
		font-size: 0.82rem;
		color: var(--ecochitas-ink);
		line-height: 1.5;
	}

	/* ── GPS Button ── */
	.gps_btn {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.55rem;
		width: 100%;
		padding: 0.85rem 1rem;
		margin-top: 1.1rem;
		border-radius: 1rem;
		border: none;
		background: var(--ecochitas-sky, #3b82f6);
		color: white;
		font-size: 0.88rem;
		font-weight: 800;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		cursor: pointer;
		transition:
			background 0.15s,
			transform 0.12s,
			box-shadow 0.15s;
		box-shadow: 0 4px 16px color-mix(in srgb, var(--ecochitas-sky, #3b82f6) 35%, transparent);
	}

	.gps_btn:hover {
		background: color-mix(in srgb, var(--ecochitas-sky, #3b82f6) 85%, #000);
		transform: translateY(-1px);
		box-shadow: 0 6px 20px color-mix(in srgb, var(--ecochitas-sky, #3b82f6) 40%, transparent);
	}

	.gps_btn_active {
		background: #16a34a;
		box-shadow: 0 4px 16px color-mix(in srgb, #16a34a 35%, transparent);
	}

	.gps_btn_active:hover {
		background: #15803d;
		box-shadow: 0 6px 20px color-mix(in srgb, #16a34a 40%, transparent);
	}

	/* ── Violation type picker ── */
	.field_label {
		font-size: 0.82rem;
		font-weight: 700;
		color: var(--ecochitas-ink);
		margin: 1.1rem 0 0.55rem;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
	}

	.field_required {
		color: #ef4444;
	}

	.violation_grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0.45rem;
	}

	.vtype_btn {
		display: flex;
		align-items: center;
		gap: 0.45rem;
		padding: 0.65rem 0.8rem;
		border-radius: 0.85rem;
		border: 1.5px solid var(--ecochitas-border);
		background: var(--ecochitas-surface);
		color: var(--ecochitas-ink);
		font-size: 0.76rem;
		font-weight: 600;
		cursor: pointer;
		text-align: left;
		transition:
			border-color 0.15s,
			background 0.15s,
			color 0.12s,
			transform 0.1s;
		line-height: 1.3;
	}

	.vtype_btn:hover {
		border-color: #fca5a5;
		background: #fff1f2;
		color: #b91c1c;
		transform: translateY(-1px);
	}

	:root[data-theme='dark'] .vtype_btn:hover {
		background: rgba(239, 68, 68, 0.12);
		border-color: rgba(239, 68, 68, 0.4);
		color: #fca5a5;
	}

	.vtype_btn_active {
		border-color: #ef4444;
		background: #fee2e2;
		color: #b91c1c;
	}

	:root[data-theme='dark'] .vtype_btn_active {
		border-color: #ef4444;
		background: rgba(239, 68, 68, 0.15);
		color: #fca5a5;
	}

	/* ── Submit button ── */
	.submit_btn {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.55rem;
		width: 100%;
		padding: 0.9rem 1rem;
		margin-top: 0.85rem;
		border-radius: 1rem;
		border: none;
		background: #ef4444;
		color: white;
		font-size: 0.88rem;
		font-weight: 800;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		cursor: pointer;
		transition:
			background 0.15s,
			transform 0.12s,
			box-shadow 0.15s,
			opacity 0.15s;
		box-shadow: 0 4px 16px rgba(239, 68, 68, 0.35);
	}

	.submit_btn:hover:not(:disabled) {
		background: #dc2626;
		transform: translateY(-1px);
		box-shadow: 0 6px 20px rgba(239, 68, 68, 0.45);
	}

	.submit_btn:disabled {
		opacity: 0.45;
		cursor: not-allowed;
		transform: none;
		box-shadow: none;
	}

	.submit_btn_sent {
		background: #16a34a !important;
		box-shadow: 0 4px 16px rgba(22, 163, 74, 0.35) !important;
		opacity: 1 !important;
		cursor: default !important;
	}

	/* ── Privacy notice ── */
	.privacy_notice {
		display: flex;
		align-items: flex-start;
		gap: 0.55rem;
		padding: 0.75rem 0.9rem;
		margin-top: 0.75rem;
		border-radius: 0.85rem;
		background: #eff6ff;
		border: 1px solid #bfdbfe;
		font-size: 0.76rem;
		color: #1d4ed8;
		line-height: 1.5;
	}

	:root[data-theme='dark'] .privacy_notice {
		background: rgba(59, 130, 246, 0.1);
		border-color: rgba(59, 130, 246, 0.25);
		color: #93c5fd;
	}

	.privacy_notice strong {
		color: #1e40af;
	}

	:root[data-theme='dark'] .privacy_notice strong {
		color: #60a5fa;
	}

	/* ── Contaminadores ── */
	.contam_panel {
		border-color: #fecaca;
		background: color-mix(in srgb, #fff5f5 80%, var(--ecochitas-surface));
	}

	:root[data-theme='dark'] .contam_panel {
		border-color: rgba(239, 68, 68, 0.25);
		background: rgba(239, 68, 68, 0.04);
	}

	.contam_title {
		color: #b91c1c;
	}

	:root[data-theme='dark'] .contam_title {
		color: #fca5a5;
	}

	.contam_list {
		display: flex;
		flex-direction: column;
		gap: 0.6rem;
		margin-top: 1rem;
	}

	.contam_card {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.9rem 1rem;
		background: white;
		border: 1.5px solid #fecaca;
		border-radius: 1.1rem;
		transition: box-shadow 0.15s;
	}

	:root[data-theme='dark'] .contam_card {
		background: rgba(255, 255, 255, 0.04);
		border-color: rgba(239, 68, 68, 0.2);
	}

	.contam_card:hover {
		box-shadow: 0 4px 16px rgba(239, 68, 68, 0.12);
	}

	.contam_avatar {
		flex-shrink: 0;
		width: 44px;
		height: 44px;
		background: #fee2e2;
		border: 1.5px solid #fecaca;
		border-radius: 0.9rem;
		display: flex;
		align-items: center;
		justify-content: center;
		color: #ef4444;
	}

	:root[data-theme='dark'] .contam_avatar {
		background: rgba(239, 68, 68, 0.15);
		border-color: rgba(239, 68, 68, 0.3);
	}

	.contam_info {
		flex: 1;
		min-width: 0;
		display: flex;
		flex-direction: column;
		gap: 0.3rem;
	}

	.contam_name_row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.5rem;
	}

	.contam_name {
		font-size: 0.85rem;
		font-weight: 800;
		color: var(--ecochitas-ink);
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
	}

	.contam_badge {
		flex-shrink: 0;
		font-size: 0.68rem;
		font-weight: 800;
		color: white;
		background: #ef4444;
		padding: 0.18rem 0.55rem;
		border-radius: 999px;
		font-family: 'Sora', sans-serif;
	}

	.contam_location {
		display: flex;
		align-items: center;
		gap: 0.3rem;
		font-size: 0.75rem;
		color: var(--ecochitas-muted);
	}

	.contam_meta_row {
		display: flex;
		align-items: center;
		gap: 0.55rem;
		flex-wrap: wrap;
	}

	.contam_status {
		display: inline-flex;
		align-items: center;
		gap: 0.3rem;
		font-size: 0.68rem;
		font-weight: 700;
		padding: 0.18rem 0.5rem;
		border-radius: 999px;
		line-height: 1.4;
	}

	.contam_status_pending {
		background: #fef2f2;
		color: #b91c1c;
		border: 1px solid #fecaca;
	}

	:root[data-theme='dark'] .contam_status_pending {
		background: rgba(239, 68, 68, 0.12);
		color: #fca5a5;
		border-color: rgba(239, 68, 68, 0.3);
	}

	.contam_status_warning {
		background: #fffbeb;
		color: #92400e;
		border: 1px solid #fde68a;
	}

	:root[data-theme='dark'] .contam_status_warning {
		background: rgba(245, 158, 11, 0.12);
		color: #fcd34d;
		border-color: rgba(245, 158, 11, 0.3);
	}

	.contam_status_investigating {
		background: #eff6ff;
		color: #1e40af;
		border: 1px solid #bfdbfe;
	}

	:root[data-theme='dark'] .contam_status_investigating {
		background: rgba(59, 130, 246, 0.12);
		color: #93c5fd;
		border-color: rgba(59, 130, 246, 0.3);
	}

	.contam_date {
		display: inline-flex;
		align-items: center;
		gap: 0.28rem;
		font-size: 0.68rem;
		color: var(--ecochitas-muted);
	}

	/* ── Responsive ── */
	@media (min-width: 500px) {
		.cat_picker {
			grid-template-columns: repeat(5, 1fr);
		}

		.items_grid {
			grid-template-columns: repeat(3, 1fr);
		}

		.rules_cols {
			grid-template-columns: 1fr 1fr;
		}
	}

	@media (min-width: 860px) {
		.hero_title {
			font-size: 2.1rem;
		}

		.items_grid {
			grid-template-columns: repeat(3, 1fr);
		}
	}
</style>
