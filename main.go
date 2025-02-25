package main

import (
	"fmt"
	"go3d/actors"
	"go3d/display"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Intercept exit and perform cleanup
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-osSig
		fmt.Println("Quitting")
		fmt.Print("\033[?25h")
		os.Exit(0)
	}()

	// Actual program
	view := display.DefaultView()
	p := actors.DefaultPoint()
	view.RegisterActor(p)

	p1 := &actors.Point{
		X: 0,
		Y: 0,
		Z: -4,
	}

	view.RegisterActor(p1)

	for {

		view.StartFrame()
		view.ClearBuffer()

		view.PrepBuffer()
		view.DrawBuffer()
		view.DrawDebug()
		view.EndFrame()
		view.FrameSync()

	}

}
