package main

import (
	"log"
	"math/rand"
	"runtime"

	"github.com/fogleman/go-maps"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/ungerik/go-cairo"
)

// http://download.geofabrik.de/north-america/us/district-of-columbia-latest.osm.pbf

const (
	m      = 4
	width  = m * 1024
	height = m * 1024

	// path   = "files/district-of-columbia-latest.osm.pbf"
	// lat    = 38.9047
	// lng    = -77.0164
	// scale  = m * 220000

	// path  = "files/nyc.osm.pbf"
	// lat   = 40.7903
	// lng   = -73.9597
	// scale = m * 220000

	// path  = "files/chicago.osm.pbf"
	// lat   = 41.87
	// lng   = -87.665
	// scale = m * 400000

	path  = "files/sfo.osm.pbf"
	lat   = 37.7507
	lng   = -122.4367
	scale = m * 300000
)

func HexColor(x string) colorful.Color {
	color, _ := colorful.Hex(x)
	return color
}

func DrawBuildings(dc *maps.Canvas, colors []colorful.Color, path string) {
	pbf, err := maps.LoadPBF(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, way := range pbf.Ways {
		if _, ok := way.Tags["building"]; ok {
			c := colors[rand.Intn(len(colors))]
			dc.DrawWay(pbf, way)
			dc.SetSourceRGB(c.R, c.G, c.B)
			dc.Fill()
		}
	}

	for _, relation := range pbf.Relations {
		if relation.Tags["type"] != "multipolygon" {
			continue
		}
		if _, ok := relation.Tags["building"]; ok {
			c := colors[rand.Intn(len(colors))]
			dc.DrawMultiPolygon(pbf, relation)
			dc.SetSourceRGB(c.R, c.G, c.B)
			dc.Fill()
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	bg := HexColor("#222028")
	colors := []colorful.Color{
		HexColor("#730046"),
		HexColor("#BFBB11"),
		HexColor("#FFC200"),
		HexColor("#E88801"),
		HexColor("#C93C00"),
	}

	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, width, height)
	dc := maps.NewCanvas(surface, lat, lng, scale, 0)

	dc.SetFillRule(cairo.FILL_RULE_EVEN_ODD)
	dc.SetLineCap(cairo.LINE_CAP_ROUND)
	dc.SetLineJoin(cairo.LINE_JOIN_ROUND)
	dc.SetSourceRGB(bg.R, bg.G, bg.B)
	dc.Paint()

	DrawBuildings(dc, colors, path)

	dc.WriteToPNG("output.png")
}
