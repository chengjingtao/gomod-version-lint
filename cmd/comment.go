package cmd

import (
	"context"
	"gomod.alauda.cn/gomod-version-lint/options"

	"github.com/spf13/cobra"
)

func NewCommentCmd(ctx context.Context, opts *options.RootOptions) *cobra.Command {
	commentOptions := &options.CommentOptions{
		Context: ctx,
	}

	var commentCmd = &cobra.Command{
		Use:   "comment",
		Short: "add comment on git server pull request",
		Long:  `add comment on git server pull request according comments file, but will delete all old comments to avoid adding times by times`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return commentOptions.Run()
		},
	}

	commentOptions.AddFlags(commentCmd.Flags())
	return commentCmd
}
