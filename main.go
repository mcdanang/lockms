package main

import (
	"database/sql"
	"go-app-be/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// jsonContentTypeMiddleware is a middleware to set the Content-Type header to application/json
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func createTablesIfNotExist(db *sql.DB) {
	// Create keys table if not exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS keys (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT
		)
	`)
	if err != nil {
		log.Fatal("Error creating keys table: ", err)
	}

	// Create key_copies table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS key_copies (
			id SERIAL PRIMARY KEY,
			key_id INTEGER REFERENCES keys(id),
			staff_id INTEGER
		)
	`)
	if err != nil {
		log.Fatal("Error creating key_copies table: ", err)
	}

	// Create staff table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS staff (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			role TEXT
		)
	`)
	if err != nil {
		log.Fatal("Error creating staff table: ", err)
	}
}

func main() {
	// Initialize the database connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tables if they don't exist
	createTablesIfNotExist(db)

	// Initialize the router
	router := mux.NewRouter()

	// Setup routes
	routes.SetupRoutes(router, db)

	// Apply JSON middleware
	router.Use(jsonContentTypeMiddleware)

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", router))
}
