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
	"os"
)

var accountApiUri = os.Getenv("ACCOUNT_API_ADDR")
const ContentType   = "application/vnd.api+json"



type Gateway struct {
	webClient http.Client
}

//Create a new account
func (g *Gateway) Create(dto AccountDto) (AccountDto, error)  {
	cnt, err := json.Marshal(dto)
	if err != nil {
		err = fmt.Errorf("error converting structure to json format: %s", err)
		log.Print(err)
		return AccountDto{}, err
	}
	resp, err := g.webClient.Post(accountApiUri, ContentType, bytes.NewBuffer(cnt))
	if err != nil {
		err = fmt.Errorf("error sending post request to account API: %s", err)
		log.Print(err)
		return AccountDto{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			err = fmt.Errorf("error parsing body: %s", err)
			log.Print(err)
			return AccountDto{}, err
		}
		err = fmt.Errorf("error creating new account: %s", msg)
		log.Print(err)
		return AccountDto{}, err
	}

	acc, err := unmarshalFromBody(resp)
	if err != nil {
		err = fmt.Errorf("error converting json format to structure: %s", err)
		log.Print(err)
		return acc, err
	}

	return acc, nil
}

//Delete an account by id
func (g *Gateway) Delete(uid uuid.UUID) error  {
	uri := fmt.Sprintf("%s/%s", accountApiUri, uid.String())
	req, err := http.NewRequest(http.MethodDelete, uri, nil)
	if err != nil {
		log.Print(err)
		return err
	}

	//add 'version' param to request
	q := url.Values{}
	q.Add("version", "0")
	req.URL.RawQuery = q.Encode()

	resp, err := g.webClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("error sending delete request to account API: %s", err)
		log.Print(err)
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		log.Printf("Deleting an account got a %d HTTP status code", resp.StatusCode)
		return fmt.Errorf("could not delete the account with uid %s", uid.String())
	}

	return nil
}

//Get an account by id
func (g *Gateway) Get(uid uuid.UUID) (AccountDto, error)  {
	resp, err := g.webClient.Get(fmt.Sprintf("%s/%s", accountApiUri, uid.String()))
	defer resp.Body.Close()

	if err != nil {
		err = fmt.Errorf("error sending get request to account API: %s", err)
		log.Print(err)
		return AccountDto{}, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Fetching an account got a %d HTTP status code", resp.StatusCode)
		return AccountDto{}, fmt.Errorf("could not fetch the account with uid %s", uid.String())
	}

	acc, err := unmarshalFromBody(resp)
	if err != nil {
		err = fmt.Errorf("error converting json format to structure: %s", err)
		log.Print(err)
		return acc, err
	}

	return acc, nil
}

func unmarshalFromBody(resp *http.Response) (AccountDto, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("error reading body content: %s", err)
		log.Print(err)
		return AccountDto{}, err
	}

	acc := AccountDto{}
	err = json.Unmarshal(body, &acc)
	if err != nil {
		err = fmt.Errorf("error converting json format to structure: %s", err)
		log.Print(err)
		return acc, err
	}
	return acc, nil
}
