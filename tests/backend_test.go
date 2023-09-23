package backend_test

import (
	"testing"

	"com.notes/notes/internal/backend"
	"com.notes/notes/internal/db"
	"com.notes/notes/internal/models"
)

func TestGetNoteOk(t *testing.T) {
	// Given
	db.GetNote = func(key string) (*models.Note, error) {
		return &models.Note{
			Data: "OK",
			Once: false,
		}, nil
	}

	// When
	note, err := backend.GetNote("key")

	// Then
	if err != nil {
		t.Errorf("GetNote() error = %v", err)
		return
	}
	if note.Data != "OK" {
		t.Errorf("GetNote() data = %v", note.Data)
	}
}
