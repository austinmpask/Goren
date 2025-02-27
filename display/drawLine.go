package display

import (
	"go3d/utils"
	"math"
)

// Adds contiguous line to framebuffer between 2 points w/ Bresenhams alg This
// will always be calculated if faces are rendered, as these lines preface the
// face area calculation
func (v *View) DrawLine(start []uint16, end []uint16) [][]uint16 {

	startX := start[0]
	startY := start[1]

	endX := end[0]
	endY := end[1]

	// Will skip vertex pixel if drawing verts is enabled
	var vertSkip uint16 = 0
	if utils.DrawVerts {
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

		// Iterate Y for slope 1 or higher as each Y coordinate will have only
		// one pixel
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
			// Iterate over X for slope < 1 as each X coordinate will have only
			// one pixel

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
	if utils.DrawWire {

		for _, p := range pixels {
			v.FrameBuffer[p[1]][p[0]] = utils.ColorMap["Cyan"][5]
		}
	}

	return pixels

}
