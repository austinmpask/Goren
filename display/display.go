package display

import (
	"fmt"
	"go3d/actors"
	"go3d/utils"
	"math"
	"slices"
	"strings"
	"time"
)

var ColorMap = map[string]string{
	"reset":   "\033[0m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"gray":    "\033[37m",
	"white":   "\033[97m",
}

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
	DepthBuffer [][]float64

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

	DrawVerts  bool
	DrawWire   bool
	RenderWire bool
	RenderFace bool

	CamMoveSpeed float64

	FrameStart time.Time
	FrameTime  time.Duration
	FrameEnd   time.Duration

	MaxFrameTime time.Duration

	Xborder   string
	Triangles []actors.Triangle
}

func CreateView(w uint16, h uint16, fps uint8, moveSpeed float64) *View {

	v := View{
		Xpx:          w,
		Ypx:          h,
		TargetFPS:    fps,
		Fov:          90,
		CamX:         0,
		CamY:         3,
		CamZ:         10,
		CamRot:       []float64{0, 0, 0},
		NearClip:     -1,
		FarClip:      -150,
		CamMoveSpeed: moveSpeed,

		DrawVerts:  false,
		DrawWire:   false,
		RenderWire: true,
		RenderFace: true,
	}

	// Calc max frame time
	v.MaxFrameTime = v.CalcMaxFrameTime()

	// Initialize buffers
	v.FrameBuffer = utils.InitFrameBuffer(v.Xpx, v.Ypx)
	v.DepthBuffer = utils.InitDepthBuffer(v.Xpx, v.Ypx)

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
func (v *View) RegisterTriangle(t actors.Triangle) {
	v.Triangles = append(v.Triangles, t)
}

func (v *View) RegisterObject(o actors.Object) {
	for _, tri := range o.Tris {
		v.RegisterTriangle(tri)
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

// Set the FrameBuffer to empty pixels and depth buffer to max depth
func (v *View) ClearBuffer() {
	for i := range v.FrameBuffer {
		for j := range v.FrameBuffer[i] {
			v.FrameBuffer[i][j] = pixel[0]
		}
	}
	for i := range v.DepthBuffer {
		for j := range v.DepthBuffer[i] {
			v.DepthBuffer[i][j] = math.MaxFloat32
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

	for _, a := range v.Triangles {
		parent := a.ObjRef

		//Save raster verts for connecting with lines & filling face after vertex pass
		var rasterVerts [][]uint16

		// Save lines drawn for filling in faces
		var rasterLines [][]uint16

		// Store depth values for an approximated zbuffer
		var depthVals []float64

		// Calculate vertecies
		for _, vert := range a.Verts() {

			// Convert to worldspace
			worldVert := utils.ApplyWorldMatrix(vert, parent.ObjX, parent.ObjY, parent.ObjZ, parent.Scale, parent.Rot)

			camSpaceVert := utils.ApplyCamMatrix(v.CamX, v.CamY, v.CamZ, v.CamRot, worldVert[0], worldVert[1], worldVert[2])

			// fmt.Printf("Camspace Vert: %v\n", camSpaceVert)
			clipVert := utils.ApplyProjectionMatrix(camSpaceVert, v.XProjConst, v.YProjConst, v.ZProjConst, v.WProjConst)

			// Discard if W out of bounds
			depthVals = append(depthVals, clipVert[3])

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
			if v.DrawVerts {
				v.TouchBuffer(ColorMap["red"], xVert, yVert)
			}

		}

		// Draw lines between 2D verts with bresenhams alg
		if v.RenderWire {

			if len(rasterVerts) > 1 {

				// Keep track of connected points
				connected := make(map[int]map[int]bool)

				// Iterate through each vertex, connecting with neighbors and skipping if the reverse has been done
				for i := range len(rasterVerts) {
					for j := range len(rasterVerts) {
						if i != j {

							// Check if i has been connected to anything
							if _, ok := connected[i]; ok {

								// Check if i has been connected to j, then skip
								if _, ok := connected[i][j]; ok {
									continue
								}
							}

							drawn := v.DrawLine(rasterVerts[i], rasterVerts[j])
							rasterLines = append(rasterLines, drawn...)

							// Record the edge as drawn
							if _, ok := connected[i]; !ok {
								connected[i] = make(map[int]bool)
							}
							connected[i][j] = true

						}

					}

				}
			}
		}

		// Fill in faces via scanlines
		if v.RenderFace {
			// Calculate the min/max X and Y in triangle verts for bounding box

			// Calculate average depth for the face
			var depth float64
			for _, w := range depthVals {
				depth += w
			}
			depth = depth / float64(len(depthVals))

			var maxX, maxY uint16
			var minX, minY uint16 = math.MaxUint16, math.MaxUint16

			for _, vert := range rasterVerts {

				// Handle X
				if vert[0] < minX {
					minX = vert[0]
				}
				if vert[0] > maxX {
					maxX = vert[0]
				}
				// Handle y
				if vert[1] < minY {
					minY = vert[1]
				}
				if vert[1] > maxY {
					maxY = vert[1]
				}
			}

			// Reshape lines slice to be more useful here
			var allPoints [][]uint16 = rasterLines
			allPoints = append(allPoints, rasterVerts...)

			linePoints := make(map[uint16][]uint16)

			// Map what x coordinates have been drawn with a given Y
			for _, v := range allPoints {
				x := v[0]
				y := v[1]

				linePoints[y] = append(linePoints[y], x)
			}

			// Offsets for skipping pixels if wireframe is drawn
			var lineOffsetLeft, lineOffsetRight uint16
			lineOffsetRight = 1
			if v.DrawWire {
				lineOffsetLeft = 1
				lineOffsetRight = 0
			}
			// Within the bounding box, find the left and right raster bounds of triangle based on drawn lines
			for y := minY; y < maxY; y++ {
				if len(linePoints[y]) > 1 {

					leftBound := slices.Min(linePoints[y])
					rightBound := slices.Max(linePoints[y])

					// Draw in the pixels inbetween these
					for x := leftBound + lineOffsetLeft; x < rightBound+lineOffsetRight; x++ {

						// Only draw if the pixel is infront of other faces, based on average face depth
						if v.DepthBuffer[y][x] > depth {
							v.TouchBuffer(ColorMap[parent.Color], x, y)
							v.DepthBuffer[y][x] = depth
						}
					}
				}

			}

		}

	}

}

// Adds contiguous line to framebuffer between 2 points w/ Bresenhams alg
func (v *View) DrawLine(start []uint16, end []uint16) [][]uint16 {

	startX := start[0]
	startY := start[1]

	endX := end[0]
	endY := end[1]

	// Will skip vertex pixel if drawing verts is enabled
	var vertSkip uint16 = 0
	if v.DrawVerts {
		vertSkip = 1
	}
	// Pixels that will be drawn to buffer after calculations
	var pixels [][]uint16

	// Case of verical line
	if startX == endX {
		yMin := min(startY, endY) + vertSkip
		yMax := max(startY, endY)

		// Save pixels to draw along vertical line
		for y := yMin; y < yMax; y++ {
			pixels = append(pixels, []uint16{startX, y})
		}
	} else if startY == endY {
		// Case of flat line
		xMin := min(startX, endX) + vertSkip
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
			yMin := min(startY, endY) + vertSkip
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
	if v.DrawWire {

		for _, p := range pixels {
			v.TouchBuffer(ColorMap["cyan"], p[0], p[1])
		}
	}

	return pixels

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
