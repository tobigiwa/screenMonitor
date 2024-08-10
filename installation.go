//go:build install

package main

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var dependencies = map[string]string{
	"pidof":       "Please install `pidof` on your local machine to continue installationg, e.g for Ubuntu: sudo apt-get pidof",
	"wails":       "You'll need the wails cli to build the desktop application, run `go install github.com/wailsapp/wails/v2/cmd/wails@latest` or visit https://wails.io/docs/gettingstarted/installation",
	"templ":       "You'll need the templ cli to compile the frontend, run `go install github.com/a-h/templ/cmd/templ@latest` or visit https://templ.guide/quick-start/installation",
	"tailwindcss": "You'll need the tailwindcss cli to compile the frontend, visit https://tailwindcss.com/blog/standalone-cli and ensure to rename the binary to `tailwindcss`",
}

func main() {

	if runtime.GOOS != "linux" {
		fmt.Println("This program is only supported on Linux for now")
		os.Exit(1)
	}

	if !isXorgRunning() {
		fmt.Println("Xorg server is not running, ensure you are running this program on a graphical environment using X11.")
		os.Exit(1)
	}

	var dependenciesUnmet bool
	for dep, msg := range dependencies {
		if !checkBinaryExists(dep) {
			fmt.Println(msg)
			dependenciesUnmet = true
		}
	}

	if dependenciesUnmet {
		os.Exit(1)
	}

	fmt.Println("You have all neccessary dependencies...")

	smDaemonDesktopFile := fmt.Sprintf(daemonAutoStart, getGOPATH())
	smDaemonTrayIconDesktopFile := fmt.Sprintf(trayIconAutoStart, getGOPATH())

	for fileName, data := range map[string][]byte{"smDaemon.desktop": []byte(smDaemonDesktopFile), "smTrayIcon.desktop": []byte(smDaemonTrayIconDesktopFile)} {
		if err := writeDesktopFile(fileName, data); err != nil {
			fmt.Printf("Failed to write desktop file for %s: err : %s\n", fileName, err.Error())
			os.Exit(1)
		}

		if err := os.Rename(fileName, autoStartDir(fileName)); err != nil {
			fmt.Printf("Failed to move %s to %s: err : %s\n", fileName, autoStartDir(""), err.Error())
			os.Exit(1)
		}
	}

	fmt.Printf("created a two .desktop files (smDaemon.desktop, smTrayIcon.desktop) and move them to %s/\n", autoStartDir(""))
	os.Exit(0)
}

func checkBinaryExists(binary string) bool {
	_, err := exec.LookPath(binary)
	return err == nil
}

func isXorgRunning() bool {
	cmd := exec.Command("pidof", "Xorg")
	err := cmd.Run()
	return err == nil
}

var daemonAutoStart = `[Desktop Entry]
Type=Application
Exec=%s/bin/smDaemon restart
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
Name[en_US]=smDaemon
Name=smDaemon
Comment[en_US]=ME...
Comment=ME...
`

var trayIconAutoStart = `[Desktop Entry]
Type=Application
Exec=%s/bin/trayIcon
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
Name[en_US]=smDesktop
Name=smDesktop
Comment[en_US]=ME...
Comment=ME...
`

func getGOPATH() string {
	var gopath string
	if gopath = build.Default.GOPATH; gopath == "" {
		if gopath = os.Getenv("GOPATH"); gopath == "" {
			os.Exit(1)
		}
	}
	return gopath
}

func writeDesktopFile(fileName string, data []byte) error {
	return os.WriteFile(fileName, data, 0644)
}

func autoStartDir(fileName string) string {
	if runtime.GOOS == "linux" { // xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		return filepath.Join(os.Getenv("HOME"), ".config", "autostart", fileName)
	}

	if runtime.GOOS == "windows" {
	}

	if runtime.GOOS == "darwin" {
	}

	return ""
}
