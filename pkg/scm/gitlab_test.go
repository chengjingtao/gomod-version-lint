package scm

import (
	"context"
	pkgctx "github.com/chengjingtao/gomod-version-lint/pkg/context"
	gogitlab "github.com/xanzy/go-gitlab"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestGitlabClient_RefreshReviewComments(t *testing.T) {
	// export GITLAB_TOKEN=7uhbhq3PCFMYnx4A_eVe
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		return
	}

	var (
		repoPath = "lab/hello-world"
		prID     = 16
		commitID = "e8384783e42540c5cad5902f11eb8574f639e4fa"
	)

	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	ctx = pkgctx.WithLogger(ctx, logger.Sugar())

	client, err := gogitlab.NewClient(token, gogitlab.WithBaseURL("https://gitlab-ce.alauda.cn/"))
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	gitlab := gitlabClient{
		Client: client,
	}

	err = gitlab.RefreshReviewComments(ctx, repoPath, prID, RefreshReviewCommentOptions{
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
