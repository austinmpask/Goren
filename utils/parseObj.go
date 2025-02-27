package utils

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Parse .obj file to convert to Object actor which is made up of Triangles
func ParseObj(path string) [][][]float64 {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	// Will store list of vertecies which will be mapped to faces
	verts := [][]float64{}
	tris := [][][]float64{}

	scanner := bufio.NewScanner(file)

	// Check each line, adda vertex if appropriate or add a triangle
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		// Skip blanks
		if len(fields) == 0 {
			continue
		}

		// Vertex line
		if fields[0] == "v" && len(fields) >= 4 {
			// Parse the X, Y, Z coordinates
			x, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil
			}

			y, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				return nil
			}

			z, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				return nil
			}

			// Add to slice of vertecies
			verts = append(verts, []float64{x, y, z})
		}

		// Face lines
		if fields[0] == "f" {

			// Works for 3-n verts, will subdivide face into triangles
			if len(fields) >= 4 {
				// Collect all verts
				vs := []int{}
				for i := 1; i < len(fields); i++ {
					vertexStr := strings.Split(fields[i], "/")[0]
					vertexInt, _ := strconv.Atoi(vertexStr)
					vertexInt-- //0 index
					vs = append(vs, vertexInt)
				}

				// Create triangles - .obj face traces face counter clockwise
				// First vert will be origin for all triangles, create triangles going around the face
				for i := 1; i < len(vs)-1; i++ {
					triangle := [][]float64{verts[vs[0]], verts[vs[i]], verts[vs[i+1]]}
					tris = append(tris, triangle)

				}

			}
		}
	}
	return tris
}
