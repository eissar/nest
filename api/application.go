package api

import (
	"fmt"
	"log"

	"github.com/eissar/nest/api/endpoints"
	"github.com/eissar/nest/config"
	f "github.com/eissar/nest/format"
	"github.com/spf13/cobra"
)

//	TODO:
//
// - [X] /api/application/info

type ApplicationInfoData struct {
	Version           string `json:"version"`
	PreReleaseVersion string `json:"preReleaseVersion,omitempty"`
	BuildVersion      string `json:"buildVersion"`
	ExecPath          string `json:"execPath"`
	Platform          string `json:"platform"`
}

// GET Get detailed information on the Eagle App currently running. In most cases, this could be used to determine whether certain functions are available on the user's device.
// <https://api.eagle.cool/application/info>
func ApplicationInfo(baseUrl string) (ApplicationInfoData, error) {
	ep := endpoints.ApplicationInfo
	uri := baseUrl + ep.Path

	var resp struct {
		EagleResponse                     // Response string `json:"response"`
		Data          ApplicationInfoData `json:"data"`
	}
	err := Request(ep.Method, uri, nil, nil, &resp)
	if err != nil {
		return resp.Data, fmt.Errorf("ApplicationInfo: err=%w", err)
	}
	if resp.Status != "success" {
		return resp.Data, fmt.Errorf("ApplicationInfo: err=%w", ErrStatusErr)
	}

	return resp.Data, nil
}

// provides commands
func ApplicationCmd() *cobra.Command {
	cfg := config.GetConfig()

	var o f.FormatType

	app := &cobra.Command{
		Use: "app",
		// Short: "Manage items",
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Println(cmd.Flags())
		// },
	}
	// return []*cobra.Command{

	func() {
		cmd := &cobra.Command{
			Use:   "info", //
			Short: "Display detailed information about the running Eagle application.",
			Long:  "Retrieves and prints detailed information about the Eagle application currently running. ",
			RunE: func(cmd *cobra.Command, args []string) error {
				v, err := ApplicationInfo(cfg.BaseURL())

				if err != nil {
					log.Fatalf("Application: %v", err)
				}

				f.Format(o, v)
				return nil
			},
		}
		app.AddCommand(cmd)
	}()
	return app
}
