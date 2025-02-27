package display

import (
	"fmt"
	"go3d/utils"
	"strings"
)

// Prints performance info as an overlay
func (v *View) DrawDebug() {

	ftMs := utils.InMs(v.FrameTime)
	maxFtMs := utils.InMs(v.MaxFrameTime)
	frEndMs := utils.InMs((v.FrameEnd))

	var util, pfps, fps float64

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

	// Calc averages, weighted towards newer values
	var aFt float64

	if v.FrameCount > 1 {
		alpha := 1.0 / float64(v.FrameCount)
		aFt = (1-alpha)*v.PrevFt + alpha*ftMs
	} else {
		aFt = ftMs
	}

	// Print directly to buffer
	c1 := "Red"
	c2 := "Cyan"

	v.DrawBigDebug(0, fmt.Sprintf("FT:      %.3fms", ftMs), c1)
	v.DrawBigDebug(1, fmt.Sprintf("FT UTIL: %.3f%%", util), c1)
	v.DrawBigDebug(2, fmt.Sprintf("P FPS:   %.3f", pfps), c1)
	v.DrawBigDebug(3, fmt.Sprintf("RL FPS:  %.3f", fps), c1)
	v.DrawBigDebug(4, fmt.Sprintf("POLYS:   %v", len(v.Triangles)), c1)
	v.DrawBigDebug(5, fmt.Sprintf("LIGHTS:  %v", len(v.PointLights)), c1)
	v.DrawBigDebug(7, fmt.Sprintf("FT AVG:     %.3fms", aFt), c2)
	v.DrawBigDebug(8, fmt.Sprintf("FT UTL AVG: %.3f%%", 100*aFt/maxFtMs), c2)
	v.DrawBigDebug(9, fmt.Sprintf("PT FPS AVG: %.3f", 1000/aFt), c2)

	// Save average for next frame for memory efficient avg frametime calc
	v.PrevFt = aFt

}

// Print debug information in 5x5 characters onto the framebuffer
func (v *View) DrawBigDebug(line uint16, text string, color string) {
	// Starting position from OverlayOrigin
	startX := v.OverlayOrigin[0]
	startY := v.OverlayOrigin[1] + 6*line

	// Loop through each character in the text
	for charIndex, char := range text {

		// Force uppercase incase i forgot
		charStr := strings.ToUpper(string(char))

		// Translate char to 5x5 pixel map
		bigChar, exists := utils.BigCharacters[charStr]
		if !exists {
			continue
		}

		// Calculate the starting position for character
		charStartX := startX + uint16(charIndex*6)

		// Build char from flattened array
		for row := 0; row < 5; row++ {
			for col := 0; col < 5; col++ {
				// Calculate the index in the flattened array
				index := row*5 + col

				// Calculate the screen position
				pixelX := charStartX + uint16(col)
				pixelY := startY + uint16(row)

				// Fill in pixel if appropriate
				if bigChar[index] == 1 {
					v.FrameBuffer[pixelY][pixelX] = utils.ColorMap[color][6]
				}
			}
		}
	}
}
