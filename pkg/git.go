package pkg

import (
	"bytes"
	"context"
	"fmt"
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

	if branchExcludeRegex == "" {
		return require, nil
	}

	res := []ModRequireAnalysis{}

	regex := branchExcludeRegex
	if !strings.HasSuffix(regex, "$") {
		regex = regex + "$"
	}
	if !strings.HasPrefix(regex, "^") {
		regex = "^" + regex
	}

	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

OUTLOOP:
	for _, item := range require {
		for _, branch := range item.Branches {
			matched := r.MatchString(branch)
			if matched {
				continue OUTLOOP
			}
		}

		// not match branchExcludeRegex
		res = append(res, item)
	}

	return res, nil
}

func BranchAnalysis(ctx context.Context, modules []modfile.Require) (require []ModRequireAnalysis) {
	logger := GetLogger(ctx)

	threshold := make(chan struct{}, 20)
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

func BranchContains(ctx context.Context, repoUrl string, commitID string) ([]string, error) {
	return branchContains(ctx, repoUrl, commitID)
}

func branchContains(ctx context.Context, repoUrl string, commitID string) ([]string, error) {
	logger := GetLogger(ctx)

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

	stdout, stderr, err := runCmd(ctx, tmp, "git", args...)
	if err != nil {
		return nil, err
	}

	if stderr != "" {
		logger.Infof(stderr)
	}

	stdout, stderr, err = runCmd(ctx, tmp, "git", []string{
		"branch",
		"-q",
		"-r",
		"--contains",
		commitID,
	}...)

	if err != nil {
		return nil, err
	}

	if stderr != "" {
		logger.Infof(stderr)
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

		branches = append(branches, branch)
	}

	return branches
}

func runCmd(ctx context.Context, workdir, name string, args ...string) (stdout string, stderr string, err error) {
	logger := GetLogger(ctx)

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
	cmd.Stdout = stdoutBf
	cmd.Stderr = stderrBf
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
