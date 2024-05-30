package chart

import (
	"fmt"
	"html/template"
	helperFuncs "pkg/helper"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func AppStatBarChart(data BarChartData) template.HTML {
	bar := charts.NewBar()

	bar.Renderer = newchartRenderer(bar, bar.Validate)

	bar.SetGlobalOptions(charts.WithInitializationOpts(
		opts.Initialization{AssetsHost: "/assets/"},
	))

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("App Screentime for: %s", data.AppName),
			TitleStyle: &opts.TextStyle{
				Color:      "black",
				FontWeight: "bolder",
				FontSize:   20,
				FontFamily: "system-ui",
			},
			Subtitle: fmt.Sprintf("from %s - %s %s, %s. Total Uptime of %s", data.XAxis[0], data.XAxis[len(data.XAxis)-1], data.Month, data.Year, helperFuncs.UsageTimeInHrsMin(data.TotalUptime)),
			SubtitleStyle: &opts.TextStyle{
				Color:      "black",
				FontWeight: "bold",
				FontSize:   14,
				FontFamily: "system-ui",
			},
			Left: "center",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      opts.Bool(true),
			Trigger:   "axis",
			TriggerOn: "mousemove",
			AxisPointer: &opts.AxisPointer{
				Type: "cross",
			},
			// Formatter: fmt.Sprintf("{b} %s, %s. <br/> {a}: {c}Hrs", data.Month, data.Year),
			// Formatter: ,
		}),
	)
	bar.SetXAxis(data.XAxis).
		AddSeries(fmt.Sprintf("%sUptime", data.AppName+" "), generateBarItems(data.YAxis, data.XAxis)).SetSeriesOptions()

	return renderToHtml(bar)
}
