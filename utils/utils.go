package utils

import (
	"math"
	"os"
	"os/exec"
	"time"
)

var ColorMap = map[string]string{
	// Reset
	"Reset": "\033[0m",

	// Red shades (using 8-bit colors)
	"Red1":  "\033[38;5;52m", // Dark red
	"Red2":  "\033[38;5;88m",
	"Red3":  "\033[38;5;124m",
	"Red4":  "\033[38;5;160m",
	"Red5":  "\033[38;5;196m", // Standard red
	"Red6":  "\033[38;5;197m",
	"Red7":  "\033[38;5;198m",
	"Red8":  "\033[38;5;199m",
	"Red9":  "\033[38;5;200m",
	"Red10": "\033[38;5;201m", // Light red/pink

	// Green shades
	"Green1":  "\033[38;5;22m", // Dark green
	"Green2":  "\033[38;5;28m",
	"Green3":  "\033[38;5;34m",
	"Green4":  "\033[38;5;40m",
	"Green5":  "\033[38;5;46m", // Standard green
	"Green6":  "\033[38;5;47m",
	"Green7":  "\033[38;5;48m",
	"Green8":  "\033[38;5;49m",
	"Green9":  "\033[38;5;50m",
	"Green10": "\033[38;5;51m", // Light green/cyan

	// Blue shades
	"Blue1":  "\033[38;5;17m", // Dark blue
	"Blue2":  "\033[38;5;18m",
	"Blue3":  "\033[38;5;19m",
	"Blue4":  "\033[38;5;20m",
	"Blue5":  "\033[38;5;21m", // Standard blue
	"Blue6":  "\033[38;5;27m",
	"Blue7":  "\033[38;5;33m",
	"Blue8":  "\033[38;5;39m",
	"Blue9":  "\033[38;5;45m",
	"Blue10": "\033[38;5;51m", // Light blue/cyan

	// Yellow shades
	"Yellow1":  "\033[38;5;58m", // Dark yellow/brown
	"Yellow2":  "\033[38;5;94m",
	"Yellow3":  "\033[38;5;136m",
	"Yellow4":  "\033[38;5;178m",
	"Yellow5":  "\033[38;5;220m", // Standard yellow
	"Yellow6":  "\033[38;5;221m",
	"Yellow7":  "\033[38;5;222m",
	"Yellow8":  "\033[38;5;223m",
	"Yellow9":  "\033[38;5;224m",
	"Yellow10": "\033[38;5;225m", // Light yellow

	// Magenta shades
	"Magenta1":  "\033[38;5;53m", // Dark magenta
	"Magenta2":  "\033[38;5;89m",
	"Magenta3":  "\033[38;5;125m",
	"Magenta4":  "\033[38;5;161m",
	"Magenta5":  "\033[38;5;197m", // Standard magenta
	"Magenta6":  "\033[38;5;198m",
	"Magenta7":  "\033[38;5;199m",
	"Magenta8":  "\033[38;5;200m",
	"Magenta9":  "\033[38;5;201m",
	"Magenta10": "\033[38;5;207m", // Light magenta

	// Cyan shades
	"Cyan1":  "\033[38;5;23m", // Dark cyan
	"Cyan2":  "\033[38;5;30m",
	"Cyan3":  "\033[38;5;37m",
	"Cyan4":  "\033[38;5;44m",
	"Cyan5":  "\033[38;5;51m", // Standard cyan
	"Cyan6":  "\033[38;5;50m",
	"Cyan7":  "\033[38;5;49m",
	"Cyan8":  "\033[38;5;48m",
	"Cyan9":  "\033[38;5;47m",
	"Cyan10": "\033[38;5;46m", // Light cyan/green

	// Gray shades
	"Gray1":  "\033[38;5;232m", // Almost black
	"Gray2":  "\033[38;5;236m",
	"Gray3":  "\033[38;5;240m",
	"Gray4":  "\033[38;5;244m",
	"Gray5":  "\033[38;5;248m", // Medium gray
	"Gray6":  "\033[38;5;252m",
	"Gray7":  "\033[38;5;253m",
	"Gray8":  "\033[38;5;254m",
	"Gray9":  "\033[38;5;255m",
	"Gray10": "\033[38;5;231m", // Almost white

	// White (just standard white, no shades)
	"White": "\033[38;5;231m",
}

// Clear terminal window TODO: make cross platform
func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// Convert a duration to millisecond value
func InMs(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}

// Convert degrees to radians
func DegToRad(deg float64) float64 {
	return (deg * math.Pi) / 180
}

// Initialize nested slices to emtpy strings
func InitFrameBuffer(x uint16, y uint16) [][]string {

	buffer := make([][]string, y)
	for i := range buffer {
		buffer[i] = make([]string, x)
	}
	return buffer

}

func InitDepthBuffer(x uint16, y uint16) [][]float64 {

	buffer := make([][]float64, y)
	for i := range buffer {
		buffer[i] = make([]float64, x)
	}
	return buffer

}

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
