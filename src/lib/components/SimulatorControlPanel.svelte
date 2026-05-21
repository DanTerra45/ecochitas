<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	type SimulationStatus = 'idle' | 'running' | 'paused';
	type SelectorOption = { value: string; label: string };

	interface SimulationPayload {
		routeId: string;
		zone: string;
		truckCount: number;
		speedMultiplier: number;
		updateIntervalMs: number;
		snapToRoute: boolean;
	}

	let {
		routeId = $bindable('route-1'),
		zone = $bindable('north-zone'),
		truckCount = $bindable(5),
		speedMultiplier = $bindable(1),
		updateIntervalMs = $bindable(1000),
		snapToRoute = $bindable(true),
		status = 'idle',
		activeTrucks = 0,
		messagesSent = 0,
		routeOptions = [
			{ value: 'route-1', label: 'Route 1 (Downtown)' },
			{ value: 'route-2', label: 'Route 2 (Suburbs)' },
			{ value: 'route-3', label: 'Route 3 (Highway)' }
		],
		zoneOptions = [
			{ value: 'north-zone', label: 'North Zone' },
			{ value: 'south-zone', label: 'South Zone' },
			{ value: 'east-zone', label: 'East Zone' },
			{ value: 'west-zone', label: 'West Zone' }
		]
	}: {
		routeId: string;
		zone: string;
		truckCount: number;
		speedMultiplier: number;
		updateIntervalMs: number;
		snapToRoute: boolean;
		status: SimulationStatus;
		activeTrucks: number;
		messagesSent: number;
		routeOptions?: SelectorOption[];
		zoneOptions?: SelectorOption[];
	} = $props();

	const dispatch = createEventDispatcher<{
		start: SimulationPayload;
		pause: SimulationPayload;
		resume: SimulationPayload;
		stop: SimulationPayload;
		change: SimulationPayload;
	}>();

	function getPayload(): SimulationPayload {
		return {
			routeId,
			zone,
			truckCount,
			speedMultiplier,
			updateIntervalMs,
			snapToRoute
		};
	}

	function handleStart() {
		dispatch('start', getPayload());
	}

	function handlePause() {
		dispatch('pause', getPayload());
	}

	function handleResume() {
		dispatch('resume', getPayload());
	}

	function handleStop() {
		dispatch('stop', getPayload());
	}

	function handleChange() {
		dispatch('change', getPayload());
	}
</script>

