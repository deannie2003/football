package controllers

import (
	"encoding/json"
	"football/models"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

// API xóa tài khoản theo ID
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Chỉ cho phép phương thức DELETE
	if r.Method != "DELETE" {
		http.Error(w, `{"message": "Chỉ hỗ trợ DELETE"}`, http.StatusMethodNotAllowed)
		return
	}

	// Lấy ID người dùng từ query string (hoặc token đăng nhập)
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, `{"message": "Thiếu ID người dùng"}`, http.StatusBadRequest)
		return
	}

	// Chuyển ID sang số nguyên
	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, `{"message": "ID không hợp lệ"}`, http.StatusBadRequest)
		return
	}

	// Tìm và xóa người dùng trong database
	var user models.User
	result := models.DB.First(&user, id)

	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, `{"message": "Người dùng không tồn tại"}`, http.StatusNotFound)
		return
	} else if result.Error != nil {
		http.Error(w, `{"message": "Lỗi truy vấn database"}`, http.StatusInternalServerError)
		return
	}

	// Xóa người dùng
	if err := models.DB.Delete(&user).Error; err != nil {
		http.Error(w, `{"message": "Lỗi khi xóa tài khoản"}`, http.StatusInternalServerError)
		return
	}

	// Phản hồi thành công
	json.NewEncoder(w).Encode(map[string]string{"message": "Tài khoản đã bị xóa thành công"})
}
