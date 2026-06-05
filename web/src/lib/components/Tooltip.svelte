<script lang="ts">
	let { text }: { text: string } = $props();

	let btn: HTMLButtonElement | undefined = $state();
	let visible = $state(false);
	let tipStyle = $state('');

	function open() {
		if (!btn) return;
		const r = btn.getBoundingClientRect();
		const width = 240;
		const left = Math.max(8, Math.min(r.left + r.width / 2 - width / 2, window.innerWidth - width - 8));
		tipStyle = `top:${r.top}px;left:${left}px;width:${width}px`;
		visible = true;
	}

	function close() {
		visible = false;
	}
</script>

<span class="wrap">
	<button
		bind:this={btn}
		type="button"
		class="tip-btn"
		aria-label="More info"
		onmouseenter={open}
		onmouseleave={close}
		onfocus={open}
		onblur={close}
	>?</button>
	{#if visible}
		<div class="tip-body" style={tipStyle} role="tooltip">{text}</div>
	{/if}
</span>

<style>
	.wrap {
		display: inline-flex;
		align-items: center;
		margin-left: 0.25rem;
	}

	.tip-btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 1rem;
		height: 1rem;
		border-radius: 50%;
		border: 1px solid var(--text-muted);
		background: transparent;
		color: var(--text-muted);
		font-size: 0.625rem;
		font-weight: 700;
		cursor: default;
		line-height: 1;
		padding: 0;
		flex-shrink: 0;
		transition: border-color 0.15s, color 0.15s;
	}

	.tip-btn:hover,
	.tip-btn:focus {
		border-color: var(--accent);
		color: var(--accent);
		outline: none;
	}

	.tip-body {
		position: fixed;
		transform: translateY(calc(-100% - 8px));
		background: var(--bg-input);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: 0.5rem 0.75rem;
		font-size: 0.75rem;
		color: var(--text-primary);
		line-height: 1.5;
		z-index: 1000;
		pointer-events: none;
		font-weight: normal;
		text-transform: none;
		letter-spacing: normal;
	}
</style>
