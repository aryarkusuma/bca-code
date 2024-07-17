package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		return
	}
	parentDir := filepath.Dir(filepath.Dir(exePath))

	directories := []string{"KTP", "SPK", "Bbtj"}
	for _, dir := range directories {
		fullPath := filepath.Join(parentDir, dir)
		err := os.MkdirAll(fullPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory %s: %v\n", fullPath, err)
			return
		}
		fmt.Printf("Created directory: %s\n", fullPath)
	}

	http.HandleFunc("/upload/KTP", handleUpload(filepath.Join(parentDir, "KTP")))
	http.HandleFunc("/upload/SPK", handleUpload(filepath.Join(parentDir, "SPK")))
	http.HandleFunc("/upload/Bbtj", handleUpload(filepath.Join(parentDir, "Bbtj")))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleUpload(directory string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(10 << 20) // 10 MB max memory
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		dst, err := os.Create(filepath.Join(directory, header.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "File uploaded successfully to %s", directory)
	}
}
