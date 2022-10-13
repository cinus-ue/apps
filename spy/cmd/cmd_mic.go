package cmd

import (
	"fmt"
	"github.com/cinus-e/spy/internal/microphone"
	"github.com/cinus-e/spy/internal/util"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
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
	go func() {
		if err := recorder.RecordWav(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()
	// Stop the stream when the user tries to quit the program.
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	recorder.Close()
	return nil
}
