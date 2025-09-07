package commandline

import (
	"fmt"
	"path/filepath"

	"github.com/eissar/nest/config"
	"github.com/spf13/cobra"
)

func Config() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Show the path to the configuration file",
		Long: `This subcommand prints the absolute path to the application's configuration file.
It does not accept any arguments.`,
		// Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := config.GetConfigPath()
			fmt.Println(filepath.Join(dir, "config.json"))
			return nil
		},
	}

	return configCmd
}
