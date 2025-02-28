package main

import (
	"go3d/actors"
	"go3d/display"
	"go3d/utils"
)

// Create a scene with some static objects
func createDemoSceneWithStatics() *display.View {

	// Create a scene
	scene := display.CreateView(30, .4)

	// Import models

	house := actors.CreateObject(utils.ParseObj("./models/house.obj"), 0, 0, 8, 10, "Gray")
	cabana := actors.CreateObject(utils.ParseObj("./models/cabana.obj"), 0, -6.5, 20, 2, "Gray")
	log := actors.CreateObject(utils.ParseObj("./models/log.obj"), 4, -5.5, 10, 2.2, "Gray")
	tele1 := actors.CreateObject(utils.ParseObj("./models/telePole.obj"), 5, -6.5, 30, 1, "Gray")
	tele2 := actors.CreateObject(utils.ParseObj("./models/telePole.obj"), 5, -6.5, 14, 1, "Gray")
	tele3 := actors.CreateObject(utils.ParseObj("./models/telePole.obj"), 5, -6.5, 0, 1, "Gray")
	tree := actors.CreateObject(utils.ParseObj("./models/tree.obj"), -17, 4, 26, 16, "Green")
	tree2 := actors.CreateObject(utils.ParseObj("./models/tree.obj"), -22, 4, 19, 20, "Green")
	grass := actors.CreateObject(utils.ParseObj("./models/grass.obj"), 1, -6.75, 9, .05, "Yellow")
	grass2 := actors.CreateObject(utils.ParseObj("./models/grass.obj"), -10.5, -6.75, 16, .05, "Yellow")
	grass3 := actors.CreateObject(utils.ParseObj("./models/grass.obj"), -12.5, -6.75, 14, .05, "Yellow")
	grass4 := actors.CreateObject(utils.ParseObj("./models/grass.obj"), -10.5, -6.75, 0, .05, "Yellow")

	house.Rotate(0, 90, 0)
	house.Translate(-15, 0, 0)
	cabana.Translate(-13, 0, 0)
	log.Rotate(0, 45, 0)

	tele1.Rotate(0, 90, 0)
	tele2.Rotate(0, 90, 0)
	tele3.Rotate(0, 90, 0)

	scene.MoveCam(8, 10, -2)
	scene.RotateCam(20, -30, 0)

	houseLight := &actors.Light{
		LightX:    0,
		LightY:    10,
		LightZ:    -15,
		Intensity: .60,
		Falloff:   30,
	}

	// Register within the scene
	scene.RegisterObject(house)
	scene.RegisterObject(cabana)
	scene.RegisterObject(log)
	scene.RegisterObject(tele1)
	scene.RegisterObject(tele2)
	scene.RegisterObject(tele3)
	scene.RegisterObject(tree)
	scene.RegisterObject(tree2)
	scene.RegisterObject(grass)
	scene.RegisterObject(grass2)
	scene.RegisterObject(grass3)
	scene.RegisterObject(grass4)

	scene.RegisterLight(houseLight)
	return scene
}
