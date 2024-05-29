package chart

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
)

const baseTpl = `
<div class="container">
    <div class="item" id="{{ .ChartID }}" style="height:{{ .Initialization.Height }};"></div>
</div>
{{- range .JSAssets.Values }}
   <script src="{{ . }}"></script>
{{- end }}
<script type="text/javascript">
    "use strict";
    let goecharts_{{ .ChartID | safeJS }} = echarts.init(document.getElementById('{{ .ChartID | safeJS }}'), "{{ .Theme }}");
    let option_{{ .ChartID | safeJS }} = {{ .JSON }};
    goecharts_{{ .ChartID | safeJS }}.setOption(option_{{ .ChartID | safeJS }});
    {{- range .JSFunctions.Fns }}
    {{ . | safeJS }}
    {{- end }}

	window.addEventListener('resize', function() {
		goecharts_{{ .ChartID | safeJS }}.resize();
	});
</script>
`

const baseTpl2 = `
<div class="container">
    <div class="item" id="{{ .ChartID }}" style="height:{{ .Initialization.Height }};"></div>
</div>
{{- range .JSAssets.Values }}
   <script src="{{ . }}"></script>
{{- end }}
<script type="text/javascript">
    "use strict";
    let goecharts_{{ .ChartID | safeJS }} = echarts.init(document.getElementById('{{ .ChartID | safeJS }}'), "{{ .Theme }}");
    let option_{{ .ChartID | safeJS }} = {{ .JSON }};
    goecharts_{{ .ChartID | safeJS }}.setOption(option_{{ .ChartID | safeJS }});
    {{- range .JSFunctions.Fns }}
    {{ . | safeJS }}
    {{- end }}

	window.addEventListener('resize', function() {
		goecharts_{{ .ChartID | safeJS }}.resize();
	});

	goecharts_{{ .ChartID | safeJS }}.setOption(%s);
</script>`

const firstOption = `{
	tooltip: {
		backgroundColor: 'rgba(50,50,50,0.01)',
		borderColor: '#000000',
		textStyle: {
			color: '#000000'
		},
		valueFormatter: function(value) {
			var hours = Math.floor(value);
			var minutes = Math.round((value - hours) * 60);
			return hours + 'Hrs:' + minutes + 'Mins';
				}
		},
		yAxis: {
			type: 'value',
			axisLabel: {
				fontSize: 15,
				align: 'middle',
				fontWeight: 'bolder'
			}
		}
	}`

var tplOne string = fmt.Sprintf(baseTpl2, firstOption)

func (r *chartRenderer) Render(w io.Writer) error {
	const tplName = "chart"
	for _, fn := range r.before {
		fn()
	}

	tpl := template.
		Must(template.New(tplName).
			Funcs(template.FuncMap{
				"safeJS": func(s interface{}) template.JS {
					return template.JS(fmt.Sprint(s))
				},
			}).
			Parse(tplOne),
		)

	err := tpl.ExecuteTemplate(w, tplName, r.c)
	return err
}

type Renderer interface {
	Render(io.Writer) error
}

func renderToHtml(c render.Renderer) template.HTML {
	var buf bytes.Buffer
	err := c.Render(&buf)
	if err != nil {
		log.Printf("Failed to render chart: %s", err)
		return ""
	}

	return template.HTML(buf.String())
}

type chartRenderer struct {
	c      interface{}
	before []func()
}

func newchartRenderer(c interface{}, before ...func()) render.Renderer {
	return &chartRenderer{c: c, before: before}
}

// To satisfy the interface
func (r *chartRenderer) RenderContent() []byte {
	return []byte{}
	// To satisfy the interface
}
func (r *chartRenderer) RenderSnippet() render.ChartSnippet {
	return render.ChartSnippet{}
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
