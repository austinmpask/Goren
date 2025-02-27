package utils

import (
	"fmt"
	"math"
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
