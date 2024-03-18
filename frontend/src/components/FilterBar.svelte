<script>
	import '../app.pcss';
	import { GradientButton, Input, Label } from 'flowbite-svelte';
	import { wrCurrData } from '$lib/stores/stores';
	import { getDateString } from '$lib/utils/utils';

	let dateFrom = getDateString(-7);
	let dateTo = getDateString();
	let currData = {};

	function getDateFrom(event) {
		dateFrom = event.target.value;
	}

	function getDateTo(event) {
		dateTo = event.target.value;
	}

	async function getData() {
		const response = await fetch(
			`http://localhost:8080/api/dashboard/v1/data?dateFrom=${dateFrom}&dateTo=${dateTo}`
		);
		currData = await response.json();
		wrCurrData.set(currData);
	}
</script>

<section id="filter-bar">
	<div class="mb-6 flex flex-row rounded-lg border-2 border-green-300 p-3">
		<div class="flex-1 pr-2">
			<Label for="date-from">Date from</Label>
			<Input type="date" value={dateFrom} id="date-from" required on:input={getDateFrom} />
		</div>
		<div class="flex-1 pr-2">
			<Label for="date-to">Date to</Label>
			<Input type="date" value={dateTo} id="date-to" required on:input={getDateTo} />
		</div>
		<div class="flex-1 pt-5">
			<GradientButton outline color="greenToBlue" on:click={getData} >Get data</GradientButton>
		</div>
	</div>
</section>
