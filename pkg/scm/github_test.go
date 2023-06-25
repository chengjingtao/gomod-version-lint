package scm

import (
	"context"
	"fmt"
	pkgctx "github.com/chengjingtao/gomod-version-lint/pkg/context"
	gogithub "github.com/google/go-github/v53/github"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestGithubClient_RefreshReviewComments(t *testing.T) {
	// GITHUB_TOKEN=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
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("skip TestGithubClient_RefreshReviewComments")
		return
	}

	var (
		repoPath = "chengjingtao/alauda-ci"
		prID     = 8
		commitID = "c2231e202d16b715f8c1cbedd6b469bc2d72accf"
	)
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	ctx = pkgctx.WithLogger(ctx, logger.Sugar())

	client := gogithub.NewTokenClient(ctx, token)

	github := githubClient{
		Client: client,
	}

	err := github.RefreshReviewComments(ctx, repoPath, prID, RefreshReviewCommentOptions{
		CommentBy: "gomod-version-lint",
		CommitID:  commitID,
		Comments: []ReviewComment{
			{
				Body: "comment-1",
				Path: "README.md",
				Line: 9,
			},
			{
				Body: "comment-2",
				Path: "README.md",
				Line: 9,
			},
			{
				Body: "comment-3",
				Path: "README.md",
				Line: 10,
			},
		},
	})
	if err != nil {
		t.Errorf("error to refreshh comment")
	}
}
