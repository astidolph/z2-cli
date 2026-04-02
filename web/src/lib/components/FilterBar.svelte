<script lang="ts">
	import { getFilters, setFilters, resetFilters } from '$lib/filters.svelte';

	let { onchange }: { onchange: () => void } = $props();

	const filters = getFilters();

	const currentYear = new Date().getFullYear();
	const yearOptions: number[] = [];
	for (let y = currentYear; y >= 2008; y--) {
		yearOptions.push(y);
	}

	// Local staging state — copied from store, committed on Apply
	let weeks = $state(filters.weeks);
	let year = $state(filters.year);
	let day = $state(filters.day);
	let minDistance = $state(filters.minDistance);
	let maxDistance = $state(filters.maxDistance);
	let maxHR = $state(filters.maxHR);
	let showAll = $state(filters.showAll);

	// Track whether weeks or year is the active time mode
	let timeMode = $derived<'weeks' | 'year'>(filters.year > 0 ? 'year' : 'weeks');

	function apply() {
		if (timeMode === 'weeks') {
			setFilters({ weeks, year: 0, day, minDistance, maxDistance, maxHR, showAll });
		} else {
			setFilters({ weeks: 0, year, day, minDistance, maxDistance, maxHR, showAll });
		}
		onchange();
	}

	function reset() {
		resetFilters();
		const f = getFilters();
		weeks = f.weeks;
		year = f.year;
		day = f.day;
		minDistance = f.minDistance;
		maxDistance = f.maxDistance;
		maxHR = f.maxHR;
		showAll = f.showAll;
		onchange();
	}

	function setTimeMode(mode: 'weeks' | 'year') {
		if (mode === 'weeks') {
			year = 0;
			if (!weeks) weeks = 12;
			setFilters({ weeks, year: 0 });
		} else {
			weeks = 0;
			if (!year) year = currentYear;
			setFilters({ weeks: 0, year });
		}
	}
</script>

<div class="card filter-bar">
	<div class="filter-row">
		<div class="time-group">
			<div class="time-toggle">
				<button class:active={timeMode === 'weeks'} onclick={() => setTimeMode('weeks')}>Weeks</button>
				<button class:active={timeMode === 'year'} onclick={() => setTimeMode('year')}>Year</button>
			</div>
			{#if timeMode === 'weeks'}
				<input type="number" bind:value={weeks} min="1" max="104" />
			{:else}
				<select bind:value={year}>
					{#each yearOptions as y}
						<option value={y}>{y}</option>
					{/each}
				</select>
			{/if}
		</div>

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
			<span>Max Distance (km)</span>
			<input type="number" bind:value={maxDistance} min="0" step="0.5" placeholder="No limit" />
		</label>

		<label>
			<span>Max Avg HR</span>
			<input type="number" bind:value={maxHR} min="0" step="1" placeholder="No limit" />
		</label>

		<label class="checkbox-label">
			<input type="checkbox" bind:checked={showAll} />
			<span>All runs</span>
		</label>

		<button class="btn btn-primary" onclick={apply}>Apply</button>
		<button class="btn btn-reset" onclick={reset}>Reset</button>
	</div>
</div>

<style>
	.filter-bar {
		padding: 0.75rem 1.25rem;
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

	.time-group {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.time-toggle {
		display: flex;
		gap: 0;
	}

	.time-toggle button {
		padding: 0.2rem 0.5rem;
		font-size: 0.7rem;
		border: 1px solid var(--border);
		background: var(--bg-input);
		color: var(--text-muted);
		cursor: pointer;
	}

	.time-toggle button:first-child {
		border-radius: var(--radius) 0 0 var(--radius);
	}

	.time-toggle button:last-child {
		border-radius: 0 var(--radius) var(--radius) 0;
		border-left: none;
	}

	.time-toggle button.active {
		background: var(--accent);
		color: #fff;
		border-color: var(--accent);
	}

	.btn-reset {
		padding: 0.375rem 0.75rem;
		background: var(--bg-input);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		color: var(--text-secondary);
		font-size: 0.875rem;
		cursor: pointer;
	}

	.btn-reset:hover {
		color: var(--text-primary);
		border-color: var(--text-secondary);
	}
</style>
