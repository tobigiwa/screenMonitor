package chart

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"

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
			Parse(baseTpl),
		)

	err := tpl.ExecuteTemplate(w, tplName, r.c)
	return err
}
