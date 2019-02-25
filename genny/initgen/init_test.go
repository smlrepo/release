package initgen

import (
	"context"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/genny/movinglater/gotools/gomods"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	gg, err := New(&Options{
		VersionFile: "foo/bar/version.go",
		MainFile:    "./main.go",
		Root:        ".",
	})
	r.NoError(err)

	run := genny.DryRunner(context.Background())
	run.WithGroup(gg)

	r.NoError(run.Run())

	var cmds []string
	res := run.Results()
	if !gomods.On() {
		cmds = []string{"git init", "go get github.com/alecthomas/gometalinter", "gometalinter --install"}
	} else {
		cmds = []string{"git init", "go mod init", "go get github.com/alecthomas/gometalinter", "gometalinter --install", "go mod tidy"}
	}

	r.NoError(gentest.CompareCommands(cmds, res.Commands))

	files := []string{
		".gitignore",
		".gometalinter.json",
		".goreleaser.yml.plush",
		"azure-pipelines.yml",
		"azure-tests.yml",
		"LICENSE",
		"Makefile",
		"foo/bar/version.go",
	}
	r.NoError(gentest.CompareFiles(files, res.Files))

	f, err := res.Find("foo/bar/version.go")
	r.NoError(err)
	r.Contains(f.String(), `const Version = "v0.0.1"`)
}
