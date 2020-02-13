package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Commands = []*cli.Command{
		{
			Name:  "put",
			Usage: "put a file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "filepath",
					Aliases:  []string{"f"},
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				log.Println("putting file", c.String("filepath"))
				return nil
			},
		},
		{
			Name:  "get",
			Usage: "get a file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "key",
					Aliases:  []string{"k"},
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				log.Println("getting file", c.String("key"))
				return nil
			},
		},
	}
	app.Run(os.Args)
}
