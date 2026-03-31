<script lang="ts">
	import type { Summary } from '$lib/types';
	import { formatDistanceKm, formatPace, formatEF, formatHR } from '$lib/format';

	let { summary, label, trend }: { summary: Summary; label: string; trend?: number } = $props();
</script>

<div class="card summary-card">
	<h3>{label}</h3>
	<div class="stats-grid">
		<div class="stat">
			<span class="stat-value">{summary.count}</span>
			<span class="stat-label">Runs</span>
		</div>
		<div class="stat">
			<span class="stat-value">{formatDistanceKm(summary.total_km)}</span>
			<span class="stat-label">Total Distance</span>
		</div>
		<div class="stat">
			<span class="stat-value">{formatEF(summary.avg_ef)}</span>
			<span class="stat-label">Avg EF</span>
			{#if trend !== undefined && trend !== 0}
				<span class="trend" class:positive={trend > 0} class:negative={trend < 0}>
					{trend > 0 ? '+' : ''}{trend.toFixed(1)}%
				</span>
			{/if}
		</div>
		<div class="stat">
			<span class="stat-value">{formatHR(summary.avg_hr)}</span>
			<span class="stat-label">Avg HR</span>
		</div>
		<div class="stat">
			<span class="stat-value">{formatPace(summary.avg_pace)}</span>
			<span class="stat-label">Avg Pace</span>
		</div>
	</div>
</div>

<style>
	.summary-card h3 {
		margin-bottom: 1rem;
		font-size: 0.875rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: var(--text-secondary);
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
		gap: 1rem;
	}

	.stat {
		display: flex;
		flex-direction: column;
		gap: 0.125rem;
	}

	.stat-value {
		font-size: 1.25rem;
		font-weight: 600;
	}

	.stat-label {
		font-size: 0.75rem;
		color: var(--text-muted);
	}

	.trend {
		font-size: 0.75rem;
		font-weight: 600;
	}

	.trend.positive {
		color: var(--positive);
	}

	.trend.negative {
		color: var(--negative);
	}
</style>
