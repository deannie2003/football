package utils

import (
	"fmt"
	"net/smtp"
	"football/configs"
)

// Gửi email thông báo
func SendEmail(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", configs.SenderEmail, configs.SenderPass, configs.SMTPServer)

	// Nội dung email
	message := []byte("Subject: " + subject + "\r\n" +
		"From: " + configs.SenderEmail + "\r\n" +
		"To: " + to + "\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n" +
		body)

	// Gửi email
	err := smtp.SendMail(configs.SMTPServer+":"+configs.SMTPPort, auth, configs.SenderEmail, []string{to}, message)
	if err != nil {
		fmt.Println("Lỗi gửi email:", err)
		return err
	}

	fmt.Println("Email đã được gửi tới", to)
	return nil
}

// Gửi email chứa mật khẩu mới
func SendEmailPassWord(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", configs.SenderEmail, configs.SenderPass, configs.SMTPServer)

	message := []byte("Subject: " + subject + "\r\n" +
		"From: " + configs.SenderEmail + "\r\n" +
		"To: " + to + "\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n" +
		body)

	err := smtp.SendMail(configs.SMTPServer+":"+configs.SMTPPort, auth, configs.SenderEmail, []string{to}, message)
	if err != nil {
		fmt.Println("Lỗi gửi email:", err)
		return err
	}

	fmt.Println("Email đã gửi thành công đến", to)
	return nil
}
