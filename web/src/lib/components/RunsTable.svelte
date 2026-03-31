<script lang="ts">
	import type { Activity } from '$lib/types';
	import { formatDate, formatDistance, formatDuration, formatHR, formatPace, efficiencyFactor, formatEF } from '$lib/format';

	let { runs }: { runs: Activity[] } = $props();
</script>

<div class="table-wrapper">
	<table>
		<thead>
			<tr>
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
			{#each runs as run}
				<tr>
					<td>{formatDate(run.start_date_local)}</td>
					<td class="name">{run.name}</td>
					<td>{formatDistance(run.distance)}</td>
					<td>{formatDuration(run.moving_time)}</td>
					<td>{run.has_heartrate ? formatHR(run.average_heartrate) : '-'}</td>
					<td>{run.distance > 0 ? formatPace(run.moving_time / (run.distance / 1000)) : '-'}</td>
					<td>{formatEF(efficiencyFactor(run.distance, run.moving_time, run.average_heartrate))}</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>

<style>
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

	tr:hover td {
		background: var(--bg-input);
	}
</style>
