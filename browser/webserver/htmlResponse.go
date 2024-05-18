package webserver

import (
	views "views"

	"github.com/a-h/templ"
)

func prepareHtTMLResponse(msg Message) templ.Component {
	switch msg.Endpoint {
	case "weekStat":
		return prepareWeekStatHTMLResponse(msg.WeekStatResponse)
	}

	return templ.NopComponent
}

func prepareWeekStatHTMLResponse(w WeekStatMessage) templ.Component {
	return views.ChartWrapper(
		weekStatBarChart(
			barChartData{
				xAxis:       w.FormattedDay,
				yAxis:       w.Values,
				month:       w.Month,
				totalUptime: w.TotalWeekUptime,
			}))
}
