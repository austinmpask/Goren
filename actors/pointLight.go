package actors

type Light struct {
	LightX float64
	LightY float64
	LightZ float64

	Falloff   float64 //How far the light reaches objects
	Intensity float64 //Maximum effect of the light 0-1
}

func (l *Light) Translate(dx float64, dy float64, dz float64) {
	l.LightX += dx
	l.LightY -= dy
	l.LightZ += dz
}
