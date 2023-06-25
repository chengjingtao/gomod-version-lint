package options

import (
	"context"
	"encoding/json"
	"fmt"
	flag "github.com/spf13/pflag"
	"gomod.alauda.cn/gomod-version-lint/pkg"
	pkgctx "gomod.alauda.cn/gomod-version-lint/pkg/context"
	"gopkg.in/yaml.v3"
	"io"
	iofs "io/fs"
	"os"
	"path"
	"strings"
)

// BranchesOptions branches command options
type BranchesOptions struct {
	RootOptions

	// ExcludeBranchesRegex branches regex
	ExcludeBranchesRegex string
	// ModuleRegex name regex
	ModuleRegex string
	// gomod directory
	ModDir string
	// OutputFmt output format, json or yaml
	OutputFmt string
	// OutputFile output file name
	OutputFile string
	// CommentsFile comments file name
	CommentsFile string
	Concurrency  int8

	FS      iofs.FS
	Context context.Context
}

func (opts *BranchesOptions) Run() error {
	logger := pkgctx.GetLogger(opts.Context)

	modDir := "./"
	if opts.ModDir != "" {
		modDir = opts.ModDir
	}
	modFilePath := path.Join(modDir, "go.mod")

	fs := os.DirFS(modDir)
	if opts.FS != nil {
		opts.FS = fs
	}

	bts, err := iofs.ReadFile(fs, "go.mod")
	if err != nil {
		logger.Errorf("read file %s error: %s", modFilePath, err.Error())
		return err
	}

	modFile, err := pkg.ParseModFile(modFilePath, bts)
	if err != nil {
		logger.Errorf("parse mod file error: %s", err.Error())
		return err
	}

	requredModules, err := pkg.MatchModules(opts.Context, modFile, opts.ModuleRegex)
	if err != nil {
		logger.Errorf("match modules by regex: %v error: %s", opts.ModuleRegex, err)
		return err
	}

	modRequireAnalysis := pkg.BranchAnalysis(opts.Context, requredModules, opts.Concurrency)
	err = opts.writeAnalysisResultV2(modRequireAnalysis)
	if err != nil {
		return err
	}

	modRequireAnalysis, err = pkg.ExcludeBranches(opts.Context, modRequireAnalysis, opts.ExcludeBranchesRegex)
	if err != nil {
		return err
	}

	if opts.CommentsFile != "" {
		err = opts.writeGitCommentsFile(modRequireAnalysis, modFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (opts *BranchesOptions) writeAnalysisResultV2(modRequireAnalysis []pkg.ModRequireAnalysis) error {
	fmt.Printf("### ANALYSIS RESULT\n")
	if len(modRequireAnalysis) == 0 {
		fmt.Printf("\n all modules %s branches matched %s", opts.ModuleRegex, opts.ExcludeBranchesRegex)
		return nil
	}

	for _, item := range modRequireAnalysis {
		matched, err := pkg.BranchMatched(opts.Context, item, opts.ExcludeBranchesRegex)
		if err != nil {
			fmt.Printf("🐛  %s \t error: %s", item.Mod.Path, err.Error())
			continue
		}
		flag := "✅️"
		if !matched {
			flag = "⚠️ "
		}
		fmt.Printf("%s  %s %s\n", flag, fillSpace(item.Mod.Path+"@"+item.Mod.Version, 100), strings.Join(item.Branches, ","))
	}

	return nil
}

func fillSpace(str string, width int) string {
	left := width - len(str)
	if left > 0 {
		return str + strings.Repeat(" ", left)
	}
	return str
}

func (opts *BranchesOptions) writeAnalysisResult(modRequireAnalysis []pkg.ModRequireAnalysis) error {
	if len(modRequireAnalysis) == 0 {
		os.Stdout.Write([]byte(fmt.Sprintf("👍🏻 all modules %s branches matched %s", opts.ModuleRegex, opts.ExcludeBranchesRegex)))
		return nil
	}

	writers := []io.Writer{os.Stdout}

	if opts.OutputFile != "" {
		f, err := os.Create(opts.OutputFile)
		if err != nil {
			return err
		}
		writers = append(writers, f)
	}

	err := opts.Output(modRequireAnalysis, io.MultiWriter(writers...))
	if err != nil {
		return err
	}
	return nil
}

func (opts *BranchesOptions) writeGitCommentsFile(modRequireAnalysis []pkg.ModRequireAnalysis, modFilePath string) error {
	commentsFile, err := os.Create(opts.CommentsFile)
	if err != nil {
		return err
	}

	comments := makeGitFileComments(modRequireAnalysis, modFilePath)
	fmt.Printf("### GIT COMMENTS\n")
	err = comments.Marshal(io.MultiWriter(os.Stdout, commentsFile))
	if err != nil {
		return err
	}
	return nil
}

func makeGitFileComments(mods []pkg.ModRequireAnalysis, modFilePath string) GitFileComments {
	comments := GitFileComments{}

	for _, item := range mods {

		body := fmt.Sprintf("⚠️ branch is %s for version: %s", strings.Join(item.Branches, ","), item.Mod.Version)
		if strings.Join(item.Branches, ",") == "" {
			body = "not found any branch for version: " + item.Mod.Version
		}
		comments = append(comments, GitFileComment{
			FilePath: modFilePath,
			Line:     item.Syntax.Start.Line,
			Comment:  body,
		})
	}
	return comments
}

func (opts *BranchesOptions) Output(requires []pkg.ModRequireAnalysis, writer io.Writer) error {
	outputFmt := "json"
	if opts.OutputFmt != "" {
		outputFmt = opts.OutputFmt
	}

	if outputFmt == "json" {
		bts, err := json.MarshalIndent(requires, "", "  ")
		if err != nil {
			return err
		}

		writer.Write(bts)
		return nil
	}

	if outputFmt == "yaml" {
		bts, err := yaml.Marshal(requires)
		if err != nil {
			return err
		}

		writer.Write(bts)
		return nil
	}

	if outputFmt == "table" {
		for _, item := range requires {
			writer.Write([]byte(item.Mod.Path + "|"))
			writer.Write([]byte(item.Mod.Version + "|" + strings.Join(item.Branches, ",")))
			writer.Write([]byte("|" + fmt.Sprint(item.Syntax.End.Line)))
			writer.Write([]byte("\n"))
		}
		return nil
	}

	return fmt.Errorf("unknown output format: %s", opts.OutputFmt)
}

func (opts *BranchesOptions) AddFlags(flags *flag.FlagSet) {
	flags.StringVar(&opts.ModuleRegex, "module", "github.com/example/.*", "modules that you want to print branches, it supports using regex")
	flags.StringVar(&opts.ExcludeBranchesRegex, "branches-exclude", "(^main$|^release-.*$)", "branch of modules that you want to exclude, it supports usiing regex")
	flags.StringVarP(&opts.ModDir, "mod-dir", "d", "./", "gomod file directory")
	flags.StringVarP(&opts.OutputFmt, "out", "o", "table", "gomod file path")
	flags.StringVar(&opts.OutputFile, "out-file", "table", "gomod file path")
	flags.StringVar(&opts.CommentsFile, "comments-file", ".git-comments", "comments file")
	flags.Int8Var(&opts.Concurrency, "concurrency", 5, "concurrency count for analysis modules")
}
