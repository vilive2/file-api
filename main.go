package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func download(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		filename := r.URL.Query().Get("filename")
		if filename == "" {
			http.Error(w, "filename is required", http.StatusBadRequest)
			return
		}

		filename = filepath.Clean(filename)
		fileDir := "uploads"
		filePath := filepath.Join(fileDir, filename)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// file, err := os.Open(filePath)
		// if err != nil {
		// 	log.Println(err)
		// 	http.Error(w, "Unable to open file", http.StatusInternalServerError)
		// 	return
		// }
		// defer file.Close()

		fileExt := strings.ToLower(filepath.Ext(filename))
		var contentType string
		switch fileExt {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		default:
			contentType = "application/octet-stream"
		}

		w.Header().Set("Content-Type", contentType)
		http.ServeFile(w, r, filePath)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodPost {

		err := r.ParseMultipartForm(4 << 20)
		if err != nil {
			log.Println(err)
			http.Error(w, "unable to parse form", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			http.Error(w, "error retrieving the file", http.StatusInternalServerError)
			return
		}

		defer file.Close()

		ext := filepath.Ext(handler.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

		filePath := filepath.Join("uploads", filename)
		out, err := os.Create(filePath)
		if err != nil {
			log.Println(err)
			http.Error(w, "unable to save file", http.StatusInternalServerError)
			return
		}

		defer out.Close()

		_, err = out.ReadFrom(file)
		if err != nil {
			log.Println(err)
			http.Error(w, "error saving file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", filename)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {

	http.HandleFunc("/file/upload", upload)
	http.HandleFunc("/file/download", download)

	http.ListenAndServe(":3001", nil)
}
