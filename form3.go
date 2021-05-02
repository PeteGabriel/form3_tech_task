package form3_task

import (
	"errors"
	"github.com/google/uuid"
	"github.com/petegabriel/form3_task/data"
	"log"
	"strconv"
)

//CreateAccount creates a new account with the given info.
//Returns an error if a problem occurs while trying to create the new account.
func CreateAccount(info *Account) (*Account, error){
	gate := data.NewGateway()
	dto := info.ToDto()
	if acc, err := gate.Create(dto); err != nil {
		log.Printf(err.Error())
		return nil, err
	}else {
		return NewAccountFromDto(acc), nil
	}
}

//DeleteAccount deletes the account with the given id.
//Given id must be a valid uuid type.
//Returns an error if a problem occurs while trying to delete the account with the given id.
func DeleteAccount(id string, vrs int) error {
	uid, isUuid := checkUuid(id)
	if !isUuid {
		invalidIdErr := errors.New("given id must be a valid uuid type")
		log.Print(invalidIdErr)
		return invalidIdErr
	}

	gate := data.NewGateway()
	err := gate.Delete(uid, strconv.Itoa(vrs))
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	return nil
}

//GetAccount retrieves an account by the given id.
//Given id must be a valid uuid type.
//Returns an error if a problem occurs while trying to delete the account with the given id.
func GetAccount(id string) (*Account, error){
	uid, isUuid := checkUuid(id)
	if !isUuid {
		invalidIdErr := errors.New("given id must be a valid uuid type")
		log.Print(invalidIdErr)
		return nil, invalidIdErr
	}

	gate := data.NewGateway()
	if found, err := gate.Get(uid); err != nil {
		log.Print(err)
		return nil, err
	}else {
		return NewAccountFromDto(found), nil
	}
}

func checkUuid(id string) (uuid.UUID, bool) {
	if uid, err := uuid.Parse(id); err != nil {
		return uuid.New(), false
	}else {
		return uid, true
	}
}
