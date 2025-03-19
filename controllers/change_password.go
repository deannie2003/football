package controllers

import (
	"encoding/json"
	"football/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// API đổi mật khẩu
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, `{"message": "Chỉ hỗ trợ phương thức POST"}`, http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Username    string `json:"username"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, `{"message": "Dữ liệu không hợp lệ"}`, http.StatusBadRequest)
		return
	}

	// Lấy thông tin người dùng từ session hoặc token (ở đây giả lập username "dean")
	var user models.User
	if err := models.DB.Where("username = ?", "dean").First(&user).Error; err != nil {
		http.Error(w, `{"message": "Người dùng không tồn tại"}`, http.StatusUnauthorized)
		return
	}

	// Kiểm tra mật khẩu cũ
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestData.OldPassword)); err != nil {
		http.Error(w, `{"message": "Mật khẩu cũ không đúng"}`, http.StatusUnauthorized)
		return
	}

	// Mã hóa mật khẩu mới
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"message": "Lỗi mã hóa mật khẩu"}`, http.StatusInternalServerError)
		return
	}

	// Cập nhật mật khẩu mới vào database
	user.Password = string(hashedPassword)
	models.DB.Save(&user)

	json.NewEncoder(w).Encode(map[string]string{"message": "Mật khẩu đã được thay đổi thành công!"})
}
