package chart

type BarChartData struct {
	XAxis       []string
	YAxis       []float64
	Keys        []string
	Month       string
	Year        string
	TotalUptime float64
	AppName     string
}

var barChartColors [7]string = [7]string{
	"rgb(255, 99, 132)",
	"rgb(255, 159, 64)",
	"rgb(255, 205, 86)",
	"rgb(75, 192, 192)",
	"rgb(54, 162, 235)",
	"rgb(153, 102, 255)",
	"rgb(201, 203, 207)"}

var barChartColorsBackGround [7]string = [7]string{
	"rgba(255, 99, 132, 0.5)",
	"rgba(255, 159, 64, 0.5)",
	"rgba(255, 205, 86, 0.5)",
	"rgba(75, 192, 192, 0.5)",
	"rgba(54, 162, 235, 0.5)",
	"rgba(153, 102, 255, 0.5)",
	"rgba(201, 203, 207, 0.5)",
}
