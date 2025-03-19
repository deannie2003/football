package controllers

import (
	"encoding/json"
	"fmt"
	"football/models"
	"football/utils"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Hàm tạo mật khẩu ngẫu nhiên
func GenerateRandomPassword() string {
	rand.Seed(time.Now().UnixNano())
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, 10)
	for i := range password {
		password[i] = chars[rand.Intn(len(chars))]
	}
	return string(password)
}

// API xử lý yêu cầu quên mật khẩu
func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, `{"message": "Chỉ hỗ trợ phương thức POST"}`, http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, `{"message": "Dữ liệu không hợp lệ"}`, http.StatusBadRequest)
		return
	}

	var user models.User
	result := models.DB.Where("email = ?", requestData.Email).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, `{"message": "Email không tồn tại trong hệ thống"}`, http.StatusNotFound)
		return
	} else if result.Error != nil {
		http.Error(w, `{"message": "Lỗi truy vấn database"}`, http.StatusInternalServerError)
		return
	}

	// Tạo mật khẩu mới
	newPassword := GenerateRandomPassword()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"message": "Lỗi mã hóa mật khẩu"}`, http.StatusInternalServerError)
		return
	}

	// Cập nhật mật khẩu mới vào database
	user.Password = string(hashedPassword)
	models.DB.Save(&user)

	// Gửi email với mật khẩu mới
	subject := "Đặt lại mật khẩu thành công"
	body := fmt.Sprintf("Mật khẩu mới của bạn là: %s\nVui lòng đăng nhập và đổi mật khẩu ngay lập tức.", newPassword)
	err = utils.SendEmailPassWord(user.Email, subject, body)
	if err != nil {
		http.Error(w, `{"message": "Lỗi gửi email"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Mật khẩu mới đã được gửi đến email của bạn!"})
}
