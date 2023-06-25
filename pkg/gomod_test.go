package pkg

import (
	"context"
	"testing"
)

var modfileString = `
module github.com/chengjingtao/gomod-version-lint

go 1.19

require (
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/example/abc v1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/example/abc/def v1
)
`

func TestMatchModules(t *testing.T) {
	ctx := context.Background()
	modfile, err := ParseModFile("./go.mod", []byte(modfileString))
	if err != nil {
		t.Errorf("should parse mod file correctly, but error: %s", err.Error())
		return
	}

	requires, err := MatchModules(ctx, modfile, "github.com/example/.*")
	if err != nil {
		t.Errorf("match modules should not return error, but error: %s", err.Error())
		return
	}

	if len(requires) != 2 {
		t.Errorf("mod requires length should be 2, but got: %d", len(requires))
		return
	}

	if requires[0].Mod.Path != "github.com/example/abc" || requires[1].Mod.Path != "github.com/example/abc/def" {
		t.Errorf(`mod requires should be "github.com/example/abc" or "github.com/example/abc/def" , but: %#v`, requires)
		return
	}
}
