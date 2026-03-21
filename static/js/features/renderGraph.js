//import { Chart, registerables } from "https://cdn.jsdelivr.net/npm/chart.js@4.5.1/dist/chart.js"
//import 'https://cdn.jsdelivr.net/npm/@kurkle/color@0.4.0/+esm'
//import 'https://cdnjs.cloudflare.com/ajax/libs/Chart.js/4.4.0/chart.umd.js';
//import Chart from '../vendor/chart/package/dist/chart.js'
import { Tooltip, LineController, LinearScale, LineElement, PointElement, CategoryScale, Chart } from 'https://cdn.jsdelivr.net/npm/chart.js@4.5.1/+esm'

export function initRenderGraph(root) {
  const container = root.querySelector('[data-feature="graph-chart"]')
  if (!container) return

  const jsonDataElement = root.querySelector('#graph-chart-json-data')
  if (!jsonDataElement) return

  const payload = JSON.parse(jsonDataElement.textContent)[0]
  const dataset_label = payload.Label
  const points = payload.Points
  const data = points.map(p => p.Data)
  const labels = points.map(p => p.Label)

  console.log(points)
  console.log(data)
  console.log(payload)
  console.log(labels)

  Chart.register(
    LineController,
    LineElement,
    PointElement,
    CategoryScale,
    LinearScale,
    Tooltip
  )

  createChart(container, dataset_label, data, labels)
}

function createChart(container, dataset_label, dataset, labels) {
  let chart = container.querySelector('#graph-chart-canvas').getContext('2d');
  let chartData = new Chart(chart, {
    type: 'line',
    data: {
      labels: labels,
      datasets: [{
        label: dataset_label,
        data: dataset,
        pointBackgroundColor: '#FF0000',
        borderColor: 'rgb(75, 192, 192)',
      }],
    },
    options: {
      plugins: {
        tooltip: {
          callbacks: {
            // Customizes the main body text
            label: function(context) {
              let label = context.dataset.label || ''; // Get the dataset label

              if (label) {
                label += ': ';
                label += context.parsed.y
                label += ' lbs'
              }

              return label;
            },

            // Customizes the title
            title: function(context) {
              // Example: use the x-axis value as the title
              return context[0].label;
            }
          }
        }
      }
    },
  });
}
