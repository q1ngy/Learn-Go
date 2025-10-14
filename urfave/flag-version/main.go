package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var (
	Revision = "fafafaf"
)

func main() {
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Printf("version=%s revision=%s\n", cmd.Root().Version, Revision)
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "print-version",
		Aliases: []string{"V"},
		Usage:   "print only the version",
	}

	cmd := &cli.Command{
		Name:    "partay",
		Version: "v19.99.0",
	}
	cmd.Run(context.Background(), os.Args)
}
