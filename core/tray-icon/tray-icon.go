package trayicon

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

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

func setIcon() {
	iconBytes, err := os.ReadFile("./assets/img/twig.ico")
	if err != nil {
		panic(err)
	}
	systray.SetIcon(iconBytes)
}
func setTitle() {
	systray.SetTitle("Nest")
	//systray.SetTooltip("Nest @" + VERSION)
	systray.SetTooltip("Nest")
}

// use this to dynamically reload menu items.
// possible impl.:?
// type func MenuItemPopulateFunc() (title string, tooltip string, func())
// func setMenuItems(...MenuItemPopulateFunc) {}
func setMenuItems() {
	mQuit := systray.AddMenuItem("Quit", "close nest background tasks and exit")
	mConfig := systray.AddMenuItem("Config", "open nest config")

	mRefresh := systray.AddMenuItem("Try Refresh?"+time.Now().String(), "test")
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
			case <-mRefresh.ClickedCh:
				systray.ResetMenu()
				fmt.Println("refreshing...")
				time.Sleep(100 * time.Millisecond)
				defer setMenuItems()
				return
			}
		}
		//set:
		//	setMenuItems()
	}()
}
func onReady() {
	setIcon()
	setTitle()
	setMenuItems()
}

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
