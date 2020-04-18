package cmd

import (
	"github.com/spf13/cobra"
	"go-admin/cmd/api"
	"go-admin/cmd/migrate"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:               "go-admin",
	Short:             "-v",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	Long:              `go-admin`,
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		usageStr := `go-admin 0.0.1`
		log.Printf("%s\n", usageStr)
	},
}

func init() {
	rootCmd.AddCommand(api.StartCmd)
	rootCmd.AddCommand(migrate.StartCmd)
}

//Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
