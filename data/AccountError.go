package data

import "fmt"

type AccountError struct {
	ErrorMsg  string `json:"error_message"`
	ErrorCode string `json:"error_code"`
}

//String represents a more readable syntax for the error coming from the account api.
func (err AccountError) String() string {
	return fmt.Sprintf("error '%s' with code '%s'", err.ErrorMsg, err.ErrorCode)
}
