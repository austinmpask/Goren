package main

import (
	"fmt"
	"go3d/actors"
	"go3d/display"
	"go3d/input"
	"go3d/utils"
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

	// Import models
	bigRat := actors.CreateObject(utils.ParseObj("./rat.obj"), 0, 0, 8, .2, "Green")
	miniRat := actors.CreateObject(utils.ParseObj("./rat.obj"), -5, 3, 23, .1, "Red")

	// Create scene lighting
	light := actors.CreatePointLight(0, 10, -10, .7, 25)

	// Register within the scene
	view.RegisterObject(*bigRat)
	view.RegisterObject(*miniRat)
	view.RegisterPointLight(light)

	// Main loop
	for {

		view.StartFrame()
		view.ClearBuffer()
		utils.ClearScreen()

		// Scene logic
		bigRat.Rotate(0, 1, 0)
		// End Scene logic

		view.HandleInput()
		view.PrepBuffer()
		view.DrawBuffer()

		view.EndFrame()
		// Wait until target frametime has expired to minimize screen tearing
		view.FrameSync("loop", 0)

	}

}
