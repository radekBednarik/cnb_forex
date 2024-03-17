<script>
	import LineChart from '../components/LineChart.svelte';
	import { wrCurrData } from '$lib/stores/stores.js';

	let currData = '';

	wrCurrData.subscribe((value) => {
		if (value.length > 0) {
			currData = JSON.parse(value);
		}
	});
</script>

<p>
	{JSON.stringify(currData)}
</p>
