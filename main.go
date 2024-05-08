package main

import (
	"fmt"
	"net/http"

	"go-in-memory-assessment/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	store := handlers.NewEmployeeStore()

	r.HandleFunc("/employees", store.ListEmployees).Methods("GET")
	r.HandleFunc("/employee", store.CreateEmployee).Methods("POST")
	r.HandleFunc("/employee/{id}", store.GetEmployee).Methods("GET")
	r.HandleFunc("/employee/{id}", store.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/employee/{id}", store.DeleteEmployee).Methods("DELETE")

	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", r)
}
