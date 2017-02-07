package main

import (
	"os"

	"github.com/yangyuqian/genus/cmd/genus/command"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "Genus"
	app.Usage = "Simple tool helping code generation in Golang"
	app.Version = "1.0"

	app.Action = func(c *cli.Context) (err error) {
		cli.ShowAppHelp(c)
		return
	}

	app.Commands = []cli.Command{command.GenerateCmd}

	app.Run(os.Args)
}
