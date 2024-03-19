<script>
	import { Chart, Card,  } from 'flowbite-svelte';
	
  export let data;

	let currency = data.currency; 
	let dates = data.dates; 
	let values = data.values; 

	let options = {
		chart: {
			height: '400px',
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
			width: 6,
			curve: 'smooth'
		},
		grid: {
			show: true,
			strokeDashArray: 4,
			padding: {
				left: 2,
				right: 2,
				top: -26
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
			show: false
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
			show: false
		}
	};
</script>

<Card>
	<div class="mb-5 flex justify-between">
		<div class="grid grid-cols-2 gap-4">
				
		</div>
	</div>
	<Chart {options} />
</Card>
