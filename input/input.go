package input

import (
	"fmt"
	"go3d/utils"
	"os"
	"syscall"
	"time"
	"unsafe"
)

// Terminal settings
var oldState, newState syscall.Termios
var fd = syscall.Stdin

// Globals for current key pressed on keyboard
var Key string
var KeyTimeStamp time.Time

// Save the current key pressed to a global variable
// TODO: Allow for key chords
func ListenKeys() {
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&oldState)))

	// Modify settings to allow for input without enter
	newState = oldState
	newState.Lflag &^= syscall.ICANON | syscall.ECHO
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&newState)))

	var b [1]byte
	go func() {

		// os.Stdin.Read is blocking, so this loop wont destroy CPU
	keyRead:
		for {
			os.Stdin.Read(b[:])

			// Save one key, and save timestamp for expiry
			Key = string(b[0])
			KeyTimeStamp = time.Now()

			switch b[0] {
			case 'q':
				// Run cleanup and exit program on 'q'
				RestoreTerminal()
				// Stop loop
				break keyRead

			// Debug toggles
			case 'e':
				utils.Debug = !utils.Debug

			case '7':
				utils.DrawVerts = !utils.DrawVerts

			case '8':
				utils.DrawWire = !utils.DrawWire

			case '9':
				utils.RenderFace = !utils.RenderFace

			case '0':
				utils.RenderLighting = !utils.RenderLighting

			}
		}

	}()
	go ManageKeys()
}

// Simulate key releases and bridge initial skip when holding one key down
// TODO: Make this more responsive
func ManageKeys() {

	// Minimum time to "hold" a key before "release"
	waitTime := 500
	for {
		elapsed := time.Since(KeyTimeStamp)

		// Compare with elapsed time in addition to sleeping for the wait time to account for different key presses
		if elapsed.Milliseconds() > int64(waitTime) {
			Key = ""
			KeyTimeStamp = time.Now()
		}
		// Wait
		time.Sleep(time.Duration(waitTime) * time.Millisecond)
	}
}

// Run cleanup upon exit
func RestoreTerminal() {
	// Exit raw mode
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&oldState)))
	// Reset terminal color
	fmt.Println("\033[0m")
	// Restore cursor to terminal
	fmt.Println("\033[?25h")
	// Wipe screen
	utils.ClearScreen()
	// Exit program
	os.Exit(0)
}
