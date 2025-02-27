package display

import (
	"go3d/utils"
	"math"
	"slices"
)

// Apply world transformations -> camera transformations -> projection transformations -> NDC transformations -> screenspace transformations
// Add results to the framebuffer (verts, lines, faces, lighting)
// Most of the meat and potatoes for rendering
func (v *View) PrepBuffer() {

triangleLoop:
	for _, a := range v.Triangles {
		// Save parent for color assignment
		parent := a.ObjRef

		//Save raster verts for connecting with lines & filling face after vertex pass
		// Save lines drawn for filling in faces
		var rasterVerts, rasterLines [][]uint16

		// Store depth values for an approximated zbuffer
		// Store worldspace verts for lighting calculations
		var depthVals []float64
		var worldVerts [][]float64
		// Calculate vertecies
		for _, vert := range a.Verts {

			// Convert to worldspace
			vert = utils.ApplyWorldMatrix(vert, parent.ObjX, parent.ObjY, parent.ObjZ, parent.Scale, parent.Rot)

			// Store world verts for lighting calculations
			worldVerts = append(worldVerts, vert)

			vert = utils.ApplyCamMatrix(v.CamX, v.CamY, v.CamZ, v.CamRot, vert[0], vert[1], vert[2])
			vert = utils.ApplyProjectionMatrix(vert, v.XProjConst, v.YProjConst, v.ZProjConst, v.WProjConst)

			// Prevent behind cam objects from drawing
			if vert[3] > math.Abs(v.FarClip) || vert[3] < math.Abs(v.NearClip) {
				continue triangleLoop
			}
			// Save depth vals for face rendering
			depthVals = append(depthVals, vert[3])

			vert = utils.ApplyNdcMatrix(vert)
			// Discard if out of bounds
			for _, v := range vert {
				if v > 1 || v < -1 {
					continue triangleLoop
				}

			}

			vert = utils.NdcToScreen(vert, v.Xpx, v.Ypx)
			// fmt.Printf("Screenspace Vert: %v\n", ssVert)

			xVert := uint16(math.Round(vert[0]))
			yVert := uint16(math.Round(vert[1]))

			// Save final 2D vertex for drawing lines
			rasterVerts = append(rasterVerts, []uint16{xVert, yVert})

		}
		// Load vertecies to buffer
		if utils.DrawVerts {
			go func() {
				for _, vertex := range rasterVerts {

					v.FrameBuffer[vertex[1]][vertex[0]] = utils.ColorMap["Red"][5]
				}
			}()
		}

		// Draw lines between 2D verts with bresenhams alg
		if v.RenderWire {

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

		// Fill in faces via scanlines
		if utils.RenderFace {

			// Calculate average depth for the face
			var depth float64
			for _, w := range depthVals {
				depth += w
			}
			depth /= float64(len(depthVals))

			// Calculate the min/max X and Y in triangle verts for bounding box
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
			if utils.DrawWire {
				lineOffsetLeft = 1
				lineOffsetRight = 0
			}

			// Calculate barycenter of face for lighting
			var xC, yC, zC float64

			for _, point := range worldVerts {
				xC += point[0]
				yC += point[1]
				zC += point[2]
			}
			xC /= float64(len(worldVerts))
			yC /= float64(len(worldVerts))
			zC /= float64(len(worldVerts))
			center := []float64{xC, yC, zC}

			// Calculate face color based on lighting and camera depth
			lum := v.CalculateFaceColor(depth, center, .3)

			// Within the bounding box, find the left and right raster bounds of triangle based on drawn lines
			for y := minY; y < maxY; y++ {
				if len(linePoints[y]) > 1 {

					leftBound := slices.Min(linePoints[y])
					rightBound := slices.Max(linePoints[y])

					// Draw in the pixels inbetween these
					for x := leftBound + lineOffsetLeft; x < rightBound+lineOffsetRight; x++ {
						// Only draw if the pixel is infront of other faces, based on average face depth
						if v.DepthBuffer[y][x] > depth {

							v.FrameBuffer[y][x] = utils.ColorMap[parent.Color][lum]
							v.DepthBuffer[y][x] = depth
						}
					}
				}

			}

		}

	}

	// Draw debug stats on the screen in big text
	if utils.Debug {
		v.FrameCount++
		v.DrawDebug()
	}

}
