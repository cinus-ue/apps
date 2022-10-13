package cmd

import (
	"github.com/cinus-e/spy/internal/webcam"
	"github.com/urfave/cli/v2"
)

var Cam = &cli.Command{
	Name:  "cam",
	Usage: "Webcam recording",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "url",
			Aliases: []string{"u"},
			Usage:   "URL to a file upload handler",
		},
		&cli.IntFlag{
			Name:    "deviceID",
			Aliases: []string{"d"},
			Value:   0,
			Usage:   "Camera ID",
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
	Action: WebcamAction,
}

func WebcamAction(c *cli.Context) error {
	return webcam.StartCapture(webcam.Params{
		Address:  c.String("url"),
		DeviceID: c.Int("deviceID"),
		Quality:  c.Int("quality"),
		Verbose:  c.Bool("verbose"),
	})
}
