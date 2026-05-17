<script lang="ts">
	export type EcochitasIconName =
		| 'brand'
		| 'home'
		| 'map'
		| 'recycling'
		| 'rewards'
		| 'profile'
		| 'upload'
		| 'schedule'
		| 'alert'
		| 'truck'
		| 'moon'
		| 'sun'
		| 'x'
		| 'linkedin'
		| 'github'
		| 'youtube';

	let { name, size = 20 } = $props<{ name: EcochitasIconName; size?: number }>();
</script>

<svg
	viewBox="0 0 24 24"
	fill="none"
	stroke="currentColor"
	stroke-linecap="round"
	stroke-linejoin="round"
	stroke-width="1.9"
	style={`width:${size}px;height:${size}px`}
	aria-hidden="true"
>
	{#if name === 'brand'}
		<path d="M19 5c-6.2 0-10 3.6-10 8.8 0 2 0.7 3.8 2 5.2" />
		<path d="M5 19c2.2-3.9 5.5-6.8 10-8.6" />
		<path d="M10.8 18c3.8 0 7.2-2.4 8.2-6.1" />
	{:else if name === 'home'}
		<path d="M4 11.5L12 5l8 6.5" />
		<path d="M6 10.8V20h12v-9.2" />
		<path d="M10 20v-4h4v4" />
	{:else if name === 'map'}
		<path d="M3.5 6.5L9.5 4l5 2.5L20.5 4v13.5L14.5 20l-5-2.5-6 2.5z" />
		<path d="M9.5 4v13.5" />
		<path d="M14.5 6.5V20" />
	{:else if name === 'recycling'}
		<path d="M3 12a9 9 0 0 1 15.2-6.4" />
		<path d="M18.2 5.6V2.4h-3.2" />
		<path d="M21 12a9 9 0 0 1-15.2 6.4" />
		<path d="M5.8 18.4v3.2H9" />
	{:else if name === 'rewards'}
		<path d="M12 3l2.7 5.5 6.1.9-4.4 4.3 1 6.1L12 17l-5.4 2.8 1-6.1L3.2 9.4l6.1-.9z" />
	{:else if name === 'profile'}
		<circle cx="12" cy="8" r="3.2" />
		<path d="M5.3 19.5c1.5-3.2 3.8-4.8 6.7-4.8s5.3 1.6 6.7 4.8" />
	{:else if name === 'upload'}
		<path d="M12 16V6" />
		<path d="M8.5 9.5L12 6l3.5 3.5" />
		<path d="M4 14.5V19h16v-4.5" />
	{:else if name === 'schedule'}
		<rect x="4" y="5" width="16" height="15" rx="2.5" />
		<path d="M8 3.8V7" />
		<path d="M16 3.8V7" />
		<path d="M4 10.2h16" />
		<path d="M8.2 14h3.6" />
	{:else if name === 'alert'}
		<path d="M12 4.3l8.3 14.4H3.7L12 4.3z" />
		<path d="M12 9.2v4.8" />
		<circle cx="12" cy="16.8" r="0.8" fill="currentColor" stroke="none" />
	{:else if name === 'truck'}
		<path d="M3.5 8.2h9.7v6.1H3.5z" />
		<path d="M13.2 10h3.7l2 2.5v1.8h-5.7z" />
		<circle cx="7.3" cy="16.3" r="1.5" />
		<circle cx="15.6" cy="16.3" r="1.5" />
	{:else if name === 'moon'}
		<path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
	{:else if name === 'sun'}
		<circle cx="12" cy="12" r="5" />
		<line x1="12" y1="1" x2="12" y2="3" />
		<line x1="12" y1="21" x2="12" y2="23" />
		<line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
		<line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
		<line x1="1" y1="12" x2="3" y2="12" />
		<line x1="21" y1="12" x2="23" y2="12" />
		<line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
		<line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
	{:else if name === 'x'}
		<line x1="18" y1="6" x2="6" y2="18"></line>
		<line x1="6" y1="6" x2="18" y2="18"></line>
	{:else if name === 'linkedin'}
		<path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z" />
		<rect x="2" y="9" width="4" height="12" />
		<circle cx="4" cy="4" r="2" />
	{:else if name === 'github'}
		<path
			d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"
		/>
	{:else if name === 'youtube'}
		<path
			d="M22.54 6.42a2.78 2.78 0 0 0-1.94-2C18.88 4 12 4 12 4s-6.88 0-8.6.42a2.78 2.78 0 0 0-1.94 2C1 8.13 1 12 1 12s0 3.87.46 5.58a2.78 2.78 0 0 0 1.94 2C5.12 20 12 20 12 20s6.88 0 8.6-.42a2.78 2.78 0 0 0 1.94-2C23 15.87 23 12 23 12s0-3.87-.46-5.58z"
		/>
		<polygon points="9.75 15.02 15.5 12 9.75 8.98 9.75 15.02" />
	{/if}
</svg>
