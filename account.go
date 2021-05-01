package form3_task

import (
	"github.com/google/uuid"
	"github.com/petegabriel/form3_task/data"
	"log"
)

//Classification of account. Can be one of 'Personal' or 'Business'.
type Classification string

const (
	Personal Classification = "Personal"
	Business Classification = "Business"
)

//Account represents a bank account that is registered with Form3 fake account api.
type Account struct {

	//Id of the account in UUID 4 format. It identifies the resource.
	Id uuid.UUID

	//CreatedOn represents the time and date on which the resource was created.
	CreatedOn string

	//ModifiedOn represents the time and date on which the resource was last modified.
	ModifiedOn string

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
	//all other fields will be automatically initialised by their
	//respective 'zero value'.
	return &Account{
		Id:                      id,
		OrganisationId:          orgId,
		Country:                 country,
		Name:                    name,
		Classification:          Personal,
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
		log.Fatal(err)
	}

	ouid, err := uuid.Parse(oid)
	if err != nil {
		log.Fatal(err)
	}

	acc := NewAccount(name, ctry, uid, ouid)
	acc.CreatedOn = dto.Data.CreatedOn
	acc.ModifiedOn = dto.Data.ModifiedOn
	acc.Version = dto.Data.Version
	acc.BaseCurrency = dto.Data.Attributes.BaseCurrency
	acc.AccountNumber = dto.Data.Attributes.AccountNumber
	acc.BankId = dto.Data.Attributes.BankID
	acc.BankIdCode = dto.Data.Attributes.BankIDCode
	acc.Bic = dto.Data.Attributes.Bic
	acc.Iban = dto.Data.Attributes.Iban
	acc.Name = dto.Data.Attributes.Name
	acc.AlternativeNames = dto.Data.Attributes.AlternativeNames
	acc.Classification = Classification(dto.Data.Attributes.AccountClassification)
	acc.IsJointAccount = dto.Data.Attributes.JointAccount
	acc.IsAccountMatchingOptOut = dto.Data.Attributes.AccountMatchingOptOut
	acc.SecondaryIdentification = dto.Data.Attributes.SecondaryIdentification
	acc.IsSwitched = dto.Data.Attributes.Switched
	return acc

}

//ToDto transforms an instance of Account into a new instance of AccountDto
func (info *Account) ToDto() data.AccountDto {
	dto := data.NewAccountDto(info.Id, info.OrganisationId, info.Country, info.Bic, info.Name)
	dto.Data.CreatedOn = info.CreatedOn
	dto.Data.ModifiedOn = info.ModifiedOn
	dto.Data.Version = info.Version
	dto.Data.Attributes.BaseCurrency = info.BaseCurrency
	dto.Data.Attributes.AccountNumber = info.AccountNumber
	dto.Data.Attributes.BankID = info.BankId
	dto.Data.Attributes.BankIDCode = info.BankIdCode
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


