package router_test

import "testing"

type UserService struct {
	Id   int
	Name string
}

type IUserService interface {
	GetId() int
	GetName() string
}

func TestSetService(t *testing.T) {

}

func TestGetService(t *testing.T) {

}
