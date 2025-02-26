package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Parse .obj file to convert to Object actor which is made up of Triangles
func ParseObj(path string) ([][][]float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Will store list of vertecies which will be mapped to faces
	verts := [][]float64{}
	tris := [][][]float64{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) == 0 {
			continue
		}

		// Vertex line
		if fields[0] == "v" && len(fields) >= 4 {
			// Parse the X, Y, Z coordinates
			x, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid X coordinate: %v", err)
			}

			y, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid Y coordinate: %v", err)
			}

			z, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid Z coordinate: %v", err)
			}

			// Add to slice of vertecies
			verts = append(verts, []float64{x, y, z})
		}

		// Face lines
		if fields[0] == "f" {

			// Case of triangle face. Ignoring others, so some detail will be lost
			if len(fields) == 5 {

				// Get vertex indecies
				v1Str := strings.Split(fields[1], "/")[0]
				v2Str := strings.Split(fields[2], "/")[0]
				v3Str := strings.Split(fields[3], "/")[0]
				v1, err1 := strconv.Atoi(v1Str)
				v2, err2 := strconv.Atoi(v2Str)
				v3, err3 := strconv.Atoi(v3Str)
				if err1 != nil || err2 != nil || err3 != nil {
					return nil, fmt.Errorf("invalid vertex")
				}
				v1-- //Convert to 0 index
				v2--
				v3--

				// Add a triangle of those vertecies
				triangle := [][]float64{verts[v1], verts[v2], verts[v3]}
				tris = append(tris, triangle)

			}
		}
	}
	return tris, nil
}
