package backend_test

import (
	"reflect"
	"testing"

	"com.notes/notes/internal/backend"
	"com.notes/notes/internal/db_postgres"
	"com.notes/notes/internal/models"
)

func TestGetNoteOk(t *testing.T) {
	// Given
	db_postgres.GetNote = func(key string) (*models.Note, error) {
		return &models.Note{
			Title:        "OK",
			Content:      "OK",
			Id:           "1",
			ReminderDate: "2023-09-29T10:00:00Z",
		}, nil
	}

	// When
	note, err := backend.GetNote("1")

	// Then
	if err != nil {
		t.Errorf("GetNote() error = %v", err)
		return
	}
	if note == nil || note.Title == "" || note.Content == "" || note.Id == "" || note.ReminderDate == "" {
		t.Errorf("GetNote() = %v", note)
	}
	if reflect.TypeOf(note) != reflect.TypeOf(&models.Note{}) {
		t.Errorf("GetNote() = %v", note)
	}
}
