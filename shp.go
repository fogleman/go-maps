package maps

import (
	"strings"

	"github.com/jonas-p/go-shp"
)

type Shape struct {
	shp.Shape
	Tags map[string]string
}

func LoadSHP(path string) ([]Shape, error) {
	file, err := shp.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fields := file.Fields()
	names := make([]string, len(fields))
	for i, field := range fields {
		names[i] = strings.Trim(field.String(), "\x00")
	}

	var result []Shape
	for file.Next() {
		n, shape := file.Shape()
		tags := make(map[string]string)
		for i, name := range names {
			value := file.ReadAttribute(n, i)
			tags[name] = value
		}
		result = append(result, Shape{shape, tags})
	}
	return result, nil
}
