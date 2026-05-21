<script lang="ts">
	let {
		is_open = $bindable(false),
		title = 'Notificación',
		message = '',
		type = 'info',
		on_close = () => {}
	}: {
		is_open: boolean;
		title?: string;
		message?: string;
		type?: 'success' | 'error' | 'info';
		on_close?: () => void;
	} = $props();

	function handle_close() {
		is_open = false;
		on_close();
	}
</script>

{#if is_open}
	<div class="modal_backdrop" onclick={handle_close} role="presentation">
		<div
			class="modal_content"
			onclick={(e) => e.stopPropagation()}
			role="dialog"
			aria-modal="true"
		>
			<div class="modal_header">
				<h3 class="modal_title" class:success={type === 'success'} class:error={type === 'error'}>
					{#if type === 'success'}
						<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
					{:else if type === 'error'}
						<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
					{:else}
						<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
					{/if}
					{title}
				</h3>
				<button class="modal_close_btn" onclick={handle_close}>
					<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
				</button>
			</div>
			
			<div class="modal_body">
				<p>{message}</p>
			</div>
			
			<div class="modal_footer">
				<button class="modal_action_btn" class:success={type === 'success'} class:error={type === 'error'} onclick={handle_close}>
					Aceptar
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.modal_backdrop {
		position: fixed;
		top: 0;
		left: 0;
		width: 100vw;
		height: 100vh;
		background: rgba(0, 0, 0, 0.7);
		backdrop-filter: blur(4px);
		z-index: 1000;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1rem;
		box-sizing: border-box;
	}

	.modal_content {
		background: #0a0a0a;
		border: 1px solid #262626;
		border-radius: 16px;
		width: 100%;
		max-width: 420px;
		box-shadow: 0 20px 40px rgba(0, 0, 0, 0.9);
		overflow: hidden;
		animation: modal_slide_in 0.3s cubic-bezier(0.16, 1, 0.3, 1);
	}

	@keyframes modal_slide_in {
		from {
			opacity: 0;
			transform: translateY(20px) scale(0.95);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}

	.modal_header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1.25rem 1.5rem;
		border-bottom: 1px solid #262626;
	}

	.modal_title {
		margin: 0;
		font-size: 1.1rem;
		font-weight: 600;
		color: #ffffff;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.modal_title.success {
		color: #22c55e;
	}

	.modal_title.error {
		color: #ef4444;
	}

	.modal_close_btn {
		background: transparent;
		border: none;
		color: #a3a3a3;
		cursor: pointer;
		padding: 0.25rem;
		border-radius: 4px;
		transition: all 0.2s;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.modal_close_btn:hover {
		color: #ffffff;
		background: #262626;
	}

	.modal_body {
		padding: 1.5rem;
		color: #d4d4d8;
		font-size: 0.95rem;
		line-height: 1.5;
	}

	.modal_body p {
		margin: 0;
	}

	.modal_footer {
		padding: 1.25rem 1.5rem;
		border-top: 1px solid #262626;
		display: flex;
		justify-content: flex-end;
	}

	.modal_action_btn {
		background: #262626;
		color: #ffffff;
		border: 1px solid #404040;
		padding: 0.6rem 1.5rem;
		border-radius: 8px;
		font-weight: 600;
		font-size: 0.9rem;
		cursor: pointer;
		transition: all 0.2s;
	}

	.modal_action_btn:hover {
		background: #404040;
	}

	.modal_action_btn.success {
		background: #22c55e;
		border-color: #16a34a;
		color: #000000;
	}

	.modal_action_btn.success:hover {
		background: #16a34a;
	}

	.modal_action_btn.error {
		background: #ef4444;
		border-color: #dc2626;
		color: #ffffff;
	}

	.modal_action_btn.error:hover {
		background: #dc2626;
	}
</style>
