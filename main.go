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

	//Cargamos las variables de entorno que están en el archivo .env por medio del package godotenv
	_ = godotenv.Load()

	//Instanciamos un logger propio
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	//Contruímos el string de conexión a la bd por medio de los envs
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	//Abrimos la instancia de base de datos por medio de GORM y la inicializamos en modo debug
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db = db.Debug()

	//Migra el "modelo" a una tabla SQL
	_ = db.AutoMigrate(user.User{})

	//Instancias de las capas: repositorio, servicio y controlador
	userRepo := user.NewRepo(l, db)
	userSrv := user.NewService(l, userRepo)
	userEnd := user.MakeEndpoints(userSrv)

	//Por medio del router de Gorilla Mux servimos los endpoints
	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	//Se crea una instancia de un servidor
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	//Se sirve la aplicación y se le vanta el servidor
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
