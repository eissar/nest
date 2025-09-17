package tray

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"fyne.io/systray"
	"github.com/eissar/nest/api"
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
func setMenuItems(libs []string) (*systray.MenuItem, *systray.MenuItem) {
	mQuit := systray.AddMenuItem("Quit", "close nest background tasks and exit")
	mConfig := systray.AddMenuItem("Config", "open nest config")

	//mRefresh := systray.AddMenuItem("Try Refresh?"+time.Now().String(), "test")
	// Sets the icon of a menu item.
	// mQuit.SetIcon(icon.Datacfg.Libraries.Paths)

	mLibraries := systray.AddMenuItem("Libraries", "Libraries")

	for _, l := range libs {
		_ = mLibraries.AddSubMenuItem(l, "")
		// TODO: add switch behavior to these subMenuItems
	}
	if len(libs) == 0 {
		mLibraries.AddSubMenuItem("no library history found.", "")
	}

	return mQuit, mConfig

}
func onReady() {
	cfg := config.GetConfig()

	// fmt.Printf("cfg.Libraries: %v\n", cfg.Libraries)
	setIcon()
	setTitle()

	// TODO: Make this use or abide by cfg.Libraries.AutoLoad preference
	libs, err := api.LibraryHistory(cfg.BaseURL())
	if err != nil {
		fmt.Printf("WARN: no library history could be found. library paths are missing.")
	}

	mQuit, mConfig := setMenuItems(libs)
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
				//case <-mRefresh.ClickedCh:
				//	systray.ResetMenu()
				//	fmt.Println("refreshing...")
				//	time.Sleep(100 * time.Millisecond)
				//	defer setMenuItems()
				//	return
			}
		}
		//set:
		//	setMenuItems()
	}()
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
