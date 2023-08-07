package database

import (
	"errors"
)

type Chirp struct {
	Id        int    `json:"id"`
	Body      string `json:"body"`
	Author_Id int    `json:"author_id"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, user_id int) (Chirp, error) {

	chirpsData, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}
	newChirp := Chirp{
		Id:        len(chirpsData.Chirps) + 1,
		Body:      body,
		Author_Id: user_id,
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
		return nil, err
	}
	chirpData := make([]Chirp, 0)

	for _, val := range data.Chirps {
		newChirp := Chirp{
			Id:        val.Id,
			Body:      val.Body,
			Author_Id: val.Author_Id,
		}
		chirpData = append(chirpData, newChirp)
	}

	return chirpData, nil
}

// function to return get chirps by author_id
func (db *DB) GetChirpsByID(user_id int) ([]Chirp, error) {
	data, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirpsData := make([]Chirp, 0)
	for _, val := range data.Chirps {
		if val.Author_Id == user_id {
			chirpsData = append(chirpsData, val)
		}
	}
	return chirpsData, nil
}

// method to delete chirps
func (db *DB) DeleteChirp(chirp_id int, user_id int) (int, error) {
	data, err := db.loadDB()
	if err != nil {
		return 500, errors.New("internal server error")
	}
	chirp_data, ok := data.Chirps[chirp_id]
	if !ok {
		return 403, errors.New("chirp with provided id not available")
	}
	if chirp_data.Author_Id != user_id {
		return 403, errors.New("couldn't delete chirp as you are not owner of the chirp")
	}
	delete(data.Chirps, chirp_id)
	db.writeDB(data)
	return 200, nil
}
