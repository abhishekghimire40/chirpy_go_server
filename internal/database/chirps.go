package database

import "log"

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {

	chirpsData, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}
	newChirp := Chirp{
		Id:   len(chirpsData.Chirps) + 1,
		Body: body,
	}
	// add the new chirp to the database
	chirpsData.Chirps[newChirp.Id] = newChirp
	// save updated chirps to our database
	err = db.writeDB(chirpsData)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {

	data, err := db.loadDB()

	if err != nil {
		log.Fatal(err)
	}
	chirpData := make([]Chirp, 0)

	for _, val := range data.Chirps {
		newChirp := Chirp{
			Id:   val.Id,
			Body: val.Body,
		}
		chirpData = append(chirpData, newChirp)
	}

	return chirpData, nil
}
