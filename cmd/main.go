package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"com.notes/notes/internal/backend"
	"com.notes/notes/internal/db"
	"com.notes/notes/internal/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	serverPort := os.Getenv("SERVER_PORT")
	db.DataBaseUrl = os.Getenv("DB_URL")
	db.DataBasePassword = os.Getenv("DB_PASSWORD")

	router := mux.NewRouter()
	router.HandleFunc("/api", GivenWellcome).Methods("GET")
	router.HandleFunc("/api/notes/{id}", ReadNote).Methods("GET")
	router.HandleFunc("/api/notes", ReadAllNote).Methods("GET")
	router.HandleFunc("/api/notes", WriteNote).Methods("POST")
	router.HandleFunc("/api/notes/{id}", DeleteNote).Methods("DELETE")
	err := http.ListenAndServe(":"+serverPort, router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func GivenWellcome(w http.ResponseWriter, r *http.Request) {
	writeResponse(http.StatusOK, map[string]string{"message": "Wellcome to Notes API"}, w)
}

func ReadNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	note, err := backend.GetNote(key)
	if err != nil {
		if err.Error() == "key not found" {
			writeResponse(http.StatusNotFound, map[string]string{"error": err.Error()}, w)
			return
		} else {
			writeResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}, w)
			return
		}
	}
	writeResponse(http.StatusOK, note, w)
}

func ReadAllNote(w http.ResponseWriter, r *http.Request) {
	notes, err := backend.GetAllNote()
	if err != nil {
		writeResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}, w)
		return
	}
	writeResponse(http.StatusOK, notes, w)
}

func WriteNote(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var note models.Note

	if err := decoder.Decode(&note); err != nil {
		writeResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}, w)
		return
	}

	fmt.Println("note ", note)

	key, err := backend.SaveNote(&note)
	if err != nil {
		writeResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}, w)
		return
	}

	writeResponse(http.StatusOK, map[string]string{"key": key}, w)
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Println("key ", key)
	err := backend.DeleteNote(key)
	if err != nil {
		writeResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}, w)
		return
	}
	writeResponse(http.StatusOK, map[string]string{"message": "Note deleted", "key": key}, w)
}

func writeResponse(status int, body interface{}, w http.ResponseWriter) {
	fmt.Println("writeResponse ", status, body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	payload, _ := json.Marshal(body)
	w.Write(payload)
}
