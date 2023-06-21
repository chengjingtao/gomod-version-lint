package options

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chengjingtao/gomod-version-lint/pkg"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
	"io"
	iofs "io/fs"
	"os"
	"path"
	"strings"
)

// BranchesOptions branches command options
type BranchesOptions struct {
	// ExcludeBranchesRegex branches regex
	ExcludeBranchesRegex string
	// ModuleRegex name regex
	ModuleRegex string

	// gomod directory
	ModDir string

	// OutputFmt output format, json or yaml
	OutputFmt string

	FS      iofs.FS
	Context context.Context
}

func (opts *BranchesOptions) Run() error {
	logger := pkg.GetLogger(opts.Context)

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

	modRequireAnalysis := pkg.BranchAnalysis(opts.Context, requredModules)
	modRequireAnalysis, err = pkg.ExcludeBranches(opts.Context, modRequireAnalysis, opts.ExcludeBranchesRegex)
	if err != nil {
		return err
	}

	err = opts.Output(modRequireAnalysis, os.Stdout)
	if err != nil {
		return err
	}

	return nil
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

	if outputFmt == "simple" {
		for _, item := range requires {
			writer.Write([]byte(item.Mod.Path + "|"))
			writer.Write([]byte(item.Mod.Version + "|" + strings.Join(item.Branches, ",")))
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
	flags.StringVarP(&opts.OutputFmt, "out", "o", "./go.mod", "gomod file path")
}
