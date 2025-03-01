package display

import (
	"fmt"
	"go3d/actors"
	"go3d/utils"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// View defines the screenspace printed to the terminal, with any debug
// information
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

	RenderWire    bool
	OverlayOrigin []uint16

	CamMoveSpeed float64

	FrameStart time.Time
	FrameTime  time.Duration
	FrameEnd   time.Duration

	MaxFrameTime time.Duration
	FrameCount   uint64
	PrevFt       float64

	Xborder     string
	Triangles   []*actors.Triangle
	PointLights []*actors.Light
}

// Create a view filling the whole terminal, and accept some custom options
func CreateView(fps uint8, moveSpeed float64) *View {

	// Get the terminal width and height
	var ws struct {
		Rows uint16
		Cols uint16
		X    uint16
		Y    uint16
	}
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, uintptr(0), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&ws)))

	v := View{
		Xpx:          ws.Cols / 2,
		Ypx:          ws.Rows - 1,
		TargetFPS:    fps,
		Fov:          90,
		CamX:         0,
		CamY:         3,
		CamZ:         10,
		CamRot:       []float64{0, 0, 0},
		NearClip:     -1,
		FarClip:      -50,
		CamMoveSpeed: moveSpeed,

		RenderWire: true,
	}

	// Calc max frame time
	v.MaxFrameTime = v.CalcMaxFrameTime()

	// Origin for drawing big text
	v.OverlayOrigin = []uint16{5, 5}
	// Initialize buffers
	v.FrameBuffer, v.DepthBuffer = utils.CreateBuffers(v.Xpx, v.Ypx)

	// Calculate projection constants
	v.CalcProjectionConstants()

	// Initialize screen border
	v.Xborder = utils.ColorMap["Blue"][4] + strings.Repeat(pixel[1], int(v.Xpx)+2) + "\033[0m\n"

	// Remove cursor
	fmt.Print("\033[?25l")

	v.ClearBuffer()

	return &v
}
