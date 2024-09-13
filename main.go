package main

import (
	"WebAssembly-Based_Image_Processing_Tool/handlers"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	// 配置 CORS 允许所有来源（可以根据需求限制来源）
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:63342"}, // 允许的前端域名
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// 使用你的处理程序
	mux := http.NewServeMux()
	mux.HandleFunc("/imageProcessing/process", handlers.ProcessMixedAlgorithms)
	mux.HandleFunc("/imageProcessing/process/arithmeticOperations", handlers.ProcessArithmeticOperations)
	mux.HandleFunc("/imageProcessing/process/bitOperations", handlers.ProcessBitOperations)
	mux.HandleFunc("/imageProcessing/process/convolution", handlers.ProcessConvolution)
	mux.HandleFunc("/imageProcessing/process/transformations", handlers.ProcessTransformations)

	// 将处理器包装在 CORS 中
	handler := c.Handler(mux)

	// 启动服务器
	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
