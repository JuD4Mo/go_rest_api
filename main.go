package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JuD4Mo/go_rest_api/internal/user"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	//Instancia de un router de Gorilla Mux
	router := mux.NewRouter()
	_ = godotenv.Load()
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db = db.Debug()

	_ = db.AutoMigrate(user.User{})

	userRepo := user.NewRepo(l, db)
	userSrv := user.NewService(l, userRepo)
	userEnd := user.MakeEndpoints(userSrv)

	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users", userEnd.Get).Methods("GET")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users", userEnd.Delete).Methods("DELETE")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	// port := ":3000"
	// http.HandleFunc("/users", getUsers)
	// http.HandleFunc("/courses", getCourses)

	// //Servir la app y levantar el servidor
	// err := http.ListenAndServe(port, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

// func getUsers(w http.ResponseWriter, r *http.Request) {
// 	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
// 	fmt.Println("got /users")
// 	io.WriteString(w, "user endpoint\n")
// }

// func getCourses(w http.ResponseWriter, r *http.Request) {
// 	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
// 	fmt.Println("got /courses")
// 	io.WriteString(w, "course endpoint\n")
// }
