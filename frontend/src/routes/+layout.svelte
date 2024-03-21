<script>
	import '../app.pcss';
	import Header from '../components/Header.svelte';
	import FilterBar from '../components/FilterBar.svelte';

	/** @type {import('./$types').LayoutData}*/
	export let data;
</script>

<div class="mx-56 my-12">
	<Header />
	<FilterBar {data} />

	<slot />
</div>
