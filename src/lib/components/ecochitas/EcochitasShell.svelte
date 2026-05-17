<script lang="ts">
	import { resolve } from '$app/paths';
	import type { Pathname } from '$app/types';
	import { page } from '$app/state';
	import EcochitasIcon, {
		type EcochitasIconName
	} from '$lib/components/ecochitas/EcochitasIcon.svelte';

	type Navigation_item = {
		href: Pathname;
		label: string;
		icon_name: EcochitasIconName;
	};

	const navigation_items: Navigation_item[] = [
		{ href: '/', label: 'Inicio', icon_name: 'home' },
		{ href: '/map', label: 'Mapa', icon_name: 'map' },
		{ href: '/recycling', label: 'Reciclaje', icon_name: 'recycling' },
		{ href: '/rewards', label: 'Puntos', icon_name: 'rewards' },
		{ href: '/profile', label: 'Perfil', icon_name: 'profile' }
	];

	const route_labels_by_path: Record<string, string> = {
		'/': 'Resumen operativo',
		'/map': 'Monitoreo en tiempo real',
		'/recycling': 'Guia de reciclaje',
		'/rewards': 'Bonificaciones por zona',
		'/profile': 'Cuenta ciudadana'
	};

	let { children } = $props();

	function is_active_route(route_path: string): boolean {
		if (route_path === '/') {
			return page.url.pathname === '/';
		}
		return page.url.pathname.startsWith(route_path);
	}

	const current_route_label = $derived(
		route_labels_by_path[page.url.pathname] ?? route_labels_by_path['/']
	);
</script>

