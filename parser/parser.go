package parser

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Repo struct {
	path    string
	version string
}

func ParseFile(name string) ([]Repo, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return Parse(f)
}

func Parse(r io.Reader) ([]Repo, error) {
	s := bufio.NewScanner(r)
	repos := make([]Repo, 0, 5)
	for s.Scan() {
		txt := s.Text()
		txt = strings.TrimSpace(txt)
		if txt == "" {
			continue
		}
		if i := strings.Index(txt, " "); i > 0 {
			ver := strings.TrimSpace(txt[i+1:])
			repos = append(repos, Repo{path: txt[:i], version: ver})
		} else {
			repos = append(repos, Repo{path: txt})
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return repos, nil
}
