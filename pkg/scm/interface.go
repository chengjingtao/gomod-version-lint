package scm

import (
	"context"
	"fmt"
)

type RefreshReviewCommentOptions struct {
	CommentBy string

	CommitID string
	Comments []ReviewComment
}

func (opts RefreshReviewCommentOptions) FmtCommentBy() string {
	return fmt.Sprintf("<!-- %s -->", opts.CommentBy)
}

func (opts ReviewComment) FmtComment(commentBy string) string {
	return fmt.Sprintf("<!-- %s -->\n%s", commentBy, opts.Body)
}

type ReviewComment struct {
	Body string
	Path string
	Line int
}

type Client interface {
	RefreshReviewComments(ctx context.Context, repoPath string, prId int, opts RefreshReviewCommentOptions) error
}

type scmClient struct {
	Client
}

func NewScmClient(ctx context.Context, t string, baseUrl string, token string) (Client, error) {
	if t == "github" {
		return &scmClient{
			Client: NewGithubClient(ctx, token),
		}, nil
	}

	if t == "gitlab" {
		client, err := NewGitlabClient(ctx, token, baseUrl)
		if err != nil {
			return nil, err
		}
		return &scmClient{
			Client: client,
		}, nil
	}

	return nil, fmt.Errorf("unknown type: %s", t)
}
