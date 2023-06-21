package cmd

import (
	"context"
	"github.com/chengjingtao/gomod-version-lint/options"
	"github.com/spf13/cobra"
)

func NewBranchesCmd(ctx context.Context) *cobra.Command {

	branchOpts := &options.BranchesOptions{
		Context: ctx,
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
