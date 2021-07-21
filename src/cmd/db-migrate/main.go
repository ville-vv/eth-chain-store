package main

import (
	"github.com/urfave/cli"
	"github.com/ville-vv/eth-chain-store/src/common/log"
	"github.com/ville-vv/eth-chain-store/src/migrate"
	"os"
)

func main() {
	app := cli.NewApp()
	log.Init()
	app.HideHelp = true
	app.Name = "migrate data"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "action", Usage: "migrate data"},
	}
	app.Commands = []cli.Command{
		{Name: "create", Usage: "create database", Action: migrate.Create},
		{Name: "drop", Usage: "drop database", Action: migrate.Drop},
		{Name: "migrate", Usage: "init tables", Action: migrate.Migrate},
	}
	app.Action = func(c *cli.Context) {
		_ = cli.ShowAppHelp(c)
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
