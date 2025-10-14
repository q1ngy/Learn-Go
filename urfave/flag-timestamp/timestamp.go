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
		Flags: []cli.Flag{
			&cli.TimestampFlag{
				Name: "meeting",
				Config: cli.TimestampConfig{
					Layouts: []string{"2006-01-02T15:04:05"},
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Printf("%s", cmd.Timestamp("meeting").String())
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
