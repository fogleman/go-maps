package main

import (
	"fmt"
	"log"
	"math"
	"runtime"

	"github.com/fogleman/go-maps"
)

const (
	Path = "files/cb_2013_us_state_5m/cb_2013_us_state_5m.shp"
	// Path      = "files/cb_2013_us_county_5m/cb_2013_us_county_5m.shp"
	StateCode = "37"
	FeedRate  = 40
	PullUp    = 0.2
	PushDown  = -0.04
	OffsetX   = 0.25
	OffsetY   = 0.25
	SizeX     = 6.5
	SizeY     = 4.5
)

type Point struct {
	X, Y float64
}

func (a Point) Min(b Point) Point {
	return Point{math.Min(a.X, b.X), math.Min(a.Y, b.Y)}
}

func (a Point) Max(b Point) Point {
	return Point{math.Max(a.X, b.X), math.Max(a.Y, b.Y)}
}

func (a Point) Add(b Point) Point {
	return Point{a.X + b.X, a.Y + b.Y}
}

func (a Point) Sub(b Point) Point {
	return Point{a.X - b.X, a.Y - b.Y}
}

func (a Point) Mul(b Point) Point {
	return Point{a.X * b.X, a.Y * b.Y}
}

func (a Point) MulScalar(b float64) Point {
	return Point{a.X * b, a.Y * b}
}

func (a Point) Div(b Point) Point {
	return Point{a.X / b.X, a.Y / b.Y}
}

func (a Point) MinComponent() float64 {
	return math.Min(a.X, a.Y)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	shapes, err := maps.LoadSHP(Path)
	if err != nil {
		log.Fatal(err)
	}

	// compute dimensions
	var points []Point
	for _, shape := range shapes {
		if shape.Tags["STATEFP"] != StateCode {
			continue
		}
		groups := shape.GetPoints()
		for _, group := range groups {
			for _, pt := range group {
				x, y := maps.Mercator(pt.Y, pt.X, 1)
				y = -y
				point := Point{x, y}
				points = append(points, point)
			}
		}
	}

	lo := points[0]
	hi := points[0]
	for _, point := range points {
		lo = lo.Min(point)
		hi = hi.Max(point)
	}

	offset := Point{OffsetX, OffsetY}
	scale := Point{SizeX, SizeY}.Div(hi.Sub(lo)).MinComponent()

	// generate code
	fmt.Println("G90")             // use absolute
	fmt.Println("G20")             // use inches
	fmt.Printf("G0 Z%f\n", PullUp) // pull up
	fmt.Println("M4")              // turn spindle on
	fmt.Println("G4 P2.0")         // pause

	for _, shape := range shapes {
		if shape.Tags["STATEFP"] != StateCode {
			continue
		}
		groups := shape.GetPoints()
		for _, group := range groups {
			fmt.Printf("G0 Z%f\n", PullUp) // pull up
			for i, pt := range group {
				x, y := maps.Mercator(pt.Y, pt.X, 1)
				y = -y
				point := Point{x, y}.Sub(lo).MulScalar(scale).Add(offset)
				if i == 0 {
					fmt.Printf("G0 X%f Y%f\n", point.X, point.Y)   // go to point
					fmt.Printf("G1 Z%f F%d\n", PushDown, FeedRate) // push down
				} else {
					fmt.Printf("G1 X%f Y%f F%d\n", point.X, point.Y, FeedRate) // cut to point
				}
			}
		}
	}

	fmt.Printf("G0 Z%f\n", PullUp) // pull up
	fmt.Println("M8")              // turn spindle off
	// fmt.Println("G0 X0 Y0")        // go home
}
