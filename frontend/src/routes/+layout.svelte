<script>
	import '../app.pcss';
	import Header from '../components/Header.svelte';
	import FilterBar from '../components/FilterBar.svelte';
</script>

<div class="mx-12 my-12">
	<Header />
	<FilterBar />

	<slot />
</div>
