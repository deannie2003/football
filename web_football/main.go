package main

import (
	"football/handlers"
	"net/http"
)

const (
	uploadDir = "./uploads/"
	outputDir = "./outputs/"
)

func main() {
	http.HandleFunc("/", handlers.UploadHandler)
	http.HandleFunc("/process", handlers.ProcessHandler)
	http.HandleFunc("/status", handlers.StatusHandler)
	http.HandleFunc("/result", handlers.ResultHandler)
	http.HandleFunc("/download", handlers.DownloadHandler)
	http.Handle("/outputs/", http.StripPrefix("/outputs/", http.FileServer(http.Dir(outputDir))))

	http.ListenAndServe(":8080", nil)
}
