package utils

import "math"

// Apply x rotation matrix
func RotateXTransform(vec []float64, rot []float64) []float64 {
	x := vec[0]
	y := vec[1]
	z := vec[2]

	xRot := DegToRad(rot[0])

	// Apply matrix
	xTrans := x
	yTrans := math.Cos(xRot)*y - math.Sin(xRot)*z
	zTrans := math.Sin(xRot)*y + math.Cos(xRot)*z

	return []float64{xTrans, yTrans, zTrans}

}

// Apply y rotation matrix
func RotateYTransform(vec []float64, rot []float64) []float64 {
	x := vec[0]
	y := vec[1]
	z := vec[2]

	yRot := DegToRad(rot[1])

	// Apply matrix
	xTrans := math.Cos(yRot)*x + math.Sin(yRot)*z
	yTrans := y
	zTrans := -math.Sin(yRot)*x + math.Cos(yRot)*z

	return []float64{xTrans, yTrans, zTrans}

}

// Apply z rotation matrix
func RotateZTransform(vec []float64, rot []float64) []float64 {
	x := vec[0]
	y := vec[1]
	z := vec[2]

	zRot := DegToRad(rot[2])

	// Apply matrix
	xTrans := math.Cos(zRot)*x - math.Sin(zRot)*y
	yTrans := math.Sin(zRot)*x + math.Cos(zRot)*y
	zTrans := z

	return []float64{xTrans, yTrans, zTrans}

}

// Rotate a 3D vector by Z -> Y -> X
func RotateVecXYZ(vec []float64, rot []float64) []float64 {

	z := RotateZTransform(vec, rot)
	y := RotateYTransform(z, rot)
	x := RotateXTransform(y, rot)
	return x
}

// Transform from object space to world space
func ApplyWorldMatrix(vert []float64, objX float64, objY float64, objZ float64, objScale float64, objRot []float64) []float64 {

	// Rotate
	rotated := RotateVecXYZ(vert, objRot)
	//Scale
	sX := rotated[0] * objScale
	sY := rotated[1] * objScale
	sZ := rotated[2] * objScale

	// Translate
	xWorld := sX - objX
	yWorld := sY - objY
	zWorld := sZ - objZ

	return []float64{xWorld, yWorld, zWorld}
}

// Apply camera translation and rotation to a vector
func ApplyCamMatrix(camX float64, camY float64, camZ float64, camRot []float64, x float64, y float64, z float64) []float64 {
	// Translate
	translatedVec := []float64{x - camX, y - camY, z - camZ}

	// Rotate
	r := RotateVecXYZ(translatedVec, camRot)

	return r

}

// Project vector from camspace to clip space
func ApplyProjectionMatrix(vert []float64, xConst float64, yConst float64, zConst float64, wConst float64) []float64 {
	// Assumed w = 1
	newX := vert[0] * xConst
	newY := vert[1] * yConst
	newZ := vert[2]*zConst + wConst
	newW := -1 * vert[2]

	return []float64{newX, newY, newZ, newW}
}

// Transform vec from clip space to norm. device coordinates
func ApplyNdcMatrix(clipLoc []float64) []float64 {
	// Divide points by w', discard y', w'
	return []float64{clipLoc[0] / clipLoc[3], clipLoc[1] / clipLoc[3]}

}

// Transform NDC to screenspace coordinates
func NdcToScreen(ndcLoc []float64, screenX uint16, screenY uint16) []float64 {
	// Transform to screensize for final location
	x := ((ndcLoc[0] + 1) / 2) * (float64(screenX) - 1)
	y := (1 - ((ndcLoc[1] + 1) / 2)) * (float64(screenY) - 1)
	return []float64{x, y}

}
