package pkg

import (
	"bytes"
	"context"
	"fmt"
	pkgctx "github.com/chengjingtao/gomod-version-lint/pkg/context"
	"golang.org/x/mod/modfile"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type ModRequireAnalysis struct {
	modfile.Require

	Branches []string
	Error    error
}

func ExcludeBranches(ctx context.Context, require []ModRequireAnalysis, branchExcludeRegex string) ([]ModRequireAnalysis, error) {

	res := []ModRequireAnalysis{}

	for _, item := range require {
		matched, err := BranchMatched(ctx, item, branchExcludeRegex)
		if err != nil {
			return nil, err
		}

		if !matched {
			res = append(res, item)
		}
	}

	return res, nil
}

func BranchMatched(ctx context.Context, require ModRequireAnalysis, branchExcludeRegex string) (bool, error) {

	if branchExcludeRegex == "" {
		return true, nil
	}

	regex := branchExcludeRegex
	if !strings.HasSuffix(regex, "$") {
		regex = regex + "$"
	}
	if !strings.HasPrefix(regex, "^") {
		regex = "^" + regex
	}

	r, err := regexp.Compile(regex)
	if err != nil {
		return false, err
	}

	for _, branch := range require.Branches {
		matched := r.MatchString(branch)
		if matched {
			return true, nil
		}
	}

	return false, nil
}

func BranchAnalysis(ctx context.Context, modules []modfile.Require, concurrency int8) (require []ModRequireAnalysis) {
	logger := pkgctx.GetLogger(ctx)

	threshold := make(chan struct{}, concurrency)
	wg := sync.WaitGroup{}
	require = []ModRequireAnalysis{}
	requireLock := sync.RWMutex{}

	for _, _module := range modules {
		module := _module
		version := module.Mod.Version
		if len(strings.Split(version, "-")) == 3 {
			version = strings.Split(version, "-")[2]
		}

		wg.Add(1)
		go func() {
			threshold <- struct{}{}
			defer func() {
				<-threshold
				wg.Done()
			}()

			// TODO: support go proxy
			branches, _err := branchContains(ctx, "https://"+module.Mod.Path, version)
			if _err != nil {
				logger.Errorw("branch contains error", "module", module.Mod.Path, "version", module.Mod.Version, "err", _err)
			}

			requireLock.Lock()
			require = append(require, ModRequireAnalysis{
				Require:  module,
				Branches: branches,
				Error:    _err,
			})
			requireLock.Unlock()
		}()
	}

	wg.Wait()
	return require
}

func branchContains(ctx context.Context, repoUrl string, commitID string) ([]string, error) {
	logger := pkgctx.GetLogger(ctx)

	dir := encodeRepoUrl(repoUrl)
	logger.Debugf("mkdir /tmp/%s", dir)
	tmp, err := os.MkdirTemp("/tmp", dir)
	if err != nil {
		logger.Errorf("mk temp dir error: %s", err.Error())
		return nil, err
	}

	args := []string{
		"clone",
		"--filter=blob:none",
		"--no-checkout",
		repoUrl,
		"./",
	}

	stdout, _, err := runCmd(ctx, tmp, "git", args...)
	if err != nil {
		return nil, err
	}

	stdout, _, err = runCmd(ctx, tmp, "git", []string{
		"branch",
		"-q",
		"-r",
		"--contains",
		commitID,
	}...)

	if err != nil {
		return nil, err
	}

	return parseStdoutOfBranchContains(stdout), nil
}

func parseStdoutOfBranchContains(stdout string) []string {
	if len(stdout) == 0 {
		return nil
	}
	items := strings.Split(stdout, "\n")
	if len(items) == 0 {
		return nil
	}

	branches := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if strings.HasPrefix(item, "origin/HEAD") {
			continue
		}
		branch := strings.TrimPrefix(item, "origin/")
		if branch == "" {
			continue
		}

		branches = append(branches, branch)
	}

	return branches
}

func runCmd(ctx context.Context, workdir, name string, args ...string) (stdout string, stderr string, err error) {
	logger := pkgctx.GetLogger(ctx)

	cmdStr := name + " " + strings.Join(args, " ")
	logger.Infof("executing \"%s\" in \"%s\" \n", cmdStr, workdir)

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = []string{
		"GIT_TERMINAL_PROMPT=false",
		"https_proxy=" + os.Getenv("https_proxy"),
		"http_proxy=" + os.Getenv("http_proxy"),
		"all_proxy=" + os.Getenv("all_proxy"),
	}

	cmd.Dir = workdir
	stdoutBf := bytes.NewBufferString("")
	stderrBf := bytes.NewBufferString("")
	cmd.Stdout = io.MultiWriter(os.Stdout, stdoutBf)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderrBf)
	err = cmd.Run()

	if ctx.Err() != nil {
		if err == nil {
			err = ctx.Err()
		} else {
			err = fmt.Errorf("err: %s, contextErr: %s", err, ctx.Err().Error())
		}
	}

	stdoutBts, _err := io.ReadAll(stdoutBf)
	if _err != nil {
		logger.Errorf("read com commd %s stdout error %s", cmdStr, _err)
	}
	stderrBts, _err := io.ReadAll(stderrBf)
	if _err != nil {
		logger.Errorf("read com commd %s stdout error %s", cmdStr, _err)
	}

	return string(stdoutBts), string(stderrBts), err
}

func encodeRepoUrl(repoUrl string) string {
	tmpDir := strings.TrimPrefix(repoUrl, "https://")
	tmpDir = strings.TrimSuffix(tmpDir, "http://")
	tmpDir = strings.Replace(tmpDir, ":", "-", -1)
	tmpDir = strings.Replace(tmpDir, "/", "-", -1)
	return tmpDir
}
