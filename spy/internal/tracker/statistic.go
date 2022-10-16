package tracker

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cinus-e/spy/internal/literr"
	"github.com/cinus-e/spy/internal/util"
)

func ShowStatistics(files []string) {
	var acti []string
	var proc []string
	for _, name := range files {
		if strings.HasPrefix(filepath.Base(name), activityPrefix) {
			acti = append(acti, name)
		}
		if strings.HasPrefix(filepath.Base(name), processPrefix) {
			proc = append(proc, name)
		}
	}
	if len(acti) > 0 {
		fmt.Println("Activity Statistics:")
		for _, s := range statistic(acti) {
			fmt.Println(s)
		}
	}
	if len(proc) > 0 {
		fmt.Println("Process Statistics:")
		for _, s := range statistic(proc) {
			fmt.Println(s)
		}
	}

}

func statistic(files []string) []string {
	var ret []string
	tmp := make(map[string]*timeEntry)
	for _, name := range files {
		file, _ := os.Open(name)
		reader := csv.NewReader(file)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if literr.CheckError(err) {
				break
			}
			name = record[0]
			start := util.ParseRFC3339(record[1])
			last := util.ParseRFC3339(record[2])
			if name != "" {
				if e, ok := tmp[name]; ok {
					if e.Start.After(start) {
						e.Start = start
					}
					if e.Last.Before(last) {
						e.Last = last
					}
				} else {
					tmp[name] = &timeEntry{start, last, nil, 0}
				}
			}
		}
	}
	for _, e := range tmp {
		e.SpendTime = e.Last.Sub(e.Start).Minutes()
	}
	keys := make([]string, 0, len(tmp))

	for key := range tmp {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return tmp[keys[i]].SpendTime > tmp[keys[j]].SpendTime
	})
	for _, name := range keys {
		ret = append(ret, fmt.Sprintf("%f minutes %s", tmp[name].SpendTime, name))
	}
	return ret
}
