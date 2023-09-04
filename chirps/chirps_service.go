package chirps

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

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

	var chirps []Chirp
	for _, chirp := range data.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}
