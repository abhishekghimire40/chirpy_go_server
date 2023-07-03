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
	return Chirp{
		Id:   len(chirpsData.Chirps) + 1,
		Body: body,
	}, nil
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

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(db.path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	// read json file
	jsonData, err := os.ReadFile(db.path)
	if err != nil {
		log.Fatal(err)
	}
	// check if the database is empty
	if len(jsonData) == 0 {

		return DBStructure{Chirps: make(map[int]Chirp)}, nil
	}
	// parse jsonfile
	response := DBStructure{}
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		log.Fatal(err)
	}

	return response, nil
}

// writeDB writes the database file into disk
func (db *DB) writeDB(dbStructure DBStructure) error {
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
