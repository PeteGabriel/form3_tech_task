package data

import "github.com/google/uuid"

type AccountApiGateway interface {

	//Create a new account
	Create(AccountDto) (AccountDto, error)

	//Delete an account by id
	Delete(uid uuid.UUID) error

	//Get an account by id
	Get(id uuid.UUID) AccountDto
}
