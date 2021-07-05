package main

import (
	"errors"
	"math/rand"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}

type IAuthService interface {
	SignIn(username, password string) (*User, error)
	SignUp(user *User) (*User, error)
}


type AuthService struct {

}

func (a *AuthService) SignIn(username, password string) (*User, error) {

	//logic ....
	if false {
		return nil, errors.New("username or password is incorrect")
	}

	return &User{
		Id:       rand.Int() % 100,
		Username: username,
		Password: password,
	}, nil
}

func (a *AuthService) SignUp(user *User) (*User, error) {

	//logic ....
	if false {
		return nil, errors.New("username or password is incorrect")
	}

	return user, nil
}

