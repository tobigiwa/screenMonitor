// NOTE: The `database package` is used by all other packages in
// daemon/internal, such it should be independent.
package database

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"strings"
	"time"
	utils "utils"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	badger "github.com/dgraph-io/badger/v4"
	"github.com/pkg/errors"
)

func (bs *BadgerDBStore) WriteUsage(data utils.ScreenTime) error {
	return bs.db.Update(func(txn *badger.Txn) error {

		var (
			app     AppInfo
			valCopy []byte
		)

		item, err := txn.Get(dbAppKey(data.AppName))

		if newApp := errors.Is(err, badger.ErrKeyNotFound); err != nil {
			if !newApp { // there is an error AND is NOT ErrkeyNotFound
				return err
			}

			app.AppName = data.AppName
			app.ScreenStat = make(dailyAppScreenTime)

			addAppInfoForNewApp(data.WindowID, &app)
			fmt.Printf("New appName:%v, time so far is: %v:%v\n\n", app.AppName, app.IsCategorySet, app.IsIconSet)
			return updateAppStats(data, &app, txn)
		}

		if valCopy, err = item.ValueCopy(nil); err != nil {
			return err
		}

		if app, err = utils.DecodeJSON[AppInfo](valCopy); err != nil {
			return err
		}

		updateAppInfoForOldApp(data.WindowID, &app)
		fmt.Printf("Existing appName:%v, time so far is: %v:%v, brought in %f\n\n", data.AppName, app.ScreenStat[utils.Today()].Active, app.ScreenStat[utils.Today()].Open, data.Duration)
		return updateAppStats(data, &app, txn)
	})
}

func updateAppStats(data utils.ScreenTime, app *AppInfo, txn *badger.Txn) error {

	todayStat, ok := app.ScreenStat[utils.Today()]

	if !ok { // we live to see a new day!!! 😎😎😎
		now := time.Now()
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		newDay := now.Sub(midnight)
		yesterdayDuration := data.Duration - newDay.Hours()

		if yesterdayDuration > 0 {
			updateYesterday(data.Type, app, yesterdayDuration)
			data.Duration = newDay.Hours()
		}
	}

	switch data.Type {
	case utils.Active:
		todayStat.Active += data.Duration
		todayStat.ActiveTimeData = append(todayStat.ActiveTimeData, data.Interval)
	case utils.Inactive:
		todayStat.Inactive += data.Duration
	case utils.Open:
		todayStat.Open += data.Duration
	}

	app.ScreenStat[utils.Today()] = todayStat

	byteData, err := utils.EncodeJSON(app)
	if err != nil {
		return err
	}
	return txn.Set(dbAppKey(data.AppName), byteData)
}

func updateYesterday(screenType utils.ScreenType, app *AppInfo, yesterdayDuration float64) {

	yesterdayStat := app.ScreenStat[yesterday()]
	switch screenType {
	case utils.Active:
		yesterdayStat.Active += yesterdayDuration
	case utils.Inactive:
		yesterdayStat.Inactive += yesterdayDuration
	case utils.Open:
		yesterdayStat.Open += yesterdayDuration
	}
	app.ScreenStat[yesterday()] = yesterdayStat
}

func addAppInfoForNewApp(windowId xproto.Window, app *AppInfo) {
	if icon, err := GetWmIcon(windowId); err == nil {
		app.Icon = icon
		app.IsIconSet = true
	}

	if r, err := getDesktopCategoryAndCmd(app.AppName); err == nil {
		app.DesktopCategories = r.desktopCategories
		if len(r.desktopCategories) != 0 {
			app.DesktopCategories = r.desktopCategories
			for _, c := range r.desktopCategories {
				if category, ok := utils.CategoryMap[strings.ToLower(c)]; ok {
					app.Category = category
					app.IsCategorySet = true
					break
				}
			}
		}
		app.CmdLine = r.cmdLine
		app.IsCmdLineSet = true
		fmt.Println("fetched info for new app", app.AppName, app.CmdLine, app.DesktopCategories)
	}
}

func updateAppInfoForOldApp(windowId xproto.Window, app *AppInfo) {
	if !app.IsIconSet {
		if icon, err := GetWmIcon(windowId); err == nil {
			app.Icon = icon
			app.IsIconSet = true
		}
	}

	if !app.IsCmdLineSet && !app.IsCategorySet {
		if r, err := getDesktopCategoryAndCmd(app.AppName); err == nil {
			if r.cmdLine != "" {
				app.CmdLine = r.cmdLine
				app.IsCmdLineSet = true
			}

			if len(r.desktopCategories) != 0 {
				app.DesktopCategories = r.desktopCategories
				for _, c := range r.desktopCategories {
					fmt.Printf("currently in category selection for app %s with c '%s'\n", app.AppName, strings.ToLower(c))
					if category, ok := utils.CategoryMap[strings.ToLower(c)]; ok {
						app.Category = category
						app.IsCategorySet = true
						break
					}
				}
			}
			fmt.Println("fetched info for old app", app.AppName, app.CmdLine, app.DesktopCategories, app.Category)
		}
	}
}

func GetWmIcon(windowID xproto.Window) ([]byte, error) {

	icons, err := ewmh.WmIconGet(utils.X11Connection, windowID)
	if err != nil {
		return nil, err
	}

	if len(icons) == 0 {
		return nil, errors.New("no icon")
	} else if len(icons) == 1 {
		return wmIcon(icons[0])
	} else {
		lastIconIndex := len(icons) - 1 // it is usually more clear
		return wmIcon(icons[lastIconIndex])
	}
}

func wmIcon(icon ewmh.WmIcon) ([]byte, error) {

	img := image.NewRGBA(image.Rect(0, 0, int(icon.Width), int(icon.Height)))
	for i, u := range icon.Data {
		x := i % int(icon.Width)
		y := i / int(icon.Width)
		r := uint8(u >> 16 & 0xFF)
		g := uint8(u >> 8 & 0xFF)
		b := uint8(u & 0xFF)
		a := uint8(u >> 24 & 0xFF)
		img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: a})
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
