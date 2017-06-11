package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/elpinal/vimperator-flavor/parser"
)

var cmdUpgrade = &Command{
	Run:       runUpgrade,
	UsageLine: "upgrade ",
	Short:     "Upgrade",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdUpgrade.Flag.BoolVar(&flagA, "a", false, "")
}

// runUpgrade executes upgrade command and return exit code.
func runUpgrade(args []string) int {
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
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		joinedPath := strings.Replace(repo.Path, "/", "_", -1)
		flavorsDir := flavorsRoot + "/" + joinedPath
		if _, err := os.Stat(flavorsDir); err != nil {
			continue
		}

		rev := "master"
		if repo.Version != "" {
			rev = repo.Version
		}

		cmd := exec.Command("git", "pull")
		cmd.Dir = srcRoot + "/" + repo.Path
		/*
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		*/

		if err := os.RemoveAll(flavorsDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if err := os.Mkdir(flavorsDir, 0777); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		cmd = exec.Command("sh", "-c", "git archive --prefix="+joinedPath+"/ "+rev+" | (cd "+flavorsRoot+" && tar xf -)")
		cmd.Dir = srcRoot + "/" + repo.Path
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}

	return 0
}
