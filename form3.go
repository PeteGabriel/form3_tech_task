package form3_task

import "github.com/petegabriel/form3_task/data"

func Tech3() bool {
	return true
}


func CreateAccount(info *Account) (*Account, error){
	gate := data.NewGateway()
	dto := info.ToDto()
	acc, err := gate.Create(dto)
	if err != nil {
		return nil, err
	}
	return NewAccountFromDto(acc), nil
}

