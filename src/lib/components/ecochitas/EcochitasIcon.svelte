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
		| 'truck';

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
	{/if}
</svg>
