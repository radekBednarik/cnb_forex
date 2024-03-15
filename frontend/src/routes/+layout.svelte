<script>
	import '../app.pcss';
	import { Hr, Heading, Span, Input, Label, Button } from 'flowbite-svelte';
	/** @type {import('./$types').LayoutData} */
	export let data;
</script>

<div class="mx-12 my-12">
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
		<form method="post" enctype="application/x-www-form-urlencoded">
			<div class="mb-6 flex flex-row">
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
	<slot />
</div>
