package display

import "time"

// Log the time the frame calculations began
func (v *View) StartFrame() {
	v.FrameStart = time.Now()

}

// Log the time once the buffer and anything else was drawn to screen
func (v *View) EndFrame() {

	v.FrameTime = time.Since(v.FrameStart)
}

// Minimize screen tearing by waiting until frametime for the target framerate
// has elapsed before continuing
func (v *View) FrameSync(method string, adj int) {

	// Calc remaining time for frame
	frameTimeSlack := v.MaxFrameTime - v.FrameTime

	targetTime := time.Now().Add(frameTimeSlack).Add(time.Duration(adj) * time.Microsecond)

	switch method {
	// Primary frame sync method, introduces some drift
	case "sleep":
		time.Sleep(time.Until(targetTime))
	// Alternate method, more accurate framerates but high cpu
	case "loop":
		for time.Now().Before(targetTime) {
			//Wait
		}

	}
	// Log the time that the entire frame ended
	v.FrameEnd = time.Since(v.FrameStart)
}
