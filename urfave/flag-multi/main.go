package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

func main() {
	//  --greeting Hello --greeting Hola
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "greeting",
				Usage: "Pass multiple greetings",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println(strings.Join(cmd.StringSlice("greeting"), `, `))
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
