package algorithms

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"math/rand"
	"mime/multipart"
)

// GeneralTransformation applies the specified transformation on the image.
func GeneralTransformation(file multipart.File, transformationType string, c float64, gamma float64) (string, error) {
	img, err := decodeImage(file)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	transformedImg := image.NewRGBA(bounds)

	switch transformationType {
	case "Logarithmic":
		for y := 0; y < bounds.Dy(); y++ {
			for x := 0; x < bounds.Dx(); x++ {
				r, g, b, a := img.At(x, y).RGBA()
				rF := float64(r >> 8)
				gF := float64(g >> 8)
				bF := float64(b >> 8)

				newR := uint8(c * math.Log(1+rF))
				newG := uint8(c * math.Log(1+gF))
				newB := uint8(c * math.Log(1+bF))

				transformedImg.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
			}
		}

	case "Power Law":
		for y := 0; y < bounds.Dy(); y++ {
			for x := 0; x < bounds.Dx(); x++ {
				r, g, b, a := img.At(x, y).RGBA()
				rF := float64(r >> 8)
				gF := float64(g >> 8)
				bF := float64(b >> 8)

				newR := uint8(c * math.Pow(rF, gamma))
				newG := uint8(c * math.Pow(gF, gamma))
				newB := uint8(c * math.Pow(bF, gamma))

				transformedImg.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
			}
		}

	case "Random LUT":
		lut := make([]uint8, 256)
		for i := 0; i < 256; i++ {
			lut[i] = uint8(rand.Intn(256))
		}

		for y := 0; y < bounds.Dy(); y++ {
			for x := 0; x < bounds.Dx(); x++ {
				r, g, b, a := img.At(x, y).RGBA()
				rF := lut[r>>8]
				gF := lut[g>>8]
				bF := lut[b>>8]

				transformedImg.Set(x, y, color.RGBA{rF, gF, bF, uint8(a >> 8)})
			}
		}
	}

	var buf bytes.Buffer
	jpeg.Encode(&buf, transformedImg, nil)
	if err != nil {
		return "", err
	}

	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64Image, nil
}
