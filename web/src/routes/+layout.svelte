<script lang="ts">
	import '../app.css';
	import NavBar from '$lib/components/NavBar.svelte';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<NavBar />
<main>
	{@render children()}
</main>

<style>
	main {
		max-width: 1200px;
		margin: 0 auto;
		padding: 1.5rem;
	}
</style>
