<script lang="ts">
	import { Chart, registerables } from 'chart.js';

	Chart.register(...registerables);

	interface Dataset {
		label: string;
		data: (number | null)[];
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

	$effect(() => {
		if (!canvas) return;

		// Clone data out of Svelte's reactive proxies — Chart.js uses
		// Object.defineProperty internally which conflicts with $state proxies.
		const currentLabels = [...labels];
		const currentDatasets = datasets.map((ds) => ({ ...ds, data: [...ds.data] }));

		const textColor = '#9b93b4';
		const gridColor = '#2d2560';

		chart = new Chart(canvas, {
			type: 'line',
			data: {
				labels: currentLabels,
				datasets: currentDatasets.map((ds) => ({
					label: ds.label,
					data: ds.data,
					borderColor: ds.color,
					backgroundColor: ds.color + '20',
					borderWidth: 2,
					pointRadius: 3,
					pointHoverRadius: 5,
					tension: 0.3,
					fill: false,
					spanGaps: true,
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

		return () => {
			chart?.destroy();
			chart = undefined;
		};
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
