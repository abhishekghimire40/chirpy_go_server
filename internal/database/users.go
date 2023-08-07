package database

import (
	"errors"

	"github.com/abhishekghimire40/chirpy_go_server/internal/auth"
)

type User struct {
	Id            int    `json:"id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Is_Chirpy_Red bool   `json:"is_chirpy_red"`
}

type PublicUser struct {
	Id            int    `json:"id"`
	Email         string `json:"email"`
	Is_Chirpy_Red bool   `json:"is_chirpy_red"`
}

// CreateUser creates a new user and saves it to our database
func (db *DB) CreateUser(email string, password string) (PublicUser, error) {
	// load all user data from database.json file
	userData, err := db.loadDB()
	if err != nil {
		return PublicUser{}, err
	}
	// check if user with provided email already exists
	_, exist := db.GetUser(email)
	if exist {
		return PublicUser{}, errors.New("User with provided email already exists")
	}

	// hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return PublicUser{}, err
	}

	// create new user
	newUser := User{
		Id:            len(userData.Users) + 1,
		Email:         email,
		Password:      hashedPassword,
		Is_Chirpy_Red: false,
	}
	userData.Users[newUser.Id] = newUser
	err = db.writeDB(userData)
	if err != nil {
		return PublicUser{}, err
	}
	return PublicUser{
		Id:            newUser.Id,
		Email:         newUser.Email,
		Is_Chirpy_Red: newUser.Is_Chirpy_Red,
	}, nil
}

// function to update user info: only password
func (db *DB) UpdateUser(id int, email string, password string) (PublicUser, error) {
	usersData, err := db.loadDB()
	if err != nil {
		return PublicUser{}, err
	}
	user, ok := usersData.Users[id]
	if !ok {
		return PublicUser{}, errors.New("internal server error")
	}
	user.Password = password
	user.Email = email
	usersData.Users[id] = user
	db.writeDB(usersData)
	return PublicUser{
		Id:            user.Id,
		Email:         user.Email,
		Is_Chirpy_Red: user.Is_Chirpy_Red,
	}, nil
}

// function to get a user info by its email
func (db *DB) GetUser(email string) (User, bool) {
	data, err := db.loadDB()
	if err != nil {
		return User{}, false
	}
	for _, val := range data.Users {
		if val.Email == email {
			return val, true
		}

	}
	return User{}, false
}

func (db *DB) UpgradeUser(user_id int) error {
	data, err := db.loadDB()
	if err != nil {
		return err
	}
	user, ok := data.Users[user_id]
	if !ok {
		return errors.New("user not found")
	}
	user.Is_Chirpy_Red = true
	data.Users[user_id] = user
	db.writeDB(data)
	return nil
}
