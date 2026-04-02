<script lang="ts">
	import { api, filtersToRunsParams } from '$lib/api';
	import { getFilters } from '$lib/filters.svelte';
	import type { ChartDataResponse } from '$lib/types';
	import LineChart from '$lib/components/LineChart.svelte';

	const filters = getFilters();

	let chartData: ChartDataResponse | null = $state(null);
	let error: string | null = $state(null);
	let loading = $state(true);

	async function load() {
		loading = true;
		error = null;
		try {
			chartData = await api.getChartData(filtersToRunsParams(filters));
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load chart data';
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		filtersToRunsParams(filters);
		load();
	});
</script>

<div class="charts-page">
	<h1>Charts</h1>

	{#if loading}
		<p class="status">Loading...</p>
	{:else if error}
		<p class="status error">{error}</p>
	{:else if chartData?.dates?.length}
		<div class="card">
			<LineChart
				labels={chartData.dates}
				datasets={[
					{ label: 'EF', data: chartData.ef ?? [], color: '#7c6ef0' },
					{ label: 'Avg HR', data: chartData.hr ?? [], color: '#f87171', yAxisID: 'y2' }
				]}
				title="Efficiency Factor & Heart Rate"
				dualAxis={true}
				secondAxisLabel="HR (bpm)"
			/>
		</div>

		<div class="card">
			<LineChart
				labels={chartData.dates}
				datasets={[
					{ label: 'Pace /km', data: chartData.pace ?? [], color: '#38bdf8' },
					{ label: 'Pace /mi', data: chartData.pace_mi ?? [], color: '#818cf8' }
				]}
				title="Pace Over Time"
			/>
		</div>

		<div class="card">
			<LineChart
				labels={chartData.dates}
				datasets={[
					{ label: 'km', data: chartData.distance ?? [], color: '#4ade80' },
					{ label: 'mi', data: chartData.distance_mi ?? [], color: '#a78bfa' }
				]}
				title="Distance Over Time"
			/>
		</div>

		<div class="card">
			<LineChart
				labels={chartData.dates}
				datasets={[{ label: 'Avg HR', data: chartData.hr ?? [], color: '#f87171' }]}
				title="Average Heart Rate Over Time"
			/>
		</div>
	{:else}
		<p class="status">No data available for the selected period.</p>
	{/if}
</div>

<style>
	.charts-page {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	h1 {
		font-size: 1.5rem;
	}

	.card {
		padding: 1.25rem;
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
