package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp `json:"chirps"`
	Users         map[int]User  `json:"users"`
	Revoked_Token map[string]time.Time
}

func NewDB(path string) (*DB, error) {
	newDB := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	newDB.ensureDB()
	return newDB, nil
}

// function to create a database
func (db *DB) createDB() error {
	dbStructure := DBStructure{Chirps: make(map[int]Chirp), Users: make(map[int]User), Revoked_Token: make(map[string]time.Time)}
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
