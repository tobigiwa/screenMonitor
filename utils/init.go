package utils

import (
	"log"
	"os"
	"path/filepath"
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("error at init fn:", err) // exit
	}

	configDir := filepath.Join(homeDir, "liScreMon")
	logDir := filepath.Join(configDir, "logs")

	for _, dir := range [2]string{configDir, logDir} {
		if err = os.MkdirAll(dir, 0755); err != nil {
			log.Fatalln("error at init fn:", err) // exit
		}
	}

	jsonConfigFile := filepath.Join(configDir, "config.json")
	file, err := os.Create(jsonConfigFile)
	if err != nil {
		log.Fatalln("error at init fn:", err) // exit
	}

	byteData, err := EncodeJSON(ConfigFile{Name: "LiScreMon", Description: "Linux Screen Monitoring", Version: "1.0.0"})
	if err != nil {
		log.Fatalln("error at init fn:", err) // exit
	}

	file.Write(byteData)
	file.Close()

	APP_LOGO_FILE_PATH = filepath.Join(configDir, "liscremon.jpeg")
	APP_CONFIG_DIR = configDir
	APP_LOGS_DIR = logDir
	APP_JSON_CONFIG_FILE_PATH = jsonConfigFile
}

var (
	APP_LOGO_FILE_PATH,
	APP_CONFIG_DIR,
	APP_LOGS_DIR,
	APP_JSON_CONFIG_FILE_PATH string
)
