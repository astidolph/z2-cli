<script lang="ts">
	import type { Activity } from '$lib/types';
	import { formatDate, formatDistance, formatDuration, formatHR, formatPace, efficiencyFactor, formatEF } from '$lib/format';

	let { runs }: { runs: Activity[] } = $props();

	// Compute a cumulative average EF for each run in date order, then map back
	// so each row shows the avg EF up to and including that run.
	let avgEFMap = $derived.by(() => {
		const dated = runs.map((r, i) => ({
			index: i,
			date: new Date(r.start_date_local).getTime(),
			ef: efficiencyFactor(r.distance, r.moving_time, r.average_heartrate)
		}));
		dated.sort((a, b) => a.date - b.date);

		const map = new Map<number, number>();
		let sum = 0;
		let count = 0;
		for (const entry of dated) {
			if (entry.ef > 0) {
				sum += entry.ef;
				count++;
			}
			map.set(entry.index, count > 0 ? sum / count : 0);
		}
		return map;
	});
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
				<th>Avg EF</th>
			</tr>
		</thead>
		<tbody>
			{#each runs as run, i}
				<tr>
					<td>{formatDate(run.start_date_local)}</td>
					<td class="name"><a href="https://www.strava.com/activities/{run.id}" target="_blank" rel="noopener noreferrer">{run.name}</a></td>
					<td>{formatDistance(run.distance)}</td>
					<td>{formatDuration(run.moving_time)}</td>
					<td>{run.has_heartrate ? formatHR(run.average_heartrate) : '-'}</td>
					<td>{run.distance > 0 ? formatPace(run.moving_time / (run.distance / 1000)) : '-'}</td>
					<td>{formatEF(efficiencyFactor(run.distance, run.moving_time, run.average_heartrate))}</td>
					<td>{formatEF(avgEFMap.get(i) ?? 0)}</td>
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
