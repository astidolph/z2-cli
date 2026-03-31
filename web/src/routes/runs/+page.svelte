<script lang="ts">
	import { api, type RunsParams } from '$lib/api';
	import type { RunsResponse } from '$lib/types';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import RunsTable from '$lib/components/RunsTable.svelte';
	import SummaryCard from '$lib/components/SummaryCard.svelte';

	let params: RunsParams = $state({ weeks: 12, sort: 'date' });
	let data: RunsResponse | null = $state(null);
	let error: string | null = $state(null);
	let loading = $state(true);

	async function load() {
		loading = true;
		error = null;
		try {
			data = await api.getRuns(params);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load runs';
		} finally {
			loading = false;
		}
	}

	function onFilterChange(newParams: RunsParams) {
		params = newParams;
		load();
	}

	$effect(() => {
		load();
	});
</script>

<div class="runs-page">
	<h1>Runs</h1>

	<FilterBar {params} onchange={onFilterChange} />

	{#if loading}
		<p class="status">Loading...</p>
	{:else if error}
		<p class="status error">{error}</p>
	{:else if data}
		<SummaryCard summary={data.current} label="Summary ({data.weeks_back} weeks)" />

		{#if data.current_runs?.length}
			<div class="card">
				<RunsTable runs={data.current_runs} />
			</div>
		{:else}
			<p class="status">No runs found for the selected filters.</p>
		{/if}
	{/if}
</div>

<style>
	.runs-page {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	h1 {
		font-size: 1.5rem;
	}

	.status {
		color: var(--text-secondary);
		padding: 2rem;
		text-align: center;
	}

	.status.error {
		color: var(--negative);
	}
</style>
