package commandline

import (
	"os"

	f "github.com/eissar/nest/format"
	"github.com/spf13/cobra"

	"github.com/eissar/nest/api"
	"github.com/eissar/nest/core"
)

// appCmdOpts groups options specific to the “app” commands.
type apiCmdOpts struct {
	format f.FormatType
}

func SubCmdApi() *cobra.Command {
	var apiCmd = &cobra.Command{Use: "api"}

	// var format string
	// apiCmd.PersistentFlags().StringVarP(&format, "format", "o", "json", "output format")

	apiCmd.AddCommand(api.ApplicationCmds()...)
	return apiCmd
}

func CmdCobra() {
	// This variable will hold the value from the --limit flag.
	var rootCmd = &cobra.Command{Use: "nest"}
	rootCmd.AddCommand(Adds())
	rootCmd.AddCommand(CmdAdd())
	rootCmd.AddCommand(Config())
	// rootCmd.AddCommand(Folder())
	rootCmd.AddCommand(api.FolderCmd())
	rootCmd.AddCommand(List())

	rootCmd.AddCommand(RecentLibraries())
	rootCmd.AddCommand(api.LibraryCmd())

	rootCmd.AddCommand(Reveal())
	rootCmd.AddCommand(Shutdown())
	rootCmd.AddCommand(Switch())
	rootCmd.AddCommand(SubCmdApi())
	rootCmd.AddCommand(
		&cobra.Command{
			Use: "start",
			Run: func(cmd *cobra.Command, args []string) {
				core.Start() // blocking
			},
		})

	if err := rootCmd.Execute(); err != nil {
		// Cobra prints the error, so we just need to exit.
		os.Exit(1)
	}
}
