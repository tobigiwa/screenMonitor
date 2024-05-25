package chart

import (
	"fmt"
	"html/template"
	helperFuncs "pkg/helper"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
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
		charts.WithYAxisOpts(opts.YAxis{
			Name:         "in Hours",
			Type:         "value",
			NameLocation: "end",
			NameGap:      5,
			Scale:        true,
		}),
		charts.WithLegendOpts(opts.Legend{
			Left:   "left",
			Orient: "vertical",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Trigger:   "axis",
			TriggerOn: "mousemove",
			AxisPointer: &opts.AxisPointer{
				Type: "cross",
			},
			Formatter: fmt.Sprintf("{b} %s, %s. <br/> {a}: {c}Hrs", data.Month, data.Year),
		}),
	)
	bar.SetXAxis(data.XAxis).
		AddSeries("Uptime", generateBarItems(data.YAxis, data.XAxis)).SetSeriesOptions()
	return renderToHtml(bar)
}

func generateBarItems(YAxis []float64, xAxis []string) []opts.BarData {
	items := make([]opts.BarData, 0, 7)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{
			Name:  xAxis[i],
			Value: YAxis[i],
			Label: &opts.Label{
				Color:     "auto",
				FontStyle: "bold",
			},
			ItemStyle: &opts.ItemStyle{
				Color:       barChartColors[i],
				BorderColor: barChartColorsBackGround[i],
			},
			Tooltip: &opts.Tooltip{},
		})
	}
	return items
}
