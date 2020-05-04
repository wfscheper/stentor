package main

import "fmt"

const (
	appName = "stentor"
)

var (
	version = "dev"
	commit  = "none"
	date    = "none"
)

func main() {
	fmt.Printf("Hello world, from %s version %s, commit %s, built at %s!\n", appName, version, commit, date)
}
