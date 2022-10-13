package screen

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type Params struct {
	Address          string
	Display, Quality int
	Verbose          bool
}

var (
	verbose = false
	running = false
)

func StartCapture(p Params) (err error) {
	verbose = p.Verbose
	running = true
	data := make(chan []byte, 200000)
	screen := NewScreenCapturer(p.Display, p.Quality, data)
	go func() {
		if err = screen.Capture(); err != nil {
			running = false
		}
	}()
	for running {
		frame := <-data
		if frame != nil {
			if err := write(frame, p.Address); err != nil {
				fmt.Println(err)
			}
		}
		time.Sleep(1 * time.Millisecond)
	}
	return
}

func write(data []byte, address string) error {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	fileWriter, _ := writer.CreateFormFile("files", "screen")
	written, err := io.Copy(fileWriter, bytes.NewReader(data))
	if err != nil {
		return err
	}
	if verbose {
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
	if verbose {
		fmt.Printf(" status:%s response:%s\n", resp.Status, string(respBody))
	}
	return nil
}
