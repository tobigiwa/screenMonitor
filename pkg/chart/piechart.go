package chart

import (
	"pkg/types"

	"github.com/go-echarts/go-echarts/v2/opts"
)

const piechart = `
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
</script>`

const pieChartOptions = `{
	tooltip: {
		backgroundColor: 'rgba(50,50,50,0.01)',
		borderColor: '#000000',
		textStyle: {
			color: '#000000'
			}
		},
	series: [{
			label:{
				formatter: function(params) {
				var hours = Math.floor(params.data.value); 
				var minutes = Math.round((params.data.value - hours) * 60); 
				var v = hours + 'Hrs:' + minutes + 'Mins';
				return params.name +": "+ v;
				}
			}
	}]			
	}`

type PieChartData struct {
	PieData  []types.AppStat
	DayTotal types.Stats
	Date     string
}

func pieItems(data PieChartData) []opts.PieData {
	pieData := make([]opts.PieData, 0, len(data.PieData))
	for i := 0; i < len(data.PieData); i++ {
		pieData = append(pieData, opts.PieData{
			Name:  data.PieData[i].AppName,
			Value: data.PieData[i].Usage.Active,
		})
	}

	return pieData

}
