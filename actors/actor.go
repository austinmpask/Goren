package actors

// Interface for anything rendered in the scene (object, plane, line, vert)
type Actor interface {
	Verts() [][]float64
}
