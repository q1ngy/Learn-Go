package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli/v3"
)

func main() {
	var mode string
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "mode",
				Value:       "+",
				Destination: &mode,
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			if command.NArg() != 2 {
				log.Fatal("需要两个参数")
			}
			arg1 := command.Args().Get(0)
			num1, err := strconv.Atoi(arg1)
			arg2 := command.Args().Get(1)
			num2, err := strconv.Atoi(arg2)

			switch mode {
			case "+":
				fmt.Println(num1 + num2)
			case "-":
				fmt.Println(num1 - num2)
			case "*":
				fmt.Println(num1 * num2)
			case "/":
				fmt.Println(num1 / num2)
			case "%":
				fmt.Println(num1 % num2)
			}
			return err
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
