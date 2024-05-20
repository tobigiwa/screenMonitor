module browser

go 1.22.1

require views v0.0.0

replace views v0.0.0 => ../frontend

require pkg v0.0.0

replace pkg v0.0.0 => ../pkg

require (
	github.com/a-h/templ v0.2.663 // indirect
	github.com/go-echarts/go-echarts/v2 v2.3.3 // indirect
)
