package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	var language string
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "lang",
				Value:       "english",
				Destination: &language,
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			name := "slim"
			if command.NArg() > 0 {
				name = command.Args().Get(0)
			}
			//if command.String("lang") == "english" {
			if language == "cn" {
				fmt.Println("你好", name)
			} else {
				fmt.Println("hello", name)
			}
			return nil
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
