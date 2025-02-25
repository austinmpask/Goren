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

	// t := actors.CreateTriangle([]float64{0, 0, -5}, []float64{5, 0, -10}, []float64{0, 5, -5})

	triangles := []actors.Actor{
		// Trunk (simplified cylinder approximation)
		actors.CreateTriangle([]float64{-0.5, 0, 0}, []float64{0.5, 0, 0}, []float64{0.5, 5, 0}),
		actors.CreateTriangle([]float64{-0.5, 0, 0}, []float64{0.5, 5, 0}, []float64{-0.5, 5, 0}),
		actors.CreateTriangle([]float64{-0.5, 0, -0.5}, []float64{0.5, 0, -0.5}, []float64{0.5, 5, -0.5}),
		actors.CreateTriangle([]float64{-0.5, 0, -0.5}, []float64{0.5, 5, -0.5}, []float64{-0.5, 5, -0.5}),
		actors.CreateTriangle([]float64{-0.5, 0, 0}, []float64{-0.5, 5, 0}, []float64{-0.5, 5, -0.5}),
		actors.CreateTriangle([]float64{-0.5, 0, 0}, []float64{-0.5, 5, -0.5}, []float64{-0.5, 0, -0.5}),
		actors.CreateTriangle([]float64{0.5, 0, 0}, []float64{0.5, 5, -0.5}, []float64{0.5, 5, 0}),
		actors.CreateTriangle([]float64{0.5, 0, 0}, []float64{0.5, 0, -0.5}, []float64{0.5, 5, -0.5}),
		actors.CreateTriangle([]float64{0, 5, 0}, []float64{1, 6, 0}, []float64{0.5, 5.5, 0.5}),
		actors.CreateTriangle([]float64{0, 5, 0}, []float64{0.5, 5.5, 0.5}, []float64{0, 6, 0.5}),
		actors.CreateTriangle([]float64{0, 5, 0}, []float64{0, 6, 0.5}, []float64{-0.5, 5.5, 0.5}),
		actors.CreateTriangle([]float64{0, 5, 0}, []float64{-0.5, 5.5, 0.5}, []float64{-1, 6, 0}),
		actors.CreateTriangle([]float64{0, 5, 0}, []float64{-1, 6, 0}, []float64{-0.5, 5.5, -0.5}),
		actors.CreateTriangle([]float64{0, 5, 0}, []float64{-0.5, 5.5, -0.5}, []float64{0, 6, -0.5}),
		actors.CreateTriangle([]float64{0, 5, 0}, []float64{0, 6, -0.5}, []float64{0.5, 5.5, -0.5}),
		actors.CreateTriangle([]float64{0, 5, 0}, []float64{0.5, 5.5, -0.5}, []float64{1, 6, 0}),
		actors.CreateTriangle([]float64{0, 4, 0}, []float64{1, 5, 1}, []float64{-1, 5, 1}),
		actors.CreateTriangle([]float64{0, 4, 0}, []float64{-1, 5, 1}, []float64{-1, 5, -1}),
		actors.CreateTriangle([]float64{0, 4, 0}, []float64{-1, 5, -1}, []float64{1, 5, -1}),
		actors.CreateTriangle([]float64{0, 4, 0}, []float64{1, 5, -1}, []float64{1, 5, 1}),
		// ... Add more branches and leaves (using similar triangle patterns)
	}

	for _, t := range triangles {
		view.RegisterActor(t)

	}

	for {

		view.StartFrame()
		view.ClearBuffer()
		view.MoveCam(view.CamMoveSpeed, 0, -view.CamMoveSpeed)

		view.PrepBuffer()
		view.DrawBuffer()
		view.DrawDebug()
		view.EndFrame()
		view.FrameSync()

	}

}
