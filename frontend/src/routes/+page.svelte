<script>
	import LineChart from '../components/LineChart.svelte';
	import { wrCurrData } from '$lib/stores/stores.js';

	let currData = {};

	wrCurrData.subscribe((value) => {
		currData = value;
	});
</script>

{#if Object.keys(currData).length > 0}
	{#key currData}
		<LineChart data={currData} />
	{/key}
{/if}
