<script>
	import '../app.pcss';
	import { GradientButton, Input, Label, Select } from 'flowbite-svelte';
	import { wrCurrData } from '$lib/stores/stores';
	import { getDateString } from '$lib/utils/utils';

	let dateFrom = getDateString(-7);
	let dateTo = getDateString();
	let currData = {};

	export let data;

	/** @type {import('flowbite-svelte').SelectOptionType<string>[]}*/
	let currSymbols = [];
	let selected = 'USD';

	data.currencies.forEach((symbol) => {
		currSymbols.push({ value: symbol, name: symbol });
	});

	function getDateFrom(event) {
		dateFrom = event.target.value;
	}

	function getDateTo(event) {
		dateTo = event.target.value;
	}

	async function getData() {
		const response = await fetch(
			`http://localhost:8080/api/dashboard/v2/data?dateFrom=${dateFrom}&dateTo=${dateTo}&currency=${selected}`
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
		<div class="flex-1 pr-2">
			<Label for="currency">Currency</Label>
			<Select items={currSymbols} bind:value={selected} id="currency"></Select>
		</div>
		<div class="flex-1 pt-5">
			<GradientButton outline color="greenToBlue" on:click={getData}>Get data</GradientButton>
		</div>
	</div>
</section>
