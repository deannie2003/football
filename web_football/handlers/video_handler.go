package handlers

import (
	"football/controllers"
	"net/http"
)

// UploadHandler xử lý request upload video
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	controllers.UploadHandler(w, r)
}

// ProcessHandler xử lý request xử lý video
func ProcessHandler(w http.ResponseWriter, r *http.Request) {
	controllers.ProcessHandler(w, r)
}

// StatusHandler xử lý request kiểm tra trạng thái xử lý video
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	controllers.StatusHandler(w, r)
}

// ResultHandler xử lý request hiển thị kết quả video
func ResultHandler(w http.ResponseWriter, r *http.Request) {
	controllers.ResultHandler(w, r)
}

// DownloadHandler xử lý request tải xuống video
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	controllers.DownloadHandler(w, r)
}
