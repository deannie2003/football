package controllers

import (
	"football/models"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	uploadDir = "./uploads/"
	outputDir = "./outputs/"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

// UploadHandler hiển thị trang upload video
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "upload.html", nil)
}

// ProcessHandler xử lý video và hiển thị trang chờ
func ProcessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Lấy file tải lên
	file, header, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Failed to upload video", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Lưu file vào thư mục uploads
	filePath := filepath.Join(uploadDir, header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save video", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, file)

	// Đánh dấu file đang được xử lý
	fileName := header.Filename
	models.SetProcessingStatus(fileName, true)

	// Xử lý video trong goroutine
	go func(inputPath, fileName string) {
		outputPath := filepath.Join(outputDir, "processed_"+fileName)
		cmd := exec.Command("python3", "scripts/process_video.py", inputPath, outputPath)
		cmd.Run()
		models.SetProcessingStatus(fileName, false) // Đánh dấu xử lý xong
	}(filePath, fileName)

	// Hiển thị trang chờ
	tpl.ExecuteTemplate(w, "waiting.html", map[string]string{"FileName": fileName})
}

// StatusHandler kiểm tra trạng thái xử lý video
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

// ResultHandler hiển thị video kết quả
func ResultHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	tpl.ExecuteTemplate(w, "result.html", map[string]string{"FileName": "processed_" + fileName})
}

// DownloadHandler tải xuống video
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	filePath := filepath.Join(outputDir, fileName)

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}
