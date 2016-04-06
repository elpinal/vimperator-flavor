package parser

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	r := strings.NewReader(`
	user/repo
	user1/repo-x
	user3/repo-y 1.0.1
	
	user5/repo-z  1.0.13  `)
	repos, err := Parse(r)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%q\n", repos)
}
