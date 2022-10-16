package cmd

import (
	"github.com/cinus-e/spy/internal/ipc"
	"github.com/urfave/cli/v2"
)

var Ste = &cli.Command{
	Name:   "stealth",
	Usage:  "Run in the background using IPC communication",
	Action: StealthAction,
}

func StealthAction(*cli.Context) error {
	return ipc.StartServer()
}
