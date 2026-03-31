<script lang="ts">
	import { api } from '$lib/api';
	import type { RunsResponse, ChartDataResponse } from '$lib/types';
	import SummaryCard from '$lib/components/SummaryCard.svelte';
	import LineChart from '$lib/components/LineChart.svelte';

	let runs: RunsResponse | null = $state(null);
	let chartData: ChartDataResponse | null = $state(null);
	let error: string | null = $state(null);
	let loading = $state(true);

	async function load() {
		loading = true;
		error = null;
		try {
			[runs, chartData] = await Promise.all([api.getRuns(), api.getChartData()]);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load data';
		} finally {
			loading = false;
		}
	}

	async function refresh() {
		await api.refresh();
		load();
	}

	$effect(() => {
		load();
	});
</script>

<div class="dashboard">
	<div class="header">
		<h1>Zone 2 Dashboard</h1>
		<button class="btn" onclick={refresh}>Refresh</button>
	</div>

	{#if loading}
		<p class="status">Loading...</p>
	{:else if error}
		<p class="status error">{error}</p>
	{:else if runs}
		<div class="summaries">
			<SummaryCard summary={runs.current} label="Current ({runs.weeks_back} weeks)" trend={runs.ef_trend} />
			<SummaryCard summary={runs.prior} label="Prior ({runs.weeks_back} weeks)" />
		</div>

		{#if chartData?.dates?.length}
			<div class="card chart-card">
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
		{/if}
	{/if}
</div>

<style>
	.dashboard {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	h1 {
		font-size: 1.5rem;
	}

	.summaries {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
	}

	.chart-card {
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

	@media (max-width: 768px) {
		.summaries {
			grid-template-columns: 1fr;
		}
	}
</style>
