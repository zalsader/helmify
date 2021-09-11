package main

import (
	"fmt"
)

// these information will be collected when build, by `-ldflags "-X main.appVersion=0.1"`.
var (
	appVersion = "development"
	buildTime  = "not set"
	gitCommit  = "not set"
	gitRef     = "not set"
)

func printVersion() {
	fmt.Printf("Version:    %s\n", appVersion)
	fmt.Printf("Build Time: %s\n", buildTime)
	fmt.Printf("Git Commit: %s\n", gitCommit)
	fmt.Printf("Git Ref:    %s\n", gitRef)
}
