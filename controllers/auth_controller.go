package controllers

import (
	"encoding/json"
	"football/models"
	"football/utils"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Xử lý đăng ký người dùng
func SignupHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json") 

    if r.Method != "POST" {
        http.Error(w, `{"message": "Chỉ hỗ trợ POST", "success": false}`, http.StatusMethodNotAllowed)
        return
    }

    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, `{"message": "Dữ liệu không hợp lệ", "success": false}`, http.StatusBadRequest)
        return
    }

    // Kiểm tra tài khoản đã tồn tại chưa
    var existingUser models.User
    result := models.DB.Where("username = ?", user.Username).First(&existingUser)

    if result.Error == nil {
        http.Error(w, `{"message": "Tên đăng nhập đã tồn tại", "success": false}`, http.StatusConflict)
        return
    } else if result.Error != gorm.ErrRecordNotFound {
        http.Error(w, `{"message": "Lỗi truy vấn database", "success": false}`, http.StatusInternalServerError)
        return
    }

    // Mã hóa mật khẩu
    hashedPassword, err := HashPassword(user.Password)
    if err != nil {
        http.Error(w, `{"message": "Lỗi mã hóa mật khẩu", "success": false}`, http.StatusInternalServerError)
        return
    }
    user.Password = hashedPassword
	if err := models.DB.Create(&user).Error; err != nil {
        log.Println("Lỗi lưu user vào DB:", err)
        http.Error(w, `{"message": "Lỗi lưu dữ liệu", "success": false}`, http.StatusInternalServerError)
        return
    }

    subject := "Chào mừng bạn đến với hệ thống của chúng tôi!"
    body := "Cảm ơn bạn đã đăng ký. Tài khoản của bạn đã được tạo thành công.\n\nChúc bạn một ngày vui vẻ!"
    err = utils.SendEmail(user.Email, subject, body)
    if err != nil {
        log.Println("Lỗi gửi email:", err)
    }

    // ✅ **Phản hồi về FE sau khi tất cả hoàn thành**
    json.NewEncoder(w).Encode(Response{Message: "Đăng ký thành công, kiểm tra email!", Success: true})
}


// Xử lý đăng nhập
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Chỉ hỗ trợ phương thức POST", http.StatusMethodNotAllowed)
		return
	}

	var loginData models.User
	json.NewDecoder(r.Body).Decode(&loginData)

	// Tìm kiếm tài khoản
	var user models.User
	if err := models.DB.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		http.Error(w, "Sai tên đăng nhập hoặc mật khẩu", http.StatusUnauthorized)
		return
	}

	// Kiểm tra mật khẩu
	if !CheckPassword(user.Password, loginData.Password) {
		http.Error(w, "Sai tên đăng nhập hoặc mật khẩu", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(Response{Message: "Đăng nhập thành công", Success: true})
}

// Lấy danh sách tất cả người dùng đã đăng ký
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := models.DB.Find(&users).Error; err != nil {
		http.Error(w, "Không thể truy xuất danh sách người dùng", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}