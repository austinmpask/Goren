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
