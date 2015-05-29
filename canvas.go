package maps

import (
	"github.com/qedus/osmpbf"
	"github.com/ungerik/go-cairo"
)

type Canvas struct {
	cairo.Surface
	Lat      float64
	Lng      float64
	Scale    float64
	Rotation float64
}

func NewCanvas(surface *cairo.Surface, lat, lng, scale, rotation float64) *Canvas {
	canvas := Canvas{*surface, lat, lng, scale, rotation}
	w := float64(surface.GetWidth())
	h := float64(surface.GetHeight())
	x, y := Mercator(lat, lng, scale)
	surface.IdentityMatrix()
	surface.Translate(w/2, h/2)
	surface.Rotate(rotation)
	surface.Translate(-x, -y)
	return &canvas
}

func (canvas *Canvas) DrawWay(pbf *PBF, way *osmpbf.Way) {
	canvas.NewSubPath()
	for _, id := range way.NodeIDs {
		node := pbf.Nodes[id]
		x, y := Mercator(node.Lat, node.Lon, canvas.Scale)
		canvas.LineTo(x, y)
	}
}

func (canvas *Canvas) DrawMultiPolygon(pbf *PBF, relation *osmpbf.Relation) {
	for _, member := range relation.Members {
		if member.Type != osmpbf.WayType {
			continue
		}
		if way, ok := pbf.Ways[member.ID]; ok {
			canvas.DrawWay(pbf, way)
		}
	}
}
