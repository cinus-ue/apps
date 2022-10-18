package cmd

import (
	"os"
	"os/signal"

	"github.com/cinus-e/spy/agent/tracker"
	"github.com/urfave/cli/v2"
)

var App = &cli.Command{
	Name:  "app",
	Usage: "Application usage tracking",
	Subcommands: []*cli.Command{
		{
			Name:  "track",
			Usage: "Start tracking your usage activity",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "watch",
					Aliases: []string{"w"},
					Value:   5,
					Usage:   "Watch interval(seconds)",
				},
				&cli.IntFlag{
					Name:    "save",
					Aliases: []string{"s"},
					Value:   60,
					Usage:   "Save interval(minutes)",
				},
			},
			Action: trackAction,
		},
		{
			Name:      "show",
			Usage:     "Show application usage statistics",
			ArgsUsage: "<file1> <file2>",
			Action:    ShowAction,
		},
	},
}

func trackAction(c *cli.Context) error {
	tkr := tracker.NewTracker(c.Int("watch"), c.Int("save"))
	go tkr.TrackingActivity()
	go tkr.TrackingProcess()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	return nil
}

func ShowAction(c *cli.Context) error {
	tracker.ShowStatistics(c.Args().Slice())
	return nil
}
