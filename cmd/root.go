package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"gomod.alauda.cn/gomod-version-lint/options"
)

func NewRootCmd(ctx context.Context) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "gomod-version-lint",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}

	rootOpts := &options.RootOptions{}
	rootCmd.PersistentFlags().BoolVar(&rootOpts.Debug, "debug", false, "enable debug log level")

	rootCmd.AddCommand(NewBranchesCmd(ctx, rootOpts))
	rootCmd.AddCommand(NewCommentCmd(ctx, rootOpts))

	return rootCmd
}
