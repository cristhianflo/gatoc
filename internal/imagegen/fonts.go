package imagegen

import (
	"fmt"
	"sync"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	_ "embed"
)

//go:embed assets/Inter-Regular.ttf
var interRegularTTF []byte

//go:embed assets/Inter-Bold.ttf
var interBoldTTF []byte

var (
	fontsOnce sync.Once
	regular   *truetype.Font
	bold      *truetype.Font
	fontsErr  error
)

// loadFonts parses the embedded Inter fonts once. truetype.Font values are
// read-only after parsing and safe to share; faces created with newFace are
// not thread safe and must only be used by a single render call.
func loadFonts() (*truetype.Font, *truetype.Font, error) {
	fontsOnce.Do(func() {
		regular, fontsErr = truetype.Parse(interRegularTTF)
		if fontsErr != nil {
			return
		}
		bold, fontsErr = truetype.Parse(interBoldTTF)
	})
	if fontsErr != nil {
		return nil, nil, fmt.Errorf("parse embedded fonts: %w", fontsErr)
	}
	return regular, bold, nil
}

func newFace(f *truetype.Font, size float64) font.Face {
	return truetype.NewFace(f, &truetype.Options{Size: size, DPI: 72})
}
