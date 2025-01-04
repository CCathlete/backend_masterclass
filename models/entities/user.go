package entities

import "fmt"

type User struct {
	Username     string
	PasswordHash string
	Email        string
}

type UserRepo interface {
	CreateUser(user User) (User, error)
	GetUser(username string) (User, error)
}

type CreateUserParams struct {
	Username, PAsswordHash, Email string
}

type UserRepoStub struct {
	Users map[string]User
}

func (r *UserRepoStub) CreateUser(params CreateUserParams,
) (user User, err error) {

	user = User{
		Username:     params.Username,
		PasswordHash: params.PAsswordHash,
		Email:        params.Email,
	}
	r.Users[params.Username] = user

	return
}

func (r *UserRepoStub) GetUser(username string) (user User, err error) {

	user, ok := r.Users[username]
	if !ok {
		err = fmt.Errorf("user (%s) not found", username)
	}

	return
}