<div class="panel">
	<header class="panel-header">
		<h2 class="title">Simulation Control</h2>
		<p class="subtitle">Generate realtime truck movement for demo/testing</p>
	</header>

	<div class="panel-body">
		<!-- Sección: Scenario -->
		<section class="section" aria-labelledby="scenario-heading">
			<h3 id="scenario-heading" class="section-title">Scenario</h3>
			<div class="field-group">
				<div class="field">
					<label for="routeId">Route</label>
					<select id="routeId" bind:value={routeId} onchange={handleChange}>
						{#each routeOptions as routeOption (routeOption.value)}
							<option value={routeOption.value}>{routeOption.label}</option>
						{/each}
					</select>
				</div>

				<div class="field">
					<label for="zone">Zone</label>
					<select id="zone" bind:value={zone} onchange={handleChange}>
						{#each zoneOptions as zoneOption (zoneOption.value)}
							<option value={zoneOption.value}>{zoneOption.label}</option>
						{/each}
					</select>
				</div>

				<div class="field">
					<label for="truckCount">Trucks ({truckCount})</label>
					<input
						id="truckCount"
						type="number"
						min="1"
						max="20"
						bind:value={truckCount}
						oninput={handleChange}
					/>
				</div>
			</div>
		</section>

		<!-- Sección: Motion -->
		<section class="section" aria-labelledby="motion-heading">
			<h3 id="motion-heading" class="section-title">Motion</h3>
			<div class="field-group">
				<div class="field">
					<label for="speedMultiplier">Speed multiplier ({speedMultiplier}x)</label>
					<input
						id="speedMultiplier"
						type="range"
						min="0.25"
						max="4"
						step="0.25"
						bind:value={speedMultiplier}
						onchange={handleChange}
					/>
				</div>

				<div class="field">
					<label for="updateIntervalMs">Update interval ({updateIntervalMs}ms)</label>
					<input
						id="updateIntervalMs"
						type="range"
						min="300"
						max="5000"
						step="100"
						bind:value={updateIntervalMs}
						onchange={handleChange}
					/>
				</div>

				<div class="field checkbox-field">
					<label>
						<input type="checkbox" bind:checked={snapToRoute} onchange={handleChange} />
						Snap to route
					</label>
				</div>
			</div>
		</section>

		<!-- Sección: Actions -->
		<section class="section actions-section" aria-label="Simulation Actions">
			<button
				class="btn btn-primary"
				onclick={handleStart}
				disabled={status === 'running' || status === 'paused'}
				aria-label="Start simulation"
			>
				Start simulation
			</button>

			<button
				class="btn btn-secondary"
				onclick={handlePause}
				disabled={status !== 'running'}
				aria-label="Pause simulation"
			>
				Pause
			</button>

			<button
				class="btn"
				onclick={handleResume}
				disabled={status !== 'paused'}
				aria-label="Resume simulation"
			>
				Resume
			</button>

			<button
				class="btn btn-danger"
				onclick={handleStop}
				disabled={status === 'idle'}
				aria-label="Stop and reset simulation"
			>
				Stop & reset
			</button>
		</section>

		<!-- Sección: Status Chips -->
		<section class="status-section" aria-live="polite">
			<div class="status-chip chip-{status}">
				Status: <strong>{status}</strong>
			</div>
			<div class="status-chip">
				Active Trucks: <strong>{activeTrucks}</strong>
			</div>
			<div class="status-chip">
				Messages Sent: <strong>{messagesSent}</strong>
			</div>
		</section>
	</div>
</div>

<style>
	/* Integración con las variables del tema, forzando estética Dark/Neon */
	.panel {
		--bg-surface: #0a0a0a;
		--border-color: #262626;
		--text-ink: #ffffff;
		--text-muted: #a3a3a3;
		--color-leaf: #22c55e;
		--color-sky: #3b82f6;
		--color-alert: #ef4444;

		/* Variables derivadas dinámicas usando color-mix para estados hover/chips */
		--color-leaf-hover: #16a34a;
		--color-sky-hover: #2563eb;
		--color-alert-hover: #dc2626;

		background-color: var(--bg-surface);
		color: var(--text-ink);
		border: 1px solid var(--border-color);
		border-radius: 16px;
		padding: 1.5rem;
		font-family:
			system-ui,
			-apple-system,
			sans-serif;
		box-shadow:
			0 10px 25px rgba(0, 0, 0, 0.8),
			0 4px 6px rgba(0, 0, 0, 0.5);
		width: 100%;
		max-width: 650px;
		margin: 0 auto;
		box-sizing: border-box;
		backdrop-filter: blur(8px);
	}

	.panel * {
		box-sizing: border-box;
	}

	.panel-header {
		margin-bottom: 1.5rem;
		padding-bottom: 1rem;
		border-bottom: 1px solid var(--border-color);
	}

	.title {
		margin: 0 0 0.25rem 0;
		font-size: 1.25rem;
		font-weight: 600;
	}

	.subtitle {
		margin: 0;
		font-size: 0.875rem;
		color: var(--text-muted);
	}

	.panel-body {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.section-title {
		font-size: 0.75rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: var(--color-leaf);
		margin: 0 0 1rem 0;
		font-weight: 700;
	}

	.field-group {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	/* Grid Responsive para campos */
	@media (min-width: 480px) {
		.field-group {
			flex-direction: row;
			flex-wrap: wrap;
		}
		.field {
			flex: 1 1 calc(50% - 0.5rem);
		}
	}

	.field {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.checkbox-field {
		flex-direction: row;
		align-items: center;
		justify-content: flex-start;
	}

	.checkbox-field label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		cursor: pointer;
	}

	label {
		font-size: 0.875rem;
		font-weight: 500;
		color: var(--text-ink);
	}

	input[type='number'],
	select {
		width: 100%;
		padding: 0.625rem;
		border: 1px solid var(--border-color);
		border-radius: 8px;
		background-color: #000000;
		color: var(--text-ink);
		font-size: 0.875rem;
		transition:
			border-color 0.2s ease,
			box-shadow 0.2s ease;
	}

	input[type='number']:focus,
	select:focus {
		outline: none;
		border-color: var(--color-leaf);
		box-shadow: 0 0 0 3px color-mix(in srgb, var(--color-leaf) 20%, transparent);
	}

	input[type='range'] {
		width: 100%;
		accent-color: var(--color-leaf);
		cursor: pointer;
		height: 1.5rem;
	}

	input[type='checkbox'] {
		width: 1.125rem;
		height: 1.125rem;
		accent-color: var(--color-leaf);
		cursor: pointer;
	}

	/* Controles de Acción (Botones) */
	.actions-section {
		display: grid;
		grid-template-columns: 1fr;
		gap: 0.75rem;
		padding-top: 1rem;
		border-top: 1px solid var(--border-color);
	}

	@media (min-width: 480px) {
		.actions-section {
			grid-template-columns: repeat(2, 1fr);
		}
	}

	@media (min-width: 640px) {
		.actions-section {
			grid-template-columns: repeat(4, 1fr);
		}
	}

	.btn {
		padding: 0.625rem 1rem;
		border-radius: 8px;
		border: 1px solid var(--border-color);
		background-color: #000000;
		color: var(--text-ink);
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s ease;
		text-align: center;
	}

	.btn:focus-visible {
		outline: 2px solid var(--color-leaf);
		outline-offset: 2px;
	}

	.btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		filter: grayscale(0.5);
	}

	.btn:not(:disabled):hover {
		background-color: #262626;
	}

	.btn-primary {
		background-color: var(--color-leaf);
		border-color: var(--color-leaf-hover);
		color: #000000;
	}

	.btn-primary:not(:disabled):hover {
		background-color: var(--color-leaf-hover);
		border-color: var(--color-leaf-hover);
	}

	.btn-secondary {
		background-color: var(--color-sky);
		border-color: var(--color-sky-hover);
		color: #000000;
	}

	.btn-secondary:not(:disabled):hover {
		background-color: var(--color-sky-hover);
		border-color: var(--color-sky-hover);
	}

	.btn-danger {
		background-color: var(--color-alert);
		border-color: var(--color-alert-hover);
		color: #ffffff;
	}

	.btn-danger:not(:disabled):hover {
		background-color: var(--color-alert-hover);
		border-color: var(--color-alert-hover);
	}

	/* Status Chips */
	.status-section {
		display: flex;
		flex-wrap: wrap;
		gap: 0.75rem;
		padding-top: 1rem;
		border-top: 1px solid var(--border-color);
	}

	.status-chip {
		display: inline-flex;
		align-items: center;
		padding: 0.375rem 0.75rem;
		border-radius: 9999px;
		font-size: 0.75rem;
		background-color: #000000;
		color: var(--text-ink);
		border: 1px solid var(--border-color);
	}

	.status-chip strong {
		margin-left: 0.375rem;
	}

	.chip-running {
		background-color: color-mix(in srgb, var(--color-leaf) 15%, transparent);
		color: var(--color-leaf);
		border-color: var(--color-leaf);
	}

	.chip-paused {
		background-color: color-mix(in srgb, var(--color-sky) 15%, transparent);
		color: var(--color-sky);
		border-color: var(--color-sky);
	}

	.chip-idle {
		background-color: #171717;
	}
</style>
