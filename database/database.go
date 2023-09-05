package database

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Database[T any] struct {
	Path string
	Mux  *sync.Mutex
}

func NewDatabase[T any](path string, initialData T) (*Database[T], error) {
	db := Database[T]{Path: path, Mux: &sync.Mutex{}}
	err := db.init(initialData)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

func (db *Database[T]) init(initialData T) error {
	_, err := os.Stat(db.Path)
	if err != nil {
		log.Println("DB file not found, creating DB file")
		_, err := os.Create(db.Path)
		if err != nil {
			log.Printf("Error creating DB file: %s", err)
			return err
		}
		db.Write(initialData)
		log.Println("DB file created")
	}
	return nil
}

func (db *Database[T]) Read() (T, error) {
	db.Mux.Lock()
	defer db.Mux.Unlock()
	var data T

	rawData, err := os.ReadFile(db.Path)
	if err != nil {
		log.Printf("Error Reading file: %s", err)
		return data, err
	}

	if err := json.Unmarshal(rawData, &data); err != nil {
		log.Printf("Error decoding JSON: %s", err)
		return data, err
	}

	return data, nil
}

func (db *Database[T]) Write(data T) error {
	db.Mux.Lock()
	defer db.Mux.Unlock()

	rawData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data: %s", err)
		return err
	}

	if err := os.WriteFile(db.Path, rawData, 0600); err != nil {
		log.Printf("Error writing to file: %s", err)
		return err
	}
	return nil
}
