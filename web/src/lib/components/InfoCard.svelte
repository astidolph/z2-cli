<script lang="ts">
	const KEY = 'z2_info_dismissed';

	let dismissed = $state(false);

	$effect(() => {
		dismissed = localStorage.getItem(KEY) === 'true';
	});

	function dismiss() {
		localStorage.setItem(KEY, 'true');
		dismissed = true;
	}
</script>

{#if !dismissed}
	<div class="info-card card">
		<div class="content">
			<p class="headline">How it works</p>
			<p class="body">
				Zone 2 training targets a comfortable aerobic effort — roughly 60–70% of your max heart
				rate — where your body builds its aerobic base most efficiently.
				<strong>EF (Efficiency Factor)</strong> measures speed ÷ heart rate × 1000. A rising EF
				score means you're running faster at the same effort: your aerobic fitness is improving.
			</p>
		</div>
		<button class="dismiss" onclick={dismiss} aria-label="Dismiss">✕</button>
	</div>
{/if}

<style>
	.info-card {
		display: flex;
		align-items: flex-start;
		gap: 1rem;
		border-left: 3px solid var(--accent);
		padding: 0.875rem 1rem;
	}

	.content {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 0.375rem;
	}

	.headline {
		font-size: 0.75rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: var(--accent);
	}

	.body {
		font-size: 0.875rem;
		color: var(--text-secondary);
		line-height: 1.5;
	}

	.body strong {
		color: var(--text-primary);
	}

	.dismiss {
		background: transparent;
		border: none;
		color: var(--text-muted);
		font-size: 0.75rem;
		cursor: pointer;
		padding: 0.125rem 0.25rem;
		line-height: 1;
		flex-shrink: 0;
	}

	.dismiss:hover {
		color: var(--text-secondary);
	}
</style>
