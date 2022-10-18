package screen

import (
	"bytes"
	"image"
	"image/jpeg"
	"time"

	"github.com/kbinani/screenshot"
)

type Screen struct {
	bounds                image.Rectangle
	displayIndex, quality int
	buf                   *bytes.Buffer
	data                  chan []byte
	running               bool
}

func NewScreenCapturer(displayIndex, quality int, data chan []byte) *Screen {
	n := screenshot.NumActiveDisplays()
	if n < displayIndex {
		panic("Active display not found")
	}
	s := &Screen{}
	s.bounds = screenshot.GetDisplayBounds(displayIndex)
	s.displayIndex = displayIndex
	s.data = data
	s.buf = new(bytes.Buffer)
	s.quality = quality
	return s
}

func (s *Screen) Capture() error {
	s.running = true
	for s.running {
		img, err := screenshot.CaptureRect(s.bounds)
		if err != nil {
			return err
		}
		s.sendImg(img)
		s.data <- nil
	}
	return nil
}

func (s *Screen) sendImg(i image.Image) {
	s.buf.Reset()
	err := jpeg.Encode(s.buf, i, &jpeg.Options{Quality: s.quality})
	if err != nil {
		return
	}
	s.data <- s.buf.Bytes()
	time.Sleep(500 * time.Millisecond)
}

func (s *Screen) Close() {
	s.running = false
}