<main class="ecochitas_page">
	<div class="app_shell">
		<header class="top_header">
			<div class="brand_group">
				<div class="brand_icon" aria-hidden="true">
					<EcochitasIcon name="brand" size={25} />
				</div>
				<div>
					<p class="eyebrow_text">EcoChitas · Cochabamba</p>
					<h1>{current_route_label}</h1>
				</div>
			</div>

			<nav class="desktop_nav" aria-label="Secciones de EcoChitas">
				{#each navigation_items as navigation_item (navigation_item.href)}
					<a
						href={resolve(navigation_item.href)}
						class="desktop_nav_item"
						class:desktop_nav_item_active={is_active_route(navigation_item.href)}
					>
						<EcochitasIcon name={navigation_item.icon_name} size={16} />
						<span>{navigation_item.label}</span>
					</a>
				{/each}
			</nav>
		</header>

		<section class="content_grid">
			{@render children()}
		</section>

		<nav class="mobile_nav" aria-label="Navegacion principal">
			{#each navigation_items as navigation_item (navigation_item.href)}
				<a
					href={resolve(navigation_item.href)}
					class="mobile_nav_item"
					class:mobile_nav_item_active={is_active_route(navigation_item.href)}
				>
					<EcochitasIcon name={navigation_item.icon_name} size={19} />
					<span>{navigation_item.label}</span>
				</a>
			{/each}
		</nav>
	</div>
</main>

<style>
	:root {
		--ecochitas-background: linear-gradient(
			160deg,
			oklch(0.93 0.06 155) 0%,
			oklch(0.97 0.03 90) 55%,
			oklch(0.94 0.05 215) 100%
		);
		--ecochitas-surface: oklch(1 0 0 / 0.8);
		--ecochitas-border: oklch(0.88 0.02 140 / 0.9);
		--ecochitas-ink: oklch(0.23 0.03 160);
		--ecochitas-muted: oklch(0.52 0.02 160);
		--ecochitas-leaf: oklch(0.55 0.13 155);
		--ecochitas-sky: oklch(0.62 0.12 220);
		--ecochitas-alert: oklch(0.66 0.18 28);
		--ecochitas-shadow-soft: 0 20px 40px -24px oklch(0.32 0.04 165 / 0.3);
		--ecochitas-shadow-float: 0 30px 60px -35px oklch(0.35 0.05 160 / 0.35);
	}

	.ecochitas_page {
		min-height: 100vh;
		background: var(--ecochitas-background);
		padding: 0;
		box-sizing: border-box;
		color: var(--ecochitas-ink);
		font-family: 'Plus Jakarta Sans', 'Manrope', 'Segoe UI', sans-serif;
	}

	.app_shell {
		position: relative;
		width: 100%;
		max-width: 460px;
		min-height: 100vh;
		margin: 0 auto;
		padding: 1.1rem 1rem 6.8rem;
		box-sizing: border-box;
	}

	.top_header {
		display: flex;
		flex-direction: column;
		gap: 0.8rem;
		margin-bottom: 1rem;
	}

	.brand_group {
		display: flex;
		align-items: center;
		gap: 0.7rem;
		min-width: 0;
	}

	.brand_icon {
		width: 2.75rem;
		height: 2.75rem;
		display: grid;
		place-items: center;
		border-radius: 1rem;
		background: linear-gradient(135deg, oklch(0.88 0.07 155), oklch(0.92 0.05 210));
		box-shadow: var(--ecochitas-shadow-soft);
		color: var(--ecochitas-ink);
	}

	.eyebrow_text {
		margin: 0;
		font-size: 0.68rem;
		color: var(--ecochitas-muted);
		letter-spacing: 0.08em;
		text-transform: uppercase;
	}

	h1 {
		margin: 0;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		letter-spacing: -0.02em;
		font-size: 1.1rem;
		line-height: 1.15;
	}

	.desktop_nav {
		display: none;
	}

	.content_grid {
		display: grid;
		gap: 0.85rem;
	}

	:global(.panel) {
		background: var(--ecochitas-surface);
		backdrop-filter: blur(14px);
		-webkit-backdrop-filter: blur(14px);
		border: 1px solid var(--ecochitas-border);
		border-radius: 1.5rem;
		padding: 1rem;
		box-shadow: var(--ecochitas-shadow-soft);
	}

	:global(.muted_text) {
		margin: 0.32rem 0 0;
		font-size: 0.82rem;
		color: var(--ecochitas-muted);
	}

	:global(.eyebrow_text) {
		margin: 0;
		font-size: 0.68rem;
		color: var(--ecochitas-muted);
		letter-spacing: 0.08em;
		text-transform: uppercase;
	}

	.mobile_nav {
		position: fixed;
		left: 0.9rem;
		right: 0.9rem;
		bottom: 0.9rem;
		display: grid;
		grid-template-columns: repeat(5, 1fr);
		align-items: end;
		gap: 0.35rem;
		padding: 0.5rem 0.45rem;
		border-radius: 1.6rem;
		background: oklch(1 0 0 / 0.68);
		backdrop-filter: blur(16px);
		-webkit-backdrop-filter: blur(16px);
		border: 1px solid oklch(0.9 0.02 160 / 0.9);
		box-shadow: var(--ecochitas-shadow-float);
	}

	.mobile_nav_item {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.25rem;
		padding: 0.45rem 0.3rem;
		border-radius: 0.9rem;
		font-size: 0.64rem;
		font-weight: 700;
		color: var(--ecochitas-muted);
		text-decoration: none;
	}

	.mobile_nav_item_active {
		color: var(--ecochitas-leaf);
		background: oklch(0.93 0.07 155 / 0.58);
	}

	@media (min-width: 1024px) {
		.ecochitas_page {
			padding: 1.2rem;
		}

		.app_shell {
			max-width: 1320px;
			min-height: calc(100vh - 2.4rem);
			border-radius: 2rem;
			padding: 1.3rem 1.3rem 1.4rem;
			background: oklch(1 0 0 / 0.5);
			border: 1px solid oklch(0.88 0.02 150 / 0.9);
			box-shadow: var(--ecochitas-shadow-float);
		}

		.top_header {
			margin-bottom: 1.2rem;
			flex-direction: row;
			align-items: center;
			justify-content: space-between;
		}

		h1 {
			font-size: 1.35rem;
		}

		.desktop_nav {
			display: flex;
			align-items: center;
			gap: 0.45rem;
			padding: 0.3rem;
			border-radius: 1rem;
			background: oklch(1 0 0 / 0.64);
			border: 1px solid var(--ecochitas-border);
		}

		.desktop_nav_item {
			display: inline-flex;
			align-items: center;
			gap: 0.35rem;
			padding: 0.5rem 0.7rem;
			border-radius: 0.75rem;
			text-decoration: none;
			font-size: 0.78rem;
			font-weight: 700;
			color: var(--ecochitas-muted);
			transition:
				background-color 120ms ease,
				color 120ms ease;
		}

		.desktop_nav_item:hover {
			background: oklch(0.95 0.02 150);
			color: var(--ecochitas-ink);
		}

		.desktop_nav_item_active {
			color: var(--ecochitas-leaf);
			background: oklch(0.93 0.07 155 / 0.5);
		}

		.content_grid {
			min-height: calc(100vh - 10rem);
		}

		.mobile_nav {
			display: none;
		}
	}
</style>
