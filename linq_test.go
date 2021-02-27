package main

import (
	"testing"
)

func TestExampleSuccess(t *testing.T) {
	l := readYaml("./test.yml")
	if len(l.Url) != 39 {
		t.Fatal("array reading error")
	}
}
