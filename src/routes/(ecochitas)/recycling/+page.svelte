<script lang="ts">
	import EcochitasIcon from '$lib/components/ecochitas/EcochitasIcon.svelte';

	type Recyclable_material = {
		category_name: string;
		examples: string;
		price_per_kilogram_bolivianos: string;
		quality_note: string;
	};

	const recyclable_materials: Recyclable_material[] = [
		{
			category_name: 'Papel y carton',
			examples: 'periodicos, revistas, cajas limpias, hojas',
			price_per_kilogram_bolivianos: '3.50 Bs/kg',
			quality_note: 'evitar papel humedo o mezclado con comida'
		},
		{
			category_name: 'Plastico PET',
			examples: 'botellas transparentes, envases de bebidas',
			price_per_kilogram_bolivianos: '4.00 Bs/kg',
			quality_note: 'entregar enjuagado y compactado'
		},
		{
			category_name: 'Vidrio',
			examples: 'botellas y frascos sin tapa',
			price_per_kilogram_bolivianos: '1.60 Bs/kg',
			quality_note: 'separar por color reduce rechazos'
		},
		{
			category_name: 'Metales',
			examples: 'latas de aluminio, chatarra ligera',
			price_per_kilogram_bolivianos: '6.20 Bs/kg',
			quality_note: 'sin residuos organicos adheridos'
		}
	];
</script>

<section class="panel">
	<div class="section_title_row">
		<EcochitasIcon name="recycling" size={18} />
		<h2>Categorias de reciclaje</h2>
	</div>
	<p class="muted_text">
		Guia de referencia para ciudadania y recolectores del turno de reciclaje de viernes.
	</p>

	<div class="materials_grid">
		{#each recyclable_materials as recyclable_material (recyclable_material.category_name)}
			<article class="material_card">
				<h3>{recyclable_material.category_name}</h3>
				<p><strong>Ejemplos:</strong> {recyclable_material.examples}</p>
				<p><strong>Valor:</strong> {recyclable_material.price_per_kilogram_bolivianos}</p>
				<p><strong>Control de calidad:</strong> {recyclable_material.quality_note}</p>
			</article>
		{/each}
	</div>
</section>

<section class="panel">
	<div class="section_title_row">
		<EcochitasIcon name="alert" size={18} />
		<h2>Reglas de penalizacion por contaminacion</h2>
	</div>
	<ul class="rule_list">
		<li>Si el contenedor reciclable tiene residuos organicos se descuenta puntaje de zona.</li>
		<li>La notificacion muestra tipo de contaminante y evidencia fotografica del recolector.</li>
		<li>Cuando la contaminacion supera 20% del volumen se congela la bonificacion semanal.</li>
		<li>Los hogares con evidencia valida siguen acumulando puntaje individual.</li>
	</ul>
</section>

<section class="panel">
	<div class="section_title_row">
		<EcochitasIcon name="upload" size={18} />
		<h2>Evidencia minima para sumar puntos</h2>
	</div>
	<ol class="rule_list">
		<li>Foto del aporte reciclable en bolsa/transporte limpio.</li>
		<li>Seleccion de categoria principal del material.</li>
		<li>Geolocalizacion aproximada del hogar o contenedor asignado.</li>
		<li>Confirmacion del recolector al momento de pesaje o vaciado.</li>
	</ol>
</section>

<style>
	h2,
	h3 {
		margin: 0;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		letter-spacing: -0.015em;
	}

	h2 {
		font-size: 1rem;
	}

	h3 {
		font-size: 0.95rem;
	}

	.section_title_row {
		display: flex;
		align-items: center;
		gap: 0.45rem;
	}

	.materials_grid {
		display: grid;
		gap: 0.65rem;
		margin-top: 0.8rem;
	}

	.material_card {
		border: 1px solid var(--ecochitas-border);
		border-radius: 1rem;
		padding: 0.8rem;
		background: oklch(1 0 0 / 0.68);
	}

	.material_card p {
		margin: 0.35rem 0 0;
		font-size: 0.82rem;
		color: var(--ecochitas-ink);
	}

	.rule_list {
		margin: 0.75rem 0 0;
		padding-left: 1.1rem;
		display: grid;
		gap: 0.45rem;
		font-size: 0.85rem;
	}

	@media (min-width: 900px) {
		.materials_grid {
			grid-template-columns: 1fr 1fr;
		}
	}
</style>
