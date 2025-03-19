package controllers

import (
	"football/configs"
	"football/models"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)


var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}


func UploadHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "upload.html", nil)
}

func ProcessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Failed to upload video", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filePath := filepath.Join(configs.UploadDir, header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save video", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, file)

	fileName := header.Filename
	log.Println("Uploaded file:", fileName)
	models.SetProcessingStatus(fileName, true)

	// Xử lý video trong goroutine
	go func(inputPath, fileName string) {
		outputPath := filepath.Join(configs.OutputDir, "processed_"+fileName)
		cmd := exec.Command("python3", "scripts/process_video.py", inputPath, outputPath) 
		cmd.Run()
		models.SetProcessingStatus(fileName, false)
	}(filePath, fileName)

	// Hiển thị trang chờ
	tpl.ExecuteTemplate(w, "waiting.html", map[string]string{"FileName": fileName})
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	if models.GetProcessingStatus(fileName) {
		w.Write([]byte("processing"))
	} else {
		w.Write([]byte("done"))
	}
}

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	tpl.ExecuteTemplate(w, "result.html", map[string]string{ "FileName": "processed_" + fileName })
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	filePath := filepath.Join(configs.OutputDir, fileName)

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}
