package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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
	router.HandleFunc("/api/notes/{id}", UpdateNote).Methods("PUT")
	router.HandleFunc("/api/notes", DeleteAllNotes).Methods("DELETE")
	err := http.ListenAndServeTLS(":"+serverPort, "certs/server.crt", "certs/server.key", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func writeResponse(status int, body interface{}, w http.ResponseWriter) {
	fmt.Println("writeResponse")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	payload, _ := json.Marshal(body)
	w.Write(payload)
}

func GivenWellcome(w http.ResponseWriter, r *http.Request) {
	writeResponse(http.StatusOK, map[string]string{"message": "Wellcome to Notes API"}, w)
}

func ReadNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	note, err := backend.GetNote(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeResponse(http.StatusNotFound, map[string]string{"error": err.Error()}, w)
		} else {
			writeResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}, w)
		}
		return
	}
	writeResponse(http.StatusOK, note, w)
}

func ReadAllNote(w http.ResponseWriter, r *http.Request) {
	notes, err := backend.GetAllNotes()
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

	id, err := backend.SaveNote(&note)
	if err != nil {
		writeResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}, w)
		return
	}

	writeResponse(http.StatusOK, map[string]string{"id": id}, w)
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var note models.Note

	if err := decoder.Decode(&note); err != nil {
		writeResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}, w)
		return
	}

	note.Id = id
	id, err := backend.UpdateNote(&note)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeResponse(http.StatusNotFound, map[string]string{"error": err.Error()}, w)
		} else {
			writeResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}, w)
		}
		return
	}

	writeResponse(http.StatusOK, map[string]string{"id": id}, w)
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := backend.DeleteNote(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeResponse(http.StatusNotFound, map[string]string{"error": err.Error()}, w)
		} else {
			writeResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}, w)
		}
		return
	}
	writeResponse(http.StatusOK, map[string]string{"message": "Note deleted", "id": id}, w)
}

func DeleteAllNotes(w http.ResponseWriter, r *http.Request) {
	err := backend.DeleteAllNotes()
	if err != nil {
		writeResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}, w)
		return
	}
	writeResponse(http.StatusOK, map[string]string{"message": "All notes deleted"}, w)
}
