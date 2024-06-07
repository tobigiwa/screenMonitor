module browser

go 1.22.1

require (
	github.com/a-h/templ v0.2.707
	github.com/google/uuid v1.6.0
	pkg v0.0.0
	views v0.0.0
)

require (
	github.com/BurntSushi/xgb v0.0.0-20210121224620-deaf085860bc // indirect
	github.com/go-echarts/go-echarts/v2 v2.4.0-rc2 // indirect
)

replace views v0.0.0 => ../frontend

replace pkg v0.0.0 => ../pkg
