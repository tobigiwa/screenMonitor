package webserver

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
)

func createPieChart() *charts.Pie {
	pie := charts.NewPie()

	pie.SetGlobalOptions(charts.WithInitializationOpts(
		opts.Initialization{AssetsHost: "/assets/"},
	))
	pie.Renderer = newchartRenderer(pie, pie.Validate)

	pieData := []opts.PieData{
		{Name: "Dead Cases", Value: 123},
		{Name: "Recovered Cases", Value: 456},
		{Name: "Active Cases", Value: 789},
	}
	pie.AddSeries("Case Distribution", pieData).SetSeriesOptions(
		charts.WithLabelOpts(opts.Label{Show: true, Formatter: "{b}: {c}"}),
	)

	return pie
}

// ###################################################THE BAR CHART##########################################################################3

type barChartData struct {
	xAxis       [7]string
	yAxis       [7]float64
	month       string
	totalUptime float64
}

func weekStatBarChart(data barChartData) template.HTML {
	bar := charts.NewBar()

	bar.SetGlobalOptions(charts.WithInitializationOpts(
		opts.Initialization{AssetsHost: "/assets/"},
	))
	bar.Renderer = newchartRenderer(bar, bar.Validate)

	hrs, min := hrsAndMinute(data.totalUptime)

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Weekly Screentime",
			TitleStyle: &opts.TextStyle{
				Color:      "black",
				FontStyle:  "bold",
				FontSize:   100,
				FontFamily: "system-ui",
			},
			Subtitle: fmt.Sprintf("from %s - %s %s. Total Uptime of %dHr:%dMin", data.xAxis[0], data.xAxis[6], data.month, hrs, min),
			SubtitleStyle: &opts.TextStyle{
				Color:      "black",
				FontStyle:  "bold",
				FontSize:   13,
				FontFamily: "system-ui",
			},
			Left: "center",
		}),
		charts.WithLegendOpts(opts.Legend{
			Type: "plain",
		},
		),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:           true,
			Trigger:        "axis",
			TriggerOn:      "mousemove",
			AxisPointer:    &opts.AxisPointer{},
			Formatter:      "Date: {b} <br/> {a0}: {c0} Hrs",
			ValueFormatter: ""},
		),
	)
	bar.SetXAxis(data.xAxis).
		AddSeries("Daily Uptime", generateBarItems(data.yAxis, data.xAxis))
	return renderToHtml(bar)
}

var barColors [7]string = [7]string{
	"rgb(255, 99, 132)",
	"rgb(255, 159, 64)",
	"rgb(255, 205, 86)",
	"rgb(75, 192, 192)",
	"rgb(54, 162, 235)",
	"rgb(153, 102, 255)",
	"rgb(201, 203, 207)"}

var barColorsBackGround [7]string = [7]string{
	"rgba(255, 99, 132, 0.5)",
	"rgba(255, 159, 64, 0.5)",
	"rgba(255, 205, 86, 0.5)",
	"rgba(75, 192, 192, 0.5)",
	"rgba(54, 162, 235, 0.5)",
	"rgba(153, 102, 255, 0.5)",
	"rgba(201, 203, 207, 0.5)",
}

func hrsAndMinute(hr float64) (int, int) {
	return int(hr), int(math.Round((hr - float64(int(hr))) * 60))
}

func generateBarItems(yAxis [7]float64, xAxis [7]string) []opts.BarData {
	items := make([]opts.BarData, 0, 7)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{
			Name:  xAxis[i],
			Value: yAxis[i],
			Label: &opts.Label{
				Color:     "auto",
				FontStyle: "bold",
			},
			ItemStyle: &opts.ItemStyle{
				Color:       barColors[i],
				BorderColor: barColorsBackGround[i],
			},
			Tooltip: &opts.Tooltip{},
		})
	}
	return items
}

func renderToHtml(c render.Renderer) template.HTML {
	var buf bytes.Buffer
	err := c.Render(&buf)
	if err != nil {
		log.Printf("Failed to render chart: %s", err)
		return ""
	}

	return template.HTML(buf.String())
}

type chartRenderer struct {
	c      interface{}
	before []func()
}

func newchartRenderer(c interface{}, before ...func()) render.Renderer {
	return &chartRenderer{c: c, before: before}
}

func (r *chartRenderer) Render(w io.Writer) error {
	const tplName = "chart"
	for _, fn := range r.before {
		fn()
	}

	tpl := template.
		Must(template.New(tplName).
			Funcs(template.FuncMap{
				"safeJS": func(s interface{}) template.JS {
					return template.JS(fmt.Sprint(s))
				},
			}).
			Parse(baseTpl),
		)

	err := tpl.ExecuteTemplate(w, tplName, r.c)
	return err
}

var baseTpl = `
<div class="container">
    <div class="item" id="{{ .ChartID }}" style="width:{{ .Initialization.Width }};height:{{ .Initialization.Height }};"></div>
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
</script>
`
