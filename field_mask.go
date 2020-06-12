package hrpc

import (
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
	"google.golang.org/genproto/protobuf/field_mask"
)

// MaskMap mask all keys of a map with provided depth.
func MaskMapPaths(m hexa.Map, mask *field_mask.FieldMask, depth int) {
	mask.Paths = gutil.MapPathExtractor{Depth: depth, Separator: "."}.Extract(m)
}
