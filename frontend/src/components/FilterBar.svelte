<script>
	import '../app.pcss';
	import { Input, Label } from 'flowbite-svelte';
</script>

<section id="filter-bar">
	<form method="get">
		<div class="mb-6 flex flex-row rounded-lg border-2 border-green-300 p-3">
			<div class="w-1/3 pr-2">
				<Label for="date-from">Date from</Label>
				<Input type="date" id="date-from" required />
			</div>
			<div class="w-1/3 pr-2">
				<Label for="date-to">Date to</Label>
				<Input type="date" id="date-to" required />
			</div>
			<div class="w-1/3 pt-5">
				<Input type="submit" value="Get data" color="green" />
			</div>
		</div>
	</form>
</section>
