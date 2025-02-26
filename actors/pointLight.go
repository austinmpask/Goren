package actors

type PointLight struct {
	LightX float64
	LightY float64
	LightZ float64

	Falloff   float64 //How far the light reaches objects
	Intensity float64 //Maximum effect of the light 0-1
}

func CreatePointLight(x float64, y float64, z float64, intensity float64, falloff float64) *PointLight {

	l := PointLight{
		LightX:    x,
		LightY:    y,
		LightZ:    z,
		Intensity: intensity,
		Falloff:   falloff,
	}

	return &l
}

func (l *PointLight) Translate(dx float64, dy float64, dz float64) {
	l.LightX += dx
	l.LightY -= dy
	l.LightZ += dz
}
