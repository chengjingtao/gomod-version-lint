package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"gomod.alauda.cn/gomod-version-lint/options"
)

func NewBranchesCmd(ctx context.Context, opts *options.RootOptions) *cobra.Command {

	branchOpts := &options.BranchesOptions{
		RootOptions: *opts,
		Context:     ctx,
	}

	cmd := &cobra.Command{
		Use:   "branches",
		Short: "output branches information for each go module dependency",
		Long:  `output branches information for each go module dependency`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return branchOpts.Run()
		},
	}

	branchOpts.AddFlags(cmd.Flags())

	return cmd
}
