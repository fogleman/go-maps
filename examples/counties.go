package main

import (
	"runtime"

	"github.com/fogleman/go-maps"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/ungerik/go-cairo"
)

const (
	m      = 1
	width  = m * 2048
	height = m * 1024
	lat    = 35.3
	lng    = -79.8
	scale  = m * 12000
)

func HexColor(x string) colorful.Color {
	color, _ := colorful.Hex(x)
	return color
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	bg := HexColor("#222028")
	land := HexColor("#FFFFFF")

	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, width, height)
	dc := maps.NewCanvas(surface, lat, lng, scale, 0)

	dc.SetLineCap(cairo.LINE_CAP_ROUND)
	dc.SetLineJoin(cairo.LINE_JOIN_ROUND)

	dc.SetSourceRGB(bg.R, bg.G, bg.B)
	dc.Paint()

	dc.SelectFontFace("Helvetica Neue", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	dc.SetFontSize(12)

	shapes, _ := maps.LoadSHP("files/cb_2013_us_county_5m/cb_2013_us_county_5m.shp")
	for _, shape := range shapes {
		tags := shape.Tags
		if tags["STATEFP"] != "37" {
			continue
		}
		name := tags["NAME"]
		dc.DrawShape(shape)
		dc.SetSourceRGB(land.R, land.G, land.B)
		dc.Stroke()
		pt := shape.Centroid()
		x, y := maps.Mercator(pt.Y, pt.X, scale)
		e := dc.TextExtents(name)
		dc.MoveTo(x-e.Width/2, y+e.Height/2)
		dc.ShowText(name)
	}

	dc.WriteToPNG("output.png")
}
