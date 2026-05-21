<script lang="ts">
	// ─── Infraction types ─────────────────────────────────────────────────────
	type InfractionType = { id: string; label: string; icon: string };

	const infraction_types: InfractionType[] = [
		{ id: 'basura_contenedores',  label: 'Basura fuera de contenedores', icon: 'M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6' },
		{ id: 'escombros',            label: 'Descarga ilegal de escombros',  icon: 'M5 7h14l-1.5 9h-11zM2 4h20M10 4V2h4v2' },
		{ id: 'quema',                label: 'Quema de basura',               icon: 'M12 22a7 7 0 0 0 7-7c0-2-1-3.9-3-5.5s-3.5-4-4-6.5c-.5 2.5-2 4.9-4 6.5C6 11.1 5 13 5 15a7 7 0 0 0 7 7z' },
		{ id: 'vertido_via',          label: 'Vertido en vía pública',        icon: 'M12 2.69l5.66 5.66a8 8 0 1 1-11.31 0z' },
		{ id: 'contaminacion_rios',   label: 'Contaminación de ríos',         icon: 'M2 12c2-4 5-6 8-6s6 2 8 6M2 18c2-4 5-6 8-6s6 2 8 6M2 6c2-4 5-6 8-6s6 2 8 6' },
		{ id: 'contenedor_desbordado',label: 'Contenedor desbordado',         icon: 'M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6M9 10v4M12 10v4M15 10v4' },
		{ id: 'otro',                 label: 'Otro',                          icon: 'M12 22c5.5 0 10-4.5 10-10S17.5 2 12 2 2 6.5 2 12s4.5 10 10 10zM12 8v4M12 16h.01' }
	];

	// ─── Multi-photo state ────────────────────────────────────────────────────
	type PhotoItem = { name: string; url: string };
	let photos = $state<PhotoItem[]>([]);
	let is_dragging = $state(false);
	const can_add_photo = $derived(photos.length < 3);

	function handle_file_add(files: FileList | null) {
		if (!files) return;
		const remaining = 3 - photos.length;
		const to_add = Array.from(files).slice(0, remaining);
		for (const file of to_add) {
			const url = URL.createObjectURL(file);
			photos = [...photos, { name: file.name, url }];
		}
	}
	function handle_photos(e: Event) { handle_file_add((e.target as HTMLInputElement).files); }
	function remove_photo(index: number) {
		URL.revokeObjectURL(photos[index].url);
		photos = photos.filter((_, i) => i !== index);
	}
	function on_dragover(e: DragEvent) { e.preventDefault(); is_dragging = true; }
	function on_dragleave() { is_dragging = false; }
	function on_drop(e: DragEvent) { e.preventDefault(); is_dragging = false; handle_file_add(e.dataTransfer?.files ?? null); }

	// ─── Form state ───────────────────────────────────────────────────────────
	let selected_type = $state('');
	let comment = $state('');
	let zone = $state('');
	let report_date = $state('');
	let report_time = $state('');
	let location_detected = $state(false);
	let report_sent = $state(false);

	function detect_location() { location_detected = true; }
	function send_report() {
		if (!selected_type) return;
		report_sent = true;
		setTimeout(() => {
			report_sent = false;
			selected_type = ''; comment = ''; zone = '';
			report_date = ''; report_time = '';
			location_detected = false; photos = [];
		}, 3500);
	}
	const can_send = $derived(!!selected_type && !report_sent);

	// ─── Report history ───────────────────────────────────────────────────────
	type ReportStatus = 'pendiente' | 'en_revision' | 'verificado' | 'atendido';
	type ReportHistory = { id: number; type: string; zone: string; date: string; status: ReportStatus; photos_count: number };

	const report_history: ReportHistory[] = [
		{ id: 1, type: 'Basura fuera de contenedores', zone: 'Av. América & Pando', date: '13 May 2026', status: 'atendido', photos_count: 2 },
		{ id: 2, type: 'Contenedor desbordado',        zone: 'Plaza Colón',         date: '10 May 2026', status: 'verificado', photos_count: 1 },
		{ id: 3, type: 'Vertido en vía pública',       zone: 'Calle Heroínas 300',  date: '05 May 2026', status: 'en_revision', photos_count: 3 }
	];
	const status_report_label: Record<ReportStatus, string> = {
		pendiente: 'Pendiente', en_revision: 'En Revisión', verificado: 'Verificado', atendido: 'Atendido'
	};

	// ─── Cochinitos ───────────────────────────────────────────────────────────
	type CochinStatus = 'multa_pendiente' | 'advertencia' | 'investigacion' | 'sancionado';
	type CochinBadge = 'responsable_identificado' | 'responsable_desconocido' | 'zona_critica' | 'reincidencia_alta';
	type Cochinito = {
		id: number; zone: string; infraction: string; reports: number;
		reincidence: boolean; status: CochinStatus; last_report: string;
		has_responsible: boolean; badges: CochinBadge[]; anonymous_desc?: string;
	};

	const cochinitos: Cochinito[] = [
		{ id: 1, zone: 'Av. América & Pando',       infraction: 'Basura fuera de contenedores',  reports: 7, reincidence: true,  status: 'multa_pendiente', last_report: '14 May 2026', has_responsible: true,  badges: ['responsable_identificado', 'reincidencia_alta'] },
		{ id: 2, zone: 'Plaza Colón',               infraction: 'Vertido en vía pública',         reports: 5, reincidence: true,  status: 'advertencia',    last_report: '12 May 2026', has_responsible: true,  badges: ['responsable_identificado', 'zona_critica'] },
		{ id: 3, zone: 'Av. Heroínas (cuadra 4–6)', infraction: 'Descarga ilegal de escombros',   reports: 4, reincidence: false, status: 'investigacion',  last_report: '10 May 2026', has_responsible: false, badges: ['responsable_desconocido', 'zona_critica'],      anonymous_desc: 'Zona con acumulación masiva de escombros sin propietario identificado. 4 reportes ciudadanos coincidentes.' },
		{ id: 4, zone: 'Parque Tunari (ingreso sur)', infraction: 'Quema de basura',              reports: 3, reincidence: false, status: 'multa_pendiente', last_report: '09 May 2026', has_responsible: false, badges: ['responsable_desconocido'],                      anonymous_desc: 'Quema nocturna recurrente. GAMC activó operativo de vigilancia en la zona.' },
		{ id: 5, zone: 'Río Rocha Sur',             infraction: 'Contaminación de ríos',          reports: 6, reincidence: true,  status: 'sancionado',     last_report: '07 May 2026', has_responsible: true,  badges: ['responsable_identificado', 'reincidencia_alta', 'zona_critica'] }
	];

	const status_label: Record<CochinStatus, string> = {
		multa_pendiente: 'Multa Pendiente', advertencia: 'Advertencia Emitida',
		investigacion: 'En Investigación', sancionado: 'Sancionado'
	};
	const badge_meta: Record<CochinBadge, { label: string; color: string; bg: string; dark_color: string; dark_bg: string }> = {
		responsable_identificado: { label: 'Responsable identificado', color: '#b91c1c', bg: '#fee2e2', dark_color: '#fca5a5', dark_bg: 'rgba(239,68,68,0.15)' },
		responsable_desconocido:  { label: 'Resp. desconocido',        color: '#92400e', bg: '#fef3c7', dark_color: '#fcd34d', dark_bg: 'rgba(245,158,11,0.15)' },
		zona_critica:             { label: 'Zona crítica',              color: '#6b21a8', bg: '#f3e8ff', dark_color: '#c4b5fd', dark_bg: 'rgba(124,58,237,0.15)' },
		reincidencia_alta:        { label: 'Reincidencia alta',         color: '#0c4a6e', bg: '#e0f2fe', dark_color: '#7dd3fc', dark_bg: 'rgba(14,116,144,0.15)' }
	};
