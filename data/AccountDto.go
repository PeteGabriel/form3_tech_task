package data

import "github.com/google/uuid"

type AccountDto struct {
	Data  Data  `json:"data"`
}
type Attributes struct {
	Country                 string   `json:"country"`
	BaseCurrency            string   `json:"base_currency"`
	AccountNumber           string   `json:"account_number"`
	BankID                  string   `json:"bank_id"`
	BankIDCode              string   `json:"bank_id_code"`
	Bic                     string   `json:"bic"`
	Iban                    string   `json:"iban"`
	Name                    []string `json:"name"`
	AlternativeNames        []string `json:"alternative_names"`
	AccountClassification   string   `json:"account_classification"`
	JointAccount            bool     `json:"joint_account"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out"`
	SecondaryIdentification string   `json:"secondary_identification"`
	Switched                bool     `json:"switched"`
}
type Data struct {
	Type           string     `json:"type"`
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Version        int        `json:"version"`
	Attributes     Attributes `json:"attributes"`
	CreatedOn      string     `json:"created_on"`
	ModifiedOn      string    `json:"modified_on"`
}


//NewAccountDto return a new account dto
func NewAccountDto(id, orgId uuid.UUID, cty string, name []string) AccountDto {
	return AccountDto{
		Data:  Data{
			Type:           "accounts",
			ID:             id.String(),
			OrganisationID: orgId.String(),
			Version:        0,
			Attributes:     Attributes{
				Country:                 cty,
				BaseCurrency:            "",
				AccountNumber:           "",
				BankID:                  "",
				BankIDCode:              "",
				Bic:                     "",
				Iban:                    "",
				Name:                    name,
				AlternativeNames:        nil,
				AccountClassification:   "Personal",
				JointAccount:            false,
				AccountMatchingOptOut:   false,
				SecondaryIdentification: "",
				Switched:                false,
			},
		},
	}
}