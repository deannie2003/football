package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
}

// func checkDBFile() {
// 	_, err := os.Stat("users.db")
// 	if os.IsNotExist(err) {
// 		fmt.Println("⚠️ File users.db không tồn tại! Hãy kiểm tra lại đường dẫn.")
// 	} else {
// 		fmt.Println("✅ File users.db đã được mở!")
// 	}
// }

// Kết nối và tạo bảng
func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Không thể kết nối database: ", err)		
	}
	//checkDBFile()
	// Tự động tạo bảng nếu chưa có
	DB.AutoMigrate(&User{})
	DB = DB.Debug()
}

