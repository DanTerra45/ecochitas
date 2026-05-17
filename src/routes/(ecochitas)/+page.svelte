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

	let active_accordion = $state('about');
</script>

<div class="hero_background"></div>

<section class="hero_section">
	<h1 class="hero_title">
		Infraestructura operativa que <br />
		<span class="highlight_text">simplifica la ciudad</span>
	</h1>
	<p class="hero_subtitle">
		Monitorea rutas en tiempo real, gestiona bonificaciones de reciclaje y valida perfiles
		ciudadanos desde una sola plataforma integrada y eficiente.
	</p>
	<div class="hero_actions">
		<a href={resolve(map_route_path)} class="btn_primary">Abrir Mapa Operativo</a>
		<a href="#flujo" class="btn_secondary">Ver Flujo</a>
	</div>
</section>

<section class="metrics_grid">
	{#each quick_metrics as quick_metric (quick_metric.label)}
		<article class="panel metric_card">
			<div class="metric_icon">
				<EcochitasIcon name={quick_metric.icon_name} size={18} />
			</div>
			<p class="metric_value">{quick_metric.value}</p>
			<p class="metric_label">{quick_metric.label}</p>
			<p class="metric_helper">{quick_metric.helper_text}</p>
		</article>
	{/each}
</section>

<section class="mission_section">
	<h2 class="mission_title"><span class="highlight_text">Nuestra</span> visión</h2>
	<div class="mission_grid">
		<div class="mission_accordion">
			<button
				class="accordion_item"
				class:active={active_accordion === 'about'}
				onclick={() => (active_accordion = 'about')}
			>
				<div class="accordion_header">
					<h3>Sobre Nosotros</h3>
					<span class="accordion_icon">{active_accordion === 'about' ? '−' : '+'}</span>
				</div>
				{#if active_accordion === 'about'}
					<div class="accordion_body">
						<p>
							EcoChitas es una iniciativa tecnológica nacida en Cochabamba con la misión de
							modernizar la gestión de residuos sólidos. Combinamos hardware (sensores IoT y GPS)
							con software avanzado para brindar herramientas a los operadores municipales y motivar
							a los ciudadanos a participar activamente en el reciclaje.
						</p>
					</div>
				{/if}
			</button>

			<button
				class="accordion_item"
				class:active={active_accordion === 'what'}
				onclick={() => (active_accordion = 'what')}
			>
				<div class="accordion_header">
					<h3>Qué Hacemos</h3>
					<span class="accordion_icon">{active_accordion === 'what' ? '−' : '+'}</span>
				</div>
				{#if active_accordion === 'what'}
					<div class="accordion_body">
						<p>
							Digitalizamos el ciclo completo de la recolección: desde el monitoreo en tiempo real
							de los camiones basureros y sensores de llenado en contenedores, hasta la validación
							de materiales reciclables entregados por los vecinos, otorgándoles puntos canjeables
							por beneficios locales.
						</p>
					</div>
				{/if}
			</button>

			<button
				class="accordion_item"
				class:active={active_accordion === 'why'}
				onclick={() => (active_accordion = 'why')}
			>
				<div class="accordion_header">
					<h3>Por Qué lo Hacemos</h3>
					<span class="accordion_icon">{active_accordion === 'why' ? '−' : '+'}</span>
				</div>
				{#if active_accordion === 'why'}
					<div class="accordion_body">
						<p>
							Creemos que una ciudad inteligente es una ciudad limpia. Al optimizar las rutas y
							recompensar las buenas prácticas ciudadanas, reducimos los costos operativos,
							disminuimos la contaminación ambiental y fomentamos una economía circular comunitaria
							sostenible a largo plazo.
						</p>
					</div>
				{/if}
			</button>
		</div>
	</div>
</section>

<footer class="page_footer">
	<div class="footer_content">
		<div class="footer_left">
			<div class="social_links">
				<a href="https://x.com" aria-label="X (Twitter)"><EcochitasIcon name="x" size={16} /></a>
				<a href="https://linkedin.com" aria-label="LinkedIn"
					><EcochitasIcon name="linkedin" size={16} /></a
				>
				<a href="https://github.com" aria-label="GitHub"
					><EcochitasIcon name="github" size={16} /></a
				>
				<a href="https://youtube.com" aria-label="YouTube"
					><EcochitasIcon name="youtube" size={16} /></a
				>
			</div>
			<p class="copyright">© EcoChitas 2026</p>
		</div>
		<div class="footer_right">
			<div class="footer_nav_column">
				<h4>Popular Examples</h4>
				<a href="/examples/stream">Serve your own data stream</a>
				<a href="/examples/custom-art">Create custom map art</a>
				<a href="/examples/analyze">Analyze zone saturation</a>
				<a href="/examples/run-hundreds">Run hundreds of simulations</a>
			</div>
			<div class="footer_nav_column">
				<h4>Company</h4>
				<a href="/about">About</a>
				<a href="/blog">Blog</a>
				<a href="/careers">Careers</a>
				<a href="/events">Events</a>
				<a href="/privacy">Privacy Policy</a>
				<a href="/terms">Terms</a>
			</div>
		</div>
	</div>
</footer>

<style>
	.hero_background {
		position: absolute;
		top: -10%;
		left: 50%;
		transform: translateX(-50%);
		width: 100vw;
		height: 80vh;
		background: var(--ecochitas-hero-bg);
		pointer-events: none;
		z-index: 0;
	}

	.hero_section {
		position: relative;
		text-align: center;
		padding: 4rem 1rem 3rem;
		z-index: 1;
		max-width: 800px;
		margin: 0 auto;
	}

	.hero_title {
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		font-size: 2.2rem;
		line-height: 1.1;
		font-weight: 700;
		letter-spacing: -0.03em;
		margin: 0 0 1.2rem;
		color: var(--ecochitas-ink);
	}

	.highlight_text {
		color: var(--ecochitas-leaf);
	}

	.hero_subtitle {
		font-size: 1.1rem;
		line-height: 1.5;
		color: var(--ecochitas-muted);
		margin: 0 auto 2rem;
		max-width: 600px;
	}

	.hero_actions {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 1rem;
		flex-wrap: wrap;
	}

	.btn_primary {
		background: var(--ecochitas-leaf);
		color: #000000;
		padding: 0.8rem 1.5rem;
		border-radius: 2rem;
		font-weight: 600;
		text-decoration: none;
		transition:
			transform 0.2s ease,
			box-shadow 0.2s ease;
		box-shadow: 0 0 20px rgba(22, 163, 74, 0.3);
	}

	:root[data-theme='dark'] .btn_primary {
		box-shadow: 0 0 20px rgba(85, 255, 119, 0.3);
	}

	.btn_primary:hover {
		transform: translateY(-2px);
		box-shadow: 0 0 30px rgba(22, 163, 74, 0.5);
	}

	:root[data-theme='dark'] .btn_primary:hover {
		box-shadow: 0 0 30px rgba(85, 255, 119, 0.5);
	}

	.btn_secondary {
		background: transparent;
		color: var(--ecochitas-ink);
		border: 1px solid var(--ecochitas-border);
		padding: 0.8rem 1.5rem;
		border-radius: 2rem;
		font-weight: 600;
		text-decoration: none;
		transition:
			background 0.2s ease,
			border-color 0.2s ease;
	}

	.btn_secondary:hover {
		background: var(--ecochitas-surface);
		border-color: var(--ecochitas-muted);
	}

	h2 {
		margin: 0;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		font-size: 1rem;
		letter-spacing: -0.015em;
		color: var(--ecochitas-ink);
	}

	.metrics_grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0.75rem;
		margin-bottom: 2rem;
		position: relative;
		z-index: 1;
	}

	.metric_card {
		padding: 1.2rem;
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
	}

	.metric_icon {
		width: 2.5rem;
		height: 2.5rem;
		display: grid;
		place-items: center;
		border-radius: 0.8rem;
		background: rgba(22, 163, 74, 0.1);
		color: var(--ecochitas-leaf);
		margin-bottom: 1rem;
	}

	:root[data-theme='dark'] .metric_icon {
		background: rgba(85, 255, 119, 0.1);
	}

	.metric_value {
		margin: 0;
		font-size: 1.5rem;
		font-weight: 800;
		letter-spacing: -0.02em;
		color: var(--ecochitas-ink);
	}

	.metric_label {
		margin: 0.3rem 0 0;
		font-size: 0.8rem;
		color: var(--ecochitas-ink);
		font-weight: 600;
	}

	.metric_helper {
		margin: 0.2rem 0 0;
		font-size: 0.72rem;
		color: var(--ecochitas-muted);
	}

	.mission_section {
		margin-bottom: 3rem;
		padding: 1rem 0;
	}

	.mission_title {
		font-size: 2rem;
		margin-bottom: 2rem;
	}

	.mission_grid {
		display: grid;
		gap: 2rem;
		align-items: center;
	}

	.mission_accordion {
		display: flex;
		flex-direction: column;
		gap: 0;
		border-top: 1px solid var(--ecochitas-border);
	}

	.accordion_item {
		background: transparent;
		border: none;
		border-bottom: 1px solid var(--ecochitas-border);
		padding: 1.2rem 0;
		text-align: left;
		cursor: pointer;
		color: var(--ecochitas-ink);
		display: flex;
		flex-direction: column;
		transition: background-color 0.2s ease;
	}

	.accordion_item:hover {
		background-color: var(--ecochitas-surface);
	}

	.accordion_header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		width: 100%;
	}

	.accordion_header h3 {
		margin: 0;
		font-size: 1.1rem;
		font-weight: 600;
		font-family: 'Plus Jakarta Sans', sans-serif;
	}

	.accordion_icon {
		font-size: 1.5rem;
		font-weight: 300;
		color: var(--ecochitas-muted);
		transition: color 0.2s ease;
	}

	.accordion_item.active .accordion_icon {
		color: var(--ecochitas-leaf);
	}

	.accordion_body {
		padding-top: 1rem;
		padding-right: 2rem;
	}

	.accordion_body p {
		margin: 0;
		color: var(--ecochitas-muted);
		font-size: 0.95rem;
		line-height: 1.6;
	}

	.page_footer {
		border-top: 1px solid var(--ecochitas-border);
		padding: 3rem 1rem;
		margin-top: 4rem;
	}

	.footer_content {
		max-width: 1000px;
		margin: 0 auto;
		display: flex;
		flex-direction: column;
		gap: 3rem;
	}

	.footer_left {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		align-items: flex-start;
	}

	.social_links {
		display: flex;
		gap: 0.8rem;
	}

	.social_links a {
		display: grid;
		place-items: center;
		width: 2.2rem;
		height: 2.2rem;
		border-radius: 50%;
		border: 1px solid var(--ecochitas-border);
		color: var(--ecochitas-muted);
		transition: all 0.2s ease;
	}

	.social_links a:hover {
		color: var(--ecochitas-ink);
		border-color: var(--ecochitas-muted);
		background: var(--ecochitas-surface);
	}

	.copyright {
		margin: 0;
		font-size: 0.8rem;
		color: var(--ecochitas-muted);
		font-weight: 600;
	}

	.footer_right {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.footer_nav_column {
		display: flex;
		flex-direction: column;
		gap: 0.8rem;
	}

	.footer_nav_column h4 {
		margin: 0 0 0.2rem;
		font-size: 0.9rem;
		color: var(--ecochitas-ink);
		font-weight: 600;
	}

	.footer_nav_column a {
		text-decoration: none;
		color: var(--ecochitas-muted);
		font-size: 0.85rem;
		transition: color 0.2s ease;
	}

	.footer_nav_column a:hover {
		color: var(--ecochitas-leaf);
	}

	@media (min-width: 900px) {
		.hero_title {
			font-size: 4rem;
			max-width: 800px;
			margin: 0 auto 1.5rem;
		}

		.hero_subtitle {
			font-size: 1.25rem;
		}

		.metrics_grid {
			grid-template-columns: repeat(4, minmax(0, 1fr));
			gap: 1rem;
			max-width: 800px;
			margin: 0 auto 2rem;
		}

		.mission_section {
			max-width: 800px;
			margin: 0 auto 3rem;
		}

		.mission_grid {
			grid-template-columns: 1fr;
		}

		.footer_content {
			flex-direction: row;
			justify-content: space-between;
			align-items: flex-end;
		}

		.footer_right {
			flex-direction: row;
			gap: 4rem;
		}

		.footer_left {
			flex-direction: row;
			align-items: center;
			gap: 2rem;
		}
	}
</style>
