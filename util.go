package maps

import "math"

func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func Mercator(lat, lng, scale float64) (x, y float64) {
	x = Radians(lng) * scale
	y = math.Asinh(math.Tan(Radians(lat))) * -scale
	return
}
