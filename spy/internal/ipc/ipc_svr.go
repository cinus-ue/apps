package ipc

import (
	"log"
	"time"

	"github.com/cinus-e/spy/internal/literr"
	ipc "github.com/james-barrow/golang-ipc"
)

const (
	success = "success"
	failure = "failure"
	msgType = 10
)

func StartServer() error {
	sc, err := ipc.StartServer("spy", nil)
	if err != nil {
		return err
	}

	for {
		m, err := sc.Read()
		if literr.CheckError(err) {
			continue
		}
		if m.MsgType > 0 {
			message := string(m.Data)
			log.Printf("msgtype: %d status: %s message: %s\n", m.MsgType, m.Status, message)
			if m.MsgType == msgType {
				ret, err := HandleCommand(message)
				if err != nil {
					literr.CheckError(sc.Write(msgType, []byte(ret+", "+err.Error())))
				} else {
					literr.CheckError(sc.Write(msgType, []byte(ret)))
				}
			}
		}
		time.Sleep(3 * time.Millisecond)
	}
}
