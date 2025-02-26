package display

import (
	"fmt"
	"go3d/actors"
	"go3d/utils"
	"math"
	"strings"
	"time"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

// Monochrome colorspace
var pixel = map[uint8]string{
	0: "  ",
	1: "██",
}

// View defines the screenspace printed to the terminal, with any debug information
type View struct {
	Xpx         uint16
	Ypx         uint16
	TargetFPS   uint8
	FrameBuffer [][]string

	Fov      uint8
	NearClip float64
	FarClip  float64

	XProjConst float64
	YProjConst float64
	ZProjConst float64
	WProjConst float64

	CamX   float64
	CamY   float64
	CamZ   float64
	CamRot []float64

	CamMoveSpeed float64

	FrameStart time.Time
	FrameTime  time.Duration
	FrameEnd   time.Duration

	MaxFrameTime time.Duration

	Xborder string
	Actors  []actors.Triangle
}

func DefaultView() *View {

	v := View{
		Xpx: 640,
		Ypx: 240,
		// Xpx:          80,
		// Ypx:          40,
		TargetFPS:    65,
		Fov:          90,
		CamX:         0,
		CamY:         3,
		CamZ:         10,
		CamRot:       []float64{0, 0, 0},
		NearClip:     -1,
		FarClip:      -150,
		CamMoveSpeed: 0.25,
	}

	// Calc max frame time
	v.MaxFrameTime = v.CalcMaxFrameTime()

	// Initialize buffer
	v.FrameBuffer = utils.InitFrameBuffer(v.Xpx, v.Ypx)

	// Calculate projection constants
	v.CalcProjectionConstants()

	// Initialize screen border
	v.Xborder = strings.Repeat(pixel[1], int(v.Xpx)+2)

	// Remove cursor
	fmt.Print("\033[?25l")

	v.ClearBuffer()

	return &v
}

// Add actors to the scene
func (v *View) RegisterActor(t actors.Triangle) {
	v.Actors = append(v.Actors, t)
}

func (v *View) RegisterObject(o actors.Object) {
	for _, tri := range o.Tris {
		v.RegisterActor(tri)
	}
}

// Calculate the maximum allowable frametime to maintain the target framerate in MS
func (v *View) CalcMaxFrameTime() time.Duration {

	fps := time.Duration(v.TargetFPS)
	return time.Second / fps
}

// TODO precompute this
func (v *View) Aspect() float64 {
	return float64(v.Xpx) / float64(v.Ypx)
}

// Set the FrameBuffer to empty pixels
func (v *View) ClearBuffer() {
	for i := range v.FrameBuffer {
		for j := range v.FrameBuffer[i] {
			v.FrameBuffer[i][j] = pixel[0]
		}
	}
	utils.ClearScreen()

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

// Precompute projection matrix constants for camera aspects which do not change
func (v *View) CalcProjectionConstants() {
	// Projection matrix [0,0]: 1/(Aspect Ratio * Tan(FOV/2))
	tanHalfFov := math.Tan((math.Pi * float64(v.Fov) / 360))
	v.XProjConst = 1 / (v.Aspect() * tanHalfFov)

	// Projection matrix [1,1]: 1/Tan(FOV/2)
	v.YProjConst = 1 / tanHalfFov

	// Projection matrix [2,2]: -1*(far + near)/(far - near)
	v.ZProjConst = -1 * ((v.FarClip + v.NearClip) / (v.FarClip - v.NearClip))

	// Projection matrix [3,2]: (2*far*near)/(far-near)
	v.WProjConst = (2 * v.FarClip * v.NearClip) / (v.FarClip - v.NearClip)

}

// Apply camera transformations -> projection transformations -> NDC transformations -> screenspace transformations
// Add results to the framebuffer (verts & lines)
func (v *View) PrepBuffer() {

	for _, a := range v.Actors {
		parent := a.ObjRef

		//Save raster verts for connecting with lines after vertex pass
		var rasterVerts [][]uint16

		// Calculate vertecies
		for _, vert := range a.Verts() {

			// Convert to worldspace
			worldVert := utils.ApplyWorldMatrix(vert, parent.ObjX, parent.ObjY, parent.ObjZ, parent.Scale, parent.Rot)

			camSpaceVert := utils.ApplyCamMatrix(v.CamX, v.CamY, v.CamZ, v.CamRot, worldVert[0], worldVert[1], worldVert[2])

			// fmt.Printf("Camspace Vert: %v\n", camSpaceVert)
			clipVert := utils.ApplyProjectionMatrix(camSpaceVert, v.XProjConst, v.YProjConst, v.ZProjConst, v.WProjConst)

			// Discard if W out of bounds

			if clipVert[3] > math.Abs(v.FarClip) || clipVert[3] < math.Abs(v.NearClip) {
				continue
			}
			// fmt.Printf("Clip Vert: %v\n", clipVert)
			ndcVert := utils.ApplyNdcMatrix(clipVert)
			// Discard if out of bounds
			ooB := false
			for _, v := range ndcVert {
				if v > 1 || v < -1 {
					ooB = true
					break
				}

			}
			if ooB {
				continue
			}
			// fmt.Printf("NDC Vert: %v\n", ndcVert)
			ssVert := utils.NdcToScreen(ndcVert, v.Xpx, v.Ypx)
			// fmt.Printf("Screenspace Vert: %v\n", ssVert)

			xVert := uint16(math.Round(ssVert[0]))
			yVert := uint16(math.Round(ssVert[1]))

			// Save final 2D vertex for drawing lines
			rasterVerts = append(rasterVerts, []uint16{xVert, yVert})

			// Load vertecies to buffer
			v.TouchBuffer(Blue, xVert, yVert)

		}

		// Draw lines between 2D verts with bresenhams alg
		if len(rasterVerts) > 1 {

			// Keep track of connected points
			connected := make(map[int]map[int]bool)

			// Iterate through each vertex, connecting with neighbors and skipping if the reverse has been done
			for i := range len(rasterVerts) {
				for j := range len(rasterVerts) {
					if i != j {
						// Is i in connected
						if c1, ok := connected[i]; ok {
							// Is j in connected[i]
							if _, ok := c1[j]; !ok {
								v.DrawLine(rasterVerts[i], rasterVerts[j])
								// Record
								// If j in connected
								if _, ok := connected[j]; !ok {
									connected[j] = make(map[int]bool)
								}
								connected[j][i] = true
							}
						} else {
							v.DrawLine(rasterVerts[i], rasterVerts[j])
							// Record
							// If j in connected
							if _, ok := connected[j]; !ok {
								connected[j] = make(map[int]bool)

							}
							connected[j][i] = true
						}
					}

				}

			}
		}

	}

}

// Adds contiguous line to framebuffer between 2 points w/ Bresenhams alg
func (v *View) DrawLine(start []uint16, end []uint16) {

	startX := start[0]
	startY := start[1]

	endX := end[0]
	endY := end[1]

	// Pixels that will be drawn to buffer after calculations
	var pixels [][]uint16

	// Case of verical line
	if startX == endX {
		yMin := min(startY, endY) + 1
		yMax := max(startY, endY)

		// Save pixels to draw along vertical line
		for y := yMin; y < yMax; y++ {
			pixels = append(pixels, []uint16{startX, y})
		}
	} else if startY == endY {
		// Case of flat line
		xMin := min(startX, endX) + 1
		xMax := max(startX, endX)

		// Save pixels to draw along horizontal line
		for x := xMin; x < xMax; x++ {
			pixels = append(pixels, []uint16{x, startY})
		}
	} else {

		// Bresenhams alg for other slopes

		// Calculate slope
		m := (float64(endY) - float64(startY)) / (float64(endX) - float64(startX))

		// Iterate Y for slope 1 or higher as each Y coordinate will have only one pixel
		if math.Abs(m) >= 1 {

			// Invert slope
			c := 1 / m

			// Skip first pixel as vertex will be drawn there
			yMin := min(startY, endY) + 1
			yMax := max(startY, endY)

			for y := yMin; y < yMax; y++ {
				x := uint16(math.Round(c*(float64(y)-float64(startY)) + float64(startX)))
				pixels = append(pixels, []uint16{x, y})
			}

		} else {
			// Iterate over X for slope < 1 as each X coordinate will have only one pixel

			// Skip first pixel for vertex
			xMin := min(startX, endX) + 1
			xMax := max(startX, endX)

			for x := xMin; x < xMax; x++ {
				y := uint16(math.Round(m*(float64(x)-float64(startX)) + float64(startY)))
				pixels = append(pixels, []uint16{x, y})
			}

		}

	}
	// Load all the pixels to framebuffer from whichever line alg was used
	for _, p := range pixels {
		v.TouchBuffer(Cyan, p[0], p[1])
	}

}

// Prints buffer contents to screen
func (v *View) DrawBuffer() {
	fmt.Println(v.Xborder)
	for _, row := range v.FrameBuffer {
		// Print with Y border added
		fmt.Printf("%v%v%v\n", pixel[1], strings.Join(row, ""), pixel[1])
	}
	fmt.Println(v.Xborder)

}

// Prints performance information below view window
func (v *View) DrawDebug() {

	ftMs := utils.InMs(v.FrameTime)
	maxFtMs := utils.InMs(v.MaxFrameTime)
	frEndMs := utils.InMs((v.FrameEnd))

	var util float64
	var pfps float64
	var fps float64

	// Only calculate these values if denominator != 0
	if maxFtMs != 0 {
		util = 100 * ftMs / maxFtMs
	}

	if ftMs != 0 {
		pfps = 1000 / ftMs
	}

	if frEndMs != 0 {
		fps = 1000 / frEndMs
	}

	fmt.Printf("Frametime: %v ms\n", ftMs)
	fmt.Printf("Frametime util: %v %% \n", util)
	fmt.Printf("Potential FPS: %v\n", pfps)
	fmt.Printf("Real FPS: %v\n", fps)

}

// Safely populate a pixel in the buffer respecting xy bounds
func (v *View) TouchBuffer(color string, x uint16, y uint16) {

	if x < v.Xpx && y < v.Ypx {

		v.FrameBuffer[y][x] = fmt.Sprintf("%s%s\033[0m", color, pixel[1])
	}

}

// Log the time the frame calculations began
func (v *View) StartFrame() {
	v.FrameStart = time.Now()

}

// Log the time once the buffer and anything else was drawn to screen
func (v *View) EndFrame() {

	v.FrameTime = time.Since(v.FrameStart)
}

// Minimize screen tearing by waiting until frametime for the target framerate has elapsed before continuing
func (v *View) FrameSync() {
	frameTimeSlack := v.MaxFrameTime - v.FrameTime

	targetTime := time.Now().Add(frameTimeSlack)

	// Use for loop instead, more accurate scheduling than sleep. Sleep introduces significant drift at high refresh rates
	for time.Now().Before(targetTime) {
		//Wait
	}

	// Log the time that the entire frame ended
	v.FrameEnd = time.Since(v.FrameStart)
}
