package cmd

import (
	"context"
	"github.com/chengjingtao/gomod-version-lint/pkg"
	"go.uber.org/zap"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gomod-version-lint",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	ctx := context.Background()
	logger, _ := zap.NewProduction()
	log := logger.Sugar()
	ctx = pkg.WithLogger(ctx, log)

	rootCmd.AddCommand(NewBranchesCmd(ctx))
}
