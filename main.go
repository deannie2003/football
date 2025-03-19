package main

import (
	"football/configs"
	"football/controllers"
	"football/models"
	"log"
	"net/http"
)

func main() {
	models.InitDB()

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))
	http.Handle("/outputs/", http.StripPrefix("/outputs/", http.FileServer(http.Dir(configs.OutputDir))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	//http.Handle("/", http.HandlerFunc(controllers.SigninHandler))
	http.Handle("/", http.FileServer(http.Dir("./templates")))
	http.Handle("/signin", http.HandlerFunc(controllers.SigninHandler))
	http.Handle("/signup", http.HandlerFunc(controllers.SignupHandler))
	http.HandleFunc("/api/users", controllers.GetUsersHandler)
	http.HandleFunc("/api/forgot-password", controllers.ForgotPasswordHandler)
	http.HandleFunc("/api/delete-user", controllers.DeleteUserHandler)
	http.HandleFunc("/api/change-password", controllers.ChangePasswordHandler)
	//http.Handle("/signout", http.HandlerFunc(controllers.SignoutHandler))
	http.HandleFunc("/upload", controllers.UploadHandler)
	http.HandleFunc("/process", controllers.ProcessHandler)
	http.HandleFunc("/status", controllers.StatusHandler)
	http.HandleFunc("/result", controllers.ResultHandler)
	http.HandleFunc("/download", controllers.DownloadHandler)

	log.Println("Server started on: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
