package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//ContentType used to make http requests to the account api.
const ContentType = "application/vnd.api+json"

//gateway represents the access point to fetch/modify data in account api.
type gateway struct {
	webClient http.Client
	apiUrl    string
}

//NewGateway creates a new instance of gateway which implements the contract
//specified by AccountApiGateway interface.
func NewGateway() AccountApiGateway {
	return &gateway{
		webClient: http.Client{},
		apiUrl: os.Getenv("ACCOUNT_API_ADDR"),
	}
}

//Create a new account
func (g *gateway) Create(dto AccountDto) (AccountDto, error) {
	cnt, err := json.Marshal(dto)
	if err != nil {
		return AccountDto{}, fmt.Errorf("error converting structure to json format: %s", err)
	}
	resp, err := g.webClient.Post(g.apiUrl, ContentType, bytes.NewBuffer(cnt))
	if err != nil {
		return AccountDto{}, fmt.Errorf("error sending post request to account API: %s", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated:
		return handleGoodResult(resp.Body)
	case http.StatusBadRequest, http.StatusConflict:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return AccountDto{}, fmt.Errorf("error reading body content: %s", err)
		}
		accError := AccountError{}
		if err = json.Unmarshal(body, &accError); err != nil {
			return AccountDto{}, fmt.Errorf("error converting json format to structure: %s", err)
		}
		//grab error from api response
		return AccountDto{}, errors.New(parseErrorMsg(accError.ErrorMsg))
	default:
		return AccountDto{}, fmt.Errorf("error creating account - code %d", resp.StatusCode)
	}
}

//Delete an account by id and version
func (g *gateway) Delete(uid uuid.UUID, vrs string) error {
	uri := fmt.Sprintf("%s/%s", g.apiUrl, uid.String())
	req, err := http.NewRequest(http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	//add 'version' param to query string
	q := url.Values{
		"version": []string{vrs},
	}
	req.URL.RawQuery = q.Encode()

	resp, err := g.webClient.Do(req)

	if err != nil {
		return fmt.Errorf("error sending delete request to account API: %s", err)
	}
	defer resp.Body.Close()


	switch resp.StatusCode {
	case http.StatusNotFound:
		return fmt.Errorf("account with uuid %s not found", uid.String())
	case http.StatusConflict:
		return errors.New("account with specified version not found")
	default:
		return nil
	}
}

//Get an account by id
func (g *gateway) Get(uid uuid.UUID) (AccountDto, error) {
	resp, err := g.webClient.Get(fmt.Sprintf("%s/%s", g.apiUrl, uid.String()))

	if err != nil {
		return AccountDto{}, fmt.Errorf("error sending get request to account API: %s", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode{
	case http.StatusNotFound:
		return AccountDto{}, fmt.Errorf("account with uid %s not found", uid.String())
	case http.StatusOK:
		return handleGoodResult(resp.Body)
	default:
		return AccountDto{}, fmt.Errorf("error getting account with uid %s - code %d", uid.String(), resp.StatusCode)
	}
}

func handleGoodResult(src io.ReadCloser) (AccountDto, error) {
	body, err := ioutil.ReadAll(src)
	if err != nil {
		return AccountDto{}, fmt.Errorf("error reading body content: %s", err)
	}
	acc := AccountDto{}
	if err = json.Unmarshal(body, &acc); err != nil {
		return acc, fmt.Errorf("error converting json format to structure: %s", err)
	}
	return acc, nil
}

/*
Error messages seem to come with different levels of
context, for example:

  `validation failure list:\nvalidation failure list:\nvalidation failure list:\nname.1 in body should be at least 1 chars long`

Decided to separate it by the char '\n' and send just the last part which seems more readable for the end user.
*/
func parseErrorMsg(msg string) string {
	msg = strings.Trim(msg, " \n \t")
	strs := strings.Split(msg, "\n")
	return strs[len(strs)-1]
}