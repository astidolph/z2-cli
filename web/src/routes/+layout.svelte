<script lang="ts">
	import '../app.css';
	import { page } from '$app/state';
	import NavBar from '$lib/components/NavBar.svelte';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();

	let showFilters = $derived(page.url.pathname !== '/settings');
</script>

<NavBar />
{#if showFilters}
	<div class="filter-container">
		<FilterBar onchange={() => {}} />
	</div>
{/if}
<main>
	{@render children()}
</main>

<style>
	.filter-container {
		max-width: 1200px;
		margin: 0 auto;
		padding: 1rem 1.5rem 0;
	}

	main {
		max-width: 1200px;
		margin: 0 auto;
		padding: 1.5rem;
	}
</style>
