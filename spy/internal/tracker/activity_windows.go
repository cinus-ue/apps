package tracker

import (
	"github.com/cinus-e/spy/internal/system"
)

func TrackingWindow() (appName, text string) {
	h := system.GetForegroundWindow()
	text = system.GetWindowText(h)
	if proc := system.FindProcessByPid(system.GetWindowThreadProcessId(h)); proc != nil {
		appName = proc.Exe
	}
	return
}
