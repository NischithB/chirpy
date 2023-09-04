package chirps

import (
	"errors"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

var ErrNotFound = errors.New("error: resource not found")

func (repo *DataRepository) CreateChirp(body string) (Chirp, error) {
	data, err := repo.read()
	if err != nil {
		return Chirp{}, err
	}

	id := len(data.Chirps) + 1
	chirp := Chirp{Id: id, Body: body}
	data.Chirps[id] = chirp
	if err := repo.write(data); err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

func (repo *DataRepository) GetChirps() ([]Chirp, error) {
	data, err := repo.read()
	if err != nil {
		return []Chirp{}, err
	}

	chirps := []Chirp{}
	for _, chirp := range data.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (repo *DataRepository) GetChirpById(chirpID int) (Chirp, error) {
	data, err := repo.read()
	if err != nil {
		return Chirp{}, err
	}

	chirp, exists := data.Chirps[chirpID]
	if !exists {
		return Chirp{}, ErrNotFound
	}
	return chirp, nil
}
