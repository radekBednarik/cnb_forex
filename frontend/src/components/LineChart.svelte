<script>
	import { Chart, Card } from 'flowbite-svelte';

	export let data;

	let currency = data.currency;
	let dates = data.dates;
	let values = data.values;

	let options = {
		chart: {
			height: '500px',
			maxWidth: '100%',
			type: 'line',
			fontFamily: 'Inter, sans-serif',
			dropShadow: {
				enabled: false
			},
			toolbar: {
				show: false
			}
		},
		tooltip: {
			enabled: true,
			x: {
				show: false
			}
		},
		dataLabels: {
			enabled: false
		},
		stroke: {
			width: 2,
			curve: 'smooth'
		},
		grid: {
			show: true,
			strokeDashArray: 4,
			padding: {
				left: 2,
				right: 4,
				top: 0
			}
		},
		series: [
			{
				name: currency,
				data: values,
				color: '#1A56DB'
			}
		],
		legend: {
			show: true,
			showForSingleSeries: true,
			position: 'bottom',
			horizontalPosition: 'right'
		},
		xaxis: {
			categories: dates,
			labels: {
				show: true,
				style: {
					fontFamily: 'Inter, sans-serif',
					cssClass: 'text-xs font-normal fill-gray-500 dark:fill-gray-400'
				}
			},
			axisBorder: {
				show: false
			},
			axisTicks: {
				show: false
			}
		},
		yaxis: {
			show: true
		}
	};
</script>

<Card padding="sm" class="max-w-screen-2xl">
	<Chart {options} />
</Card>
