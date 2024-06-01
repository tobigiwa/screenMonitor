package frontend

import (
	"errors"
	"io/fs"
	"math/rand"
	"os"
	"time"
)

var (
	art1 string = `
┓ •┏┓     ┳┳┓
┃ ┓┗┓┏┏┓┏┓┃┃┃┏┓┏┓
┗┛┗┗┛┗┛ ┗ ┛ ┗┗┛┛┗`

	art2 string = `
╦ ┬ ╔═╗┌─┐┬─┐┌─┐╔╦╗┌─┐┌┐┌
║ │ ╚═╗│  ├┬┘├┤ ║║║│ ││││
╩═╝┴╚═╝└─┘┴└─└─┘╩ ╩└─┘┘└┘`

	art3 string = `
+-+-+-+-+-+-+-+-+-+
|L|i|S|c|r|e|M|o|n|
+-+-+-+-+-+-+-+-+-+`
)

var arrayOfArt = []string{art1, art2, art3}

func getRandomArt() string {
	randonIndex := rand.Intn(len(arrayOfArt))
	return arrayOfArt[randonIndex]
}

func monthDropDownSelectArray(n int) []string {
	today := time.Now()
	firstDayOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	past3Month := make([]string, n)
	for i := 0; i < n; i++ {
		m := firstDayOfMonth.AddDate(0, -(i + 1), 0)
		past3Month[i] = m.Month().String()
	}
	return past3Month
}

func writeImageToFile(imageData []byte, filename string) bool {
	var (
		tmpFolder = "/tmp/LiScreMon/"
		file      *os.File
		err       error
	)
	if _, err = os.Stat(tmpFolder); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if err = os.MkdirAll(tmpFolder, 0755); err != nil {
				return false
			}
			return false
		}
	}
	if file, err = os.Create(tmpFolder + filename + ".png"); err != nil {
		return false
	}
	defer file.Close()

	if _, err = file.Write(imageData); err != nil {
		return false
	}

	return true
}

func getImageFilePath(filename string) string {
	if _, err := os.Stat("/tmp/LiScreMon/" + filename + ".png"); err != nil {
		return "/assets/image/noAppImage.jpg"
	}
	return "/tmp/LiScreMon/" + filename + ".png"
}
