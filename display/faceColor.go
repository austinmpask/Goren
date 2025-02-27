package display

import (
	"go3d/utils"
	"math"
)

// Returns a value between 1-10 referring to a color shade based on scene
// lighting and camera depth
func (v *View) CalculateFaceColor(depth float64, center []float64, falloff float64) int {
	var baseIntensity = 1

	if utils.RenderLighting {
		// 1 is the minimum luminance for a given color

		// Apply scene lighting
		for _, light := range v.PointLights {
			// Calculate worldspace distance from light to face
			d := math.Sqrt(math.Pow(light.LightX-center[0], 2) + math.Pow(light.LightY-center[1], 2) + math.Pow(light.LightZ-center[2], 2))

			// Check if face is within the lights falloff
			if d <= light.Falloff {
				lightFactor := int(5 - math.Round((5/light.Intensity)*(d/light.Falloff)))

				// Bound the light between 1-5
				lightFactor = max(0, lightFactor)
				lightFactor = min(5, lightFactor)
				baseIntensity += lightFactor
			}
		}

		// Adjust based on camera depth

		// Reduce luminance a bit for extremely far objects Apply falloff to
		// depth range, this is where the effect will be applied
		depthRange := falloff * (math.Abs(v.FarClip) - math.Abs(v.NearClip))
		minDepth := math.Abs(v.FarClip) - depthRange //Minimum depth at which depth will be applied

		maxEffect := 2

		if depth >= minDepth {
			dNorm := depth - minDepth
			maxNorm := math.Abs(v.FarClip) - minDepth
			baseIntensity -= int(math.Round(float64(maxEffect) * (dNorm / maxNorm)))
		}

		// Bound final luminance to colorspace
		baseIntensity = min(baseIntensity, 10)
		baseIntensity = max(baseIntensity, 1)

	}
	return baseIntensity

}
