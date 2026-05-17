<script lang="ts">
	import EcochitasIcon from '$lib/components/ecochitas/EcochitasIcon.svelte';
</script>

<section class="profile_grid">
	<article class="panel">
		<div class="section_title_row">
			<EcochitasIcon name="profile" size={18} />
			<h2>Datos del ciudadano</h2>
		</div>
		<div class="field_grid">
			<label>
				Nombre completo
				<input value="Maria Quispe" />
			</label>
			<label>
				Celular
				<input value="71456789" />
			</label>
			<label>
				Zona
				<input value="Zona Sarco" />
			</label>
			<label>
				Direccion de referencia
				<input value="Calle 12 de Octubre #123" />
			</label>
		</div>
	</article>

	<article class="panel">
		<div class="section_title_row">
			<EcochitasIcon name="rewards" size={18} />
			<h2>Resumen de puntos</h2>
		</div>
		<div class="stats_grid">
			<div class="stat_card">
				<p>Puntos personales</p>
				<strong>184</strong>
			</div>
			<div class="stat_card">
				<p>Puntos zona</p>
				<strong>1320</strong>
			</div>
			<div class="stat_card">
				<p>Aportes validados</p>
				<strong>9</strong>
			</div>
			<div class="stat_card">
				<p>Penalizaciones</p>
				<strong>1</strong>
			</div>
		</div>
	</article>

	<article class="panel">
		<div class="section_title_row">
			<EcochitasIcon name="upload" size={18} />
			<h2>Denuncias y evidencia</h2>
		</div>
		<p class="muted_text">
			Sube fotos cuando exista basura fuera de los puntos autorizados. El sistema agrupa
			reincidencias por zona y notifica al encargado municipal.
		</p>
		<div class="upload_actions">
			<button type="button">Subir foto de evidencia</button>
			<button type="button">Ver historial de denuncias</button>
		</div>
	</article>

	<article class="panel">
		<div class="section_title_row">
			<EcochitasIcon name="alert" size={18} />
			<h2>Cochinitos de la semana</h2>
		</div>
		<ul class="rule_list">
			<li>Persona reportada 3 veces por vecinos distintos pasa a revision del encargado.</li>
			<li>La publicacion requiere validacion manual para evitar denuncias falsas.</li>
			<li>Se muestra solo informacion minima y evidencia de lugar/hora.</li>
		</ul>
	</article>
</section>

<style>
	h2 {
		margin: 0;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		font-size: 1rem;
		letter-spacing: -0.015em;
	}

	.profile_grid {
		display: grid;
		gap: 0.85rem;
	}

	.section_title_row {
		display: flex;
		align-items: center;
		gap: 0.45rem;
	}

	.field_grid {
		margin-top: 0.7rem;
		display: grid;
		gap: 0.58rem;
	}

	label {
		display: block;
		font-size: 0.8rem;
		font-weight: 600;
	}

	input {
		width: 100%;
		margin-top: 0.3rem;
		padding: 0.63rem 0.72rem;
		font-size: 0.9rem;
		border-radius: 0.8rem;
		border: 1px solid oklch(0.84 0.02 140);
		background: oklch(0.99 0.005 120);
		color: var(--ecochitas-ink);
		box-sizing: border-box;
	}

	.stats_grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0.6rem;
		margin-top: 0.7rem;
	}

	.stat_card {
		border: 1px solid var(--ecochitas-border);
		border-radius: 0.85rem;
		padding: 0.65rem;
		background: oklch(1 0 0 / 0.68);
	}

	.stat_card p {
		margin: 0;
		font-size: 0.74rem;
		color: var(--ecochitas-muted);
	}

	.stat_card strong {
		display: block;
		margin-top: 0.26rem;
		font-size: 1.1rem;
		letter-spacing: -0.02em;
	}

	.upload_actions {
		margin-top: 0.75rem;
		display: grid;
		gap: 0.5rem;
	}

	button {
		padding: 0.6rem 0.68rem;
		border-radius: 0.9rem;
		border: 1px solid oklch(0.86 0.03 145);
		background: linear-gradient(135deg, oklch(0.98 0.01 150), oklch(0.95 0.03 220));
		color: var(--ecochitas-ink);
		font-size: 0.8rem;
		font-weight: 700;
		cursor: pointer;
	}

	.rule_list {
		margin: 0.75rem 0 0;
		padding-left: 1.1rem;
		display: grid;
		gap: 0.45rem;
		font-size: 0.85rem;
	}

	@media (min-width: 900px) {
		.profile_grid {
			grid-template-columns: 1fr 1fr;
		}
	}
</style>
