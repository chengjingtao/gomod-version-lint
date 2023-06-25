package options

import (
	"bufio"
	"context"
	"fmt"
	flag "github.com/spf13/pflag"
	pkgctx "gomod.alauda.cn/gomod-version-lint/pkg/context"
	pkgscm "gomod.alauda.cn/gomod-version-lint/pkg/scm"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type CommentOptions struct {
	File       string
	ServerType string
	Repository string
	PrID       int
	CommitID   string

	Context context.Context
}

func (opts *CommentOptions) Run() error {
	serverType := "github"
	baseUrl := ""
	if strings.Contains(strings.ToLower(opts.Repository), "gitlab") {
		serverType = "gitlab"
	}
	if opts.ServerType != "" {
		serverType = opts.ServerType
	}
	u, err := url.Parse(opts.Repository)
	if err != nil {
		return err
	}
	baseUrl = opts.Repository[0 : strings.LastIndex(opts.Repository, u.Host)+len(u.Host)]

	token := os.Getenv("TOKEN")
	if token == "" {
		return fmt.Errorf("should provide private access token by env: TOKEN")
	}

	client, err := pkgscm.NewScmClient(opts.Context, serverType, baseUrl, token)
	if err != nil {
		return err
	}

	comments, err := loadCommentFromFile(opts.Context, opts.File)
	if err != nil {
		return err
	}

	repo := strings.TrimSuffix(u.Path, ".git")
	repo = strings.TrimPrefix(repo, "/")
	repo = strings.TrimSuffix(repo, "/")

	err = client.RefreshReviewComments(opts.Context, repo, opts.PrID, pkgscm.RefreshReviewCommentOptions{
		CommentBy: "gomod-version-lint",
		CommitID:  opts.CommitID,
		Comments:  comments.ToReviewComments(),
	})

	return err
}

func loadCommentFromFile(ctx context.Context, file string) (*GitFileComments, error) {

	fs, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	comments := &GitFileComments{}
	err = comments.UnMarshal(ctx, fs)
	if err != nil {
		return nil, nil
	}
	return comments, err
}

func (opts *CommentOptions) AddFlags(flags *flag.FlagSet) {
	flags.StringVar(&opts.Repository, "repo", "https://gitub.com/example/example", "repository address")
	flags.StringVar(&opts.File, "file", "./comments", "comments file, the file format of each line should be file-path|line-number|comment-body")
	flags.StringVar(&opts.ServerType, "server-type", "", "git server type, eg. github gitlab, it is optional, "+
		"if you do not provide it, it will set to gitlab when repository contains 'gitlab'")
	flags.StringVar(&opts.CommitID, "commit-id", "", "current commit id of branch")
	flags.IntVar(&opts.PrID, "pr-id", -1, "id of pull request")
}

// file-path|line|comment
type GitFileComment struct {
	FilePath string
	Line     int
	Comment  string
}
type GitFileComments []GitFileComment

func (comments *GitFileComments) Marshal(writer io.Writer) error {
	for _, item := range *comments {
		line := fmt.Sprintf("%s|%d|%s\n", item.FilePath, item.Line, item.Comment)
		_, err := writer.Write([]byte(line))
		if err != nil {
			return err
		}
	}

	return nil
}

func (comments *GitFileComments) UnMarshal(ctx context.Context, reader io.Reader) error {
	logger := pkgctx.GetLogger(ctx)

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	var index = 0
	for scanner.Scan() {
		index++
		line := scanner.Text()
		if line == "" {
			continue
		}
		segments := strings.Split(line, "|")
		if len(segments) < 3 {
			logger.Warnf("error format of line: %d", index)
			continue
		}
		lineNum, err := strconv.Atoi(segments[1])
		if err != nil {
			logger.Errorf("error parse %s to int, error:%s", segments[1], err.Error())
			return err
		}
		*comments = append(*comments, GitFileComment{
			FilePath: segments[0],
			Line:     lineNum,
			Comment:  segments[2],
		})
	}

	return nil
}
func (comments GitFileComments) ToReviewComments() []pkgscm.ReviewComment {
	res := []pkgscm.ReviewComment{}

	for _, item := range comments {
		com := pkgscm.ReviewComment{
			Body: item.Comment,
			Path: item.FilePath,
			Line: item.Line,
		}
		res = append(res, com)
	}

	return res
}
