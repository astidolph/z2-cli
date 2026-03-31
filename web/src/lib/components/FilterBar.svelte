<script lang="ts">
	import type { RunsParams } from '$lib/api';

	let { params: initialParams, onchange }: { params: RunsParams; onchange: (p: RunsParams) => void } = $props();

	let weeks = $state(initialParams.weeks ?? 12);
	let day = $state(initialParams.day ?? '');
	let minDistance = $state(initialParams.minDistance ?? 0);
	let showAll = $state(initialParams.all ?? false);
	let sort = $state(initialParams.sort ?? 'date');
	let asc = $state(initialParams.asc ?? false);

	function apply() {
		onchange({
			weeks,
			day: day || undefined,
			minDistance: minDistance || undefined,
			all: showAll || undefined,
			sort,
			asc: asc || undefined
		});
	}
</script>

<div class="card filter-bar">
	<div class="filter-row">
		<label>
			<span>Weeks</span>
			<input type="number" bind:value={weeks} min="1" max="104" />
		</label>

		<label>
			<span>Day</span>
			<select bind:value={day}>
				<option value="">All days</option>
				<option value="monday">Monday</option>
				<option value="tuesday">Tuesday</option>
				<option value="wednesday">Wednesday</option>
				<option value="thursday">Thursday</option>
				<option value="friday">Friday</option>
				<option value="saturday">Saturday</option>
				<option value="sunday">Sunday</option>
			</select>
		</label>

		<label>
			<span>Min Distance (km)</span>
			<input type="number" bind:value={minDistance} min="0" step="0.5" />
		</label>

		<label>
			<span>Sort by</span>
			<select bind:value={sort}>
				<option value="date">Date</option>
				<option value="distance">Distance</option>
				<option value="time">Time</option>
				<option value="hr">Heart Rate</option>
				<option value="pace">Pace</option>
				<option value="ef">EF</option>
			</select>
		</label>

		<label class="checkbox-label">
			<input type="checkbox" bind:checked={asc} />
			<span>Ascending</span>
		</label>

		<label class="checkbox-label">
			<input type="checkbox" bind:checked={showAll} />
			<span>All runs</span>
		</label>

		<button class="btn btn-primary" onclick={apply}>Apply</button>
	</div>
</div>

<style>
	.filter-bar {
		padding: 1rem 1.25rem;
	}

	.filter-row {
		display: flex;
		flex-wrap: wrap;
		align-items: flex-end;
		gap: 1rem;
	}

	label {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		font-size: 0.75rem;
		color: var(--text-muted);
	}

	input[type='number'],
	select {
		padding: 0.375rem 0.5rem;
		background: var(--bg-input);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		color: var(--text-primary);
		font-size: 0.875rem;
		width: 120px;
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
</style>
