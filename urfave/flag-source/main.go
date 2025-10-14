package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	// export APP_LANG=en_US
	// export LEGACY_COMPAT_LANG=en_CN
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "lang",
				Aliases: []string{"l"},
				Value:   "english",
				Usage:   "language for the greeting",
				Sources: cli.EnvVars("LEGACY_COMPAT_LANG", "APP_LANG", "LANG"),
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			fmt.Println(command.String("lang"))
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
