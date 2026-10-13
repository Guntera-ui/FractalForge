package main

// First time writing keystrokes dont JUDGE bitte
func HandleKey(view Viewport, key byte) (Viewport, bool) {
	switch key {
	case 'w':
		distance := view.Width * 0.1
		return Pan(view, 0, distance), true
	case 's':
		distance := view.Width * (-0.1)
		return Pan(view, 0, distance), true
	case 'a':
		distance := view.Width * (-0.1)
		return Pan(view, distance, 0), true

	case 'd':
		distance := view.Width * (0.1)
		return Pan(view, distance, 0), true

	case '+':
		return Zoom(view, 2), true

	case '-':
		return Zoom(view, 0.5), true
	case 'q':
		return view, false

	default:
		return view, true
	}
}
