package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "greet",
		Usage: "just greet, nothing else to say.",
		Action: func(ctx context.Context, command *cli.Command) error {
			fmt.Println("greeting!")
			return nil
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
