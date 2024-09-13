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

// Convolution 通用卷积函数
func Convolution(file multipart.File, algorithm string) (string, error) {
	img, err := decodeImage(file)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	convolutionImage := image.NewRGBA(bounds)

	var kernel [][]int
	var divisor int

	// 根据算法选择不同的卷积核
	switch algorithm {
	case "Convolution - Averaging":
		kernel = AveragingKernel
		divisor = 9

	case "Convolution - Weighted averaging":
		kernel = WeightedAveragingKernel
		divisor = 16

	case "Convolution - Four Neighbour Laplacian":
		kernel = FourNeighbourLaplacianKernel
		divisor = 1

	case "Convolution - Eight Neighbour Laplacian":
		kernel = EightNeighbourLaplacianKernel
		divisor = 1

	case "Convolution - Four Neighbour Laplacian Enhancement":
		kernel = FourNeighbourLaplacianEnhancementKernel
		divisor = 1

	case "Convolution - Eight Neighbour Laplacian Enhancement":
		kernel = EightNeighbourLaplacianEnhancementKernel
		divisor = 1

	case "Convolution - Roberts One":
		kernel = RobertsOneKernel
		divisor = 1

	case "Convolution - Roberts Two":
		kernel = RobertsTwoKernel
		divisor = 1

	case "Convolution - Sobel X":
		kernel = SobelXKernel
		divisor = 1

	case "Convolution - Sobel Y":
		kernel = SobelYKernel
		divisor = 1

	default:
		return "", errors.New("unsupported convolution algorithm")
	}

	// 卷积核的尺寸
	kernelSize := len(kernel)
	offset := kernelSize / 2

	// 遍历图像的每个像素
	for y := offset; y < bounds.Dy()-offset; y++ {
		for x := offset; x < bounds.Dx()-offset; x++ {
			var rSum, gSum, bSum int

			// 遍历卷积核
			for ky := 0; ky < kernelSize; ky++ {
				for kx := 0; kx < kernelSize; kx++ {
					px := x + kx - offset
					py := y + ky - offset
					r, g, b, _ := img.At(px, py).RGBA()

					rSum += int(r>>8) * kernel[ky][kx]
					gSum += int(g>>8) * kernel[ky][kx]
					bSum += int(b>>8) * kernel[ky][kx]
				}
			}

			// 卷积后的新像素值
			rSum /= divisor
			gSum /= divisor
			bSum /= divisor

			// 确保值在0-255之间
			rNew := uint8(clamp(rSum, 0, 255))
			gNew := uint8(clamp(gSum, 0, 255))
			bNew := uint8(clamp(bSum, 0, 255))

			convolutionImage.Set(x, y, color.RGBA{rNew, gNew, bNew, 255})
		}
	}

	//processedImagePath := "static/uploads/convolution_result.jpg"
	//out, err := os.Create(processedImagePath)
	//if err != nil {
	//	return "", err
	//}
	//defer out.Close()
	//
	//jpeg.Encode(out, convolutionImage, nil)
	//return processedImagePath, nil
	// 将结果图像写入内存中的 buffer
	var buf bytes.Buffer
	jpeg.Encode(&buf, convolutionImage, nil)
	if err != nil {
		return "", err
	}

	// 将图像转换为 base64 编码
	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64Image, nil
}

// 各种卷积核定义

// AveragingKernel 均值滤波器
var AveragingKernel = [][]int{
	{1, 1, 1},
	{1, 1, 1},
	{1, 1, 1},
}

// WeightedAveragingKernel 加权均值滤波器
var WeightedAveragingKernel = [][]int{
	{1, 2, 1},
	{2, 4, 2},
	{1, 2, 1},
}

// FourNeighbourLaplacianKernel 四邻域拉普拉斯核
var FourNeighbourLaplacianKernel = [][]int{
	{0, -1, 0},
	{-1, 4, -1},
	{0, -1, 0},
}

// EightNeighbourLaplacianKernel 八邻域拉普拉斯核
var EightNeighbourLaplacianKernel = [][]int{
	{-1, -1, -1},
	{-1, 8, -1},
	{-1, -1, -1},
}

// FourNeighbourLaplacianEnhancementKernel 四邻域拉普拉斯增强
var FourNeighbourLaplacianEnhancementKernel = [][]int{
	{0, -1, 0},
	{-1, 5, -1},
	{0, -1, 0},
}

// EightNeighbourLaplacianEnhancementKernel 八邻域拉普拉斯增强
var EightNeighbourLaplacianEnhancementKernel = [][]int{
	{-1, -1, -1},
	{-1, 9, -1},
	{-1, -1, -1},
}

// RobertsOneKernel 罗伯茨算子1
var RobertsOneKernel = [][]int{
	{1, 0},
	{0, -1},
}

// RobertsTwoKernel 罗伯茨算子2
var RobertsTwoKernel = [][]int{
	{0, 1},
	{-1, 0},
}

// SobelXKernel Sobel X 滤波器
var SobelXKernel = [][]int{
	{-1, 0, 1},
	{-2, 0, 2},
	{-1, 0, 1},
}

// SobelYKernel Sobel Y 滤波器
var SobelYKernel = [][]int{
	{-1, -2, -1},
	{0, 0, 0},
	{1, 2, 1},
}

// 限制数值在0到255之间
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
