<script lang="ts">
	import EcochitasShell from '$lib/components/ecochitas/EcochitasShell.svelte';

	let { children } = $props();
</script>

<EcochitasShell>
	{@render children()}
</EcochitasShell>
