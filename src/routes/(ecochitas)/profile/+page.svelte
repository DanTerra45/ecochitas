<script lang="ts">
	let active_tab: 'individual' | 'condominio' = $state('individual');

	// ── Individual data ──
	const user = {
		name: 'María González',
		zone: 'Zona Norte',
		address: 'Casa #234',
		points: 1245,
		cash_equiv: 124.5,
		monthly_goal: 2000,
		deliveries: 18,
		recycled_kg: 42,
		co2_avoided_kg: 84,
		trees_equiv: 4
	};

	const monthly_pct = $derived(Math.min((user.points / user.monthly_goal) * 100, 100));

	// ── Condominio data ──
	const building = {
		name: 'Torres Sofer',
		address: 'Av. América',
		units: 24,
		ranking: 3,
		reward_pct: 15,
		community_points: 8450,
		community_goal: 12000
	};

	const community_pct = $derived(
		Math.min((building.community_points / building.community_goal) * 100, 100)
	);

	const top_contributors = [
		{ unit: 'Depto 301', points: 850 },
		{ unit: 'Depto 205', points: 720 },
		{ unit: 'Depto 102', points: 680 },
		{ unit: 'Depto 401', points: 540 },
		{ unit: 'Depto 104', points: 480 }
	];

	const community_benefits = [
		'Descuentos grupales en impuestos municipales',
		'Reconocimiento público como Edificio Eco-Responsable',
		'Acceso prioritario a programas de mejora urbana'
	];
</script>

<!-- ── Header ── -->
<div class="profile_header panel">
	<div class="avatar_wrap">
		<svg
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-linecap="round"
			stroke-linejoin="round"
			stroke-width="1.75"
			style="width:28px;height:28px"
			aria-hidden="true"
		>
			<circle cx="12" cy="8" r="3.5" />
			<path d="M5 20c1.6-3.4 4-5.2 7-5.2s5.4 1.8 7 5.2" />
		</svg>
	</div>
	<div class="header_info">
		<h1 class="user_name">{user.name}</h1>
		<p class="user_sub">{user.zone} · {user.address}</p>
	</div>
</div>

<!-- ── Tab switcher ── -->
<div class="tab_switcher">
	<button
		class="tab_btn"
		class:tab_active={active_tab === 'individual'}
		onclick={() => (active_tab = 'individual')}
	>
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:15px;height:15px" aria-hidden="true">
			<path d="M4 11.5L12 5l8 6.5" /><path d="M6 10.8V20h12v-9.2" /><path d="M10 20v-4h4v4" />
		</svg>
		Individual
	</button>
	<button
		class="tab_btn"
		class:tab_active={active_tab === 'condominio'}
		onclick={() => (active_tab = 'condominio')}
	>
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:15px;height:15px" aria-hidden="true">
			<rect x="3" y="3" width="18" height="18" rx="2"/><path d="M9 3v18M15 3v18M3 9h18M3 15h18"/>
		</svg>
		Condominio
	</button>
</div>

