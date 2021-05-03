package data

import (
	"fmt"
	"github.com/google/uuid"
	is2 "github.com/matryer/is"
	"testing"
)

func TestGet(t *testing.T){
	is := is2.New(t)
	dto := setupNewAccount([]string{"Peter", "Devos"})

	gate := NewGateway()
	id, _ := uuid.Parse(dto.Data.ID)
	accFound, err := gate.Get(id)

	is.NoErr(err)
	assertAccounts(is, accFound, dto)

	resetState(id)
}

func TestGetNotFoundID(t *testing.T){
	is := is2.New(t)
	gate := NewGateway()
	var uid uuid.UUID

	acc, err := gate.Get(uid)
	is.True(err != nil)
	is.Equal(err.Error(), fmt.Sprintf("account with uid %s not found", uid.String()))
	is.Equal(acc, AccountDto{})
}

func TestDelete(t *testing.T) {
	is := is2.New(t)
	gate := NewGateway()
	dto := setupNewAccount([]string{"Kim", "Emma"})
	uid, _ := uuid.Parse(dto.Data.ID)
	err := gate.Delete(uid, "0")
	is.NoErr(err)
}

func TestDeleteNotFoundID(t *testing.T) {
	is := is2.New(t)
	gate := NewGateway()
	var uid uuid.UUID

	err := gate.Delete(uid, "0")
	is.True(err != nil)
	is.Equal(err.Error(), fmt.Sprintf("account with uuid %s not found", uid.String()))
}

func TestDeleteConflictID(t *testing.T) {
	is := is2.New(t)
	gate := NewGateway()
	dto := setupNewAccount([]string{"Kim", "Emma"})
	uid, _ := uuid.Parse(dto.Data.ID)
	err := gate.Delete(uid, "1")
	is.True(err != nil)
	is.Equal(err.Error(), "account with specified version not found")

	resetState(uid)
}

func TestCreate(t *testing.T) {
	is := is2.New(t)
	dto, _ := newAccount([]string{"Martin", "Fuchs"})
	gate := NewGateway()
	created, err := gate.Create(dto)
	is.NoErr(err)
	assertAccounts(is, created, dto)

	id, _ := uuid.Parse(created.Data.ID)
	resetState(id)
}

func TestCreateWithConflict(t *testing.T) {
	is := is2.New(t)
	gate := NewGateway()

	dto, _ := newAccount([]string{"Kim", "Emma", "First"})
	_, err := gate.Create(dto)
	is.NoErr(err)
	_, err = gate.Create(dto)

	is.True(err != nil)
	//reuse the error message from account api, assert that the string is not empty. Content may vary
	is.True(len(err.Error()) > 0)

	id, _ := uuid.Parse(dto.Data.ID)
	resetState(id)
}

func TestCreateWithBadRequest(t *testing.T)  {
	is := is2.New(t)
	gate := NewGateway()
	//one example is empty line of name
	dto, _ := newAccount([]string{"Kim", ""})
	_, err := gate.Create(dto)
	is.True(err != nil)
	//reuse the error message from account api, assert that the string is not empty. Content may vary
	is.True(len(err.Error()) > 0)
}


func setupNewAccount(name []string) AccountDto{
	gate := NewGateway()

	dto, err := newAccount(name)

	acc, err := gate.Create(dto)
	if err != nil {
		panic(err)
	}
	return acc
}

func newAccount(name []string) (AccountDto, error) {
	//create new account for testing purposes
	idAcc, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	idOrg, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	dto := NewAccountDto(idAcc, idOrg, "GB", name)
	return dto, err
}

func resetState(uid uuid.UUID){
	gate := NewGateway()
	//reset state previously to testing
	err := gate.Delete(uid, "0")
	if err != nil {
		panic(err)
	}
}

func assertAccounts(is *is2.I, created AccountDto, dto AccountDto) {
	is.True(created.Data.Type == dto.Data.Type)
	is.True(created.Data.ID == dto.Data.ID)
	is.True(created.Data.OrganisationID == dto.Data.OrganisationID)
	is.True(created.Data.Version == dto.Data.Version)
	is.True(created.Data.Attributes.Country == dto.Data.Attributes.Country)
	is.True(created.Data.Attributes.BaseCurrency == dto.Data.Attributes.BaseCurrency)
	is.True(created.Data.Attributes.AccountNumber == dto.Data.Attributes.AccountNumber)
	is.True(created.Data.Attributes.BankID == dto.Data.Attributes.BankID)
	is.True(created.Data.Attributes.BankIDCode == dto.Data.Attributes.BankIDCode)
	is.True(created.Data.Attributes.Bic == dto.Data.Attributes.Bic)
	is.True(created.Data.Attributes.Iban == dto.Data.Attributes.Iban)
	is.Equal(created.Data.Attributes.Name, dto.Data.Attributes.Name)
	is.Equal(created.Data.Attributes.AlternativeNames, nil)
	is.True(created.Data.Attributes.AccountClassification == "Personal")
	is.Equal(created.Data.Attributes.SecondaryIdentification, dto.Data.Attributes.SecondaryIdentification)
	is.Equal(created.Data.Attributes.Switched, dto.Data.Attributes.Switched)
}

func TestParseErrorMsg(t *testing.T){
	is := is2.New(t)
	err := "validation failure list:\nvalidation failure list:\nvalidation failure list:\nname.1 in body should be at least 1 chars long\n"
	expected :=  "name.1 in body should be at least 1 chars long"
	msg := parseErrorMsg(err)
	is.Equal(msg, expected)
}