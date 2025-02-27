package actors

type Triangle struct {
	Verts  [][]float64
	ObjRef *Object
}

// Basic triangle for rendering which inherits the Actor interface
func CreateTriangle(a []float64, b []float64, c []float64, obj *Object) Triangle {
	t := Triangle{
		Verts:  [][]float64{a, b, c},
		ObjRef: obj,
	}

	return t
}
