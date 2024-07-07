module browser

go 1.22.1

require (
	github.com/a-h/templ v0.2.731
	github.com/google/uuid v1.6.0
	pkg v0.0.0
	views v0.0.0
)

require (
	github.com/BurntSushi/xgb v0.0.0-20210121224620-deaf085860bc // indirect
	github.com/gen2brain/beeep v0.0.0-20240516210008-9c006672e7f4 // indirect
	github.com/go-echarts/go-echarts/v2 v2.4.0-rc2 // indirect
	github.com/go-toast/toast v0.0.0-20190211030409-01e6764cf0a4 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/tadvi/systray v0.0.0-20190226123456-11a2b8fa57af // indirect
	golang.org/x/sys v0.19.0 // indirect
)

replace views v0.0.0 => ../frontend

replace pkg v0.0.0 => ../pkg
