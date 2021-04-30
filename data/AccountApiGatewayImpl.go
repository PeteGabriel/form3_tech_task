package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
		apiUrl:    "http://0.0.0.0:8080/v1/organisation/accounts", //os.Getenv("ACCOUNT_API_ADDR"),
	}
}

//Create a new account
func (g *gateway) Create(dto AccountDto) (AccountDto, error) {
	cnt, err := json.Marshal(dto)
	if err != nil {
		err = fmt.Errorf("error converting structure to json format: %s", err)
		log.Print(err)
		return AccountDto{}, err
	}
	resp, err := g.webClient.Post(g.apiUrl, ContentType, bytes.NewBuffer(cnt))
	if err != nil {
		err = fmt.Errorf("error sending post request to account API: %s", err)
		log.Print(err)
		return AccountDto{}, err
	}
	defer resp.Body.Close()

	acc, err := unmarshalResponse(resp)
	if err != nil {
		log.Print(err)
		return acc, err
	}

	return acc, nil
}

//Delete an account by id
func (g *gateway) Delete(uid uuid.UUID) error {
	uri := fmt.Sprintf("%s/%s", g.apiUrl, uid.String())
	req, err := http.NewRequest(http.MethodDelete, uri, nil)
	if err != nil {
		log.Print(err)
		return err
	}

	//add 'version' param to request
	//TODO make it variable in function signature
	q := url.Values{}
	q.Add("version", "0")
	req.URL.RawQuery = q.Encode()

	resp, err := g.webClient.Do(req)

	if err != nil {
		err = fmt.Errorf("error sending delete request to account API: %s", err)
		log.Print(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		log.Printf("Deleting an account got a %d HTTP status code", resp.StatusCode)
		return fmt.Errorf("could not delete the account with uid %s", uid.String())
	}

	return nil
}

//Get an account by id
func (g *gateway) Get(uid uuid.UUID) (AccountDto, error) {
	resp, err := g.webClient.Get(fmt.Sprintf("%s/%s", g.apiUrl, uid.String()))

	if err != nil {
		err = fmt.Errorf("error sending get request to account API: %s", err)
		log.Print(err)
		return AccountDto{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Fetching an account got a %d HTTP status code", resp.StatusCode)
		return AccountDto{}, fmt.Errorf("could not fetch the account with uid %s", uid.String())
	}

	acc, err := unmarshalResponse(resp)
	if err != nil {
		err = fmt.Errorf("error converting json format to structure: %s", err)
		log.Print(err)
		return acc, err
	}

	return acc, nil
}

func unmarshalResponse(resp *http.Response) (AccountDto, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return AccountDto{}, fmt.Errorf("error reading body content: %s", err)
	}

	if resp.StatusCode == http.StatusBadRequest {
		accError := AccountError{}
		if err = json.Unmarshal(body, &accError); err != nil {
			return AccountDto{}, fmt.Errorf("error converting json format to structure: %s", err)
		}
		//grab error from api response
		return AccountDto{}, fmt.Errorf("bad request: %s", parseErrorMsg(accError.ErrorMsg))
	}

	//handle success case. Did not handled redirects as those seem not to be part of any
	//response and if implemented as HTTP spec it would need a specific handling.
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

Decided to separate ir by the \n char and send just the last part which seems more readable for the end user.
*/
func parseErrorMsg(msg string) string {
	strs := strings.Split(msg, "\n")
	return strs[len(strs)-1]
}