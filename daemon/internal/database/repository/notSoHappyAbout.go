// This file was named "notSoHappyAbout.go" because it violates
// idiomatic Go code, the package should not have ay x11 related
// shit, just ontop that another X connection.
// This is shitty code right here, Please help!!!.
package repository

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/pkg/errors"
)

var (
	xconnDb      *xgbutil.XUtil
	isXconnDbSet bool
)

func notSoHappyAbout() {
	var err error
	if xconnDb, err = xgbutil.NewConn(); err != nil {
		isXconnDbSet = false
		return
	}
	isXconnDbSet = true
}

func GetWmIcon(windowID xproto.Window) ([]byte, error) {

	if !isXconnDbSet {
		notSoHappyAbout()
	}

	icons, err := ewmh.WmIconGet(xconnDb, windowID)
	if err != nil {
		return nil, err
	}

	if len(icons) == 0 {
		return nil, errors.New("no icon")
	}
	return wmIcon(icons[0])
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
