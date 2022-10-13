package webcam

import (
	"bytes"
	"fmt"
	"gocv.io/x/gocv"
	"image/jpeg"
	"time"
)

type Device struct {
	capture           *gocv.VideoCapture
	deviceID, quality int
	buf               *bytes.Buffer
	data              chan []byte
}

func NewVideoCapturer(deviceID, quality int, data chan []byte) (*Device, error) {
	capture, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		return nil, err
	}
	return &Device{capture, deviceID, quality, new(bytes.Buffer), data}, nil
}

func (d *Device) Capture() error {
	mat := gocv.NewMat()
	defer mat.Close()
	for {
		if ok := d.capture.Read(&mat); !ok {
			return fmt.Errorf("device closed: %d", d.deviceID)
		}
		if mat.Empty() {
			continue
		}
		img, err := mat.ToImage()
		if err != nil {
			return err
		}
		d.buf.Reset()
		if err = jpeg.Encode(d.buf, img, &jpeg.Options{Quality: d.quality}); err != nil {
			return err
		}
		d.data <- d.buf.Bytes()
		time.Sleep(500 * time.Millisecond)
	}
}
