package core

import (
	"bytes"
	"encoding/json"
	"github.com/juju/errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	GET  = "GET"
	POST = "POST"
	PUT  = "PUT"
)

type (
	Api struct {
		BaseUrl string
		Token   string
		client  *http.Client
	}
)

func NewClient(baseUrl string, token string) *Api {
	return &Api{
		BaseUrl: baseUrl,
		Token:   token,
		client:  &http.Client{},
	}
}

func (a *Api) CallAPI(method string, endpoint string, body string) (string, error) {

	var b io.Reader
	if len(body) > 0 {
		b = bytes.NewReader([]byte(body))
	}
	request, _ := http.NewRequest(method, a.BaseUrl+endpoint, b)
	request.Header.Set("Authorization", "Bearer "+a.Token)

	resp, err := a.client.Do(request)
	if err != nil {
		return "", err
	}
	if resp.StatusCode > 206 {
		return "", errors.New(resp.Status)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	return string(bodyText), nil
}

func (a *Api) BasicAuth(username string, passwd string) error {

	req, err := http.NewRequest("POST", a.BaseUrl+"/auth/token", nil)
	req.SetBasicAuth(username, passwd)
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode > 206 {
		return errors.New(resp.Status)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s, _ := strconv.Unquote(string(bodyText))

	a.Token = s
	return nil
}

func (a *Api) GetSinks(agentId string) (anyJson []string, err error) {

	sinkList, err := a.CallAPI(GET, "/collectors/"+agentId+"/sinks", "")
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(sinkList), &anyJson)
	if err != nil {
		return
	}

	return

}

func (a *Api) GetSink(agentId string, sinkName string) (anyJson map[string]interface{}, err error) {

	sink, err := a.CallAPI(GET, "/collectors/"+agentId+"/sinks/"+sinkName, "")
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(sink), &anyJson)
	if err != nil {
		return
	}
	return

}
