package utils

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Clear terminal window
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// Convert a duration to millisecond value
func InMs(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}

// Convert degrees to radians
func DegToRad(deg float64) float64 {
	return (deg * math.Pi) / 180
}

// Initialize frame buffer and depth buffer
func CreateBuffers(x uint16, y uint16) ([][]string, [][]float64) {

	fb := make([][]string, y)
	db := make([][]float64, y)
	for i := range fb {
		fb[i] = make([]string, x)
		db[i] = make([]float64, x)
	}
	return fb, db
}

// Cross platform terminal window determination
func GetTerminalSize() (rows uint16, cols uint16) {
	// Determine OS
	switch runtime.GOOS {
	// Do different window size methods on Windows
	case "windows":
		return winTerminalSize()
	default:
		return unixTerminalSize() // Mac/Linux
	}
}

// Powershell cmd for windows terminal
func winTerminalSize() (rows, cols uint16) {
	cmd := exec.Command("powershell", "-command",
		"&{$host.ui.rawui.WindowSize.Width}, {$host.ui.rawui.WindowSize.Height}")
	out, _ := cmd.Output()

	// Split into row and col
	parts := strings.Split(strings.TrimSpace(string(out)), "\n")

	r, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	c, _ := strconv.Atoi(strings.TrimSpace(parts[1]))

	return uint16(r), uint16(c)
}

// Do same thing but with stty for mac/linux
func unixTerminalSize() (rows, cols uint16) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, _ := cmd.Output()

	split := strings.Split(strings.TrimSpace(string(out)), " ")

	r, _ := strconv.Atoi(split[0])
	c, _ := strconv.Atoi(split[1])

	return uint16(r), uint16(c)
}
