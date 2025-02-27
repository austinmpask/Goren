package display

import (
	"go3d/input"
	"go3d/utils"
	"math"
)

// Camera rotation + transalation relative to normal
func (v *View) HandleInput() {
	switch input.Key {
	case "w":
		v.MoveCam(math.Cos(utils.DegToRad(90-v.CamRot[1]))*v.CamMoveSpeed, 0, math.Sin(utils.DegToRad(90-v.CamRot[1]))*-v.CamMoveSpeed)
	case "s":
		v.MoveCam(math.Cos(utils.DegToRad(90-v.CamRot[1]))*-v.CamMoveSpeed, 0, math.Sin(utils.DegToRad(90-v.CamRot[1]))*v.CamMoveSpeed)
	case "d":
		v.MoveCam(math.Cos(utils.DegToRad(-v.CamRot[1]))*v.CamMoveSpeed, 0, math.Sin(utils.DegToRad(-v.CamRot[1]))*-v.CamMoveSpeed)
	case "a":
		v.MoveCam(math.Cos(utils.DegToRad(-v.CamRot[1]))*-v.CamMoveSpeed, 0, math.Sin(utils.DegToRad(-v.CamRot[1]))*v.CamMoveSpeed)
	case " ":
		v.MoveCam(0, v.CamMoveSpeed, 0)
	case "z":
		v.MoveCam(0, -v.CamMoveSpeed, 0)
	case "i":
		v.RotateCam(-5*v.CamMoveSpeed, 0, 0)
	case "k":
		v.RotateCam(5*v.CamMoveSpeed, 0, 0)
	case "l":
		v.RotateCam(0, 10*v.CamMoveSpeed, 0)
	case "j":
		v.RotateCam(0, -10*v.CamMoveSpeed, 0)
	}
}

// Apply a 3d translation to camera
func (v *View) MoveCam(dx float64, dy float64, dz float64) {
	v.CamX += dx
	v.CamY += dy
	v.CamZ += dz
}

// Apply a new rotation value to camera. Rotation transformations occur during rendering
func (v *View) RotateCam(rx float64, ry float64, rz float64) {
	v.CamRot[0] += rx
	v.CamRot[1] += ry
	v.CamRot[2] += rz
}
