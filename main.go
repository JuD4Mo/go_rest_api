package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	port := ":3000"
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/courses", getCourses)

	//Servir la app y levantar el servidor
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /users")
	io.WriteString(w, "user endpoint\n")
}

func getCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /courses")
	io.WriteString(w, "course endpoint\n")
}
