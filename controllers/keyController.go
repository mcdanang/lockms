package controllers

import (
	"database/sql"
	"encoding/json"
	"go-app-be/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Get all keys
func GetKeys(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve page from query parameters
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page <= 0 {
			page = 1 // Default to page 1 if invalid or non-positive page
		}

		// Set page size and calculate offset
		pageSize := 10
		offset := (page - 1) * pageSize

		// Query the database with LIMIT and OFFSET
		rows, err := db.Query("SELECT * FROM keys LIMIT $1 OFFSET $2", pageSize, offset)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var keys []models.Key
		for rows.Next() {
			var k models.Key
			if err := rows.Scan(&k.ID, &k.Name, &k.Description); err != nil {
				log.Fatal(err)
			}
			keys = append(keys, k)
		}

		// Return the result as JSON
		json.NewEncoder(w).Encode(keys)
	}
}

// Get a specific key by ID
func GetKey(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var k models.Key
		err := db.QueryRow("SELECT * FROM keys WHERE id = $1", id).Scan(&k.ID, &k.Name, &k.Description)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(k)
	}
}

// Create a new key
func CreateKey(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var k models.Key
		json.NewDecoder(r.Body).Decode(&k)

		err := db.QueryRow("INSERT INTO keys (name, description) VALUES ($1, $2) RETURNING id", k.Name, k.Description).Scan(&k.ID)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(k)
	}
}

// Update a key
func UpdateKey(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var k models.Key
		json.NewDecoder(r.Body).Decode(&k)

		vars := mux.Vars(r)
		id := vars["id"]

		var existingKey models.Key
		err := db.QueryRow("SELECT * FROM keys WHERE id = $1", id).Scan(&existingKey.ID, &existingKey.Name, &existingKey.Description)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("UPDATE keys SET name = $1, description = $2 WHERE id = $3", k.Name, k.Description, id)
			if err != nil {
				log.Fatal(err)
			}
		}

		json.NewEncoder(w).Encode(k)
	}
}

// Delete a key
func DeleteKey(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var existingKey models.Key
		err := db.QueryRow("SELECT * FROM keys WHERE id = $1", id).Scan(&existingKey.ID, &existingKey.Name, &existingKey.Description)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM keys WHERE id = $1", id)
			if err != nil {
				log.Fatal(err)
			}
		}

		json.NewEncoder(w).Encode("Key deleted")

	}
}
