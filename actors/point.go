package actors

type Point struct {
	X float64
	Y float64
	Z float64
}

func DefaultPoint() Actor {

	p := Point{
		X: 0,
		Y: 0,
		Z: -10,
	}
	return &p
}

func (p *Point) Verts() [][]float64 {
	return [][]float64{{p.X, p.Y, p.Z}}
}
