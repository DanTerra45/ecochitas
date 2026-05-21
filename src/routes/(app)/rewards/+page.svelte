<script lang="ts">
	import EcochitasIcon from '$lib/components/ecochitas/EcochitasIcon.svelte';

	type Zone_score_card = {
		zone_name: string;
		recyclable_kilograms: number;
		contamination_percentage: number;
		houses_with_evidence: number;
		monthly_points: number;
	};

	const zone_scoreboard: Zone_score_card[] = [
		{
			zone_name: 'Sarco',
			recyclable_kilograms: 420,
			contamination_percentage: 9,
			houses_with_evidence: 84,
			monthly_points: 1320
		},
		{
			zone_name: 'Queru Queru',
			recyclable_kilograms: 388,
			contamination_percentage: 12,
			houses_with_evidence: 73,
			monthly_points: 1184
		},
		{
			zone_name: 'Tiquipaya Sur',
			recyclable_kilograms: 356,
			contamination_percentage: 7,
			houses_with_evidence: 67,
			monthly_points: 1152
		}
	];
</script>

<section class="panel">
	<div class="section_title_row">
		<EcochitasIcon name="rewards" size={18} />
		<h2>Modelo de bonificacion por zonas</h2>
	</div>
	<p class="muted_text">
		Se prioriza puntaje barrial para simplificar el MVP y reducir fraude individual.
	</p>
	<div class="formula_box">
		<p><strong>Puntaje base:</strong> kilogramos validos x factor_material</p>
		<p>
			<strong>Multiplicador de participacion:</strong> hogares con evidencia / hogares asignados
		</p>
		<p><strong>Penalizacion:</strong> contaminacion detectada en contenedor reciclable</p>
		<p><strong>Bono mensual:</strong> +12% si zona mantiene contaminacion menor o igual a 10%</p>
	</div>
</section>

<section class="panel">
	<div class="section_title_row">
		<EcochitasIcon name="map" size={18} />
		<h2>Ranking mensual de zonas</h2>
	</div>
	<div class="scoreboard_table">
		<div class="scoreboard_row scoreboard_header">
			<span>Zona</span>
			<span>Kg validos</span>
			<span>Contam.</span>
			<span>Hogares activos</span>
			<span>Puntos</span>
		</div>
		{#each zone_scoreboard as zone_score (zone_score.zone_name)}
			<div class="scoreboard_row">
				<span>{zone_score.zone_name}</span>
				<span>{zone_score.recyclable_kilograms}</span>
				<span>{zone_score.contamination_percentage}%</span>
				<span>{zone_score.houses_with_evidence}</span>
				<span>{zone_score.monthly_points}</span>
			</div>
		{/each}
	</div>
</section>

<section class="panel">
	<div class="section_title_row">
		<EcochitasIcon name="schedule" size={18} />
		<h2>Distribucion de recompensas</h2>
	</div>
	<ul class="rule_list">
		<li>70% del puntaje para hogares con evidencia valida en la semana.</li>
		<li>20% para fondo comun de la zona (mejoras barriales y limpieza).</li>
		<li>10% para incentivo del recolector por calidad de clasificacion.</li>
		<li>Si una zona no reporta evidencia, solo conserva puntaje comun minimo.</li>
	</ul>
</section>

<style>
	h2 {
		margin: 0;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		font-size: 1rem;
		letter-spacing: -0.015em;
	}

	.section_title_row {
		display: flex;
		align-items: center;
		gap: 0.45rem;
	}

	.formula_box {
		margin-top: 0.75rem;
		padding: 0.75rem;
		border-radius: 0.9rem;
		border: 1px solid var(--ecochitas-border);
		background: oklch(1 0 0 / 0.65);
		display: grid;
		gap: 0.35rem;
		font-size: 0.84rem;
	}

	.formula_box p {
		margin: 0;
	}

	.scoreboard_table {
		margin-top: 0.8rem;
		display: grid;
		gap: 0.45rem;
	}

	.scoreboard_row {
		display: grid;
		grid-template-columns: minmax(120px, 2fr) repeat(4, 1fr);
		gap: 0.4rem;
		padding: 0.6rem 0.65rem;
		border-radius: 0.8rem;
		border: 1px solid var(--ecochitas-border);
		background: oklch(1 0 0 / 0.68);
		font-size: 0.8rem;
	}

	.scoreboard_header {
		font-weight: 700;
		background: oklch(0.95 0.02 160 / 0.72);
	}

	.rule_list {
		margin: 0.75rem 0 0;
		padding-left: 1.1rem;
		display: grid;
		gap: 0.45rem;
		font-size: 0.85rem;
	}
</style>
