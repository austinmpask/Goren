package main

import (
	"fmt"
	"go3d/actors"
	"go3d/display"
	"go3d/input"
	"go3d/utils"
	"math"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Intercept exit and perform cleanup
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM)

	// Run cleanup functions
	go func() {
		<-osSig
		fmt.Println("Quitting")
		// Exit raw mode
		input.RestoreTerminal()
		// Put cursor back
		fmt.Println("\033[?25h")
		os.Exit(0)
	}()

	// Listen to keys and manage key expiry to simulate "key depress"
	go input.ListenKeys()
	go input.ManageKeys()

	// Actual program
	view := display.CreateView(256, 224, 60, .1, true)

	panda := actors.CreateObject(utils.ParseObj("./rat.obj"), 0, 3, 8, .2, "Green")
	miniRat := actors.CreateObject(utils.ParseObj("./rat.obj"), -5, 3, 23, .1, "Red")
	view.RegisterObject(*panda)
	view.RegisterObject(*miniRat)

	// Main loop
	for {

		view.StartFrame()
		view.ClearBuffer()

		// panda.Translate(0, 0, .1)
		panda.Rotate(0, 1, 0)

		// Scene logic

		switch input.Key {
		case "w":
			view.MoveCam(math.Cos(utils.DegToRad(90-view.CamRot[1]))*view.CamMoveSpeed, 0, math.Sin(utils.DegToRad(90-view.CamRot[1]))*-view.CamMoveSpeed)
		case "s":
			view.MoveCam(math.Cos(utils.DegToRad(90-view.CamRot[1]))*-view.CamMoveSpeed, 0, math.Sin(utils.DegToRad(90-view.CamRot[1]))*view.CamMoveSpeed)
		case "d":
			view.MoveCam(math.Cos(utils.DegToRad(-view.CamRot[1]))*view.CamMoveSpeed, 0, math.Sin(utils.DegToRad(-view.CamRot[1]))*-view.CamMoveSpeed)
		case "a":
			view.MoveCam(math.Cos(utils.DegToRad(-view.CamRot[1]))*-view.CamMoveSpeed, 0, math.Sin(utils.DegToRad(-view.CamRot[1]))*view.CamMoveSpeed)
		case " ":
			view.MoveCam(0, view.CamMoveSpeed, 0)
		case "z":
			view.MoveCam(0, -view.CamMoveSpeed, 0)
		case "i":
			view.RotateCam(-5*view.CamMoveSpeed, 0, 0)
		case "k":
			view.RotateCam(5*view.CamMoveSpeed, 0, 0)
		case "l":
			view.RotateCam(0, 10*view.CamMoveSpeed, 0)
		case "j":
			view.RotateCam(0, -10*view.CamMoveSpeed, 0)
		}

		view.PrepBuffer()
		view.DrawBuffer()
		// view.DrawDebug()
		view.EndFrame()
		view.FrameSync()

	}

}
