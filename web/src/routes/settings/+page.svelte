<script lang="ts">
	import { api } from '$lib/api';
	import type { AuthStatusResponse } from '$lib/types';

	let zone2HR: number = $state(0);
	let age: number = $state(0);
	let authStatus: AuthStatusResponse | null = $state(null);
	let loading = $state(true);
	let saving = $state(false);
	let message: string | null = $state(null);
	let error: string | null = $state(null);

	// Check for OAuth callback results in URL params
	function checkAuthResult() {
		const params = new URLSearchParams(window.location.search);
		const authError = params.get('auth_error');
		const authSuccess = params.get('auth_success');
		if (authError) {
			error = authError;
		} else if (authSuccess) {
			message = 'Successfully connected to Strava!';
		}
		// Clean up URL params
		if (authError || authSuccess) {
			window.history.replaceState({}, '', window.location.pathname);
		}
	}

	async function load() {
		loading = true;
		checkAuthResult();
		try {
			const [config, auth] = await Promise.all([api.getConfig(), api.getAuthStatus()]);
			zone2HR = config.zone2_hr;
			authStatus = auth;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load settings';
		} finally {
			loading = false;
		}
	}

	async function saveHR() {
		saving = true;
		message = null;
		error = null;
		try {
			const result = await api.putConfig({ zone2_hr: zone2HR });
			zone2HR = result.zone2_hr;
			message = `Zone 2 HR ceiling set to ${result.zone2_hr} bpm`;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to save';
		} finally {
			saving = false;
		}
	}

	async function saveAge() {
		if (age <= 0) {
			error = 'Please enter a valid age';
			return;
		}
		saving = true;
		message = null;
		error = null;
		try {
			const result = await api.putConfig({ age });
			zone2HR = result.zone2_hr;
			message = `Zone 2 HR ceiling calculated as ${result.zone2_hr} bpm (180 - ${age})`;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to save';
		} finally {
			saving = false;
		}
	}

	$effect(() => {
		load();
	});
</script>

<div class="settings-page">
	<h1>Settings</h1>

	{#if loading}
		<p class="status">Loading...</p>
	{:else}
		<div class="card">
			<h2>Strava Connection</h2>
			{#if authStatus?.authenticated}
				<p class="auth-ok">Connected</p>
				<a href="/api/auth/login" class="btn btn-secondary btn-small">Reconnect</a>
			{:else}
				<a href="/api/auth/login" class="btn btn-primary">Connect Strava</a>
			{/if}
		</div>

		<div class="card">
			<h2>Zone 2 Heart Rate Ceiling</h2>
			<p class="description">Runs with an average HR at or below this value are considered Zone 2.</p>

			<div class="form-group">
				<label>
					<span>Direct value (bpm)</span>
					<div class="input-row">
						<input type="number" bind:value={zone2HR} min="1" max="220" />
						<button class="btn btn-primary" onclick={saveHR} disabled={saving}>Save</button>
					</div>
				</label>
			</div>

			<div class="divider">or</div>

			<div class="form-group">
				<label>
					<span>Calculate from age (Maffetone: 180 - age)</span>
					<div class="input-row">
						<input type="number" bind:value={age} min="1" max="120" placeholder="Your age" />
						<button class="btn btn-primary" onclick={saveAge} disabled={saving}>Calculate & Save</button>
					</div>
				</label>
			</div>
		</div>

		{#if message}
			<p class="feedback success">{message}</p>
		{/if}
		{#if error}
			<p class="feedback error">{error}</p>
		{/if}
	{/if}
</div>

<style>
	.settings-page {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
		max-width: 600px;
	}

	h1 {
		font-size: 1.5rem;
	}

	h2 {
		font-size: 1rem;
		margin-bottom: 0.75rem;
	}

	.description {
		font-size: 0.875rem;
		color: var(--text-secondary);
		margin-bottom: 1rem;
	}

	.auth-ok {
		color: var(--positive);
		font-weight: 600;
	}

	.btn-secondary {
		background: var(--bg-input);
		border: 1px solid var(--border);
		color: var(--text-secondary);
	}

	.btn-secondary:hover {
		color: var(--text-primary);
		border-color: var(--text-secondary);
	}

	.btn-small {
		padding: 0.25rem 0.75rem;
		font-size: 0.75rem;
	}

	code {
		background: var(--bg-input);
		padding: 0.125rem 0.375rem;
		border-radius: 4px;
		font-size: 0.8125rem;
	}

	.form-group {
		margin-bottom: 0.5rem;
	}

	label span {
		display: block;
		font-size: 0.75rem;
		color: var(--text-muted);
		margin-bottom: 0.375rem;
	}

	.input-row {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	input[type='number'] {
		padding: 0.375rem 0.5rem;
		background: var(--bg-input);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		color: var(--text-primary);
		font-size: 0.875rem;
		width: 120px;
	}

	.divider {
		text-align: center;
		color: var(--text-muted);
		font-size: 0.875rem;
		margin: 0.5rem 0;
	}

	.feedback {
		font-size: 0.875rem;
		padding: 0.75rem 1rem;
		border-radius: var(--radius);
	}

	.feedback.success {
		color: var(--positive);
		background: #4ade8015;
	}

	.feedback.error {
		color: var(--negative);
		background: #f8717115;
	}

	.status {
		color: var(--text-secondary);
		padding: 2rem;
		text-align: center;
	}
</style>
