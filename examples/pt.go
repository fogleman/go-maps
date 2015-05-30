package main

import (
	"log"
	"math"
	"strconv"

	"github.com/fogleman/go-maps"
	"github.com/fogleman/pt/pt"
	"github.com/qedus/osmpbf"
)

const radius = 6371000

func point(lat, lng, altitude float64) pt.Vector {
	cosLat := math.Cos(lat * math.Pi / 180)
	sinLat := math.Sin(lat * math.Pi / 180)
	cosLng := math.Cos(lng * math.Pi / 180)
	sinLng := math.Sin(lng * math.Pi / 180)
	r := radius + altitude
	x := r * cosLat * cosLng
	y := r * cosLat * sinLng
	z := r * sinLat
	return pt.Vector{x, y, z}
}

func main() {
	scene := pt.Scene{}
	material := pt.GlossyMaterial(pt.HexColor(0xE4DFD3), 1.2, pt.Radians(20))

	pbf, err := maps.LoadPBF("files/manhattan.osm.pbf")
	if err != nil {
		log.Fatal(err)
	}

	var triangles []*pt.Triangle
	for _, way := range pbf.Ways {
		if _, ok := way.Tags["building"]; ok {
			height, err := strconv.ParseFloat(way.Tags["height"], 64)
			if err != nil {
				height = 8 // mean=8.4, dev=5.7
			}
			for i := range way.NodeIDs {
				a := pbf.Nodes[way.NodeIDs[i]]
				var b *osmpbf.Node
				if i == 0 {
					b = pbf.Nodes[way.NodeIDs[len(way.NodeIDs)-1]]
				} else {
					b = pbf.Nodes[way.NodeIDs[i-1]]
				}
				if a == nil || b == nil {
					continue
				}
				{
					v1 := point(a.Lat, a.Lon, 0)
					v2 := point(b.Lat, b.Lon, 0)
					v3 := point(b.Lat, b.Lon, height)
					t := pt.NewTriangle(v1, v2, v3, material)
					triangles = append(triangles, t)
				}
				{
					v1 := point(a.Lat, a.Lon, 0)
					v2 := point(b.Lat, b.Lon, height)
					v3 := point(a.Lat, a.Lon, height)
					t := pt.NewTriangle(v1, v2, v3, material)
					triangles = append(triangles, t)
				}
			}
		}
	}

	mesh := pt.NewMesh(triangles)
	scene.Add(mesh)

	earth := pt.GlossyMaterial(pt.HexColor(0xB1DAFF), 1.2, pt.Radians(20))
	shape := pt.NewSphere(pt.Vector{}, radius, earth)
	scene.Add(shape)

	look := point(40.735958, -74.001035, -1000)
	eye := point(40.687651, -74.001850, 800)
	up := eye.Normalize()
	camera := pt.LookAt(eye, look, up, 40)

	pt.IterativeRender("out%03d.png", 1000, &scene, &camera, 2048, 1024, -1, 4, 4)
}
