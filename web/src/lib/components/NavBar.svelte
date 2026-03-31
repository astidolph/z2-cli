<script lang="ts">
	import { page } from '$app/state';

	const links = [
		{ href: '/', label: 'Dashboard' },
		{ href: '/runs', label: 'Runs' },
		{ href: '/charts', label: 'Charts' },
		{ href: '/settings', label: 'Settings' }
	];
</script>

<nav>
	<div class="nav-brand">z2-cli</div>
	<div class="nav-links">
		{#each links as link}
			<a href={link.href} class:active={page.url.pathname === link.href}>
				{link.label}
			</a>
		{/each}
	</div>
</nav>

<style>
	nav {
		display: flex;
		align-items: center;
		gap: 2rem;
		padding: 0.75rem 1.5rem;
		background: var(--bg-card);
		border-bottom: 1px solid var(--border);
	}

	.nav-brand {
		font-weight: 700;
		font-size: 1.125rem;
		color: var(--accent);
	}

	.nav-links {
		display: flex;
		gap: 0.25rem;
	}

	.nav-links a {
		padding: 0.375rem 0.75rem;
		border-radius: var(--radius);
		color: var(--text-secondary);
		font-size: 0.875rem;
		transition: color 0.15s, background 0.15s;
	}

	.nav-links a:hover {
		color: var(--text-primary);
		background: var(--bg-input);
	}

	.nav-links a.active {
		color: var(--text-primary);
		background: var(--accent);
	}
</style>
