package models

var processingStatus = map[string]bool{} // process video status

func SetProcessingStatus(fileName string, status bool) {
	processingStatus[fileName] = status
}

func GetProcessingStatus(fileName string) bool {
	return processingStatus[fileName]
}
