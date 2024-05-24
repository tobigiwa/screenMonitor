package webserver

import (
	"pkg/types"
	views "views"

	"github.com/a-h/templ"
)

func prepareHtTMLResponse(msg types.Message) templ.Component {
	switch msg.Endpoint {
	case "weekStat":
		return weekStatResponse(msg.WeekStatResponse)
	case "appStat":
		return templ.NopComponent
	}

	return templ.NopComponent
}



func weekStatResponse(w types.WeekStatMessage) templ.Component {
	return views.WeekStatChartAndHighlight(
		weekStatBarChart(barChartData{
			xAxis:       w.FormattedDay,
			yAxis:       w.Values,
			month:       w.Month,
			year:        w.Year,
			totalUptime: w.TotalWeekUptime,
		}),
		w.AppDetail,
	)
}

func appStatResponse() templ.Component
