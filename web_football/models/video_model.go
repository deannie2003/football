package models

var processingStatus = map[string]bool{} // Trạng thái xử lý video

// SetProcessingStatus cập nhật trạng thái xử lý video
func SetProcessingStatus(fileName string, status bool) {
	processingStatus[fileName] = status
}

// GetProcessingStatus lấy trạng thái xử lý video
func GetProcessingStatus(fileName string) bool {
	return processingStatus[fileName]
}
