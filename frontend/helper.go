package frontend

import (
	"os"
	"time"
)

func monthDropDownSelectArray() [3]string {
	today := time.Now()
	past4Month := [3]string{}
	for i := 0; i < 3; i++ {
		m := today.AddDate(0, -(i + 1), 0)
		past4Month[i] = m.Month().String()
	}
	return past4Month
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
