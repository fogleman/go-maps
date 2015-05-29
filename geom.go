package maps

type Point struct {
	X, Y float64
}

func Centroid(points []Point) Point {
	centroid := Point{}
	totalArea := 0.0
	for i, a := range points {
		var b Point
		if i == 0 {
			b = points[len(points)-1]
		} else {
			b = points[i-1]
		}
		area := a.X*b.Y - b.X*a.Y
		totalArea += area
		centroid.X += (a.X + b.X) * area
		centroid.Y += (a.Y + b.Y) * area
	}
	centroid.X /= totalArea * 3
	centroid.Y /= totalArea * 3
	return centroid
}
