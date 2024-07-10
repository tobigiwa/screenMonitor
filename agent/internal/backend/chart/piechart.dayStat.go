package chart

import (
	"fmt"
	"html/template"
	helperFuncs "pkg/helper"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func DayStatPieChart(data PieChartData) template.HTML {
	pie := charts.NewPie()

	pie.SetGlobalOptions(charts.WithInitializationOpts(
		opts.Initialization{AssetsHost: "/assets/libraries/"},
	))

	pie.Renderer = newchartRenderer(
		pie,
		"dayStat",
		fmt.Sprintf(piechart, pieChartOptions),
		pie.Validate)

	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Daily Screentime",
			TitleStyle: &opts.TextStyle{
				Color:      "black",
				FontWeight: "bolder",
				FontSize:   20,
				FontFamily: "system-ui",
			},
			Subtitle: fmt.Sprintf("For %s. Total Uptime of %s", data.Date, helperFuncs.UsageTimeInHrsMin(data.DayTotal.Active)),
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

		charts.WithToolboxOpts(opts.Toolbox{
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  opts.Bool(true),
					Title: "Download chart as Image",
					Name:  fmt.Sprintf("daily-screentime-for-%s", data.Date),
				},
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      opts.Bool(true),
			Trigger:   "item",
			Formatter: types.FuncStr(fmt.Sprintf("<b>%s <br/><i>{b}: ({d}%%)</i></b>", data.Date)),
		}),
	)
	pie.AddSeries("Daily uptime", pieItems(data)).SetSeriesOptions(
		charts.WithLabelOpts(
			opts.Label{
				Show:         opts.Bool(true),
				FontStyle: "italic",
			}),

		charts.WithPieChartOpts(opts.PieChart{
			Radius:   []string{"40", "140"},
			Center:   []string{"55%", "50%"},
			RoseType: "radius",
		}),
		charts.WithEmphasisOpts(opts.Emphasis{
			Label: &opts.Label{
				Show:       opts.Bool(true),
				FontSize:   16,
				FontWeight: "bold",
				FontStyle:  "italic",
			},
			ItemStyle: &opts.ItemStyle{
				ShadowBlur:    10,
				ShadowOffsetX: 0,
				ShadowColor:   "rgba(0, 0, 0, 0.5)",
			},
		}),
	)

	return renderToHtml(pie)
}
