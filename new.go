package itswizard_m_berlinIdm

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const apiEndpoint = "schulportal.berlin.de/blusd/api/"
const apiEndpointDev = "idm-dev.lernraum-berlin.de/blusd/api/"
const auth = "authentication/authorization_header"
const users = "export_benutzer"
const user = "export_benutzer/"
const class = "export_klasse"
const course = "export_kurs"
const school = "export_schule"
const timeLayout = "2006-01-02 15:04:05.000"

type blusdError struct {
	ListError []struct {
		ErrorSource string `json:"errorSource"`
		ErrorText   string `json:"errorText"`
		Reference   string `json:"reference"`
	} `json:"listError"`
}

type BlusdConnection struct {
	clientId               string
	clientSecret           string
	authentificationHeader string
	url                    *url.URL
	dev                    bool
	client                 *http.Client
}

/*
Creates a new Connection to the lusd services
*/
func NewBlusdConnectionTestEnv(client string, secret string, dev bool) (blusd *BlusdConnection, err error) {
	blusd = new(BlusdConnection)
	blusd.clientId = client
	blusd.clientSecret = secret
	blusd.dev = dev
	if blusd.dev {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		blusd.client = &http.Client{Transport: tr}
	} else {
		blusd.client = &http.Client{}
	}
	err = blusd.authenficate()
	if err != nil {
		return
	}
	return
}

/*
Return the AuthorizartionToken which is in the
*/
func (p *BlusdConnection) GetAuthorizationToken() (string, error) {
	if p.authentificationHeader == "" {
		return "", errors.New("No token in system. Use ReAuthentificate()")
	} else {
		return p.authentificationHeader, nil
	}

}

/*
Authenticate to LUSD
*/
func (p *BlusdConnection) authenficate() error {
	login := make(map[string]string)
	login["client_id"] = p.clientId
	login["client_secret"] = p.clientSecret

	url, err := createUrl(login, auth, p.dev)
	p.url = url
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("GET", p.url.String(), nil)
	req.Header.Add("Accept", "application/json")
	resp, err := p.client.Do(req)
	log.Println(err)
	resp_body, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(resp_body), "listError") {
		var errorResp blusdError
		err = json.Unmarshal(resp_body, &errorResp)
		if err != nil {
			return errors.New(err.Error() + " " + string(resp_body))
		}
		var errorString string
		for _, k := range errorResp.ListError {
			errorString = errorString + " --- " + k.ErrorText + " : " + k.ErrorSource + " : " + " : " + k.Reference
		}
		p.authentificationHeader = ""
		err = resp.Body.Close()
		if err != nil {
			return err
		}
		return errors.New(errorString)
	}
	p.authentificationHeader = string(resp_body)

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

/*
ReAuthentificate the system
*/
func (p *BlusdConnection) ReAuthenficate() error {
	return p.authenficate()
}
