package actors

type Object struct {
	Tris  []Triangle
	ObjX  float64
	ObjY  float64
	ObjZ  float64
	Rot   []float64
	Color string

	Scale float64
}

// Constructor that associates triangles with the object
func CreateObject(triangles [][][]float64, objX float64, objY float64, objZ float64, scale float64, color string) *Object {

	o := Object{
		Tris:  []Triangle{},
		Scale: scale,
		ObjX:  -objX,
		ObjY:  -objY,
		ObjZ:  objZ,
		Rot:   []float64{0, 0, 0},
		Color: color,
	}
	// Make triangles out of the coordinates
	for _, t := range triangles {
		tri := CreateTriangle(t[0], t[1], t[2], &o)
		o.Tris = append(o.Tris, tri)

	}

	return &o
}

func (o *Object) Translate(dx float64, dy float64, dz float64) {
	o.ObjX -= dx
	o.ObjY -= dy
	o.ObjZ -= dz
}

func (o *Object) Rotate(rx float64, ry float64, rz float64) {
	o.Rot[0] += rx
	o.Rot[1] += ry
	o.Rot[2] += rz
}
