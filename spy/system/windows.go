//go:build windows

package system

import (
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

//TH32CS_SNAPPROCESS https://msdn.microsoft.com/de-de/library/windows/desktop/ms682489(v=vs.85).aspx
const TH32CS_SNAPPROCESS = 0x00000002

var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetWindowTextW           = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW     = user32.NewProc("GetWindowTextLengthW")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
)

//Get Active Window Title
func GetForegroundWindow() syscall.Handle {
	r1, _, _ := procGetForegroundWindow.Call()
	return syscall.Handle(r1)
}

func GetWindowTextLength(h syscall.Handle) int {
	ret, _, _ := procGetWindowTextLengthW.Call(
		uintptr(h))
	return int(ret)
}

func GetWindowText(h syscall.Handle) string {
	length := GetWindowTextLength(h) + 1
	buf := make([]uint16, length)
	procGetWindowTextW.Call(
		uintptr(h),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(length))
	return syscall.UTF16ToString(buf)
}

func GetWindowThreadProcessId(h syscall.Handle) int {
	var processId int
	procGetWindowThreadProcessId.Call(
		uintptr(h),
		uintptr(unsafe.Pointer(&processId)))

	return processId
}

// WindowsProcess is an implementation of Process for Windows.
type WindowsProcess struct {
	ProcessID       int
	ParentProcessID int
	Exe             string
}

func newWindowsProcess(e *windows.ProcessEntry32) WindowsProcess {
	// Find when the string ends for decoding
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return WindowsProcess{
		ProcessID:       int(e.ProcessID),
		ParentProcessID: int(e.ParentProcessID),
		Exe:             syscall.UTF16ToString(e.ExeFile[:end]),
	}
}

func Processes() ([]WindowsProcess, error) {
	handle, err := windows.CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(handle)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	// get the first process
	err = windows.Process32First(handle, &entry)
	if err != nil {
		return nil, err
	}

	results := make([]WindowsProcess, 0, 50)
	for {
		results = append(results, newWindowsProcess(&entry))

		err = windows.Process32Next(handle, &entry)
		if err != nil {
			// windows sends ERROR_NO_MORE_FILES on last process
			if err == syscall.ERROR_NO_MORE_FILES {
				return results, nil
			}
			return nil, err
		}
	}
}

func FindProcessByName(name string) *WindowsProcess {
	procs, err := Processes()
	if err != nil {
		return nil
	}
	for _, p := range procs {
		if strings.ToLower(p.Exe) == strings.ToLower(name) {
			return &p
		}
	}
	return nil
}

func FindProcessByPid(pid int) *WindowsProcess {
	procs, err := Processes()
	if err != nil {
		return nil
	}
	for _, p := range procs {
		if p.ProcessID == pid {
			return &p
		}
	}
	return nil
}
