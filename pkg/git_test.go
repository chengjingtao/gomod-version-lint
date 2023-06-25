package pkg

import (
	"context"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"testing"
)

func TestExcludeBranches(t *testing.T) {
	ctx := context.Background()

	modules := []ModRequireAnalysis{
		{
			Require: modfile.Require{
				Mod: module.Version{
					Path:    "git.example.com/demo/demo",
					Version: "v1.0.0-20201130134442-10cb98267c6c",
				},
			},
			Branches: []string{
				"feat/test1",
				"main",
			},
		},
		{
			Require: modfile.Require{
				Mod: module.Version{
					Path:    "git.example.com/demo/demo",
					Version: "v1.0.0-20201130134442-10cb98267c6c",
				},
			},
			Branches: []string{
				"main",
			},
		},
		{
			Require: modfile.Require{
				Mod: module.Version{
					Path:    "git.example.com/demo/demo",
					Version: "v1.0.0-20201130134442-10cb98267c6c",
				},
			},
			Branches: []string{
				"feat/test1",
				"feat/test2",
			},
		},
	}

	res, err := ExcludeBranches(ctx, modules, "(^main$|^release-.*$)")
	if err != nil {
		t.Errorf("err should be nil, error: %s", err.Error())
		return
	}

	if len(res) != 1 {
		t.Error("modrequire after exclude branches should be zero")
		return
	}

	if res[0].Mod.Path != "git.example.com/demo/demo" {
		t.Errorf("modrequire after exclude branch should return mod: git.example.com/demo/demo, but: %s", res[0].Mod.Path)
	}

	if (res[0].Branches[0] != "feat/test1") || (res[0].Branches[1] != "feat/test2") {
		t.Errorf(`modrequire after exclude branch should return branches:[]{"feat/test1", "feat/test2"}, but: %#v`, res[0].Branches)
	}
}
