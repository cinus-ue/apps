package tracker

import (
	"encoding/csv"
	"os"
	"strings"
	"time"

	"github.com/cinus-e/spy/internal/literr"
	"github.com/cinus-e/spy/internal/system"
	"github.com/cinus-e/spy/internal/util"
)

const (
	activityPrefix = "activity"
	processPrefix  = "process"
)

type timeEntry struct {
	Start     time.Time
	Last      time.Time
	text      []string
	SpendTime float64 //minutes
}

type Tracker struct {
	activity                    map[string]*timeEntry
	process                     map[string]*timeEntry
	actistate, procstate        bool
	actisave, procsave          time.Time
	watchInterval, saveInterval int
}

func NewTracker(watchInterval, saveInterval int) *Tracker {
	return &Tracker{activity: make(map[string]*timeEntry),
		process:       make(map[string]*timeEntry),
		actisave:      time.Now(),
		procsave:      time.Now(),
		watchInterval: watchInterval,
		saveInterval:  saveInterval,
	}
}

func (t *Tracker) TrackingActivity() {
	t.actistate = true
	for t.actistate {
		appName, windowText := TrackingWindow()
		if appName != "" {
			if e, ok := t.activity[appName]; ok {
				e.Last = time.Now()
				e.text = util.Unique(append(e.text, windowText))
			} else {
				t.activity[appName] = &timeEntry{time.Now(), time.Now(), []string{windowText}, 0}
			}
		}
		if time.Now().Sub(t.actisave).Minutes() > float64(t.saveInterval) {
			t.saveActivity()
		}
		time.Sleep(time.Duration(t.watchInterval) * time.Second)
	}
}

func (t *Tracker) TrackingProcess() {
	t.procstate = true
	for t.procstate {
		processList, err := system.Processes()
		if !literr.CheckError(err) {
			for _, p := range processList {
				name := p.Exe
				if e, ok := t.process[name]; ok {
					e.Last = time.Now()
				} else {
					t.process[name] = &timeEntry{time.Now(), time.Now(), nil, 0}
				}
			}
			if time.Now().Sub(t.procsave).Minutes() > float64(t.saveInterval) {
				t.saveProcess()
			}
		}
		time.Sleep(time.Duration(t.watchInterval) * time.Second)
	}
}

func (t *Tracker) Stop() {
	t.actistate = false
	t.procstate = false
}

func (t *Tracker) saveActivity() {
	t.actisave = time.Now()
	file, err := os.Create("./" + util.FileNameFormat(activityPrefix, ".csv"))
	if literr.CheckError(err) {
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	for name, entry := range t.activity {
		literr.CheckError(writer.Write([]string{name,
			util.FormatRFC3339(entry.Start),
			util.FormatRFC3339(entry.Last),
			strings.Join(entry.text, "; "),
		}))
	}
	writer.Flush()
	t.activity = make(map[string]*timeEntry)
}

func (t *Tracker) saveProcess() {
	t.procsave = time.Now()
	file, err := os.Create("./" + util.FileNameFormat(processPrefix, ".csv"))
	if literr.CheckError(err) {
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	for name, entry := range t.process {
		literr.CheckError(writer.Write([]string{name,
			util.FormatRFC3339(entry.Start),
			util.FormatRFC3339(entry.Last)},
		))
	}
	writer.Flush()
	t.process = make(map[string]*timeEntry)
}
