package webserver

import (
	"fmt"
	"log"
	"pkg/chart"
	"pkg/types"
	views "views/components"

	"github.com/a-h/templ"
)

func prepareHtTMLResponse(msg types.Message) templ.Component {
	switch msg.Endpoint {
	case "weekStat":
		return weekStatResponse(msg.WeekStatResponse)
	case "appStat":
		return appStatResponse(msg.AppStatResponse)
	case "dayStat":
		return dayStatResponse(msg.DayStatResponse)
	default:
		return templ.NopComponent
	}
}

func weekStatResponse(w types.WeekStatMessage) templ.Component {
	if len(w.FormattedDay) != len(w.Values) {
		log.Fatal(fmt.Errorf("length of YAxis and XAxis must be equal, instead YAxis - %d and XAxis - %d", len(w.Values), len(w.FormattedDay)))
	}
	return views.WeekStatTempl(
		chart.WeekStatBarChart(chart.BarChartData{
			XAxis:       w.FormattedDay[:],
			YAxis:       w.Values[:],
			Keys:        w.Keys[:],
			Month:       w.Month,
			Year:        w.Year,
			TotalUptime: w.TotalWeekUptime,
		}),
		w.TotalWeekUptime,
		w.AppDetail,
		w.AllCategory,
		w.Keys[6],
	)
}

func appStatResponse(w types.AppStatMessage) templ.Component {
	if len(w.FormattedDay) != len(w.Values) {
		log.Fatal(fmt.Errorf("length of YAxis and XAxis must be equal, instead YAxis - %d and XAxis - %d", len(w.Values), len(w.FormattedDay)))
	}
	return views.AppStatTempl(
		chart.AppStatBarChart(
			chart.BarChartData{
				AppName:     w.AppInfo.AppName,
				YAxis:       w.Values,
				XAxis:       w.FormattedDay,
				Month:       w.Month,
				Year:        w.Year,
				TotalUptime: w.TotalRangeUptime,
			}),
		w.AppInfo)
}

func dayStatResponse(w types.DayStatMessage) templ.Component {
	return views.DayStatTempl(
		chart.DayStatPieChart(
			chart.PieChartData{
				PieData:  w.EachApp,
				DayTotal: w.DayTotal,
				Date:     w.Date,
			}),
		nil,
	)
}
