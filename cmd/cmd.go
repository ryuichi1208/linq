package cmd

import "io"

type Url struct {
	Url []string `yaml:"url"`
}

type withGoroutineID struct {
	out io.Writer
}

func CommandNameUsage() string {
	return "Usage: linq"
}
