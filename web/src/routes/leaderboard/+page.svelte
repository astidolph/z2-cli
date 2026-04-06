<script lang="ts">
	import { api, filtersToLeaderboardParams } from '$lib/api';
	import { getFilters } from '$lib/filters.svelte';
	import type { LeaderboardResponse } from '$lib/types';
	import { formatDate, formatDistance, formatDuration, formatHR, formatPace, efficiencyFactor, formatEF } from '$lib/format';

	const filters = getFilters();

	let data: LeaderboardResponse | null = $state(null);
	let error: string | null = $state(null);
	let loading = $state(true);
	let refreshing = $state(false);
	let page = $state(1);

	async function load() {
		loading = true;
		error = null;
		try {
			data = await api.getLeaderboard(filtersToLeaderboardParams(filters, { page }));
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load leaderboard';
		} finally {
			loading = false;
		}
	}

	async function refresh() {
		refreshing = true;
		error = null;
		try {
			await api.refreshLeaderboard();
			page = 1;
			await load();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to refresh leaderboard';
		} finally {
			refreshing = false;
		}
	}

	function goToPage(p: number) {
		page = p;
	}

	$effect(() => {
		// Reset to page 1 when any filter changes
		void [filters.weeks, filters.year, filters.day, filters.minDistance, filters.maxDistance, filters.maxHR, filters.showAll];
		page = 1;
	});

	$effect(() => {
		// Load data whenever page changes (covers both filter resets and manual pagination)
		void page;
		load();
	});

	let totalPages = $derived.by(() => {
		if (!data) return 0;
		return Math.ceil(data.total_count / data.page_size);
	});
	let rankOffset = $derived.by(() => {
		if (!data) return 0;
		return (data.page - 1) * data.page_size;
	});
</script>

<div class="leaderboard-page">
	<div class="header">
		<h1>Leaderboard</h1>
		<button class="refresh-btn" onclick={refresh} disabled={refreshing}>
			{refreshing ? 'Syncing...' : 'Sync from Strava'}
		</button>
	</div>

	{#if refreshing}
		<p class="status">Syncing your run history from Strava — this may take a while on the first load...</p>
	{:else if loading}
		<p class="status">Loading...</p>
	{:else if error}
		<p class="status error">{error}</p>
	{:else if data}
		{#if !data.runs?.length}
			<p class="status">No runs found. Click "Sync from Strava" to load your run history.</p>
		{:else}
			<p class="total">{data.total_count} runs ranked by Efficiency Factor</p>
			<div class="card">
				<div class="table-wrapper">
					<table>
						<thead>
							<tr>
								<th class="rank">#</th>
								<th>Date</th>
								<th>Name</th>
								<th>Distance</th>
								<th>Time</th>
								<th>Avg HR</th>
								<th>Pace</th>
								<th>EF</th>
							</tr>
						</thead>
						<tbody>
							{#each data.runs as run, i}
								<tr>
									<td class="rank">{rankOffset + i + 1}</td>
									<td>{formatDate(run.start_date_local)}</td>
									<td class="name"><a href="https://www.strava.com/activities/{run.id}" target="_blank" rel="noopener noreferrer">{run.name}</a></td>
									<td>{formatDistance(run.distance)}</td>
									<td>{formatDuration(run.moving_time)}</td>
									<td>{formatHR(run.average_heartrate)}</td>
									<td>{run.distance > 0 ? formatPace(run.moving_time / (run.distance / 1000)) : '-'}</td>
									<td class="ef">{formatEF(efficiencyFactor(run.distance, run.moving_time, run.average_heartrate))}</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>

			{#if totalPages > 1}
				<div class="pagination">
					<button onclick={() => goToPage(page - 1)} disabled={page <= 1}>Previous</button>
					<span>Page {page} of {totalPages}</span>
					<button onclick={() => goToPage(page + 1)} disabled={page >= totalPages}>Next</button>
				</div>
			{/if}
		{/if}
	{/if}
</div>

<style>
	.leaderboard-page {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	h1 {
		font-size: 1.5rem;
	}

	.refresh-btn {
		padding: 0.5rem 1rem;
		background: var(--accent);
		color: var(--bg-card);
		border: none;
		border-radius: var(--radius);
		cursor: pointer;
		font-size: 0.875rem;
		font-weight: 500;
	}

	.refresh-btn:hover:not(:disabled) {
		opacity: 0.9;
	}

	.refresh-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.total {
		color: var(--text-secondary);
		font-size: 0.875rem;
	}

	.table-wrapper {
		overflow-x: auto;
	}

	table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.875rem;
	}

	th {
		text-align: left;
		padding: 0.625rem 0.75rem;
		color: var(--text-muted);
		font-size: 0.75rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		border-bottom: 1px solid var(--border);
	}

	td {
		padding: 0.625rem 0.75rem;
		border-bottom: 1px solid var(--border);
		white-space: nowrap;
	}

	td.name {
		white-space: normal;
		max-width: 200px;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	td.rank, th.rank {
		text-align: center;
		color: var(--text-muted);
		font-weight: 600;
		min-width: 3rem;
	}

	td.ef {
		font-weight: 600;
	}

	tr:hover td {
		background: var(--bg-input);
	}

	.status {
		color: var(--text-secondary);
		padding: 2rem;
		text-align: center;
	}

	.status.error {
		color: var(--negative);
	}

	.pagination {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 1rem;
	}

	.pagination button {
		padding: 0.375rem 0.75rem;
		background: var(--bg-input);
		color: var(--text-primary);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		cursor: pointer;
		font-size: 0.875rem;
	}

	.pagination button:hover:not(:disabled) {
		background: var(--border);
	}

	.pagination button:disabled {
		opacity: 0.4;
		cursor: not-allowed;
	}

	.pagination span {
		color: var(--text-secondary);
		font-size: 0.875rem;
	}
</style>
