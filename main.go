package main

import (
	"go3d/actors"
	"go3d/input"
	"go3d/utils"
)

func main() {

	// Make a view with some static objects for demo
	scene := createDemoSceneWithStatics()

	// Add some dynamic lighting
	headLights := &actors.Light{
		LightX:    0,
		LightY:    -4,
		LightZ:    -40,
		Intensity: .8,
		Falloff:   25,
	}
	carLight := &actors.Light{
		LightX:    0,
		LightY:    0,
		LightZ:    -50,
		Intensity: .3,
		Falloff:   10,
	}

	// Add some dynamic objects
	car := actors.CreateObject(utils.ParseObj("./models/car.obj"), 0, -6.5, 45, .1, "Red")
	hl1 := actors.CreateObject(utils.ParseObj("./models/headlight.obj"), -1, -5.2, 40.7, .08, "Yellow")
	hl2 := actors.CreateObject(utils.ParseObj("./models/headlight.obj"), 1, -5.2, 40.7, .08, "Yellow")
	hl1.Rotate(90, 0, 0)

	// Add these to the scene too
	scene.RegisterObject(hl1)
	scene.RegisterObject(hl2)
	scene.RegisterLight(headLights)
	scene.RegisterLight(carLight)
	scene.RegisterObject(car)

	// Listen to keyboard input
	input.ListenKeys()
	// Main demo loop
	for scene.CamY > 0 {

		scene.StartFrame()
		scene.ClearBuffer()

		// Scene logic for demo

		scene.MoveCam(0, -.05, -.05)
		scene.RotateCam(-.05, -.13, 0)
		car.Translate(0, 0, .15)
		hl1.Translate(0, 0, .15)
		hl2.Translate(0, 0, .15)
		headLights.Translate(0, 0, .15)
		carLight.Translate(0, 0, .15)
		// End Scene logic

		scene.HandleInput()
		scene.PrepBuffer()
		scene.DrawBuffer()
		scene.EndFrame()
		// Wait until target frametime has expired to minimize screen tearing
		scene.FrameSync("sleep", 0)

	}
	// Restore terminal settings after demo
	input.RestoreTerminal()

}
