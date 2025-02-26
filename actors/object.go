package actors

type Object struct {
	Tris []Actor
	ObjX float64
	ObjY float64
	ObjZ float64
}

func CreateObject(triangles [][][]float64) Object {

	o := Object{
		Tris: []Actor{},
	}
	// Make triangles out of the coordinates
	for _, t := range triangles {
		tri := CreateTriangle(t[0], t[1], t[2])
		o.Tris = append(o.Tris, tri)

	}

	return o
}
