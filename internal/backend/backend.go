package backend

import (
	"errors"
	"fmt"

	db_postgres "com.notes/notes/internal/db_postgres"
	"com.notes/notes/internal/models"
)

func GetNote(key string) (*models.Note, error) {
	if len(key) == 0 || len(key) > 36 {
		return nil, errors.New("invalid key size")
	}
	note, err := db_postgres.GetNote(key)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func GetAllNotes() ([]*models.Note, error) {
	notes, err := db_postgres.GetAllNotes()
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func SaveNote(note *models.Note) (string, error) {
	fmt.Println("SaveNote")
	byteSize := len([]byte(note.Title)) + len([]byte(note.Content))
	if byteSize == 0 || byteSize > (64*1024) {
		return "", errors.New("invalid data sizeee")
	}
	key, err := db_postgres.SaveNote(note)
	if err != nil {
		return "", err
	}
	return key, nil
}

func UpdateNote(note *models.Note) (string, error) {
	fmt.Println("UpdateNote")
	byteSize := len([]byte(note.Title)) + len([]byte(note.Content))
	if byteSize == 0 || byteSize > (64*1024) {
		return "", errors.New("invalid data sizeee")
	}
	id, err := db_postgres.UpdateNote(note)
	fmt.Println("backend")
	fmt.Println("-------\n", id, "\n-------")
	if err != nil {
		return "", err
	}
	return id, nil
}

func DeleteNote(key string) error {
	fmt.Println("DeleteNote")
	err := db_postgres.DeleteNote(key)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllNotes() error {
	fmt.Println("DeleteAllNote")
	err := db_postgres.DeleteAllNotes()
	if err != nil {
		return err
	}
	return nil
}
