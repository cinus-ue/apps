package cmd

import (
	"github.com/cinus-e/spy/agent/screen"
	"github.com/urfave/cli/v2"
)

var Scr = &cli.Command{
	Name:  "scr",
	Usage: "Screen recording",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "url",
			Aliases: []string{"u"},
			Usage:   "URL to a file upload handler",
		},
		&cli.IntFlag{
			Name:    "display",
			Aliases: []string{"d"},
			Value:   0,
			Usage:   "Screen display index",
		},
		&cli.IntFlag{
			Name:    "quality",
			Aliases: []string{"q"},
			Value:   75,
			Usage:   "JPEG compress quality",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Value:   false,
			Usage:   "Enable verbose output",
		},
	},
	Action: ScreenAction,
}

func ScreenAction(c *cli.Context) error {
	cap := screen.Capture{}
	return cap.StartCapture(screen.Params{
		Address: c.String("url"),
		Display: c.Int("display"),
		Quality: c.Int("quality"),
		Verbose: c.Bool("verbose"),
	})
}
