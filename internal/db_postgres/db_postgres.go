package db_postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"com.notes/notes/internal/models"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectWithDB() *sql.DB {
	godotenv.Load()
	dbPawword := os.Getenv("POSTGRES_PASSWORD")
	connectionString := fmt.Sprintf("user=postgres dbname=notes_api password=%s host=localhost sslmode=disable", dbPawword)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Close(db *sql.DB) {
	db.Close()
}

var GetNote = func(key string) (*models.Note, error) {
	db := ConnectWithDB()
	defer Close(db)
	var note models.Note
	err := db.QueryRow("SELECT id, title, content, reminderdate FROM Notas WHERE id=$1", key).Scan(&note.Id, &note.Title, &note.Content, &note.ReminderDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("id not found")
		}

		return nil, err
	}
	return &note, nil
}

var GetAllNotes = func() ([]*models.Note, error) {

	db := ConnectWithDB()
	defer Close(db)
	rows, err := db.Query("SELECT * FROM Notas")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("id not found")
		}

		return nil, err
	}

	notes := []*models.Note{}

	for rows.Next() {
		var note models.Note
		err := rows.Scan(&note.Id, &note.Title, &note.Content, &note.ReminderDate)
		if err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	return notes, nil
}

var SaveNote = func(note *models.Note) (string, error) {
	db := ConnectWithDB()
	defer Close(db)

	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	uuidKey := uuid.New().String()
	note.Id = uuidKey

	_, err = tx.Exec("INSERT INTO Notas (id, title, content, reminderdate) VALUES ($1, $2, $3, $4)", note.Id, note.Title, note.Content, note.ReminderDate)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return note.Id, nil
}

var UpdateNote = func(note *models.Note) (string, error) {
	db := ConnectWithDB()
	defer Close(db)

	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	result, err := tx.Exec("UPDATE Notas SET title=$1, content=$2, reminderdate=$3 WHERE id=$4", note.Title, note.Content, note.ReminderDate, note.Id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return "", fmt.Errorf("id not found")
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return note.Id, nil
}

var DeleteNote = func(id string) error {
	db := ConnectWithDB()
	defer Close(db)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	result, err := tx.Exec("DELETE FROM Notas WHERE id=$1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("id not found")
	}

	return tx.Commit()
}

var DeleteAllNotes = func() error {
	db := ConnectWithDB()
	defer Close(db)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM Notas")
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
