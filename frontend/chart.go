package frontend

// On the hope that someone would make sense of go-echart.
// I had too much high hope, tried and tried lotta times to
// use it, frustrated the hell deep-into me (yeah...maybe I'm a dullard,
// buh the best glories in software are tools/library usable by dullard).
// So in my opinion, it was a shitty experience, great docs by the way but
//  a good show on how to make the damn thing responsive should have been extra-ly simplified.

// import (
// 	"context"
// 	"io"
// 	"math/rand"

// 	"github.com/a-h/templ"
// 	"github.com/go-echarts/go-echarts/v2/charts"
// 	"github.com/go-echarts/go-echarts/v2/opts"
// )

// type Renderable interface {
// 	Render(w io.Writer) error
// }

// func ConvertChartToTemplComponent(chart Renderable) templ.Component {
// 	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
// 		return chart.Render(w)
// 	})
// }

// func generateBarItems() []opts.BarData {
// 	items := make([]opts.BarData, 0)
// 	for i := 0; i < 7; i++ {
// 		items = append(items, opts.BarData{Value: rand.Intn(300)})
// 	}
// 	return items
// }

// func createBarChart() *charts.Bar {
// 	bar := charts.NewBar()
// 	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
// 		Title:    "Bar chart",
// 		Subtitle: "That works well with templ",
// 	}))
// 	bar.SetXAxis([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
// 		AddSeries("Category A", generateBarItems()).
// 		AddSeries("Category B", generateBarItems())
// 	return bar
// }
