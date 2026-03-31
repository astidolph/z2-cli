<script lang="ts">
	import { api } from '$lib/api';
	import type { ChartDataResponse } from '$lib/types';
	import LineChart from '$lib/components/LineChart.svelte';

	let chartData: ChartDataResponse | null = $state(null);
	let error: string | null = $state(null);
	let loading = $state(true);
	let weeks = $state(12);

	async function load() {
		loading = true;
		error = null;
		try {
			chartData = await api.getChartData({ weeks });
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load chart data';
		} finally {
			loading = false;
		}
	}

	function changeWeeks() {
		load();
	}

	$effect(() => {
		load();
	});
</script>

<div class="charts-page">
	<div class="header">
		<h1>Charts</h1>
		<div class="controls">
			<label>
				<span>Weeks</span>
				<input type="number" bind:value={weeks} min="1" max="104" onchange={changeWeeks} />
			</label>
		</div>
	</div>

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

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	h1 {
		font-size: 1.5rem;
	}

	.controls label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.875rem;
		color: var(--text-secondary);
	}

	.controls input {
		padding: 0.375rem 0.5rem;
		background: var(--bg-input);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		color: var(--text-primary);
		font-size: 0.875rem;
		width: 80px;
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
