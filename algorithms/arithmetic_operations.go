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

// ArithmeticOperations performs basic arithmetic operations on two images
func ArithmeticOperations(file1 multipart.File, file2 multipart.File, operation string) (string, error) {
	img1, err := decodeImage(file1)
	if err != nil {
		return "", err
	}

	img2, err := decodeImage(file2)
	if err != nil {
		return "", err
	}

	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	if bounds1 != bounds2 {
		return "", errors.New("images must have the same dimensions")
	}

	resultImage := image.NewRGBA(bounds1)

	for y := 0; y < bounds1.Dy(); y++ {
		for x := 0; x < bounds1.Dx(); x++ {
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			r2, g2, b2, a2 := img2.At(x, y).RGBA()

			var r, g, b, a uint8

			switch operation {
			case "Addition":
				r = uint8(clamp(int(r1>>8)+int(r2>>8), 0, 255))
				g = uint8(clamp(int(g1>>8)+int(g2>>8), 0, 255))
				b = uint8(clamp(int(b1>>8)+int(b2>>8), 0, 255))
				a = uint8(clamp(int(a1>>8)+int(a2>>8), 0, 255))
			case "Substraction":
				r = uint8(clamp(int(r1>>8)-int(r2>>8), 0, 255))
				g = uint8(clamp(int(g1>>8)-int(g2>>8), 0, 255))
				b = uint8(clamp(int(b1>>8)-int(b2>>8), 0, 255))
				a = uint8(clamp(int(a1>>8)-int(a2>>8), 0, 255))
			case "Multiplication":
				r = uint8(clamp(int(r1>>8)*int(r2>>8), 0, 255))
				g = uint8(clamp(int(g1>>8)*int(g2>>8), 0, 255))
				b = uint8(clamp(int(b1>>8)*int(b2>>8), 0, 255))
				a = uint8(clamp(int(a1>>8)*int(a2>>8), 0, 255))
			case "Division":
				if r2 == 0 || g2 == 0 || b2 == 0 {
					r, g, b, a = uint8(r1>>8), uint8(g1>>8), uint8(b1>>8), uint8(a1>>8) // 保持原有值
				} else {
					r = uint8(clamp(int(r1>>8)/int(r2>>8), 0, 255))
					g = uint8(clamp(int(g1>>8)/int(g2>>8), 0, 255))
					b = uint8(clamp(int(b1>>8)/int(b2>>8), 0, 255))
					a = uint8(clamp(int(a1>>8)/int(a2>>8), 0, 255))
				}
			default:
				return "", errors.New("unsupported operation")
			}

			resultImage.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	//processedImagePath := "static/uploads/arithmetic_operations_result.jpg"
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
	err = jpeg.Encode(&buf, resultImage, nil)
	if err != nil {
		return "", err
	}

	// 将图像转换为 base64 编码
	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64Image, nil
}
