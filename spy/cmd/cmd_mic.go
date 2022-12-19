package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/cinus-ue/spy/agent/microphone"
	"github.com/cinus-ue/spy/literr"
	"github.com/cinus-ue/spy/util"
	"github.com/urfave/cli/v2"
)

var Mic = &cli.Command{
	Name:   "mic",
	Usage:  "Microphone recording",
	Action: MicrophoneAction,
}

func MicrophoneAction(*cli.Context) error {
	recorder, err := microphone.NewMicRecorder(util.FileNameFormat("mic", ".wav"))
	if err != nil {
		return err
	}
	fmt.Printf("Recording.\nPress Ctrl-C to stop.\n")
	go literr.CheckFatal(recorder.RecordWav())
	// Stop the stream when the user tries to quit the program.
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	recorder.Close()
	return nil
}
