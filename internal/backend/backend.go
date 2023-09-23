package backend

import (
	"errors"
	"fmt"

	"com.notes/notes/internal/db"
	"com.notes/notes/internal/models"
)

func GetNote(key string) (*models.Note, error) {
	if len(key) == 0 || len(key) > 36 {
		return nil, errors.New("invalid key size")
	}
	note, err := db.GetNote(key)
	if err != nil {
		return nil, err
	}
	if note.Once {
		err := db.DeleteNote(key)
		if err != nil {
			return nil, err
		}
	}
	return note, nil
}

func GetAllNote() ([]*models.Note, error) {
	notes, err := db.GetAllNote()
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func SaveNote(note *models.Note) (string, error) {
	fmt.Println("SaveNote")
	fmt.Println(note)
	byteSize := len([]byte(note.Data))
	if byteSize == 0 || byteSize > (32*1024) {
		return "", errors.New("invalid data sizeee")
	}
	key, err := db.SaveNote(note)
	if err != nil {
		return "", err
	}
	return key, nil
}

func DeleteNote(key string) error {
	fmt.Println("DeleteNote")
	err := db.DeleteNote(key)
	if err != nil {
		return err
	}
	return nil
}
