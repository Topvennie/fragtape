package image

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // Support jpegs
	_ "image/png"  // Support png's

	"github.com/chai2010/webp"
)

func ToWebp(imageBytes []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("decode image %w", err)
	}

	var out bytes.Buffer
	if err := webp.Encode(&out, img, &webp.Options{Quality: 90}); err != nil {
		return nil, fmt.Errorf("encode to webp %w", err)
	}

	return out.Bytes(), nil
}
