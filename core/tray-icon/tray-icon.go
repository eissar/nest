package trayicon

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
	"github.com/eissar/nest/config"
	"github.com/eissar/nest/plugins/launch"
	"github.com/labstack/echo/v4"
)

// may panic
func handleQuit() {
	log.Fatalf("quit requested")
	// script := `Add-Type -AssemblyName System.Windows.Forms && [System.Windows.Forms.MessageBox]::Show("Message", "Title", "OK", "Information")`
	// cmd := exec.Command("pwsh", "-NoProfile", "-c", script)
	//_, err := cmd.CombinedOutput()
}

func onReady() {
	iconBytes, err := os.ReadFile("./assets/img/twig.ico")
	if err != nil {
		panic(err)
	}
	systray.SetIcon(iconBytes)
	systray.SetTitle("Nest")
	//systray.SetTooltip("Nest @" + VERSION)
	systray.SetTooltip("Nest")
	mQuit := systray.AddMenuItem("Quit", "close nest background tasks and exit")
	mConfig := systray.AddMenuItem("Config", "open nest config")
	// Sets the icon of a menu item.
	mQuit.SetIcon(icon.Data)

	// event listeners for menu items
	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			case <-mConfig.ClickedCh:
				cfgPath := filepath.Join(config.GetConfigPath(), "config.json")
				launch.Open(cfgPath)
				return
			}
		}
	}()
}
func test() {
	launch.OpenURI(config.GetConfigPath())

}

//func onExit() { // clean up here }

func Quit() {
	systray.Quit()
}

// make sure to build with flag go build -ldflags -H=windowsgui
// onExit can be emitted by calling systray.Quit.
func RunOld(onExit func()) {
	go systray.Run(onReady, onExit)
}

// param s echo server.
func Run(s *echo.Echo, onExit func()) {
	go systray.Run(onReady, onExit)

}
