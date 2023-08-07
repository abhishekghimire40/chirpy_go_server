package database

import (
	"errors"
	"time"
)

func (db *DB) RevokeToken(token string) (string, error) {
	data, err := db.loadDB()
	if err != nil {
		return "", err
	}
	if _, ok := db.IsRevoked(token); ok {
		return "", errors.New("token already revoked")
	}
	data.Revoked_Token[token] = time.Now()
	db.writeDB(data)
	return token, nil
}

func (db *DB) IsRevoked(token string) (string, bool) {
	data, err := db.loadDB()
	if err != nil {
		return "", true
	}
	for key := range data.Revoked_Token {
		if key == token {
			return key, true
		}
	}
	return token, false
}
