package data

import (
	"github.com/google/uuid"
	is2 "github.com/matryer/is"
	"os"
	"testing"
)

var uidAcc, orgId uuid.UUID


func TestGet(t *testing.T){
	is := is2.New(t)
	gate := new(Gateway)

	acc, err := gate.Get(uidAcc)
	is.NoErr(err)
	is.True(acc.Data.ID == uidAcc.String())
	is.True(acc.Data.OrganisationID == orgId.String())
	//by default, Personal
	is.True(acc.Data.Attributes.AccountClassification == "Personal")
	is.Equal(acc.Data.Attributes.Name, []string{"Peter", "Devos"})
}

func TestDelete(t *testing.T) {
	is := is2.New(t)
	gate := new(Gateway)
	dto, _ := newAccount([]string{"Kim", "Emma"})
	acc, _ := gate.Create(dto)
	uid, _ := uuid.Parse(acc.Data.ID)
	err := gate.Delete(uid)
	is.NoErr(err)
}

func TestCreate(t *testing.T) {
	is := is2.New(t)
	gate := new(Gateway)
	dto, _ := newAccount([]string{"Martin", "Fuchs"})
	acc, err := gate.Create(dto)
	is.NoErr(err)
	is.True(acc.Data.ID == dto.Data.ID)
	is.True(acc.Data.Attributes.AccountClassification == "Personal")
	is.Equal(acc.Data.Attributes.Name, []string{"Martin", "Fuchs"})
}

func TestGetWithInvalidUID(t *testing.T){
	is := is2.New(t)
	gate := new(Gateway)
	var uid uuid.UUID

	acc, err := gate.Get(uid)
	is.True(err != nil)
	is.Equal(acc, AccountDto{})
}

func TestDeleteWithInvalidUID(t *testing.T) {
	is := is2.New(t)
	gate := new(Gateway)
	var uid uuid.UUID

	err := gate.Delete(uid)
	is.True(err != nil)
	is.True(len(err.Error()) > 0)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup(){
	gate := new(Gateway)

	dto, err := newAccount([]string{"Peter", "Devos"})

	acc, err := gate.Create(dto)
	if err != nil {
		panic(err)
	}
	uidAcc, err = uuid.Parse(acc.Data.ID)
	if err != nil {
		panic(err)
	}
	orgId, err = uuid.Parse(acc.Data.OrganisationID)
	if err != nil {
		panic(err)
	}
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
	dto := NewAccountDto(idAcc, idOrg, "GB", "GBDSC")
	dto.Data.Attributes.Name = name
	return dto, err
}

func shutdown(){
	gate := new(Gateway)
	//reset state previously to testing
	err := gate.Delete(uidAcc)
	if err != nil {
		panic(err)
	}
}
