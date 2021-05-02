package form3_task

import (
	"fmt"
	"github.com/google/uuid"
	is2 "github.com/matryer/is"
	"testing"
)

func TestGetAccount(t *testing.T) {
	is := is2.New(t)
	dto := NewAccount([]string{"Hugh", "Grant"}, "GB", getRandomId(), getRandomId())
	_, _ = CreateAccount(dto)

	acc, err := GetAccount(dto.Id.String())
	is.NoErr(err)
	assertAccountData(is, acc, dto)

	//clean up data
	_ = DeleteAccount(acc.Id.String(), acc.Version)
}

func TestGetAccountUuidNotFound(t *testing.T) {
	is := is2.New(t)

	id := getRandomId().String()
	acc, err := GetAccount(id)
	is.True(err != nil)
	is.Equal(err.Error(), fmt.Sprintf("account with uid %s not found", id))
	is.True(acc == nil)
}

func TestGetWithInvalidUUID(t *testing.T) {
	is := is2.New(t)
	id := "c1-70-41-9a-e21"
	acc, err := GetAccount(id)
	is.True(err != nil)
	is.True(acc == nil)
	is.Equal(err.Error(), "given id must be a valid uuid type")
}

func TestCreateAccount(t *testing.T) {
	is := is2.New(t)

	country := "GB"
	name := []string{"Samantha Holder"}
	id := getRandomId()
	orgId := getRandomId()
	dto := NewAccount(name, country, id, orgId)
	dto.BaseCurrency = "GBP"
	dto.AccountNumber = "41426819"
	dto.BankId = "400300"
	dto.BankIdCode = "GBDSC"
	dto.Bic = "NWBKGB22"
	dto.Iban = "GB11NWBK40030041426819"
	dto.AlternativeNames = []string{"Sam Holder"}
	dto.Classification = Personal
	dto.SecondaryIdentification = "A1B2C3D4"

	acc, err := CreateAccount(dto)

	is.NoErr(err)
	assertAccountData(is, acc, dto)
	//clean up data
	_ = DeleteAccount(acc.Id.String(), acc.Version)
}

func TestCreateAccountWithMinimumInfo(t *testing.T) {
	is := is2.New(t)
	country := "PT"
	name := []string{"Pedro", "Almeida"}
	id := getRandomId()
	orgId := getRandomId()
	dto := NewAccount(name, country, id, orgId)

	acc, err := CreateAccount(dto)
	is.NoErr(err)

	assertAccountData(is, acc, dto)

	//clean up data
	err = DeleteAccount(acc.Id.String(), acc.Version)
	is.NoErr(err)
}

func TestCreateAccountConflict(t *testing.T) {
	is := is2.New(t)
	country := "PT"
	name := []string{"Pedro", "Almeida"}
	id := getRandomId()
	orgId := getRandomId()
	dto := NewAccount(name, country, id, orgId)

	_, err := CreateAccount(dto)
	is.NoErr(err)

	//same id
	dto2 := NewAccount(name, country, id, getRandomId())
	_, err = CreateAccount(dto2)
	is.True(err != nil)
	is.True(len(err.Error()) > 0)

	//clean up data
	err = DeleteAccount(id.String(), 0)
	is.NoErr(err)
}


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

func TestDeleteWithInvalidUUID(t *testing.T) {
	is := is2.New(t)
	id := "c1-70-41-9a-e21"
	err := DeleteAccount(id, 0)
	is.True(err != nil)
	is.Equal(err.Error(), "given id must be a valid uuid type")
}

func TestDeleteWithNonexistentUUID(t *testing.T) {
	is := is2.New(t)
	id := getRandomId().String()
	err := DeleteAccount(id, 0)
	is.True(err != nil)
	is.Equal(err.Error(), fmt.Sprintf("account with uuid %s not found", id))
}

func TestDeleteWithInvalidVersion(t *testing.T) {
	is := is2.New(t)

	country := "GB"
	name := []string{"Peter Devos"}
	id := getRandomId()
	orgId := getRandomId()
	dto := NewAccount(name, country, id, orgId)

	_, err := CreateAccount(dto)

	err = DeleteAccount(id.String(), 10)
	is.True(err != nil)
	is.Equal(err.Error(), "account with specified version not found")

	//clean up data
	_ = DeleteAccount(id.String(), 0)
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
	is.Equal(accInfo.AlternativeNames, nil)
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
	is.True(len(acc.CreatedOn)>0)
	is.True(len(acc.ModifiedOn)>0)
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