package main

import (
	"fmt"
	"github.com/rtmelsov/metrigger/internal/agent"
	// "github.com/rtmelsov/metrigger/cmd/staticlint"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\r\n", buildVersion)
	fmt.Printf("Build date: %s\r\n", buildDate)
	fmt.Printf("Build commit: %s\r\n", buildCommit)
	// staticlint.Check()

	a := agent.NewAgent()

	a.Run()
}
