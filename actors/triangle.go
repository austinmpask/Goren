package actors

type Triangle struct {
	A []float64
	B []float64
	C []float64
}

// Basic triangle for rendering which inherits the Actor interface
func CreateTriangle(a []float64, b []float64, c []float64) Actor {
	t := Triangle{
		A: a,
		B: b,
		C: c,
	}

	return &t
}

func (t *Triangle) Verts() [][]float64 {
	return [][]float64{t.A, t.B, t.C}
}
