package form3_task

import (
	"github.com/google/uuid"
	"github.com/petegabriel/form3_task/data"
)

//Classification of account. Can be one of 'Personal' or 'Business'.
type Classification string

const (
	Personal Classification = "Personal"
	Business Classification = "Business"
)

//Account represents a bank account that is registered with Form3 fake account api.
type Account struct {

	//Id of the account.
	Id uuid.UUID

	//OrganisationId of the account.
	OrganisationId uuid.UUID

	//Version number
	Version int

	//Country is an ISO code used to identify the domicile of the account.
	Country string

	//BaseCurrency is an ISO code used to identify the base currency of the account.
	BaseCurrency string

	//AccountNumber identifies uniquely the account. Generated if not provided.
	AccountNumber string

	//BankId is the local country bank identifier.
	BankId string

	//BankIdCode identifies the type of bank ID being used
	BankIdCode string

	//Bic is a swift bic code in either 8 or 11 character format
	Bic string

	//Iban of the account. Will be calculated from other fields if not supplied.
	Iban string

	//Name of the account holder, up to four lines possible.
	Name []string

	//AlternativeNames are account's alternative names.
	AlternativeNames []string

	//Classification of account. Can be one of 'Personal' or 'Business'.
	Classification Classification

	//IsJointAccount flag to indicate if the account is a joint account.
	IsJointAccount bool

	//IsAccountMatchingOptOut flag to indicate if the account has opted out of account matching.
	IsAccountMatchingOptOut bool

	//Additional information to identify the account and account holder.
	SecondaryIdentification string

	//IsSwitched flag to indicate if the account has been switched away from this organisation.
	IsSwitched bool
}

//NewAccount creates an instance of Account with default values assigned.
func NewAccount(name []string, country string, id, orgId uuid.UUID) *Account {
	return &Account{
		Id:                      id,
		Version:                 0,
		OrganisationId:          orgId,
		Country:                 country,
		BaseCurrency:            "",
		AccountNumber:           "",
		BankId:                  "",
		BankIdCode:              "",
		Bic:                     "",
		Iban:                    "",
		Name:                    name,
		AlternativeNames:        nil,
		Classification:          Personal,
		IsJointAccount:          false,
		IsAccountMatchingOptOut: false,
		SecondaryIdentification: "",
		IsSwitched:              false,
	}
}

//NewAccountFromDto creates a new instance of Account based upon
//an instance of AccountDto.
func NewAccountFromDto(dto data.AccountDto) *Account {
	name := dto.Data.Attributes.Name
	ctry := dto.Data.Attributes.Country
	id, oid := dto.Data.ID, dto.Data.OrganisationID

	uid, err := uuid.Parse(id)
	if err != nil {
		//TODO
	}

	ouid, err := uuid.Parse(oid)
	if err != nil {
		//TODO
	}

	return NewAccount(name, ctry, uid, ouid)
}


//ToDto transforms an instance of Account into a new instance of AccountDto
func (info *Account) ToDto() data.AccountDto {
	dto := data.NewAccountDto(info.Id, info.OrganisationId, info.Country, info.Bic)
	dto.Data.Version = info.Version
	dto.Data.Attributes.BaseCurrency = info.BaseCurrency
	dto.Data.Attributes.AccountNumber = info.AccountNumber
	dto.Data.Attributes.BankID = info.BankId
	dto.Data.Attributes.Bic = info.Bic
	dto.Data.Attributes.Iban = info.Iban
	dto.Data.Attributes.Name = info.Name[:]
	dto.Data.Attributes.AlternativeNames = info.AlternativeNames
	dto.Data.Attributes.AccountClassification = string(info.Classification)
	dto.Data.Attributes.JointAccount = info.IsJointAccount
	dto.Data.Attributes.AccountMatchingOptOut = info.IsAccountMatchingOptOut
	dto.Data.Attributes.SecondaryIdentification = info.SecondaryIdentification
	dto.Data.Attributes.Switched = info.IsSwitched
	return dto
}


