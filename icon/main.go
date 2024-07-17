package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"runtime"
	"strings"
	"time"
	"utils"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
)

var (
	broswerProcess *os.Process
)

func main() {
	onExit := func() {
		now := time.Now()
		fmt.Println("Exit at", now.String())
	}

	systray.Run(onReady, onExit)
}

func addQuitItem() {
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	mQuit.Enable()
	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()
	systray.AddSeparator()
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("LiScreMon")
	systray.SetTooltip("Linux Screen Monitor")
	addQuitItem()

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("LiScreMon")
		systray.SetTooltip("Linux Screen Monitor")

		launchBrowser := systray.AddMenuItem("Launch browser view", "Launch browser view")
		launchBrowserSubOne := launchBrowser.AddSubMenuItem("View in browser", "")
		launchBrowserSubTwo := launchBrowser.AddSubMenuItem("Close browser view", "")
		launchBrowserSubTwo.Disable()

		launchDesktop := systray.AddMenuItem("Launch desktop view", "Launch desktop view")

		for {
			select {
			case <-launchBrowserSubOne.ClickedCh:
				if strings.Contains(launchBrowser.String(), "running 🟢") {
					jumpToBrowserView()
				} else {
					launchBrowserView()
					launchBrowserSubTwo.Enable()
					launchBrowser.SetTitle("Browser view: running 🟢")
				}

			case <-launchBrowserSubTwo.ClickedCh:
				broswerProcess.Signal(os.Interrupt)
				broswerProcess.Release()
				launchBrowserSubTwo.Disable()
				launchBrowser.SetTitle("Launch browser view")

			case <-launchDesktop.ClickedCh:
				if launcDesktopView() != nil {
					utils.NotifyWithBeep("Operation failed", "Could not launch LiScreMon desktop app.")
				}

			}
		}
	}()
}

func jumpToBrowserView() {
	var (
		portAddres string
		err        error
		cmd        *exec.Cmd
	)

	if portAddres, err = getBrowserRunningAddr(); err != nil {
		utils.NotifyWithBeep("Operation failed", "Could not launch LiScreMon broswer view.")
		fmt.Println(err)
		return
	}

	if runtime.GOOS == "linux" {
		if cmd = exec.Command("xdg-open", portAddres); cmd.Start() != nil {
			utils.NotifyWithBeep("Operation failed", "Could not launch LiScreMon broswer view.")
			fmt.Println(err)
			return
		}

		cmd.Wait()
	}

	if runtime.GOOS == "windows" {
		notImplemented()
	}
}

func launchBrowserView() {
	var cmd *exec.Cmd

	gopath := getGOPATH()

	if runtime.GOOS == "linux" {
		gopathBin := filepath.Join(gopath, "bin", "browser")
		cmd = exec.Command(gopathBin)
	}

	if runtime.GOOS == "windows" {
		notImplemented()
	}

	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}

	broswerProcess = cmd.Process
}

func launcDesktopView() error {
	var cmd *exec.Cmd

	gopath := getGOPATH()

	if runtime.GOOS == "linux" {
		gopathBin := filepath.Join(gopath, "bin", "desktop")
		cmd = exec.Command(gopathBin)
	}

	if runtime.GOOS == "windows" {
		notImplemented()
	}

	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func notImplemented() {}

func getBrowserRunningAddr() (string, error) {
	byteData, err := os.ReadFile(utils.APP_JSON_CONFIG_FILE_PATH)
	if err != nil {
		return "", err
	}
	config, err := utils.DecodeJSON[utils.ConfigFile](byteData)
	if err != nil {
		return "", err
	}
	return config.BrowserAddr, nil
}

func getGOPATH() string {
	var gopath string
	if gopath = build.Default.GOPATH; gopath == "" {
		if gopath = os.Getenv("GOPATH"); gopath == "" {
			log.Fatalln("cannot build program, unable to determine GOPATH")
		}
	}
	return gopath
}
