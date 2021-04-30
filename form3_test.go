package form3_task

import (
	"github.com/google/uuid"
	is2 "github.com/matryer/is"
	"testing"
)


func TestCreateAccountWithInvalidParams(t *testing.T) {
	is := is2.New(t)
	country := "PT"
	//invalid name
	name := []string{"Pedro", "", "Almeida"}
	id := getRandomId()
	orgId := getRandomId()
	dto := NewAccount(name, country, id, orgId)
	acc, err := CreateAccount(dto)
	is.True(err != nil)
	is.True(len(err.Error()) > 0)
	is.Equal(acc, nil)
}

func TestCreateAccount(t *testing.T) {
	is := is2.New(t)
	country := "PT"
	name := []string{"Pedro", "Almeida"}
	id := getRandomId()
	orgId := getRandomId()
	dto := NewAccount(name, country, id, orgId)
	acc, err := CreateAccount(dto)
	is.NoErr(err)

	assertAccountData(is, acc, dto)

	//TODO clean up data

}

func TestNewAccount(t *testing.T) {
	is := is2.New(t)
	country := "PT"
	name := []string{"Pedro", "Almeida"}
	id := getRandomId()
	orgId := getRandomId()
	accInfo := NewAccount(name, country, id, orgId)

	is.Equal(accInfo.Id, id)
	is.Equal(accInfo.OrganisationId, orgId)
	is.Equal(accInfo.Country, country)
	is.Equal(accInfo.BaseCurrency, "")
	is.Equal(accInfo.AccountNumber, "")
	is.Equal(accInfo.BankId, "")
	is.Equal(accInfo.BankIdCode, "")
	is.Equal(accInfo.Bic, "")
	is.Equal(accInfo.Iban, "")
	is.Equal(accInfo.Name, name)
	is.Equal(accInfo.AlternativeNames, [3]string{})
	is.Equal(accInfo.Classification, Personal)
	is.Equal(accInfo.IsJointAccount, false)
	is.Equal(accInfo.IsAccountMatchingOptOut, false)
	is.Equal(accInfo.SecondaryIdentification, "")
	is.Equal(accInfo.IsSwitched, false)
}

func getRandomId() uuid.UUID {
	id, err := uuid.NewRandom()
	if err != nil {
		panic("error generating random id for testing purposes")
	}
	return id
}

func assertAccountData(is *is2.I, acc *Account, dto *Account) {
	is.Equal(acc.Id, dto.Id)
	is.Equal(acc.OrganisationId, dto.OrganisationId)
	is.Equal(acc.Country, dto.Country)
	is.Equal(acc.BaseCurrency, dto.BaseCurrency)
	is.Equal(acc.AccountNumber, dto.AccountNumber)
	is.Equal(acc.BankId, dto.BankId)
	is.Equal(acc.BankIdCode, dto.BankIdCode)
	is.Equal(acc.Bic, dto.Bic)
	is.Equal(acc.Iban, dto.Iban)
	is.Equal(acc.Name, dto.Name)
	is.Equal(acc.AlternativeNames, dto.AlternativeNames)
	is.Equal(acc.Classification, dto.Classification)
	is.Equal(acc.IsJointAccount, dto.IsJointAccount)
	is.Equal(acc.IsAccountMatchingOptOut, dto.IsAccountMatchingOptOut)
	is.Equal(acc.SecondaryIdentification, dto.SecondaryIdentification)
	is.Equal(acc.IsSwitched, dto.IsSwitched)
}