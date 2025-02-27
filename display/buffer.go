package display

import (
	"fmt"
	"math"
	"strings"
)

var pixel = map[uint8]string{
	0: "  ",
	1: "██",
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

}

// Prints buffer contents to screen
func (v *View) DrawBuffer() {

	var sb strings.Builder

	sb.WriteString(v.Xborder)
	for _, row := range v.FrameBuffer {
		for _, pxl := range row {
			sb.WriteString(pxl)
		}
		sb.WriteByte('\n')
	}
	sb.WriteString(v.Xborder)

	s := sb.String()
	fmt.Print(s)

}
