package utils

import (
	"os"
	"os/exec"
	"sync"
	"time"
)

var Wg = sync.WaitGroup{}

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func InMs(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}

func InitFrameBuffer(x uint8, y uint8) [][]string {

	buffer := make([][]string, y)
	for i := range buffer {
		buffer[i] = make([]string, x)
	}
	return buffer

}

func ApplyCamMatrix(camX float64, camY float64, camZ float64, x float64, y float64, z float64) []float64 {
	return []float64{x - camX, y - camY, z - camZ}
}

func ApplyProjectionMatrix(vert []float64, xConst float64, yConst float64, zConst float64, wConst float64) []float64 {
	// Assumed w = 1
	newX := vert[0] * xConst
	newY := vert[1] * yConst
	newZ := vert[2]*zConst + wConst
	newW := -1 * vert[2]

	return []float64{newX, newY, newZ, newW}
}

func ApplyNdcMatrix(clipLoc []float64) []float64 {
	// Divide points by w', discard y', w'
	return []float64{clipLoc[0] / clipLoc[3], clipLoc[1] / clipLoc[3]}

}

func NdcToScreen(ndcLoc []float64, screenX uint8, screenY uint8) []float64 {
	// Transform to screensize for final location
	x := ((ndcLoc[0] + 1) / 2) * (float64(screenX) - 1)
	y := (1 - ((ndcLoc[1] + 1) / 2)) * (float64(screenY) - 1)
	return []float64{x, y}

}
