package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	newDB := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	newDB.ensureDB()
	return newDB, nil
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

// function to create a database
func (db *DB) createDB() error {
	dbStructure := DBStructure{Chirps: make(map[int]Chirp)}
	return db.writeDB(dbStructure)
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if errors.Is(err, os.ErrNotExist) {
		db.createDB()
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	// read json file
	db.mux.RLock()
	defer db.mux.RUnlock()
	response := DBStructure{}
	jsonData, err := os.ReadFile(db.path)
	if err != nil {
		return response, err
	}
	// check if the database is empty
	if len(jsonData) == 0 {
		return response, nil
	}
	// parse jsonfile

	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

// writeDB writes the database file into disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	updatedData, err := json.MarshalIndent(dbStructure, "", "\t")
	if err != nil {
		return err
	}
	err1 := os.WriteFile(db.path, updatedData, 0644)
	if err1 != nil {
		return err
	}
	return nil
}
