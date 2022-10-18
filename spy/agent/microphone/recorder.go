package microphone

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cinus-e/spy/literr"
	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/gen2brain/malgo"
)

type MicRecorder struct {
	name    string
	stream  *Streamer
	format  beep.Format
	context *malgo.AllocatedContext
}

func NewMicRecorder(name string) (*MicRecorder, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		log.Printf("LOG <%v>\n", message)
	})
	if err != nil {
		return nil, err
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS24
	deviceConfig.Capture.Channels = 2
	deviceConfig.SampleRate = 44100

	stream, format, err := OpenStream(ctx, deviceConfig)
	if err != nil {
		return nil, err
	}
	return &MicRecorder{name: name, stream: stream, format: format, context: ctx}, nil
}

func (m *MicRecorder) RecordWav() error {
	path, err := filepath.Abs(m.name)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	m.stream.Start()
	return wav.Encode(file, m.stream, m.format)
}

func (m *MicRecorder) Close() {
	literr.CheckError(m.stream.Close())
	_ = m.context.Uninit()
	m.context.Free()
}
