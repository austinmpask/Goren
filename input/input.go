package input

import (
	"os"
	"syscall"
	"time"
	"unsafe"
)

var oldState, newState syscall.Termios
var fd = syscall.Stdin

// Globals for current key pressed on keyboard
var Key string
var KeyTimeStamp time.Time

// Save the current key pressed to a global variable
// TODO: Allow for simultaneous keypresses
func ListenKeys() {
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&oldState)))

	// Modify settings for raw input
	newState = oldState
	newState.Lflag &^= syscall.ICANON | syscall.ECHO
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&newState)))

	var b [1]byte
	for {
		os.Stdin.Read(b[:])

		// Save one key, and save timestamp for expiry
		Key = string(b[0])
		KeyTimeStamp = time.Now()

		if b[0] == 'q' {
			RestoreTerminal()
			break
		}
	}
}

// Maintain natural inputs with the terminal input deadzone limitation
// TODO: Make this better
func ManageKeys() {
	for {
		elapsed := time.Since(KeyTimeStamp)

		if elapsed.Milliseconds() > 500 {
			Key = ""
			KeyTimeStamp = time.Now()
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func RestoreTerminal() {
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&oldState)))
}
