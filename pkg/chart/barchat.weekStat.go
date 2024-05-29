package chart

import (
	"bytes"
	"fmt"
	"html/template"
	helperFuncs "pkg/helper"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
)

func WeekStatBarChart(data BarChartData) template.HTML {
	bar := charts.NewBar()

	bar.SetGlobalOptions(charts.WithInitializationOpts(
		opts.Initialization{AssetsHost: "/assets/"},
	))

	bar.Renderer = newchartRenderer(bar, bar.Validate)

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Weekly Screentime",
			TitleStyle: &opts.TextStyle{
				Color:      "black",
				FontStyle:  "bold",
				FontSize:   100,
				FontFamily: "system-ui",
			},
			Subtitle: fmt.Sprintf("from %s - %s %s. Total Uptime of %s", data.XAxis[0], data.XAxis[6], data.Month, helperFuncs.UsageTimeInHrsMin(data.TotalUptime)),
			SubtitleStyle: &opts.TextStyle{
				Color:      "black",
				FontStyle:  "bold",
				FontSize:   13,
				FontFamily: "system-ui",
			},
			Left: "center",
		}),
		charts.WithLegendOpts(opts.Legend{
			Left:   "0%",
			Orient: "vertical",
			Bottom: "50%",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      opts.Bool(true),
			Trigger:   "axis",
			TriggerOn: "mousemove",
			AxisPointer: &opts.AxisPointer{
				Type: "cross",
			},
		}),
	)
	bar.SetXAxis(data.XAxis).
		AddSeries("Uptime in Hours", generateBarItems(data.YAxis, data.XAxis)).SetSeriesOptions()
	return renderToHtml(bar)
}

func renderToHtml2(chart render.Renderer) template.HTML {
	var buf bytes.Buffer
	chartSnippet := chart.RenderSnippet()

	tmpl := "{{.Element  }} {{.Script }} {{.Option}}"
	t := template.New("snippet")
	t, err := t.Parse(tmpl)
	if err != nil {
		panic(fmt.Errorf("crash from renderToHtml2:t.Parse error: %w", err))
	}

	// fmt.Printf("chartSnippet\n%+v\n\n", chartSnippet)

	data := struct {
		Element template.HTML
		Script  template.HTML
		Option  template.HTML
	}{
		Element: template.HTML(baseTpl),
		Script:  template.HTML(chartSnippet.Script),
		Option:  template.HTML(chartSnippet.Option),
	}

	err = t.Execute(&buf, data)
	if err != nil {
		panic(fmt.Errorf("crash from renderToHtml2:t.Execute error: %w", err))
	}

	// fmt.Println(buf.String())

	return template.HTML(buf.String())
}
