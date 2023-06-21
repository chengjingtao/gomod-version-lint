package pkg

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/mod/modfile"
	"regexp"
	"strings"
)

func ParseModFile(modFilePath string, bts []byte) (*modfile.File, error) {
	file, err := modfile.Parse(modFilePath, bts, nil)

	return file, err
}

func MatchModules(ctx context.Context, file *modfile.File, modulesRegex string) (requires []modfile.Require, err error) {
	if file == nil {
		return nil, errors.New("modfile should not be nil")
	}

	if modulesRegex != "" {
		if !strings.HasPrefix(modulesRegex, "^") {
			modulesRegex = "^" + modulesRegex
		}
		if !strings.HasSuffix(modulesRegex, "$") {
			modulesRegex = modulesRegex + "$"
		}
	}

	reg, err := regexp.Compile(modulesRegex)
	if err != nil {
		return nil, fmt.Errorf("regex '%s' error: %s ", modulesRegex, err.Error())
	}

	matchedRequires := []modfile.Require{}

	for _, item := range file.Require {
		if modulesRegex != "" {
			matched := reg.MatchString(item.Mod.Path)
			if !matched {
				continue
			}
		}

		matchedRequires = append(matchedRequires, *item)
	}

	return matchedRequires, nil
}
