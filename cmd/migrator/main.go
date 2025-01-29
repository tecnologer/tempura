package main

import (
	"os"

	"github.com/tecnologer/tempura/cmd/migrator/cli"
	"github.com/tecnologer/tempura/pkg/utils/log"
)

var version string

func main() {
	newCLI := cli.NewCLI(version)
	if err := newCLI.Run(os.Args); err != nil {
		log.Error(err.Error())
	}
}
