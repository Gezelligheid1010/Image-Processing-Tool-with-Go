package algorithms

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"mime/multipart"
)

// BitOperations performs bitwise operations on two images
func BitOperations(file1 multipart.File, file2 multipart.File, operation string) (string, error) {
	img1, err := decodeImage(file1)
	if err != nil {
		return "", err
	}

	var img2 image.Image
	if file2 != nil {
		img2, err = decodeImage(file2)
		if err != nil {
			return "", err
		}
	}

	bounds1 := img1.Bounds()

	if operation != "Bitwise Not" && img2 != nil {
		bounds2 := img2.Bounds()
		if bounds1 != bounds2 {
			return "", errors.New("images must have the same dimensions")
		}
	}

	resultImage := image.NewRGBA(bounds1)

	for y := 0; y < bounds1.Dy(); y++ {
		for x := 0; x < bounds1.Dx(); x++ {
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			var r2, g2, b2, a2 uint32

			if operation != "Bitwise Not" && img2 != nil {
				r2, g2, b2, a2 = img2.At(x, y).RGBA()
			}

			var r, g, b, a uint8

			switch operation {
			case "Bitwise Not":
				r = ^uint8(r1 >> 8)
				g = ^uint8(g1 >> 8)
				b = ^uint8(b1 >> 8)
				a = uint8(a1 >> 8)
				//log.Printf("Pixel at (%d, %d) - r1: %d, g1: %d, b1: %d, a1: %d\n", x, y, r1, g1, b1, a1)
			case "Bitwise And":
				r = uint8((r1 >> 8) & (r2 >> 8))
				g = uint8((g1 >> 8) & (g2 >> 8))
				b = uint8((b1 >> 8) & (b2 >> 8))
				a = uint8((a1 >> 8) & (a2 >> 8))
			case "Bitwise Or":
				r = uint8((r1 >> 8) | (r2 >> 8))
				g = uint8((g1 >> 8) | (g2 >> 8))
				b = uint8((b1 >> 8) | (b2 >> 8))
				a = uint8((a1 >> 8) | (a2 >> 8))
			case "Bitwise Xor":
				r = uint8((r1 >> 8) ^ (r2 >> 8))
				g = uint8((g1 >> 8) ^ (g2 >> 8))
				b = uint8((b1 >> 8) ^ (b2 >> 8))
				a = uint8((a1 >> 8) ^ (a2 >> 8))
			default:
				return "", errors.New("unsupported operation")
			}

			resultImage.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	//processedImagePath := "static/uploads/bit_operations_result.jpg"
	//out, err := os.Create(processedImagePath)
	//if err != nil {
	//	return "", err
	//}
	//defer out.Close()
	//
	//jpeg.Encode(out, resultImage, nil)
	//return processedImagePath, nil

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
