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
