<script>
	import LineChart from '../components/LineChart.svelte';
	import { wrCurrData } from '$lib/stores/stores.js';

	let currData = {};

	wrCurrData.subscribe((value) => {
		currData = value;
	});
</script>

{#key currData}
	<LineChart data={currData} />
{/key}
