<script>
	import '../app.pcss';
	import { Hr, Heading, Span, Input, Label, Button } from 'flowbite-svelte';
	/** @type {import('./$types').LayoutData} */
	export let data;
</script>

<header>
	<Heading tag="h1"><Span gradient>Daily forex data</Span></Heading>
	<Heading tag="h4">
		provided by <a
			href="https://www.cnb.cz"
			target="_blank"
			referrerpolicy="no-referrer-when-downgrade">Czech national bank</a
		>
	</Heading>
	<Hr />
</header>

<section id="filter-bar">
	<form>
		<div>
			<div>
				<Label for="date-from">Date from</Label>
				<Input type="date" id="date-from" required />
			</div>
			<div>
				<Label for="date-to">Date to</Label>
				<Input type="date" id="date-to" required />
			</div>
			<div>
				<Button color="green">Filter data</Button>
			</div>
		</div>
	</form>
</section>
<slot />

<style>
	#filter-bar {
		flex: auto;
		flex-direction: row;
	}
</style>