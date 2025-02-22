package trayicon

import (
	"log"
	"os"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
	"github.com/eissar/nest/core"
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
	// Sets the icon of a menu item.
	mQuit.SetIcon(icon.Data)

	// event listeners for menu items
	go func() {
		for range mQuit.ClickedCh {
			// just close the systray.
			//the registered onExit func will handle closing logic.
			systray.Quit()
		}
	}()
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
func Run(s *echo.Echo) {

	onExit := func() {
		core.Shutdown(s)
	}

	go systray.Run(onReady, onExit)

}