</script>

<!-- ── Hero ── -->
<div class="report_hero">
	<div class="report_hero_eyebrow">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:13px;height:13px" aria-hidden="true"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0zM12 9v4M12 17h.01"/></svg>
		<span>Denuncias Ambientales · Cochabamba GAMC</span>
	</div>
	<h1 class="report_hero_title">Reportá una<br />infracción ambiental</h1>
	<p class="report_hero_sub">Toma una fotografía clara de la situación e infracción ambiental y repórtala para ayudar a mantener limpia tu zona y generar impacto ciudadano.</p>
	<div class="hero_chips">
		<span class="hero_chip">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" style="width:11px;height:11px" aria-hidden="true"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
			100% anónimo
		</span>
		<span class="hero_chip">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" style="width:11px;height:11px" aria-hidden="true"><path d="M12 21s-8-7.5-8-12a8 8 0 0 1 16 0c0 4.5-8 12-8 12z"/><circle cx="12" cy="9" r="2.5"/></svg>
			Zona Norte, Cbba.
		</span>
		<span class="hero_chip">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" style="width:11px;height:11px" aria-hidden="true"><polyline points="20 6 9 17 4 12"/></svg>
			Procesado en 24–48h
		</span>
	</div>
</div>

<!-- ── Formulario ── -->
<section class="panel">
	<div class="form_section_head">
		<svg viewBox="0 0 24 24" fill="none" stroke="var(--ecochitas-alert,#ef4444)" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:18px;height:18px;flex-shrink:0" aria-hidden="true"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0zM12 9v4M12 17h.01"/></svg>
		<h2 class="form_title">Nuevo Reporte de Contaminación</h2>
	</div>
	<p class="form_sub">Los campos con <span class="req">*</span> son necesarios para procesar tu denuncia.</p>

	<!-- ── Multi-photo upload ── -->
	<div class="photo_section">
		<div class="photo_section_head">
			<span class="field_label" style="margin:0">Evidencia fotográfica</span>
			<span class="photo_counter" class:photo_counter_full={photos.length >= 3}>{photos.length}/3 fotos</span>
		</div>

		{#if photos.length > 0}
			<div class="photos_preview">
				{#each photos as photo, i (photo.name + i)}
					<div class="photo_thumb">
						<img src={photo.url} alt="Evidencia {i + 1}" class="thumb_img" />
						<button class="thumb_remove" onclick={() => remove_photo(i)} aria-label="Eliminar foto">
							<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" style="width:9px;height:9px" aria-hidden="true"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
						</button>
					</div>
				{/each}
				{#if can_add_photo}
					<label class="photo_add_slot" for="photo_input_more">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:22px;height:22px;color:var(--ecochitas-muted)" aria-hidden="true"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
						<span class="photo_add_lbl">Agregar</span>
					</label>
					<input type="file" accept="image/*" id="photo_input_more" multiple class="upload_input" onchange={handle_photos} />
				{/if}
			</div>
		{:else}
			<div
				class="upload_zone"
				class:upload_drag={is_dragging}
				ondragover={on_dragover}
				ondragleave={on_dragleave}
				ondrop={on_drop}
				role="region"
				aria-label="Zona de carga de fotos"
			>
				<input type="file" accept="image/*" id="photo_input" multiple class="upload_input" onchange={handle_photos} />
				<label for="photo_input" class="upload_label">
					<div class="upload_icon_wrap">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.6" style="width:26px;height:26px" aria-hidden="true"><rect x="3" y="3" width="18" height="18" rx="2.5"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
					</div>
					<span class="upload_text">Arrastrá fotos aquí o tocá para subir</span>
					<span class="upload_hint">JPG, PNG o WEBP · máx. 10 MB · hasta 3 fotos</span>
				</label>
			</div>
		{/if}
	</div>

	<!-- ── Infraction type ── -->
	<p class="field_label">Tipo de Infracción <span class="req">*</span></p>
	<div class="infraction_grid">
		{#each infraction_types as inf (inf.id)}
			<button
				class="inf_btn"
				class:inf_btn_active={selected_type === inf.id}
				onclick={() => (selected_type = inf.id)}
				aria-pressed={selected_type === inf.id}
			>
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:14px;height:14px;flex-shrink:0" aria-hidden="true"><path d={inf.icon}/></svg>
				{inf.label}
			</button>
		{/each}
	</div>

	<!-- ── Meta fields ── -->
	<div class="meta_fields">
		<div class="field_group">
			<label class="field_label" for="zone_input">Zona / Referencia <span class="req">*</span></label>
			<input id="zone_input" class="field_input" type="text" placeholder="Ej: Av. América y Pando, Zona Norte" bind:value={zone} />
		</div>
		<div class="meta_row">
			<div class="field_group">
				<label class="field_label" for="date_input">Fecha</label>
				<input id="date_input" class="field_input" type="date" bind:value={report_date} />
			</div>
			<div class="field_group">
				<label class="field_label" for="time_input">Hora</label>
				<input id="time_input" class="field_input" type="time" bind:value={report_time} />
			</div>
		</div>
		<div class="field_group">
			<label class="field_label" for="comment_input">Comentario adicional</label>
			<textarea id="comment_input" class="field_textarea" placeholder="Describe brevemente lo que observaste..." rows={3} bind:value={comment}></textarea>
		</div>
	</div>

	<!-- ── GPS ── -->
	<button class="gps_btn" class:gps_btn_active={location_detected} onclick={detect_location} aria-pressed={location_detected}>
		{#if location_detected}
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.2" style="width:15px;height:15px;flex-shrink:0" aria-hidden="true"><polyline points="20 6 9 17 4 12"/></svg>
			Ubicación detectada · Zona Norte, Cochabamba
		{:else}
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:15px;height:15px;flex-shrink:0" aria-hidden="true"><path d="M12 21s-8-7.5-8-12a8 8 0 0 1 16 0c0 4.5-8 12-8 12z"/><circle cx="12" cy="9" r="2.5"/></svg>
			Detectar Mi Ubicación (GPS)
		{/if}
	</button>

	<!-- ── Submit ── -->
	<button class="submit_btn" class:submit_btn_sent={report_sent} disabled={!can_send} onclick={send_report}>
		{#if report_sent}
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.2" style="width:16px;height:16px;flex-shrink:0" aria-hidden="true"><polyline points="20 6 9 17 4 12"/></svg>
			Reporte enviado — ¡Gracias por tu aporte!
		{:else}
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:15px;height:15px;flex-shrink:0" aria-hidden="true"><line x1="22" y1="2" x2="11" y2="13"/><polygon points="22 2 15 22 11 13 2 9 22 2" fill="currentColor" stroke="none"/></svg>
			Enviar Reporte Anónimo
		{/if}
	</button>

	<div class="privacy_notice">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:13px;height:13px;flex-shrink:0;margin-top:1px" aria-hidden="true"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
		<span><strong>Privacidad garantizada:</strong> Los reportes son completamente anónimos y se usan exclusivamente para fines de cumplimiento municipal y mejora ambiental.</span>
	</div>
</section>

<!-- ── Mis reportes recientes ── -->
<section class="panel">
	<div class="form_section_head">
		<svg viewBox="0 0 24 24" fill="none" stroke="var(--ecochitas-sky,#3b82f6)" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:18px;height:18px;flex-shrink:0" aria-hidden="true"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8zM14 2v6h6M8 13h8M8 17h5"/></svg>
		<h2 class="form_title">Mis Reportes Recientes</h2>
	</div>
	<div class="history_list">
		{#each report_history as r (r.id)}
			<div class="history_item">
				<div class="history_left">
					<div class="history_icon">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:16px;height:16px" aria-hidden="true"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0zM12 9v4M12 17h.01"/></svg>
					</div>
					<div class="history_info">
						<strong class="history_type">{r.type}</strong>
						<span class="history_zone">
							<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:10px;height:10px;flex-shrink:0" aria-hidden="true"><path d="M12 21s-8-7.5-8-12a8 8 0 0 1 16 0c0 4.5-8 12-8 12z"/><circle cx="12" cy="9" r="2"/></svg>
							{r.zone}
						</span>
						<div class="history_meta">
							<span class="history_date">{r.date}</span>
							{#if r.photos_count > 0}
								<span class="history_photos_chip">
									<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:9px;height:9px" aria-hidden="true"><rect x="3" y="3" width="18" height="18" rx="2.5"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
									{r.photos_count}
								</span>
							{/if}
						</div>
					</div>
				</div>
				<span class="history_status history_status_{r.status}">{status_report_label[r.status]}</span>
			</div>
		{/each}
	</div>
</section>

<!-- ── Cochinitos de la semana ── -->
<section class="panel cochinitos_panel">
	<div class="coch_head_row">
		<div class="coch_head_icon">
			<svg viewBox="0 0 24 24" fill="none" stroke="#ef4444" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:18px;height:18px" aria-hidden="true">
				<path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0zM12 9v4M12 17h.01"/>
			</svg>
		</div>
		<div>
			<h2 class="coch_title">Zonas Críticas de la Semana</h2>
			<p class="coch_sub">Infracciones ambientales con mayor actividad ciudadana — con y sin responsable identificado.</p>
		</div>
	</div>

	<div class="coch_info_box">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:13px;height:13px;flex-shrink:0;margin-top:1px" aria-hidden="true"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
		<span>Esta sección agrupa zonas con mayor cantidad de reportes ciudadanos. Cuando hay responsable identificado la GAMC puede proceder con sanciones. Sin responsable, se activan operativos de limpieza y vigilancia.</span>
	</div>

	<div class="coch_warning">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:13px;height:13px;flex-shrink:0;margin-top:1px" aria-hidden="true"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0zM12 9v4M12 17h.01"/></svg>
		<span>Si existe reincidencia y se identifica al responsable, puede procederse con sanciones o multas municipales según ordenanzas vigentes de la GAMC.</span>
	</div>

	<div class="coch_list">
		{#each cochinitos as c (c.id)}
			<div class="coch_card" class:coch_card_anon={!c.has_responsible}>
				<div class="coch_avatar" class:coch_avatar_anon={!c.has_responsible}>
					{#if c.has_responsible}
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:18px;height:18px" aria-hidden="true">
							<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
							<circle cx="9" cy="7" r="4"/>
							<line x1="17" y1="8" x2="23" y2="14"/>
							<line x1="23" y1="8" x2="17" y2="14"/>
						</svg>
					{:else}
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" style="width:18px;height:18px" aria-hidden="true">
							<path d="M12 21s-8-7.5-8-12a8 8 0 0 1 16 0c0 4.5-8 12-8 12z"/>
							<circle cx="12" cy="9" r="2.5"/>
						</svg>
					{/if}
				</div>
				<div class="coch_info">
					<div class="coch_top_row">
						<strong class="coch_name">{c.has_responsible ? `Infractor #${c.id}` : 'Zona sin responsable'}</strong>
						<span class="coch_reports_badge" class:coch_badge_amber={!c.has_responsible}>{c.reports} Rep.</span>
					</div>
					<div class="coch_zone_row">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:10px;height:10px;flex-shrink:0" aria-hidden="true"><path d="M12 21s-8-7.5-8-12a8 8 0 0 1 16 0c0 4.5-8 12-8 12z"/><circle cx="12" cy="9" r="2"/></svg>
						{c.zone}
					</div>
					<div class="coch_infraction_label">{c.infraction}</div>
					{#if !c.has_responsible && c.anonymous_desc}
						<p class="coch_anon_desc">{c.anonymous_desc}</p>
					{/if}
					<div class="coch_badges_row">
						{#each c.badges as badge (badge)}
							<span class="coch_badge" style="background:{badge_meta[badge].bg};color:{badge_meta[badge].color}">{badge_meta[badge].label}</span>
						{/each}
					</div>
					<div class="coch_meta_row">
						<span class="coch_status coch_status_{c.status}">{status_label[c.status]}</span>
						{#if c.reincidence}
							<span class="coch_reincidence">
								<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" style="width:9px;height:9px;flex-shrink:0" aria-hidden="true"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 1 0 .49-4.5"/></svg>
								Reincidente
							</span>
						{/if}
						<span class="coch_date">
							<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.9" style="width:9px;height:9px;flex-shrink:0" aria-hidden="true"><rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
							{c.last_report}
						</span>
					</div>
				</div>
			</div>
		{/each}
	</div>
</section>

<style>
	h1, h2 { margin: 0; font-family: 'Sora', 'Plus Jakarta Sans', sans-serif; letter-spacing: -0.02em; }
	p { margin: 0; }

	/* ── Hero ── */
	.report_hero {
		background: linear-gradient(145deg, #7f1d1d, #991b1b, #b91c1c);
		border-radius: 1.5rem; padding: 1.75rem 1.5rem 1.5rem;
		display: flex; flex-direction: column; gap: 0.75rem;
	}
	.report_hero_eyebrow {
		display: inline-flex; align-items: center; gap: 0.4rem;
		font-size: 0.7rem; font-weight: 700; letter-spacing: 0.06em; text-transform: uppercase;
		color: oklch(1 0 0 / 0.65); background: oklch(1 0 0 / 0.1);
		border: 1px solid oklch(1 0 0 / 0.2); border-radius: 999px;
		padding: 0.25rem 0.75rem; width: fit-content;
	}
	.report_hero_title { font-size: 1.7rem; font-weight: 900; line-height: 1.1; color: white; }
	.report_hero_sub { font-size: 0.84rem; color: oklch(1 0 0 / 0.72); line-height: 1.55; max-width: 38ch; }
	.hero_chips { display: flex; flex-wrap: wrap; gap: 0.4rem; }
	.hero_chip {
		display: inline-flex; align-items: center; gap: 0.3rem;
		font-size: 0.68rem; font-weight: 700;
		background: oklch(1 0 0 / 0.12); color: oklch(1 0 0 / 0.85);
		border: 1px solid oklch(1 0 0 / 0.2); border-radius: 999px;
		padding: 0.2rem 0.6rem;
	}

	/* ── Form ── */
	.form_section_head { display: flex; align-items: center; gap: 0.5rem; }
	.form_title { font-size: 1.05rem; font-weight: 800; color: var(--ecochitas-ink); }
	.form_sub { font-size: 0.8rem; color: var(--ecochitas-muted); margin-top: 0.2rem; }
	.req { color: #ef4444; font-weight: 800; }

	/* ── Photo section ── */
	.photo_section { margin-top: 1.1rem; }
	.photo_section_head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 0.55rem; }
	.photo_counter { font-size: 0.72rem; font-weight: 700; color: var(--ecochitas-muted); font-family: 'Sora', sans-serif; }
	.photo_counter_full { color: #16a34a; }

	.photos_preview { display: flex; flex-wrap: wrap; gap: 0.55rem; }
	.photo_thumb {
		position: relative; width: 82px; height: 82px;
		border-radius: 1rem; overflow: hidden;
		border: 1.5px solid var(--ecochitas-border);
		box-shadow: 0 2px 8px rgba(0,0,0,0.1);
	}
	.thumb_img { width: 100%; height: 100%; object-fit: cover; display: block; }
	.thumb_remove {
		position: absolute; top: 4px; right: 4px;
		width: 20px; height: 20px; border-radius: 50%;
		background: rgba(0,0,0,0.65); border: none; color: white;
		cursor: pointer; display: flex; align-items: center; justify-content: center;
		padding: 0; transition: background 0.15s;
	}
	.thumb_remove:hover { background: rgba(239,68,68,0.85); }
	.photo_add_slot {
		width: 82px; height: 82px; border-radius: 1rem;
		border: 2px dashed var(--ecochitas-border);
		display: flex; flex-direction: column;
		align-items: center; justify-content: center; gap: 0.2rem;
		cursor: pointer; transition: border-color 0.15s, background 0.15s;
	}
	.photo_add_slot:hover { border-color: #ef4444; background: rgba(239,68,68,0.04); }
	.photo_add_lbl { font-size: 0.6rem; color: var(--ecochitas-muted); font-weight: 700; }

	.upload_zone {
		border: 2px dashed var(--ecochitas-border); border-radius: 1.1rem;
		transition: border-color 0.15s, background 0.15s;
	}
	.upload_zone:hover { border-color: #fca5a5; }
	.upload_drag { border-color: #ef4444; background: rgba(239,68,68,0.04); }
	.upload_input { display: none; }
	.upload_label {
		display: flex; flex-direction: column; align-items: center; justify-content: center;
		gap: 0.5rem; padding: 1.6rem 1rem; cursor: pointer; text-align: center;
	}
	.upload_icon_wrap {
		width: 52px; height: 52px; border-radius: 1rem;
		background: var(--ecochitas-surface); border: 1.5px solid var(--ecochitas-border);
		display: flex; align-items: center; justify-content: center;
		color: var(--ecochitas-muted);
	}
	.upload_text { font-size: 0.84rem; font-weight: 700; color: var(--ecochitas-ink); }
	.upload_hint { font-size: 0.71rem; color: var(--ecochitas-muted); }

	/* ── Infraction picker ── */
	.field_label {
		display: block; font-size: 0.82rem; font-weight: 700;
		color: var(--ecochitas-ink); margin-top: 1rem;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
	}
	.infraction_grid {
		display: grid; grid-template-columns: 1fr 1fr;
		gap: 0.45rem; margin-top: 0.55rem;
	}
	.inf_btn {
		display: flex; align-items: center; gap: 0.45rem;
		padding: 0.62rem 0.75rem; border-radius: 0.85rem;
		border: 1.5px solid var(--ecochitas-border);
		background: var(--ecochitas-surface);
		color: var(--ecochitas-ink);
		font-size: 0.73rem; font-weight: 600; cursor: pointer;
		text-align: left; line-height: 1.3;
		transition: border-color 0.15s, background 0.15s, color 0.12s;
		font-family: 'Plus Jakarta Sans', sans-serif;
	}
	.inf_btn:hover { border-color: #fca5a5; background: #fff1f2; color: #b91c1c; }
	:root[data-theme='dark'] .inf_btn:hover { background: rgba(239,68,68,0.1); border-color: rgba(239,68,68,0.4); color: #fca5a5; }
	.inf_btn_active { border-color: #ef4444; background: #fee2e2; color: #b91c1c; }
	:root[data-theme='dark'] .inf_btn_active { border-color: #ef4444; background: rgba(239,68,68,0.14); color: #fca5a5; }

	/* ── Meta fields ── */
	.meta_fields { display: flex; flex-direction: column; gap: 0.75rem; margin-top: 0.85rem; }
	.field_group { display: flex; flex-direction: column; gap: 0.35rem; }
	.meta_row { display: grid; grid-template-columns: 1fr 1fr; gap: 0.65rem; }
	.field_input, .field_textarea {
		padding: 0.65rem 0.85rem; border-radius: 0.85rem;
		border: 1.5px solid var(--ecochitas-border);
		background: var(--ecochitas-surface);
		color: var(--ecochitas-ink); font-size: 0.85rem;
		font-family: 'Plus Jakarta Sans', sans-serif;
		transition: border-color 0.15s; box-sizing: border-box; width: 100%;
	}
	.field_textarea { resize: vertical; }
	.field_input:focus, .field_textarea:focus { outline: none; border-color: #ef4444; }

	/* ── GPS & Submit ── */
	.gps_btn {
		display: flex; align-items: center; justify-content: center; gap: 0.55rem;
		width: 100%; padding: 0.82rem 1rem; margin-top: 0.75rem;
		border-radius: 1rem; border: none;
		background: var(--ecochitas-sky, #3b82f6); color: white;
		font-size: 0.86rem; font-weight: 800; cursor: pointer;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		transition: background 0.15s, transform 0.12s;
		box-shadow: 0 4px 14px rgba(59,130,246,0.3);
	}
	.gps_btn:hover { background: #2563eb; transform: translateY(-1px); }
	.gps_btn_active { background: #16a34a; box-shadow: 0 4px 14px rgba(22,163,74,0.3); }
	.gps_btn_active:hover { background: #15803d; }
	.submit_btn {
		display: flex; align-items: center; justify-content: center; gap: 0.55rem;
		width: 100%; padding: 0.9rem 1rem; margin-top: 0.6rem;
		border-radius: 1rem; border: none;
		background: #ef4444; color: white;
		font-size: 0.88rem; font-weight: 800; cursor: pointer;
		font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
		transition: background 0.15s, transform 0.12s, opacity 0.15s;
		box-shadow: 0 4px 16px rgba(239,68,68,0.35);
	}
	.submit_btn:hover:not(:disabled) { background: #dc2626; transform: translateY(-1px); }
	.submit_btn:disabled { opacity: 0.42; cursor: not-allowed; transform: none; box-shadow: none; }
	.submit_btn_sent { background: #16a34a !important; box-shadow: 0 4px 14px rgba(22,163,74,0.3) !important; opacity: 1 !important; cursor: default !important; }

	/* ── Privacy ── */
	.privacy_notice {
		display: flex; align-items: flex-start; gap: 0.55rem;
		padding: 0.72rem 0.9rem; margin-top: 0.75rem;
		border-radius: 0.85rem; background: #eff6ff;
		border: 1px solid #bfdbfe; font-size: 0.75rem; color: #1d4ed8; line-height: 1.5;
	}
	:root[data-theme='dark'] .privacy_notice { background: rgba(59,130,246,0.08); border-color: rgba(59,130,246,0.25); color: #93c5fd; }
	.privacy_notice strong { color: #1e40af; }
	:root[data-theme='dark'] .privacy_notice strong { color: #60a5fa; }

	/* ── Report history ── */
	.history_list { display: flex; flex-direction: column; gap: 0.5rem; margin-top: 0.85rem; }
	.history_item {
		display: flex; align-items: center; justify-content: space-between; gap: 0.75rem;
		padding: 0.8rem 0.95rem; border: 1.5px solid var(--ecochitas-border);
		border-radius: 1rem; background: var(--ecochitas-surface);
		transition: box-shadow 0.15s;
	}
	.history_item:hover { box-shadow: 0 3px 12px rgba(0,0,0,0.05); }
	.history_left { display: flex; align-items: flex-start; gap: 0.65rem; min-width: 0; flex: 1; }
	.history_icon {
		flex-shrink: 0; width: 36px; height: 36px; border-radius: 0.75rem;
		background: rgba(59,130,246,0.08); border: 1.5px solid rgba(59,130,246,0.18);
		display: flex; align-items: center; justify-content: center; color: var(--ecochitas-sky);
	}
	.history_info { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 0.18rem; }
	.history_type { font-size: 0.8rem; font-weight: 800; color: var(--ecochitas-ink); font-family: 'Sora', sans-serif; line-height: 1.2; }
	.history_zone { display: flex; align-items: center; gap: 0.28rem; font-size: 0.71rem; color: var(--ecochitas-muted); }
	.history_meta { display: flex; align-items: center; gap: 0.5rem; }
	.history_date { font-size: 0.63rem; color: var(--ecochitas-muted); }
	.history_photos_chip { display: inline-flex; align-items: center; gap: 0.2rem; font-size: 0.63rem; color: var(--ecochitas-muted); }
	.history_status {
		flex-shrink: 0; font-size: 0.63rem; font-weight: 800;
		padding: 0.22rem 0.6rem; border-radius: 999px;
		white-space: nowrap; font-family: 'Sora', sans-serif;
	}
	.history_status_pendiente { background: #fef9c3; color: #854d0e; border: 1px solid #fde68a; }
	:root[data-theme='dark'] .history_status_pendiente { background: rgba(234,179,8,0.12); color: #fbbf24; border-color: rgba(234,179,8,0.3); }
	.history_status_en_revision { background: #eff6ff; color: #1e40af; border: 1px solid #bfdbfe; }
	:root[data-theme='dark'] .history_status_en_revision { background: rgba(59,130,246,0.12); color: #93c5fd; border-color: rgba(59,130,246,0.3); }
	.history_status_verificado { background: #f5f3ff; color: #4c1d95; border: 1px solid #ddd6fe; }
	:root[data-theme='dark'] .history_status_verificado { background: rgba(124,58,237,0.12); color: #c4b5fd; border-color: rgba(124,58,237,0.3); }
	.history_status_atendido { background: #dcfce7; color: #166534; border: 1px solid #bbf7d0; }
	:root[data-theme='dark'] .history_status_atendido { background: rgba(34,197,94,0.12); color: #4ade80; border-color: rgba(34,197,94,0.3); }

	/* ── Cochinitos ── */
	.cochinitos_panel { border-color: #fecaca; background: color-mix(in srgb, #fff5f5 80%, var(--ecochitas-surface)); }
	:root[data-theme='dark'] .cochinitos_panel { border-color: rgba(239,68,68,0.22); background: rgba(239,68,68,0.04); }
	.coch_head_row { display: flex; align-items: flex-start; gap: 0.65rem; }
	.coch_head_icon { flex-shrink: 0; width: 38px; height: 38px; border-radius: 0.75rem; background: #fee2e2; border: 1.5px solid #fecaca; display: flex; align-items: center; justify-content: center; }
	:root[data-theme='dark'] .coch_head_icon { background: rgba(239,68,68,0.12); border-color: rgba(239,68,68,0.25); }
	.coch_title { font-size: 1.02rem; font-weight: 800; color: #b91c1c; }
	:root[data-theme='dark'] .coch_title { color: #fca5a5; }
	.coch_sub { font-size: 0.78rem; color: var(--ecochitas-muted); margin-top: 0.12rem; line-height: 1.4; }
	.coch_info_box {
		display: flex; align-items: flex-start; gap: 0.5rem;
		padding: 0.7rem 0.9rem; margin-top: 0.85rem;
		border-radius: 0.85rem; background: #eff6ff;
		border: 1px solid #bfdbfe; font-size: 0.74rem; color: #1d4ed8; line-height: 1.5;
	}
	:root[data-theme='dark'] .coch_info_box { background: rgba(59,130,246,0.08); border-color: rgba(59,130,246,0.22); color: #93c5fd; }
	.coch_warning {
		display: flex; align-items: flex-start; gap: 0.5rem;
		padding: 0.68rem 0.9rem; margin-top: 0.55rem;
		border-radius: 0.85rem; background: #fffbeb;
		border: 1px solid #fde68a; font-size: 0.74rem; color: #92400e; line-height: 1.5;
	}
	:root[data-theme='dark'] .coch_warning { background: rgba(245,158,11,0.08); border-color: rgba(245,158,11,0.28); color: #fcd34d; }
	.coch_list { display: flex; flex-direction: column; gap: 0.6rem; margin-top: 0.85rem; }
	.coch_card {
		display: flex; align-items: flex-start; gap: 0.75rem;
		padding: 0.9rem 1rem; background: white;
		border: 1.5px solid #fecaca; border-radius: 1.15rem;
		transition: box-shadow 0.15s;
	}
	:root[data-theme='dark'] .coch_card { background: rgba(255,255,255,0.04); border-color: rgba(239,68,68,0.18); }
	.coch_card:hover { box-shadow: 0 4px 16px rgba(239,68,68,0.1); }
	.coch_card_anon { border-color: #fde68a !important; }
	:root[data-theme='dark'] .coch_card_anon { border-color: rgba(245,158,11,0.25) !important; }
	.coch_avatar {
		flex-shrink: 0; width: 42px; height: 42px;
		background: #fee2e2; border: 1.5px solid #fecaca;
		border-radius: 0.85rem; display: flex; align-items: center; justify-content: center; color: #ef4444;
	}
	:root[data-theme='dark'] .coch_avatar { background: rgba(239,68,68,0.12); border-color: rgba(239,68,68,0.25); }
	.coch_avatar_anon { background: #fffbeb !important; border-color: #fde68a !important; color: #d97706 !important; }
	:root[data-theme='dark'] .coch_avatar_anon { background: rgba(245,158,11,0.12) !important; border-color: rgba(245,158,11,0.28) !important; color: #fbbf24 !important; }
	.coch_info { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 0.28rem; }
	.coch_top_row { display: flex; align-items: center; justify-content: space-between; gap: 0.5rem; }
	.coch_name { font-size: 0.84rem; font-weight: 800; color: var(--ecochitas-ink); font-family: 'Sora', sans-serif; }
	.coch_reports_badge { flex-shrink: 0; font-size: 0.63rem; font-weight: 800; color: white; background: #ef4444; padding: 0.18rem 0.5rem; border-radius: 999px; font-family: 'Sora', sans-serif; }
	.coch_badge_amber { background: #d97706 !important; }
	.coch_zone_row { display: flex; align-items: center; gap: 0.3rem; font-size: 0.73rem; color: var(--ecochitas-muted); }
	.coch_infraction_label { font-size: 0.74rem; font-weight: 600; color: var(--ecochitas-ink); }
	.coch_anon_desc { font-size: 0.71rem; color: var(--ecochitas-muted); line-height: 1.5; margin: 0.05rem 0; }
	.coch_badges_row { display: flex; flex-wrap: wrap; gap: 0.3rem; margin: 0.25rem 0; }
	.coch_badge { display: inline-flex; align-items: center; font-size: 0.61rem; font-weight: 700; padding: 0.15rem 0.48rem; border-radius: 999px; font-family: 'Sora', sans-serif; }
	:root[data-theme='dark'] .coch_badge { opacity: 0.88; }
	.coch_meta_row { display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap; }
	.coch_status { display: inline-flex; align-items: center; font-size: 0.62rem; font-weight: 700; padding: 0.18rem 0.5rem; border-radius: 999px; }
	.coch_status_multa_pendiente { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; }
	:root[data-theme='dark'] .coch_status_multa_pendiente { background: rgba(239,68,68,0.1); color: #fca5a5; border-color: rgba(239,68,68,0.28); }
	.coch_status_advertencia { background: #fffbeb; color: #92400e; border: 1px solid #fde68a; }
	:root[data-theme='dark'] .coch_status_advertencia { background: rgba(245,158,11,0.1); color: #fcd34d; border-color: rgba(245,158,11,0.28); }
	.coch_status_investigacion { background: #eff6ff; color: #1e40af; border: 1px solid #bfdbfe; }
	:root[data-theme='dark'] .coch_status_investigacion { background: rgba(59,130,246,0.1); color: #93c5fd; border-color: rgba(59,130,246,0.28); }
	.coch_status_sancionado { background: #f5f3ff; color: #4c1d95; border: 1px solid #ddd6fe; }
	:root[data-theme='dark'] .coch_status_sancionado { background: rgba(124,58,237,0.1); color: #c4b5fd; border-color: rgba(124,58,237,0.28); }
	.coch_reincidence { display: inline-flex; align-items: center; gap: 0.28rem; font-size: 0.62rem; font-weight: 700; color: #dc2626; }
	:root[data-theme='dark'] .coch_reincidence { color: #f87171; }
	.coch_date { display: inline-flex; align-items: center; gap: 0.25rem; font-size: 0.62rem; color: var(--ecochitas-muted); }
</style>