<!-- ══════════════ INDIVIDUAL TAB ══════════════ -->
{#if active_tab === 'individual'}

<!-- Eco-Wallet -->
<div class="wallet_card">
	<div class="wallet_header">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:18px;height:18px;opacity:0.85" aria-hidden="true">
			<rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/><path d="M6 11h.01M6 15h.01"/>
		</svg>
		<span class="wallet_label">Eco-Wallet</span>
	</div>

	<p class="points_label">Puntos Acumulados</p>
	<strong class="points_num">{user.points.toLocaleString()}</strong>

	<div class="cash_sub">
		<span class="cash_lbl">Valor Equivalente en Efectivo</span>
		<strong class="cash_val">{user.cash_equiv.toFixed(2)} Bs</strong>
	</div>

	<button class="redeem_btn">Canjear Puntos</button>
</div>

<!-- Progreso mensual -->
<div class="panel">
	<div class="section_row">
		<svg viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:18px;height:18px;flex-shrink:0" aria-hidden="true">
			<path d="M12 3l2.7 5.5 6.1.9-4.4 4.3 1 6.1L12 17l-5.4 2.8 1-6.1L3.2 9.4l6.1-.9z"/>
		</svg>
		<h2 class="section_title">Tu Progreso Este Mes</h2>
	</div>

	<div class="progress_row">
		<span class="progress_label">Meta Mensual</span>
		<span class="progress_val" style="color:var(--ecochitas-sky)">{user.points.toLocaleString()} / {user.monthly_goal.toLocaleString()}</span>
	</div>
	<div class="progress_track">
		<div class="progress_fill progress_fill_green" style="width:{monthly_pct}%"></div>
	</div>

	<div class="stats_duo">
		<div class="stat_pill stat_pill_green">
			<strong class="stat_big">{user.deliveries}</strong>
			<span class="stat_small">Entregas</span>
		</div>
		<div class="stat_pill stat_pill_blue">
			<strong class="stat_big stat_blue">{user.recycled_kg}kg</strong>
			<span class="stat_small">Reciclado</span>
		</div>
	</div>
</div>

<!-- Impacto ambiental -->
<div class="impact_card">
	<div class="impact_header">
		<span class="impact_emoji" aria-hidden="true">🌱</span>
		<h2 class="impact_title">Impacto Ambiental</h2>
	</div>
	<p class="impact_desc">Tu contribución este mes ha evitado la emisión de</p>
	<strong class="impact_co2">{user.co2_avoided_kg}kg CO<sub>2</sub></strong>
	<p class="impact_equiv">Equivalente a plantar {user.trees_equiv} árboles</p>
</div>

<!-- ══════════════ CONDOMINIO TAB ══════════════ -->
{:else}

<!-- Building header -->
<div class="panel building_header_card">
	<div class="building_icon_wrap">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:22px;height:22px" aria-hidden="true">
			<rect x="3" y="3" width="18" height="18" rx="2"/><path d="M9 3v18M15 3v18M3 9h18M3 15h18"/>
		</svg>
	</div>
	<div>
		<h2 class="building_name">{building.name}</h2>
		<p class="building_sub">{building.address} · {building.units} departamentos</p>
	</div>
</div>

<!-- Recompensa top -->
<div class="reward_card">
	<div class="reward_header">
		<div class="reward_icon">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:18px;height:18px" aria-hidden="true">
				<circle cx="12" cy="8" r="6"/><path d="M15.477 12.89L17 22l-5-3-5 3 1.523-9.11"/>
			</svg>
		</div>
		<h3 class="reward_title">Recompensa de Condominio Top</h3>
	</div>
	<p class="reward_desc">
		{building.reward_pct}% de descuento en impuesto predial para todos los residentes este mes
	</p>
	<div class="reward_rank">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.2" style="width:13px;height:13px;flex-shrink:0" aria-hidden="true">
			<polyline points="22 7 13.5 15.5 8.5 10.5 2 17"/><polyline points="16 7 22 7 22 13"/>
		</svg>
		Ranking Actual: #{building.ranking}
	</div>
</div>

<!-- Meta comunitaria -->
<div class="panel">
	<div class="progress_row">
		<span class="progress_label">Meta Comunitaria Mensual</span>
		<span class="progress_val" style="color:var(--ecochitas-sky)">{building.community_points.toLocaleString()} / {building.community_goal.toLocaleString()}</span>
	</div>
	<div class="progress_track">
		<div class="progress_fill progress_fill_green" style="width:{community_pct}%"></div>
	</div>
	<p class="progress_hint">Faltan {(building.community_goal - building.community_points).toLocaleString()} puntos para alcanzar la meta</p>
</div>

<!-- Top contribuyentes -->
<div class="panel">
	<h2 class="section_title" style="margin-bottom:0.85rem">Top Contribuyentes del Edificio</h2>
	<div class="rank_list">
		{#each top_contributors as contrib, i (contrib.unit)}
			<div class="rank_row" class:rank_row_gold={i === 0}>
				<div class="rank_badge" class:rank_badge_gold={i === 0} class:rank_badge_silver={i === 1} class:rank_badge_bronze={i === 2}>
					{i + 1}
				</div>
				<span class="rank_unit">{contrib.unit}</span>
				<div class="rank_right">
					<strong class="rank_pts" class:rank_pts_gold={i === 0}>{contrib.points.toLocaleString()}</strong>
					<span class="rank_pts_lbl">puntos</span>
				</div>
			</div>
		{/each}
	</div>
</div>

<!-- Beneficios comunitarios -->
<div class="benefits_card">
	<div class="benefits_header">
		<span class="benefits_emoji" aria-hidden="true">🏆</span>
		<h2 class="benefits_title">Beneficios Comunitarios</h2>
	</div>
	<ul class="benefits_list">
		{#each community_benefits as benefit}
			<li class="benefit_item">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.2" style="width:12px;height:12px;flex-shrink:0;margin-top:2px;opacity:0.8" aria-hidden="true">
					<polyline points="20 6 9 17 4 12"/>
				</svg>
				{benefit}
			</li>
		{/each}
	</ul>
</div>

{/if}

<style>
	h1, h2, h3 { margin: 0; font-family: 'Sora', 'Plus Jakarta Sans', sans-serif; letter-spacing: -0.02em; }
	p { margin: 0; }
	sub { font-size: 0.6em; }

	/* ── Profile Header ── */
	.profile_header {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 1.1rem 1.25rem;
	}

	.avatar_wrap {
		flex-shrink: 0;
		width: 56px;
		height: 56px;
		border-radius: 50%;
		background: linear-gradient(135deg, #0d9488, #16a34a);
		display: flex;
		align-items: center;
		justify-content: center;
		color: white;
		box-shadow: 0 4px 14px rgba(13, 148, 136, 0.4);
	}

	.user_name {
		font-size: 1.1rem;
		font-weight: 900;
		color: var(--ecochitas-ink);
		line-height: 1.2;
	}

	.user_sub {
		font-size: 0.78rem;
		color: var(--ecochitas-muted);
		margin-top: 0.18rem;
	}

	.header_info { flex: 1; min-width: 0; }

	/* ── Tab Switcher ── */
	.tab_switcher {
		display: grid;
		grid-template-columns: 1fr 1fr;
		background: var(--ecochitas-surface);
		border: 1px solid var(--ecochitas-border);
		border-radius: 1.2rem;
		padding: 0.3rem;
		gap: 0.3rem;
	}

	.tab_btn {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.45rem;
		padding: 0.65rem 1rem;
		border-radius: 0.9rem;
		border: none;
		background: transparent;
		color: var(--ecochitas-muted);
		font-size: 0.83rem;
		font-weight: 700;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		cursor: pointer;
		transition: background 0.15s, color 0.15s, box-shadow 0.15s;
	}

	.tab_active {
		background: var(--ecochitas-leaf);
		color: white;
		box-shadow: 0 3px 12px color-mix(in srgb, var(--ecochitas-leaf) 40%, transparent);
	}

	/* ── Eco-Wallet ── */
	.wallet_card {
		background: linear-gradient(145deg, oklch(0.30 0.11 163), oklch(0.22 0.09 175));
		border-radius: 1.5rem;
		padding: 1.4rem 1.5rem 1.2rem;
		display: flex;
		flex-direction: column;
		gap: 0.3rem;
		position: relative;
		overflow: hidden;
	}

	.wallet_card::before {
		content: '';
		position: absolute;
		top: -30px;
		right: -30px;
		width: 120px;
		height: 120px;
		border-radius: 50%;
		background: oklch(1 0 0 / 0.04);
		pointer-events: none;
	}

	.wallet_header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: oklch(1 0 0 / 0.7);
		margin-bottom: 0.5rem;
	}

	.wallet_label {
		font-size: 0.82rem;
		font-weight: 800;
		font-family: 'Sora', sans-serif;
		letter-spacing: 0.04em;
		text-transform: uppercase;
	}

	.points_label {
		font-size: 0.75rem;
		color: oklch(1 0 0 / 0.55);
	}

	.points_num {
		font-size: 2.8rem;
		font-weight: 900;
		color: white;
		line-height: 1;
		letter-spacing: -0.04em;
		font-family: 'Sora', sans-serif;
		margin: 0.1rem 0 0.6rem;
	}

	.cash_sub {
		background: oklch(1 0 0 / 0.1);
		border: 1px solid oklch(1 0 0 / 0.15);
		border-radius: 0.9rem;
		padding: 0.6rem 0.9rem;
		display: flex;
		flex-direction: column;
		gap: 0.1rem;
		margin-bottom: 0.5rem;
	}

	.cash_lbl {
		font-size: 0.7rem;
		color: oklch(1 0 0 / 0.55);
	}

	.cash_val {
		font-size: 1.3rem;
		font-weight: 900;
		color: white;
		font-family: 'Sora', sans-serif;
		letter-spacing: -0.02em;
	}

	.redeem_btn {
		width: 100%;
		padding: 0.78rem 1rem;
		border-radius: 0.9rem;
		border: 1.5px solid oklch(1 0 0 / 0.3);
		background: oklch(1 0 0 / 0.08);
		color: white;
		font-size: 0.88rem;
		font-weight: 800;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		cursor: pointer;
		transition: background 0.15s, border-color 0.15s;
		margin-top: 0.15rem;
	}

	.redeem_btn:hover {
		background: oklch(1 0 0 / 0.15);
		border-color: oklch(1 0 0 / 0.5);
	}

	/* ── Section commons ── */
	.section_row {
		display: flex;
		align-items: center;
		gap: 0.45rem;
		margin-bottom: 0.85rem;
	}

	.section_title {
		font-size: 1rem;
		font-weight: 800;
		color: var(--ecochitas-ink);
	}

	/* ── Progress ── */
	.progress_row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.5rem;
		margin-bottom: 0.55rem;
	}

	.progress_label {
		font-size: 0.82rem;
		color: var(--ecochitas-muted);
		font-weight: 600;
	}

	.progress_val {
		font-size: 0.82rem;
		font-weight: 800;
		font-family: 'Sora', sans-serif;
		white-space: nowrap;
	}

	.progress_track {
		height: 10px;
		background: var(--ecochitas-border);
		border-radius: 999px;
		overflow: hidden;
	}

	.progress_fill {
		height: 100%;
		border-radius: 999px;
		transition: width 0.5s ease;
	}

	.progress_fill_green {
		background: linear-gradient(90deg, var(--ecochitas-leaf), #22c55e);
	}

	.progress_hint {
		font-size: 0.74rem;
		color: var(--ecochitas-muted);
		margin-top: 0.45rem;
	}

	/* ── Stats duo ── */
	.stats_duo {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0.6rem;
		margin-top: 0.85rem;
	}

	.stat_pill {
		border-radius: 1rem;
		padding: 0.9rem 0.75rem;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.2rem;
		text-align: center;
	}

	.stat_pill_green {
		background: #f0fdf4;
		border: 1.5px solid #bbf7d0;
	}

	:root[data-theme='dark'] .stat_pill_green {
		background: rgba(34, 197, 94, 0.08);
		border-color: rgba(34, 197, 94, 0.2);
	}

	.stat_pill_blue {
		background: #eff6ff;
		border: 1.5px solid #bfdbfe;
	}

	:root[data-theme='dark'] .stat_pill_blue {
		background: rgba(59, 130, 246, 0.08);
		border-color: rgba(59, 130, 246, 0.2);
	}

	.stat_big {
		font-size: 1.6rem;
		font-weight: 900;
		color: var(--ecochitas-leaf);
		line-height: 1;
		font-family: 'Sora', sans-serif;
		letter-spacing: -0.03em;
	}

	.stat_blue { color: var(--ecochitas-sky) !important; }

	.stat_small {
		font-size: 0.74rem;
		color: var(--ecochitas-muted);
		font-weight: 600;
	}

	/* ── Impact card ── */
	.impact_card {
		background: linear-gradient(145deg, #16a34a, #15803d);
		border-radius: 1.5rem;
		padding: 1.4rem 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
		position: relative;
		overflow: hidden;
	}

	.impact_card::after {
		content: '';
		position: absolute;
		bottom: -20px;
		right: -20px;
		width: 100px;
		height: 100px;
		border-radius: 50%;
		background: oklch(1 0 0 / 0.05);
		pointer-events: none;
	}

	.impact_header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		margin-bottom: 0.15rem;
	}

	.impact_emoji { font-size: 1.1rem; }

	.impact_title {
		font-size: 1rem;
		font-weight: 900;
		color: white;
	}

	.impact_desc {
		font-size: 0.8rem;
		color: oklch(1 0 0 / 0.7);
	}

	.impact_co2 {
		font-size: 2.4rem;
		font-weight: 900;
		color: white;
		line-height: 1.1;
		letter-spacing: -0.04em;
		font-family: 'Sora', sans-serif;
		margin: 0.1rem 0;
	}

	.impact_equiv {
		font-size: 0.78rem;
		color: oklch(1 0 0 / 0.65);
	}

	/* ── Building header ── */
	.building_header_card {
		display: flex;
		align-items: center;
		gap: 0.9rem;
		padding: 1rem 1.25rem;
	}

	.building_icon_wrap {
		flex-shrink: 0;
		width: 48px;
		height: 48px;
		border-radius: 0.9rem;
		background: color-mix(in srgb, var(--ecochitas-leaf) 12%, transparent);
		border: 1.5px solid color-mix(in srgb, var(--ecochitas-leaf) 25%, transparent);
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--ecochitas-leaf);
	}

	.building_name {
		font-size: 1.05rem;
		font-weight: 900;
		color: var(--ecochitas-ink);
	}

	.building_sub {
		font-size: 0.78rem;
		color: var(--ecochitas-muted);
		margin-top: 0.18rem;
	}

	/* ── Reward card ── */
	.reward_card {
		background: linear-gradient(145deg, #fef9c3, #fef3c7);
		border: 1.5px solid #fde68a;
		border-radius: 1.5rem;
		padding: 1.1rem 1.25rem;
		display: flex;
		flex-direction: column;
		gap: 0.45rem;
	}

	:root[data-theme='dark'] .reward_card {
		background: rgba(251, 191, 36, 0.08);
		border-color: rgba(251, 191, 36, 0.25);
	}

	.reward_header {
		display: flex;
		align-items: center;
		gap: 0.55rem;
	}

	.reward_icon {
		width: 34px;
		height: 34px;
		border-radius: 0.65rem;
		background: #fde68a;
		border: 1px solid #fcd34d;
		display: flex;
		align-items: center;
		justify-content: center;
		color: #92400e;
		flex-shrink: 0;
	}

	:root[data-theme='dark'] .reward_icon {
		background: rgba(251, 191, 36, 0.2);
		border-color: rgba(251, 191, 36, 0.35);
		color: #fcd34d;
	}

	.reward_title {
		font-size: 0.88rem;
		font-weight: 800;
		color: #92400e;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
	}

	:root[data-theme='dark'] .reward_title { color: #fcd34d; }

	.reward_desc {
		font-size: 0.8rem;
		color: #78350f;
		line-height: 1.5;
	}

	:root[data-theme='dark'] .reward_desc { color: #fde68a; }

	.reward_rank {
		display: inline-flex;
		align-items: center;
		gap: 0.35rem;
		font-size: 0.78rem;
		font-weight: 800;
		color: #b45309;
		font-family: 'Sora', sans-serif;
	}

	:root[data-theme='dark'] .reward_rank { color: #fbbf24; }

	/* ── Rank list ── */
	.rank_list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.rank_row {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 0.85rem;
		border-radius: 1rem;
		border: 1.5px solid var(--ecochitas-border);
		background: var(--ecochitas-surface);
		transition: box-shadow 0.15s;
	}

	.rank_row:hover { box-shadow: 0 3px 12px rgba(0,0,0,0.06); }

	.rank_row_gold {
		background: #fefce8;
		border-color: #fde68a;
	}

	:root[data-theme='dark'] .rank_row_gold {
		background: rgba(251, 191, 36, 0.07);
		border-color: rgba(251, 191, 36, 0.25);
	}

	.rank_badge {
		flex-shrink: 0;
		width: 30px;
		height: 30px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 0.78rem;
		font-weight: 900;
		font-family: 'Sora', sans-serif;
		background: var(--ecochitas-border);
		color: var(--ecochitas-muted);
	}

	.rank_badge_gold   { background: #f59e0b; color: white; }
	.rank_badge_silver { background: #94a3b8; color: white; }
	.rank_badge_bronze { background: #c2855a; color: white; }

	.rank_unit {
		flex: 1;
		font-size: 0.85rem;
		font-weight: 700;
		color: var(--ecochitas-ink);
	}

	.rank_right {
		display: flex;
		flex-direction: column;
		align-items: flex-end;
		gap: 0.05rem;
	}

	.rank_pts {
		font-size: 1rem;
		font-weight: 900;
		color: var(--ecochitas-sky);
		font-family: 'Sora', sans-serif;
		letter-spacing: -0.02em;
	}

	.rank_pts_gold { color: #d97706; }

	.rank_pts_lbl {
		font-size: 0.65rem;
		color: var(--ecochitas-muted);
	}

	/* ── Benefits card ── */
	.benefits_card {
		background: linear-gradient(145deg, #7c3aed, #6d28d9);
		border-radius: 1.5rem;
		padding: 1.3rem 1.4rem;
		display: flex;
		flex-direction: column;
		gap: 0.85rem;
		position: relative;
		overflow: hidden;
	}

	.benefits_card::before {
		content: '';
		position: absolute;
		top: -40px;
		right: -40px;
		width: 130px;
		height: 130px;
		border-radius: 50%;
		background: oklch(1 0 0 / 0.05);
		pointer-events: none;
	}

	.benefits_header {
		display: flex;
		align-items: center;
		gap: 0.55rem;
	}

	.benefits_emoji { font-size: 1.15rem; }

	.benefits_title {
		font-size: 1rem;
		font-weight: 900;
		color: white;
	}

	.benefits_list {
		margin: 0;
		padding: 0;
		list-style: none;
		display: flex;
		flex-direction: column;
		gap: 0.55rem;
	}

	.benefit_item {
		display: flex;
		align-items: flex-start;
		gap: 0.5rem;
		font-size: 0.82rem;
		color: oklch(1 0 0 / 0.85);
		line-height: 1.5;
	}
</style>