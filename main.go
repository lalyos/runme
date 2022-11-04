package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/itchyny/gojq/cli"
	fzf "github.com/junegunn/fzf/src"
	"github.com/progrium/go-basher"
	"github.com/stateful/runme/internal/cmd"
)

func fakeArgs(cmd string, args []string) func() {
	oldargs := os.Args
	os.Args = []string{cmd}
	os.Args = append(os.Args, args...)

	return func() {
		os.Args = oldargs
	}
}
func myJq(args []string) {
	restoreArgs := fakeArgs("jq", args)
	defer restoreArgs()
	cli.Run()
}

func myFzf(args []string) {
	restoreArgs := fakeArgs("fzf", args)
	defer restoreArgs()
	fzf.Run(fzf.ParseOptions(), "", "")
}

func myRunme(args []string) {
	restoreArgs := fakeArgs("runme", args)
	defer restoreArgs()

	root := cmd.Root()
	root.Version = fmt.Sprintf("stateful %s (%s) on %s", BuildVersion, Commit, BuildDate)
	root.Execute()
}

//go:embed bash/*.bash
var f embed.FS

func Asset(name string) ([]byte, error) {
	return f.ReadFile(name)
}

// These are variables so that they can be set during the build time.
var (
	BuildDate    = "unknown"
	BuildVersion = "0.0.0"
	Commit       = "unknown"
)

func main() {
	loader := Asset
	if os.Getenv("DEBUG") != "" {
		loader = nil
	}
	basher.Application(map[string]func([]string){
		"_jq":    myJq,
		"_fzf":   myFzf,
		"_runme": myRunme,
	}, []string{
		"bash/runme.bash",
		"bash/fn.bash",
		"bash/cmd.bash",
	}, loader, true)

}
