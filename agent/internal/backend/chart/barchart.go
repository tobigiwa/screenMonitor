package chart

import (
	"github.com/go-echarts/go-echarts/v2/opts"
)

const baseTpl = `
<div class="container">
    <div class="item" id="{{ .ChartID }}" style="height:30rem;"></div>
</div>
{{- range .JSAssets.Values }}
   <script src="{{ . }}"></script>
{{- end }}
<script type="text/javascript">
    "use strict";
    let goecharts_{{ .ChartID | safeJS }} = echarts.init(document.getElementById('{{ .ChartID | safeJS }}'), "{{ .Theme }}");
    let option_{{ .ChartID | safeJS }} = {{ .JSON }};
    goecharts_{{ .ChartID | safeJS }}.setOption(option_{{ .ChartID | safeJS }});
    {{- range .JSFunctions.Fns }}
    {{ . | safeJS }}
    {{- end }}

	window.addEventListener('resize', function() {
		goecharts_{{ .ChartID | safeJS }}.resize();
	});

	goecharts_{{ .ChartID | safeJS }}.setOption(%s);

	goecharts_{{ .ChartID | safeJS }}.on('click', function (params) {
		const thatDay = goecharts_{{ .ChartID | safeJS }}.getOption().legend[0].data[params.dataIndex]
		htmx.ajax('GET','/dayStat?day=' + thatDay + '&nonce=' + Math.random(), {target:'#echart', swap:'innerHTML'})
	});
	
</script>`

const barChartOptions = `{
	tooltip: {
		backgroundColor: 'rgba(50,50,50,0.01)',
		borderColor: '#000000',
		textStyle: {
			color: '#000000'
		},
		valueFormatter: function(value) {
			var hours = Math.floor(value);
			var minutes = Math.trunc((value - hours) * 60);
			return hours + 'Hrs:' + minutes + 'Mins';
				}
		},
		xAxis: {
			name: 'Days',
			axisLabel: {
				fontSize: 10,
				fontWeight: 'bold',
			},
			nameTextStyle: {
				color: '#000000',
				fontWeight: 'bolder',
				fontSize: 10,
			}
		},
		yAxis: {
			type: 'value',
			axisLabel: {
				fontSize: 10,
				align: 'middle',
				fontWeight: 'bold',
				formatter: '{value}Hr'
			}
		}
	}`

type BarChartData struct {
	XAxis       []string
	YAxis       []float64
	Keys        []string
	Month       string
	Year        string
	TotalUptime float64
	AppName     string
}

func barItems(YAxis []float64, xAxis []string) []opts.BarData {
	items := make([]opts.BarData, 0, 7)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{
			Name:  xAxis[i],
			Value: YAxis[i],
			ItemStyle: &opts.ItemStyle{
				Color:       barChartColors[i],
				BorderColor: barChartColorsBackGround[i],
			},
		})
	}
	return items
}

var barChartColors [7]string = [7]string{
	"rgb(255, 99, 132)",
	"rgb(255, 159, 64)",
	"rgb(255, 205, 86)",
	"rgb(75, 192, 192)",
	"rgb(54, 162, 235)",
	"rgb(153, 102, 255)",
	"rgb(201, 203, 207)"}

var barChartColorsBackGround [7]string = [7]string{
	"rgba(255, 99, 132, 0.5)",
	"rgba(255, 159, 64, 0.5)",
	"rgba(255, 205, 86, 0.5)",
	"rgba(75, 192, 192, 0.5)",
	"rgba(54, 162, 235, 0.5)",
	"rgba(153, 102, 255, 0.5)",
	"rgba(201, 203, 207, 0.5)",
}
