package db

import (
	"context"
	"encoding/json"

	"errors"

	"com.notes/notes/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var GetNote = func(key string) (*models.Note, error) {
	db := GetDataBase()
	val, err := db.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.New("key not found")
	} else if err != nil {
		return nil, err
	}
	note := &models.Note{}
	err = json.Unmarshal([]byte(val), note)
	if err != nil {
		return nil, err
	}
	return note, nil
}

var GetAllNote = func() ([]*models.Note, error) {
	db := GetDataBase()
	keys, err := db.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}
	notes := []*models.Note{}
	for _, key := range keys {
		note, err := GetNote(key)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

var SaveNote = func(note *models.Note) (string, error) {
	uuidKey := uuid.New().String()
	db := GetDataBase()
	note.Id = uuidKey
	json, err := json.Marshal(note)
	if err != nil {
		return "", err
	}
	err = db.Set(ctx, uuidKey, json, 0).Err()
	if err != nil {
		return "", err
	}
	return uuidKey, nil
}

var DeleteNote = func(key string) error {
	db := GetDataBase()
	err := db.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

var (
	DataBaseUrl      string
	DataBasePassword string
	redisClient      *redis.Client
	ctx              = context.Background()
)

var GetDataBase = func() *redis.Client {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     DataBaseUrl,
			Password: DataBasePassword,
			DB:       0,
		})
	}
	return redisClient
}
