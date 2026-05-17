<script lang="ts">
	import { resolve } from '$app/paths';
	import type { Pathname } from '$app/types';
	import EcochitasIcon from '$lib/components/ecochitas/EcochitasIcon.svelte';

	type Quick_metric = {
		label: string;
		value: string;
		helper_text: string;
		icon_name: 'truck' | 'recycling' | 'alert' | 'rewards';
	};

	const quick_metrics: Quick_metric[] = [
		{
			label: 'Camiones conectados',
			value: '12',
			helper_text: 'seguimiento GPS activo',
			icon_name: 'truck'
		},
		{
			label: 'Contenedores monitoreados',
			value: '186',
			helper_text: 'sensores de llenado',
			icon_name: 'recycling'
		},
		{
			label: 'Alertas hoy',
			value: '17',
			helper_text: 'zonas con saturacion',
			icon_name: 'alert'
		},
		{
			label: 'Puntos comunitarios',
			value: '9,240',
			helper_text: 'bonificaciones del mes',
			icon_name: 'rewards'
		}
	];

	const map_route_path = '/map' as Pathname;
	const recycling_route_path = '/recycling' as Pathname;
	const rewards_route_path = '/rewards' as Pathname;
	const profile_route_path = '/profile' as Pathname;
</script>

<section class="panel hero_panel">
	<p class="eyebrow_text">Panel base del prototipo</p>
	<h2>Operacion municipal en una sola vista</h2>
	<p class="muted_text">
		Este frontend ya esta dividido por rutas reales para validar backend, GPS en tiempo real,
		reciclaje, bonificaciones y perfil ciudadano sin depender de una sola pantalla.
	</p>
</section>

<section class="metrics_grid">
	{#each quick_metrics as quick_metric (quick_metric.label)}
		<article class="panel metric_card">
			<div class="metric_icon">
				<EcochitasIcon name={quick_metric.icon_name} size={18} />
			</div>
			<p class="metric_label">{quick_metric.label}</p>
			<p class="metric_value">{quick_metric.value}</p>
			<p class="metric_helper">{quick_metric.helper_text}</p>
		</article>
	{/each}
</section>

<section class="panel action_panel">
	<div class="section_title_row">
		<EcochitasIcon name="map" size={18} />
		<h2>Flujo operativo recomendado</h2>
	</div>
	<ol class="action_list">
		<li>
			Monitorea ubicaciones en <a href={resolve(map_route_path)}>/map</a> y revisa estado de stream.
		</li>
		<li>
			Configura categorias y materiales en
			<a href={resolve(recycling_route_path)}>/recycling</a>.
		</li>
		<li>
			Ajusta reglas de puntos barriales en <a href={resolve(rewards_route_path)}>/rewards</a>.
		</li>
		<li>
			Valida datos de usuario y evidencia en <a href={resolve(profile_route_path)}>/profile</a>.
		</li>
	</ol>
</section>

<section class="panel updates_panel">
	<div class="section_title_row">
		<EcochitasIcon name="schedule" size={18} />
		<h2>Consideraciones del MVP</h2>
	</div>
	<ul class="update_list">
		<li>Martes/Jueves: recoleccion regular de residuos mixtos.</li>
		<li>Viernes: ruta exclusiva de material reciclable por zona.</li>
		<li>Puntos barriales se otorgan solo con evidencia valida por hogar.</li>
		<li>Contaminacion alta en contenedor reciclable descuenta puntaje comunitario.</li>
	</ul>
</section>

<style>
	h2 {
		margin: 0;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		font-size: 1rem;
		letter-spacing: -0.015em;
	}

	.hero_panel h2 {
		font-size: 1.22rem;
	}

	.metrics_grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0.75rem;
	}

	.metric_card {
		padding: 0.9rem;
	}

	.metric_icon {
		width: 2rem;
		height: 2rem;
		display: grid;
		place-items: center;
		border-radius: 0.7rem;
		background: oklch(0.93 0.06 160 / 0.45);
		color: var(--ecochitas-leaf);
	}

	.metric_label {
		margin: 0.6rem 0 0;
		font-size: 0.76rem;
		color: var(--ecochitas-muted);
	}

	.metric_value {
		margin: 0.2rem 0 0;
		font-size: 1.2rem;
		font-weight: 800;
		letter-spacing: -0.02em;
	}

	.metric_helper {
		margin: 0.12rem 0 0;
		font-size: 0.72rem;
		color: var(--ecochitas-muted);
	}

	.section_title_row {
		display: flex;
		align-items: center;
		gap: 0.45rem;
	}

	.action_list,
	.update_list {
		margin: 0.7rem 0 0;
		padding-left: 1.1rem;
		display: grid;
		gap: 0.45rem;
		color: var(--ecochitas-ink);
		font-size: 0.86rem;
	}

	.action_list a {
		color: var(--ecochitas-sky);
		font-weight: 700;
		text-decoration: none;
	}

	.action_list a:hover {
		text-decoration: underline;
	}

	@media (min-width: 900px) {
		.metrics_grid {
			grid-template-columns: repeat(4, minmax(0, 1fr));
		}
	}
</style>
