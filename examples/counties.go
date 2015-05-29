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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	bg, _ := colorful.Hex("#A7C9AE")
	land, _ := colorful.Hex("#FFE7AD")
	stroke, _ := colorful.Hex("#FFAB48")
	text, _ := colorful.Hex("#CC6B32")

	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, width, height)
	dc := maps.NewCanvas(surface, lat, lng, scale, 0)

	dc.SetLineCap(cairo.LINE_CAP_ROUND)
	dc.SetLineJoin(cairo.LINE_JOIN_ROUND)

	dc.SetSourceRGB(bg.R, bg.G, bg.B)
	dc.Paint()

	dc.SelectFontFace("Helvetica Neue", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
	dc.SetFontSize(12)

	shapes, _ := maps.LoadSHP("files/cb_2013_us_county_5m/cb_2013_us_county_5m.shp")

	dc.SetSourceRGB(land.R, land.G, land.B)
	for _, shape := range shapes {
		if shape.Tags["STATEFP"] != "37" {
			continue
		}
		dc.DrawShape(shape)
		dc.Fill()
	}

	for _, shape := range shapes {
		if shape.Tags["STATEFP"] != "37" {
			continue
		}
		name := shape.Tags["NAME"]
		dc.SetSourceRGB(stroke.R, stroke.G, stroke.B)
		dc.DrawShape(shape)
		dc.Stroke()
		pt := shape.Centroid()
		x, y := maps.Mercator(pt.Y, pt.X, scale)
		e := dc.TextExtents(name)
		dc.SetSourceRGB(text.R, text.G, text.B)
		dc.MoveTo(x-e.Width/2, y+e.Height/2)
		dc.ShowText(name)
	}

	dc.WriteToPNG("output.png")
}
