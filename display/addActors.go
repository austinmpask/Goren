package display

import "go3d/actors"

// Add actors to the scene
func (v *View) RegisterTriangle(t *actors.Triangle) {
	v.Triangles = append(v.Triangles, t)
}

func (v *View) RegisterObject(o *actors.Object) {
	// Objects are converted to triangles in the scene, which can be associated
	// to a parent object for transformations
	for _, tri := range o.Tris {
		v.RegisterTriangle(&tri)
	}
}
func (v *View) RegisterLight(l *actors.Light) {
	v.PointLights = append(v.PointLights, l)
}
