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
		opts.Initialization{AssetsHost: "/assets/libraries/"},
	))

	bar.Renderer = newchartRenderer(
		bar,
		"weekStat",
		fmt.Sprintf(baseTpl, barChartOptions),
		bar.Validate)

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Weekly Screentime",
			TitleStyle: &opts.TextStyle{
				Color:      "black",
				FontWeight: "bolder",
				FontSize:   20,
				FontFamily: "system-ui",
			},
			Subtitle: fmt.Sprintf("from %s - %s %s. Total Uptime of %s", data.XAxis[0], data.XAxis[6], data.Month, helperFuncs.UsageTimeInHrsMin(data.TotalUptime)),
			SubtitleStyle: &opts.TextStyle{
				Color:      "black",
				FontWeight: "bold",
				FontSize:   14,
				FontFamily: "system-ui",
			},
			Left: "center",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(true),
			Left: "20%",
			Data: data.Keys,
			Icon: "triangle",
		}),

		charts.WithToolboxOpts(opts.Toolbox{
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  opts.Bool(true),
					Title: "Download chart as Image",
					Name:  fmt.Sprintf("weekly-screentime-%s-%s-%s-%s", data.XAxis[0], data.XAxis[6], data.Month, data.Year),
				},
			},
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
		AddSeries("Daily uptime", barItems(data.YAxis, data.XAxis)).SetSeriesOptions()
	return renderToHtml(bar)
}
