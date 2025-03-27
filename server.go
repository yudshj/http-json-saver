package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// allowedOrigins 是一个包含允许访问的来源的集合
var allowedOrigins = map[string]bool{
	"https://yudshj.synology.me": true,
	"http://127.0.0.1":           true,
}

var enableOriginCheck bool

// 定义一个结构体来存储请求数据
type Payload struct {
	Name       string      `json:"name"`
	MajorRunId string      `json:"majorRunId"`
	MinorRunId string      `json:"minorRunId"`
	Data       interface{} `json:"data"` // Use interface{} if the data field can hold various types
}

type requestData struct {
	Name    string
	Body    []byte
	Payload Payload
}

var (
	requestQueue []requestData
	queueMutex   sync.Mutex
)

func saveJSONHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	// 如果启用了 Origin 验证，检查请求的来源是否在允许的来源列表中
	if enableOriginCheck && !allowedOrigins[origin] {
		http.Error(w, "Forbidden: Invalid origin", http.StatusForbidden)
		return
	}

	// 设置 CORS 头
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 处理预检请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the JSON body into the Payload struct
	var payload Payload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate the 'name' field
	if payload.Name == "" {
		http.Error(w, "Missing or invalid 'name' field in JSON", http.StatusBadRequest)
		return
	}

	// Add the request to the queue
	queueMutex.Lock()
	requestQueue = append(requestQueue, requestData{Name: payload.Name, Body: body, Payload: payload})
	queueMutex.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("JSON received successfully"))
}

func batchSaveToDisk() {
	// 准备输出目录
	outputDir := "./json_out"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Println("Unable to create output directory:", err)
		os.Exit(1)
	}

	for {
		time.Sleep(5 * time.Second)
		if len(requestQueue) == 0 {
			continue
		}

		queueMutex.Lock()

		// 处理队列中的请求
		for _, request := range requestQueue {
			outputDir2 := filepath.Join(outputDir, request.Payload.MajorRunId)
			if err := os.MkdirAll(outputDir2, os.ModePerm); err != nil {
				fmt.Println("Unable to create output directory:", err)
				continue
			}

			filePath := filepath.Join(outputDir2, request.Name+".json")
			if err := os.WriteFile(filePath, request.Body, 0644); err != nil {
				fmt.Println("Unable to write file:", err)
			}
		}

		// 清空队列
		requestQueue = nil
		queueMutex.Unlock()
	}
}

func main() {
	// 使用 flag 包来解析命令行参数
	var host string
	var port int
	flag.StringVar(&host, "host", "127.0.0.1", "The host to listen on")
	flag.IntVar(&port, "port", 3000, "The port to listen on")
	flag.BoolVar(&enableOriginCheck, "enable-origin-check", true, "Enable or disable origin check")
	flag.Parse()

	// 启动批处理保存的 Goroutine
	go batchSaveToDisk()

	http.HandleFunc("/save", saveJSONHandler)

	address := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Starting server on %s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
