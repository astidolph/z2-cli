<script lang="ts">
	import { api, filtersToRunsParams } from '$lib/api';
	import { getFilters } from '$lib/filters.svelte';
	import type { RunsResponse } from '$lib/types';
	import RunsTable from '$lib/components/RunsTable.svelte';
	import SummaryCard from '$lib/components/SummaryCard.svelte';

	const filters = getFilters();

	let sort = $state('date');
	let asc = $state(false);
	let data: RunsResponse | null = $state(null);
	let error: string | null = $state(null);
	let loading = $state(true);

	async function load() {
		loading = true;
		error = null;
		try {
			data = await api.getRuns(filtersToRunsParams(filters, { sort, asc }));
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load runs';
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		filtersToRunsParams(filters);
		load();
	});
</script>

<div class="runs-page">
	<div class="header">
		<h1>Runs</h1>
		<div class="sort-controls">
			<label>
				<span>Sort by</span>
				<select bind:value={sort} onchange={() => load()}>
					<option value="date">Date</option>
					<option value="distance">Distance</option>
					<option value="time">Time</option>
					<option value="hr">Heart Rate</option>
					<option value="pace">Pace</option>
					<option value="ef">EF</option>
				</select>
			</label>

			<label class="checkbox-label">
				<input type="checkbox" bind:checked={asc} onchange={() => load()} />
				<span>Ascending</span>
			</label>
		</div>
	</div>

	{#if loading}
		<p class="status">Loading...</p>
	{:else if error}
		<p class="status error">{error}</p>
	{:else if data}
		<SummaryCard summary={data.current} label="Summary ({data.weeks_back} weeks)" trend={data.ef_trend} />

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

	.header {
		display: flex;
		justify-content: space-between;
		align-items: flex-end;
	}

	h1 {
		font-size: 1.5rem;
	}

	.sort-controls {
		display: flex;
		align-items: flex-end;
		gap: 1rem;
	}

	.sort-controls label {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		font-size: 0.75rem;
		color: var(--text-muted);
	}

	.sort-controls select {
		padding: 0.375rem 0.5rem;
		background: var(--bg-input);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		color: var(--text-primary);
		font-size: 0.875rem;
	}

	.checkbox-label {
		flex-direction: row;
		align-items: center;
		gap: 0.375rem;
		font-size: 0.875rem;
		color: var(--text-secondary);
		padding-bottom: 0.375rem;
	}

	.checkbox-label input[type='checkbox'] {
		accent-color: var(--accent);
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
