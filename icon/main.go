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
	"utils"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
)

var (
	broswerCmd *exec.Cmd
	desktopCmd *exec.Cmd
)

func main() {
	onExit := func() { systray.Quit() }

	systray.Run(onReady, onExit)
}

// func addQuitItem() {
// 	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
// 	mQuit.Enable()
// 	go func() {
// 		<-mQuit.ClickedCh
// 		fmt.Println("Requesting quit")
// 		systray.Quit()
// 		fmt.Println("Finished quitting")
// 	}()
// 	systray.AddSeparator()
// }

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("LiScreMon")
	systray.SetTooltip("Linux Screen Monitor")
	// addQuitItem()

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("LiScreMon")
		systray.SetTooltip("Linux Screen Monitor")

		launchBrowser := systray.AddMenuItem("Launch browser view", "Launch browser view")
		launchBrowserSubOne := launchBrowser.AddSubMenuItem("View in browser", "")
		launchBrowserSubTwo := launchBrowser.AddSubMenuItem("Close browser view", "")
		launchBrowserSubTwo.Disable()

		systray.AddSeparator()

		launchDesktop := systray.AddMenuItem("Launch desktop view", "Launch desktop view")

		for {
			select {
			case <-launchBrowserSubOne.ClickedCh:
				if strings.Contains(launchBrowser.String(), "running 🟢") {
					if err := jumpToBrowserView(); err != nil { // opens **another** browser tab of the `browser server's` port Addr
						utils.NotifyWithBeep("Operation failed", "Could not launch LiScreMon browser view.")
					}
					continue
				}

				if err := launchBrowserView(); err != nil { // starts the browser server
					fmt.Println(err)
					utils.NotifyWithBeep("Operation failed", "Could not launch LiScreMon browser view.")
					continue
				}

				launchBrowserSubTwo.Enable()
				launchBrowser.SetTitle("Browser view: running 🟢")

			case <-launchBrowserSubTwo.ClickedCh: // closes the browser server
				broswerCmd.Process.Signal(os.Interrupt)
				broswerCmd.Wait()
				launchBrowserSubTwo.Disable()
				launchBrowser.SetTitle("Launch browser view")

			case <-launchDesktop.ClickedCh:
				if err := launcDesktopView(); err != nil { // desktop app is launched
					utils.NotifyWithBeep("Operation failed", "Could not launch LiScreMon desktop view.")
					continue
				}

				launchDesktop.Disable()
				launchDesktop.SetTitle("Desktop view: running 🟢")

				go func() {
					if err := desktopCmd.Wait(); err != nil { // waits for desktop app to be closed
						log.Println("error releasing cmd resource:err ", err)

					}
					launchDesktop.Enable()
					launchDesktop.SetTitle("Launch desktop view")
				}()
			}
		}
	}()
}

func launchBrowserView() error {

	gopath := getGOPATH()

	if runtime.GOOS == "linux" {
		gopathBin := filepath.Join(gopath, "bin", "browser")
		broswerCmd = exec.Command(gopathBin)
	}

	if runtime.GOOS == "windows" {
		notImplemented()
	}

	return broswerCmd.Start()
}

func launcDesktopView() error {
	gopath := getGOPATH()

	if runtime.GOOS == "linux" {
		gopathBin := filepath.Join(gopath, "bin", "desktop")
		desktopCmd = exec.Command(gopathBin)
	}

	if runtime.GOOS == "windows" {
		notImplemented()
	}

	return desktopCmd.Start()
}

func notImplemented() {}

func jumpToBrowserView() error {
	var (
		portAddres string
		err        error
		cmd        *exec.Cmd
	)

	if portAddres, err = getBrowserRunningAddr(); err != nil {
		return err
	}

	if runtime.GOOS == "linux" {
		cmd = exec.Command("xdg-open", portAddres)
	}

	if runtime.GOOS == "windows" {
		notImplemented()
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

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
