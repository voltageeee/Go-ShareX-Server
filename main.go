package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const uploadDir = "./uploads"

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "Error parsing the form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving the file: %v", err)
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	os.MkdirAll(uploadDir, os.ModePerm)

	filePath := filepath.Join(uploadDir, filepath.Base(fileHeader.Filename))
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file %s: %v", filePath, err)
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Printf("Error copying the file: %v", err)
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	fileURL := "http://youripordomain.here/uploads/" + filepath.Base(fileHeader.Filename)
	w.Write([]byte(fileURL))
}

func serveUploads(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))).ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/uploads/", serveUploads)
	log.Println("Starting server on :8040")
	log.Fatal(http.ListenAndServe(":8040", nil))
}
