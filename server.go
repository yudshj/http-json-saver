package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// allowedOrigins 是一个包含允许访问的来源的集合
var allowedOrigins = map[string]bool{
	"https://yudshj.synology.me": true,
	"http://127.0.0.1":           true,
}

func saveJSONHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	// 检查请求的来源是否在允许的来源列表中
	if !allowedOrigins[origin] {
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

	// Parse the JSON body to extract the "name" field
	var requestBody map[string]interface{}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	name, ok := requestBody["name"].(string)
	if !ok || name == "" {
		http.Error(w, "Missing or invalid 'name' field in JSON", http.StatusBadRequest)
		return
	}

	// Prepare the output directory
	outputDir := "./json_out"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		http.Error(w, "Unable to create output directory", http.StatusInternalServerError)
		return
	}

	// Write body to a file in the output directory
	filePath := filepath.Join(outputDir, fmt.Sprintf("%s.json", name))
	if err := os.WriteFile(filePath, body, 0644); err != nil {
		http.Error(w, "Unable to write file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("JSON saved successfully"))
}

func main() {
	http.HandleFunc("/save", saveJSONHandler)

	fmt.Println("Starting server on 127.0.0.1:3000")
	if err := http.ListenAndServe("127.0.0.1:3000", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
