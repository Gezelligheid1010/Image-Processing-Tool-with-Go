package algorithms

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png" // 必须导入以支持 PNG 解码
	"math/rand"
	"mime/multipart"
)

// MixedAlgorithms processes various mixed algorithms on a single image
func MixedAlgorithms(file multipart.File, algorithm string, scalingFactor float64, shiftingValue float64, nBit int) (string, error) {
	img, err := decodeImage(file)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	resultImage := image.NewRGBA(bounds)

	switch algorithm {
	case "Negative":
		for y := 0; y < bounds.Dy(); y++ {
			for x := 0; x < bounds.Dx(); x++ {
				r, g, b, a := img.At(x, y).RGBA()
				resultImage.Set(x, y, color.RGBA{
					R: uint8(255 - (r >> 8)),
					G: uint8(255 - (g >> 8)),
					B: uint8(255 - (b >> 8)),
					A: uint8(a >> 8),
				})
			}
		}

	case "Rescaling":
		for y := 0; y < bounds.Dy(); y++ {
			for x := 0; x < bounds.Dx(); x++ {
				r, g, b, a := img.At(x, y).RGBA()
				resultImage.Set(x, y, color.RGBA{
					R: uint8(clamp(int(float64(r>>8)*scalingFactor), 0, 255)),
					G: uint8(clamp(int(float64(g>>8)*scalingFactor), 0, 255)),
					B: uint8(clamp(int(float64(b>>8)*scalingFactor), 0, 255)),
					A: uint8(a >> 8),
				})
			}
		}

	case "Shift&Rescale":
		for y := 0; y < bounds.Dy(); y++ {
			for x := 0; x < bounds.Dx(); x++ {
				r, g, b, a := img.At(x, y).RGBA()
				resultImage.Set(x, y, color.RGBA{
					R: uint8(clamp(int(float64(r>>8)*scalingFactor+shiftingValue), 0, 255)),
					G: uint8(clamp(int(float64(g>>8)*scalingFactor+shiftingValue), 0, 255)),
					B: uint8(clamp(int(float64(b>>8)*scalingFactor+shiftingValue), 0, 255)),
					A: uint8(a >> 8),
				})
			}
		}

	case "Bit Plane Slicing":
		for y := 0; y < bounds.Dy(); y++ {
			for x := 0; x < bounds.Dx(); x++ {
				r, g, b, a := img.At(x, y).RGBA()
				resultImage.Set(x, y, color.RGBA{
					R: uint8(((r >> 8) >> nBit) & 1 * 255),
					G: uint8(((g >> 8) >> nBit) & 1 * 255),
					B: uint8(((b >> 8) >> nBit) & 1 * 255),
					A: uint8(a >> 8),
				})
			}
		}

	case "Salt&Pepper noise":
		// Implement Salt & Pepper noise addition
		noiseProbability := 0.02 // Adjust as needed
		for y := 0; y < bounds.Dy(); y++ {
			for x := 0; x < bounds.Dx(); x++ {
				r, g, b, a := img.At(x, y).RGBA()
				if rand.Float64() < noiseProbability {
					if rand.Intn(2) == 0 {
						resultImage.Set(x, y, color.RGBA{0, 0, 0, uint8(a >> 8)}) // Salt
					} else {
						resultImage.Set(x, y, color.RGBA{255, 255, 255, uint8(a >> 8)}) // Pepper
					}
				} else {
					resultImage.Set(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
				}
			}
		}

	default:
		return "", errors.New("unsupported algorithm")
	}

	//processedImagePath := "static/uploads/mixed_algorithms_result.jpg"
	//out, err := os.Create(processedImagePath)
	//if err != nil {
	//	return "", err
	//}
	//defer out.Close()

	// 将结果图像写入内存中的 buffer
	var buf bytes.Buffer
	jpeg.Encode(&buf, resultImage, nil)
	if err != nil {
		return "", err
	}

	// 将图像转换为 base64 编码
	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64Image, nil
}

// 解码图像
func decodeImage(file multipart.File) (image.Image, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}
