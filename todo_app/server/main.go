package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // React's default port
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Apply CORS middleware to the router
	handler := c.Handler(router)

	// Define your routes here
	router.HandleFunc("/api/todos", handleTodos).Methods("GET", "POST")
	router.HandleFunc("/api/todos/{id}", handleTodo).Methods("GET", "PUT", "DELETE")

	log.Println("Server starting on port 5000...")
	log.Fatal(http.ListenAndServe(":5000", handler))
}

func handleTodos(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement todo handlers
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Todo endpoint"}`))
}

func handleTodo(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement single todo handlers
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Single todo endpoint"}`))
}
