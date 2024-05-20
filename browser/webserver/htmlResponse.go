package webserver

import (
	"fmt"
	"pkg/types"
	views "views"

	"github.com/a-h/templ"
)

func prepareHtTMLResponse(msg types.Message) templ.Component {
	switch msg.Endpoint {
	case "weekStat":
		return prepareWeekStatHTMLResponse(msg.WeekStatResponse)
	}

	return templ.NopComponent
}

func prepareWeekStatHTMLResponse(w types.WeekStatMessage) templ.Component {
	for i, v := range w.AppDetail {
		fmt.Println(i, v.AppInfo.AppName, v.Usage, v.AppInfo.IsIconSet)
	}
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
