<script lang="ts">
	import { onMount } from 'svelte';
	import { Chart, registerables } from 'chart.js';

	Chart.register(...registerables);

	interface Dataset {
		label: string;
		data: number[];
		color: string;
		yAxisID?: string;
	}

	let {
		labels,
		datasets,
		title = '',
		dualAxis = false,
		secondAxisLabel = ''
	}: {
		labels: string[];
		datasets: Dataset[];
		title?: string;
		dualAxis?: boolean;
		secondAxisLabel?: string;
	} = $props();

	let canvas: HTMLCanvasElement;
	let chart: Chart | undefined;

	function buildChart() {
		if (chart) chart.destroy();

		const textColor = '#9b93b4';
		const gridColor = '#2d2560';

		chart = new Chart(canvas, {
			type: 'line',
			data: {
				labels,
				datasets: datasets.map((ds) => ({
					label: ds.label,
					data: ds.data,
					borderColor: ds.color,
					backgroundColor: ds.color + '20',
					borderWidth: 2,
					pointRadius: 3,
					pointHoverRadius: 5,
					tension: 0.3,
					fill: false,
					yAxisID: ds.yAxisID || 'y'
				}))
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				interaction: {
					mode: 'index',
					intersect: false
				},
				plugins: {
					title: {
						display: !!title,
						text: title,
						color: textColor,
						font: { size: 14 }
					},
					legend: {
						labels: { color: textColor }
					}
				},
				scales: {
					x: {
						ticks: { color: textColor },
						grid: { color: gridColor }
					},
					y: {
						ticks: { color: textColor },
						grid: { color: gridColor }
					},
					...(dualAxis
						? {
								y2: {
									position: 'right' as const,
									title: {
										display: !!secondAxisLabel,
										text: secondAxisLabel,
										color: textColor
									},
									ticks: { color: textColor },
									grid: { drawOnChartArea: false }
								}
							}
						: {})
				}
			}
		});
	}

	onMount(() => {
		buildChart();
		return () => chart?.destroy();
	});

	$effect(() => {
		// Re-render when data changes
		if (canvas && labels && datasets) {
			buildChart();
		}
	});
</script>

<div class="chart-wrapper">
	<canvas bind:this={canvas}></canvas>
</div>

<style>
	.chart-wrapper {
		position: relative;
		height: 350px;
	}
</style>
