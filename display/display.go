package display

import (
	"fmt"
	"go3d/utils"
	"strings"
	"time"
)

// Monochrome colorspace
var pixel = map[uint8]string{
	0: "  ",
	1: "██",
}

// View defines the screenspace printed to the terminal, with any debug information
type View struct {
	Xpx         uint8
	Ypx         uint8
	TargetFPS   uint8
	FrameBuffer [][]string

	FrameStart time.Time
	FrameTime  time.Duration
	FrameEnd   time.Duration

	MaxFrameTime time.Duration

	Xborder string
}

func DefaultView() *View {

	v := View{
		Xpx:       80,
		Ypx:       40,
		TargetFPS: 144,
	}

	// Calc max frame time
	v.MaxFrameTime = v.CalcMaxFrameTime()

	// Initialize buffer
	v.FrameBuffer = make([][]string, v.Ypx)
	for i := range v.FrameBuffer {
		v.FrameBuffer[i] = make([]string, v.Xpx)
	}

	// Initialize screen border
	v.Xborder = strings.Repeat(pixel[1], int(v.Xpx)+2)

	// Remove cursor
	fmt.Print("\033[?25l")

	v.ClearBuffer()

	return &v
}

// Calculate the maximum allowable frametime to maintain the target framerate in MS
func (v *View) CalcMaxFrameTime() time.Duration {

	fps := time.Duration(v.TargetFPS)
	return time.Second / fps
}

func (v *View) Aspect() float32 {
	return float32(v.Xpx) / float32(v.Ypx)
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

// Translate worldspace and cameraspace data into screenspace pixels
func (v *View) PrepBuffer() {
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

func (v *View) DrawDebug() {
	ftMs := utils.InMs(v.FrameTime)
	maxFtMs := utils.InMs(v.MaxFrameTime)
	frEndMs := utils.InMs((v.FrameEnd))

	var util float64
	var pfps float64
	var fps float64

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
func (v *View) TouchBuffer() {

}

func (v *View) StartFrame() {
	v.FrameStart = time.Now()

}

func (v *View) EndFrame() {

	v.FrameTime = time.Since(v.FrameStart)
}

func (v *View) FrameSync() {
	frameTimeSlack := v.MaxFrameTime - v.FrameTime

	targetTime := time.Now().Add(frameTimeSlack)

	// Use for loop instead, more accurate scheduling than sleep
	for time.Now().Before(targetTime) {
		//Wait
	}

	v.FrameEnd = time.Since(v.FrameStart)
}
