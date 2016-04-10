package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/susp/vimperator-flavor/parser"
)

var cmdInstall = &Command{
	Run:       runInstall,
	UsageLine: "install ",
	Short:     "Install",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdInstall.Flag.BoolVar(&flagA, "a", false, "")
}

// runInstall executes install command and return exit code.
func runInstall(args []string) int {
	repos, err := parser.ParseFile("VimperatorFlavor")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// TODO: fix in case GOPATH is set to multiple paths
	srcRoot := filepath.Join(os.Getenv("GOPATH"), "src")

	usr, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	flavorsRoot := usr.HomeDir + "/.vimperator/flavors"

	for _, repo := range repos {
		if _, err := os.Stat(srcRoot + "/" + repo.Path); err != nil {
			cmd := exec.Command("git", "clone", "https://"+repo.Path, srcRoot+"/"+repo.Path)
			if err := cmd.Run(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}
		}

		flavorsDir := flavorsRoot + "/" + strings.Replace(repo.Path, "/", "_", -1)
		if _, err := os.Stat(flavorsDir); err == nil {
			continue
		}

		rev := "master"
		if repo.Version != "" {
			rev = repo.Version
		}

		if err := os.MkdirAll(flavorsDir, 0777); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		cmd := exec.Command("sh", "-c", "git archive "+rev+" | (cd "+flavorsDir+" && tar xf -)")
		cmd.Dir = srcRoot + "/" + repo.Path
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}

	return 0
}
