package main

import (
	"go3d/actors"
	"go3d/display"
	"go3d/input"
	"go3d/utils"
)

func main() {

	// Create a scene
	scene := display.CreateView(256, 224, 60, .1)

	// Import models
	bigRat := actors.CreateObject(utils.ParseObj("./rat.obj"), 0, 0, 8, .2, "Green")
	miniRat := actors.CreateObject(utils.ParseObj("./rat.obj"), -5, 3, 23, .1, "Red")

	// Create scene lighting
	light := actors.CreatePointLight(0, 10, -10, .7, 25)

	// Register within the scene
	scene.RegisterObject(*bigRat)
	scene.RegisterObject(*miniRat)
	scene.RegisterPointLight(light)

	// Listen to keyboard input
	input.ListenKeys()
	// Main loop
	for {

		scene.StartFrame()
		scene.ClearBuffer()
		utils.ClearScreen()

		// Scene logic
		bigRat.Rotate(0, 1, 0)
		// End Scene logic

		scene.HandleInput()
		scene.PrepBuffer()
		scene.DrawBuffer()

		scene.EndFrame()
		// Wait until target frametime has expired to minimize screen tearing
		scene.FrameSync("sleep", -3)

	}

}
