package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/elpinal/vimperator-flavor/parser"
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

	if _, err := os.Stat(flavorsRoot + "/bootstrap.js"); err != nil {
		data := []byte(`
		(function() {
			if (liberator.globalVariables.loaded_my_flavors == "true") {
				return;
			}

			let rtps = options.runtimepath.split(",");
			let flavorDirs = [];
			for (let dir in io.File("~/.vimperator/flavors").iterDirectory()) {
				if (!dir.isDirectory()) {
					continue;
				}
				flavorDirs.push(dir.path);
			}
			let newRtps = rtps.concat(flavorDirs);
			options.runtimepath = newRtps.join(",");
		})();`)
		if err := ioutil.WriteFile(flavorsRoot+"/bootstrap.js", data, 0644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}

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
