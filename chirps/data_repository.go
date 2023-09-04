package chirps

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type DataSchema struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type DataRepository struct {
	path string
	mux  *sync.Mutex
}

func NewChirpRepository() (*DataRepository, error) {
	db := DataRepository{path: "database.json", mux: &sync.Mutex{}}
	err := db.init()
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func (db *DataRepository) init() error {
	_, err := os.Stat(db.path)
	if err != nil {
		log.Println("DB file not found, creating DB file")
		_, err := os.Create(db.path)
		if err != nil {
			log.Printf("Error creating DB file: %s", err)
			return err
		}
		db.write(DataSchema{Chirps: map[int]Chirp{}})
		log.Println("DB file created")
	}
	return nil
}

func (db *DataRepository) read() (DataSchema, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	rawData, err := os.ReadFile(db.path)
	if err != nil {
		log.Printf("Error Reading file: %s", err)
		return DataSchema{}, err
	}

	var data DataSchema
	if err := json.Unmarshal(rawData, &data); err != nil {
		log.Printf("Error decoding JSON: %s", err)
		return DataSchema{}, err
	}

	return data, nil
}

func (db *DataRepository) write(data DataSchema) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	rawData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data: %s", err)
		return err
	}

	if err := os.WriteFile(db.path, rawData, 0600); err != nil {
		log.Printf("Error writing to file: %s", err)
		return err
	}
	return nil
}
