package handlers

import (
	"WebAssembly-Based_Image_Processing_Tool/algorithms"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

// ProcessMixedAlgorithms 处理混合算法 (Rescaling, Negative, Shift&Rescale, etc.)
func ProcessMixedAlgorithms(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return
	}
	algorithm := r.FormValue("algorithm")
	log.Printf("Algorithm selected: %s", algorithm)

	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("Error reading file: %v", err)
		http.Error(w, "Invalid image upload", http.StatusBadRequest)
		return
	}
	log.Printf("Uploaded file: %s, size: %d, header: %v", header.Filename, header.Size, header.Header)
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var base64Image string
	switch algorithm {
	case "Rescaling":
		scalingFactor, err := strconv.ParseFloat(r.FormValue("scalingFactor"), 64)
		if err != nil {
			http.Error(w, "Invalid scaling factor", http.StatusBadRequest)
			return
		}
		base64Image, err = algorithms.MixedAlgorithms(file, algorithm, scalingFactor, 0, 0)
	case "Negative":
		base64Image, err = algorithms.MixedAlgorithms(file, algorithm, 0, 0, 0)
		if err != nil {
			http.Error(w, "Error processing image: "+err.Error(), http.StatusInternalServerError)
			return
		}

	case "Shift&Rescale":
		scalingFactor, err := strconv.ParseFloat(r.FormValue("scalingFactor"), 64)
		shiftingValue, err := strconv.ParseFloat(r.FormValue("shiftingValue"), 64)
		if err != nil {
			http.Error(w, "Invalid shift value", http.StatusBadRequest)
			return
		}
		base64Image, err = algorithms.MixedAlgorithms(file, algorithm, scalingFactor, shiftingValue, 0)
	case "Bit Plane Slicing":
		n, err := strconv.Atoi(r.FormValue("nBit"))
		if err != nil {
			http.Error(w, "invalid number of bit plane", http.StatusBadRequest)
			return
		}
		base64Image, err = algorithms.MixedAlgorithms(file, algorithm, 0, 0, n)
	case "Salt&Pepper noise":
		base64Image, err = algorithms.MixedAlgorithms(file, algorithm, 0, 0, 0)
	default:
		http.Error(w, "Unknown algorithm", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Error processing image", http.StatusInternalServerError)
		return
	}

	// 返回 base64 编码的图像
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(base64Image))

	//http.ServeFile(w, r, base64Image)
}

// ProcessArithmeticOperations 处理算术运算 (Addition, Substraction, etc.)
func ProcessArithmeticOperations(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("Error parsing multipart form:", err)
		http.Error(w, "ParseMultipartForm", http.StatusBadRequest)
		return
	}

	file1, _, err := r.FormFile("image")
	if err != nil {
		log.Println("Error uploading first image:", err)
		http.Error(w, "Invalid image upload", http.StatusBadRequest)
		return
	}
	defer file1.Close()

	file2, _, err := r.FormFile("secondImage")
	if err != nil {
		log.Println("Error uploading second image:", err)
		http.Error(w, "Invalid second image upload", http.StatusBadRequest)
		return
	}
	defer file2.Close()

	algorithm := r.FormValue("algorithm")
	log.Println("Selected algorithm:", algorithm)

	var base64Image string
	switch algorithm {
	case "Addition":
		base64Image, err = algorithms.ArithmeticOperations(file1, file2, "Addition")
	case "Substraction":
		base64Image, err = algorithms.ArithmeticOperations(file1, file2, "Substraction")
	case "Multiplication":
		base64Image, err = algorithms.ArithmeticOperations(file1, file2, "Multiplication")
	case "Division":
		base64Image, err = algorithms.ArithmeticOperations(file1, file2, "Division")
	default:
		log.Println("Unknown arithmetic operation:", algorithm)
		http.Error(w, "Unknown arithmetic operation", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Println("Error processing image:", err)
		http.Error(w, "Error processing image", http.StatusInternalServerError)
		return
	}

	// 返回 base64 编码的图像
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(base64Image))

	//http.ServeFile(w, r, processedImagePath)
}

// ProcessBitOperations 处理位运算
func ProcessBitOperations(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return
	}

	file1, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid image upload", http.StatusBadRequest)
		return
	}
	defer file1.Close()

	algorithm := r.FormValue("algorithm")

	var base64Image string
	switch algorithm {
	case "Bitwise Not":
		base64Image, err = algorithms.BitOperations(file1, nil, algorithm)
	case "Bitwise And", "Bitwise Or", "Bitwise Xor":
		file2, _, err := r.FormFile("secondImage")
		if err != nil {
			http.Error(w, "Invalid second image upload", http.StatusBadRequest)
			return
		}
		defer file2.Close()

		base64Image, err = algorithms.BitOperations(file1, file2, algorithm)
	default:
		http.Error(w, "Unknown bit operation", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Error processing image", http.StatusInternalServerError)
		return
	}

	// 返回 base64 编码的图像
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(base64Image))

	//http.ServeFile(w, r, processedImagePath)
}

// ProcessConvolution 处理卷积操作
func ProcessConvolution(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid image upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	algorithm := r.FormValue("algorithm")

	var base64Image string
	base64Image, err = algorithms.Convolution(file, algorithm)
	if err != nil {
		http.Error(w, "Error processing image", http.StatusInternalServerError)
		return
	}

	//http.ServeFile(w, r, processedImagePath)

	// 返回 base64 编码的图像
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(base64Image))
}

// ProcessTransformations 处理图像变换 (Logarithmic Transformation, Power Law, etc.)
func ProcessTransformations(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid image upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	algorithm := r.FormValue("algorithm")
	param, err := strconv.ParseFloat(r.FormValue("param"), 64)
	if err != nil {
		http.Error(w, "Invalid parameter", http.StatusBadRequest)
		return
	}

	var base64Image string
	switch algorithm {
	case "Logarithmic Transformation":

		base64Image, err = algorithms.GeneralTransformation(file, algorithm, param, 0)
	case "Power Law":
		base64Image, err = algorithms.GeneralTransformation(file, algorithm, param, 0)
	case "Random LUT":
		base64Image, err = algorithms.GeneralTransformation(file, algorithm, 0, 0)
	default:
		http.Error(w, "Unknown transformation", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Error processing image", http.StatusInternalServerError)
		return
	}

	//http.ServeFile(w, r, processedImagePath)

	// 返回 base64 编码的图像
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(base64Image))
}
