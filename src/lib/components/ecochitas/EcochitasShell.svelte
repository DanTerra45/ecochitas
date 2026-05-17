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
		{ href: '/map', label: 'Mapa', icon_name: 'map' },
		{ href: '/recycling', label: 'EcoHub', icon_name: 'recycling' },
		{ href: '/profile', label: 'Perfil', icon_name: 'profile' },
		{ href: '/reportar', label: 'Reportar', icon_name: 'alert' }
	];

	const route_labels_by_path: Record<string, string> = {
		'/': 'Resumen operativo',
		'/map': 'Monitoreo en tiempo real',
		'/recycling': 'EcoHub · Reciclaje',
		'/profile': 'Mi Perfil',
		'/reportar': 'Reportar Infracción'
	};

	let { children } = $props();

	let is_dark_mode = $state(true);
	let show_notifications = $state(false);

	type Notif = { id: number; type: 'success' | 'alert' | 'info' | 'warn'; title: string; desc: string; time: string; read: boolean };
	const notifications: Notif[] = [
		{ id: 1, type: 'success', title: 'Validación completada', desc: 'Tu entrega del 15 May fue validada. +120 EcoPoints', time: 'hace 2h', read: false },
		{ id: 2, type: 'alert', title: 'Reporte procesado', desc: 'Tu denuncia en Av. América fue derivada a la GAMC.', time: 'hace 5h', read: false },
		{ id: 3, type: 'info', title: 'Nueva recompensa disponible', desc: 'Café La Casona: 20% off disponible para canjear.', time: 'ayer', read: true },
		{ id: 4, type: 'warn', title: 'Nivel en riesgo', desc: 'Sin actividad en 12 días. Tu nivel podría reducirse.', time: 'hace 2d', read: true }
	];
	const unread_count = $derived(notifications.filter((n) => !n.read).length);

	$effect(() => {
		const saved_theme = localStorage.getItem('theme');
		if (saved_theme) {
			is_dark_mode = saved_theme === 'dark';
		} else {
			is_dark_mode = window.matchMedia('(prefers-color-scheme: dark)').matches;
		}
		document.documentElement.dataset.theme = is_dark_mode ? 'dark' : 'light';
	});

	function toggle_theme() {
		is_dark_mode = !is_dark_mode;
		const new_theme = is_dark_mode ? 'dark' : 'light';
		document.documentElement.dataset.theme = new_theme;
		localStorage.setItem('theme', new_theme);
	}

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

			<div class="nav_and_actions">
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

				<div class="notif_wrap">
					<button
						class="notif_btn"
						onclick={() => (show_notifications = !show_notifications)}
						aria-label="Notificaciones"
						aria-expanded={show_notifications}
					>
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:18px;height:18px" aria-hidden="true"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>
						{#if unread_count > 0}
							<span class="notif_dot">{unread_count}</span>
						{/if}
					</button>
					{#if show_notifications}
						<div class="notif_panel">
							<div class="notif_panel_head">
								<span class="notif_head_title">Notificaciones</span>
								{#if unread_count > 0}<span class="notif_head_count">{unread_count} nuevas</span>{/if}
							</div>
							<div class="notif_list">
								{#each notifications as n (n.id)}
									<div class="notif_item" class:notif_item_unread={!n.read}>
										<div class="notif_icon notif_icon_{n.type}">
											{#if n.type === 'success'}
												<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.2" style="width:12px;height:12px" aria-hidden="true"><polyline points="20 6 9 17 4 12"/></svg>
											{:else if n.type === 'alert'}
												<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:12px;height:12px" aria-hidden="true"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0zM12 9v4M12 17h.01"/></svg>
											{:else if n.type === 'warn'}
												<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:12px;height:12px" aria-hidden="true"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
											{:else}
												<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:12px;height:12px" aria-hidden="true"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
											{/if}
										</div>
										<div class="notif_body">
											<strong class="notif_title">{n.title}</strong>
											<p class="notif_desc">{n.desc}</p>
											<span class="notif_time">{n.time}</span>
										</div>
										{#if !n.read}<div class="notif_unread_dot"></div>{/if}
									</div>
								{/each}
							</div>
						</div>
					{/if}
				</div>

				<button
					class="theme_toggle_btn"
					onclick={toggle_theme}
					aria-label="Cambiar tema oscuro/claro"
				>
					<EcochitasIcon name={is_dark_mode ? 'sun' : 'moon'} size={18} />
				</button>
			</div>
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
	:global(html),
	:global(body) {
		margin: 0;
		padding: 0;
		background-color: var(--ecochitas-background, #000000);
		color: var(--ecochitas-ink, #ffffff);
		transition:
			background-color 0.3s ease,
			color 0.3s ease;
	}

	:root {
		/* Light Mode Variables */
		--ecochitas-background: #f4f4f5;
		--ecochitas-surface: #ffffff;
		--ecochitas-border: rgba(0, 0, 0, 0.08);
		--ecochitas-ink: #18181b;
		--ecochitas-muted: #71717a;
		--ecochitas-leaf: #16a34a;
		--ecochitas-sky: #3b82f6;
		--ecochitas-alert: #ef4444;
		--ecochitas-nav-bg: rgba(255, 255, 255, 0.85);
		--ecochitas-nav-border: rgba(0, 0, 0, 0.05);
		--ecochitas-shadow-soft: 0 10px 30px -10px rgba(0, 0, 0, 0.1);
		--ecochitas-shadow-float: 0 20px 40px -10px rgba(0, 0, 0, 0.15);
		--ecochitas-hero-bg: radial-gradient(
			circle at 50% 30%,
			rgba(22, 163, 74, 0.1) 0%,
			transparent 60%
		);
	}

	:root[data-theme='dark'] {
		/* Dark Mode Variables */
		--ecochitas-background: #000000;
		--ecochitas-surface: rgba(255, 255, 255, 0.03);
		--ecochitas-border: rgba(255, 255, 255, 0.1);
		--ecochitas-ink: #ffffff;
		--ecochitas-muted: #a1a1aa;
		--ecochitas-leaf: #55ff77;
		--ecochitas-sky: #60a5fa;
		--ecochitas-alert: #f87171;
		--ecochitas-nav-bg: rgba(30, 30, 30, 0.6);
		--ecochitas-nav-border: rgba(255, 255, 255, 0.08);
		--ecochitas-shadow-soft: 0 10px 30px -10px rgba(0, 0, 0, 0.5);
		--ecochitas-shadow-float: 0 20px 40px -10px rgba(0, 0, 0, 0.8);
		--ecochitas-hero-bg: radial-gradient(
			circle at 50% 30%,
			rgba(85, 255, 119, 0.08) 0%,
			transparent 60%
		);
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
		padding: 6rem 1rem 6.8rem;
		box-sizing: border-box;
	}

	.top_header {
		position: fixed;
		top: 1rem;
		left: 50%;
		transform: translateX(-50%);
		width: calc(100% - 2rem);
		max-width: 440px;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.8rem;
		padding: 0.75rem 1rem;
		border-radius: 2rem;
		background: var(--ecochitas-nav-bg);
		backdrop-filter: blur(16px);
		-webkit-backdrop-filter: blur(16px);
		border: 1px solid var(--ecochitas-nav-border);
		box-shadow: var(--ecochitas-shadow-float);
		z-index: 50;
	}

	.brand_group {
		display: flex;
		align-items: center;
		gap: 0.7rem;
		min-width: 0;
	}

	.brand_icon {
		width: 2.2rem;
		height: 2.2rem;
		display: grid;
		place-items: center;
		border-radius: 0.8rem;
		background: linear-gradient(135deg, var(--ecochitas-leaf), #22c55e);
		color: #000;
	}

	.eyebrow_text {
		margin: 0;
		font-size: 0.68rem;
		color: var(--ecochitas-muted);
		letter-spacing: 0.08em;
		text-transform: uppercase;
		display: none;
	}

	h1 {
		margin: 0;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		letter-spacing: -0.02em;
		font-size: 1rem;
		line-height: 1.15;
		color: var(--ecochitas-ink);
	}

	.nav_and_actions {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.desktop_nav {
		display: none;
	}

	.theme_toggle_btn {
		background: transparent;
		border: none;
		color: var(--ecochitas-muted);
		cursor: pointer;
		display: grid;
		place-items: center;
		padding: 0.5rem;
		border-radius: 50%;
		transition:
			color 0.2s ease,
			background-color 0.2s ease;
	}

	.theme_toggle_btn:hover {
		color: var(--ecochitas-ink);
		background-color: var(--ecochitas-border);
	}

	.content_grid {
		display: grid;
		gap: 0.85rem;
	}

	:global(.panel) {
		background: var(--ecochitas-surface);
		border: 1px solid var(--ecochitas-border);
		border-radius: 1.5rem;
		padding: 1.5rem;
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
		grid-template-columns: repeat(4, 1fr);
		align-items: end;
		gap: 0.35rem;
		padding: 0.5rem 0.45rem;
		border-radius: 1.6rem;
		background: var(--ecochitas-nav-bg);
		backdrop-filter: blur(16px);
		-webkit-backdrop-filter: blur(16px);
		border: 1px solid var(--ecochitas-nav-border);
		box-shadow: var(--ecochitas-shadow-float);
		z-index: 50;
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
		background: rgba(22, 163, 74, 0.1);
	}

	:root[data-theme='dark'] .mobile_nav_item_active {
		background: rgba(85, 255, 119, 0.1);
	}

	/* ── Notification bell ── */
	.notif_wrap { position: relative; }
	.notif_btn {
		position: relative;
		background: transparent; border: none;
		color: var(--ecochitas-muted); cursor: pointer;
		display: grid; place-items: center;
		padding: 0.5rem; border-radius: 50%;
		transition: color 0.2s ease, background-color 0.2s ease;
	}
	.notif_btn:hover { color: var(--ecochitas-ink); background-color: var(--ecochitas-border); }
	.notif_dot {
		position: absolute; top: 2px; right: 2px;
		min-width: 16px; height: 16px; padding: 0 3px;
		background: #ef4444; color: white;
		font-size: 0.58rem; font-weight: 800;
		border-radius: 999px; display: flex; align-items: center; justify-content: center;
		font-family: 'Sora', sans-serif;
		border: 1.5px solid var(--ecochitas-nav-bg);
		pointer-events: none;
	}
	.notif_panel {
		position: absolute; top: calc(100% + 0.75rem); right: -0.5rem;
		width: min(310px, calc(100vw - 2rem));
		background: var(--ecochitas-nav-bg);
		backdrop-filter: blur(20px); -webkit-backdrop-filter: blur(20px);
		border: 1px solid var(--ecochitas-nav-border);
		border-radius: 1.2rem;
		box-shadow: var(--ecochitas-shadow-float);
		overflow: hidden; z-index: 100;
	}
	.notif_panel_head {
		display: flex; align-items: center; justify-content: space-between;
		padding: 0.85rem 1rem 0.65rem;
		border-bottom: 1px solid var(--ecochitas-nav-border);
	}
	.notif_head_title { font-size: 0.85rem; font-weight: 800; color: var(--ecochitas-ink); font-family: 'Sora', sans-serif; }
	.notif_head_count { font-size: 0.68rem; font-weight: 800; background: #ef4444; color: white; padding: 0.15rem 0.5rem; border-radius: 999px; font-family: 'Sora', sans-serif; }
	.notif_list { display: flex; flex-direction: column; max-height: 280px; overflow-y: auto; }
	.notif_item {
		display: flex; align-items: flex-start; gap: 0.65rem;
		padding: 0.72rem 1rem;
		border-bottom: 1px solid var(--ecochitas-nav-border);
		position: relative;
		transition: background 0.15s;
	}
	.notif_item:last-child { border-bottom: none; }
	.notif_item_unread { background: color-mix(in srgb, var(--ecochitas-leaf) 5%, transparent); }
	.notif_icon {
		flex-shrink: 0; width: 28px; height: 28px;
		border-radius: 50%; display: flex; align-items: center; justify-content: center;
	}
	.notif_icon_success { background: #dcfce7; color: #16a34a; }
	:root[data-theme='dark'] .notif_icon_success { background: rgba(34,197,94,0.15); }
	.notif_icon_alert { background: #fef2f2; color: #ef4444; }
	:root[data-theme='dark'] .notif_icon_alert { background: rgba(239,68,68,0.15); }
	.notif_icon_warn { background: #fffbeb; color: #d97706; }
	:root[data-theme='dark'] .notif_icon_warn { background: rgba(217,119,6,0.15); }
	.notif_icon_info { background: #eff6ff; color: #3b82f6; }
	:root[data-theme='dark'] .notif_icon_info { background: rgba(59,130,246,0.15); }
	.notif_body { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 0.1rem; }
	.notif_title { font-size: 0.76rem; font-weight: 800; color: var(--ecochitas-ink); font-family: 'Sora', sans-serif; line-height: 1.2; }
	.notif_desc { font-size: 0.7rem; color: var(--ecochitas-muted); line-height: 1.4; margin: 0; }
	.notif_time { font-size: 0.62rem; color: var(--ecochitas-muted); }
	.notif_unread_dot { flex-shrink: 0; width: 7px; height: 7px; border-radius: 50%; background: var(--ecochitas-leaf); margin-top: 5px; }

	@media (min-width: 1024px) {
		.ecochitas_page {
			padding: 0;
		}

		.app_shell {
			max-width: 1320px;
			min-height: 100vh;
			border-radius: 0;
			padding: 8rem 1.3rem 1.4rem;
			background: transparent;
			border: none;
			box-shadow: none;
		}

		.top_header {
			top: 2rem;
			width: auto;
			max-width: none;
			padding: 0.5rem 1rem;
		}

		.brand_group h1 {
			display: none;
		}

		.brand_group::after {
			content: 'EcoChitas';
			font-family: 'Sora', sans-serif;
			font-weight: 700;
			font-size: 1.2rem;
			color: var(--ecochitas-ink);
			margin-left: 0.5rem;
			margin-right: 2rem;
		}

		.nav_and_actions {
			gap: 1rem;
		}

		.desktop_nav {
			display: flex;
			align-items: center;
			gap: 0.2rem;
			padding: 0;
			border-radius: 0;
			background: transparent;
			border: none;
		}

		.desktop_nav_item {
			display: inline-flex;
			align-items: center;
			gap: 0.35rem;
			padding: 0.5rem 1rem;
			border-radius: 2rem;
			text-decoration: none;
			font-size: 0.85rem;
			font-weight: 600;
			color: var(--ecochitas-muted);
			transition:
				background-color 150ms ease,
				color 150ms ease;
		}

		.desktop_nav_item:hover {
			color: var(--ecochitas-ink);
		}

		.desktop_nav_item_active {
			color: var(--ecochitas-leaf);
			background: rgba(22, 163, 74, 0.1);
		}

		:root[data-theme='dark'] .desktop_nav_item_active {
			background: rgba(85, 255, 119, 0.1);
		}

		.content_grid {
			min-height: calc(100vh - 10rem);
		}

		.mobile_nav {
			display: none;
		}
	}
</style>
