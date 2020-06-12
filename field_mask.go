package hrpc

import (
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
	"google.golang.org/genproto/protobuf/field_mask"
)

// MaskMap mask all keys of a map with provided depth.
func MaskMapKeys(inp hexa.Map, mask *field_mask.FieldMask, depth int) {
	extractor := gutil.MapKeysExtractor{Depth: depth, Separator: "."}
	keys := extractor.Extract(inp)
	mask.Paths = keys
	return
}
