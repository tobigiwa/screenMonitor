package frontend

import (
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

func writeImageToFile(imageData []byte, filename string) (string, error) {
	file, err := os.Create("/tmp/" + filename + ".png")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(imageData)
	if err != nil {
		return "", err
	}

	return "", nil
}
